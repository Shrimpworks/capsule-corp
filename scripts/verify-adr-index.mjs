import { readdir, readFile } from "node:fs/promises";
import { fileURLToPath } from "node:url";

const adrDirectory = new URL("../docs/adr/", import.meta.url);
const entries = await readdir(adrDirectory);

// Anything named like an ADR is an ADR for coverage purposes. This is
// deliberately looser than adrFilenamePattern below so that a numbered file
// the strict pattern cannot parse fails the run instead of being filtered
// out of every check that follows.
const numberedFilenamePattern = /^\d{4}-.*\.md$/;
// The strict naming rule. The slug admits "." because ADR numbers can appear
// in a decision's own title (0039-license-capsule-under-apache-2.0.md).
const adrFilenamePattern = /^(\d{4})-[a-z0-9.-]+\.md$/;

const adrFiles = entries.filter((entry) => numberedFilenamePattern.test(entry)).sort();

if (adrFiles.length === 0) {
  throw new Error("No ADR files were found under docs/adr/");
}

const unparsableFilenames = adrFiles.filter((filename) => !adrFilenamePattern.test(filename));
if (unparsableFilenames.length > 0) {
  throw new Error(
    `ADR filenames do not match ${adrFilenamePattern}:\n${unparsableFilenames
      .map((filename) => `  ${filename}`)
      .join("\n")}`,
  );
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
