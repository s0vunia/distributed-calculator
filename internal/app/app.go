package app

import (
	"log/slog"
	"myproject/internal/app/httpapp"
	"myproject/internal/services/orchestrator"
)

type App struct {
	ServerHTTP *httpapp.ServerHTTP
}

func New(
	log *slog.Logger,
	orchestrator orchestrator.IOrchestrator,
	httpPort int,
) *App {
	serverHTTP := httpapp.New(log, orchestrator, httpPort)

	return &App{
		ServerHTTP: serverHTTP,
	}
}
