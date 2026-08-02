import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import test from "node:test";

import {
  constructExecutionPlan,
  createTrustedExecutionPlanBindings,
  createTrustedExecutionPlanDigest,
  createTrustedExecutionPlanInstallationId,
  createTrustedProfileRegistryContext,
  createTrustedUserPolicyContext,
  decodeJobProposal,
  type ExecutionPlanRetainedByteRole,
  type ResolvedJobProposalPlanInputs,
  resolveJobProposal,
  type TrustedExecutionPlanBindings,
  type TrustedExecutionPlanDigest,
  type TrustedExecutionPlanDigestRole,
} from "./index.js";

const corpusRoot = new URL("../../../schemas/conformance/v0/", import.meta.url);
const expectedPlanDigest = "627f9524479000dab6f3cee1d70c0428c63285bcadbc2cb3c6e8018b2dea008c";

test("Task 3C emits the retained preferred-CBOR ExecutionPlan known answer", async () => {
  const planInputs = await ordinaryPlanInputs();
  const result = constructExecutionPlan(planInputs, ordinaryTrustedBindings());
  assert.equal(result.ok, true);
  if (!result.ok) {
    return;
  }

  const expected = new Uint8Array(
    await readFile(new URL("execution-plan/ordinary.cbor", corpusRoot)),
  );
  assert.equal(result.plan.exactBytes.byteLength, 530);
  assert.deepEqual(result.plan.copyExactBytes(), expected);
  assert.equal(Buffer.from(result.plan.digest.bytes).toString("hex"), expectedPlanDigest);
  assert.deepEqual(result.plan.candidate, {
    objectType: "capsule.execution-plan",
    objectVersion: 0,
    installationId: retained("installation-id", repeatedBytes(0x11, 16)),
    epochSequence: 7,
    epochDigest: retained("trust-epoch", repeatedBytes(0x22, 32)),
    sourceManifestDigest: retained(
      "source-manifest",
      hexBytes("e5e09b2435baedf897526a89c698c0b0531437a69472372ae426f62d801fc171"),
    ),
    sourceEntrypoint: "src/main.ts",
    sourceByteLength: 91,
    inputSlot: "primary-data",
    inlineInputDigest: retained(
      "inline-input",
      hexBytes("bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72"),
    ),
    inlineInputByteLength: 118,
    runtimeProfileAlias: "fixture-active@1",
    runtimeBundleManifestDigest: retained("runtime-bundle-manifest", repeatedBytes(0x55, 32)),
    profileReviewAttestationDigests: [
      retained("profile-review-attestation", repeatedBytes(0x66, 32)),
      retained("profile-review-attestation", repeatedBytes(0x67, 32)),
    ],
    profileRegistryEntryDigest: retained("profile-registry-entry", repeatedBytes(0x77, 32)),
    backendValidationRecordDigest: retained("backend-validation-record", repeatedBytes(0x88, 32)),
    backendConfigurationDigest: retained("backend-configuration", repeatedBytes(0x99, 32)),
    trustSnapshotDigest: retained("trust-snapshot", repeatedBytes(0xaa, 32)),
    policyDecisionDigest: retained("policy-decision", repeatedBytes(0xbb, 32)),
    wallTimeMs: 5_000,
    wallTimeOrigin: "requested",
    outputSlot: "transformed-json",
    outputMaxJsonBytes: 65_536,
    expiresAt: 1_785_456_300,
  });
});

test("builder retains deeply immutable values and returns only defensive byte copies", async () => {
  const inputBytes = repeatedBytes(0x11, 16);
  const installationId = createTrustedExecutionPlanInstallationId(inputBytes);
  inputBytes[0] = 0xff;
  const bindings = ordinaryTrustedBindings({ installationId });
  const result = constructExecutionPlan(await ordinaryPlanInputs(), bindings);
  assert.equal(result.ok, true);
  if (!result.ok) {
    return;
  }

  assert.equal(Object.isFrozen(bindings), true);
  assert.equal(Object.isFrozen(bindings.profileReviewAttestationDigests), true);
  assert.equal(Object.isFrozen(result.plan), true);
  assert.equal(Object.isFrozen(result.plan.candidate), true);
  assert.equal(Object.isFrozen(result.plan.candidate.installationId.bytes), true);
  assert.equal(Object.isFrozen(result.plan.candidate.profileReviewAttestationDigests), true);
  assert.equal(Object.isFrozen(result.plan.exactBytes.bytes), true);
  assert.equal(Object.isFrozen(result.plan.digest.bytes), true);
  assert.equal(result.plan.candidate.installationId.bytes[0], 0x11);

  const exactCopy = result.plan.copyExactBytes();
  exactCopy[0] = (exactCopy[0] ?? 0) ^ 0xff;
  const digestCopy = result.plan.copyDigestBytes();
  digestCopy[0] = (digestCopy[0] ?? 0) ^ 0xff;
  assert.equal(result.plan.copyExactBytes()[0], 0xb8);
  assert.equal(Buffer.from(result.plan.copyDigestBytes()).toString("hex"), expectedPlanDigest);
  assert.equal(Reflect.set(result.plan.exactBytes.bytes, "0", 0), false);
  assert.equal(Reflect.set(result.plan.digest.bytes, "0", 0), false);
});

