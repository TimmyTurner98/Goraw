package service

import (
	"github.com/TimmyTurner98/Goraw/pkg/modules"
	"github.com/TimmyTurner98/Goraw/pkg/repository"
)

type UserService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) CreateUser(user modules.User) (int, error) {
	// Хэшируем пароль
	user.Password = generatePasswordHash(user.Password)
	// Передаем обработанные данные в репозиторий
	return s.repo.CreatUser(user)
}

func (s *UserService) GetUserByID(id int) (*modules.UserWithoutPassword, error) {
	return s.repo.GetUserByID(id)
}

func (s *UserService) GetAllUsers() ([]modules.UserWithoutPassword, error) {
	return s.repo.GetAllUsers()
}

func (s *UserService) DeleteUser(userID int) error {
	// Можно добавить дополнительную логику (например, проверку, существует ли пользователь)
	return s.repo.DeleteUser(userID)
}
