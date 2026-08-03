import { readdir, readFile } from "node:fs/promises";
import { fileURLToPath } from "node:url";

const adrDirectory = new URL("../docs/adr/", import.meta.url);
const entries = await readdir(adrDirectory);

const adrFilenamePattern = /^(\d{4})-[a-z0-9-]+\.md$/;
const adrFiles = entries.filter((entry) => adrFilenamePattern.test(entry)).sort();

if (adrFiles.length === 0) {
  throw new Error("No ADR files were found under docs/adr/");
}

const numbersByFile = new Map(
  adrFiles.map((filename) => [filename, adrFilenamePattern.exec(filename)[1]]),
);

const filesByNumber = new Map();
for (const [filename, number] of numbersByFile) {
  if (!filesByNumber.has(number)) {
    filesByNumber.set(number, []);
  }
  filesByNumber.get(number).push(filename);
}

const duplicateNumbers = [...filesByNumber.entries()].filter(([, files]) => files.length > 1);
if (duplicateNumbers.length > 0) {
  const detail = duplicateNumbers
    .map(([number, files]) => `  ADR-${number}: ${files.join(", ")}`)
    .join("\n");
  throw new Error(`Duplicate ADR numbers found:\n${detail}`);
}

const readmePath = fileURLToPath(new URL("README.md", adrDirectory));
const readme = await readFile(readmePath, "utf8");

const missingFromIndex = adrFiles.filter((filename) => !readme.includes(`(${filename})`));
if (missingFromIndex.length > 0) {
  throw new Error(
    `docs/adr/README.md is missing an index entry for:\n${missingFromIndex.map((filename) => `  ${filename}`).join("\n")}`,
  );
}

const indexedButMissingFile = [...readme.matchAll(/\((\d{4}-[a-z0-9-]+\.md)\)/g)]
  .map((match) => match[1])
  .filter((filename) => !numbersByFile.has(filename));
if (indexedButMissingFile.length > 0) {
  throw new Error(
    `docs/adr/README.md links to a file that does not exist:\n${indexedButMissingFile.map((filename) => `  ${filename}`).join("\n")}`,
  );
}

process.stdout.write(`validated ${adrFiles.length} ADR files against docs/adr/README.md\n`);
