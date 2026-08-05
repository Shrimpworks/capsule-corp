import { readFile } from "node:fs/promises";
import { resolve } from "node:path";
import { fileURLToPath, pathToFileURL } from "node:url";
import Ajv2020 from "ajv/dist/2020.js";

const repositoryRoot = fileURLToPath(new URL("../", import.meta.url));
const defaultManifestPath = "schemas/authority/field-authority-manifest.json";
const defaultSchemaPath = "schemas/authority/field-authority-manifest.schema.json";

export async function verifyFieldAuthorityManifest(options = {}) {
  const rootDirectory = normalizeRoot(options.rootDirectory ?? repositoryRoot);
  const schema =
    options.schema ??
    JSON.parse(
      await readFile(resolve(rootDirectory, options.schemaPath ?? defaultSchemaPath), "utf8"),
    );
  const manifest =
    options.manifest ??
    JSON.parse(
      await readFile(resolve(rootDirectory, options.manifestPath ?? defaultManifestPath), "utf8"),
    );

  const ajv = new Ajv2020({ allErrors: true, strict: true });
  const validate = ajv.compile(schema);
  if (!validate(manifest)) {
    throw new Error(`field-authority manifest schema failure: ${ajv.errorsText(validate.errors)}`);
  }

  const targetIdentities = new Set();
  const usedProfiles = new Set();
  let fieldCount = 0;

  for (const target of manifest.targets) {
    const identity = `${target.object}@${target.version}`;
    if (targetIdentities.has(identity)) {
      throw new Error(`duplicate field-authority target: ${identity}`);
    }
    targetIdentities.add(identity);

    const paths = new Set();
    const sourceFields = new Set();
    for (const field of target.fields) {
      if (paths.has(field.path)) {
        throw new Error(`duplicate classified path for ${identity}: ${field.path}`);
      }
      paths.add(field.path);
      const sourceIdentity = String(field.sourceField);
      if (sourceFields.has(sourceIdentity)) {
        throw new Error(`duplicate canonical source field for ${identity}: ${sourceIdentity}`);
      }
      sourceFields.add(sourceIdentity);
      if (!(field.profile in manifest.profiles)) {
        throw new Error(
          `unknown classification profile for ${identity}${field.path}: ${field.profile}`,
        );
      }
      if (target.object === "capsule.governed-deno-core-c2b-passive-binding") {
        const profile = manifest.profiles[field.profile];
        const expectedConsumers =
          target.version === 2
            ? [
                "governed-runtime-c2b-v2-passive-go-validator",
                "governed-runtime-c2b-v2-passive-typescript-validator",
                "separately-authorized-composed-profile-owned-guest-task-after-gate",
              ]
            : [
                "governed-runtime-c2b-passive-go-validator",
                "governed-runtime-c2b-passive-typescript-validator",
                "separately-authorized-composed-profile-owned-guest-task-after-gate",
              ];
        if (
          profile.retainer !== "capsule-repository-conformance-fixture" ||
          JSON.stringify(profile.allowedConsumers) !== JSON.stringify(expectedConsumers)
        ) {
          throw new Error(
            `incomplete retention or consumer classification for ${identity}${field.path}`,
          );
        }
      }
      usedProfiles.add(field.profile);
      fieldCount += 1;
    }

    const canonical = await inspectCanonicalTarget(rootDirectory, target);
    if (canonical.object !== target.object || canonical.version !== target.version) {
      throw new Error(
        `stale object version or identity for ${identity}: canonical definition is ${canonical.object}@${canonical.version}`,
      );
    }
    assertExactCoverage(identity, canonical.fields, sourceFields);
  }

  const unusedProfiles = Object.keys(manifest.profiles).filter(
    (profile) => !usedProfiles.has(profile),
  );
  if (unusedProfiles.length > 0) {
    throw new Error(`unused classification profiles: ${unusedProfiles.join(", ")}`);
  }

  return {
    fieldCount,
    profileCount: Object.keys(manifest.profiles).length,
    targetCount: manifest.targets.length,
  };
}

