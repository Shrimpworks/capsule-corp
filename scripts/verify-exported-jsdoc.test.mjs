import assert from "node:assert/strict";
import test from "node:test";

import { findViolations } from "./verify-exported-jsdoc.mjs";

function names(violations) {
  return violations.map((violation) => violation.name).sort();
}

test("flags exported declarations with no leading doc comment", () => {
  const text = `
export function undocumented() {}
export class UndocumentedClass {}
export interface UndocumentedInterface {}
export type UndocumentedType = string;
export enum UndocumentedEnum { A }
export const undocumentedArrow = () => {};
`;
  const violations = findViolations("fixture.ts", text);
  assert.deepEqual(
    names(violations),
    [
      "UndocumentedClass",
      "UndocumentedEnum",
      "UndocumentedInterface",
      "UndocumentedType",
      "undocumentedArrow",
      "undocumented",
    ].sort(),
  );
});

test("accepts a leading /** */ doc comment on every exported form", () => {
  const text = `
/** Docs. */
export function documented() {}
/** Docs. */
export class DocumentedClass {}
/** Docs. */
export interface DocumentedInterface {}
/** Docs. */
export type DocumentedType = string;
/** Docs. */
export enum DocumentedEnum { A }
/** Docs. */
export const documentedArrow = () => {};
/** Docs. */
export const documentedTypedArrow: (x: number) => number = (x) => x;
`;
  assert.deepEqual(findViolations("fixture.ts", text), []);
});

test("rejects a plain // comment as insufficient", () => {
  const text = `
// not a doc comment
export function undocumented() {}
`;
  assert.deepEqual(names(findViolations("fixture.ts", text)), ["undocumented"]);
});

test("ignores unexported declarations, re-exports, and plain data constants", () => {
  const text = `
function internalHelper() {}
export { helper } from "./other.js";
export const PLAIN_CONSTANT = 42;
export default function () {}
`;
  assert.deepEqual(findViolations("fixture.ts", text), []);
});

test("reports the exact 1-based line of the declaration", () => {
  const text = `import { x } from "./x.js";\n\nexport function undocumented() {}\n`;
  const violations = findViolations("fixture.ts", text);
  assert.equal(violations.length, 1);
  assert.equal(violations[0].line, 3);
});
