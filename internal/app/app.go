package app

import (
	"log/slog"
	grpcapp "myproject/internal/app/grpc"
	"myproject/internal/app/httpapp"
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
	auth auth.IOAuth,
	httpPort int,
	grpcPort int,
	tokenTTL time.Duration,
) *App {
	serverHTTP := httpapp.New(log, orchestrator, httpPort)
	grpcServer := grpcapp.New(log, auth, grpcPort)
	return &App{
		ServerHTTP: serverHTTP,
		GRPCServer: grpcServer,
	}
}
