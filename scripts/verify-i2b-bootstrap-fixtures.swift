#!/usr/bin/env swift
import CryptoKit
import Foundation

enum Failure: Error { case refused(String) }

indirect enum CBOR {
    case unsigned(UInt64), negative(Int64), bytes(Data), text(String), array([CBOR]), map([Int: CBOR]), bool(Bool), tag(UInt64, CBOR)
}

struct Parser {
    let bytes: [UInt8]
    let maxDepth: Int
    let maxItems: Int
    let maxMap: Int
    let maxArray: Int
    let allowTag18: Bool
    var offset = 0
    var items = 0

    mutating func parse() throws -> CBOR {
        let value = try item(depth: 1)
        guard offset == bytes.count else { throw Failure.refused("trailing-cbor") }
        return value
    }

    mutating func item(depth: Int) throws -> CBOR {
        guard depth <= maxDepth else { throw Failure.refused("depth") }
        items += 1
        guard items <= maxItems else { throw Failure.refused("items") }
        let (major, argument) = try head()
        switch major {
        case 0:
            guard argument <= 9_007_199_254_740_991 else { throw Failure.refused("unsafe-integer") }
            return .unsigned(argument)
        case 1:
            guard argument < 9_007_199_254_740_991 else { throw Failure.refused("unsafe-negative") }
            return .negative(-1 - Int64(argument))
        case 2:
            return .bytes(try take(Int(argument)))
        case 3:
            let data = try take(Int(argument))
            guard let string = String(data: data, encoding: .utf8) else { throw Failure.refused("utf8") }
            return .text(string)
        case 4:
            guard argument <= UInt64(maxArray) else { throw Failure.refused("array-cap") }
            return .array(try (0..<Int(argument)).map { _ in try item(depth: depth + 1) })
        case 5:
            guard argument <= UInt64(maxMap) else { throw Failure.refused("map-cap") }
            var result: [Int: CBOR] = [:]
            var previous: [UInt8]? = nil
            for _ in 0..<Int(argument) {
                let start = offset
                let keyValue = try item(depth: depth + 1)
                let keyBytes = Array(bytes[start..<offset])
                if let old = previous, deterministicCompare(old, keyBytes) >= 0 { throw Failure.refused("duplicate-or-order") }
                previous = keyBytes
                let key: Int
                switch keyValue {
                case .unsigned(let value): key = Int(value)
                case .negative(let value): key = Int(value)
                default: throw Failure.refused("map-key")
                }
                guard result[key] == nil else { throw Failure.refused("duplicate") }
                result[key] = try item(depth: depth + 1)
            }
            return .map(result)
        case 6:
            guard allowTag18, depth == 1, argument == 18 else { throw Failure.refused("tag") }
            return .tag(argument, try item(depth: depth + 1))
        case 7:
            guard argument == 20 || argument == 21 else { throw Failure.refused("simple") }
            return .bool(argument == 21)
        default: throw Failure.refused("major")
        }
    }

    mutating func head() throws -> (UInt8, UInt64) {
        guard offset < bytes.count else { throw Failure.refused("truncated") }
        let initial = bytes[offset]; offset += 1
        let major = initial >> 5, additional = initial & 31
        if additional < 24 { return (major, UInt64(additional)) }
        let width: Int
        switch additional { case 24: width = 1; case 25: width = 2; case 26: width = 4; case 27: width = 8; default: throw Failure.refused("indefinite") }
        let raw = try take(width)
        var value: UInt64 = 0
        for byte in raw { value = (value << 8) | UInt64(byte) }
        let minimum: UInt64 = width == 1 ? 24 : width == 2 ? 256 : width == 4 ? 65_536 : 4_294_967_296
        guard value >= minimum else { throw Failure.refused("nonpreferred") }
        return (major, value)
    }

    mutating func take(_ count: Int) throws -> Data {
        guard count >= 0, offset + count <= bytes.count else { throw Failure.refused("truncated") }
        defer { offset += count }
        return Data(bytes[offset..<(offset + count)])
    }
}

struct Envelope {
    let protected: Data
    let payload: Data
    let signature: Data
}

