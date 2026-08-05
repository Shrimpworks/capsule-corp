#!/usr/bin/env node
import assert from "node:assert/strict";
import { createHash } from "node:crypto";
import { mkdir, readFile, writeFile } from "node:fs/promises";
import { dirname, join, resolve } from "node:path";
import { fileURLToPath } from "node:url";

const repositoryRoot = resolve(dirname(fileURLToPath(import.meta.url)), "..");
const outputRoot = join(repositoryRoot, "schemas/conformance/macos-i2b2-unsigned-installation");
const checkOnly = process.argv.includes("--check");

const fixed = Object.freeze({
  profileId: "capsule.macos-installation.unsigned-installation-only/i2b2",
  coordinatorRoleId: "capsule.role.trust-bootstrap-coordinator",
  coordinatorBundleId: "com.capsulecorp.capsule.trust-bootstrap.v1",
  bootstrapAppGroup: "3DDR84M4JS.com.capsulecorp.capsule.bootstrap.v0",
  bootstrapService: "3DDR84M4JS.com.capsulecorp.capsule.bootstrap.v0.supervisor",
  coordinatorKeychainGroup:
    "3DDR84M4JS.com.capsulecorp.capsule.trust-bootstrap.installation-root.epoch-1",
  supervisorKeychainGroup: "3DDR84M4JS.com.capsulecorp.capsule.supervisor.bootstrap-anchor.epoch-1",
});

function canonicalJSON(value) {
  return `${JSON.stringify(value, null, 2)}\n`;
}

function sha256(bytes) {
  return createHash("sha256").update(bytes).digest("hex");
}

function schemaFor(value) {
  if (Array.isArray(value)) {
    assert.ok(value.length > 0, "profile arrays must be nonempty for recursive schema coverage");
    return {
      type: "array",
      minItems: value.length,
      maxItems: value.length,
      items: schemaFor(value[0]),
    };
  }
  if (value !== null && typeof value === "object") {
    const properties = Object.fromEntries(
      Object.entries(value).map(([key, child]) => [key, schemaFor(child)]),
    );
    return {
      type: "object",
      additionalProperties: false,
      required: Object.keys(properties),
      properties,
    };
  }
  if (typeof value === "string") return { type: "string" };
  if (typeof value === "boolean") return { type: "boolean" };
  if (typeof value === "number") return { type: "integer", minimum: 0 };
  throw new Error(`unsupported profile schema value: ${typeof value}`);
}

async function retained(path, expectedSha256 = "") {
  const bytes = await readFile(join(repositoryRoot, path));
  const digest = sha256(bytes);
  if (expectedSha256) assert.equal(digest, expectedSha256, `${path}: retained digest changed`);
  return { path, bytes: bytes.length, sha256: digest };
}

function entitlement(roleId, key, disposition, valueIdentity) {
  return { roleId, key, disposition, valueIdentity, activationState: "inactive-unsigned-refusing" };
}

const i0Profile = JSON.parse(
  await readFile(join(repositoryRoot, "internal/installation/macosplan/testdata/profile.json")),
);
assert.equal(i0Profile.profileId, "capsule.macos-installation.no-guest/i0");
assert.equal(i0Profile.roles.length, 7);
assert.equal(i0Profile.signing.state, "inactive-refusing");

