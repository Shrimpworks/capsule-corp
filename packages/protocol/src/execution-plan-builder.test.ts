import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import { readFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";

import {
  constructExecutionPlan,
  createTrustedExecutionPlanBindings,
  createTrustedExecutionPlanDigest,
  createTrustedExecutionPlanInstallationId,
  createTrustedExecutionPlanRegistrationRoleBindings,
  createTrustedProfileRegistryContext,
  createTrustedUserPolicyContext,
  decodeJobProposal,
  type ExecutionPlanRegistrationHandoff,
  type ExecutionPlanRetainedByteRole,
  prepareExecutionPlanRegistrationHandoff,
  type ResolvedJobProposalPlanInputs,
  resolveJobProposal,
  type TrustedExecutionPlanBindings,
  type TrustedExecutionPlanDigest,
  type TrustedExecutionPlanDigestRole,
  type TrustedExecutionPlanRegistrationRoleBindings,
} from "./index.js";

const corpusRoot = new URL("../../../schemas/conformance/v0/", import.meta.url);
const expectedPlanDigest = "ef268a0b829adc1ce1307203f4b805f63379954ccf41e8e20a7487b6e5acf241";

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
  assert.equal(result.plan.exactBytes.byteLength, 527);
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
      hexBytes("c387c80094027ffbcacb573f44f5f6b4dec4d243bb436b24dd644434feaa1d14"),
    ),
    sourceEntrypoint: "main.mjs",
    sourceByteLength: 50,
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

test("unwired registration handoff retains only defensive exact-byte and complete role copies", async () => {
  const constructed = constructExecutionPlan(await ordinaryPlanInputs(), ordinaryTrustedBindings());
  assert.equal(constructed.ok, true);
  if (!constructed.ok) {
    return;
  }
  const roleBindings = ordinaryRegistrationRoleBindings();
  const prepared = prepareExecutionPlanRegistrationHandoff(constructed.plan, roleBindings);
  assert.equal(prepared.ok, true);
  if (!prepared.ok) {
    return;
  }

  assert.equal(prepared.handoff.exactPlanByteLength, 527);
  assert.equal(Object.isFrozen(prepared.handoff), true);
  assert.equal(Object.isFrozen(prepared.handoff.roleBindings), true);
  assert.equal(
    Object.isFrozen(prepared.handoff.roleBindings.profileReviewAttestationDigests),
    true,
  );
  const callerCopy = prepared.handoff.copyExactPlanBytes();
  callerCopy[0] = (callerCopy[0] ?? 0) ^ 0xff;
  assert.equal(prepared.handoff.copyExactPlanBytes()[0], 0xb8);

  assert.deepEqual(
    prepareExecutionPlanRegistrationHandoff(
      { ...constructed.plan } as typeof constructed.plan,
      roleBindings,
    ),
    {
      ok: false,
      refusal: {
        owner: "execution-plan-registration-handoff",
        classification: "BINDING",
        code: "PLAN_PROVENANCE",
      },
    },
  );
  assert.deepEqual(
    prepareExecutionPlanRegistrationHandoff(constructed.plan, {
      ...roleBindings,
    } as TrustedExecutionPlanRegistrationRoleBindings),
    {
      ok: false,
      refusal: {
        owner: "execution-plan-registration-handoff",
        classification: "BINDING",
        code: "ROLE_BINDING_PROVENANCE",
      },
    },
  );
  assert.deepEqual(
    prepareExecutionPlanRegistrationHandoff(
      constructed.plan,
      ordinaryRegistrationRoleBindings({
        policyDecisionDigest: trustedDigest(0xbc, "policy-decision"),
      }),
    ),
    {
      ok: false,
      refusal: {
        owner: "execution-plan-registration-handoff",
        classification: "BINDING",
        code: "ROLE_BINDING_MISMATCH",
      },
    },
  );
});