async function inspectCanonicalTarget(rootDirectory, target) {
  const definitionPath = resolveDefinitionPath(rootDirectory, target.definition.path);
  const source = await readFile(definitionPath, "utf8");
  switch (target.definition.kind) {
    case "json-schema":
      return inspectJsonSchema(source);
    case "cddl-map":
      return inspectCddlMap(source, target.definition.rule);
    case "go-struct":
      return inspectGoStruct(source, target.definition.type);
    default:
      throw new Error(`unsupported canonical definition kind: ${target.definition.kind}`);
  }
}

function inspectJsonSchema(source) {
  const schema = JSON.parse(source);
  const versionMatch = /:v(\d+)$/u.exec(schema.$id ?? "");
  if (!versionMatch) {
    throw new Error("canonical JSON Schema $id does not end in a numeric version");
  }
  const titleMatch = /\b([A-Z][A-Za-z0-9]+) v\d+\b/u.exec(schema.title ?? "");
  if (!titleMatch) {
    throw new Error("canonical JSON Schema title does not identify its object");
  }
  const fields = new Set();
  collectJsonSchemaFields(schema, schema, "", fields, new Set());
  return {
    object: `capsule.${toKebabCase(titleMatch[1])}`,
    version: Number(versionMatch[1]),
    fields,
  };
}

function collectJsonSchemaFields(root, node, basePath, fields, activeRefs) {
  const resolved = resolveLocalSchemaRef(root, node, activeRefs);
  if (!resolved) {
    return;
  }
  if (resolved.properties) {
    for (const [name, child] of Object.entries(resolved.properties)) {
      const path = `${basePath}/${name}`;
      fields.add(path);
      collectJsonSchemaFields(root, child, path, fields, new Set(activeRefs));
    }
  }
  if (resolved.items) {
    collectJsonSchemaFields(root, resolved.items, `${basePath}/*`, fields, new Set(activeRefs));
  }
  if (
    resolved.additionalProperties &&
    typeof resolved.additionalProperties === "object" &&
    basePath.length > 0
  ) {
    fields.add(`${basePath}/*`);
  }
}

function resolveLocalSchemaRef(root, node, activeRefs) {
  if (!node || typeof node !== "object" || Array.isArray(node) || !node.$ref) {
    return node;
  }
  if (!node.$ref.startsWith("#/") || activeRefs.has(node.$ref)) {
    return null;
  }
  activeRefs.add(node.$ref);
  return node.$ref
    .slice(2)
    .split("/")
    .map((part) => part.replaceAll("~1", "/").replaceAll("~0", "~"))
    .reduce((value, part) => value?.[part], root);
}

function inspectCddlMap(source, rule) {
  const rulePattern = new RegExp(`^${escapeRegExp(rule)}\\s*=\\s*\\{`, "mu");
  const start = rulePattern.exec(source);
  if (!start) {
    throw new Error(`canonical CDDL rule is absent: ${rule}`);
  }
  const bodyStart = source.indexOf("{", start.index) + 1;
  const bodyEnd = findClosingBrace(source, bodyStart);
  const body = source.slice(bodyStart, bodyEnd);
  const entries = [...body.matchAll(/^\s*(\d+)\s*:\s*([^\n]+)$/gmu)];
  if (entries.length === 0) {
    throw new Error(`canonical CDDL rule has no numbered fields: ${rule}`);
  }
  const topLevelFields = entries.map((entry) => entry[1]);
  const nestedFields = [...body.matchAll(/^\s*;\s*field-authority-nested:\s*([^\s]+)\s*$/gmu)].map(
    (entry) => entry[1],
  );
  const fields = new Set([...topLevelFields, ...nestedFields]);
  if (fields.size !== topLevelFields.length + nestedFields.length) {
    throw new Error(`canonical CDDL rule has duplicate field labels: ${rule}`);
  }
  const objectEntry = entries.find((entry) => entry[1] === "1");
  const versionEntry = entries.find((entry) => entry[1] === "2");
  const objectMatch = /"([^"]+)"/u.exec(objectEntry?.[2] ?? "");
  const versionMatch = /^\s*(\d+)\s*(?:,|;|$)/u.exec(versionEntry?.[2] ?? "");
  if (!objectMatch || !versionMatch) {
    throw new Error(`canonical CDDL rule lacks fixed object type/version fields: ${rule}`);
  }
  return { object: objectMatch[1], version: Number(versionMatch[1]), fields };
}

