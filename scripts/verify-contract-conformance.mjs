import { createHash } from "node:crypto";
import { lstat, readdir, readFile } from "node:fs/promises";
import { isAbsolute, relative, resolve, sep } from "node:path";
import { pathToFileURL } from "node:url";
import Ajv2020 from "ajv/dist/2020.js";

const defaultRootDirectory = new URL("../schemas/conformance/v0/", import.meta.url);
const metadataFiles = new Set(["manifest.json", "manifest.schema.json"]);

export async function verifyConformanceCorpus({ rootDirectory = defaultRootDirectory } = {}) {
  const root = resolveDirectory(rootDirectory);
  const schema = await readJson(resolve(root, "manifest.schema.json"), "manifest schema");
  const manifest = await readJson(resolve(root, "manifest.json"), "manifest");
  const ajv = new Ajv2020({ allErrors: true, strict: true });
  const validate = ajv.compile(schema);

  if (!validate(manifest)) {
    throw new Error(`manifest schema validation failed: ${ajv.errorsText(validate.errors)}`);
  }

  assertUnique(manifest.rules, "rule");
  assertUnique(manifest.cases, "case");
  const rules = new Map(manifest.rules.map((rule) => [rule.id, rule]));
  const listedFixtures = new Set();

  for (const entry of manifest.cases) {
    for (const ruleId of entry.ruleIds) {
      if (!rules.has(ruleId)) {
        throw new Error(`case ${entry.id} references unknown rule ${ruleId}`);
      }
    }

    if (entry.expected.decision === "reject" && entry.expected.authorityStateChanged) {
      throw new Error(`rejected case ${entry.id} cannot change authority-bearing state`);
    }
    if (
      entry.expected.decision === "reject" &&
      entry.expected.owner === "source-validator-passive-contract" &&
      (!entry.expected.effects || Object.values(entry.expected.effects).some(Boolean))
    ) {
      throw new Error(
        `rejected passive Source Validator case ${entry.id} must assert zero state/approval/key/endpoint/process/runtime/backend/guest effects`,
      );
    }

    assertScalarRoleContext(entry);
    await verifyFixture(root, entry.fixture, entry.id);
    listedFixtures.add(entry.fixture.path);

    for (const [label, fixture] of contextFixtures(entry.context)) {
      await verifyFixture(root, fixture, `${entry.id} ${label}`);
      listedFixtures.add(fixture.path);
    }
    if (entry.expected.stateDelta) {
      await verifyFixture(root, entry.expected.stateDelta.after, `${entry.id} expected post-state`);
      listedFixtures.add(entry.expected.stateDelta.after.path);
    }
    assertProposalResolutionContext(entry);
    await assertApprovalAttemptStateContext(root, entry);
    await assertRegistrationStateContext(root, entry);
  }

  for (const rule of manifest.rules) {
    for (const requirement of rule.requiredCases) {
      const covered = manifest.cases.some(
        (entry) =>
          entry.ruleIds.includes(rule.id) &&
          entry.expected.decision === requirement.decision &&
          entry.variant === requirement.variant,
      );
      if (!covered) {
        throw new Error(
          `required case is not covered: ${rule.id} ${requirement.decision}/${requirement.variant}`,
        );
      }
    }
  }

  const corpusFiles = await listCorpusFiles(root);
  for (const path of corpusFiles) {
    if (!listedFixtures.has(path)) {
      throw new Error(`fixture is not listed in the manifest: ${path}`);
    }
  }

  return {
    caseCount: manifest.cases.length,
    fixtureCount: listedFixtures.size,
    ruleCount: manifest.rules.length,
  };
}

