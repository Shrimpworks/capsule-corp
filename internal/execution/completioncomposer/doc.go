// Package completioncomposer implements an unwired terminal projection and a
// bounded fixed-file completion-last transaction for the no-guest FakeBackend.
// It composes only Supervisor-retained attempt, completion, lifecycle,
// teardown, and fake authoritative-absence facts. It does not expose an
// endpoint, sign evidence, create a receipt, call a backend, or create a guest.
package completioncomposer
