package main

import (
	"JWTGRPC/internal/app"
	"JWTGRPC/internal/config"
	"JWTGRPC/internal/lib/logger/handlers/slogpretty"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

//===================================================================================================================//

const (
	envLocal = 	"local"
	envDev 	 = 	"dev"
	envProd  = 	"prod"
)

//===================================================================================================================//

func main() {
	// Иннициализация объекта конфига Микросервиса
	cfg := config.MustLoad()
	// go run cmd/sso/main.go --config=./config/local.yaml

	// Иннциализация логгера
	log := setupLogger(cfg.Env)
	log.Info("starting application", slog.String("env", cfg.Env))

	//Запуск gRPC сервер приложения
	application := app.New(log, cfg.GRPC.Port, cfg.StoragePath, cfg.TokenTTL)
	go application.GRPCSrv.MustRun()

	// TODO: Иннициализировать приложение (app)

	stop := make(chan os.Signal, 1) // блокирует горутину до прихода ответа
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)
	<- stop

	application.GRPCSrv.Stop() // Вручную прерываем сервер по нажатию клавиш
	log.Info("application stopped") 

}

//===================================================================================================================//

func setupLogger(env string) *slog.Logger {
	var log *slog.Logger // Указатель на Структуру Logger

	switch env {
	case envLocal:

		log = setupPrettySlog()

		// log = slog.New(
		// 	// Debug -> info -> Warm -> Error
		// 	slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		// ) 
	case envDev: 
		log = slog.New(
			// Debug -> info -> Warm -> Error
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}),
		) 
	case envProd:
		log = slog.New(
			// Debug -> info -> Warm -> Error
			slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}),
		)
	}

	return log
}

func setupPrettySlog() *slog.Logger {
	opts := slogpretty.PrettyHandlerOptions{
		SlogOpts: &slog.HandlerOptions{
			Level: slog.LevelDebug,
		},
	}

	handler := opts.NewPrettyHandler(os.Stdout)

	return slog.New(handler)
}