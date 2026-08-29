import assert from "node:assert/strict";
import { spawnSync } from "node:child_process";
import {
  chmod,
  lstat,
  mkdir,
  mkdtemp,
  readdir,
  readFile,
  rename,
  rm,
  symlink,
  unlink,
  writeFile,
} from "node:fs/promises";
import { tmpdir } from "node:os";
import { dirname, join, relative, resolve } from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";
import { sha256Hex } from "./lib/fixture-bytes.mjs";
import { verifyCompiledArtifactArchive } from "./verify-compiled-artifact-archive.mjs";

const repositoryRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const remoteScript = join(repositoryRoot, "scripts/verify-compiled-artifact-archive-remote.sh");
const localVerifier = join(repositoryRoot, "scripts/verify-compiled-artifact-archive.mjs");

async function inventory(root, current = root) {
  const entries = [];
  for (const name of (await readdir(current)).sort()) {
    const path = join(current, name);
    const stat = await lstat(path);
    if (stat.isDirectory()) entries.push(...(await inventory(root, path)));
    else {
      entries.push({
        path: relative(root, path),
        mode: stat.mode & 0o777,
        bytes: stat.size,
        sha256: sha256Hex(await readFile(path)),
      });
    }
  }
  return entries;
}

async function writeIndexes(root) {
  const entries = [
    ...(await inventory(root, join(root, "payloads"))),
    ...(await inventory(root, join(root, "source-bindings"))),
  ].sort((left, right) => left.path.localeCompare(right.path));
  await writeFile(
    join(root, "SOURCE_FILES.txt"),
    `${entries
      .map((entry) => `${entry.mode.toString(8).padStart(3, "0")} ${entry.bytes} ${entry.path}`)
      .join("\n")}\n`,
  );
  await writeFile(
    join(root, "SHA256SUMS"),
    `${entries.map((entry) => `${entry.sha256}  ${entry.path}`).join("\n")}\n`,
  );
}

async function createArchiveFixture() {
  const root = await mkdtemp(join(tmpdir(), "capsule-archive-verifier-test-"));
  const payloadRoot = join(root, "payloads/capsule-corp/artifacts");
  const executable = Buffer.from("bin");
  const signedEnrollment = {
    status: "PASSED",
    developerIdUsed: false,
    notarizationUsed: false,
    signedComponents: [{ signedExecutableSha256: "a".repeat(64) }],
  };
  await mkdir(join(payloadRoot, "r3/evidence"), { recursive: true });
  await mkdir(join(root, "source-bindings"), { recursive: true });
  await writeFile(join(payloadRoot, "bin"), executable);
  await chmod(join(payloadRoot, "bin"), 0o755);
  await writeFile(join(payloadRoot, "copy"), executable);
  await chmod(join(payloadRoot, "copy"), 0o755);
  await writeFile(
    join(payloadRoot, "r3/evidence/signed-enrollment.json"),
    JSON.stringify(signedEnrollment),
  );
  await writeFile(join(root, "source-bindings/binding.txt"), "binding\n");
  await writeIndexes(root);

  return {
    root,
    contract: {
      includedRoots: ["payloads", "source-bindings"],
      payloadRoot: "payloads/capsule-corp/artifacts",
      expectedFileCount: 4,
      expectedMachO: new Map([["bin", [executable.length, sha256Hex(executable)]]]),
      crossCopies: new Map([["copy", "bin"]]),
      policyCopies: new Map(),
      signedEnrollment: "r3/evidence/signed-enrollment.json",
      signedComponentCount: 1,
      trackedMachOPlacementCount: 2,
    },
  };
}

test("Capsule-owned verifier accepts a closed archive and rejects inventory mutations", async (t) => {
  const fixture = await createArchiveFixture();
  t.after(() => rm(fixture.root, { recursive: true, force: true }));
  assert.deepEqual(await verifyCompiledArtifactArchive(fixture.root, fixture.contract), {
    copiedFileCount: 4,
    trackedMachOPlacementCount: 2,
    uniqueCompiledIdentityCount: 1,
  });

  await writeFile(join(fixture.root, "source-bindings/extra.txt"), "extra\n");
  await assert.rejects(() => verifyCompiledArtifactArchive(fixture.root, fixture.contract));
  await unlink(join(fixture.root, "source-bindings/extra.txt"));

  await unlink(join(fixture.root, "source-bindings/binding.txt"));
  await assert.rejects(() => verifyCompiledArtifactArchive(fixture.root, fixture.contract));
});

