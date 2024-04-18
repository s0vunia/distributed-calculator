package httpapp

import (
	"context"
	log "github.com/sirupsen/logrus"
	"log/slog"
	"myproject/internal/config"
	"myproject/internal/services/orchestrator"
	"net/http"
	"strconv"
)

type ServerHTTP struct {
	server       *http.Server
	log          *slog.Logger
	orchestrator orchestrator.IOrchestrator
	port         int
	timeouts     config.CalculationTimeoutsConfig
}

// New creates new gRPC server app.
func New(
	log *slog.Logger,
	orchestrator orchestrator.IOrchestrator,
	port int,
	timeouts config.CalculationTimeoutsConfig,
) *ServerHTTP {

	server := &http.Server{Addr: ":" + strconv.Itoa(port)}

	return &ServerHTTP{
		server:       server,
		log:          log,
		orchestrator: orchestrator,
		port:         port,
		timeouts:     timeouts,
	}
}

// MustRun runs gRPC server and panics if any error occurs.
func (a *ServerHTTP) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Run runs gRPC server.
func (a *ServerHTTP) Run() error {
	mux := a.Routing()
	a.server.Handler = mux

	a.log.Info("Starting server on %s", a.port)
	err := a.server.ListenAndServe()
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
	return err
}

func (a *ServerHTTP) Routing() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/api/v0/expression/", orchestrator.EndpointExpression(a.orchestrator))
	mux.HandleFunc("/api/v0/expression", orchestrator.EndpointExpression(a.orchestrator))
	mux.HandleFunc("/api/v0/expressions", orchestrator.GetExpressions(a.orchestrator))
	mux.HandleFunc("/api/v0/sub_expressions", orchestrator.GetSubExpressions(a.orchestrator))
	mux.HandleFunc("/api/v0/agents", orchestrator.GetAgents(a.orchestrator))
	mux.HandleFunc("/api/v0/operators", orchestrator.GetOperators(a.timeouts))
	return mux
}

// Stop stops gRPC server.
func (a *ServerHTTP) Stop() {
	a.server.Shutdown(context.Background())
}
