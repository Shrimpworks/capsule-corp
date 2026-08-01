import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { cp, mkdir, mkdtemp, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { verifyConformanceCorpus } from "./verify-contract-conformance.mjs";

const schemaPath = new URL("../schemas/conformance/v0/manifest.schema.json", import.meta.url);

test("verifies a closed, complete, byte-exact conformance corpus", async (t) => {
  const corpus = await createCorpus(t);

  const result = await verifyConformanceCorpus({ rootDirectory: corpus.root });

  assert.deepEqual(result, { caseCount: 2, fixtureCount: 2, ruleCount: 1 });
});

test("rejects an unknown manifest field", async (t) => {
  const corpus = await createCorpus(t, (manifest) => {
    manifest.unreviewed = true;
  });

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /manifest schema validation failed/u,
  );
});

test("rejects a listed fixture that is missing", async (t) => {
  const corpus = await createCorpus(t);
  await rm(join(corpus.root, "shared/reject.bin"));

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /listed fixture does not exist/u,
  );
});

test("rejects a fixture that is present but unlisted", async (t) => {
  const corpus = await createCorpus(t);
  await writeFile(join(corpus.root, "shared/unlisted.bin"), "unlisted");

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /fixture is not listed/u,
  );
});

test("rejects fixture bytes that do not match the retained hash", async (t) => {
  const corpus = await createCorpus(t);
  await writeFile(join(corpus.root, "shared/accept.bin"), "mutated");

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /fixture byte length mismatch|fixture SHA-256 mismatch/u,
  );
});

test("rejects deletion of a required rule outcome and its fixture", async (t) => {
  const corpus = await createCorpus(t, (manifest) => {
    manifest.cases = manifest.cases.filter((entry) => entry.expected.decision !== "reject");
  });
  await rm(join(corpus.root, "shared/reject.bin"));

  await assert.rejects(
    verifyConformanceCorpus({ rootDirectory: corpus.root }),
    /required case is not covered/u,
  );
});

async function createCorpus(t, mutateManifest = () => {}) {
  const root = await mkdtemp(join(tmpdir(), "capsule-conformance-"));
  t.after(async () => rm(root, { recursive: true, force: true }));
  await mkdir(join(root, "shared"));
  await cp(schemaPath, join(root, "manifest.schema.json"));

  const acceptBytes = Buffer.from("accepted", "utf8");
  const rejectBytes = Buffer.from("rejected", "utf8");
  await writeFile(join(root, "shared/accept.bin"), acceptBytes);
  await writeFile(join(root, "shared/reject.bin"), rejectBytes);

  const implementations = {
    "fixture-integrity": "verified",
    go: "pending",
    typescript: "pending",
    swift: "not-applicable",
  };
  const manifest = {
    manifestVersion: "capsule.conformance/v0",
    rules: [
      {
        id: "job-proposal.raw.example",
        source: "ADR-0023#strict-jobproposal-raw-profile",
        description: "Test-only representative rule.",
        requiredCases: [
          { decision: "accept", variant: "ordinary" },
          { decision: "reject", variant: "malformed" },
        ],
      },
    ],
    cases: [
      {
        id: "job-proposal.raw.example.accept",
        description: "Accept representative bytes.",
        ruleIds: ["job-proposal.raw.example"],
        object: "JobProposal",
        wireFormat: "raw-bytes",
        mediaType: "application/capsule.job-proposal+json;v=0",
        variant: "ordinary",
        fixture: reference("shared/accept.bin", acceptBytes),
        context: { kind: "none" },
        expected: {
          decision: "accept",
          classification: null,
          owner: "public-raw-decoder",
          authorityStateChanged: false,
        },
        implementations,
      },
      {
        id: "job-proposal.raw.example.reject",
        description: "Reject representative bytes.",
        ruleIds: ["job-proposal.raw.example"],
        object: "JobProposal",
        wireFormat: "raw-bytes",
        mediaType: "application/capsule.job-proposal+json;v=0",
        variant: "malformed",
        fixture: reference("shared/reject.bin", rejectBytes),
        context: { kind: "none" },
        expected: {
          decision: "reject",
          classification: "MALFORMED",
          owner: "public-raw-decoder",
          authorityStateChanged: false,
        },
        implementations,
      },
    ],
  };
  mutateManifest(manifest);
  await writeFile(join(root, "manifest.json"), `${JSON.stringify(manifest, null, 2)}\n`);

  return { root };
}

function reference(path, bytes) {
  return {
    path,
    sha256: createHash("sha256").update(bytes).digest("hex"),
    byteLength: bytes.length,
  };
}