func parseEnvelope(_ data: Data, object: String) throws -> Envelope {
    let rawCap = object == "request" ? 4096 : 6144
    guard !data.isEmpty, data.count <= rawCap else { throw Failure.refused("envelope-raw-cap") }
    var parser = Parser(bytes: Array(data), maxDepth: 3, maxItems: 7, maxMap: 0, maxArray: 4, allowTag18: true)
    guard case .tag(18, .array(let values)) = try parser.parse(), values.count == 4,
          case .bytes(let protected) = values[0], case .map(let unprotected) = values[1], unprotected.isEmpty,
          case .bytes(let payload) = values[2], case .bytes(let signature) = values[3], signature.count == 64 else {
        throw Failure.refused("sign1-shape")
    }
    let payloadRaw = object == "request" ? 2048 : 4096
    let protectedRaw = 256
    let envelopeCalculated = object == "request" ? 1033 : 1698
    let payloadCalculated = object == "request" ? 861 : 1527
    let protectedCalculated = object == "request" ? 98 : 97
    guard protected.count <= protectedRaw, payload.count <= payloadRaw else { throw Failure.refused("nested-raw-cap") }
    guard data.count <= envelopeCalculated, payload.count <= payloadCalculated, protected.count <= protectedCalculated else { throw Failure.refused("calculated-cap") }
    return Envelope(protected: protected, payload: payload, signature: signature)
}

func parseMap(_ data: Data, object: String) throws -> [Int: CBOR] {
    let entries: Int
    switch object { case "request": entries = 34; case "record": entries = 69; case "protected": entries = 3; case "public-key": entries = 5; default: throw Failure.refused("unknown-map-profile") }
    var parser = Parser(bytes: Array(data), maxDepth: 2, maxItems: 1 + entries * 2, maxMap: entries, maxArray: 0, allowTag18: false)
    guard case .map(let map) = try parser.parse(), map.count == entries else { throw Failure.refused("closed-map") }
    return map
}

