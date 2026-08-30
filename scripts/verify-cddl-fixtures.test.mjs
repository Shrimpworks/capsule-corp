// verify-cddl-fixtures.mjs is a top-level, unparameterized script: it always
// reads schemas/fixtures/*.json relative to its own file location and has no
// exported function to call with substitute fixtures. To exercise a
// deliberate mutation without changing that production script, this test
// copies it unmodified into an isolated temp directory alongside the real
// fixtures (preserving the "../schemas/fixtures/" relative layout the
// script depends on), mutates one field, and spawns the script there as a
// subprocess.
import assert from "node:assert/strict";
import { execFile } from "node:child_process";
import { cp, mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);
const scriptPath = fileURLToPath(new URL("./verify-cddl-fixtures.mjs", import.meta.url));
const fixturesDir = fileURLToPath(new URL("../schemas/fixtures/", import.meta.url));

test("verifies the checked-in CDDL fixtures cleanly", async () => {
  const { stdout } = await execFileAsync(process.execPath, [scriptPath]);
  assert.match(stdout, /validated .*approval-grant-v0\.json/u);
  assert.match(stdout, /validated .*plan-registration-v0\.json/u);
  assert.match(stdout, /validated .*execution-plan-v0\.json/u);
  assert.match(stdout, /rejected invalid internal candidate ID/u);
});

test("catches a mutated ApprovalGrant installationIdHex the verifier must reject", async () => {
  await withMutatedFixture(
    "approval-grant-v0.json",
    (fixture) => {
      fixture.installationIdHex = flipHexNibble(fixture.installationIdHex);
    },
    /ApprovalGrant payload bytes: got/u,
  );
});

test("catches a mutated PlanRegistration objectVersion the verifier must reject", async () => {
  await withMutatedFixture(
    "plan-registration-v0.json",
    (fixture) => {
      fixture.objectVersion = 1;
    },
    /object version: got 1, want 0/u,
  );
});

test("catches a mutated ExecutionPlan wallTimeOrigin the verifier must reject", async () => {
  await withMutatedFixture(
    "execution-plan-v0.json",
    (fixture) => {
      fixture.wallTimeOrigin = "invented-origin";
    },
    /wallTimeOrigin is unsupported/u,
  );
});

/**
 * Copies the whole fixtures directory and the unmodified verifier script
 * into an isolated temp directory, applies `mutate` to the named fixture,
 * and asserts the script rejects it with `expectedMessage`.
 */
async function withMutatedFixture(fixtureName, mutate, expectedMessage) {
  const dir = await mkdtemp(join(tmpdir(), "verify-cddl-fixtures-"));
  try {
    const scriptsDir = join(dir, "scripts");
    const fixturesOut = join(dir, "schemas", "fixtures");
    await mkdir(scriptsDir, { recursive: true });
    await mkdir(fixturesOut, { recursive: true });
    await cp(scriptPath, join(scriptsDir, "verify-cddl-fixtures.mjs"));
    await cp(fixturesDir, fixturesOut, { recursive: true });

    const fixturePath = join(fixturesOut, fixtureName);
    const fixture = JSON.parse(await readFile(fixturePath, "utf8"));
    mutate(fixture);
    await writeFile(fixturePath, `${JSON.stringify(fixture, null, 2)}\n`);

    await assert.rejects(
      () => execFileAsync(process.execPath, [join(scriptsDir, "verify-cddl-fixtures.mjs")]),
      expectedMessage,
    );
  } finally {
    await rm(dir, { recursive: true, force: true });
  }
}

function flipHexNibble(hex) {
  const first = hex[0] === "0" ? "1" : "0";
  return first + hex.slice(1);
}
