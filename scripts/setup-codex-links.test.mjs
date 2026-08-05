import assert from "node:assert/strict";
import { execFile } from "node:child_process";
import { mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { promisify } from "node:util";
import {
  evaluatePin,
  findGitRoot,
  readPin,
  resolveAiCentralCommit,
  writePin,
} from "./setup-codex-links.mjs";

const execFileAsync = promisify(execFile);

async function makeGitCheckout(root) {
  await mkdir(root, { recursive: true });
  await execFileAsync("git", ["init", "--quiet"], { cwd: root });
  await execFileAsync("git", ["config", "user.email", "test@example.invalid"], { cwd: root });
  await execFileAsync("git", ["config", "user.name", "Test"], { cwd: root });
  await writeFile(join(root, "seed.txt"), "seed\n");
  await execFileAsync("git", ["add", "."], { cwd: root });
  await execFileAsync("git", ["-c", "commit.gpgsign=false", "commit", "--quiet", "-m", "seed"], {
    cwd: root,
  });
  const { stdout } = await execFileAsync("git", ["rev-parse", "HEAD"], { cwd: root });
  return stdout.trim();
}

async function withTempDir(run) {
  const root = await mkdtemp(join(tmpdir(), "codex-pin-"));
  try {
    await run(root);
  } finally {
    await rm(root, { recursive: true, force: true });
  }
}

test("findGitRoot walks up from a nested directory to the checkout root", async () => {
  await withTempDir(async (root) => {
    await makeGitCheckout(root);
    const nested = join(root, "templates", "skills");
    await mkdir(nested, { recursive: true });
    assert.equal(await findGitRoot(nested), root);
  });
});

test("findGitRoot returns undefined when no ancestor is a git checkout", async () => {
  await withTempDir(async (root) => {
    const nested = join(root, "templates");
    await mkdir(nested, { recursive: true });
    assert.equal(await findGitRoot(nested), undefined);
  });
});

test("resolveAiCentralCommit reports the checkout's current commit", async () => {
  await withTempDir(async (root) => {
    const commit = await makeGitCheckout(root);
    const templates = join(root, "templates");
    await mkdir(templates, { recursive: true });
    assert.equal(await resolveAiCentralCommit(templates), commit);
  });
});

test("resolveAiCentralCommit reports undefined for a non-git directory", async () => {
  await withTempDir(async (root) => {
    assert.equal(await resolveAiCentralCommit(root), undefined);
  });
});

test("readPin reports an unset pin when the file is missing", async () => {
  await withTempDir(async (root) => {
    assert.deepEqual(await readPin(join(root, "ai-central-pin.json")), { expectedCommit: null });
  });
});

test("writePin then readPin round-trips the recorded commit", async () => {
  await withTempDir(async (root) => {
    const pinPath = join(root, "ai-central-pin.json");
    await writePin(pinPath, "abc123");
    const pin = await readPin(pinPath);
    assert.equal(pin.expectedCommit, "abc123");
    const raw = await readFile(pinPath, "utf8");
    assert.match(raw, /"expectedCommit": "abc123"/u);
  });
});

test("evaluatePin classifies unpinned, matching, mismatched, and non-git states", () => {
  assert.equal(evaluatePin({ expectedCommit: null }, "abc").status, "unpinned");
  assert.equal(evaluatePin({ expectedCommit: null }, undefined).status, "unpinned");
  assert.equal(evaluatePin({ expectedCommit: "abc" }, "abc").status, "ok");
  assert.equal(evaluatePin({ expectedCommit: "abc" }, "def").status, "mismatch");
  assert.equal(evaluatePin({ expectedCommit: "abc" }, undefined).status, "not-git");
});