async function verifyFixture(root, fixture, owner) {
  const fixturePath = resolveFixturePath(root, fixture.path);
  let metadata;
  try {
    metadata = await lstat(fixturePath);
  } catch (error) {
    if (error?.code === "ENOENT") {
      throw new Error(`listed fixture does not exist for ${owner}: ${fixture.path}`);
    }
    throw error;
  }
  if (!metadata.isFile() || metadata.isSymbolicLink()) {
    throw new Error(`listed fixture is not a regular file for ${owner}: ${fixture.path}`);
  }

  const bytes = await readFile(fixturePath);
  if (bytes.length !== fixture.byteLength) {
    throw new Error(
      `fixture byte length mismatch for ${owner}: expected ${fixture.byteLength}, got ${bytes.length}`,
    );
  }
  const digest = createHash("sha256").update(bytes).digest("hex");
  if (digest !== fixture.sha256) {
    throw new Error(
      `fixture SHA-256 mismatch for ${owner}: expected ${fixture.sha256}, got ${digest}`,
    );
  }
}

function resolveFixturePath(root, path) {
  if (isAbsolute(path) || path.includes("\\")) {
    throw new Error(`fixture path is not a safe relative path: ${path}`);
  }
  const segments = path.split("/");
  if (segments.some((segment) => segment === "" || segment === "." || segment === "..")) {
    throw new Error(`fixture path is not a safe relative path: ${path}`);
  }
  const resolved = resolve(root, path);
  if (!resolved.startsWith(`${root}${sep}`)) {
    throw new Error(`fixture path escapes the corpus root: ${path}`);
  }
  return resolved;
}

async function listCorpusFiles(root) {
  const files = [];
  await visit(root);
  return files.sort();

  async function visit(directory) {
    const entries = await readdir(directory, { withFileTypes: true });
    for (const entry of entries) {
      const absolutePath = resolve(directory, entry.name);
      const corpusPath = relative(root, absolutePath).split(sep).join("/");
      if (entry.isSymbolicLink()) {
        throw new Error(`corpus contains a symbolic link: ${corpusPath}`);
      }
      if (entry.isDirectory()) {
        await visit(absolutePath);
      } else if (entry.isFile() && !metadataFiles.has(corpusPath)) {
        files.push(corpusPath);
      } else if (!entry.isFile()) {
        throw new Error(`corpus contains a non-regular entry: ${corpusPath}`);
      }
    }
  }
}

function assertUnique(entries, kind) {
  const ids = new Set();
  for (const entry of entries) {
    if (ids.has(entry.id)) {
      throw new Error(`duplicate ${kind} ID: ${entry.id}`);
    }
    ids.add(entry.id);
  }
}

function assertScalarRoleContext(entry) {
  if (entry.context.kind !== "scalar-role") {
    return;
  }
  const sameRole = entry.context.expectedRole === entry.context.providedRole;
  if (entry.expected.decision === "accept" && !sameRole) {
    throw new Error(`accepted scalar-role case ${entry.id} must use the expected role`);
  }
  if (
    entry.expected.decision === "reject" &&
    entry.expected.classification === "DOMAIN" &&
    sameRole
  ) {
    throw new Error(`DOMAIN case ${entry.id} must substitute a different scalar role`);
  }
}

function contextFixtures(context) {
  if (context.kind === "fixture") {
    return [["context", context.fixture]];
  }
  if (context.kind === "source-manifest") {
    return [["source bytes", context.source]];
  }
  if (context.kind === "source-validator") {
    return [
      ...(context.source ? [["source bytes", context.source]] : []),
      ...(context.request ? [["validator request", context.request]] : []),
      ...(context.artifactProfile ? [["validator artifact profile", context.artifactProfile]] : []),
      ...(context.engineeringCandidate
        ? [["validator engineering candidate", context.engineeringCandidate]]
        : []),
      ...(context.sourceManifest ? [["source manifest", context.sourceManifest]] : []),
    ];
  }
  if (context.kind !== "proposal-resolution") {
    if (context.kind === "approval-attempt-state") {
      return [
        ["approval/attempt operation", context.operation],
        ["approval/attempt pre-state", context.before],
        ...(context.payload ? [["approval payload", context.payload]] : []),
        ...(context.protectedHeader
          ? [["approval protected header", context.protectedHeader]]
          : []),
      ];
    }
    return context.kind === "registration-state"
      ? [
          ["registration operation", context.operation],
          ["registration pre-state", context.before],
        ]
      : [];
  }
  return [
    ["profile registry", context.profileRegistry],
    ["user policy", context.userPolicy],
    ...(context.oracle.sourceManifest
      ? [["source manifest oracle", context.oracle.sourceManifest]]
      : []),
    ...(context.oracle.canonicalInlineInput
      ? [["canonical inline-input oracle", context.oracle.canonicalInlineInput]]
      : []),
  ];
}

