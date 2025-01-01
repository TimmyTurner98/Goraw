package main

import (
	"database/sql"
	goraw "github.com/TimmyTurner98/Goraw"
	"github.com/TimmyTurner98/Goraw/pkg/handlers"
	"github.com/TimmyTurner98/Goraw/pkg/repository"
	"github.com/TimmyTurner98/Goraw/pkg/service"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

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

	// Создание репозитория
	userRepo := repository.NewGorawUserPostgres(db)
	// Создание сервисов
	userService := service.NewUserService(userRepo)
	// Создание обработчиков
	handler := handlers.NewHandler(userService)

	// Регистрация обработчика для пути "/users"
	http.HandleFunc("/users", handler.CreateUserHandler)

	// Запуск сервера на порту 8443
	err = server.Run("8443", nil) // nil — это заглушка, так как маршруты уже обработаны через http.HandleFunc
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

}
