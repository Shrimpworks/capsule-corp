import { createHash } from "node:crypto";
import {
  predecodeStrictJobProposalJson,
  type StrictJsonValue,
} from "./strict-job-proposal-json.js";

export const C2B_PASSIVE_BINDING_V3_KNOWN_ANSWER = Object.freeze({
  bytes: 18_357,
  sha256: "d72327bba369484a56db7d543a32e8bbd4eac403230ac65d63709ac3ba3bbdfb",
  contractSha256: "8b1ec936a7b56370716d28557125e46866dea8f21a149704a01f251a0dddbcc1",
  semanticBytes: 18_523,
  semanticSha256: "e550212a9e10f2ecc30cb862fe4c228e0ddbcbb0dc807b14b154e3eead7a89d0",
  c1Bytes: 9_289,
  c1Sha256: "d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e",
  c2aBytes: 26_850,
  c2aSha256: "d4ce88888186266f5d251e6246c889b1fd46d7746bb0ba56bcc4b3ce4675992f",
  c2bV1Bytes: 8_221,
  c2bV1Sha256: "3540d5224bdc81edbceafa1f0f17ac119904a70feab604957ab349dd116961a6",
  c2bV2Bytes: 7_115,
  c2bV2Sha256: "c59f7fdd27834dd7be5a05a3c44a973d6ffa99869b9b99e2531045926827190a",
} as const);

export interface C2BPassiveBindingV3 {
  readonly objectType: "capsule.governed-deno-core-c2b-passive-binding";
  readonly schemaVersion: 3;
  readonly domain: "capsule.governed-deno-core.c2b.passive-binding/v3";
  readonly identity: "capsule.governed-deno-core-c2b-passive-binding/fixed-owned-guest-passive-successor-v3";
  readonly status: "PASSED-passive-successor-contract-only";
  readonly predecessors: Readonly<Record<string, StrictJsonValue>>;
  readonly governedInputs: Readonly<Record<string, StrictJsonValue>>;
  readonly artifactRoles: Readonly<Record<string, StrictJsonValue>>;
  readonly hostRunnerContract: Readonly<Record<string, StrictJsonValue>>;
  readonly descriptorAndDeviceContract: Readonly<Record<string, StrictJsonValue>>;
  readonly runtimeContract: Readonly<Record<string, StrictJsonValue>>;
  readonly resourceContract: Readonly<Record<string, StrictJsonValue>>;
  readonly transportContract: Readonly<Record<string, StrictJsonValue>>;
  readonly teardownContract: Readonly<Record<string, StrictJsonValue>>;
  readonly composedProfile: Readonly<Record<string, StrictJsonValue>>;
  readonly workStatus: Readonly<Record<string, StrictJsonValue>>;
  readonly nextGuestGate: Readonly<Record<string, StrictJsonValue>>;
  readonly effects: Readonly<Record<string, false>>;
}

export function decodeC2BPassiveBindingV3(exact: Uint8Array): C2BPassiveBindingV3 {
  checkExact(exact);
  const result = predecodeStrictJobProposalJson(exact);
  if (!result.ok) throw new Error(`C2B_V3_BINDING_${result.refusal.code}`);
  validateC2BPassiveBindingV3(result.value);
  return deepFreeze(structuredClone(result.value)) as unknown as C2BPassiveBindingV3;
}

export function validateC2BPassiveBindingV3(value: StrictJsonValue): void {
  if (containsNull(value)) throw new Error("C2B_V3_AUTHORITY_NULL_FORBIDDEN");
  if (!isObject(value)) throw new Error("C2B_V3_SEMANTIC_KNOWN_ANSWER_SHAPE");
  const encoded = new TextEncoder().encode(`${JSON.stringify(value, null, 2)}\n`);
  if (encoded.length !== C2B_PASSIVE_BINDING_V3_KNOWN_ANSWER.semanticBytes) {
    throw new Error("C2B_V3_SEMANTIC_KNOWN_ANSWER_LENGTH");
  }
  if (sha256(encoded) !== C2B_PASSIVE_BINDING_V3_KNOWN_ANSWER.semanticSha256) {
    throw new Error("C2B_V3_SEMANTIC_KNOWN_ANSWER_DIGEST");
  }
  if (
    !isObject(value.composedProfile) ||
    value.composedProfile.contractDigestSha256 !==
      C2B_PASSIVE_BINDING_V3_KNOWN_ANSWER.contractSha256
  ) {
    throw new Error("C2B_V3_COMPOSED_PROFILE_DIGEST");
  }
}

function checkExact(value: Uint8Array): void {
  if (value.length > C2B_PASSIVE_BINDING_V3_KNOWN_ANSWER.bytes) {
    throw new Error("C2B_V3_BINDING_CAP");
  }
  if (value.length !== C2B_PASSIVE_BINDING_V3_KNOWN_ANSWER.bytes) {
    throw new Error("C2B_V3_BINDING_LENGTH");
  }
  if (sha256(value) !== C2B_PASSIVE_BINDING_V3_KNOWN_ANSWER.sha256) {
    throw new Error("C2B_V3_BINDING_DIGEST");
  }
  const source = new TextDecoder().decode(value);
  const zeroForm = new TextEncoder().encode(
    source.replace(C2B_PASSIVE_BINDING_V3_KNOWN_ANSWER.contractSha256, "0".repeat(64)),
  );
  if (
    new TextDecoder().decode(zeroForm) === source ||
    sha256(zeroForm) !== C2B_PASSIVE_BINDING_V3_KNOWN_ANSWER.contractSha256
  ) {
    throw new Error("C2B_V3_COMPOSED_PROFILE_DIGEST");
  }
}

function isObject(
  value: StrictJsonValue | undefined,
): value is Readonly<Record<string, StrictJsonValue>> {
  return (
    value !== undefined && value !== null && typeof value === "object" && !Array.isArray(value)
  );
}

function containsNull(value: StrictJsonValue): boolean {
  if (value === null) return true;
  if (Array.isArray(value)) return value.some(containsNull);
  if (typeof value === "object") return Object.values(value).some(containsNull);
  return false;
}

function sha256(value: Uint8Array): string {
  return createHash("sha256").update(value).digest("hex");
}

function deepFreeze<T>(value: T): T {
  if (value !== null && typeof value === "object") {
    for (const child of Object.values(value as Record<string, unknown>)) deepFreeze(child);
    Object.freeze(value);
  }
  return value;
}