async function assertRegistrationStateContext(root, entry) {
  if (entry.context.kind !== "registration-state") {
    if (entry.context.kind !== "approval-attempt-state" && entry.expected.stateDelta) {
      throw new Error(`non-state case ${entry.id} cannot retain a registration state delta`);
    }
    return;
  }
  if (entry.expected.stateDelta?.kind !== "exact") {
    throw new Error(`registration-state case ${entry.id} requires an exact post-state`);
  }

  const operation = await readFixtureJson(root, entry.context.operation, `${entry.id} operation`);
  const before = await readFixtureJson(root, entry.context.before, `${entry.id} pre-state`);
  const after = await readFixtureJson(
    root,
    entry.expected.stateDelta.after,
    `${entry.id} post-state`,
  );
  assertExactKeys(
    operation,
    [
      "contextType",
      "contextVersion",
      "method",
      "caller",
      "trustedClockObservationUnixSeconds",
      "concurrentHighWaterUnixSeconds",
      "identifier",
      "mutation",
      "fault",
    ],
    `${entry.id} operation`,
  );
  if (
    operation.contextType !== "capsule.conformance.registration-operation" ||
    operation.contextVersion !== 0 ||
    !["register-plan", "use-registration"].includes(operation.method) ||
    !Number.isSafeInteger(operation.trustedClockObservationUnixSeconds) ||
    operation.trustedClockObservationUnixSeconds < 0 ||
    !(
      operation.concurrentHighWaterUnixSeconds === null ||
      (Number.isSafeInteger(operation.concurrentHighWaterUnixSeconds) &&
        operation.concurrentHighWaterUnixSeconds >= 0)
    )
  ) {
    throw new Error(`registration-state case ${entry.id} has an unknown operation context`);
  }
  assertExactKeys(
    operation.caller,
    ["authenticated", "role", "purpose"],
    `${entry.id} operation caller`,
  );
  if (
    typeof operation.caller.authenticated !== "boolean" ||
    !["daemon", "broker", "updater"].includes(operation.caller.role) ||
    operation.caller.purpose !== "register-plan"
  ) {
    throw new Error(`registration-state case ${entry.id} has an invalid caller context`);
  }
  for (const [label, state] of [
    ["pre-state", before],
    ["post-state", after],
  ]) {
    assertRegistrationStateShape(entry.id, label, state);
  }
  if (after.timeHighWaterUnixSeconds < before.timeHighWaterUnixSeconds) {
    throw new Error(`registration-state case ${entry.id} moved effective time backward`);
  }
  if (
    after.installationIdHex !== before.installationIdHex ||
    after.supervisorIdHex !== before.supervisorIdHex ||
    after.epochSequence !== before.epochSequence ||
    after.epochDigestHex !== before.epochDigestHex
  ) {
    throw new Error(`registration-state case ${entry.id} changed bound installation identity`);
  }

  const beforePopulation = before.registrationPopulation;
  const afterPopulation = after.registrationPopulation;
  if (entry.expected.decision === "reject") {
    if (entry.expected.authorityStateChanged) {
      throw new Error(`rejected registration-state case ${entry.id} created authority`);
    }
    if (
      after.lastRegistrationSequence !== before.lastRegistrationSequence ||
      afterPopulation.storedCount !== beforePopulation.storedCount ||
      afterPopulation.unexpiredCount !== beforePopulation.unexpiredCount ||
      afterPopulation.setDigest !== beforePopulation.setDigest ||
      JSON.stringify(after.materializedRecords) !== JSON.stringify(before.materializedRecords)
    ) {
      throw new Error(
        `rejected registration-state case ${entry.id} changed registration authority`,
      );
    }
    if (trustStateWidened(before, after)) {
      throw new Error(`rejected registration-state case ${entry.id} widened trust state`);
    }
    return;
  }

  if (operation.method !== "register-plan" || !entry.expected.authorityStateChanged) {
    throw new Error(`accepted registration-state case ${entry.id} must create one registration`);
  }
  if (
    after.trustPhase !== before.trustPhase ||
    after.trustReason !== before.trustReason ||
    after.quarantined !== before.quarantined ||
    after.lastRegistrationSequence !== before.lastRegistrationSequence + 1 ||
    afterPopulation.storedCount !== beforePopulation.storedCount + 1 ||
    afterPopulation.unexpiredCount !== beforePopulation.unexpiredCount + 1 ||
    afterPopulation.setDigest === beforePopulation.setDigest ||
    after.materializedRecords.length !== before.materializedRecords.length + 1
  ) {
    throw new Error(`accepted registration-state case ${entry.id} lacks one atomic append`);
  }
  const record = after.materializedRecords.at(-1);
  assertExactKeys(
    record,
    [
      "wireRegistrationHex",
      "exactPlanHex",
      "recomputedPlanDigestHex",
      "registeredAtUnixSeconds",
      "storageFormatVersion",
      "retentionState",
    ],
    `${entry.id} stored record`,
  );
  const planBytes = await readFile(resolveFixturePath(root, entry.fixture.path));
  if (
    !/^[0-9a-f]+$/u.test(record.wireRegistrationHex) ||
    record.wireRegistrationHex.length % 2 !== 0 ||
    record.exactPlanHex !== planBytes.toString("hex") ||
    record.recomputedPlanDigestHex !== createHash("sha256").update(planBytes).digest("hex") ||
    !record.wireRegistrationHex.includes(`055820${record.recomputedPlanDigestHex}`) ||
    record.registeredAtUnixSeconds !== after.timeHighWaterUnixSeconds ||
    record.storageFormatVersion !== 0 ||
    record.retentionState !== "retained"
  ) {
    throw new Error(`accepted registration-state case ${entry.id} has an invalid stored record`);
  }
}

