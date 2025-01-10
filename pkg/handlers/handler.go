package handlers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/TimmyTurner98/Goraw/pkg/modules"
	"github.com/TimmyTurner98/Goraw/pkg/service"
	"net/http"
	"strconv"
)

type Handler struct {
	userService *service.UserService
}

func NewHandler(userService *service.UserService) *Handler {
	return &Handler{userService: userService}
}

func (h *Handler) CreateUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	var user modules.User
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	id, err := h.userService.CreateUser(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(fmt.Sprintf("User created with ID: %d", id)))
}

func (h *Handler) GetUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid method", http.StatusMethodNotAllowed)
		return
	}

	idStr := r.URL.Query().Get("id")
	id, err := strconv.Atoi(idStr) // Преобразуем строку в int
	if err != nil {
		http.Error(w, "Invalid ID format", http.StatusBadRequest)
		return
	}

	user, err := h.userService.GetUserByID(id) // передаем id в сервис
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "User not found", http.StatusNotFound) //400
		} else {
			http.Error(w, "Internal server error", http.StatusInternalServerError) //500
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user) // отправляем ответ в формате JSON
}

func (h *Handler) DeleteUserHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	userIDStr := r.URL.Query().Get("id")
	if userIDStr == "" {
		http.Error(w, "Missing user ID", http.StatusBadRequest)
		return
	}

	userID, err := strconv.Atoi(userIDStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	err = h.userService.DeleteUser(userID)
	if err != nil {
		http.Error(w, fmt.Sprintf("Error deleting user: %v", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("User with ID %d deleted successfully", userID)))
}