test("Capsule-owned verifier rejects cross-copy substitution and symlinks", async (t) => {
  const fixture = await createArchiveFixture();
  t.after(() => rm(fixture.root, { recursive: true, force: true }));

  await writeFile(join(fixture.root, "payloads/capsule-corp/artifacts/copy"), "bad");
  await writeIndexes(fixture.root);
  await assert.rejects(
    () => verifyCompiledArtifactArchive(fixture.root, fixture.contract),
    /cross-artifact copy changed/,
  );

  await unlink(join(fixture.root, "payloads/capsule-corp/artifacts/copy"));
  await symlink("bin", join(fixture.root, "payloads/capsule-corp/artifacts/copy"));
  await assert.rejects(
    () => verifyCompiledArtifactArchive(fixture.root, fixture.contract),
    /symlink forbidden/,
  );
});

for (const includedRoot of ["payloads", "source-bindings"]) {
  test(`Capsule-owned verifier rejects symlinked ${includedRoot} root`, async (t) => {
    const fixture = await createArchiveFixture();
    t.after(() => rm(fixture.root, { recursive: true, force: true }));

    const movedRoot = `${includedRoot}-real`;
    await rename(join(fixture.root, includedRoot), join(fixture.root, movedRoot));
    await symlink(movedRoot, join(fixture.root, includedRoot));

    await assert.rejects(
      () => verifyCompiledArtifactArchive(fixture.root, fixture.contract),
      new RegExp(`${includedRoot}: symlink forbidden`),
    );
  });
}

async function writeExecutable(path, body) {
  await writeFile(path, body);
  await chmod(path, 0o755);
}

async function runRemoteHarness(nodeStatus) {
  const root = await mkdtemp(join(tmpdir(), "capsule-archive-shell-test-"));
  const fakeBin = join(root, "bin");
  const verificationRoot = join(root, "capsule-artifact-archive.abc123");
  const nodeArgs = join(root, "node-args");
  await mkdir(fakeBin);
  await writeExecutable(
    join(fakeBin, "mktemp"),
    '#!/bin/sh\nmkdir -p "$CAPSULE_TEST_ROOT"\nprintf "%s\\n" "$CAPSULE_TEST_ROOT"\n',
  );
  await writeExecutable(
    join(fakeBin, "git"),
    '#!/bin/sh\ncase "$*" in *"rev-parse FETCH_HEAD"*) printf "%s\\n" "0944ffd8cfd01ec23e4ae99138b0931d56804077";; esac\nexit 0\n',
  );
  await writeExecutable(
    join(fakeBin, "node"),
    '#!/bin/sh\nprintf "%s\\n" "$*" > "$CAPSULE_NODE_ARGS"\nexit "$CAPSULE_NODE_STATUS"\n',
  );

  const result = spawnSync("/bin/sh", [remoteScript], {
    encoding: "utf8",
    env: {
      ...process.env,
      PATH: `${fakeBin}:/usr/bin:/bin`,
      TMPDIR: root,
      CAPSULE_TEST_ROOT: verificationRoot,
      CAPSULE_NODE_ARGS: nodeArgs,
      CAPSULE_NODE_STATUS: String(nodeStatus),
    },
  });
  const invokedArgs = await readFile(nodeArgs, "utf8");
  const checkoutExists = await lstat(verificationRoot).then(
    () => true,
    () => false,
  );
  await rm(root, { recursive: true, force: true });
  return { result, invokedArgs, checkoutExists };
}

test("remote harness executes only the Capsule-owned verifier and cleans success", async () => {
  const { result, invokedArgs, checkoutExists } = await runRemoteHarness(0);
  assert.equal(result.status, 0, result.stderr);
  assert.equal(checkoutExists, false);
  assert.ok(invokedArgs.startsWith(`${localVerifier} `));
  assert.equal(invokedArgs.includes("/scripts/verify.mjs"), false);
});

test("remote harness preserves verifier failure and cleans the checkout", async () => {
  const { result, invokedArgs, checkoutExists } = await runRemoteHarness(7);
  assert.equal(result.status, 7, result.stderr);
  assert.equal(checkoutExists, false);
  assert.ok(invokedArgs.startsWith(`${localVerifier} `));
  assert.equal(invokedArgs.includes("/scripts/verify.mjs"), false);
});
