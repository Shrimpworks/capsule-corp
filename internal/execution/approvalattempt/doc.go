// Package approvalattempt retains the ADR-0024 candidate contracts, the
// fixture verifier used by the unwired colocated registration store, and the
// ADR-0043 public-key-only production-shaped verifier. The latter is not wired
// to the store. The package remains unwired from any deployed Supervisor,
// consumer, key operation, backend, or guest.
package approvalattempt
