#!/usr/bin/env node

import { execFile } from "node:child_process";
import fs from "node:fs/promises";
import path from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);

const scriptDirectory = path.dirname(fileURLToPath(import.meta.url));
const repositoryRoot = path.resolve(scriptDirectory, "..");
const codexDirectory = path.join(repositoryRoot, ".codex");
const skillsDirectory = path.join(codexDirectory, "skills");
const steeringDirectory = path.join(codexDirectory, "steering");
const pinPath = path.join(codexDirectory, "ai-central-pin.json");
const templatesRoot = resolveTemplatesRoot();
const dryRun = process.argv.includes("--dry-run");
const recordPin = process.argv.includes("--record-pin");
const sharedSteeringFiles = ["javascript-esm-steering.md"];
const isEntrypoint = Boolean(
  process.argv[1] && import.meta.url === pathToFileURL(process.argv[1]).href,
);

function usage() {
  process.stdout.write(`Usage: node scripts/setup-codex-links.mjs [--dry-run] [--record-pin]

Creates local .codex symlinks from AI Central while preserving repo-owned files.

Environment:
  AI_CENTRAL_HOME  Path to ai-central or ai-central/templates.
                   Defaults to ../ai-central/templates.

Options:
  --dry-run        Report changes without writing links.
  --record-pin     Record the AI Central checkout's current git commit into
                   .codex/ai-central-pin.json instead of linking. Run this
                   after reviewing the linked checkout; subsequent runs
                   refuse to link if the checkout's commit no longer matches.
  --help           Show this help.
`);
}

// This argument parsing/validation only makes sense for a direct CLI
// invocation; guard it so importing this module for its exported functions
// (e.g. from a test) never inspects the host process's own argv or exits it.
if (isEntrypoint) {
  if (process.argv.includes("--help") || process.argv.includes("-h")) {
    usage();
    process.exit(0);
  }

  const unknownArguments = process.argv
    .slice(2)
    .filter((argument) => !["--", "--dry-run", "--record-pin"].includes(argument));

  if (unknownArguments.length > 0) {
    process.stderr.write(`Unknown option: ${unknownArguments[0]}\n`);
    usage();
    process.exit(2);
  }
}

function resolveTemplatesRoot() {
  const input =
    process.env.AI_CENTRAL_HOME ?? path.resolve(repositoryRoot, "../ai-central/templates");
  const absolute = path.resolve(input);

  return path.basename(absolute) === "ai-central" ? path.join(absolute, "templates") : absolute;
}

export async function findGitRoot(startDirectory) {
  let current = startDirectory;

  for (;;) {
    if (await pathExists(path.join(current, ".git"))) {
      return current;
    }

    const parent = path.dirname(current);
    if (parent === current) {
      return undefined;
    }

    current = parent;
  }
}

export async function resolveAiCentralCommit(templatesRootPath) {
  const gitRoot = await findGitRoot(templatesRootPath);
  if (!gitRoot) {
    return undefined;
  }

  try {
    const { stdout } = await execFileAsync("git", ["-C", gitRoot, "rev-parse", "HEAD"]);
    return stdout.trim();
  } catch {
    return undefined;
  }
}

export async function readPin(pinFilePath) {
  try {
    const raw = await fs.readFile(pinFilePath, "utf8");
    return JSON.parse(raw);
  } catch (error) {
    if (error.code === "ENOENT") {
      return { expectedCommit: null };
    }

    throw error;
  }
}

export async function writePin(pinFilePath, expectedCommit) {
  const payload = {
    expectedCommit,
    note:
      "Recorded by `pnpm codex:links -- --record-pin` after reviewing the linked AI Central " +
      "checkout. setup-codex-links.mjs refuses to link when the checkout's current commit no " +
      "longer matches this value.",
  };

  await fs.writeFile(pinFilePath, `${JSON.stringify(payload, null, 2)}\n`);
}

// evaluatePin classifies the recorded pin against the checkout's current
// commit without performing any I/O, so the decision logic itself is a pure,
// directly testable function.
export function evaluatePin(pin, currentCommit) {
  if (pin.expectedCommit && currentCommit && pin.expectedCommit !== currentCommit) {
    return { status: "mismatch" };
  }

  if (!pin.expectedCommit) {
    return { status: "unpinned" };
  }

  if (!currentCommit) {
    return { status: "not-git" };
  }

  return { status: "ok" };
}

async function pathExists(target) {
  try {
    await fs.lstat(target);
    return true;
  } catch (error) {
    if (error.code === "ENOENT") {
      return false;
    }

    if (["EACCES", "EPERM"].includes(error.code)) {
      process.stderr.write(`warning: permission denied reading ${target}, treating as missing\n`);
      return false;
    }

    throw error;
  }
}

async function* walkDirectories(root) {
  let entries;

  try {
    entries = await fs.readdir(root, { withFileTypes: true });
  } catch (error) {
    if (error.code === "ENOENT") {
      return;
    }

    if (["EACCES", "EPERM"].includes(error.code)) {
      process.stderr.write(`warning: permission denied reading ${root}, skipping\n`);
      return;
    }

    throw error;
  }

  yield root;

  for (const entry of entries) {
    if (entry.isDirectory()) {
      yield* walkDirectories(path.join(root, entry.name));
    }
  }
}

