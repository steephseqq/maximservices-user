package app

import (
	"log/slog"
	"time"
	grpcapp "user/internal/app/grpc"
	"user/internal/services/user"
	"user/internal/storage/postgres"
)

type App struct {
	GRPCSrv *grpcapp.App
}

func New(
	log *slog.Logger,
	grpcPort int,
	tokenTTL time.Duration,
) *App {
	storage := postgres.MustLoad()

	userService := user.New(log, storage, storage, storage, tokenTTL)
	grpcApp := grpcapp.New(log, userService, grpcPort, tokenTTL)
	return &App{
		GRPCSrv: grpcApp,
	}
}