const i1aManifest = await retained(
  "artifacts/macos-i1a-unsigned-app-shell/dist/Capsule.app/Contents/Resources/CapsuleConstruction/bundle-manifest.json",
  "5bd80097775908031b1a4c90680e8c7656cc5e9f97df2cc187592f75ee67a56f",
);
const i1aEvidence = await retained(
  "artifacts/macos-i1a-unsigned-app-shell/evidence/construction.json",
  "31f79bdbd3dae29f6cfa340683ce59bc445041db0da12a66b1c051abc3db6ae5",
);
const i1bEnrollment = await retained(
  "artifacts/macos-i1b-r3-signed-development-composition/evidence/signed-enrollment.json",
  "afc7002032fc1ff4ead29269e7a370d94524aff42ca9181827a03233a31fbc94",
);
const i2b1Manifest = await retained(
  "schemas/conformance/i2b-bootstrap-v0/manifest.json",
  "9f1b8a86be9ada8e6afa4b913aef027dfe031d9ab69b0d0913c4f63132163203",
);
const coordinatorInfo = await retained(
  "artifacts/macos-i2b2-unsigned-installation-bundle/templates/coordinator-Info.plist",
);
const coordinatorPlaceholder = await retained(
  "artifacts/macos-i2b2-unsigned-installation-bundle/placeholders/CapsuleTrustBootstrap",
);
const supervisorLaunchAgent = await retained(
  "artifacts/macos-i2b2-unsigned-installation-bundle/templates/supervisor-LaunchAgent.plist",
);
const coordinatorEntitlements = await retained(
  "artifacts/macos-i2b2-unsigned-installation-bundle/Entitlements/coordinator.plist",
);
const supervisorEntitlements = await retained(
  "artifacts/macos-i2b2-unsigned-installation-bundle/Entitlements/supervisor.plist",
);
const bootstrapConstraint = await retained(
  "artifacts/macos-i2b2-unsigned-installation-bundle/Constraints/coordinator-supervisor-bootstrap.json",
);

const roles = i0Profile.roles.map((role) => ({
  roleId: role.roleId,
  containingRoleId: role.containingRoleId,
  bundlePath: role.bundlePath,
  executablePath: role.executablePath,
  bundleIdentifier: role.bundleIdentifier,
  signingIdentifier: role.signingIdentifier,
  required: true,
  sourceProfile: "capsule.macos-installation.no-guest/i0",
}));
roles.push({
  roleId: fixed.coordinatorRoleId,
  containingRoleId: "capsule.role.approval-broker/v0",
  bundlePath: "Capsule.app/Contents/XPCServices/CapsuleTrustBootstrap.xpc",
  executablePath:
    "Capsule.app/Contents/XPCServices/CapsuleTrustBootstrap.xpc/Contents/MacOS/CapsuleTrustBootstrap",
  bundleIdentifier: fixed.coordinatorBundleId,
  signingIdentifier: fixed.coordinatorBundleId,
  required: true,
  sourceProfile: fixed.profileId,
});

