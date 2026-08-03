// Package lifecyclestate defines the passive, runtime-neutral durable
// lifecycle contract proposed by ADR-0025. It performs no I/O, persists no
// state, calls no adapter, and authorizes no runtime, backend, or guest.
package lifecyclestate
