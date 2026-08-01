package execution

import "context"

type BackendTier string

const (
	BackendTierDevelopment        BackendTier = "development"
	BackendTierAuthoritativeLocal BackendTier = "authoritative-local"
	BackendTierHostedHardened     BackendTier = "hosted-hardened"
	BackendTierHighAssurance      BackendTier = "high-assurance"
)

type BackendDescription struct {
	Name string
	Tier BackendTier
}

type PreparedJob struct {
	ID             string
	RuntimeProfile string
	SourceDigest   string
}

type Result struct {
	Status   string
	ExitCode *int
}

// Backend is the original build scaffold only.
//
// Deprecated: a single Execute call cannot satisfy Capsule's durable
// prepare/create/stage/start/collect/destroy/reconcile contract. Do not wire
// this interface to SupervisorCore or use it to launch hostile code.
type Backend interface {
	Describe() BackendDescription
	Execute(ctx context.Context, job PreparedJob) (Result, error)
}

// RuntimeAdapter is the original build scaffold only. Runtime preparation is
// not a security boundary and cannot grant backend or host authority.
type RuntimeAdapter interface {
	Name() string
	Prepare(ctx context.Context, job PreparedJob) (PreparedJob, error)
}
