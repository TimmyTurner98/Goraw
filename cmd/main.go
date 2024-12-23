package main

import (
	"database/sql"
	"fmt"
	goraw "github.com/TimmyTurner98/Goraw"
	"github.com/TimmyTurner98/Goraw/pkg/handlers"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Product struct {
	Name      string
	Price     float64
	Available bool
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

	createProductTable(db)

	data := []Product{}
	rows, err := db.Query("SELECT name, available, price FROM product")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()
	// to scan DB values
	var name string
	var available bool
	var price float64

	for rows.Next() {
		err = rows.Scan(&name, &available, &price)
		if err != nil {
			log.Fatal(err)
		}
		data = append(data, Product{name, price, available})
	}

	fmt.Println(data)
	/*
		product := Product{"Book", 15.55, true}
		pk := insertProduct(db, product)
		//___________________________________________________
		var name string
		var available bool
		var price float64

		query := "SELECT name, available, price FROM product WHERE id = $1"
		err = db.QueryRow(query, pk).Scan(&name, &available, &price)
		if err != nil {
			log.Fatal(err)

	*/

	//___________________________________________________________
	fmt.Printf("Name: %s\n", name)
	fmt.Printf("Price: %t\n", available)
	fmt.Printf("Price: %f\n", price)

	http.HandleFunc("/fullname", handlers.FullnameHandler)
	http.HandleFunc("/timmy", handlers.TimmyHandler)
	http.HandleFunc("/", handlers.Handler)

	// Запуск сервера на порту 8443
	err = server.Run("8443", nil) // nil — это заглушка, так как маршруты уже обработаны через http.HandleFunc
	if err != nil {
		log.Fatalf("Ошибка при запуске сервера: %v", err)
	}

}

func createProductTable(db *sql.DB) {
	query := `CREATE TABLE IF NOT EXISTS product (
    	Id SERIAL PRIMARY KEY,
    	name VARCHAR(100) NOT NULL,
    	price NUMERIC(6,2) NOT NULL,
    	available BOOLEAN,
    	created_at timestamp DEFAULT NOW()
        )`
	_, err := db.Exec(query)
	if err != nil {
		log.Fatal(err)
	}
}

func insertProduct(db *sql.DB, product Product) int {
	query := `INSERT INTO product (name, price, available) VALUES ($1, $2, $3) returning id`

	var pk int
	err := db.QueryRow(query, product.Name, product.Price, product.Available).Scan(&pk)
	if err != nil {
		log.Fatal(err)
	}
	return pk
}
