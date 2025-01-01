package repository

import (
	"database/sql"
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