const profile = {
  schema: "capsule.macos-installation.i2b2-unsigned-profile/v0",
  checkpoint: "I2B2",
  status: "unsigned-declared-inputs-only",
  profileId: fixed.profileId,
  intendedDevelopmentTeamId: "3DDR84M4JS",
  containingRelease: {
    identity: "capsule.macos-installation.i2b2-containing-release/v0",
    baseProfileId: i0Profile.profileId,
    i1aBundleManifest: i1aManifest,
    i1aConstructionEvidence: i1aEvidence,
    i1bSignedDevelopmentEnrollment: i1bEnrollment,
    signedI2B3ReleaseDigestState: "unavailable-no-signing-performed",
    signedI2B3ComponentProfileDigestState: "unavailable-no-signed-readback",
  },
  roles,
  executableInputs: [
    {
      roleId: fixed.coordinatorRoleId,
      identityClass: "inert-non-executable-placeholder",
      ...coordinatorPlaceholder,
    },
    { roleId: fixed.coordinatorRoleId, identityClass: "closed-info-plist", ...coordinatorInfo },
    {
      roleId: "capsule.role.execution-supervisor/v0",
      identityClass: "inactive-launch-agent-descriptor",
      ...supervisorLaunchAgent,
    },
  ],
  services: [
    {
      serviceId: "capsule.service.trust-bootstrap.private-xpc/v0",
      serviceName: fixed.coordinatorBundleId,
      ownerRoleId: fixed.coordinatorRoleId,
      peerRoleId: "capsule.role.approval-broker/v0",
      protocolIdentity: "private-containing-application-invocation-only",
      methodSet: "fixed-setup-invocation-no-generic-signing",
      activationState: "inactive-not-launched",
    },
    {
      serviceId: "capsule.service.supervisor-protected-root-bootstrap/v0",
      serviceName: fixed.bootstrapService,
      ownerRoleId: "capsule.role.execution-supervisor/v0",
      peerRoleId: fixed.coordinatorRoleId,
      protocolIdentity: "capsule.supervisor-protected-root-bootstrap/xpc-v0",
      methodSet: "prepare-observe-and-finalize-only",
      activationState: "inactive-not-registered",
    },
  ],
  entitlements: [
    entitlement(
      fixed.coordinatorRoleId,
      "com.apple.security.app-sandbox",
      "required-true-when-signed",
      "true",
    ),
    entitlement(
      fixed.coordinatorRoleId,
      "com.apple.security.application-groups",
      "required-exact-when-signed",
      fixed.bootstrapAppGroup,
    ),
    entitlement(
      fixed.coordinatorRoleId,
      "keychain-access-groups",
      "required-exact-when-signed",
      fixed.coordinatorKeychainGroup,
    ),
    entitlement(
      "capsule.role.execution-supervisor/v0",
      "com.apple.security.app-sandbox",
      "required-true-when-signed",
      "true",
    ),
    entitlement(
      "capsule.role.execution-supervisor/v0",
      "com.apple.security.application-groups",
      "required-exact-when-signed",
      fixed.bootstrapAppGroup,
    ),
    entitlement(
      "capsule.role.execution-supervisor/v0",
      "keychain-access-groups",
      "required-exact-when-signed",
      fixed.supervisorKeychainGroup,
    ),
    entitlement(
      "capsule.role.approval-broker/v0",
      "com.apple.security.application-groups",
      "required-absent",
      fixed.bootstrapAppGroup,
    ),
    entitlement(
      "capsule.role.agent-daemon/v0",
      "com.apple.security.application-groups",
      "required-absent",
      fixed.bootstrapAppGroup,
    ),
  ],
  entitlementInputs: [coordinatorEntitlements, supervisorEntitlements],
  constraints: {
    projection: bootstrapConstraint,
    teamIdentifier: "3DDR84M4JS",
    coordinatorSigningIdentifier: fixed.coordinatorBundleId,
    supervisorSigningIdentifier: "com.capsulecorp.capsule.supervisor",
    activeCodeDirectoryHashSetState: "unavailable-no-signing-performed",
    activeEffectiveEntitlementDigestSetState: "unavailable-no-signing-performed",
    activationDecision: "refuse",
  },
  bootstrapObjects: {
    fixtureManifest: i2b1Manifest,
    request: {
      objectType: "capsule.supervisor-bootstrap-request",
      version: 0,
      purpose: "capsule.installation.bootstrap.request",
      audience: "capsule.execution-supervisor.bootstrap",
      mediaType: "application/capsule.supervisor-bootstrap-request+cbor;v=0",
      containingReleaseBinding: "signed-i2b3-containing-release-digest-required",
      i1ProfileBinding: "exact-i1b-signed-development-enrollment-source-required",
      componentProfileBinding: "signed-i2b3-component-profile-digest-required",
      signingState: "unavailable-no-key-or-signer",
    },
    record: {
      objectType: "capsule.supervisor-bootstrap-record",
      version: 0,
      purpose: "capsule.installation.bootstrap.record",
      audience: "capsule.execution-supervisor",
      mediaType: "application/capsule.supervisor-bootstrap-record+cbor;v=0",
      containingReleaseBinding: "exact-request-binding-required",
      i1ProfileBinding: "exact-request-binding-required",
      componentProfileBinding: "exact-request-binding-required",
      signingState: "unavailable-no-key-or-signer",
    },
  },
  protectedState: {
    stateRootClass: "supervisor-private-app-sandbox",
    privateParentName: "CapsuleSupervisor",
    stateRootName: "supervisor.state",
    pendingJournalName: "CapsuleSupervisor.bootstrap-pending",
    stagingRootPrefix: "supervisor.state.staging.",
    publishIntentName: "supervisor.state.publish-intent",
    ownerEntryName: "supervisor.owner",
    ownerMode: 384,
    ownerLinkCount: 1,
    ownerProfile: "darwin-openat-flock-v0",
    storeEntryName: "supervisor.store",
    storeFormat: "capsule.supervisor-store/fixed-v1",
    storeProfile: "conformance-non-product-no-guest",
    requestEntryName: "supervisor.bootstrap-request",
    recordEntryName: "supervisor.bootstrap-record",
    ordinaryStartupOperation: "open-without-create",
  },
  stateProjection: {
    signingState: "unsigned-no-apple-identity",
    provisioningProfileState: "absent",
    coordinatorExecutableState: "inert-non-executable-placeholder",
    supervisorExecutableState: "preserved-i1a-inert-non-executable-placeholder",
    serviceRegistrationState: "inactive-not-registered",
    serviceLaunchState: "inactive-not-launched",
    coordinatorSigningState: "unavailable",
    installationRootKeyState: "absent",
    bootstrapRequestState: "absent",
    bootstrapRecordState: "absent",
    protectedRootState: "absent-no-create",
    ownerState: "absent-no-create",
    storeState: "absent-no-create",
    attemptsEnabled: false,
    runtimePresent: false,
    backendPresent: false,
    guestPresent: false,
    activationDecision: "refuse",
    activationReason: "unsigned-profile-inactive",
  },
  cleanupRepair: {
    cleanupProjection: "remove-unsigned-bundle-copy-only",
    protectedStateMutation: "none",
    appGroupMutation: "none",
    keychainMutation: "none",
    serviceMutation: "none",
    repairProjection: "missing-or-mixed-input-refuse-before-side-effect",
    nextGate: "I2B3-separate-authorization-signing-key-container-service-handoff",
  },
  readbackCaps: {
    profileRawBytes: 65536,
    bundleManifestRawBytes: 262144,
    bundleFileCount: 31,
    bundlePathUtf8Bytes: 1024,
  },
};