test("registration handoff accepts the same profile review attestation set in a different order", async () => {
  const constructed = constructExecutionPlan(await ordinaryPlanInputs(), ordinaryTrustedBindings());
  assert.equal(constructed.ok, true);
  if (!constructed.ok) {
    return;
  }
  assert.deepEqual(
    constructed.plan.candidate.profileReviewAttestationDigests.map((digest) => digest.bytes[0]),
    [0x66, 0x67],
  );

  const reordered = ordinaryRegistrationRoleBindings({
    profileReviewAttestationDigests: [
      trustedDigest(0x67, "profile-review-attestation"),
      trustedDigest(0x66, "profile-review-attestation"),
    ],
  });
  const prepared = prepareExecutionPlanRegistrationHandoff(constructed.plan, reordered);
  assert.equal(prepared.ok, true);

  const missingOneAddedAnother = ordinaryRegistrationRoleBindings({
    profileReviewAttestationDigests: [
      trustedDigest(0x66, "profile-review-attestation"),
      trustedDigest(0x68, "profile-review-attestation"),
    ],
  });
  assert.deepEqual(
    prepareExecutionPlanRegistrationHandoff(constructed.plan, missingOneAddedAnother),
    {
      ok: false,
      refusal: {
        owner: "execution-plan-registration-handoff",
        classification: "BINDING",
        code: "ROLE_BINDING_MISMATCH",
      },
    },
  );
});

test("TypeScript exact-plan output reaches Go registrationstate through the local conformance harness", async () => {
  const constructed = constructExecutionPlan(await ordinaryPlanInputs(), ordinaryTrustedBindings());
  assert.equal(constructed.ok, true);
  if (!constructed.ok) {
    return;
  }
  const prepared = prepareExecutionPlanRegistrationHandoff(
    constructed.plan,
    ordinaryRegistrationRoleBindings(),
  );
  assert.equal(prepared.ok, true);
  if (!prepared.ok) {
    return;
  }

  const ordinary = serializeHandoff(prepared.handoff);
  const accepted = runGoRegistrationHandoff(ordinary);
  assert.deepEqual(accepted, {
    accepted: true,
    stateUnchanged: false,
    retainedExactPlanHex: Buffer.from(prepared.handoff.copyExactPlanBytes()).toString("hex"),
    recomputedPlanDigestHex: expectedPlanDigest,
    wireRegistrationPlanDigestHex: expectedPlanDigest,
  });

  const wrongRole = structuredClone(ordinary);
  wrongRole.roleBindings.policyDecisionDigestHex = "bc".repeat(32);
  assert.deepEqual(runGoRegistrationHandoff(wrongRole), {
    accepted: false,
    classification: "DOMAIN",
    stateUnchanged: true,
  });

  const missingRole = structuredClone(ordinary);
  delete missingRole.roleBindings.sourceManifestDigestHex;
  assert.deepEqual(runGoRegistrationHandoff(missingRole), {
    accepted: false,
    classification: "DOMAIN",
    stateUnchanged: true,
  });

  const mutatedPlan = structuredClone(ordinary);
  const mutatedBytes = Buffer.from(mutatedPlan.exactPlanHex, "hex");
  const policyDigestOffset = mutatedBytes.indexOf(Buffer.alloc(32, 0xbb));
  assert.notEqual(policyDigestOffset, -1);
  mutatedBytes[policyDigestOffset] = (mutatedBytes[policyDigestOffset] ?? 0) ^ 0x01;
  mutatedPlan.exactPlanHex = mutatedBytes.toString("hex");
  assert.deepEqual(runGoRegistrationHandoff(mutatedPlan), {
    accepted: false,
    classification: "DOMAIN",
    stateUnchanged: true,
  });
});

