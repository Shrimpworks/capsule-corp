// Enumerates the .ts files that fall under the exported-JSDoc documentation
// rule (docs/JAVASCRIPT_TYPESCRIPT_ENGINEERING_STANDARDS.md, issue #217):
// every `packages/*/src/**/*.ts` file, excluding tests and declaration
// files. Used by verify-exported-jsdoc.mjs's full-repository audit mode.
import { readdir } from "node:fs/promises";
import { join } from "node:path";
import { fileURLToPath } from "node:url";

const packagesRoot = fileURLToPath(new URL("../../packages/", import.meta.url));

async function walk(directory) {
  const entries = await readdir(directory, { withFileTypes: true });
  const files = [];
  for (const entry of entries) {
    if (entry.name === "node_modules" || entry.name === "dist") {
      continue;
    }
    const path = join(directory, entry.name);
    if (entry.isDirectory()) {
      files.push(...(await walk(path)));
      continue;
    }
    if (
      (path.endsWith(".ts") || path.endsWith(".tsx")) &&
      !path.endsWith(".test.ts") &&
      !path.endsWith(".d.ts")
    ) {
      files.push(path);
    }
  }
  return files;
}

/** Every in-scope `packages/*\/src/**\/*.ts` file, as absolute paths. */
export async function listInScopeFiles() {
  const packageNames = await readdir(packagesRoot, { withFileTypes: true });
  const files = [];
  for (const entry of packageNames) {
    if (!entry.isDirectory()) {
      continue;
    }
    const srcDirectory = join(packagesRoot, entry.name, "src");
    try {
      files.push(...(await walk(srcDirectory)));
    } catch (error) {
      if (error.code !== "ENOENT") {
        throw error;
      }
    }
  }
  return files.sort();
}
