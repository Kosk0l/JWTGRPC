package grpcapp

import (
	authgrpc "JWTGRPC/internal/grpc/auth"
	"fmt"
	"log/slog"
	"net"

	"google.golang.org/grpc"
)

//
type App struct {
	log 		*slog.Logger
	gRPCServer 	*grpc.Server
	port 		int
}

// Создание нового gRPC Сервера
func New(
	log *slog.Logger,
	port int,
) *App {
	gRPCServer := grpc.NewServer() // Создание сервера
	authgrpc.Register(gRPCServer) // Подключение обработчика

	return &App {
		log: 		log,
		gRPCServer: gRPCServer,
		port: 		port,
	}
}

// Запуск Сервера
func (a *App) Run() error {
	const op = "grpcapp.Run" 

	log := a.log.With( // Логи
		slog.String("op", op),
		slog.Int("port", a.port),
	)

	l, err := net.Listen("tcp", fmt.Sprintf(":%d",a.port)) // Слушатель 
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	log.Info("gRPC server is running", slog.String("addr", l.Addr().String())) // Лог

	if err := a.gRPCServer.Serve(l); err != nil { // Коннект сервера
		return fmt.Errorf("%s: %w", op, err) // Возврат возможной ошибки
	}

	return nil
}

// Запускает gRPC server и паникует, если ошибки есть.
func (a *App) MustRun() {
	if err := a.Run(); err != nil {
		panic(err)
	}
}

// Функция остановки сервера
func (a *App) Stop() {
	const op = "grpcapp.Stop"
	a.log.With(slog.String("op", op)).Info("stopping gRPC server", slog.Int("port", a.port))
	a.gRPCServer.GracefulStop() // Остановка сервера 
	// Прекращает прием новых запросов и заканчивает полностью старые
}