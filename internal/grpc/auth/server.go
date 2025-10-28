package auth

import (
	"JWTGRPC/internal/services/auth"
	"JWTGRPC/internal/storage"
	"context"
	"errors"

	ssov1 "github.com/Kosk0l/Protos/gen/go/sso"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

//===================================================================================================================//

// Интерфейс Авторизации
type Auth interface {
	Login(
		ctx context.Context,
		email string,
		password string,
		appID int,
	) ( token string, err error)

	RegisterNewUser(
		ctx context.Context,
		email string,
		password string,
	) (userID int64,	err error)
}

// Структура, наследующая интерфейс со всеми методами - хендлерами
type serverAPI struct {
	ssov1.UnimplementedAuthServer // насследование из прото файла
	auth Auth // ссылка на любой объект, который реализует интерфейс Auth.
	//хранит объект, который реализует интерфейс и к которому сервер делегирует работу.
}

// Регистрирует обработчки
func Register(gRPC *grpc.Server, auth Auth) {
	ssov1.RegisterAuthServer(gRPC, &serverAPI{auth: auth}) 
}

// Константы сервиса
const (
	emptyValue = 0 
)

//===================================================================================================================//

// Handler
func (s *serverAPI) Login( ctx context.Context, req *ssov1.LoginRepuest) (*ssov1.LoginResponse, error) {
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "bad Email")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Bad password")
	}

	if req.GetAppID() == emptyValue {
		return nil, status.Error(codes.InvalidArgument, "Bad AppID")
	}


	token, err := s.auth.Login(ctx, req.GetEmail(), req.GetPassword(), int(req.GetAppID()))
	if err != nil {
		if errors.Is(err, auth.ErrInvalidCredentials) {
			return nil, status.Error(codes.InvalidArgument, "invalid Argument")
		}
		return nil, status.Error(codes.Internal, "Internal error")
	}

	return &ssov1.LoginResponse{
		Token: token,
	}, nil
}

// Handler
func (s *serverAPI) Register( ctx context.Context, req *ssov1.RegisterRepuest) (*ssov1.RegisterResponse, error) {
	
	if req.GetEmail() == "" {
		return nil, status.Error(codes.InvalidArgument, "bad Email")
	}

	if req.GetPassword() == "" {
		return nil, status.Error(codes.InvalidArgument, "Bad password")
	}

	userID, err := s.auth.RegisterNewUser(ctx, req.GetEmail(), req.GetPassword())
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			return nil, status.Error(codes.AlreadyExists, "user already exists")
		}
		return nil, status.Error(codes.Internal,"Bad Internal")
	}

	return &ssov1.RegisterResponse{
		UserID: int64(userID),
	}, nil
}

