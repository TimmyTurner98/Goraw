package repository

import "github.com/TimmyTurner98/Goraw/pkg/modules"

type UserRepository interface {
	CreatUser(user modules.User) (int, error)
	DeleteUser(id int) error
	GetUserByID(id int) (*modules.UserWithoutPassword, error)
	GetAllUsers() ([]modules.UserWithoutPassword, error)
	UpdateUser(id int, user modules.User) error
}
