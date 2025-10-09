package main

import (
	"errors"
	"flag"
	"fmt"

	// Драйвер миграции основной
	"github.com/golang-migrate/migrate/v4"

	// Драйвер для получения миграции из файла
	_ "github.com/golang-migrate/migrate/v4/source/file"

	// Драйвер для выполнения миграция sqlite3
	_ "github.com/golang-migrate/migrate/v4/database/sqlite3"
)

// go run ./cmd/migrator --storage-path=./storage/sso.db --migrations-path=./migrations
func main() {
	// Путь Бд // путь файлов Миграции // Имя таблицы, куда мигратор будет сохранять изменения
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to storages")
	flag.StringVar(&migrationsPath,"migrations-path", "", "path to migrations")
	flag.StringVar(&migrationsTable,"migrations-table", "migrations", "name of migration")
	flag.Parse()

	if storagePath == "" {
		panic("storage-path is required")
	}

	if migrationsPath == "" {
		panic("migration path is required")
	}

	m, err := migrate.New(
		"file://"+migrationsPath,
		fmt.Sprintf("sqlite3://%s?x-migrations-table=%s", storagePath, migrationsTable),
	)
	if err != nil {
		panic(err)
	}

	if err := m.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
			return
		}
		panic(err)
	}

	fmt.Println("migrations applied successfully")
}