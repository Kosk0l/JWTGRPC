package suite

import (
	"JWTGRPC/internal/config"
	"testing"

	ssov1 "github.com/Kosk0l/Protos/gen/go/sso"
)

type suite struct {
	*testing.T // Потребуется для вызова методов *testing.T внутри suite
	Cfg *config.Config // Конфигурация приложения
	AuthClient ssov1.AuthClient // Клиент для взаимодействия с grpc-сервером
}