test("builder refuses copied or generic values with fixed provenance data", async () => {
  const planInputs = await ordinaryPlanInputs();
  const bindings = ordinaryTrustedBindings();
  assert.deepEqual(
    constructExecutionPlan({ ...planInputs } as ResolvedJobProposalPlanInputs, bindings),
    {
      ok: false,
      refusal: {
        owner: "execution-plan-builder",
        classification: "BINDING",
        code: "PLAN_INPUT_PROVENANCE",
      },
    },
  );
  assert.deepEqual(
    constructExecutionPlan(planInputs, { ...bindings } as TrustedExecutionPlanBindings),
    {
      ok: false,
      refusal: {
        owner: "execution-plan-builder",
        classification: "BINDING",
        code: "TRUSTED_BINDING_PROVENANCE",
      },
    },
  );

  const decoded = decodeJobProposal(
    new Uint8Array(await readFile(new URL("job-proposal/ordinary.json", corpusRoot))),
  );
  assert.equal(decoded.ok, true);
  if (decoded.ok) {
    assert.deepEqual(
      constructExecutionPlan(
        decoded.proposal as unknown as ResolvedJobProposalPlanInputs,
        bindings,
      ),
      {
        ok: false,
        refusal: {
          owner: "execution-plan-builder",
          classification: "BINDING",
          code: "PLAN_INPUT_PROVENANCE",
        },
      },
    );
  }
});

test("trusted binding construction enforces closed roles, widths, counts, and UInt53 bounds", () => {
  assert.throws(() => createTrustedExecutionPlanInstallationId(new Uint8Array(16)), /nonzero/u);
  assert.throws(
    () => createTrustedExecutionPlanDigest(new Uint8Array(31), "trust-epoch"),
    /32 bytes/u,
  );
  assert.doesNotThrow(() =>
    createTrustedExecutionPlanDigest(new Uint8Array(32), "policy-decision"),
  );
  const ordinary = ordinaryTrustedBindingsInput();
  assert.throws(
    () =>
      createTrustedExecutionPlanBindings({
        ...ordinary,
        runtimeBundleManifestDigest:
          ordinary.profileRegistryEntryDigest as unknown as TrustedExecutionPlanDigest<"runtime-bundle-manifest">,
      }),
    /runtime-bundle-manifest digest provenance/u,
  );
  assert.throws(
    () => createTrustedExecutionPlanBindings({ ...ordinary, profileReviewAttestationDigests: [] }),
    /one through eight/u,
  );
  assert.throws(
    () =>
      createTrustedExecutionPlanBindings({
        ...ordinary,
        profileReviewAttestationDigests: Array.from({ length: 9 }, () =>
          trustedDigest(0x66, "profile-review-attestation"),
        ),
      }),
    /one through eight/u,
  );
  assert.throws(
    () => createTrustedExecutionPlanBindings({ ...ordinary, epochSequence: -1 }),
    /unsigned safe integer/u,
  );
  assert.throws(
    () =>
      createTrustedExecutionPlanBindings({
        ...ordinary,
        expiresAt: Number.MAX_SAFE_INTEGER + 1,
      }),
    /unsigned safe integer/u,
  );
  assert.throws(
    () =>
      createTrustedExecutionPlanBindings({
        ...ordinary,
        unexpected: true,
      } as typeof ordinary),
    /closed binding shape/u,
  );
});

test("preferred CBOR encoding is deterministic across the complete UInt53 range", async () => {
  const maximum = Number.MAX_SAFE_INTEGER;
  const bindings = ordinaryTrustedBindings({ epochSequence: maximum, expiresAt: maximum });
  const planInputs = await ordinaryPlanInputs();
  const first = constructExecutionPlan(planInputs, bindings);
  const second = constructExecutionPlan(planInputs, bindings);
  assert.equal(first.ok, true);
  assert.equal(second.ok, true);
  if (!first.ok || !second.ok) {
    return;
  }
  assert.deepEqual(first.plan.copyExactBytes(), second.plan.copyExactBytes());
  const maximumEncoding = Buffer.from("1b001fffffffffffff", "hex");
  const bytes = Buffer.from(first.plan.copyExactBytes());
  assert.notEqual(bytes.indexOf(Buffer.concat([Buffer.from([0x04]), maximumEncoding])), -1);
  assert.notEqual(bytes.indexOf(Buffer.concat([Buffer.from([0x18, 0x18]), maximumEncoding])), -1);
});

