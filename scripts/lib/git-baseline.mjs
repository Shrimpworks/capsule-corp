import { execFileSync } from "node:child_process";

/** Run Git without a shell; fail rather than treating an invalid diff as empty. */
export function gitOutput(args) {
  return execFileSync("git", args, { encoding: "utf8", maxBuffer: 16 * 1024 * 1024 });
}

/** Resolve an explicit commit/ref to its merge base with HEAD from the repository root. */
export function resolveBaseline(ref) {
  if (!ref || ref.startsWith("-")) {
    throw new Error("A nonempty Git baseline is required");
  }
  try {
    if (gitOutput(["rev-parse", "--show-prefix"]).trim() !== "") {
      throw new Error("Run from the repository root");
    }
    const commit = gitOutput([
      "rev-parse",
      "--verify",
      "--end-of-options",
      `${ref}^{commit}`,
    ]).trim();
    return gitOutput(["merge-base", commit, "HEAD"]).trim();
  } catch (cause) {
    throw new Error(
      `Cannot resolve baseline ${JSON.stringify(ref)}; fetch it and run from the repository root`,
      { cause },
    );
  }
}

/** Select complete changed TS/TSX files, preserving the CI packages-wide scope. */
export function changedTypeScriptFiles(ref, includeWorktree) {
  const base = resolveBaseline(ref);
  const paths = gitOutput([
    "diff",
    "--no-ext-diff",
    "--name-only",
    "-z",
    "--diff-filter=d",
    base,
    ...(includeWorktree ? [] : ["HEAD"]),
    "--",
    "packages/",
  ]).split("\0");
  if (includeWorktree) {
    paths.push(
      ...gitOutput(["ls-files", "--others", "--exclude-standard", "-z", "--", "packages/"]).split(
        "\0",
      ),
    );
  }
  return [...new Set(paths)]
    .filter(
      (path) =>
        /^packages\/.+\.tsx?$/s.test(path) &&
        !/\.test\.tsx?$/.test(path) &&
        !path.endsWith(".d.ts"),
    )
    .sort();
}
