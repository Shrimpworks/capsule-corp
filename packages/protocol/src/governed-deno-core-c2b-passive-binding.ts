import { createHash } from "node:crypto";
import {
  predecodeStrictJobProposalJson,
  type StrictJsonValue,
} from "./strict-job-proposal-json.js";

export const C2B_PASSIVE_BINDING_KNOWN_ANSWER = Object.freeze({
  bytes: 8_221,
  sha256: "3540d5224bdc81edbceafa1f0f17ac119904a70feab604957ab349dd116961a6",
  c1Bytes: 9_289,
  c1Sha256: "d5d75e638a15be6c9f4a3230d17309d085f6ec103a73b64d9e0fd656a5423c9e",
  c2aBytes: 26_850,
  c2aSha256: "d4ce88888186266f5d251e6246c889b1fd46d7746bb0ba56bcc4b3ce4675992f",
  forkBytes: 1_895,
  forkSha256: "41350bcfc854338ded5e62f77475daf86486351356104dbbf647a8f8b5f11946",
  evidenceBytes: 6_026,
  evidenceSha256: "7de51e88ccc7cf35ee04b822d028f6c2184b72eaa833352fe9fc656b4a7baa13",
  evidenceSelfDigest: "732301bf8553b0c59b3fe0e4f2b9e070dcc3a1b478e742dc13bd438873b7e488",
} as const);

export interface C2BPassiveBindingInputs {
  readonly c1: Uint8Array;
  readonly c2a: Uint8Array;
  readonly forkSupplement: Uint8Array;
  readonly buildEvidence: Uint8Array;
}

export interface C2BPassiveBinding {
  readonly objectType: "capsule.governed-deno-core-c2b-passive-binding";
  readonly schemaVersion: 1;
  readonly domain: "capsule.governed-deno-core.c2b.passive-binding/v1";
  readonly identity: "capsule.governed-deno-core-c2b-passive-binding/c1-c2a-v1";
  readonly status: "PASSED-passive-reconciliation-only";
  readonly predecessors: Readonly<Record<"c1" | "c2a", Readonly<Record<string, StrictJsonValue>>>>;
  readonly forkSupplement: Readonly<Record<string, StrictJsonValue>>;
  readonly buildEvidence: Readonly<Record<string, StrictJsonValue>>;
  readonly fixedFixture: Readonly<{
    source: Readonly<Record<string, StrictJsonValue>>;
    sourceManifest: Readonly<Record<string, StrictJsonValue>>;
    input: Readonly<Record<string, StrictJsonValue>>;
    completion: Readonly<Record<string, StrictJsonValue>>;
    authority: Readonly<Record<string, StrictJsonValue>>;
  }>;
  readonly artifacts: readonly Readonly<Record<string, StrictJsonValue>>[];
  readonly dependencyMergePolicy: Readonly<{
    state: "BLOCKED";
    dependencies: readonly Readonly<Record<string, StrictJsonValue>>[];
    changeRule: string;
    inPlaceMutation: "forbidden";
  }>;
  readonly workStatus: Readonly<Record<string, StrictJsonValue>>;
  readonly limitations: Readonly<Record<string, StrictJsonValue>>;
  readonly nextConsumerGate: Readonly<{
    state: "BLOCKED";
    consumer: string;
    requirements: readonly string[];
    implemented: false;
    admissionEffect: false;
  }>;
  readonly effects: Readonly<Record<string, false>>;
}

export function decodeC2BPassiveBinding(
  exact: Uint8Array,
  inputs: C2BPassiveBindingInputs,
): C2BPassiveBinding {
  checkExact(
    "binding",
    exact,
    C2B_PASSIVE_BINDING_KNOWN_ANSWER.bytes,
    C2B_PASSIVE_BINDING_KNOWN_ANSWER.sha256,
  );
  const value = decodeStrict("binding", exact);
  validateC2BPassiveBinding(value);
  validateInputs(value, inputs);
  return deepFreeze(structuredClone(value)) as unknown as C2BPassiveBinding;
}

export function validateC2BPassiveBinding(value: StrictJsonValue): void {
  const encoded = new TextEncoder().encode(`${JSON.stringify(value, null, 2)}\n`);
  checkExact(
    "semantic-known-answer",
    encoded,
    C2B_PASSIVE_BINDING_KNOWN_ANSWER.bytes,
    C2B_PASSIVE_BINDING_KNOWN_ANSWER.sha256,
  );
}