async function assertApprovalAttemptStateContext(root, entry) {
  if (entry.context.kind !== "approval-attempt-state") {
    return;
  }
  if (entry.expected.stateDelta?.kind !== "exact") {
    throw new Error(`approval/attempt case ${entry.id} requires an exact post-state`);
  }
  const operation = await readFixtureJson(root, entry.context.operation, `${entry.id} operation`);
  const before = await readFixtureJson(root, entry.context.before, `${entry.id} pre-state`);
  const after = await readFixtureJson(
    root,
    entry.expected.stateDelta.after,
    `${entry.id} post-state`,
  );
  assertApprovalAttemptOperation(entry.id, operation);
  assertApprovalAttemptState(entry.id, "pre-state", before);
  assertApprovalAttemptState(entry.id, "post-state", after);
  if (operation.mode !== "store-transition") {
    for (const field of [
      "timeHighWaterChanged",
      "trustStateTightened",
      "fakeBackendEffectPermitted",
    ]) {
      if (entry.expected[field] !== false) {
        throw new Error(`passive Slice A case ${entry.id} must set ${field} to false`);
      }
    }
    if (entry.expected.authorityStateChanged) {
      throw new Error(`passive Slice A case ${entry.id} cannot change authority state`);
    }
    if (JSON.stringify(before) !== JSON.stringify(after)) {
      throw new Error(`passive Slice A case ${entry.id} changed approval or attempt state`);
    }
    return;
  }
  assertDurableApprovalAttemptTransition(entry, operation, before, after);
}

