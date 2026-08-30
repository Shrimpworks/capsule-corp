// verify-typed-guest-transport-fixtures.mjs is a top-level, unparameterized
// script: it always reads schemas/conformance/typed-guest-transport-v1/
// relative to its own file location and has no exported function to call
// with substitute fixtures. To exercise a deliberate mutation without
// changing that production script, this test copies it (and the
// scripts/lib helper it imports) unmodified into an isolated temp directory
// alongside the real conformance fixtures (preserving the relative layout
// the script depends on), mutates one fixture, and spawns the script there
// as a subprocess.
import assert from "node:assert/strict";
import { execFile } from "node:child_process";
import { cp, mkdir, mkdtemp, readFile, rm, writeFile } from "node:fs/promises";
import { tmpdir } from "node:os";
import { join } from "node:path";
import test from "node:test";
import { fileURLToPath } from "node:url";
import { promisify } from "node:util";

const execFileAsync = promisify(execFile);
const scriptPath = fileURLToPath(
  new URL("./verify-typed-guest-transport-fixtures.mjs", import.meta.url),
);
const libDir = fileURLToPath(new URL("./lib/", import.meta.url));
const fixturesDir = fileURLToPath(
  new URL("../schemas/conformance/typed-guest-transport-v1/", import.meta.url),
);

test("verifies the checked-in typed-guest-transport conformance fixtures cleanly", async () => {
  await execFileAsync(process.execPath, [scriptPath]);
});

test("catches a byte-corrupted accepted fixture the verifier must reject", async () => {
  await withFixtures(async (fixturesOut) => {
    const path = join(fixturesOut, "accept-source-ordinary.bin");
    const bytes = await readFile(path);
    bytes[0] ^= 0xff;
    await writeFile(path, bytes);
  }, /accept-source-ordinary/u);
});

test("catches a manifest state-case disposition swap the verifier must reject", async () => {
  await withFixtures(async (fixturesOut) => {
    const manifestPath = join(fixturesOut, "manifest.json");
    const manifest = JSON.parse(await readFile(manifestPath, "utf8"));
    const target = manifest.stateCases.find(({ id }) => id === "cancel-before-g");
    assert.ok(target, "fixture setup: expected cancel-before-g state case to exist");
    target.disposition = "REFUSED_LIFECYCLE";
    await writeFile(manifestPath, `${JSON.stringify(manifest, null, 2)}\n`);
  }, /AssertionError/u);
});

/**
 * Copies the fixtures directory, scripts/lib, and the unmodified verifier
 * script into an isolated temp directory, applies `mutate`, and asserts the
 * script rejects with `expectedMessage`.
 */
async function withFixtures(mutate, expectedMessage) {
  const dir = await mkdtemp(join(tmpdir(), "verify-typed-guest-transport-"));
  try {
    const scriptsDir = join(dir, "scripts");
    const fixturesOut = join(dir, "schemas", "conformance", "typed-guest-transport-v1");
    await mkdir(scriptsDir, { recursive: true });
    await mkdir(fixturesOut, { recursive: true });
    await cp(scriptPath, join(scriptsDir, "verify-typed-guest-transport-fixtures.mjs"));
    await cp(libDir, join(scriptsDir, "lib"), { recursive: true });
    await cp(fixturesDir, fixturesOut, { recursive: true });

    await mutate(fixturesOut);

    await assert.rejects(
      () =>
        execFileAsync(process.execPath, [
          join(scriptsDir, "verify-typed-guest-transport-fixtures.mjs"),
        ]),
      expectedMessage,
    );
  } finally {
    await rm(dir, { recursive: true, force: true });
  }
}
