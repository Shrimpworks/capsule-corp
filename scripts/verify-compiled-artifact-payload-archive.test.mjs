import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { dirname, join, resolve } from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";
import { verifyCompiledArtifactManifest } from "./verify-compiled-artifact-payload-archive.mjs";

const repositoryRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const manifest = JSON.parse(
  await readFile(
    join(repositoryRoot, "schemas/conformance/compiled-artifact-payloads/manifest.json"),
    "utf8",
  ),
);

function mutate(mutator) {
  const candidate = structuredClone(manifest);
  mutator(candidate);
  return candidate;
}

async function rejectsMutation(mutator, readFixture) {
  await assert.rejects(() => verifyCompiledArtifactManifest(mutate(mutator), readFixture));
}

test("accepts the closed compact fixture and retained source bytes", async () => {
  await verifyCompiledArtifactManifest(manifest);
});

test("rejects archive and source pin mutations", async () => {
  await rejectsMutation((candidate) => {
    candidate.archive.commit = "0".repeat(40);
  });
  await rejectsMutation((candidate) => {
    candidate.archive.pullRequest = "https://github.com/Shrimpworks/capsule-experiments/releases";
  });
  await rejectsMutation((candidate) => {
    candidate.source.commit = "f".repeat(40);
  });
});

test("rejects compiled identity, placement, role, and mode mutations", async () => {
  await rejectsMutation((candidate) => {
    candidate.payloads.uniqueCompiledIdentities[0].sha256 = "0".repeat(64);
  });
  await rejectsMutation((candidate) => {
    candidate.payloads.uniqueCompiledIdentities[1].placements.pop();
  });
  await rejectsMutation((candidate) => {
    candidate.payloads.uniqueCompiledIdentities[1].placements[0] =
      candidate.payloads.uniqueCompiledIdentities[2].placements[0];
  });
  await rejectsMutation((candidate) => {
    candidate.payloads.uniqueCompiledIdentities[0].mode = "644";
  });
  await rejectsMutation((candidate) => {
    candidate.payloads.uniqueCompiledIdentities[0].placements[0] =
      "mjs-source-validator-v1/dist/renamed-validator";
  });
  await rejectsMutation((candidate) => {
    const daemon = candidate.payloads.uniqueCompiledIdentities[1];
    const broker = candidate.payloads.uniqueCompiledIdentities[3];
    [daemon.sha256, broker.sha256] = [broker.sha256, daemon.sha256];
  });
});

test("rejects inactive-policy and binary-vector mutations", async () => {
  await rejectsMutation((candidate) => {
    candidate.payloads.binaryVectorIdentities[1].sha256 =
      candidate.payloads.binaryVectorIdentities[2].sha256;
  });
  await rejectsMutation((candidate) => {
    candidate.payloads.binaryVectorIdentities[2].placements.push(
      "macos-i2b2-unsigned-installation-bundle/extra-policy.bin",
    );
  });
});

test("rejects R3 payload, distribution, and Release-asset claim mutations", async () => {
  await rejectsMutation((candidate) => {
    candidate.r3Evidence.trackedSignedPayload = true;
  });
  await rejectsMutation((candidate) => {
    candidate.r3Evidence.releaseAssetPublished = true;
  });
  await rejectsMutation((candidate) => {
    candidate.r3Evidence.distribution = "Developer ID";
  });
  await rejectsMutation((candidate) => {
    candidate.r3Evidence.signedComponentExecutableSha256[0] = "0".repeat(64);
  });
});

test("rejects historical binding and archived mutation-oracle omissions", async () => {
  await rejectsMutation((candidate) => {
    candidate.historicalBindings.i1aBundleManifest.sha256 = "0".repeat(64);
  });
  await rejectsMutation((candidate) => {
    candidate.archivedBindings.tests.pop();
  });
  await rejectsMutation((candidate) => {
    candidate.archivedBindings.reproductionRequirements.pop();
  });
  await rejectsMutation((candidate) => {
    candidate.archivedBindings.reproductionRequirements[0] = "latest toolchain";
  });
  await rejectsMutation((candidate) => {
    candidate.currentExactByteRecords.pop();
  });
});

test("rejects unsafe, missing, and mutated local source fixtures", async () => {
  await rejectsMutation((candidate) => {
    candidate.sourceFixtures[0].localPath = "../outside";
  });
  await rejectsMutation((candidate) => {
    candidate.sourceFixtures.pop();
  });
  const tamperedPath = manifest.sourceFixtures[0].localPath;
  await assert.rejects(() =>
    verifyCompiledArtifactManifest(manifest, (path) => {
      if (path === tamperedPath) return Buffer.from("tampered");
      return readFile(join(repositoryRoot, path));
    }),
  );
});
