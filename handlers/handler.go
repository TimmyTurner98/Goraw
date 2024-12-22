package handlers

import (
	"encoding/json"
	"fmt"
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
