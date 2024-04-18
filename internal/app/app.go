package app

import (
	"log/slog"
	grpcapp "myproject/internal/app/grpc"
	"myproject/internal/app/httpapp"
	"myproject/internal/config"
	"myproject/internal/repositories/app"
	"myproject/internal/services/auth"
	"myproject/internal/services/orchestrator"
	"time"
)

type App struct {
	ServerHTTP *httpapp.ServerHTTP
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	orchestrator orchestrator.IOrchestrator,
	appRepo app.Repository,
	auth auth.IOAuth,
	httpPort int,
	grpcPort int,
	timeouts config.CalculationTimeoutsConfig,
	tokenTTL time.Duration,
) *App {
	serverHTTP := httpapp.New(log, orchestrator, httpPort, timeouts)
	grpcServer := grpcapp.New(log, auth, orchestrator, appRepo, grpcPort, timeouts)
	return &App{
		ServerHTTP: serverHTTP,
		GRPCServer: grpcServer,
	}
}
