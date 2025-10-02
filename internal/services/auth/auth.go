package auth

import (
	"context"
	"log/slog"
)

// Пакет взаимодействия с БД и работой с Бизнес Логикой Сервиса

type Auth struct {
	log 	*slog.Logger
}

type UserSaver interface {
	SaveUser(
		ctx context.Context,
		email string,
		passHach []byte,
	) (
		uid int64,
		err error,
	)
}