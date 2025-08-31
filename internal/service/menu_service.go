// internal/service/menu_service.go
package service

import (
	"template-backend/internal/model"
	"template-backend/internal/repository"
)

type MenuService struct {
	repo *repository.MenuRepository
}

func NewMenuService(repo *repository.MenuRepository) *MenuService {
	return &MenuService{repo: repo}
}

func (s *MenuService) GetMenuTree(menu model.Menu) ([]*model.Menu, error) {
	return s.repo.GetMenuTree(menu)
}

func (s *MenuService) CreateMenu(menu *model.Menu) error {
	return s.repo.Create(menu)
}

func (s *MenuService) UpdateMenu(menu *model.Menu) error {
	return s.repo.Update(menu)
}

func (s *MenuService) DeleteMenu(id uint) error {
	return s.repo.Delete(id)
}

func (s *MenuService) GetByID(id uint) (*model.Menu, error) {
	return s.repo.GetByID(id)
}
