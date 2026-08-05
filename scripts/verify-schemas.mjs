import { readdir, readFile } from "node:fs/promises";
import { join } from "node:path";
import Ajv2020 from "ajv/dist/2020.js";
import addFormats from "ajv-formats";

const schemaDirectory = new URL("../schemas/", import.meta.url);
const entries = (await readdir(schemaDirectory)).sort();
const schemaFiles = [
  ...entries.filter((entry) => entry.endsWith(".schema.json")),
  "candidates/job-proposal-v0.schema.json",
];
const schemas = new Map();

if (schemaFiles.length === 0) {
  throw new Error("No JSON Schema files were found");
}

for (const filename of schemaFiles) {
  const source = await readFile(new URL(filename, schemaDirectory), "utf8");
  const schema = JSON.parse(source);

  if (schema.$schema !== "https://json-schema.org/draft/2020-12/schema") {
    throw new Error(`${filename} does not declare JSON Schema draft 2020-12`);
  }

  if (typeof schema.$id !== "string" || schema.$id.length === 0) {
    throw new Error(`${filename} does not declare a stable $id`);
  }

  schemas.set(filename, schema);
  process.stdout.write(`validated ${join("schemas", filename)}\n`);
}

const ajv = new Ajv2020({ allErrors: true, strict: true });
addFormats(ajv);

for (const schema of schemas.values()) {
  ajv.addSchema(schema);
}

const fixtures = [
  ["job.schema.json", new URL("../examples/jobs/hello.job.json", import.meta.url)],
  [
    "candidates/job-proposal-v0.schema.json",
    new URL("../examples/jobs/hello.proposal.json", import.meta.url),
  ],
  ["profile.schema.json", new URL("../profiles/bun-data/profile.json", import.meta.url)],
  [
    "governed-deno-core-c1-composition.schema.json",
    new URL(
      "../schemas/conformance/c1-governed-deno-core/controlled-development-profile.json",
      import.meta.url,
    ),
  ],
  [
    "governed-deno-core-c2a-execution-profile.schema.json",
    new URL(
      "../schemas/conformance/c2a-governed-deno-core/passive-execution-profile.json",
      import.meta.url,
    ),
  ],
  [
    "governed-deno-core-c2b-passive-binding.schema.json",
    new URL("../schemas/conformance/c2b-governed-deno-core/passive-binding.json", import.meta.url),
  ],
  [
    "governed-deno-core-release-candidate.schema.json",
    new URL(
      "../schemas/conformance/governed-deno-core-release-candidate/candidate-manifest.json",
      import.meta.url,
    ),
  ],
];

for (const [schemaName, fixtureUrl] of fixtures) {
  const schema = schemas.get(schemaName);
  const validate = ajv.getSchema(schema.$id);
  const fixture = JSON.parse(await readFile(fixtureUrl, "utf8"));

  if (!validate(fixture)) {
    throw new Error(`${fixtureUrl.pathname} failed validation: ${ajv.errorsText(validate.errors)}`);
  }

  process.stdout.write(`validated ${fixtureUrl.pathname}\n`);
}

const proposalSchema = schemas.get("candidates/job-proposal-v0.schema.json");
const validateProposal = ajv.getSchema(proposalSchema.$id);
const validProposal = JSON.parse(
  await readFile(new URL("../examples/jobs/hello.proposal.json", import.meta.url), "utf8"),
);
const invalidProposals = [
  ["unknown object version", { ...validProposal, apiVersion: "capsule.dev/v1" }],
  ["unsupported capability union", { ...validProposal, capabilities: { network: [] } }],
  [
    "agent-selected runtime digest",
    {
      ...validProposal,
      runtimeProfile: { alias: "bun-data@1", digest: `sha256:${"0".repeat(64)}` },
    },
  ],
  [
    "agent-selected guest output path",
    {
      ...validProposal,
      outputs: [{ ...validProposal.outputs[0], guestPath: "/output/result.json" }],
    },
  ],
  [
    "absolute source path",
    { ...validProposal, source: { ...validProposal.source, entrypoint: "/main.ts" } },
  ],
  [
    "wrong input slot role",
    { ...validProposal, input: { ...validProposal.input, slot: "transformed-json" } },
  ],
  [
    "unsafe requested integer",
    { ...validProposal, requestedLimits: { wallTimeMs: 9_007_199_254_740_992 } },
  ],
  [
    "fractional inline JSON number",
    { ...validProposal, input: { ...validProposal.input, value: { ratio: 1.5 } } },
  ],
];

