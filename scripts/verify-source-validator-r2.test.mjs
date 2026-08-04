import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { lstat, readFile } from "node:fs/promises";
import test from "node:test";

const artifactRoot = new URL("../artifacts/mjs-source-validator-r2/", import.meta.url);

test("retains two closed unsigned role-specific Source Validator bundles", async () => {
  const construction = JSON.parse(
    await readFile(new URL("evidence/construction.json", artifactRoot), "utf8"),
  );

  assert.equal(construction.schema, "capsule.source-validator.unsigned-construction/v1");
  assert.equal(construction.status, "PASSED");
  assert.equal(construction.enrollment, "not-enrolled");
  assert.equal(construction.signing.appleIdentityUsed, false);
  assert.equal(construction.build.network, "offline");
  assert.equal(construction.roles.length, 2);

  const expected = new Map([
    [
      "daemon",
      {
        service: "com.capsulecorp.capsule.source-validator.daemon.v1",
        bundle: "dist/CapsuleSourceValidatorDaemon.xpc",
        launcher: "Contents/MacOS/CapsuleSourceValidatorDaemonLauncher",
        parser: "Contents/Resources/capsule-mjs-source-validator-daemon",
      },
    ],
    [
      "approval-broker",
      {
        service: "com.capsulecorp.capsule.source-validator.approval-broker.v1",
        bundle: "dist/CapsuleSourceValidatorBroker.xpc",
        launcher: "Contents/MacOS/CapsuleSourceValidatorBrokerLauncher",
        parser: "Contents/Resources/capsule-mjs-source-validator-approval-broker",
      },
    ],
  ]);

  const digests = new Set();
  for (const role of construction.roles) {
    const wanted = expected.get(role.role);
    assert.ok(wanted, `unexpected role ${role.role}`);
    assert.equal(role.serviceIdentifier, wanted.service);
    assert.equal(role.bundlePath, wanted.bundle);
    assert.equal(role.launcher.path, `${wanted.bundle}/${wanted.launcher}`);
    assert.equal(role.parser.path, `${wanted.bundle}/${wanted.parser}`);
    assert.equal(role.resourcePolicy.activation, "inactive");
    assert.equal(role.resourcePolicy.activeMeasurements, false);
    assert.equal(role.launcher.appleIdentity, null);
    assert.equal(role.parser.appleIdentity, null);

    for (const retained of [role.infoPlist, role.launcher, role.parser, role.resourcePolicy]) {
      const path = new URL(retained.path, artifactRoot);
      const metadata = await lstat(path);
      assert.equal(metadata.isFile(), true);
      assert.equal(metadata.isSymbolicLink(), false);
      const bytes = await readFile(path);
      assert.equal(bytes.length, retained.bytes);
      assert.equal(createHash("sha256").update(bytes).digest("hex"), retained.sha256);
    }
    digests.add(role.launcher.sha256);
    digests.add(role.parser.sha256);
  }
  assert.equal(digests.size, 4, "both launchers and parser children must be role-distinct bytes");
});

test("binds exact source, dependency, notice, SBOM, and unsigned provenance evidence", async () => {
  for (const path of [
    "evidence/build-manifest.json",
    "evidence/source-manifest.json",
    "evidence/license-report.json",
    "evidence/sbom.cdx.json",
    "evidence/provenance.intoto.json",
    "evidence/reproduction.json",
  ]) {
    const bytes = await readFile(new URL(path, artifactRoot));
    assert.ok(bytes.length > 0, `${path} is empty`);
  }

  const reproduction = JSON.parse(
    await readFile(new URL("evidence/reproduction.json", artifactRoot), "utf8"),
  );
  assert.equal(reproduction.byteIdentical, true);
  assert.equal(reproduction.sameHost, true);
  assert.equal(reproduction.independentBuilder, false);
  assert.equal(reproduction.roles.length, 2);
});
