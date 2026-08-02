import { createHash } from "node:crypto";
import { lstat, readdir, readFile } from "node:fs/promises";
import { isAbsolute, relative, resolve, sep } from "node:path";
import { pathToFileURL } from "node:url";
import Ajv2020 from "ajv/dist/2020.js";

const defaultRootDirectory = new URL("../schemas/conformance/v0/", import.meta.url);
const metadataFiles = new Set(["manifest.json", "manifest.schema.json"]);

export async function verifyConformanceCorpus({ rootDirectory = defaultRootDirectory } = {}) {
  const root = resolveDirectory(rootDirectory);
  const schema = await readJson(resolve(root, "manifest.schema.json"), "manifest schema");
  const manifest = await readJson(resolve(root, "manifest.json"), "manifest");
  const ajv = new Ajv2020({ allErrors: true, strict: true });
  const validate = ajv.compile(schema);

  if (!validate(manifest)) {
    throw new Error(`manifest schema validation failed: ${ajv.errorsText(validate.errors)}`);
  }

  assertUnique(manifest.rules, "rule");
  assertUnique(manifest.cases, "case");
  const rules = new Map(manifest.rules.map((rule) => [rule.id, rule]));
  const listedFixtures = new Set();

  for (const entry of manifest.cases) {
    for (const ruleId of entry.ruleIds) {
      if (!rules.has(ruleId)) {
        throw new Error(`case ${entry.id} references unknown rule ${ruleId}`);
      }
    }

    if (entry.expected.decision === "reject" && entry.expected.authorityStateChanged) {
      throw new Error(`rejected case ${entry.id} cannot change authority-bearing state`);
    }

    assertScalarRoleContext(entry);
    await verifyFixture(root, entry.fixture, entry.id);
    listedFixtures.add(entry.fixture.path);

    for (const [label, fixture] of contextFixtures(entry.context)) {
      await verifyFixture(root, fixture, `${entry.id} ${label}`);
      listedFixtures.add(fixture.path);
    }
    assertProposalResolutionContext(entry);
  }

  for (const rule of manifest.rules) {
    for (const requirement of rule.requiredCases) {
      const covered = manifest.cases.some(
        (entry) =>
          entry.ruleIds.includes(rule.id) &&
          entry.expected.decision === requirement.decision &&
          entry.variant === requirement.variant,
      );
      if (!covered) {
        throw new Error(
          `required case is not covered: ${rule.id} ${requirement.decision}/${requirement.variant}`,
        );
      }
    }
  }

  const corpusFiles = await listCorpusFiles(root);
  for (const path of corpusFiles) {
    if (!listedFixtures.has(path)) {
      throw new Error(`fixture is not listed in the manifest: ${path}`);
    }
  }

  return {
    caseCount: manifest.cases.length,
    fixtureCount: listedFixtures.size,
    ruleCount: manifest.rules.length,
  };
}

async function verifyFixture(root, fixture, owner) {
  const fixturePath = resolveFixturePath(root, fixture.path);
  let metadata;
  try {
    metadata = await lstat(fixturePath);
  } catch (error) {
    if (error?.code === "ENOENT") {
      throw new Error(`listed fixture does not exist for ${owner}: ${fixture.path}`);
    }
    throw error;
  }
  if (!metadata.isFile() || metadata.isSymbolicLink()) {
    throw new Error(`listed fixture is not a regular file for ${owner}: ${fixture.path}`);
  }

  const bytes = await readFile(fixturePath);
  if (bytes.length !== fixture.byteLength) {
    throw new Error(
      `fixture byte length mismatch for ${owner}: expected ${fixture.byteLength}, got ${bytes.length}`,
    );
  }
  const digest = createHash("sha256").update(bytes).digest("hex");
  if (digest !== fixture.sha256) {
    throw new Error(
      `fixture SHA-256 mismatch for ${owner}: expected ${fixture.sha256}, got ${digest}`,
    );
  }
}

