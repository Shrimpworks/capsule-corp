import { mkdir, readFile, writeFile } from "node:fs/promises";
import { fromHex, sha256, sha256Hex } from "./lib/fixture-bytes.mjs";

const root = new URL("../schemas/conformance/typescript-approved-byte-v0/", import.meta.url);
const check = process.argv.includes("--check");
const original = Buffer.from(
  "aW50ZXJmYWNlIENhcHN1bGVJbnB1dCB7CiAgcmVhZG9ubHkgdmFsdWVzOiByZWFkb25seSBudW1iZXJbXTsKICByZWFkb25seSBsYWJlbDogc3RyaW5nOwp9Cgpjb25zdCBpbnB1dDogQ2Fwc3VsZUlucHV0ID0gewogIHZhbHVlczogWzEsIDIsIDNdLAogIGxhYmVsOiAiY2Fwc3VsZS1vd25lZCIsCn07Cgpjb25zdCBvdXRwdXQ6IHsgcmVhZG9ubHkgc3VtOiBudW1iZXI7IHJlYWRvbmx5IGxhYmVsOiBzdHJpbmcgfSA9IHsKICBzdW06IGlucHV0LnZhbHVlcy5yZWR1Y2UoKHRvdGFsOiBudW1iZXIsIHZhbHVlOiBudW1iZXIpOiBudW1iZXIgPT4gdG90YWwgKyB2YWx1ZSwgMCksCiAgbGFiZWw6IGlucHV0LmxhYmVsLAp9OwoKZ2xvYmFsVGhpcy5fX2NhcHN1bGVSZXN1bHQgPSBvdXRwdXQ7Cg==",
  "base64",
);
const emitted = Buffer.from(
  "ICAgICAgICAgICAgICAgICAgICAgICAgCiAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAKICAgICAgICAgICAgICAgICAgICAgICAgIAogCgpjb25zdCBpbnB1dCAgICAgICAgICAgICAgID0gewogIHZhbHVlczogWzEsIDIsIDNdLAogIGxhYmVsOiAiY2Fwc3VsZS1vd25lZCIsCn07Cgpjb25zdCBvdXRwdXQgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICAgICA9IHsKICBzdW06IGlucHV0LnZhbHVlcy5yZWR1Y2UoKHRvdGFsICAgICAgICAsIHZhbHVlICAgICAgICApICAgICAgICAgPT4gdG90YWwgKyB2YWx1ZSwgMCksCiAgbGFiZWw6IGlucHV0LmxhYmVsLAp9OwoKZ2xvYmFsVGhpcy5fX2NhcHN1bGVSZXN1bHQgPSBvdXRwdXQ7Cg==",
  "base64",
);
const tsMedia = "application/capsule.typescript-source;v=0;module=esm";
const jsMedia = "application/capsule.javascript-source;v=0;module=esm";
const path = "ordinary.ts";
const toolchain = {
  source: "87104b07e7acee748bcc5391e1bc69cf3571caa0fdfb8b1d6b5fd3f9599b7849",
  distribution: "261da057fb25ff2912dd6abb7842fc915ddf7947a2cb3c8cce90875d2b9bb667",
  executable: "245e0321af97d3c21dd4e7104457334dfe3c3ba7982d0db75363e354565f8cbb",
};

const profile = map([
  [1, "capsule.typescript-transformer-profile"],
  [2, 0],
  [3, "node:module.stripTypeScriptTypes"],
  [4, "22.22.1"],
  [5, "1.1.5"],
  [6, fromHex(toolchain.source)],
  [7, "darwin"],
  [8, "arm64"],
  [9, fromHex(toolchain.distribution)],
  [10, fromHex(toolchain.executable)],
]);
const options = map([
  [1, "capsule.typescript-normalized-options"],
  [2, 0],
  [3, "strip"],
  [4, "reject-any"],
  [5, "absent"],
  [6, "absent"],
  [7, tsMedia],
  [8, jsMedia],
]);
const record = map([
  [1, "capsule.typescript-transformation-record"],
  [2, 0],
  [3, path],
  [4, tsMedia],
  [5, original.length],
  [6, sha256(original)],
  [7, jsMedia],
  [8, emitted.length],
  [9, sha256(emitted)],
  [10, sha256(profile)],
  [11, "22.22.1"],
  [12, "1.1.5"],
  [13, fromHex(toolchain.source)],
  [14, fromHex(toolchain.distribution)],
  [15, fromHex(toolchain.executable)],
  [16, options],
  [17, sha256(options)],
  [18, "absent"],
  [19, "absent"],
  [20, "reject-any"],
  [21, 0],
]);
const originalManifest = sourceManifest(
  "capsule.original-authoring-source-manifest",
  tsMedia,
  original,
);
const executableManifest = sourceManifest(
  "capsule.executable-javascript-source-manifest",
  jsMedia,
  emitted,
);
const recordSet = map([
  [1, "capsule.typescript-transformation-record-set"],
  [2, 0],
  [3, [record]],
]);
const planBindings = map([
  [1, "capsule.typescript-plan-source-bindings"],
  [2, 0],
  [3, sha256(originalManifest)],
  [4, sha256(executableManifest)],
  [5, sha256(recordSet)],
]);

