/**
 * Task 3A's package-internal strict JSON predecoder. It exists separately from
 * the JobProposal schema projection so conformance tests can assert the first
 * owning raw boundary without treating a generic decoded value as a proposal.
 */

/** Hard structural ceilings enforced during predecode, before any schema/domain check runs. */
export const JOB_PROPOSAL_RAW_LIMITS = Object.freeze({
  rawBytes: 2_097_152,
  depth: 32,
  nodes: 8_193,
  totalMembers: 4_096,
  objectMembers: 256,
  totalElements: 4_096,
  arrayElements: 256,
  decodedTextBytes: 1_572_864,
  keyBytes: 128,
  stringBytes: 65_536,
} as const);

/** A generic JSON value produced by the raw predecoder, before any JobProposal-specific schema check. */
export type StrictJsonValue =
  | null
  | boolean
  | number
  | string
  | readonly StrictJsonValue[]
  | { readonly [key: string]: StrictJsonValue };

/** The exhaustive set of raw-boundary refusal codes {@link predecodeStrictJobProposalJson} can return. */
export type JobProposalRawRefusalCode =
  | "INPUT_TYPE"
  | "RAW_BYTES"
  | "EMPTY"
  | "UTF8_BOM"
  | "UTF8"
  | "SYNTAX"
  | "TRAILING_DATA"
  | "DUPLICATE_KEY"
  | "SURROGATE"
  | "INTEGER_GRAMMAR"
  | "INTEGER_RANGE"
  | "DEPTH"
  | "NODES"
  | "TOTAL_MEMBERS"
  | "OBJECT_MEMBERS"
  | "TOTAL_ELEMENTS"
  | "ARRAY_ELEMENTS"
  | "DECODED_TEXT"
  | "KEY_BYTES"
  | "STRING_BYTES";

export interface JobProposalRawRefusal {
  readonly owner: "public-raw-decoder";
  readonly classification: "MALFORMED";
  readonly code: JobProposalRawRefusalCode;
}

/** The outcome of {@link predecodeStrictJobProposalJson}: either a generic JSON value or a classified raw refusal. */
export type StrictJobProposalJsonResult =
  | { readonly ok: true; readonly value: StrictJsonValue }
  | { readonly ok: false; readonly refusal: JobProposalRawRefusal };

interface ParserMetrics {
  nodes: number;
  totalMembers: number;
  totalElements: number;
  decodedTextBytes: number;
}

class RawDecodeFailure extends Error {
  constructor(readonly code: JobProposalRawRefusalCode) {
    super(code);
  }
}

class StrictJsonParser {
  private index = 0;
  private readonly metrics: ParserMetrics = {
    nodes: 0,
    totalMembers: 0,
    totalElements: 0,
    decodedTextBytes: 0,
  };

  constructor(private readonly document: string) {}

  parse(): StrictJsonValue {
    this.skipWhitespace();
    if (this.index === this.document.length) {
      this.fail("EMPTY");
    }
    const value = this.parseValue(1, []);
    this.skipWhitespace();
    if (this.index !== this.document.length) {
      this.fail("TRAILING_DATA");
    }
    return value;
  }

  private parseValue(depth: number, path: readonly string[]): StrictJsonValue {
    if (depth > JOB_PROPOSAL_RAW_LIMITS.depth) {
      this.fail("DEPTH");
    }
    this.metrics.nodes += 1;
    if (this.metrics.nodes > JOB_PROPOSAL_RAW_LIMITS.nodes) {
      this.fail("NODES");
    }

    const character = this.document[this.index];
    if (character === "{") {
      return this.parseObject(depth, path);
    }
    if (character === "[") {
      return this.parseArray(depth, path);
    }
    if (character === '"') {
      return this.parseString(
        false,
        isSourceFileContentPath(path)
          ? JOB_PROPOSAL_RAW_LIMITS.decodedTextBytes
          : JOB_PROPOSAL_RAW_LIMITS.stringBytes,
      );
    }
    if (character === "t") {
      this.consumeLiteral("true");
      return true;
    }
    if (character === "f") {
      this.consumeLiteral("false");
      return false;
    }
    if (character === "n") {
      this.consumeLiteral("null");
      return null;
    }
    if (character === "-" || isAsciiDigit(character)) {
      return this.parseInteger();
    }
    this.fail("SYNTAX");
  }

