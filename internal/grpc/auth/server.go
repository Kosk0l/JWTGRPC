package auth

import (
	"context"

	ssov1 "github.com/Kosk0l/Protos/gen/go/sso"
	"google.golang.org/grpc"
)

// Структура, наследующая интерфейс со всеми методами - хендлерами
type serverAPI struct {
	ssov1.UnimplementedAuthServer
}

// Регистрирует обработчки
func Register(gRPC *grpc.Server) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{}) 
}

// Handler
func (s *serverAPI) Login(
	ctx context.Context,
	req *ssov1.LoginRepuest,
) (*ssov1.LoginResponse, error) {
	return &ssov1.LoginResponse{
		Token: "token123",
	}, nil
}

// Handler
func (s *serverAPI) Register(
	ctx context.Context,
	req *ssov1.RegisterRepuest,
) (*ssov1.RegisterResponse, error) {
	panic("aboba")
}

