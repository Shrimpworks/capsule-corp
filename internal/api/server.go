package api

import (
	"encoding/json"
	"net/http"
)

type Options struct {
	Version   string
	Commit    string
	BuildDate string
}

type Server struct {
	options Options
	mux     *http.ServeMux
}

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