const fixtures = new Map([
  ["files/ordinary.ts", original],
  ["files/ordinary.js", emitted],
  ["objects/transformer-profile.cbor", profile],
  ["objects/normalized-options.cbor", options],
  ["objects/original-manifest.cbor", originalManifest],
  ["objects/executable-manifest.cbor", executableManifest],
  ["objects/transformation-record.cbor", record],
  ["objects/transformation-record-set.cbor", recordSet],
  ["objects/plan-source-bindings.cbor", planBindings],
]);

const mediaTypes = {
  "objects/transformer-profile.cbor": "application/capsule.typescript-transformer-profile+cbor;v=0",
  "objects/normalized-options.cbor": "application/capsule.typescript-normalized-options+cbor;v=0",
  "objects/original-manifest.cbor":
    "application/capsule.original-authoring-source-manifest+cbor;v=0",
  "objects/executable-manifest.cbor":
    "application/capsule.executable-javascript-source-manifest+cbor;v=0",
  "objects/transformation-record.cbor":
    "application/capsule.typescript-transformation-record+cbor;v=0",
  "objects/transformation-record-set.cbor":
    "application/capsule.typescript-transformation-record-set+cbor;v=0",
  "objects/plan-source-bindings.cbor":
    "application/capsule.typescript-plan-source-bindings+cbor;v=0",
};
const knownAnswers = [...fixtures].map(([fixturePath, bytes]) => ({
  path: fixturePath,
  byteLength: bytes.length,
  sha256: sha256Hex(bytes),
  ...(mediaTypes[fixturePath] ? { mediaType: mediaTypes[fixturePath] } : {}),
}));
const mutations = [
  ["original-byte", "BYTE_BINDING"],
  ["emitted-byte", "BYTE_BINDING"],
  ["path-order", "ORDERING"],
  ["profile-digest", "TRANSFORMER_BINDING"],
  ["profile-node", "TRANSFORMER_BINDING"],
  ["options-bytes", "OPTIONS_BINDING"],
  ["options-digest", "OPTIONS_BINDING"],
  ["input-media-type", "MEDIA_TYPE"],
  ["output-media-type", "MEDIA_TYPE"],
  ["source-map-present", "DISPOSITION"],
  ["source-url-present", "DISPOSITION"],
  ["diagnostic-count-one", "DIAGNOSTIC"],
  ["cross-domain-original-as-executable", "DOMAIN"],
  ["cross-domain-record-set-as-original", "DOMAIN"],
].map(([id, classification]) => ({ id, decision: "reject", classification }));
const manifest = Buffer.from(
  `${JSON.stringify(
    {
      manifestVersion: "capsule.typescript-approved-byte-conformance/v0",
      status: "passive-unwired-fixture-only",
      fixtureCount: knownAnswers.length,
      mutationCount: mutations.length,
      caps: {
        sourceFiles: 32,
        originalFileBytes: 262144,
        originalAggregateBytes: 1048576,
        emittedFileBytes: 262144,
        emittedAggregateBytes: 1048576,
      },
      knownAnswers,
      mutations,
    },
    null,
    2,
  )}\n`,
);

for (const [fixturePath, bytes] of [...fixtures, ["manifest.json", manifest]]) {
  const destination = new URL(fixturePath, root);
  if (check) {
    let existing;
    try {
      existing = await readFile(destination);
    } catch {
      throw new Error(`missing generated approved-byte fixture: ${fixturePath}`);
    }
    if (!existing.equals(bytes)) {
      throw new Error(`stale generated approved-byte fixture: ${fixturePath}`);
    }
  } else {
    await mkdir(new URL("./", destination), { recursive: true });
    await writeFile(destination, bytes);
  }
}
process.stdout.write(
  `${check ? "verified" : "generated"} approved-byte corpus: ${knownAnswers.length} fixtures, ${mutations.length} mutations\n`,
);

function sourceManifest(objectType, mediaType, bytes) {
  return map([
    [1, objectType],
    [2, 0],
    [3, path],
    [4, [[path, mediaType, sha256(bytes), bytes.length]]],
    [5, bytes.length],
  ]);
}

function encode(value) {
  if (Buffer.isBuffer(value) || value instanceof Uint8Array) {
    return Buffer.concat([argument(2, value.length), Buffer.from(value)]);
  }
  if (typeof value === "number") return argument(0, value);
  if (typeof value === "string") {
    const bytes = Buffer.from(value);
    return Buffer.concat([argument(3, bytes.length), bytes]);
  }
  if (Array.isArray(value)) return Buffer.concat([argument(4, value.length), ...value.map(encode)]);
  throw new Error(`unsupported fixture value: ${String(value)}`);
}

function map(entries) {
  const encodedEntries = entries
    .map(([key, value]) => [argument(0, key), encode(value)])
    .sort(([left], [right]) => {
      if (left.length !== right.length) {
        return left.length - right.length;
      }
      return Buffer.compare(left, right);
    });
  return Buffer.concat([argument(5, encodedEntries.length), ...encodedEntries.flat()]);
}

function argument(major, value) {
  if (value < 24) return Buffer.from([(major << 5) | value]);
  if (value <= 0xff) return Buffer.from([(major << 5) | 24, value]);
  if (value <= 0xffff) return Buffer.from([(major << 5) | 25, value >> 8, value & 0xff]);
  if (value <= 0xffffffff) {
    return Buffer.from([
      (major << 5) | 26,
      Math.floor(value / 0x1000000) & 0xff,
      Math.floor(value / 0x10000) & 0xff,
      Math.floor(value / 0x100) & 0xff,
      value & 0xff,
    ]);
  }
  throw new Error("fixture integer exceeds uint32");
}

