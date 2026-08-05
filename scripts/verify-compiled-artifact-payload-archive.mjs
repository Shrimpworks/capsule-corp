#!/usr/bin/env node
import assert from "node:assert/strict";
import { readFile } from "node:fs/promises";
import { dirname, isAbsolute, join, normalize, resolve } from "node:path";
import { fileURLToPath } from "node:url";
import { sha256Hex } from "./lib/fixture-bytes.mjs";

const repositoryRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const manifestPath = join(
  repositoryRoot,
  "schemas/conformance/compiled-artifact-payloads/manifest.json",
);

const expected = Object.freeze({
  sourceCommit: "bd926f436003d61a70c0312d9605804b2735449e",
  archiveCommit: "0944ffd8cfd01ec23e4ae99138b0931d56804077",
  archivePullRequest: "https://github.com/Shrimpworks/capsule-experiments/pull/5",
  archivePath: "experiments/completed-compiled-artifact-payloads",
  archiveVerifier: "experiments/completed-compiled-artifact-payloads/scripts/verify.mjs",
  closureDigests: Object.freeze({
    compiled: "045824bd3dcaa10f6ae16dbb94aea2c09f1f5b6e9fb57854dd6015361059e23a",
    vectors: "7d8878c2c4d1cf3ae7619bf79d75da4386a8414099de6aaa121e3556b82a12b7",
    r3: "71d216db9c8da0d64b95ba1457a8428f25f84380df784c8216c35729859b6720",
    fixtures: "14c28499b3a07d56df9d025b05d659a3fe6cda3cb66f1b7c2db3fc7a27928e7c",
    requirements: "11f9e9b6a6aff7108ed285b3e2276084de9b6928c1966ae3bd60b322910a21e0",
    coverage: "3f65266ab888dfed506208237e9e527a05963b5aa9921cfe4847ee90e9076f72",
    limitations: "9fdfb81cc91527bfc5a9f8121de81625b99290c41636b67bfb6775baf290d7ca",
  }),
  compiled: new Map([
    [
      "source-validator-v1-parser",
      [1146656, "ba2a6b38be6b4eea8c067887cf80988756e2f4a551d128bf2dabdaf7f2ecb600", 1],
    ],
    [
      "source-validator-r2-daemon-launcher",
      [35464, "4bc270c84f166dfb077d84458940411073f3c70a7f70db2e4af48601500b36cc", 3],
    ],
    [
      "source-validator-r2-daemon-parser",
      [1146560, "f54c349e3a61b06e0b4d482bc1ed28924ffe712a7ff2531f504e7b57917defc7", 3],
    ],
    [
      "source-validator-r2-broker-launcher",
      [35464, "81284de5ba54e2288602bee4e9aca4e4513211b560bacfd1286b7ab57c922613", 3],
    ],
    [
      "source-validator-r2-broker-parser",
      [1146560, "7abac7da99f4b9edef77bb5ecfff135e8b752e5ed656664632272079b5408577", 3],
    ],
    [
      "macos-i1a-app-shell",
      [59928, "365b8ebb5bb7dbd8823db7cc292c1b5807baa0fda4d09ba2d2905df7bee3cd5f", 2],
    ],
  ]),
  vectors: new Map([
    [
      "source-validator-v1-artifact-profile",
      [160, "8c9a96c1e835ea870da7b3bbabb4cf995ec34550257c8daed1381ced5aac4571", 1],
    ],
    [
      "source-validator-r2-daemon-inactive-policy",
      [256, "c198dac71f3b5c2d2e8cca34fc3e9c01ff7b8093ef1a881d8160a34800ff1098", 3],
    ],
    [
      "source-validator-r2-broker-inactive-policy",
      [256, "b0ce8504190b5fe9b0a0296c22340a6439ab453cb32f32c19ddb6e594698568d", 3],
    ],
  ]),
  historicalBindings: new Map([
    [
      "i1aBundleManifest",
      [
        "artifacts/macos-i1a-unsigned-app-shell/dist/Capsule.app/Contents/Resources/CapsuleConstruction/bundle-manifest.json",
        5546,
        "5bd80097775908031b1a4c90680e8c7656cc5e9f97df2cc187592f75ee67a56f",
      ],
    ],
    [
      "i1aConstructionEvidence",
      [
        "artifacts/macos-i1a-unsigned-app-shell/evidence/construction.json",
        2848,
        "31f79bdbd3dae29f6cfa340683ce59bc445041db0da12a66b1c051abc3db6ae5",
      ],
    ],
    [
      "i1bSignedDevelopmentEnrollment",
      [
        "artifacts/macos-i1b-r3-signed-development-composition/evidence/signed-enrollment.json",
        10643,
        "afc7002032fc1ff4ead29269e7a370d94524aff42ca9181827a03233a31fbc94",
      ],
    ],
  ]),
  archivedTests: [
    "scripts/verify-macos-i1a.test.mjs",
    "scripts/verify-macos-i1b-r3.test.mjs",
    "scripts/verify-macos-i2b2.test.mjs",
    "scripts/verify-mjs-source-validator-v1.test.mjs",
    "scripts/verify-source-validator-r2.test.mjs",
  ],
  archivedDocuments: [
    "docs/EXPERIMENT_ARCHIVE.md",
    "docs/MACOS_INSTALLATION_I1A_UNSIGNED_CONSTRUCTION.md",
    "docs/MACOS_INSTALLATION_I2B2_UNSIGNED_CONSTRUCTION.md",
    "docs/SOURCE_VALIDATOR_R3_EXECUTION_PACKET.md",
  ],
  archivedAdrs: [
    "docs/adr/0034-freeze-mjs-first-release-contract.md",
    "docs/adr/0035-select-disposable-mjs-source-validator.md",
    "docs/adr/0036-select-role-separated-source-validator-launchers.md",
    "docs/adr/0037-freeze-passive-macos-installation-i0-contract.md",
    "docs/adr/0038-select-one-shot-coordinator-supervisor-bootstrap.md",
    "docs/adr/0040-freeze-owner-only-internal-alpha-posture.md",
  ],
  currentExactByteRecords: [
    "docs/MACOS_INSTALLATION_I1A_UNSIGNED_CONSTRUCTION.md",
    "docs/MACOS_INSTALLATION_I2B2_UNSIGNED_CONSTRUCTION.md",
    "docs/MJS_SOURCE_VALIDATOR_IMPLEMENTATION_PLAN.md",
    "docs/PROJECT.md",
    "docs/SOURCE_VALIDATOR_R3_EXECUTION_PACKET.md",
    "docs/WORKSTREAM_EVIDENCE_LEDGER.md",
    "docs/adr/0035-select-disposable-mjs-source-validator.md",
    "schemas/conformance/macos-i2b2-unsigned-installation/profile.json",
    "scripts/generate-macos-i2b2-profile.mjs",
  ],
});

