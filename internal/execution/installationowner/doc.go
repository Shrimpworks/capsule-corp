// Package installationowner implements the passive installation-global
// Supervisor owner capability selected by Proposed ADR-0033.
//
// The package is internal and unwired. It serializes cooperating Darwin
// processes around one installer-enrolled object; it does not authorize store,
// archive, IPC, adapter, backend, runtime, or guest work and does not establish
// same-UID storage containment.
package installationowner