function ordinaryPlanInputs(): Promise<ResolvedJobProposalPlanInputs> {
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

type RegistrationRoleBindingOverrides = Partial<
  Parameters<typeof createTrustedExecutionPlanRegistrationRoleBindings>[0]
>;

function ordinaryRegistrationRoleBindings(
  overrides: RegistrationRoleBindingOverrides = {},
): TrustedExecutionPlanRegistrationRoleBindings {
  return createTrustedExecutionPlanRegistrationRoleBindings({
    installationId: createTrustedExecutionPlanInstallationId(repeatedBytes(0x11, 16)),
    epochDigest: trustedDigest(0x22, "trust-epoch"),
    sourceManifestDigest: createTrustedExecutionPlanDigest(
      hexBytes("c387c80094027ffbcacb573f44f5f6b4dec4d243bb436b24dd644434feaa1d14"),
      "source-manifest",
    ),
    inlineInputDigest: createTrustedExecutionPlanDigest(
      hexBytes("bd9968c72c34a6779dfe3259937a1d9a9e558036c7cd4895ef634fbf76181e72"),
      "inline-input",
    ),
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

interface SerializedRegistrationHandoff {
  exactPlanHex: string;
  roleBindings: {
    installationIdHex: string;
    epochDigestHex: string;
    sourceManifestDigestHex?: string;
    inlineInputDigestHex: string;
    runtimeBundleManifestDigestHex: string;
    profileReviewAttestationDigestHexes: string[];
    profileRegistryEntryDigestHex: string;
    backendValidationRecordDigestHex: string;
    backendConfigurationDigestHex: string;
    trustSnapshotDigestHex: string;
    policyDecisionDigestHex: string;
  };
}

interface GoRegistrationHandoffResponse {
  accepted: boolean;
  classification?: string;
  stateUnchanged: boolean;
  retainedExactPlanHex?: string;
  recomputedPlanDigestHex?: string;
  wireRegistrationPlanDigestHex?: string;
}

function serializeHandoff(
  handoff: ExecutionPlanRegistrationHandoff,
): SerializedRegistrationHandoff {
  const bindings = handoff.roleBindings;
  return {
    exactPlanHex: Buffer.from(handoff.copyExactPlanBytes()).toString("hex"),
    roleBindings: {
      installationIdHex: retainedHex(bindings.installationId),
      epochDigestHex: retainedHex(bindings.epochDigest),
      sourceManifestDigestHex: retainedHex(bindings.sourceManifestDigest),
      inlineInputDigestHex: retainedHex(bindings.inlineInputDigest),
      runtimeBundleManifestDigestHex: retainedHex(bindings.runtimeBundleManifestDigest),
      profileReviewAttestationDigestHexes:
        bindings.profileReviewAttestationDigests.map(retainedHex),
      profileRegistryEntryDigestHex: retainedHex(bindings.profileRegistryEntryDigest),
      backendValidationRecordDigestHex: retainedHex(bindings.backendValidationRecordDigest),
      backendConfigurationDigestHex: retainedHex(bindings.backendConfigurationDigest),
      trustSnapshotDigestHex: retainedHex(bindings.trustSnapshotDigest),
      policyDecisionDigestHex: retainedHex(bindings.policyDecisionDigest),
    },
  };
}

function retainedHex(value: { readonly bytes: readonly number[] }): string {
  return Buffer.from(value.bytes).toString("hex");
}

function runGoRegistrationHandoff(
  request: SerializedRegistrationHandoff,
): GoRegistrationHandoffResponse {
  const repositoryRoot = fileURLToPath(new URL("../../../", import.meta.url));
  const result = spawnSync(
    "go",
    ["run", "./internal/execution/registrationstate/conformancehandoff"],
    {
      cwd: repositoryRoot,
      encoding: "utf8",
      env: {
        ...process.env,
        GOCACHE: process.env.GOCACHE ?? join(tmpdir(), "capsule-registration-handoff-gocache"),
      },
      input: JSON.stringify(request),
      maxBuffer: 1 << 20,
    },
  );
  assert.equal(result.status, 0, result.stderr || result.stdout);
  return JSON.parse(result.stdout) as GoRegistrationHandoffResponse;
}