function assertDurableApprovalAttemptTransition(entry, operation, before, after) {
  if (entry.expected.fakeBackendEffectPermitted) {
    throw new Error(`Slice B case ${entry.id} cannot permit a fake-backend effect`);
  }
  if (
    after.installationIdHex !== before.installationIdHex ||
    after.supervisorIdHex !== before.supervisorIdHex ||
    after.epochSequence !== before.epochSequence ||
    after.epochDigestHex !== before.epochDigestHex ||
    after.timeHighWaterUnixSeconds < before.timeHighWaterUnixSeconds
  ) {
    throw new Error(`Slice B case ${entry.id} changed identity or moved time backward`);
  }
  const approvalChanged =
    JSON.stringify(after.approvalPopulation) !== JSON.stringify(before.approvalPopulation) ||
    JSON.stringify(after.materializedApprovals) !== JSON.stringify(before.materializedApprovals);
  const attemptChanged =
    JSON.stringify(after.attemptPopulation) !== JSON.stringify(before.attemptPopulation) ||
    JSON.stringify(after.materializedAttempts) !== JSON.stringify(before.materializedAttempts);
  if (entry.expected.decision === "reject" && entry.expected.authorityStateChanged) {
    throw new Error(`rejected Slice B case ${entry.id} claims an authority change`);
  }
  if (
    entry.expected.classification !== "RECOVERY_REQUIRED" &&
    entry.expected.authorityStateChanged !== (approvalChanged || attemptChanged)
  ) {
    throw new Error(`Slice B case ${entry.id} has a mismatched authority-state oracle`);
  }
  if (
    entry.expected.decision === "reject" &&
    entry.expected.classification !== "RECOVERY_REQUIRED" &&
    (approvalChanged || attemptChanged)
  ) {
    throw new Error(`rejected Slice B case ${entry.id} changed approval/attempt authority`);
  }
  if (operation.recovery !== after.recoveryFence && operation.scenario.includes("indeterminate")) {
    throw new Error(`indeterminate Slice B case ${entry.id} lacks its recovery fence`);
  }
}

function assertApprovalAttemptOperation(caseId, operation) {
  if (
    operation?.contextType !== "capsule.conformance.approval-attempt-operation" ||
    operation.contextVersion !== 0 ||
    ![
      "identifier",
      "reference",
      "classification-vocabulary",
      "fixture-verifier",
      "store-transition",
    ].includes(operation.mode)
  ) {
    throw new Error(`approval/attempt case ${caseId} has an invalid operation context`);
  }
  const keysByMode = {
    identifier: ["contextType", "contextVersion", "mode", "expectedDomain", "providedDomain"],
    reference: ["contextType", "contextVersion", "mode", "referenceKind", "providedDomain"],
    "classification-vocabulary": ["contextType", "contextVersion", "mode"],
    "fixture-verifier": [
      "contextType",
      "contextVersion",
      "mode",
      "vector",
      "bindingMutation",
      "callerMutation",
    ],
    "store-transition": ["contextType", "contextVersion", "mode", "method", "scenario", "recovery"],
  };
  assertExactKeys(operation, keysByMode[operation.mode], `${caseId} approval/attempt operation`);
  const domains = ["approval", "attempt", "attempt-nonce"];
  if (
    operation.mode === "identifier" &&
    (!domains.includes(operation.expectedDomain) || !domains.includes(operation.providedDomain))
  ) {
    throw new Error(`approval/attempt case ${caseId} has invalid identifier roles`);
  }
  if (
    operation.mode === "reference" &&
    (!domains.includes(operation.providedDomain) ||
      !["approval-reference", "attempt-reference"].includes(operation.referenceKind))
  ) {
    throw new Error(`approval/attempt case ${caseId} has invalid reference roles`);
  }
  if (
    operation.mode === "fixture-verifier" &&
    (typeof operation.vector !== "string" ||
      typeof operation.bindingMutation !== "string" ||
      typeof operation.callerMutation !== "boolean")
  ) {
    throw new Error(`approval/attempt case ${caseId} has invalid verifier context`);
  }
  if (
    operation.mode === "store-transition" &&
    (!["submit-approval", "request-attempt"].includes(operation.method) ||
      typeof operation.scenario !== "string" ||
      typeof operation.recovery !== "boolean")
  ) {
    throw new Error(`approval/attempt case ${caseId} has invalid store transition context`);
  }
}

