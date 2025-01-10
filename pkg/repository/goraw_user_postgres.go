package repository

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/TimmyTurner98/Goraw/pkg/modules"
)

type GorawUserPostgres struct {
	db *sql.DB
}

func NewGorawUserPostgres(db *sql.DB) *GorawUserPostgres {
	return &GorawUserPostgres{db: db}
}

func (r *GorawUserPostgres) CreatUser(user modules.User) (int, error) {
	query := `INSERT INTO users (name, email, password) VALUES ($1, $2, $3) RETURNING id`
	var id int
	err := r.db.QueryRow(query, user.Name, user.Email, user.Password).Scan(&id)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func (r *GorawUserPostgres) GetUserByID(id int) (*modules.UserWithoutPassword, error) {
	query := `SELECT name, email FROM users WHERE id = $1`

	var user modules.UserWithoutPassword
	err := r.db.QueryRow(query, id).Scan(&user.Name, &user.Email)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("user with id %d not found: %w", id, err)
		}
		return nil, fmt.Errorf("error querying user by id: %w", err)
	}

	return &modules.UserWithoutPassword{user.Name, user.Email}, nil
}

func (r *GorawUserPostgres) DeleteUser(userID int) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user with id %d: %w", userID, err)
	}
	return nil
}