for (const [name, fixture] of invalidProposals) {
  if (validateProposal(fixture)) {
    throw new Error(`invalid JobProposal fixture was accepted: ${name}`);
  }
  process.stdout.write(`rejected invalid JobProposal: ${name}\n`);
}

const c2aSchema = schemas.get("governed-deno-core-c2a-execution-profile.schema.json");
const validateC2A = ajv.getSchema(c2aSchema.$id);
const validC2A = JSON.parse(
  await readFile(
    new URL(
      "../schemas/conformance/c2a-governed-deno-core/passive-execution-profile.json",
      import.meta.url,
    ),
    "utf8",
  ),
);
const invalidC2A = [
  [
    "invented final artifact digest",
    mutate(validC2A, (fixture) => {
      fixture.artifactClosure.required[0].sha256 = "0".repeat(64);
    }),
  ],
  [
    "missing runner descriptor",
    mutate(validC2A, (fixture) => {
      fixture.runnerDescriptorProfile.entries.pop();
    }),
  ],
  [
    "silent resource substitution",
    mutate(validC2A, (fixture) => {
      fixture.machineProfile.guestRamMiB = 512;
    }),
  ],
  [
    "invented runtime-profile identity",
    mutate(validC2A, (fixture) => {
      fixture.c1PlanAndProfileBindings.runtimeProfile.selectedIdentity = "candidate";
    }),
  ],
  [
    "activated runtime admission",
    mutate(validC2A, (fixture) => {
      fixture.effects.admission = true;
    }),
  ],
  [
    "unknown authority field",
    mutate(validC2A, (fixture) => {
      fixture.runnerDescriptorProfile.entries[0].hostPath = "/tmp";
    }),
  ],
];

for (const [name, fixture] of invalidC2A) {
  if (validateC2A(fixture)) {
    throw new Error(`invalid C2A fixture was accepted: ${name}`);
  }
  process.stdout.write(`rejected invalid C2A fixture: ${name}\n`);
}

const c2bSchema = schemas.get("governed-deno-core-c2b-passive-binding.schema.json");
const validateC2B = ajv.getSchema(c2bSchema.$id);
const validC2B = JSON.parse(
  await readFile(
    new URL("../schemas/conformance/c2b-governed-deno-core/passive-binding.json", import.meta.url),
    "utf8",
  ),
);
const invalidC2B = [
  [
    "wrong domain",
    mutate(validC2B, (fixture) => {
      fixture.domain = "capsule.other/v1";
    }),
  ],
  [
    "unknown nested authority field",
    mutate(validC2B, (fixture) => {
      fixture.fixedFixture.authority.callerBytes = true;
    }),
  ],
  [
    "changed fixed-fixture cap",
    mutate(validC2B, (fixture) => {
      fixture.fixedFixture.source.maximumBytes += 1;
    }),
  ],
  [
    "invented admission",
    mutate(validC2B, (fixture) => {
      fixture.effects.admission = true;
    }),
  ],
  [
    "missing merge dependency",
    mutate(validC2B, (fixture) => {
      fixture.dependencyMergePolicy.dependencies.pop();
    }),
  ],
];

for (const [name, fixture] of invalidC2B) {
  if (validateC2B(fixture)) {
    throw new Error(`invalid C2B passive binding was accepted: ${name}`);
  }
  process.stdout.write(`rejected invalid C2B passive binding: ${name}\n`);
}

function mutate(value, mutation) {
  const copy = structuredClone(value);
  mutation(copy);
  return copy;
}