const profileBytes = canonicalJSON(profile);
assert.ok(Buffer.byteLength(profileBytes) <= profile.readbackCaps.profileRawBytes);

const mutations = [
  { id: "profile-raw-cap-plus-one", expected: "refuse-profile-raw-cap" },
  { id: "bundle-path-cap-plus-one", expected: "refuse-bundle-path-cap" },
  { id: "role-missing", expected: "refuse-role-inventory" },
  { id: "role-extra", expected: "refuse-role-inventory" },
  { id: "role-duplicate", expected: "refuse-role-inventory" },
  { id: "role-mixed", expected: "refuse-role-inventory" },
  { id: "wrong-containing-release", expected: "refuse-release-binding" },
  { id: "wrong-i1-profile", expected: "refuse-profile-binding" },
  { id: "wrong-bootstrap-service", expected: "refuse-service-binding" },
  { id: "unsafe-entitlement", expected: "refuse-entitlement" },
  { id: "active-signing", expected: "refuse-inactive-boundary" },
  { id: "bootstrap-created", expected: "refuse-no-create-boundary" },
  { id: "store-created", expected: "refuse-no-create-boundary" },
  { id: "unknown-field", expected: "refuse-schema" },
];

const schema = {
  $schema: "https://json-schema.org/draft/2020-12/schema",
  $id: "https://capsule.local/schemas/macos-i2b2-unsigned-installation-profile:v0",
  title: "MacosI2b2UnsignedInstallationProfile v0",
  ...schemaFor(profile),
};
const outputs = new Map([
  ["profile.json", profileBytes],
  ["profile.schema.json", canonicalJSON(schema)],
  ["mutation-corpus.json", canonicalJSON(mutations)],
]);
const generated = [];
for (const name of ["profile.json", "profile.schema.json", "mutation-corpus.json"]) {
  const bytes = Buffer.from(outputs.get(name));
  generated.push({
    path: `schemas/conformance/macos-i2b2-unsigned-installation/${name}`,
    bytes: bytes.length,
    sha256: sha256(bytes),
  });
}
const manifest = {
  schema: "capsule.macos-installation.i2b2-fixture-manifest/v0",
  status: "unsigned-inactive-no-side-effects",
  profileId: fixed.profileId,
  cases: mutations.length,
  files: generated,
};
outputs.set("manifest.json", canonicalJSON(manifest));
await mkdir(outputRoot, { recursive: true });
for (const [name, bytes] of outputs) {
  const path = join(outputRoot, name);
  if (checkOnly) {
    assert.deepEqual(await readFile(path), Buffer.from(bytes), `${name}: generated fixture drift`);
  } else {
    await writeFile(path, bytes, { mode: 0o644 });
  }
}
process.stdout.write(
  canonicalJSON({
    profileSha256: sha256(profileBytes),
    profileBytes: Buffer.byteLength(profileBytes),
    cases: mutations.length,
  }),
);
