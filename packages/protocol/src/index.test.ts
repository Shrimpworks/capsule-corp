import assert from "node:assert/strict";
import test from "node:test";

import {
  asPositiveSafeInteger,
  asRuntimeProfileAlias,
  asSafeJsonInteger,
  createDenyByDefaultCapabilities,
  JOB_PROPOSAL_API_VERSION,
  JOB_PROPOSAL_KIND,
  type JobProposal,
  PRIMARY_DATA_INPUT_SLOT,
  TRANSFORMED_JSON_OUTPUT_SLOT,
} from "./index.js";
// These candidate helpers/types are intentionally not part of the public
// barrel (see internal-contract-candidates.ts's header): no product
// component consumes them yet, so this test imports the internal module
// directly rather than reopening it as public surface.
import {
  asCandidateDigest,
  asCandidateInstallationId,
  asCandidatePositiveUInt53,
  asCandidateRegistrationId,
  asCandidateSupervisorId,
  asCandidateUInt53,
  type PlanRegistrationCandidate,
} from "./internal-contract-candidates.js";

test("deny-by-default capabilities grant no ambient authority", () => {
  assert.deepEqual(createDenyByDefaultCapabilities(), {
    network: [],
    subprocesses: [],
    environment: [],
    nativeAddons: false,
    ffi: false,
    packageInstallation: false,
  });
});

test("passive JobProposal candidate exposes only fixed first-slice roles", () => {
  const proposal: JobProposal = {
    apiVersion: JOB_PROPOSAL_API_VERSION,
    kind: JOB_PROPOSAL_KIND,
    source: {
      entrypoint: "main.mjs",
      files: {
        "main.mjs": "console.log('{}');",
      },
    },
    runtimeProfile: asRuntimeProfileAlias("bun-data@1"),
    input: {
      slot: PRIMARY_DATA_INPUT_SLOT,
      kind: "inline-json",
      value: { message: "hello" },
    },
    requestedLimits: {
      wallTimeMs: asPositiveSafeInteger(5000),
    },
    outputs: [
      {
        slot: TRANSFORMED_JSON_OUTPUT_SLOT,
        kind: "inline-json",
        maxBytes: asPositiveSafeInteger(65536),
      },
    ],
  };

  assert.equal(proposal.kind, "JobProposal");
  assert.equal(proposal.input.slot, "primary-data");
  assert.equal(proposal.outputs[0].slot, "transformed-json");

  // "only" fixed first-slice roles: no extra top-level fields, no extra
  // input fields, and exactly one output carrying no extra fields either.
  assert.deepEqual(
    Object.keys(proposal).sort(),
    [
      "apiVersion",
      "kind",
      "source",
      "runtimeProfile",
      "input",
      "requestedLimits",
      "outputs",
    ].sort(),
  );
  assert.deepEqual(Object.keys(proposal.input).sort(), ["slot", "kind", "value"].sort());
  assert.equal(proposal.outputs.length, 1);
  assert.deepEqual(Object.keys(proposal.outputs[0]).sort(), ["slot", "kind", "maxBytes"].sort());
});

test("candidate scalar constructors reject ambiguous or unsafe values", () => {
  assert.throws(() => asRuntimeProfileAlias("Bun@latest"), /runtime profile alias/u);
  assert.throws(() => asPositiveSafeInteger(0), /positive safe integer/u);
  assert.throws(() => asPositiveSafeInteger(Number.MAX_SAFE_INTEGER + 1), /positive safe integer/u);
  assert.throws(() => asSafeJsonInteger(1.5), /safe JSON integer/u);
});

test("PlanRegistration candidate keeps identity domains distinct", () => {
  const registration: PlanRegistrationCandidate = {
    objectType: "capsule.plan-registration",
    objectVersion: 0,
    registrationId: asCandidateRegistrationId(Uint8Array.from({ length: 16 }, () => 0x77)),
    registrationSequence: asCandidatePositiveUInt53(1),
    planDigest: asCandidateDigest(
      Uint8Array.from({ length: 32 }, () => 0x88),
      "execution-plan",
    ),
    installationId: asCandidateInstallationId(Uint8Array.from({ length: 16 }, () => 0x11)),
    epochSequence: asCandidateUInt53(7),
    epochDigest: asCandidateDigest(
      Uint8Array.from({ length: 32 }, () => 0x22),
      "trust-epoch",
    ),
    supervisorId: asCandidateSupervisorId(Uint8Array.from({ length: 16 }, () => 0x55)),
    expiresAt: asCandidateUInt53(1_785_456_300),
  };

  assert.equal(registration.objectType, "capsule.plan-registration");
  assert.equal(registration.planDigest.length, 32);
  assert.notDeepEqual(registration.registrationId, registration.installationId);

  // Real cross-domain misuse: these candidate types carry no runtime role
  // tag (asCandidateDigest/copyCandidateId only check width — see
  // internal-contract-candidates.ts), so the *only* mechanism that keeps,
  // say, a trust-epoch digest out of an execution-plan digest slot is the
  // nominal branding TypeScript enforces at compile time. Exercise that
  // mechanism directly: assigning a wrong-domain value must fail to
  // compile. If a future change collapsed the brands (or dropped them),
  // these assignments would silently type-check and `@ts-expect-error`
  // would report an unused directive, failing `tsc`/`pnpm build`.
  // @ts-expect-error - a trust-epoch digest must not satisfy an execution-plan digest slot
  const misusedPlanDigest: PlanRegistrationCandidate["planDigest"] = registration.epochDigest;
  // @ts-expect-error - a registration ID must not satisfy an installation ID slot
  const misusedInstallationId: PlanRegistrationCandidate["installationId"] =
    registration.registrationId;
  // @ts-expect-error - a Supervisor ID must not satisfy a registration ID slot
  const misusedRegistrationId: PlanRegistrationCandidate["registrationId"] =
    registration.supervisorId;
  void misusedPlanDigest;
  void misusedInstallationId;
  void misusedRegistrationId;
});

test("candidate internal scalar constructors copy and bound inputs", () => {
  const source = Uint8Array.from({ length: 16 }, () => 0x11);
  const installationId = asCandidateInstallationId(source);
  source[0] = 0xff;

  assert.equal(installationId[0], 0x11);
  assert.throws(() => asCandidateRegistrationId(new Uint8Array(15)), /16 bytes/u);
  assert.throws(() => asCandidateSupervisorId(new Uint8Array(16)), /nonzero/u);
  assert.throws(() => asCandidateDigest(new Uint8Array(31), "execution-plan"), /32 bytes/u);
  assert.throws(() => asCandidateUInt53(-1), /unsigned safe integer/u);
});
