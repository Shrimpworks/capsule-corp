// Package macosplan defines the passive macOS installation Slice I0 contract.
//
// It performs no filesystem, signing, Keychain, Service Management, XPC,
// process, update, repair, uninstall, runtime, backend, or guest operation.
// The canonical profile is deliberately inactive until later signed installed
// slices replace every refusing dependency with retained evidence.
//
//go:generate go run ./cmd/generatefixtures
package macosplan
