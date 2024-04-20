package app

import (
	"log/slog"
	grpcapp "myproject/internal/app/grpc"
	"myproject/internal/config"
	"myproject/internal/repositories/app"
	"myproject/internal/services/auth"
	"myproject/internal/services/orchestrator"
	"time"
)

type App struct {
	GRPCServer *grpcapp.App
}

func New(
	log *slog.Logger,
	orchestrator orchestrator.IOrchestrator,
	appRepo app.Repository,
	auth auth.IOAuth,
	grpcPort int,
	timeouts config.CalculationTimeoutsConfig,
	tokenTTL time.Duration,
) *App {
	grpcServer := grpcapp.New(log, auth, orchestrator, appRepo, grpcPort, timeouts)
	return &App{
		GRPCServer: grpcServer,
	}
}