  private parseObject(depth: number, path: readonly string[]): StrictJsonValue {
    this.index += 1;
    this.skipWhitespace();
    const value: Record<string, StrictJsonValue> = Object.create(null) as Record<
      string,
      StrictJsonValue
    >;
    const keys = new Set<string>();
    let members = 0;
    if (this.document[this.index] === "}") {
      this.index += 1;
      return value;
    }

    while (true) {
      members += 1;
      this.metrics.totalMembers += 1;
      if (members > JOB_PROPOSAL_RAW_LIMITS.objectMembers) {
        this.fail("OBJECT_MEMBERS");
      }
      if (this.metrics.totalMembers > JOB_PROPOSAL_RAW_LIMITS.totalMembers) {
        this.fail("TOTAL_MEMBERS");
      }
      if (this.document[this.index] !== '"') {
        this.fail("SYNTAX");
      }
      const key = this.parseString(
        true,
        isSourceFilesObjectPath(path)
          ? JOB_PROPOSAL_RAW_LIMITS.decodedTextBytes
          : JOB_PROPOSAL_RAW_LIMITS.keyBytes,
      );
      if (keys.has(key)) {
        this.fail("DUPLICATE_KEY");
      }
      keys.add(key);
      this.skipWhitespace();
      if (this.document[this.index] !== ":") {
        this.fail("SYNTAX");
      }
      this.index += 1;
      this.skipWhitespace();
      value[key] = this.parseValue(depth + 1, [...path, key]);
      this.skipWhitespace();
      const separator = this.document[this.index];
      if (separator === "}") {
        this.index += 1;
        return value;
      }
      if (separator !== ",") {
        this.fail("SYNTAX");
      }
      this.index += 1;
      this.skipWhitespace();
    }
  }

  private parseArray(depth: number, path: readonly string[]): StrictJsonValue {
    this.index += 1;
    this.skipWhitespace();
    const value: StrictJsonValue[] = [];
    if (this.document[this.index] === "]") {
      this.index += 1;
      return value;
    }

    while (true) {
      this.metrics.totalElements += 1;
      if (value.length >= JOB_PROPOSAL_RAW_LIMITS.arrayElements) {
        this.fail("ARRAY_ELEMENTS");
      }
      if (this.metrics.totalElements > JOB_PROPOSAL_RAW_LIMITS.totalElements) {
        this.fail("TOTAL_ELEMENTS");
      }
      value.push(this.parseValue(depth + 1, [...path, String(value.length)]));
      this.skipWhitespace();
      const separator = this.document[this.index];
      if (separator === "]") {
        this.index += 1;
        return value;
      }
      if (separator !== ",") {
        this.fail("SYNTAX");
      }
      this.index += 1;
      this.skipWhitespace();
    }
  }

