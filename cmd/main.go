package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

type Person struct {
	Name  string `json: "name"`
	Age   int    `json: "age"`
	Email string `json: "email"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello HTTPS World!")
}

func TimmyHandler(w http.ResponseWriter, r *http.Request) {
	// Создаем экземпляр структуры Person
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	person := Person{
		Name: "Timmy",
		Age:  26,
	}
	// Устанавливаем Content-Type как application/json
	w.Header().Set("Content-Type", "application/json")

	jsonData, err := json.MarshalIndent(person, "", "    ")
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}
	// Отправляем JSON-ответ
	w.Write(jsonData)
}

func fullnameHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	person := Person{
		Name: "Temirlan",
		Age:  26,
	}
	w.Header().Set("Content-Type", "application/json")
	jsonData, err := json.MarshalIndent(person, "", "    ")
	if err != nil {
		http.Error(w, "Unable to marshal JSON", http.StatusInternalServerError)
		return
	}
	w.Write(jsonData)
}

type Product struct {
	Name      string
	Price     float64
	Available bool
}

func main() {
	server := &http.Server{}

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

	http.HandleFunc("/fullname", fullnameHandler)
	http.HandleFunc("/timmy", TimmyHandler)
	http.HandleFunc("/", Handler)

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