test("builder preserves the Task 3B trusted-default origin without rewriting exact values", async () => {
  const planInputs = await resolvedPlanInputs("job-proposal/wall-time-defaulted.json");
  const first = constructExecutionPlan(planInputs, ordinaryTrustedBindings());
  const second = constructExecutionPlan(planInputs, ordinaryTrustedBindings());
  assert.equal(first.ok, true);
  assert.equal(second.ok, true);
  if (!first.ok || !second.ok) {
    return;
  }
  assert.equal(first.plan.candidate.wallTimeMs, 5_000);
  assert.equal(first.plan.candidate.wallTimeOrigin, "trusted-default");
  assert.equal(first.plan.candidate.outputMaxJsonBytes, planInputs.outputMaxJsonBytes);
  assert.deepEqual(first.plan.copyExactBytes(), second.plan.copyExactBytes());
});

async function ordinaryPlanInputs(): Promise<ResolvedJobProposalPlanInputs> {
  return resolvedPlanInputs("job-proposal/ordinary.json");
}

async function resolvedPlanInputs(path: string): Promise<ResolvedJobProposalPlanInputs> {
  const decoded = decodeJobProposal(new Uint8Array(await readFile(new URL(path, corpusRoot))));
  assert.equal(decoded.ok, true);
  if (!decoded.ok) {
    throw new Error("ordinary proposal did not cross Task 3A");
  }
  const result = resolveJobProposal(
    decoded.proposal,
    createTrustedProfileRegistryContext([
      {
        alias: "fixture-active@1",
        status: "active",
        exactWallTimeMs: { minimum: 1, maximum: 10_000 },
      },
    ]),
    createTrustedUserPolicyContext({
      wallTimeMs: { trustedDefault: 5_000, ceiling: 10_000 },
    }),
  );
  assert.equal(result.ok, true);
  if (!result.ok) {
    throw new Error("ordinary proposal did not cross Task 3B");
  }
  return result.resolved.planInputs;
}

function ordinaryTrustedBindings(
  overrides: Partial<TrustedExecutionPlanBindingsInputForTest> = {},
): TrustedExecutionPlanBindings {
  return createTrustedExecutionPlanBindings({
    ...ordinaryTrustedBindingsInput(),
    ...overrides,
  });
}

type TrustedExecutionPlanBindingsInputForTest = ReturnType<typeof ordinaryTrustedBindingsInput>;

function ordinaryTrustedBindingsInput() {
  return {
    installationId: createTrustedExecutionPlanInstallationId(repeatedBytes(0x11, 16)),
    epochSequence: 7,
    epochDigest: trustedDigest(0x22, "trust-epoch"),
    runtimeBundleManifestDigest: trustedDigest(0x55, "runtime-bundle-manifest"),
    profileReviewAttestationDigests: [
      trustedDigest(0x66, "profile-review-attestation"),
      trustedDigest(0x67, "profile-review-attestation"),
    ],
    profileRegistryEntryDigest: trustedDigest(0x77, "profile-registry-entry"),
    backendValidationRecordDigest: trustedDigest(0x88, "backend-validation-record"),
    backendConfigurationDigest: trustedDigest(0x99, "backend-configuration"),
    trustSnapshotDigest: trustedDigest(0xaa, "trust-snapshot"),
    policyDecisionDigest: trustedDigest(0xbb, "policy-decision"),
    expiresAt: 1_785_456_300,
  };
}

function trustedDigest<Role extends TrustedExecutionPlanDigestRole>(
  value: number,
  role: Role,
): TrustedExecutionPlanDigest<Role> {
  return createTrustedExecutionPlanDigest(repeatedBytes(value, 32), role);
}

function retained<Role extends ExecutionPlanRetainedByteRole>(
  role: Role,
  bytes: Uint8Array,
): { readonly role: Role; readonly byteLength: number; readonly bytes: readonly number[] } {
  return {
    role,
    byteLength: bytes.byteLength,
    bytes: Array.from(bytes),
  };
}

function repeatedBytes(value: number, length: number): Uint8Array {
  return Uint8Array.from({ length }, () => value);
}

function hexBytes(value: string): Uint8Array {
  return Uint8Array.from(Buffer.from(value, "hex"));
}