function assertSafeRelativePath(path, label) {
  assert.equal(typeof path, "string", `${label}: path must be a string`);
  assert.equal(isAbsolute(path), false, `${label}: absolute path forbidden`);
  assert.equal(normalize(path), path, `${label}: non-canonical path forbidden`);
  assert.equal(path.startsWith("../"), false, `${label}: traversal forbidden`);
  assert.equal(path.includes("\\"), false, `${label}: backslash forbidden`);
}

function verifyIdentities(identities, expectedIdentities, expectedPlacementCount, executable) {
  assert.equal(identities.length, expectedIdentities.size, "identity closure changed");
  const placements = [];
  for (const identity of identities) {
    const identityExpected = expectedIdentities.get(identity.id);
    assert.ok(identityExpected, `${identity.id}: unexpected identity`);
    const [bytes, digest, placementCount] = identityExpected;
    assert.equal(identity.bytes, bytes, `${identity.id}: byte length changed`);
    assert.equal(identity.sha256, digest, `${identity.id}: digest changed`);
    assert.equal(
      identity.placements.length,
      placementCount,
      `${identity.id}: placement closure changed`,
    );
    if (executable) assert.equal(identity.mode, "755", `${identity.id}: executable mode changed`);
    for (const path of identity.placements) {
      assertSafeRelativePath(path, identity.id);
      placements.push(path);
    }
  }
  assert.equal(placements.length, expectedPlacementCount, "placement count changed");
  assert.equal(new Set(placements).size, placements.length, "duplicate placement forbidden");
}

function assertExactArray(actual, expectedArray, label) {
  assert.deepEqual(actual, expectedArray, `${label}: closed list changed`);
}

function assertClosureDigest(value, expectedDigest, label) {
  assert.equal(
    sha256Hex(Buffer.from(JSON.stringify(value))),
    expectedDigest,
    `${label}: closure changed`,
  );
}

