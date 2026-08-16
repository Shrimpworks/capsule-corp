import assert from "node:assert/strict";
import test from "node:test";

import { validateControlEvidenceMatrix } from "./verify-control-evidence-matrix.mjs";

const header = `| Claim ID | Claim | Threat | Owner | Mechanism | Live verification | Required attack tests | Receipt/transcript evidence | Known limitation | Status |
| --- | --- | --- | --- | --- | --- | --- | --- | --- | --- |`;

function row(id = "TEST-001", status = "proposed") {
  return `| ${id} | claim | threat | owner | mechanism | verification | tests | evidence | limitation | ${status} |`;
}

test("accepts unique canonical rows and provenance that does not override them", () => {
  const text = `${header}\n${row()}\n\n## Provenance\nHistorical evidence remains linked.\n`;
  assert.equal(validateControlEvidenceMatrix(text), 1);
});

test("treats a pipe inside inline code as cell content", () => {
  const coded = row().replace("mechanism", "`userPresence|privateKeyUsage`");
  assert.equal(validateControlEvidenceMatrix(`${header}\n${coded}\n`), 1);
});

test("rejects prose that supersedes canonical row text", () => {
  const text = `${header}\n${row()}\n\nThis supersedes the older TEST-001 row text.\n`;
  assert.throws(() => validateControlEvidenceMatrix(text), /must not supersede/);
});

test("rejects duplicate claim IDs", () => {
  const text = `${header}\n${row()}\n${row()}\n`;
  assert.throws(() => validateControlEvidenceMatrix(text), /duplicate.*TEST-001/);
});

test("rejects unknown evidence states", () => {
  const text = `${header}\n${row("TEST-001", "passed")}\n`;
  assert.throws(() => validateControlEvidenceMatrix(text), /unsupported.*passed/);
});
