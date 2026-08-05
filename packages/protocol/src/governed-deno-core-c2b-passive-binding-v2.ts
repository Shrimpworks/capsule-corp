import { createHash } from "node:crypto";
import {
  predecodeStrictJobProposalJson,
  type StrictJsonValue,
} from "./strict-job-proposal-json.js";

export const C2B_PASSIVE_BINDING_V2_KNOWN_ANSWER = Object.freeze({
  bytes: 7_115,
  sha256: "c59f7fdd27834dd7be5a05a3c44a973d6ffa99869b9b99e2531045926827190a",
  c1Bytes: 9_289,
  c1Sha256: "d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e",
  c2aBytes: 26_850,
  c2aSha256: "d4ce88888186266f5d251e6246c889b1fd46d7746bb0ba56bcc4b3ce4675992f",
  c2bV1Bytes: 8_221,
  c2bV1Sha256: "3540d5224bdc81edbceafa1f0f17ac119904a70feab604957ab349dd116961a6",
} as const);

export interface C2BPassiveBindingV2 {
  readonly objectType: "capsule.governed-deno-core-c2b-passive-binding";
  readonly schemaVersion: 2;
  readonly domain: "capsule.governed-deno-core.c2b.passive-binding/v2";
  readonly identity: "capsule.governed-deno-core-c2b-passive-binding/c1-c2a-no-guest-build-closure-v2";
  readonly status: "PASSED-passive-no-guest-build-closure-only";
  readonly predecessors: Readonly<Record<string, StrictJsonValue>>;
  readonly archiveEvidence: Readonly<Record<string, StrictJsonValue>>;
  readonly constructedArtifacts: readonly Readonly<Record<string, StrictJsonValue>>[];
  readonly runtimeManifestCandidate: Readonly<Record<string, StrictJsonValue>>;
  readonly hostPreflightHarness: Readonly<Record<string, StrictJsonValue>>;
  readonly unresolved: Readonly<Record<string, StrictJsonValue>>;
  readonly workStatus: Readonly<Record<string, StrictJsonValue>>;
  readonly nextConsumerGate: Readonly<Record<string, StrictJsonValue>>;
  readonly effects: Readonly<Record<string, false>>;
}

export function decodeC2BPassiveBindingV2(exact: Uint8Array): C2BPassiveBindingV2 {
  checkExact(exact);
  const result = predecodeStrictJobProposalJson(exact);
  if (!result.ok) throw new Error(`C2B_V2_BINDING_${result.refusal.code}`);
  validateC2BPassiveBindingV2(result.value);
  return deepFreeze(structuredClone(result.value)) as unknown as C2BPassiveBindingV2;
}

export function validateC2BPassiveBindingV2(value: StrictJsonValue): void {
  const encoded = new TextEncoder().encode(`${JSON.stringify(value, null, 2)}\n`);
  if (encoded.length !== C2B_PASSIVE_BINDING_V2_KNOWN_ANSWER.bytes) {
    throw new Error("C2B_V2_SEMANTIC_KNOWN_ANSWER_LENGTH");
  }
  if (sha256(encoded) !== C2B_PASSIVE_BINDING_V2_KNOWN_ANSWER.sha256) {
    throw new Error("C2B_V2_SEMANTIC_KNOWN_ANSWER_DIGEST");
  }
}

function checkExact(value: Uint8Array): void {
  if (value.length > C2B_PASSIVE_BINDING_V2_KNOWN_ANSWER.bytes) {
    throw new Error("C2B_V2_BINDING_CAP");
  }
  if (value.length !== C2B_PASSIVE_BINDING_V2_KNOWN_ANSWER.bytes) {
    throw new Error("C2B_V2_BINDING_LENGTH");
  }
  if (sha256(value) !== C2B_PASSIVE_BINDING_V2_KNOWN_ANSWER.sha256) {
    throw new Error("C2B_V2_BINDING_DIGEST");
  }
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
