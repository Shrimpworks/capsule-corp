// verify-adr-index.mjs is a top-level, unparameterized script: it always
// reads docs/adr/ relative to its own file location and has no exported
// function to call with substitute fixtures. To exercise a deliberate
// mutation without changing that production script, this test copies it
// unmodified into an isolated temp directory alongside a small synthetic
// docs/adr/ fixture (preserving the "../docs/adr/" relative layout the
// script depends on) and spawns it there as a subprocess.
import assert from "node:assert/strict";
import { execFile } from "node:child_process";
import { cp, mkdir, mkdtemp, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);
const scriptPath = fileURLToPath(new URL("./verify-adr-index.mjs", import.meta.url));

test("verifies the checked-in ADR index cleanly", async () => {
  const { stdout } = await execFileAsync(process.execPath, [scriptPath]);
  assert.match(stdout, /^validated \d+ ADR files against docs\/adr\/README\.md$/mu);
});

test("catches a duplicate ADR number the verifier must reject", async () => {
  await withFixture(async (adrDir) => {
    await writeFile(join(adrDir, "0001-first.md"), "# ADR 1\n");
    await writeFile(join(adrDir, "0001-second.md"), "# ADR 1 duplicate\n");
    await writeFile(
      join(adrDir, "README.md"),
      "## Index\n\n- [First](0001-first.md)\n- [Second](0001-second.md)\n",
    );
  }, /Duplicate ADR numbers found/u);
});

test("catches an ADR missing from the README index the verifier must reject", async () => {
  await withFixture(async (adrDir) => {
    await writeFile(join(adrDir, "0001-first.md"), "# ADR 1\n");
    await writeFile(join(adrDir, "0002-second.md"), "# ADR 2\n");
    await writeFile(join(adrDir, "README.md"), "## Index\n\n- [First](0001-first.md)\n");
  }, /is missing an index entry for/u);
});

test("catches a README index entry with no backing file the verifier must reject", async () => {
  await withFixture(async (adrDir) => {
    await writeFile(join(adrDir, "0001-first.md"), "# ADR 1\n");
    await writeFile(
      join(adrDir, "README.md"),
      "## Index\n\n- [First](0001-first.md)\n- [Missing](0002-missing.md)\n",
    );
  }, /links to a file that does not exist/u);
});

/**
 * Runs `seed` against an isolated `docs/adr/` fixture backing an unmodified
 * copy of verify-adr-index.mjs, then asserts the script rejects with
 * `expectedMessage`.
 */
async function withFixture(seed, expectedMessage) {
  const dir = await mkdtemp(join(tmpdir(), "verify-adr-index-"));
  try {
    const scriptsDir = join(dir, "scripts");
    const adrDir = join(dir, "docs", "adr");
    await mkdir(scriptsDir, { recursive: true });
    await mkdir(adrDir, { recursive: true });
    await cp(scriptPath, join(scriptsDir, "verify-adr-index.mjs"));
    await seed(adrDir);

    await assert.rejects(
      () => execFileAsync(process.execPath, [join(scriptsDir, "verify-adr-index.mjs")]),
      expectedMessage,
    );
  } finally {
    await rm(dir, { recursive: true, force: true });
  }
}