function assertApprovalAttemptState(caseId, label, state) {
  assertExactKeys(
    state,
    [
      "contextType",
      "contextVersion",
      "installationIdHex",
      "supervisorIdHex",
      "epochSequence",
      "epochDigestHex",
      "trustPhase",
      "trustReason",
      "attemptsEnabled",
      "recoveryFence",
      "timeHighWaterUnixSeconds",
      "approvalPopulation",
      "attemptPopulation",
      "materializedApprovals",
      "materializedAttempts",
    ],
    `${caseId} ${label}`,
  );
  if (
    state.contextType !== "capsule.conformance.approval-attempt-state" ||
    state.contextVersion !== 0 ||
    !/^[0-9a-f]{32}$/u.test(state.installationIdHex) ||
    !/^[0-9a-f]{32}$/u.test(state.supervisorIdHex) ||
    !/^[0-9a-f]{64}$/u.test(state.epochDigestHex) ||
    !Number.isSafeInteger(state.epochSequence) ||
    state.epochSequence < 0 ||
    !["stable", "transition-fenced", "repair-required"].includes(state.trustPhase) ||
    !(state.trustReason === null || typeof state.trustReason === "string") ||
    typeof state.attemptsEnabled !== "boolean" ||
    typeof state.recoveryFence !== "boolean" ||
    !Number.isSafeInteger(state.timeHighWaterUnixSeconds) ||
    state.timeHighWaterUnixSeconds < 0 ||
    !Array.isArray(state.materializedApprovals) ||
    !Array.isArray(state.materializedAttempts)
  ) {
    throw new Error(`approval/attempt case ${caseId} has invalid ${label} scalars`);
  }
  assertPopulation(state.approvalPopulation, "usableCount", `${caseId} ${label} approvals`);
  assertPopulation(state.attemptPopulation, "nonterminalCount", `${caseId} ${label} attempts`);
}

function assertPopulation(population, liveField, label) {
  assertExactKeys(population, [liveField, "retainedCount", "setDigest"], label);
  if (
    !Number.isSafeInteger(population[liveField]) ||
    population[liveField] < 0 ||
    !Number.isSafeInteger(population.retainedCount) ||
    population.retainedCount < population[liveField] ||
    !/^[0-9a-f]{64}$/u.test(population.setDigest)
  ) {
    throw new Error(`${label} is invalid`);
  }
}

function trustStateWidened(before, after) {
  if (before.quarantined && !after.quarantined) {
    return true;
  }
  if (before.trustPhase === "stable") {
    return false;
  }
  return after.trustPhase !== before.trustPhase || after.trustReason !== before.trustReason;
}

