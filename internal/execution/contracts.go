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

type Backend interface {
	Describe() BackendDescription
	Execute(ctx context.Context, job PreparedJob) (Result, error)
}

type RuntimeAdapter interface {
	Name() string
	Prepare(ctx context.Context, job PreparedJob) (PreparedJob, error)
}
