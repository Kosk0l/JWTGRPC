package migrator

import (
	"errors"
	"flag"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
)

func main() {
	// Путь Бд // путь файлов Миграции // Имя таблицы, куда мигратор будет сохранять изменения
	var storagePath, migrationsPath, migrationsTable string

	flag.StringVar(&storagePath, "storage-path", "", "path to storages")
	flag.StringVar(&migrationsPath,"migratiions-path", "", "path to migrations")
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