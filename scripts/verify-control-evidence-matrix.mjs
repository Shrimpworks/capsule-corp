import { readFile } from "node:fs/promises";
import { resolve } from "node:path";
import { pathToFileURL } from "node:url";

const matrixUrl = new URL("../docs/security/CONTROL_EVIDENCE_MATRIX.md", import.meta.url);
const allowedStatuses = new Set([
  "proposed",
  "local-mechanic",
  "spike-observed",
  "implemented",
  "validated",
  "degraded",
  "unsupported",
]);

export function validateControlEvidenceMatrix(text) {
  const lines = text.split("\n");
  const headerIndex = lines.findIndex((line) => line.startsWith("| Claim ID |"));
  if (headerIndex === -1) {
    throw new Error("control-evidence matrix table header is missing");
  }

  const rows = [];
  for (const line of lines.slice(headerIndex + 2)) {
    if (!line.startsWith("|")) {
      break;
    }
    const cells = parseMarkdownTableRow(line);
    if (cells.length !== 10) {
      throw new Error(`control-evidence row must have 10 columns: ${line}`);
    }
    rows.push(cells);
  }

  if (rows.length === 0) {
    throw new Error("control-evidence matrix has no claim rows");
  }

  const seen = new Set();
  for (const row of rows) {
    const claimId = row[0];
    if (seen.has(claimId)) {
      throw new Error(`duplicate control-evidence claim ID: ${claimId}`);
    }
    seen.add(claimId);
    if (!allowedStatuses.has(row[9])) {
      throw new Error(`unsupported control-evidence status for ${claimId}: ${row[9]}`);
    }
  }

  const prose = lines.slice(headerIndex + 2 + rows.length).join("\n");
  if (/supersed(?:e|es|ed|ing)[^\n]{0,120}row text/iu.test(prose)) {
    throw new Error(
      "prose must not supersede canonical control-evidence row text; update the row in place",
    );
  }

  return rows.length;
}

function parseMarkdownTableRow(line) {
  const cells = [];
  let cell = "";
  let inCode = false;
  let escaped = false;

  for (const character of line.slice(1, -1)) {
    if (escaped) {
      cell += character;
      escaped = false;
      continue;
    }
    if (character === "\\") {
      cell += character;
      escaped = true;
      continue;
    }
    if (character === "`") {
      inCode = !inCode;
      cell += character;
      continue;
    }
    if (character === "|" && !inCode) {
      cells.push(cell.trim());
      cell = "";
      continue;
    }
    cell += character;
  }
  cells.push(cell.trim());
  return cells;
}

async function main() {
  const text = await readFile(matrixUrl, "utf8");
  const count = validateControlEvidenceMatrix(text);
  process.stdout.write(`validated ${count} canonical control-evidence rows\n`);
}

if (process.argv[1] && pathToFileURL(resolve(process.argv[1])).href === import.meta.url) {
  await main();
}
