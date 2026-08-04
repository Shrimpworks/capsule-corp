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
