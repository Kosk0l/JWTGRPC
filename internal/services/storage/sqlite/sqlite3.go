package sqlite

import (
	"JWTGRPC/internal/services/storage"
	"context"
	"database/sql" // Общие правила для всех sql БД
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	_ "github.com/mattn/go-sqlite3" // Не используется напрямую, поэтому "_"
)

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


// Сохранить нового юзера в бд в БД
func (s *Storage) SaveUser(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "storage.sqlite.SaveuUer"

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

