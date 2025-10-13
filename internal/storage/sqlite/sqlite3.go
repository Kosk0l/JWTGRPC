package sqlite

import (
	"JWTGRPC/internal/domain/models"
	"JWTGRPC/internal/storage"
	"context"
	"database/sql" // Общие правила для всех sql БД
	"errors"
	"fmt"

	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3" // Не используется напрямую, поэтому "_"
)

//===================================================================================================================//

type Storage struct {
	db *sql.DB // Коннект с бд
}


// Конструктор структуры Storage
func New(storagePath string) (*Storage, error) {
	const op = "storage.sqlite.New"

	// Указывает путь бд
	db1, err := sql.Open("sqlite3", storagePath)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	// Возвращает указатель на объект структуры с db параметром,
	// заполненным локальной db1 переменной
	return &Storage{db: db1}, nil 
}

//===================================================================================================================//


// Сохранить нового юзера в бд в БД
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.sqlite.SaveUser"

	// Запрос к бд
	stmt, err := s.db.Prepare("INSERT INTO users(email, pass_hash) VALUES(?, ?)")
	if err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	res, err := stmt.ExecContext(ctx, email, passHash)
	if err != nil {
		var sqliteErr sqlite3.Error

		if errors.As(err, &sqliteErr) && sqliteErr.ExtendedCode == sqlite3.ErrConstraintUnique {
			return 0, fmt.Errorf("%s: %w", op, storage.ErrUserExist)
		}

		return 0, fmt.Errorf("%s: %w", op, err)
	}

	// Получение id созданной записи
	id, err := res.LastInsertId()
	if err != nil {
			return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}


// Возвращает юзера по email данным
func (s *Storage) User(ctx context.Context, email string) (models.User, error) {
	const op = "storage.sqlite.User"

	// Подготавливаем запрос
	stmt, err := s.db.Prepare("SELECT id, email, pass_hash FROM users WHERE email = ?")
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	// Выполняем запрос
	row := stmt.QueryRowContext(ctx, email)

	// заполняем объект
	var user models.User
	err = row.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, storage.ErrUserNotFound)
		}
		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (s *Storage) App(ctx context.Context, id int) (models.App, error) {
	const op = "storage.sqlite.App"

	// Запрос к бд
	stmt, err := s.db.Prepare("SELECT id, name, secret FROM apps WHERE id = ?")
	if err != nil {
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	// Выполняем запрос
	row := stmt.QueryRowContext(ctx, id)

	// заполняем объект
	var app models.App
	err = row.Scan(&app.ID, &app.Name, &app.Secret)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.App{}, fmt.Errorf("%s: %w", op, storage.ErrAppNotFound)
		}
		return models.App{}, fmt.Errorf("%s: %w", op, err)
	}

	return app, nil
}