function resolveFixturePath(root, path) {
  if (isAbsolute(path) || path.includes("\\")) {
    throw new Error(`fixture path is not a safe relative path: ${path}`);
  }
  const segments = path.split("/");
  if (segments.some((segment) => segment === "" || segment === "." || segment === "..")) {
    throw new Error(`fixture path is not a safe relative path: ${path}`);
  }
  const resolved = resolve(root, path);
  if (!resolved.startsWith(`${root}${sep}`)) {
    throw new Error(`fixture path escapes the corpus root: ${path}`);
  }
  return resolved;
}

async function listCorpusFiles(root) {
  const files = [];
  await visit(root);
  return files.sort();

  async function visit(directory) {
    const entries = await readdir(directory, { withFileTypes: true });
    for (const entry of entries) {
      const absolutePath = resolve(directory, entry.name);
      const corpusPath = relative(root, absolutePath).split(sep).join("/");
      if (entry.isSymbolicLink()) {
        throw new Error(`corpus contains a symbolic link: ${corpusPath}`);
      }
      if (entry.isDirectory()) {
        await visit(absolutePath);
      } else if (entry.isFile() && !metadataFiles.has(corpusPath)) {
        files.push(corpusPath);
      } else if (!entry.isFile()) {
        throw new Error(`corpus contains a non-regular entry: ${corpusPath}`);
      }
    }
  }
}

function assertUnique(entries, kind) {
  const ids = new Set();
  for (const entry of entries) {
    if (ids.has(entry.id)) {
      throw new Error(`duplicate ${kind} ID: ${entry.id}`);
    }
    ids.add(entry.id);
  }
}

function assertScalarRoleContext(entry) {
  if (entry.context.kind !== "scalar-role") {
    return;
  }
  const sameRole = entry.context.expectedRole === entry.context.providedRole;
  if (entry.expected.decision === "accept" && !sameRole) {
    throw new Error(`accepted scalar-role case ${entry.id} must use the expected role`);
  }
  if (
    entry.expected.decision === "reject" &&
    entry.expected.classification === "DOMAIN" &&
    sameRole
  ) {
    throw new Error(`DOMAIN case ${entry.id} must substitute a different scalar role`);
  }
}

function contextFixtures(context) {
  if (context.kind === "fixture") {
    return [["context", context.fixture]];
  }
  if (context.kind !== "proposal-resolution") {
    return [];
  }
  return [
    ["profile registry", context.profileRegistry],
    ["user policy", context.userPolicy],
    ...(context.oracle.sourceManifest
      ? [["source manifest oracle", context.oracle.sourceManifest]]
      : []),
    ...(context.oracle.canonicalInlineInput
      ? [["canonical inline-input oracle", context.oracle.canonicalInlineInput]]
      : []),
  ];
}

function assertProposalResolutionContext(entry) {
  if (entry.context.kind !== "proposal-resolution") {
    return;
  }
  const { oracle } = entry.context;
  for (const [label, fixture, digest] of [
    ["source manifest", oracle.sourceManifest, oracle.sourceManifestDigest],
    ["canonical inline input", oracle.canonicalInlineInput, oracle.inlineInputDigest],
  ]) {
    if ((fixture === null) !== (digest === null)) {
      throw new Error(`${entry.id} ${label} fixture and digest must be present together`);
    }
    if (fixture && fixture.sha256 !== digest) {
      throw new Error(`${entry.id} ${label} digest must match the retained fixture digest`);
    }
  }
  if (
    entry.expected.decision === "reject" &&
    (oracle.sourceManifest || oracle.canonicalInlineInput || oracle.wallTime)
  ) {
    throw new Error(`rejected proposal case ${entry.id} cannot retain a resolution result`);
  }
}

async function readJson(path, label) {
  try {
    return JSON.parse(await readFile(path, "utf8"));
  } catch (error) {
    throw new Error(`failed to read ${label} at ${path}: ${error.message}`, { cause: error });
  }
}

function resolveDirectory(value) {
  return value instanceof URL ? resolve(value.pathname) : resolve(value);
}

if (process.argv[1] && pathToFileURL(resolve(process.argv[1])).href === import.meta.url) {
  const result = await verifyConformanceCorpus();
  process.stdout.write(
    `validated conformance corpus: ${result.ruleCount} rules, ${result.caseCount} cases, ${result.fixtureCount} fixtures\n`,
  );
}