func verifyEnvelope(_ data: Data, object: String, ordinaryPayload: Data, ordinaryRequestEnvelope: Data, ordinaryRequestPayload: Data, selfExpected: Bool, trustedNow: UInt64?, replay: String) throws -> String {
    let envelope = try parseEnvelope(data, object: object)
    let protected = try parseMap(envelope.protected, object: "protected")
    guard protected.count == 3, int(protected[1]) == -7,
          text(protected[3]) == (object == "request" ? "application/capsule.supervisor-bootstrap-request+cbor;v=0" : "application/capsule.supervisor-bootstrap-record+cbor;v=0"),
          bytes(protected[4])?.count == 32 else { throw Failure.refused("protected-profile") }
    let payload = try parseMap(envelope.payload, object: object)
    if !selfExpected && envelope.payload != ordinaryPayload { throw Failure.refused("ordinary-binding") }

    let publicKeyBytes = bytes(payload[object == "request" ? 6 : 10])
    let payloadKeyID = bytes(payload[object == "request" ? 7 : 11])
    guard let publicKeyBytes, publicKeyBytes.count == 77, let payloadKeyID,
          Data(SHA256.hash(data: publicKeyBytes)) == payloadKeyID, payloadKeyID == bytes(protected[4]) else {
        throw Failure.refused("key-binding")
    }
    guard let installationID = bytes(payload[object == "request" ? 5 : 9]),
          let epoch = uint(payload[object == "request" ? 15 : 19]),
          let epochDigest = bytes(payload[object == "request" ? 16 : 20]),
          let authorization = bytes(payload[object == "request" ? 8 : 12]) else {
        throw Failure.refused("authorization-identity-shape")
    }
    var authorizationInput = Data("capsule.installation-root-key-authorization/v0\0".utf8)
    authorizationInput.append(payloadKeyID)
    authorizationInput.append(installationID)
    authorizationInput.append(uint64BE(epoch))
    authorizationInput.append(epochDigest)
    guard Data(SHA256.hash(data: authorizationInput)) == authorization else {
        throw Failure.refused("authorization-identity-binding")
    }
    let keyMap = try parseMap(publicKeyBytes, object: "public-key")
    guard int(keyMap[1]) == 2, int(keyMap[3]) == -7, int(keyMap[-1]) == 1,
          let x = bytes(keyMap[-2]), x.count == 32, let y = bytes(keyMap[-3]), y.count == 32 else { throw Failure.refused("public-key") }
    let x963 = Data([4]) + x + y
    let publicKey = try P256.Signing.PublicKey(x963Representation: x963)
    let signature = try P256.Signing.ECDSASignature(rawRepresentation: envelope.signature)
    guard publicKey.isValidSignature(signature, for: sigStructure(envelope.protected, envelope.payload)) else { throw Failure.refused("signature") }

    if object == "request" {
        guard text(payload[1]) == "capsule.supervisor-bootstrap-request", int(payload[2]) == 0,
              text(payload[3]) == "capsule.installation.bootstrap.request", text(payload[4]) == "capsule.execution-supervisor.bootstrap",
              bytes(payload[5])?.count == 16, bytes(payload[9])?.count == 16, bytes(payload[30])?.count == 32,
              int(payload[15]) == 1, bool(payload[28]) == false, bool(payload[29]) == false else { throw Failure.refused("request-profile") }
        let issued = uint(payload[31]), notBefore = uint(payload[32]), expires = uint(payload[33])
        guard let issued, let notBefore, let expires, issued <= notBefore, notBefore < expires, expires - issued <= 300 else { throw Failure.refused("time-window") }
        guard let now = trustedNow, now >= notBefore, now < expires else { throw Failure.refused("time-admission") }
        switch replay {
        case "fresh": return "admit-once"
        case "pending-exact": return "resume-exact"
        case "pending-other", "completed-exact": throw Failure.refused("request-replay")
        default: throw Failure.refused("request-replay-state")
        }
    }

    guard text(payload[1]) == "capsule.supervisor-bootstrap-record", int(payload[2]) == 0,
          text(payload[3]) == "capsule.installation.bootstrap.record", text(payload[4]) == "capsule.execution-supervisor",
          bytes(payload[5]) == Data(SHA256.hash(data: ordinaryRequestPayload)),
          bytes(payload[6]) == Data(SHA256.hash(data: ordinaryRequestEnvelope)),
          uint(payload[7]) == UInt64(ordinaryRequestEnvelope.count), int(payload[19]) == 1,
          bool(payload[53]) == false, bool(payload[54]) == false, bool(payload[55]) == false, bool(payload[56]) == false else {
        throw Failure.refused("record-profile-binding")
    }
    let boundRequest = try parseMap(ordinaryRequestPayload, object: "request")
    guard uint(payload[59]) == uint(boundRequest[31]), uint(payload[60]) == uint(boundRequest[32]),
          uint(payload[61]) == uint(boundRequest[33]), text(payload[62]) == text(boundRequest[34]),
          text(payload[63]) == text(boundRequest[1]), uint(payload[64]) == uint(boundRequest[2]),
          text(payload[65]) == text(boundRequest[3]), text(payload[66]) == text(boundRequest[4]) else {
        throw Failure.refused("record-request-stable-binding")
    }
    let finalized = uint(payload[52]), requestExpires = uint(payload[61])
    let issued = uint(payload[67]), notBefore = uint(payload[68]), expires = uint(payload[69])
    guard let finalized, let requestExpires, let issued, let notBefore, let expires,
          let now = trustedNow, finalized <= issued, issued <= notBefore, notBefore < expires,
          expires - issued <= 300, expires <= requestExpires, now >= notBefore, now < expires else {
        throw Failure.refused("record-time-window")
    }
    switch replay {
    case "fresh", "pending": return "commit-once"
    case "completed-exact": return "return-retained-envelope"
    case "completed-other": throw Failure.refused("record-replay")
    default: throw Failure.refused("record-replay-state")
    }
}

