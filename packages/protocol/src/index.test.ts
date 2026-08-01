import assert from "node:assert/strict";
import test from "node:test";

import {
  asCandidateDigest,
  asCandidateInstallationId,
  asCandidateRegistrationId,
  asCandidateSupervisorId,
  asCandidateUInt53,
  asPositiveSafeInteger,
  asRuntimeProfileAlias,
  asSourceEntrypoint,
  asSourcePath,
  createDenyByDefaultCapabilities,
  JOB_PROPOSAL_API_VERSION,
  JOB_PROPOSAL_KIND,
  type JobProposal,
  type PlanRegistrationCandidate,
  PRIMARY_DATA_INPUT_SLOT,
  TRANSFORMED_JSON_OUTPUT_SLOT,
} from "./index.js";

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
      entrypoint: asSourceEntrypoint("main.ts"),
      files: {
        [asSourcePath("main.ts")]: "console.log('{}');",
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
});

test("candidate scalar constructors reject ambiguous or unsafe values", () => {
  assert.throws(() => asSourcePath("/main.ts"), /source path/u);
  assert.throws(() => asSourcePath("lib/../main.ts"), /source path/u);
  assert.throws(() => asSourceEntrypoint("README.md"), /entrypoint/u);
  assert.throws(() => asRuntimeProfileAlias("Bun@latest"), /runtime profile alias/u);
  assert.throws(() => asPositiveSafeInteger(0), /positive safe integer/u);
  assert.throws(() => asPositiveSafeInteger(Number.MAX_SAFE_INTEGER + 1), /positive safe integer/u);
});

test("PlanRegistration candidate keeps identity domains distinct", () => {
  const registration: PlanRegistrationCandidate = {
    objectType: "capsule.plan-registration",
    objectVersion: 0,
    registrationId: asCandidateRegistrationId(Uint8Array.from({ length: 16 }, () => 0x77)),
    registrationSequence: asCandidateUInt53(1),
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
