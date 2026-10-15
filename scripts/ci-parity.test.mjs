import assert from "node:assert/strict";
import { execFileSync, spawnSync } from "node:child_process";
import {
  chmodSync,
  cpSync,
  mkdirSync,
  mkdtempSync,
  readFileSync,
  realpathSync,
  rmSync,
  symlinkSync,
  writeFileSync,
} from "node:fs";
import { tmpdir } from "node:os";
import { dirname, join, resolve } from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";

const root = fileURLToPath(new URL("../", import.meta.url));
const verifier = join(root, "scripts/verify-exported-jsdoc.mjs");
const goLint = join(root, "scripts/lint-go.mjs");
const missing = "export function undocumented() {}\n";
const documented = "/** Documented. */\nexport function documented() {}\n";

function fixture(t) {
  const cwd = mkdtempSync(join(tmpdir(), "capsule-ci-parity-"));
  t.after(() => rmSync(cwd, { recursive: true, force: true }));
  const git = (...args) => execFileSync("git", args, { cwd, encoding: "utf8" }).trim();
  const write = (path, content) => {
    mkdirSync(dirname(join(cwd, path)), { recursive: true });
    writeFileSync(join(cwd, path), content);
  };
  git("init", "--quiet");
  git("config", "core.hooksPath", "/dev/null");
  git("config", "commit.gpgsign", "false");
  git("config", "user.name", "Local test fixture");
  git("config", "user.email", "fixture@example.invalid");
  write("packages/example/src/backlog.ts", missing);
  write("packages/example/src/deleted.ts", missing);
  git("add", ".");
  git("commit", "--quiet", "-m", "baseline");
  const base = git("rev-parse", "HEAD");
  const run = (...args) =>
    spawnSync(process.execPath, [verifier, ...args], { cwd, encoding: "utf8" });
  return { cwd, git, write, base, run };
}

function output(result) {
  return `${result.stdout}\n${result.stderr}`;
}

test("changed JSDoc rejects undocumented exports and accepts documented exports", (t) => {
  const f = fixture(t);
  f.write("packages/example/src/changed.ts", missing);
  f.git("add", ".");
  f.git("commit", "--quiet", "-m", "changed source");
  let result = f.run("--changed-since", f.base);
  assert.notEqual(result.status, 0);
  assert.match(output(result), /changed.ts:1: exported undocumented/);
  f.write("packages/example/src/changed.ts", documented);
  f.git("add", ".");
  f.git("commit", "--quiet", "-m", "document source");
  result = f.run("--changed-since", f.base);
  assert.equal(result.status, 0, output(result));
});

test("empty changed selection succeeds without auditing unchanged backlog", (t) => {
  const f = fixture(t);
  const result = f.run("--changed-since", f.base);
  assert.equal(result.status, 0, output(result));
  assert.match(result.stdout, /no in-scope/);
});

test("invalid or missing baseline fails even when no TypeScript files changed", (t) => {
  const f = fixture(t);
  for (const args of [
    ["--changed-since", "not-a-ref"],
    ["--changed-since"],
    ["--include-worktree"],
  ]) {
    const result = f.run(...args);
    assert.notEqual(result.status, 0);
    assert.match(output(result), /baseline|requires|usage/i);
  }
});

test("changed selection handles spaces, newlines, renames, deletion and TS/TSX test exclusions", (t) => {
  const f = fixture(t);
  f.git("mv", "packages/example/src/backlog.ts", "packages/example/src/renamed source\nfile.tsx");
  f.git("rm", "packages/example/src/deleted.ts");
  f.write("packages/example/outside src.ts", documented);
  f.write("packages/example/src/ignored.test.ts", missing);
  f.write("packages/example/src/ignored.test.tsx", missing);
  f.write("packages/example/src/ignored.d.ts", missing);
  f.git("add", ".");
  f.git("commit", "--quiet", "-m", "rename and delete");
  let result = f.run("--changed-since", f.base);
  assert.notEqual(result.status, 0);
  assert.match(output(result), /renamed source\nfile.tsx:1: exported undocumented/);
  f.write("packages/example/src/renamed source\nfile.tsx", documented);
  f.git("add", ".");
  f.git("commit", "--quiet", "-m", "document renamed source");
  result = f.run("--changed-since", f.base);
  assert.equal(result.status, 0, output(result));
  f.write("packages/example/outside src.ts", missing);
  f.git("add", ".");
  f.git("commit", "--quiet", "-m", "undocumented outside src");
  result = f.run("--changed-since", f.base);
  assert.notEqual(result.status, 0);
  assert.match(output(result), /outside src.ts:1: exported undocumented/);
});

