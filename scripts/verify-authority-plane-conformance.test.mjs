// verify-authority-plane-conformance.mjs is a top-level, unparameterized
// script: it always reads schemas/conformance/authority-plane-v0/ relative
// to its own file location and has no exported function to call with
// substitute fixtures. To exercise a deliberate mutation without changing
// that production script, this test copies it (and the scripts/lib helper
// it imports) unmodified into an isolated temp directory alongside the real
// conformance fixtures (preserving the relative layout the script depends
// on), mutates one fixture, and spawns the script there as a subprocess.
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
  new URL("./verify-authority-plane-conformance.mjs", import.meta.url),
);
const libDir = fileURLToPath(new URL("./lib/", import.meta.url));
const fixturesDir = fileURLToPath(
  new URL("../schemas/conformance/authority-plane-v0/", import.meta.url),
);

test("verifies the checked-in authority-plane conformance fixtures cleanly", async () => {
  const { stdout } = await execFileAsync(process.execPath, [scriptPath]);
  assert.match(stdout, /verified independent TypeScript authority-plane known answers/u);
});

test("catches a byte-corrupted role-bindings.bin the verifier must reject", async () => {
  await withFixtures(async (fixturesOut) => {
    const path = join(fixturesOut, "role-bindings.bin");
    const bytes = await readFile(path);
    bytes[0] ^= 0xff;
    await writeFile(path, bytes);
  }, /role-bindings\.bin/u);
});

test("catches a mutated manifest cap the verifier must reject", async () => {
  await withFixtures(async (fixturesOut) => {
    const manifestPath = join(fixturesOut, "manifest.json");
    const manifest = JSON.parse(await readFile(manifestPath, "utf8"));
    manifest.caps.executionPlan += 1;
    await writeFile(manifestPath, `${JSON.stringify(manifest, null, 2)}\n`);
  }, /AssertionError/u);
});

/**
 * Copies the fixtures directory, scripts/lib, and the unmodified verifier
 * script into an isolated temp directory, applies `mutate`, and asserts the
 * script rejects with `expectedMessage`.
 */
async function withFixtures(mutate, expectedMessage) {
  const dir = await mkdtemp(join(tmpdir(), "verify-authority-plane-"));
  try {
    const scriptsDir = join(dir, "scripts");
    const fixturesOut = join(dir, "schemas", "conformance", "authority-plane-v0");
    await mkdir(scriptsDir, { recursive: true });
    await mkdir(fixturesOut, { recursive: true });
    await cp(scriptPath, join(scriptsDir, "verify-authority-plane-conformance.mjs"));
    await cp(libDir, join(scriptsDir, "lib"), { recursive: true });
    await cp(fixturesDir, fixturesOut, { recursive: true });

    await mutate(fixturesOut);

    await assert.rejects(
      () =>
        execFileAsync(process.execPath, [
          join(scriptsDir, "verify-authority-plane-conformance.mjs"),
        ]),
      expectedMessage,
    );
  } finally {
    await rm(dir, { recursive: true, force: true });
  }
}
