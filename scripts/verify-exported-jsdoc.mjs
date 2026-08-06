// Verifies that top-level exported types and functions in the given
// TypeScript files carry a `/** ... */` doc comment, per
// docs/JAVASCRIPT_TYPESCRIPT_ENGINEERING_STANDARDS.md's "Every exported type
// and function gets a doc comment" rule (see issue #217).
//
// This is a new-code-only gate, not a full-repository audit: CI invokes it
// with only the .ts files changed in the current pull request (see
// .github/workflows/ci.yml), so it enforces the rule going forward without
// requiring the existing repository-wide backlog to be documented in one
// change. Run with explicit file arguments locally to check specific files;
// run with no arguments to audit every in-scope file and report full counts
// (informational — this mode is not wired into CI because of that backlog).
//
// In scope: exported function/class/interface/type-alias/enum declarations,
// and exported `const`/`let` bindings whose initializer is a function or
// arrow function. Plain exported data constants, re-exports (`export {
// x } from "./y"`), and `export default` are out of scope, matching the
// standard's literal "type and function" wording.

import { readFile } from "node:fs/promises";
import { relative, resolve } from "node:path";
import { pathToFileURL } from "node:url";
import ts from "typescript";

import { listInScopeFiles } from "./lib/exported-jsdoc-files.mjs";

/**
 * @param {ts.SourceFile} sourceFile
 * @param {ts.Node} node
 */
function hasLeadingDocComment(sourceFile, node) {
  const ranges = ts.getLeadingCommentRanges(sourceFile.text, node.getFullStart()) ?? [];
  return ranges.some((range) => sourceFile.text.slice(range.pos, range.pos + 3) === "/**");
}

/** @param {readonly ts.ModifierLike[] | undefined} modifiers */
function isExported(modifiers) {
  return (modifiers ?? []).some((modifier) => modifier.kind === ts.SyntaxKind.ExportKeyword);
}

/** @param {ts.VariableDeclaration} declaration */
function isFunctionValuedBinding(declaration) {
  const initializer = declaration.initializer;
  if (!initializer) {
    return false;
  }
  return (
    ts.isArrowFunction(initializer) ||
    ts.isFunctionExpression(initializer) ||
    // Typed function bindings, e.g. `export const f: (x: number) => number = ...`.
    (declaration.type !== undefined && ts.isFunctionTypeNode(declaration.type))
  );
}

/**
 * @param {string} path
 * @param {string} text
 * @returns {{path: string, name: string, line: number}[]}
 */
export function findViolations(path, text) {
  const sourceFile = ts.createSourceFile(path, text, ts.ScriptTarget.Latest, true);
  const violations = [];

  const record = (node, name) => {
    if (!hasLeadingDocComment(sourceFile, node)) {
      const { line } = sourceFile.getLineAndCharacterOfPosition(node.getStart(sourceFile));
      violations.push({ path, name, line: line + 1 });
    }
  };

  for (const statement of sourceFile.statements) {
    if (
      ts.isFunctionDeclaration(statement) ||
      ts.isClassDeclaration(statement) ||
      ts.isInterfaceDeclaration(statement) ||
      ts.isTypeAliasDeclaration(statement) ||
      ts.isEnumDeclaration(statement)
    ) {
      if (isExported(statement.modifiers) && statement.name) {
        record(statement, statement.name.getText(sourceFile));
      }
      continue;
    }
    if (ts.isVariableStatement(statement) && isExported(statement.modifiers)) {
      for (const declaration of statement.declarationList.declarations) {
        if (ts.isIdentifier(declaration.name) && isFunctionValuedBinding(declaration)) {
          // The doc comment belongs on the `export const ...` statement, not
          // the individual declarator.
          record(statement, declaration.name.getText(sourceFile));
        }
      }
    }
  }

  return violations;
}

async function main() {
  const explicitFiles = process.argv.slice(2);
  const files = explicitFiles.length > 0 ? explicitFiles : await listInScopeFiles();

  if (files.length === 0) {
    console.log("verify-exported-jsdoc: no in-scope .ts files to check");
    return;
  }

  const allViolations = [];
  for (const path of files) {
    if (!path.endsWith(".ts") && !path.endsWith(".tsx")) {
      continue;
    }
    if (path.endsWith(".test.ts") || path.endsWith(".d.ts")) {
      continue;
    }
    const text = await readFile(path, "utf8");
    allViolations.push(...findViolations(path, text));
  }

  if (allViolations.length === 0) {
    console.log(`verify-exported-jsdoc: ${files.length} file(s) checked, no violations`);
    return;
  }

  const detail = allViolations
    .map(
      (violation) =>
        `  ${relative(process.cwd(), violation.path)}:${violation.line}: exported ${violation.name} has no doc comment`,
    )
    .join("\n");
  throw new Error(
    `${allViolations.length} exported type/function(s) missing a doc comment:\n${detail}`,
  );
}

if (process.argv[1] && pathToFileURL(resolve(process.argv[1])).href === import.meta.url) {
  await main();
}