function assertRegistrationStateShape(caseId, label, state) {
  assertExactKeys(
    state,
    [
      "contextType",
      "contextVersion",
      "installationIdHex",
      "supervisorIdHex",
      "epochSequence",
      "epochDigestHex",
      "trustPhase",
      "trustReason",
      "quarantined",
      "timeHighWaterUnixSeconds",
      "lastRegistrationSequence",
      "registrationPopulation",
      "materializedRecords",
    ],
    `${caseId} ${label}`,
  );
  if (
    state.contextType !== "capsule.conformance.registration-state" ||
    state.contextVersion !== 0 ||
    !/^[0-9a-f]{32}$/u.test(state.installationIdHex) ||
    !/^[0-9a-f]{32}$/u.test(state.supervisorIdHex) ||
    !/^[0-9a-f]{64}$/u.test(state.epochDigestHex) ||
    !Number.isSafeInteger(state.epochSequence) ||
    state.epochSequence < 0 ||
    !["stable", "transition-fenced", "repair-required"].includes(state.trustPhase) ||
    !(state.trustReason === null || typeof state.trustReason === "string") ||
    typeof state.quarantined !== "boolean" ||
    !Number.isSafeInteger(state.timeHighWaterUnixSeconds) ||
    state.timeHighWaterUnixSeconds < 0 ||
    !Number.isSafeInteger(state.lastRegistrationSequence) ||
    state.lastRegistrationSequence < 0 ||
    !Array.isArray(state.materializedRecords)
  ) {
    throw new Error(`registration-state case ${caseId} has invalid ${label} scalars`);
  }
  assertExactKeys(
    state.registrationPopulation,
    ["storedCount", "unexpiredCount", "setDigest"],
    `${caseId} ${label} population`,
  );
  if (
    !Number.isSafeInteger(state.registrationPopulation.storedCount) ||
    state.registrationPopulation.storedCount < 0 ||
    !Number.isSafeInteger(state.registrationPopulation.unexpiredCount) ||
    state.registrationPopulation.unexpiredCount < 0 ||
    state.registrationPopulation.unexpiredCount > state.registrationPopulation.storedCount ||
    !/^[0-9a-f]{64}$/u.test(state.registrationPopulation.setDigest)
  ) {
    throw new Error(`registration-state case ${caseId} has an invalid ${label} population`);
  }
}

function assertExactKeys(value, expected, label) {
  if (
    value === null ||
    typeof value !== "object" ||
    JSON.stringify(Object.keys(value)) !== JSON.stringify(expected)
  ) {
    throw new Error(`${label} does not use the closed fixture shape`);
  }
}

async function readFixtureJson(root, fixture, label) {
  try {
    return JSON.parse(await readFile(resolveFixturePath(root, fixture.path), "utf8"));
  } catch (error) {
    throw new Error(`failed to read ${label}: ${error.message}`, { cause: error });
  }
}

function assertProposalResolutionContext(entry) {
  if (entry.context.kind !== "proposal-resolution") {
    return;
  }
  const { oracle } = entry.context;
  for (const [label, fixture, digest] of [
    ["source manifest", oracle.sourceManifest, oracle.sourceManifestDigest],
    ["canonical inline input", oracle.canonicalInlineInput, oracle.inlineInputDigest],
  ]) {
    if ((fixture === null) !== (digest === null)) {
      throw new Error(`${entry.id} ${label} fixture and digest must be present together`);
    }
    if (fixture && fixture.sha256 !== digest) {
      throw new Error(`${entry.id} ${label} digest must match the retained fixture digest`);
    }
  }
  if (
    entry.expected.decision === "reject" &&
    (oracle.sourceManifest || oracle.canonicalInlineInput || oracle.wallTime)
  ) {
    throw new Error(`rejected proposal case ${entry.id} cannot retain a resolution result`);
  }
}

async function readJson(path, label) {
  try {
    return JSON.parse(await readFile(path, "utf8"));
  } catch (error) {
    throw new Error(`failed to read ${label} at ${path}: ${error.message}`, { cause: error });
  }
}

function resolveDirectory(value) {
  return value instanceof URL ? resolve(value.pathname) : resolve(value);
}

if (process.argv[1] && pathToFileURL(resolve(process.argv[1])).href === import.meta.url) {
  const result = await verifyConformanceCorpus();
  process.stdout.write(
    `validated conformance corpus: ${result.ruleCount} rules, ${result.caseCount} cases, ${result.fixtureCount} fixtures\n`,
  );
}
