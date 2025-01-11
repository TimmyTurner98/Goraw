package main

import (
	"database/sql"
	"fmt"
	goraw "github.com/TimmyTurner98/Goraw"
	"github.com/TimmyTurner98/Goraw/pkg/handlers"
	"github.com/TimmyTurner98/Goraw/pkg/repository"
	"github.com/TimmyTurner98/Goraw/pkg/service"
	_ "github.com/lib/pq"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
)

func createMigrationsTableIfNotExist(db *sql.DB) error {
	_, err := db.Exec(`
		CREATE TABLE IF NOT EXISTS migrations (
			id SERIAL PRIMARY KEY,
			name VARCHAR(255) NOT NULL,
			applied_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`)
	return err
}

func applyMigrations(db *sql.DB) error {
	// Получаем список всех файлов миграций в каталоге 'migrations'
	files, err := ioutil.ReadDir("migrations")
	if err != nil {
		return fmt.Errorf("failed to read migrations directory: %w", err)
	}

	// Загружаем уже примененные миграции
	appliedMigrations, err := getAppliedMigrations(db)
	if err != nil {
		return fmt.Errorf("failed to get applied migrations: %w", err)
	}

	// Применяем новые миграции
	for _, file := range files {
		if !file.IsDir() {
			migrationName := file.Name()
			if !isMigrationApplied(migrationName, appliedMigrations) {
				// Читаем содержимое файла миграции
				migrationContent, err := ioutil.ReadFile(filepath.Join("migrations", migrationName))
				if err != nil {
					return fmt.Errorf("failed to read migration file: %w", err)
				}

				// Выполняем миграцию
				if _, err := db.Exec(string(migrationContent)); err != nil {
					return fmt.Errorf("failed to apply migration %s: %w", migrationName, err)
				}

				// Записываем, что миграция была применена
				if err := markMigrationAsApplied(db, migrationName); err != nil {
					return fmt.Errorf("failed to mark migration as applied: %w", err)
				}
				log.Printf("Applied migration: %s", migrationName)
			}
		}
	}

	return nil
}

// Функция для получения всех примененных миграций
func getAppliedMigrations(db *sql.DB) ([]string, error) {
	rows, err := db.Query("SELECT name FROM migrations")
	if err != nil {
		return nil, fmt.Errorf("failed to get applied migrations: %w", err)
	}
	defer rows.Close()

	var migrations []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, fmt.Errorf("failed to scan migration row: %w", err)
		}
		migrations = append(migrations, name)
	}

	return migrations, nil
}

// Проверяем, была ли уже применена миграция
func isMigrationApplied(migrationName string, appliedMigrations []string) bool {
	for _, applied := range appliedMigrations {
		if migrationName == applied {
			return true
		}
	}
	return false
}

// Функция для записи миграции как примененной
func markMigrationAsApplied(db *sql.DB, migrationName string) error {
	_, err := db.Exec("INSERT INTO migrations (name) VALUES ($1)", migrationName)
	return err
}

func main() {
	server := &goraw.Server{}

	connStr := "postgres://postgres:qwerty@localhost:5437/gopgtest?sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	defer db.Close()

	if err != nil {
		log.Fatal(err, "Ошибка при запуске БД")
	}

	if err := db.Ping(); err != nil {
		log.Fatal(err, "Ошибка при соединение к БД")
	}
	log.Println("Подключение к базе данных успешно установлено")

	// Создаем таблицу миграций, если ее еще нет
	if err := createMigrationsTableIfNotExist(db); err != nil {
		log.Fatalf("Ошибка при создании таблицы миграций: %v", err)
	}

	// Применяем миграции
	if err := applyMigrations(db); err != nil {
		log.Fatalf("Ошибка при применении миграций: %v", err)
	}

	// Создание репозитория
	userRepo := repository.NewGorawUserPostgres(db)
	// Создание сервисов
	userService := service.NewUserService(userRepo)
	// Создание обработчиков
	handler := handlers.NewHandler(userService)

	// Регистрация обработчика для пути "/users"
	http.HandleFunc("/user", handler.CreateUserHandler)
	http.HandleFunc("/delete-user", handler.DeleteUserHandler)
	http.HandleFunc("/getuserbyid", handler.GetUserByID)
	http.HandleFunc("/getallusers", handler.GetAllUsers)

	// Запуск сервера на порту 8443
	err = server.Run("8443", nil) // nil — это заглушка, так как маршруты уже обработаны через http.HandleFunc
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}
}
