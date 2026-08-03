// Package api is the capsule daemon's local HTTP surface. It currently
// exposes only read-only, unauthenticated diagnostic endpoints (health,
// version, runtime listing) — no agent-facing job proposal, plan,
// approval, or execution endpoint is wired here yet. Per AGENTS.md, the
// daemon must never itself hold Approval/Supervisor authority; when
// job-facing endpoints are added, they proxy to the Execution Supervisor
// rather than implementing that authority in this package.
package api

import (
	"encoding/json"
	"net/http"
)

// Options carries build-time identity reported by the version endpoint.
type Options struct {
	Version   string
	Commit    string
	BuildDate string
}

// Server is the daemon's HTTP handler. Construct it with NewServer.
type Server struct {
	options Options
	mux     *http.ServeMux
}

// NewServer builds the daemon's HTTP handler with its fixed route table.
func NewServer(options Options) http.Handler {
	server := &Server{
		options: options,
		mux:     http.NewServeMux(),
	}

	server.mux.HandleFunc("GET /healthz", server.handleHealth)
	server.mux.HandleFunc("GET /v1/version", server.handleVersion)
	server.mux.HandleFunc("GET /v1/runtimes", server.handleRuntimes)

	return server
}

// ServeHTTP sets a fixed set of defensive response headers (no caching, a
// closed CSP, no MIME sniffing) before dispatching to the route table.
func (server *Server) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	response.Header().Set("Cache-Control", "no-store")
	response.Header().Set("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")
	response.Header().Set("X-Content-Type-Options", "nosniff")
	server.mux.ServeHTTP(response, request)
}

func (server *Server) handleHealth(response http.ResponseWriter, _ *http.Request) {
	writeJSON(response, http.StatusOK, map[string]string{"status": "ok"})
}

func (server *Server) handleVersion(response http.ResponseWriter, _ *http.Request) {
	writeJSON(response, http.StatusOK, map[string]string{
		"version":   server.options.Version,
		"commit":    server.options.Commit,
		"buildDate": server.options.BuildDate,
	})
}

// handleRuntimes reports a hardcoded placeholder profile. It is not yet
// backed by the profiles/ registry on disk (see profiles/bun-data/profile.json)
// — wiring that up is follow-up work, not a claim that discovery is implemented.
func (server *Server) handleRuntimes(response http.ResponseWriter, _ *http.Request) {
	type runtimeProfile struct {
		Name      string `json:"name"`
		Runtime   string `json:"runtime"`
		Status    string `json:"status"`
		Available bool   `json:"available"`
	}

	writeJSON(response, http.StatusOK, map[string][]runtimeProfile{
		"profiles": {
			{
				Name:      "bun-data@1",
				Runtime:   "bun",
				Status:    "draft",
				Available: false,
			},
		},
	})
}

func writeJSON(response http.ResponseWriter, status int, value any) {
	response.Header().Set("Content-Type", "application/json; charset=utf-8")
	response.WriteHeader(status)

	if err := json.NewEncoder(response).Encode(value); err != nil {
		return
	}
}
