package auth

import (
	"JWTGRPC/internal/domain/models"
	"JWTGRPC/internal/lib/jwt"
	"JWTGRPC/internal/lib/logger/sl"
	"JWTGRPC/internal/storage"
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"golang.org/x/crypto/bcrypt"
)

//===================================================================================================================//

// Пакет взаимодействия с БД и работой с Бизнес Логикой Сервиса
type Auth struct {
	log 		*slog.Logger
	usrSaver 	UserSaver
	usrProvider UserProvider
	appProvider AppProvider 
	tokenTTL 	time.Duration
}

type UserSaver interface {
	SaveUser(ctx context.Context, email string, passHach []byte) ( uid int64, err error)
}

type UserProvider interface {
	User(ctx context.Context, email string) (models.User, error)
}


type AppProvider interface {
	App(ctx context.Context, appID int) (models.App, error)
}

// Глобальная переменная ошибки
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrUserExist = errors.New("user already exists")
)



//===================================================================================================================//


// New возвращает новый истенс Аутх Сервиса
func New(
	log slog.Logger,
	userSaver UserSaver,
	userProvider UserProvider,
	appProvider AppProvider,
	tokenTTL time.Duration,
) *Auth {
	return &Auth{
		usrSaver: 		userSaver,
		usrProvider: 	userProvider,
		log: 			&log,
		appProvider: 	appProvider,
		tokenTTL: 		tokenTTL,
	}
}

//===================================================================================================================//

// Login checks if user with given credentials exists in the system and returns access token.
//
// If user exists, but password is incorrect, returns error.
// 
// If user doesn't exist, returns error. 
func (a *Auth) Login(ctx context.Context, email string,password string,appID int) (string, error) {
	const op = "auth.Login"

	// Иннициализация логгера
	log := a.log.With(
		slog.String("op:", op),
		slog.String("username:", email),
	)
	log.Info("Attempting to Login user")
	
	user, err := a.usrProvider.User(ctx, email)
	if err != nil {
		// Проверка на особую ошибку
		if errors.Is(err, storage.ErrUserNotFound) { // Паттерн для сравнивания (==) ошибок
			a.log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials) // позволяет "обернуть" одну ошибку в другую,
			// сохраняя оригинальную ошибку для последующего извлечения с помощью функции
		}

		// Лог на другие ошибки
		a.log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Проверка правильности пароля 
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		a.log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	// Проверка правильности приложения
	app, err := a.appProvider.App(ctx, appID)
	if err != nil {
		return "", fmt.Errorf("%s: %w", op, err) 
	}

	log.Info("User logged in successfully")

	token, err := jwt.NewToken(user, app, a.tokenTTL)
	if err != nil {
		a.log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	return token, nil
}


// RegisterNewUser registers new user in the system and returns user ID.
// If user with given username already exists, returns error.
func (a *Auth) RegisterNewUser(ctx context. Context, email string, pass string) (int64, error) {
	const op = "auth.RegisterNewUser"

	// Иннициализация логгера
	log := a.log.With(
		slog.String("op:", op),
		slog.String("email:", email),
	)
	log.Info("Registering user")

	// Запись пароля - хеширование
	passHach, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost) // Хеширование данных
	if err != nil {
		log.Error("failed to generation hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Сохранение данных пользователя по 3-м параметрам
	id, err := a.usrSaver.SaveUser(ctx, email, passHach)
	if err != nil {
		if errors.Is(err, storage.ErrUserExist) {
			log.Warn("user already exists", sl.Err(err))
			return 0, fmt.Errorf("%s: %w", op, ErrUserExist)
		}

		log.Error("Failed to safe user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	log.Info("User was registered")
	return id, nil
}