function validateInputs(binding: StrictJsonValue, inputs: C2BPassiveBindingInputs): void {
  const values = [
    [
      "c1",
      inputs.c1,
      C2B_PASSIVE_BINDING_KNOWN_ANSWER.c1Bytes,
      C2B_PASSIVE_BINDING_KNOWN_ANSWER.c1Sha256,
    ],
    [
      "c2a",
      inputs.c2a,
      C2B_PASSIVE_BINDING_KNOWN_ANSWER.c2aBytes,
      C2B_PASSIVE_BINDING_KNOWN_ANSWER.c2aSha256,
    ],
    [
      "fork",
      inputs.forkSupplement,
      C2B_PASSIVE_BINDING_KNOWN_ANSWER.forkBytes,
      C2B_PASSIVE_BINDING_KNOWN_ANSWER.forkSha256,
    ],
    [
      "evidence",
      inputs.buildEvidence,
      C2B_PASSIVE_BINDING_KNOWN_ANSWER.evidenceBytes,
      C2B_PASSIVE_BINDING_KNOWN_ANSWER.evidenceSha256,
    ],
  ] as const;
  for (const [name, value, bytes, digest] of values) {
    checkExact(name, value, bytes, digest);
  }

  const c1 = decodeStrict("c1", inputs.c1);
  const c2a = decodeStrict("c2a", inputs.c2a);
  const fork = decodeStrict("fork", inputs.forkSupplement);
  const evidence = decodeStrict("evidence", inputs.buildEvidence);
  validateEvidenceSelfDigest(inputs.buildEvidence);

  equal(at(binding, "predecessors", "c1", "bytes"), inputs.c1.length, "C1 length cross-link");
  equal(at(binding, "predecessors", "c1", "sha256"), sha256(inputs.c1), "C1 digest cross-link");
  equal(at(binding, "predecessors", "c2a", "bytes"), inputs.c2a.length, "C2A length cross-link");
  equal(at(binding, "predecessors", "c2a", "sha256"), sha256(inputs.c2a), "C2A digest cross-link");
  equal(
    at(c1, "contract"),
    at(binding, "predecessors", "c1", "identity"),
    "C1 identity cross-link",
  );
  equal(
    at(c2a, "contract"),
    at(binding, "predecessors", "c2a", "identity"),
    "C2A identity cross-link",
  );

  const knownAnswerLinks = [
    [at(c2a, "knownAnswer", "source", "utf8"), at(binding, "fixedFixture", "source", "utf8")],
    [at(c2a, "knownAnswer", "source", "bytes"), at(binding, "fixedFixture", "source", "bytes")],
    [at(c2a, "knownAnswer", "source", "sha256"), at(binding, "fixedFixture", "source", "sha256")],
    [
      at(c2a, "knownAnswer", "source", "sourceManifestBytes"),
      at(binding, "fixedFixture", "sourceManifest", "bytes"),
    ],
    [
      at(c2a, "knownAnswer", "source", "sourceManifestSha256"),
      at(binding, "fixedFixture", "sourceManifest", "sha256"),
    ],
    [
      at(c2a, "knownAnswer", "canonicalInput", "utf8"),
      at(binding, "fixedFixture", "input", "utf8"),
    ],
    [
      at(c2a, "knownAnswer", "canonicalInput", "bytes"),
      at(binding, "fixedFixture", "input", "bytes"),
    ],
    [
      at(c2a, "knownAnswer", "canonicalInput", "sha256"),
      at(binding, "fixedFixture", "input", "sha256"),
    ],
    [
      at(c2a, "knownAnswer", "expectedCompletion", "jsonUtf8"),
      at(binding, "fixedFixture", "completion", "utf8"),
    ],
    [
      at(c2a, "knownAnswer", "expectedCompletion", "jsonBytes"),
      at(binding, "fixedFixture", "completion", "bytes"),
    ],
    [
      at(c2a, "knownAnswer", "expectedCompletion", "jsonSha256"),
      at(binding, "fixedFixture", "completion", "sha256"),
    ],
  ] as const;
  for (const [left, right] of knownAnswerLinks) equal(left, right, "C2A known-answer cross-link");

  equal(at(fork, "identity"), at(binding, "forkSupplement", "identity"), "fork identity");
  for (const role of ["source", "input", "completion"] as const) {
    equal(at(fork, role, "bytes"), at(binding, "fixedFixture", role, "bytes"), `${role} length`);
    equal(at(fork, role, "sha256"), at(binding, "fixedFixture", role, "sha256"), `${role} digest`);
  }

  equal(at(evidence, "identity"), at(binding, "buildEvidence", "identity"), "evidence identity");
  equal(
    at(evidence, "selfDigest", "sha256"),
    at(binding, "buildEvidence", "selfDigest"),
    "evidence self-digest",
  );
  equal(
    at(evidence, "immutableBuildOnlySupplement", "identity"),
    at(binding, "forkSupplement", "identity"),
    "evidence supplement identity",
  );
  equal(
    at(evidence, "immutableBuildOnlySupplement", "sha256"),
    at(binding, "forkSupplement", "sha256"),
    "evidence supplement digest",
  );
  equal(
    at(evidence, "sources", "deno", "commit"),
    at(binding, "forkSupplement", "head"),
    "fork head",
  );
  equal(
    at(evidence, "sources", "deno", "tree"),
    at(binding, "forkSupplement", "tree"),
    "fork tree",
  );

  const artifactLinks = [
    ["binary", 1],
    ["snapshot", 2],
    ["twoFileBundle", 3],
    ["bundleManifest", 4],
  ] as const;
  for (const [evidenceName, index] of artifactLinks) {
    equal(
      at(evidence, "artifacts", evidenceName, "bytes"),
      at(binding, "artifacts", index, "bytes"),
      `${evidenceName} length`,
    );
    equal(
      at(evidence, "artifacts", evidenceName, "sha256"),
      at(binding, "artifacts", index, "sha256"),
      `${evidenceName} digest`,
    );
  }
}