func sigStructure(_ protected: Data, _ payload: Data) -> Data {
    var result = Data([0x84, 0x6a]); result.append(Data("Signature1".utf8)); result.append(cborBytes(protected)); result.append(0x40); result.append(cborBytes(payload)); return result
}
func cborBytes(_ value: Data) -> Data {
    var result = Data(); let count = value.count
    if count < 24 { result.append(UInt8(0x40 | count)) }
    else if count <= 255 { result.append(contentsOf: [0x58, UInt8(count)]) }
    else { result.append(contentsOf: [0x59, UInt8(count >> 8), UInt8(count & 255)]) }
    result.append(value); return result
}
func uint64BE(_ value: UInt64) -> Data {
    var bigEndian = value.bigEndian
    return Data(bytes: &bigEndian, count: MemoryLayout<UInt64>.size)
}
func deterministicCompare(_ left: [UInt8], _ right: [UInt8]) -> Int {
    if left.count != right.count { return left.count < right.count ? -1 : 1 }
    for (a, b) in zip(left, right) where a != b { return a < b ? -1 : 1 }
    return 0
}
func int(_ value: CBOR?) -> Int64? { if case .unsigned(let v) = value { return Int64(v) }; if case .negative(let v) = value { return v }; return nil }
func uint(_ value: CBOR?) -> UInt64? { if case .unsigned(let v) = value { return v }; return nil }
func text(_ value: CBOR?) -> String? { if case .text(let v) = value { return v }; return nil }
func bytes(_ value: CBOR?) -> Data? { if case .bytes(let v) = value { return v }; return nil }
func bool(_ value: CBOR?) -> Bool? { if case .bool(let v) = value { return v }; return nil }

let root = URL(fileURLWithPath: FileManager.default.currentDirectoryPath).appendingPathComponent("schemas/conformance/i2b-bootstrap-v0")
let manifestData = try Data(contentsOf: root.appendingPathComponent("manifest.json"))
guard let manifest = try JSONSerialization.jsonObject(with: manifestData) as? [String: Any],
      manifest["manifestVersion"] as? String == "capsule.i2b-bootstrap-conformance/v0",
      let cases = manifest["cases"] as? [[String: Any]], cases.count == 71,
      let effects = manifest["effects"] as? [String: Int], effects.values.allSatisfy({ $0 == 0 }) else { throw Failure.refused("manifest") }
let ordinaryRequestEnvelope = try Data(contentsOf: root.appendingPathComponent("request/ordinary.cose"))
let ordinaryRequestPayload = try Data(contentsOf: root.appendingPathComponent("request/ordinary.payload.cbor"))
let maximumRequestEnvelope = try Data(contentsOf: root.appendingPathComponent("request/calculated-maximum.cose"))
let maximumRequestPayload = try Data(contentsOf: root.appendingPathComponent("request/calculated-maximum.payload.cbor"))
let ordinaryRecordPayload = try Data(contentsOf: root.appendingPathComponent("record/ordinary.payload.cbor"))
var accepted = 0, refused = 0
for testCase in cases {
    let id = testCase["id"] as! String, object = testCase["object"] as! String, fixture = testCase["fixture"] as! String
    let expected = testCase["expected"] as! String, replay = testCase["replay"] as! String, selfExpected = testCase["selfExpected"] as! Bool
    let trustedNow = (testCase["trustedNow"] as? NSNumber)?.uint64Value
    let requestVariant = testCase["requestVariant"] as? String
    let boundRequestEnvelope = requestVariant == "maximum" ? maximumRequestEnvelope : ordinaryRequestEnvelope
    let boundRequestPayload = requestVariant == "maximum" ? maximumRequestPayload : ordinaryRequestPayload
    let data = try Data(contentsOf: root.appendingPathComponent(fixture))
    do {
        let decision = try verifyEnvelope(data, object: object, ordinaryPayload: object == "request" ? ordinaryRequestPayload : ordinaryRecordPayload, ordinaryRequestEnvelope: boundRequestEnvelope, ordinaryRequestPayload: boundRequestPayload, selfExpected: selfExpected, trustedNow: trustedNow, replay: replay)
        guard expected == "ACCEPT", decision == testCase["decision"] as? String else { throw Failure.refused("unexpected-accept-\(id)") }
        accepted += 1
    } catch {
        guard expected == "REFUSE" else { throw Failure.refused("unexpected-refusal-\(id)-\(error)") }
        refused += 1
    }
}
print("verified independent Swift I2B bootstrap corpus: \(accepted) accepts, \(refused) refusals")