export async function verifyCompiledArtifactManifest(
  manifest,
  readFixture = async (path) => readFile(join(repositoryRoot, path)),
) {
  assert.equal(manifest.schema, "capsule.compiled-artifact-payload-archive/v1");
  assert.equal(manifest.status, "PASSED");
  assert.equal(manifest.source.repository, "Shrimpworks/capsule-corp");
  assert.equal(manifest.source.commit, expected.sourceCommit);
  assert.equal(manifest.archive.repository, "Shrimpworks/capsule-experiments");
  assert.equal(manifest.archive.commit, expected.archiveCommit);
  assert.equal(manifest.archive.pullRequest, expected.archivePullRequest);
  assert.equal(manifest.archive.path, expected.archivePath);
  assert.equal(manifest.archive.verifierPath, expected.archiveVerifier);
  assert.equal(manifest.archive.copiedFileCount, 210);
  assert.equal(manifest.archive.freshNetworkCloneVerified, true);

  assert.equal(manifest.payloads.trackedMachOPlacementCount, 15);
  verifyIdentities(manifest.payloads.uniqueCompiledIdentities, expected.compiled, 15, true);
  assertClosureDigest(
    manifest.payloads.uniqueCompiledIdentities,
    expected.closureDigests.compiled,
    "compiled identities",
  );
  assert.equal(manifest.payloads.trackedBinaryVectorPlacementCount, 7);
  verifyIdentities(manifest.payloads.binaryVectorIdentities, expected.vectors, 7, false);
  assertClosureDigest(
    manifest.payloads.binaryVectorIdentities,
    expected.closureDigests.vectors,
    "binary vectors",
  );

  assert.equal(manifest.r3Evidence.trackedSignedPayload, false);
  assert.equal(manifest.r3Evidence.releaseAssetPublished, false);
  assert.equal(manifest.r3Evidence.distribution, "Apple Development");
  assert.equal(manifest.r3Evidence.developerIdUsed, false);
  assert.equal(manifest.r3Evidence.notarizationUsed, false);
  assert.equal(manifest.r3Evidence.signedComponentExecutableSha256.length, 7);
  for (const digest of manifest.r3Evidence.signedComponentExecutableSha256) {
    assert.match(digest, /^[0-9a-f]{64}$/);
  }
  assert.equal(
    new Set(manifest.r3Evidence.signedComponentExecutableSha256).size,
    7,
    "R3 signed-component digest closure changed",
  );
  assertClosureDigest(
    manifest.r3Evidence.signedComponentExecutableSha256,
    expected.closureDigests.r3,
    "R3 signed components",
  );

  const historicalBindings = Object.entries(manifest.historicalBindings);
  assert.equal(
    historicalBindings.length,
    expected.historicalBindings.size,
    "historical binding closure changed",
  );
  for (const [key, binding] of historicalBindings) {
    const expectedBinding = expected.historicalBindings.get(key);
    assert.ok(expectedBinding, `${key}: unexpected historical binding`);
    assertSafeRelativePath(binding.path, "historical binding");
    assert.deepEqual(
      [binding.path, binding.bytes, binding.sha256],
      expectedBinding,
      `${key}: historical binding changed`,
    );
  }

  assert.equal(manifest.sourceFixtures.length, 6, "source-fixture closure changed");
  const fixturePaths = [];
  for (const fixture of manifest.sourceFixtures) {
    assertSafeRelativePath(fixture.historicalPath, "historical fixture");
    assertSafeRelativePath(fixture.localPath, "local fixture");
    assert.ok(
      fixture.localPath.startsWith(
        "schemas/conformance/compiled-artifact-payloads/source-fixtures/",
      ),
      `${fixture.localPath}: fixture escaped conformance root`,
    );
    const bytes = await readFixture(fixture.localPath);
    assert.equal(bytes.length, fixture.bytes, `${fixture.localPath}: byte length changed`);
    assert.equal(sha256Hex(bytes), fixture.sha256, `${fixture.localPath}: digest changed`);
    fixturePaths.push(fixture.localPath);
  }
  assert.equal(
    new Set(fixturePaths).size,
    fixturePaths.length,
    "duplicate source fixture forbidden",
  );
  assertClosureDigest(manifest.sourceFixtures, expected.closureDigests.fixtures, "source fixtures");

  assertExactArray(manifest.archivedBindings.tests, expected.archivedTests, "archived tests");
  assertExactArray(
    manifest.archivedBindings.documents,
    expected.archivedDocuments,
    "archived documents",
  );
  assertExactArray(manifest.archivedBindings.adrs, expected.archivedAdrs, "archived ADRs");
  assert.equal(manifest.archivedBindings.reproductionRequirements.length, 7);
  assert.equal(new Set(manifest.archivedBindings.reproductionRequirements).size, 7);
  assertClosureDigest(
    manifest.archivedBindings.reproductionRequirements,
    expected.closureDigests.requirements,
    "reproduction requirements",
  );
  assertExactArray(
    manifest.currentExactByteRecords,
    expected.currentExactByteRecords,
    "current exact-byte records",
  );
  assert.equal(manifest.mutationCoverage.length, 8);
  assert.equal(manifest.limitations.length, 3);
  assertClosureDigest(
    manifest.mutationCoverage,
    expected.closureDigests.coverage,
    "mutation coverage",
  );
  assertClosureDigest(manifest.limitations, expected.closureDigests.limitations, "limitations");
}

if (process.argv[1] && resolve(process.argv[1]) === fileURLToPath(import.meta.url)) {
  const manifest = JSON.parse(await readFile(manifestPath, "utf8"));
  await verifyCompiledArtifactManifest(manifest);
  console.log(
    `verified pinned archive ${manifest.archive.commit}, 15 Mach-O placements, seven binary-vector placements, and six local source fixtures`,
  );
}
