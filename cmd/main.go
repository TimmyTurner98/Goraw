package main

import (
	"database/sql"
	goraw "github.com/TimmyTurner98/Goraw"
	"github.com/TimmyTurner98/Goraw/pkg/handlers"
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

	http.HandleFunc("/fullname", handlers.FullnameHandler)
	http.HandleFunc("/timmy", handlers.TimmyHandler)
	http.HandleFunc("/", handlers.Handler)

	// Запуск сервера на порту 8443
	err = server.Run("8443", nil) // nil — это заглушка, так как маршруты уже обработаны через http.HandleFunc
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

}
