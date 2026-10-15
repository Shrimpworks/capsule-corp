// Explicit local integration suite: requires Go from go.mod and golangci-lint 2.12.2.
// Kept outside *.test.mjs: the TypeScript CI job does not install the Go linter.
import assert from "node:assert/strict";
import { execFileSync, spawnSync } from "node:child_process";
import { copyFileSync, mkdtempSync, readFileSync, rmSync, writeFileSync } from "node:fs";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";

const root = fileURLToPath(new URL("../", import.meta.url));

function fixture(t, source = "package fixture\n\nfunc Backlog() {}\n") {
  const cwd = mkdtempSync(join(tmpdir(), "capsule-go-parity-"));
  t.after(() => rmSync(cwd, { recursive: true, force: true }));
  const git = (...args) => execFileSync("git", args, { cwd, encoding: "utf8" }).trim();
  git("init", "--quiet");
  git("config", "core.hooksPath", "/dev/null");
  git("config", "commit.gpgsign", "false");
  git("config", "user.name", "Local fixture");
  git("config", "user.email", "fixture@example.invalid");
  const version = readFileSync(join(root, "go.mod"), "utf8").match(/^go (.+)$/m)[1];
  writeFileSync(join(cwd, "go.mod"), `module example.invalid/parity\n\ngo ${version}\n`);
  writeFileSync(join(cwd, "source.go"), source);
  copyFileSync(join(root, ".golangci.yml"), join(cwd, ".golangci.yml"));
  git("add", ".");
  git("commit", "--quiet", "-m", "baseline");
  const base = git("rev-parse", "HEAD");
  const run = (env = process.env) =>
    spawnSync(process.execPath, [join(root, "scripts/lint-go.mjs"), "--base", base], {
      cwd,
      env,
      encoding: "utf8",
      timeout: 120000,
    });
  return { cwd, git, run };
}

function output(result) {
  return `${result.stdout}\n${result.stderr}`;
}

test("real pinned linter excludes old docs but rejects committed, staged and unstaged new docs", (t) => {
  const f = fixture(t);
  let result = f.run();
  assert.equal(result.status, 0, output(result));
  for (const state of ["committed", "staged", "unstaged"]) {
    writeFileSync(
      join(f.cwd, "source.go"),
      `package fixture\n\nfunc Backlog() {}\n\nfunc Added${state}() {}\n`,
    );
    if (state !== "unstaged") f.git("add", "source.go");
    if (state === "committed") f.git("commit", "--quiet", "-m", "new undocumented function");
    result = f.run();
    assert.notEqual(result.status, 0);
    assert.match(output(result), new RegExp(`Added${state}.*comment|comment.*Added${state}`));
  }
  writeFileSync(
    join(f.cwd, "source.go"),
    "package fixture\n\nfunc Backlog() {}\n\n// Added documents a new export.\nfunc Added() {}\n",
  );
  result = f.run();
  assert.equal(result.status, 0, output(result));
});

test("untracked Go source is refused until staged, then undocumented export is detected", (t) => {
  const f = fixture(t);
  writeFileSync(join(f.cwd, "new.go"), "package fixture\n\nfunc Untracked() {}\n");
  let result = f.run();
  assert.notEqual(result.status, 0);
  assert.match(output(result), /Stage new Go files/);
  f.git("add", "new.go");
  result = f.run();
  assert.notEqual(result.status, 0);
  assert.match(output(result), /Untracked.*comment|comment.*Untracked/);
});

test("unchanged security findings still fail the whole-code scan", (t) => {
  const f = fixture(
    t,
    'package fixture\n\nimport "crypto/tls"\n\nfunc Insecure() *tls.Config { return &tls.Config{InsecureSkipVerify: true} }\n',
  );
  const result = f.run();
  assert.notEqual(result.status, 0);
  assert.match(output(result), /G402/);
});

test("identical modules in two retained checkouts enforce new docs through a shared cache parent on repeat runs", (t) => {
  const cacheParent = mkdtempSync(join(tmpdir(), "capsule-shared-lint-cache-"));
  t.after(() => rmSync(cacheParent, { recursive: true, force: true }));
  const fixtures = [fixture(t), fixture(t)];
  for (const f of fixtures) {
    writeFileSync(
      join(f.cwd, "source.go"),
      "package fixture\n\nfunc Backlog() {}\n\nfunc Added() {}\n",
    );
    f.git("add", "source.go");
    f.git("commit", "--quiet", "-m", "add undocumented export");
  }
  for (const round of [1, 2]) {
    for (const [index, f] of fixtures.entries()) {
      const result = f.run({ ...process.env, GOLANGCI_LINT_CACHE: cacheParent });
      assert.notEqual(result.status, 0, `round ${round}, checkout ${index}: ${output(result)}`);
      assert.match(output(result), /Added.*comment|comment.*Added/);
    }
  }
});