test("local selection includes staged, unstaged and untracked sources", (t) => {
  const f = fixture(t);
  for (const state of ["staged", "unstaged", "untracked"]) {
    const path = `packages/example/src/${state}.ts`;
    if (state === "unstaged") {
      f.write(path, documented);
      f.git("add", ".");
      f.git("commit", "--quiet", "-m", "prepare tracked source");
    }
    f.write(path, missing);
    if (state === "staged") f.git("add", path);
    const result = f.run("--changed-since", f.base, "--include-worktree");
    assert.notEqual(result.status, 0);
    assert.match(output(result), new RegExp(`${state}.ts:1: exported undocumented`));
    f.write(path, documented);
  }
  assert.equal(f.run("--changed-since", f.base, "--include-worktree").status, 0);
});

test("explicit-file JSDoc mode remains available", (t) => {
  const f = fixture(t);
  const result = f.run("packages/example/src/backlog.ts");
  assert.notEqual(result.status, 0);
  assert.match(output(result), /backlog.ts:1: exported undocumented/);
});

function commandRecorder(f, name) {
  const path = join(f.cwd, "bin", name);
  f.write(
    `bin/${name}`,
    `#!${process.execPath}\nimport { appendFileSync } from 'node:fs';\nconst args = process.argv.slice(2);\nappendFileSync(process.env.CI_TEST_LOG, JSON.stringify([${JSON.stringify(name)}, ...args]) + '\\n');\nif (${JSON.stringify(name)} === 'golangci-lint' && args[0] === 'version') console.log('golangci-lint has version ' + (process.env.CI_TEST_LINT_VERSION || '2.12.2'));\nif (args[0] === 'run' && process.env.CI_TEST_CACHE_LOG) appendFileSync(process.env.CI_TEST_CACHE_LOG, JSON.stringify(process.env.GOLANGCI_LINT_CACHE) + '\\n');\nif (process.env.CI_TEST_FAIL === args[0]) process.exit(17);\n`,
  );
  chmodSync(path, 0o755);
}

function recordedRun(f, command, args, extraEnv = {}) {
  const log = join(f.cwd, "calls.jsonl");
  writeFileSync(log, "");
  const result = spawnSync(command, args, {
    cwd: f.cwd,
    encoding: "utf8",
    env: {
      ...process.env,
      PATH: `${join(f.cwd, "bin")}:${process.env.PATH}`,
      CI_TEST_LOG: log,
      ...extraEnv,
    },
  });
  const calls = readFileSync(log, "utf8")
    .trim()
    .split("\n")
    .filter(Boolean)
    .map((line) => JSON.parse(line));
  return { result, calls };
}

test("Go lint validates its version and baseline, then retains full correctness and changed-line docs", (t) => {
  const f = fixture(t);
  commandRecorder(f, "golangci-lint");
  let run = recordedRun(f, process.execPath, [goLint, "--base", f.base]);
  assert.equal(run.result.status, 0, output(run.result));
  assert.deepEqual(run.calls, [
    ["golangci-lint", "version"],
    ["golangci-lint", "config", "verify"],
    ["golangci-lint", "run", "--disable", "revive", "./..."],
    ["golangci-lint", "run", "--enable-only", "revive", "--new-from-rev", f.base, "./..."],
  ]);
  run = recordedRun(f, process.execPath, [goLint, "--base", f.base], {
    CI_TEST_LINT_VERSION: "2.11.0",
  });
  assert.notEqual(run.result.status, 0);
  assert.match(output(run.result), /2.12.2/);
  assert.equal(run.calls.length, 1);
  run = recordedRun(f, process.execPath, [goLint, "--base", "missing-ref"]);
  assert.notEqual(run.result.status, 0);
  assert.equal(run.calls.filter((call) => call[1] === "run").length, 0);
  run = recordedRun(f, process.execPath, [goLint, "--full"]);
  assert.equal(run.result.status, 0, output(run.result));
  assert.deepEqual(run.calls.at(-1), ["golangci-lint", "run", "./..."]);
});

test("make ci covers missing gates once and stops on a failed site build", (t) => {
  const f = fixture(t);
  for (const name of ["pnpm", "go", "gofmt", "node", "golangci-lint"]) commandRecorder(f, name);
  const args = ["-j4", "-f", resolve(root, "Makefile"), "ci", `CI_BASE=${f.base}`];
  let run = recordedRun(f, "make", args);
  assert.equal(run.result.status, 0, output(run.result));
  for (const gate of [
    "audit:dependencies",
    "site:build",
    "verify:schemas",
    "verify:adrs",
    "lint",
    "verify:jsdoc",
    "check",
    "test:built",
  ]) {
    assert.equal(
      run.calls.filter((call) => call[0] === "pnpm" && call[1] === gate).length,
      1,
      gate,
    );
  }
  assert.equal(
    run.calls.filter((call) => call[0] === "pnpm" && ["build", "test"].includes(call[1])).length,
    0,
  );
  run = recordedRun(f, "make", args, { CI_TEST_FAIL: "site:build" });
  assert.notEqual(run.result.status, 0);
  assert.equal(
    run.calls.some((call) => call[1] === "test:built"),
    false,
  );
});