function decodeStrict(name: string, exact: Uint8Array): StrictJsonValue {
  const result = predecodeStrictJobProposalJson(exact);
  if (!result.ok) {
    throw new Error(`C2B_${name.toUpperCase()}_${result.refusal.code}`);
  }
  return result.value;
}

function validateEvidenceSelfDigest(exact: Uint8Array): void {
  const text = new TextDecoder("utf-8", { fatal: true }).decode(exact);
  const needle = `"sha256": "${C2B_PASSIVE_BINDING_KNOWN_ANSWER.evidenceSelfDigest}"`;
  if (text.split(needle).length !== 2) throw new Error("C2B_EVIDENCE_SELF_DIGEST_FIELD");
  const normalized = new TextEncoder().encode(text.replace(needle, '"sha256": null'));
  if (sha256(normalized) !== C2B_PASSIVE_BINDING_KNOWN_ANSWER.evidenceSelfDigest) {
    throw new Error("C2B_EVIDENCE_SELF_DIGEST_MISMATCH");
  }
}

function checkExact(name: string, value: Uint8Array, bytes: number, digest: string): void {
  if (value.length !== bytes) throw new Error(`C2B_${name.toUpperCase()}_LENGTH`);
  if (sha256(value) !== digest) throw new Error(`C2B_${name.toUpperCase()}_DIGEST`);
}

function sha256(value: Uint8Array): string {
  return createHash("sha256").update(value).digest("hex");
}

function at(value: StrictJsonValue, ...path: readonly (string | number)[]): StrictJsonValue {
  let current = value;
  for (const part of path) {
    if (typeof part === "number") {
      if (!Array.isArray(current) || part < 0 || part >= current.length) {
        throw new Error(`C2B_CROSS_LINK_MISSING:${path.join("/")}`);
      }
      current = current[part] as StrictJsonValue;
      continue;
    }
    if (
      current === null ||
      typeof current !== "object" ||
      Array.isArray(current) ||
      !(part in current)
    ) {
      throw new Error(`C2B_CROSS_LINK_MISSING:${path.join("/")}`);
    }
    current = (current as { readonly [key: string]: StrictJsonValue })[part] as StrictJsonValue;
  }
  return current;
}

function equal(left: StrictJsonValue, right: StrictJsonValue, name: string): void {
  if (JSON.stringify(left) !== JSON.stringify(right)) throw new Error(`C2B_CROSS_LINK:${name}`);
}

function deepFreeze<T>(value: T): T {
  if (value !== null && typeof value === "object") {
    for (const child of Object.values(value as Record<string, unknown>)) deepFreeze(child);
    Object.freeze(value);
  }
  return value;
}
