package repository

import "github.com/TimmyTurner98/Goraw/pkg/modules"

type UserRepository interface {
	CreatUser(user modules.User) (int, error)
	//GetUserByID(id int) (*modules.User, error)
	//GetAllUsers() ([]modules.User, error)
	//UpdateUser(user modules.User) error
	//DeleteUser(id int) error
}
