// Match CI's full security/correctness scan and changed-line revive policy.
// --full retains the separate issue #217 documentation-backlog audit.
import { execFileSync } from "node:child_process";
import { createHash } from "node:crypto";
import { accessSync, constants, mkdirSync, realpathSync } from "node:fs";
import { join, resolve } from "node:path";
import { gitOutput, resolveBaseline } from "./lib/git-baseline.mjs";

const args = process.argv.slice(2);
const full = args.length === 1 && args[0] === "--full";
if (!full && !(args.length === 2 && args[0] === "--base")) {
  throw new Error("Usage: node scripts/lint-go.mjs --base REF | --full");
}
const version = execFileSync("golangci-lint", ["version"], { encoding: "utf8" });
if (!/\bversion 2\.12\.2(?:\s|$)/.test(version)) {
  throw new Error("Local CI requires golangci-lint 2.12.2, matching .github/workflows/ci.yml");
}
const base = full ? undefined : resolveBaseline(args[1]);
if (!full) {
  const untracked = gitOutput(["ls-files", "--others", "--exclude-standard", "-z"])
    .split("\0")
    .filter((path) => path.endsWith(".go"));
  if (untracked.length > 0) {
    throw new Error(
      "Stage new Go files before changed-line lint; untracked Go files cannot be included in the Git baseline diff",
    );
  }
}
// golangci-lint 2.12.2 caches absolute finding paths. Sharing those entries
// between checkouts can silently discard findings during changed-line filtering.
// Keep reusable entries per canonical checkout, under a caller-selected parent
// or Git's per-checkout metadata directory. Never clear the caller's cache.
const checkoutRoot = realpathSync(gitOutput(["rev-parse", "--show-toplevel"]).trim());
const cacheParent =
  process.env.GOLANGCI_LINT_CACHE ||
  gitOutput(["rev-parse", "--git-path", "capsule-golangci-lint-cache"]).trim();
const cacheKey = createHash("sha256").update(checkoutRoot).digest("hex");
const cachePath = join(resolve(cacheParent), `checkout-${cacheKey}`);
try {
  mkdirSync(cachePath, { recursive: true, mode: 0o700 });
  accessSync(cachePath, constants.R_OK | constants.W_OK | constants.X_OK);
} catch (cause) {
  throw new Error(
    `Cannot prepare checkout-scoped golangci-lint cache at ${JSON.stringify(cachePath)}`,
    { cause },
  );
}
const lint = (...lintArgs) =>
  execFileSync("golangci-lint", lintArgs, {
    stdio: "inherit",
    env: { ...process.env, GOLANGCI_LINT_CACHE: cachePath },
  });
lint("config", "verify");
if (full) {
  lint("run", "./...");
} else {
  lint("run", "--disable", "revive", "./...");
  lint("run", "--enable-only", "revive", "--new-from-rev", base, "./...");
}
