package app

import (
	grpcapp "JWTGRPC/internal/app/grpc"
	"JWTGRPC/internal/services/auth"

	"JWTGRPC/internal/storage/sqlite"
	"log/slog"
	"time"
)

//===================================================================================================================//

type App struct {
	GRPCSrv *grpcapp.App
}

//===================================================================================================================//

// Создание нового gRPC Сервера
func New(
	log *slog.Logger,
	grpcPort int,
	storagePath string,
	tokenTTL time.Duration,
) *App {
	storage, err := sqlite.New(storagePath)
	if err != nil {
		panic(err)
	}

	authService := auth.New(*log, storage, storage, storage, tokenTTL)
	
	grpcApp := grpcapp.New(log, authService, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}