  private parseString(isKey: boolean, maximumBytes: number): string {
    this.index += 1;
    let rawStart = this.index;
    let decodedBytes = 0;
    const pieces: string[] = [];

    while (this.index < this.document.length) {
      const codeUnit = this.document.charCodeAt(this.index);
      if (codeUnit === 0x22) {
        pieces.push(this.document.slice(rawStart, this.index));
        this.index += 1;
        this.metrics.decodedTextBytes += decodedBytes;
        return pieces.join("");
      }
      if (codeUnit === 0x5c) {
        pieces.push(this.document.slice(rawStart, this.index));
        this.index += 1;
        const escaped = this.parseEscape();
        pieces.push(escaped.value);
        decodedBytes += escaped.utf8Bytes;
        this.assertTextBudgets(decodedBytes, maximumBytes, isKey);
        rawStart = this.index;
        continue;
      }
      if (codeUnit < 0x20) {
        this.fail("SYNTAX");
      }
      if (isHighSurrogate(codeUnit)) {
        const low = this.document.charCodeAt(this.index + 1);
        if (!isLowSurrogate(low)) {
          this.fail("SURROGATE");
        }
        decodedBytes += 4;
        this.index += 2;
      } else {
        if (isLowSurrogate(codeUnit)) {
          this.fail("SURROGATE");
        }
        decodedBytes += utf8BytesForCodePoint(codeUnit);
        this.index += 1;
      }
      this.assertTextBudgets(decodedBytes, maximumBytes, isKey);
    }
    this.fail("SYNTAX");
  }

  private parseEscape(): { value: string; utf8Bytes: number } {
    const escapeCode = this.document[this.index];
    this.index += 1;
    const simple = SIMPLE_ESCAPES[escapeCode ?? ""];
    if (simple !== undefined) {
      return { value: simple, utf8Bytes: 1 };
    }
    if (escapeCode !== "u") {
      this.fail("SYNTAX");
    }
    const first = this.parseHexCodeUnit();
    if (isLowSurrogate(first)) {
      this.fail("SURROGATE");
    }
    if (!isHighSurrogate(first)) {
      return { value: String.fromCharCode(first), utf8Bytes: utf8BytesForCodePoint(first) };
    }
    if (this.document[this.index] !== "\\" || this.document[this.index + 1] !== "u") {
      this.fail("SURROGATE");
    }
    this.index += 2;
    const second = this.parseHexCodeUnit();
    if (!isLowSurrogate(second)) {
      this.fail("SURROGATE");
    }
    const codePoint = 0x1_0000 + ((first - 0xd800) << 10) + (second - 0xdc00);
    return { value: String.fromCodePoint(codePoint), utf8Bytes: 4 };
  }

  private parseHexCodeUnit(): number {
    if (this.index + 4 > this.document.length) {
      this.fail("SYNTAX");
    }
    let value = 0;
    for (let offset = 0; offset < 4; offset += 1) {
      const digit = hexDigit(this.document.charCodeAt(this.index + offset));
      if (digit < 0) {
        this.fail("SYNTAX");
      }
      value = value * 16 + digit;
    }
    this.index += 4;
    return value;
  }

  private parseInteger(): number {
    const start = this.index;
    const negative = this.document[this.index] === "-";
    if (negative) {
      this.index += 1;
    }
    const first = this.document[this.index];
    if (first === "0") {
      this.index += 1;
      if (negative || isAsciiDigit(this.document[this.index])) {
        this.fail("INTEGER_GRAMMAR");
      }
    } else if (first !== undefined && first >= "1" && first <= "9") {
      do {
        this.index += 1;
      } while (isAsciiDigit(this.document[this.index]));
    } else {
      this.fail("INTEGER_GRAMMAR");
    }
    if ([".", "e", "E"].includes(this.document[this.index] ?? "")) {
      this.fail("INTEGER_GRAMMAR");
    }
    const token = this.document.slice(start, this.index);
    const magnitude = negative ? token.slice(1) : token;
    if (magnitude.length > 16 || (magnitude.length === 16 && magnitude > "9007199254740991")) {
      this.fail("INTEGER_RANGE");
    }
    const value = Number(token);
    if (!Number.isSafeInteger(value)) {
      this.fail("INTEGER_RANGE");
    }
    return value;
  }

  private consumeLiteral(literal: string): void {
    if (this.document.slice(this.index, this.index + literal.length) !== literal) {
      this.fail("SYNTAX");
    }
    this.index += literal.length;
  }

