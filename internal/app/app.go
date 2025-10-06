package app

import (
	grpcapp "JWTGRPC/internal/app/grpc"
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
	//TODO: Иннициализировать хранилище
	
	//TODO: init auth service(auth)
	
	grpcApp := grpcapp.New(log, grpcPort)

	return &App{
		GRPCSrv: grpcApp,
	}
}