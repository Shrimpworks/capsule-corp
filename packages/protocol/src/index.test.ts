import assert from "node:assert/strict";
import test from "node:test";

import {
  asPositiveSafeInteger,
  asRuntimeProfileAlias,
  asSourceEntrypoint,
  asSourcePath,
  createDenyByDefaultCapabilities,
  JOB_PROPOSAL_API_VERSION,
  JOB_PROPOSAL_KIND,
  type JobProposal,
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