test("changed-file policy checks an old undocumented export when another line changes", (t) => {
  const f = fixture(t);
  f.write("packages/example/src/backlog.ts", `${missing}\n// Unrelated new comment.\n`);
  f.git("add", ".");
  f.git("commit", "--quiet", "-m", "touch backlog file");
  const result = f.run("--changed-since", f.base);
  assert.notEqual(result.status, 0);
  assert.match(output(result), /backlog.ts:1: exported undocumented/);
});

test("locally deleted branch additions are omitted and unrelated histories fail", (t) => {
  const f = fixture(t);
  f.write("packages/example/src/added.ts", missing);
  f.git("add", ".");
  f.git("commit", "--quiet", "-m", "add source");
  rmSync(join(f.cwd, "packages/example/src/added.ts"));
  const deleted = f.run("--changed-since", f.base, "--include-worktree");
  assert.equal(deleted.status, 0, output(deleted));
  f.git("checkout", "--orphan", "unrelated");
  f.git("rm", "-rf", ".");
  f.write("README.md", "Separate history\n");
  f.git("add", ".");
  f.git("commit", "--quiet", "-m", "unrelated root");
  assert.notEqual(f.run("--changed-since", f.base).status, 0);
});

test("no-argument mode continues to audit all in-scope source files", (t) => {
  const f = fixture(t);
  for (const path of [
    "scripts/verify-exported-jsdoc.mjs",
    "scripts/lib/exported-jsdoc-files.mjs",
    "scripts/lib/git-baseline.mjs",
  ]) {
    mkdirSync(dirname(join(f.cwd, path)), { recursive: true });
    cpSync(join(root, path), join(f.cwd, path));
  }
  symlinkSync(join(root, "node_modules"), join(f.cwd, "node_modules"), "dir");
  const result = spawnSync(process.execPath, ["scripts/verify-exported-jsdoc.mjs"], {
    cwd: f.cwd,
    encoding: "utf8",
  });
  assert.notEqual(result.status, 0);
  assert.match(output(result), /backlog.ts:1: exported undocumented/);
});

test("changed mode refuses subdirectory and nonrepository invocations", (t) => {
  const f = fixture(t);
  for (const cwd of [join(f.cwd, "packages"), dirname(f.cwd)]) {
    const result = spawnSync(process.execPath, [verifier, "--changed-since", f.base], {
      cwd,
      encoding: "utf8",
    });
    assert.notEqual(result.status, 0);
    assert.match(output(result), /repository root/);
  }
});

test("Go lint namespaces caller cache parents by canonical checkout and preserves repeat/full reuse", (t) => {
  const first = fixture(t);
  const second = fixture(t);
  for (const f of [first, second]) commandRecorder(f, "golangci-lint");
  const cacheParent = join(first.cwd, "shared-cache");
  const cacheLog = join(first.cwd, "cache-paths.jsonl");
  const env = { GOLANGCI_LINT_CACHE: cacheParent, CI_TEST_CACHE_LOG: cacheLog };
  const invoke = (f, args, override = {}) => {
    writeFileSync(cacheLog, "");
    const run = recordedRun(f, process.execPath, [goLint, ...args], { ...env, ...override });
    assert.equal(run.result.status, 0, output(run.result));
    const paths = readFileSync(cacheLog, "utf8")
      .trim()
      .split("\n")
      .map((line) => JSON.parse(line));
    assert.ok(paths.every((path) => path === paths[0]));
    return realpathSync(paths[0]);
  };
  const firstPath = invoke(first, ["--base", first.base]);
  assert.ok(firstPath.startsWith(`${realpathSync(cacheParent)}/`));
  assert.equal(invoke(first, ["--full"]), firstPath);
  assert.notEqual(invoke(second, ["--base", second.base]), firstPath);
  assert.equal(
    invoke(first, ["--base", first.base], { GOLANGCI_LINT_CACHE: "shared-cache" }),
    firstPath,
  );
  const alias = join(second.cwd, "first-alias");
  symlinkSync(first.cwd, alias, "dir");
  assert.equal(invoke({ ...first, cwd: alias }, ["--base", first.base]), firstPath);
  const defaultPath = invoke(first, ["--full"], { GOLANGCI_LINT_CACHE: "" });
  assert.ok(defaultPath.startsWith(`${realpathSync(first.cwd)}/.git/capsule-golangci-lint-cache/`));
  const badParent = join(first.cwd, "cache-is-a-file");
  writeFileSync(badParent, "do not overwrite");
  const failed = recordedRun(first, process.execPath, [goLint, "--full"], {
    GOLANGCI_LINT_CACHE: badParent,
  });
  assert.notEqual(failed.result.status, 0);
  assert.match(output(failed.result), /Cannot prepare checkout-scoped/);
  assert.equal(readFileSync(badParent, "utf8"), "do not overwrite");
});