function skillLinkName(parts, name) {
  if (!name || parts[0] === undefined) {
    return undefined;
  }

  if (parts[0] === "adapted" || parts[0] !== "imported") {
    return name;
  }

  switch (parts[1]) {
    case "agent-skills":
      return name;
    case "pm-skills":
      return `pm-${name}`;
    case "claude-skills":
      return `claude-${name}`;
    case "agent-toolkit":
      return `toolkit-${name}`;
    case "web-quality-skills":
      return `web-${name}`;
    default:
      return name;
  }
}

async function findSkillLinks() {
  const skillRoot = path.join(templatesRoot, "skills");
  const links = new Map();

  for await (const directory of walkDirectories(skillRoot)) {
    if (!(await pathExists(path.join(directory, "SKILL.md")))) {
      continue;
    }

    const relativeDirectory = path.relative(skillRoot, directory);
    const parts = relativeDirectory.split(path.sep);
    const linkName = skillLinkName(parts, parts.at(-1));

    if (!linkName) {
      continue;
    }

    const existingTarget = links.get(linkName);
    if (existingTarget && existingTarget !== directory) {
      throw new Error(`AI Central has duplicate skill link name '${linkName}'`);
    }

    links.set(linkName, directory);
  }

  return [...links.entries()]
    .map(([linkName, target]) => ({ linkName, target }))
    .sort((left, right) => left.linkName.localeCompare(right.linkName));
}

async function findSteeringLinks() {
  const root = path.join(templatesRoot, "steering");
  const links = [];

  for (const fileName of sharedSteeringFiles) {
    const target = path.join(root, fileName);

    if (await pathExists(target)) {
      links.push({ linkName: fileName, target });
    }
  }

  return links;
}

async function ensureSymlink(directory, linkName, target) {
  const linkPath = path.join(directory, linkName);
  let existing;

  try {
    existing = await fs.lstat(linkPath);
  } catch (error) {
    if (error.code !== "ENOENT") {
      throw error;
    }
  }

  if (existing && !existing.isSymbolicLink()) {
    return { action: "preserved", linkPath, target };
  }

  if (existing?.isSymbolicLink()) {
    const currentTarget = await fs.readlink(linkPath);

    if (path.resolve(path.dirname(linkPath), currentTarget) === target) {
      return { action: "unchanged", linkPath, target };
    }

    if (!dryRun) {
      await fs.unlink(linkPath);
    }
  }

  if (!dryRun) {
    await fs.symlink(target, linkPath);
  }

  return { action: existing ? "updated" : "created", linkPath, target };
}

async function main() {
  if (!(await pathExists(path.join(templatesRoot, "skills")))) {
    process.stderr.write(`AI Central templates not found: ${templatesRoot}\n`);
    process.stderr.write(
      "Set AI_CENTRAL_HOME to your ai-central checkout or templates directory.\n",
    );
    process.exitCode = 1;
    return;
  }

  if (recordPin) {
    const commit = await resolveAiCentralCommit(templatesRoot);
    if (!commit) {
      process.stderr.write(
        `AI Central at ${templatesRoot} is not a git checkout; cannot record a pinned commit.\n`,
      );
      process.exitCode = 1;
      return;
    }

    await writePin(pinPath, commit);
    process.stdout.write(`Recorded AI Central pin: ${commit}\n`);
    return;
  }

  const pin = await readPin(pinPath);
  const currentCommit = await resolveAiCentralCommit(templatesRoot);
  const evaluation = evaluatePin(pin, currentCommit);

  if (evaluation.status === "mismatch") {
    process.stderr.write(
      `AI Central checkout commit ${currentCommit} does not match the recorded pin ` +
        `${pin.expectedCommit} in .codex/ai-central-pin.json. Review what changed, then run ` +
        "`pnpm codex:links -- --record-pin` to accept it, or restore the pinned commit.\n",
    );
    process.exitCode = 1;
    return;
  }

  if (evaluation.status === "unpinned") {
    process.stderr.write(
      "warning: no pinned AI Central revision recorded (.codex/ai-central-pin.json); linked " +
        "content is not verified against a reviewed commit. Run " +
        "`pnpm codex:links -- --record-pin` after reviewing the linked checkout.\n",
    );
  } else if (evaluation.status === "not-git") {
    process.stderr.write(
      "warning: AI Central checkout is not a git repository; the recorded pin cannot be verified.\n",
    );
  }

  if (!dryRun) {
    await fs.mkdir(skillsDirectory, { recursive: true });
    await fs.mkdir(steeringDirectory, { recursive: true });
  }

  const links = [
    ...(await findSkillLinks()).map((link) => ({ ...link, directory: skillsDirectory })),
    ...(await findSteeringLinks()).map((link) => ({ ...link, directory: steeringDirectory })),
  ];
  const results = [];

  for (const link of links) {
    results.push(await ensureSymlink(link.directory, link.linkName, link.target));
  }

  const counts = results.reduce((summary, result) => {
    summary[result.action] = (summary[result.action] ?? 0) + 1;
    return summary;
  }, {});

  for (const result of results.filter(
    (item) => !["unchanged", "preserved"].includes(item.action),
  )) {
    process.stdout.write(
      `${result.action}: ${path.relative(repositoryRoot, result.linkPath)} -> ${result.target}\n`,
    );
  }

  process.stdout.write(
    `AI Central links checked: ${results.length} ` +
      `(created ${counts.created ?? 0}, updated ${counts.updated ?? 0}, ` +
      `unchanged ${counts.unchanged ?? 0}, preserved ${counts.preserved ?? 0})\n`,
  );
}

if (isEntrypoint) {
  await main();
}
