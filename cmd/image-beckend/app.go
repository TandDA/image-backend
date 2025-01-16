package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/TandDA/image-beckend/internal/handler"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "file:mydatabase.db?cache=shared&_fk=1")
	if err != nil {
		log.Fatal(err)
	}
	doMigration(db)
	hndlr := handler.NewHandler(db)
	hndlr.Start()
	fmt.Scanf("h")
}

func doMigration(db *sql.DB) {
	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		log.Fatalf("Не удалось создать экземпляр драйвера: %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://migrations", // Путь к директории с файлами миграции.1
		"sqlite3",           // Имя драйвера базы данных. 2
		driver,
	)
	if err != nil {
		log.Fatalf("Не удалось создать экземпляр миграции: %v", err)
	}

	// Применяем миграции.
	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		log.Fatalf("Ошибка при выполнении миграции: %v", err)
	}

	fmt.Println("Миграция успешно применена")
}