function inspectGoStruct(source, typeName) {
  const typePattern = new RegExp(`^type ${escapeRegExp(typeName)} struct \\{`, "mu");
  const typeMatch = typePattern.exec(source);
  if (!typeMatch) {
    throw new Error(`canonical Go struct is absent: ${typeName}`);
  }
  const markers = [
    ...source
      .slice(0, typeMatch.index)
      .matchAll(/^\/\/ field-authority-object: ([a-z0-9.-]+) v(\d+)$/gmu),
  ];
  const marker = markers.at(-1);
  if (!marker) {
    throw new Error(
      `canonical Go struct lacks a field-authority object/version marker: ${typeName}`,
    );
  }
  const structStart = source.indexOf("{", typeMatch.index) + 1;
  const structEnd = source.indexOf("\n}", structStart);
  if (structEnd < 0) {
    throw new Error(`canonical Go struct is unterminated: ${typeName}`);
  }
  const fields = new Set();
  for (const line of source.slice(structStart, structEnd).split("\n")) {
    const field = /^\s*([A-Za-z_][A-Za-z0-9_]*)\s+[^/]/u.exec(line);
    if (field) {
      fields.add(field[1]);
    }
  }
  if (fields.size === 0) {
    throw new Error(`canonical Go struct has no fields: ${typeName}`);
  }
  return { object: marker[1], version: Number(marker[2]), fields };
}

function assertExactCoverage(identity, canonicalFields, classifiedFields) {
  const missing = [...canonicalFields].filter((field) => !classifiedFields.has(field));
  if (missing.length > 0) {
    throw new Error(`missing field classifications for ${identity}: ${missing.join(", ")}`);
  }
  const absent = [...classifiedFields].filter((field) => !canonicalFields.has(field));
  if (absent.length > 0) {
    throw new Error(
      `classified fields absent from canonical target ${identity}: ${absent.join(", ")}`,
    );
  }
}

function findClosingBrace(source, bodyStart) {
  let depth = 1;
  for (let index = bodyStart; index < source.length; index += 1) {
    if (source[index] === "{") depth += 1;
    if (source[index] === "}") depth -= 1;
    if (depth === 0) return index;
  }
  throw new Error("canonical CDDL map is unterminated");
}

function resolveDefinitionPath(rootDirectory, relativePath) {
  const absolute = resolve(rootDirectory, relativePath);
  const root = resolve(rootDirectory);
  if (absolute !== root && !absolute.startsWith(`${root}/`)) {
    throw new Error(`canonical definition escapes repository root: ${relativePath}`);
  }
  return absolute;
}

function normalizeRoot(value) {
  const path = value instanceof URL ? fileURLToPath(value) : value;
  return resolve(path);
}

function toKebabCase(value) {
  return value.replace(/([a-z0-9])([A-Z])/gu, "$1-$2").toLowerCase();
}

function escapeRegExp(value) {
  return value.replace(/[.*+?^${}()|[\]\\]/gu, "\\$&");
}

if (process.argv[1] && import.meta.url === pathToFileURL(process.argv[1]).href) {
  const result = await verifyFieldAuthorityManifest();
  process.stdout.write(
    `validated ${result.fieldCount} authority field classifications across ${result.targetCount} passive targets with ${result.profileCount} closed profiles\n`,
  );
}