  private assertTextBudgets(decodedBytes: number, maximumBytes: number, isKey: boolean): void {
    if (decodedBytes > maximumBytes) {
      this.fail(isKey ? "KEY_BYTES" : "STRING_BYTES");
    }
    if (this.metrics.decodedTextBytes + decodedBytes > JOB_PROPOSAL_RAW_LIMITS.decodedTextBytes) {
      this.fail("DECODED_TEXT");
    }
  }

  private skipWhitespace(): void {
    while (isJsonWhitespace(this.document.charCodeAt(this.index))) {
      this.index += 1;
    }
  }

  private fail(code: JobProposalRawRefusalCode): never {
    throw new RawDecodeFailure(code);
  }
}

const SIMPLE_ESCAPES: Readonly<Record<string, string>> = Object.freeze({
  '"': '"',
  "\\": "\\",
  "/": "/",
  b: "\b",
  f: "\f",
  n: "\n",
  r: "\r",
  t: "\t",
});

/** @internal Exact raw-boundary entry point used by the package decoder and corpus tests. */
export function predecodeStrictJobProposalJson(bytes: Uint8Array): StrictJobProposalJsonResult {
  if (!(bytes instanceof Uint8Array)) {
    return rejected("INPUT_TYPE");
  }
  if (bytes.length > JOB_PROPOSAL_RAW_LIMITS.rawBytes) {
    return rejected("RAW_BYTES");
  }
  if (bytes.length === 0) {
    return rejected("EMPTY");
  }

  const receivedBytes = Uint8Array.from(bytes);
  if (receivedBytes[0] === 0xef && receivedBytes[1] === 0xbb && receivedBytes[2] === 0xbf) {
    return rejected("UTF8_BOM");
  }

  let document: string;
  try {
    document = new TextDecoder("utf-8", { fatal: true }).decode(receivedBytes);
  } catch {
    return rejected("UTF8");
  }

  try {
    return Object.freeze({ ok: true, value: new StrictJsonParser(document).parse() });
  } catch (error) {
    if (error instanceof RawDecodeFailure) {
      return rejected(error.code);
    }
    throw error;
  }
}

function rejected(code: JobProposalRawRefusalCode): StrictJobProposalJsonResult {
  return Object.freeze({
    ok: false,
    refusal: Object.freeze({
      owner: "public-raw-decoder",
      classification: "MALFORMED",
      code,
    }),
  });
}

function isSourceFilesObjectPath(path: readonly string[]): boolean {
  return path.length === 2 && path[0] === "source" && path[1] === "files";
}

function isSourceFileContentPath(path: readonly string[]): boolean {
  return path.length === 3 && path[0] === "source" && path[1] === "files";
}

function isAsciiDigit(value: string | undefined): boolean {
  return value !== undefined && value >= "0" && value <= "9";
}

function isJsonWhitespace(codeUnit: number): boolean {
  return codeUnit === 0x20 || codeUnit === 0x09 || codeUnit === 0x0a || codeUnit === 0x0d;
}

function isHighSurrogate(codeUnit: number): boolean {
  return codeUnit >= 0xd800 && codeUnit <= 0xdbff;
}

function isLowSurrogate(codeUnit: number): boolean {
  return codeUnit >= 0xdc00 && codeUnit <= 0xdfff;
}

function utf8BytesForCodePoint(codePoint: number): number {
  if (codePoint <= 0x7f) {
    return 1;
  }
  if (codePoint <= 0x7ff) {
    return 2;
  }
  return 3;
}

function hexDigit(codeUnit: number): number {
  if (codeUnit >= 0x30 && codeUnit <= 0x39) {
    return codeUnit - 0x30;
  }
  if (codeUnit >= 0x41 && codeUnit <= 0x46) {
    return codeUnit - 0x41 + 10;
  }
  if (codeUnit >= 0x61 && codeUnit <= 0x66) {
    return codeUnit - 0x61 + 10;
  }
  return -1;
}
