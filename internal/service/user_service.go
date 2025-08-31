// service/user_service.go
package service

import (
	"template-backend/internal/repository"

	"template-backend/internal/model"
)

type UserService struct {
	userDAO *repository.UserRepository
}

func NewUserService(userDAO *repository.UserRepository) *UserService {
	return &UserService{userDAO: userDAO}
}

func (s *UserService) GetList(page, pageSize int, filters map[string]interface{}) ([]model.User, int64, error) {
	return s.userDAO.GetList(page, pageSize, filters)
}

func (s *UserService) GetByID(id uint) (*model.User, error) {
	return s.userDAO.GetByID(id)
}

func (s *UserService) Create(user *model.User) error {
	return s.userDAO.Create(user)
}

func (s *UserService) Update(user *model.User) error {
	return s.userDAO.Update(user)
}

func (s *UserService) Delete(id uint) error {
	return s.userDAO.Delete(id)
}
