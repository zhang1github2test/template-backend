// service/role_service.go
package service

import (
	"template-backend/internal/model"
	"template-backend/internal/repository"
)

type RoleService struct {
	roleRepo *repository.RoleRepository
}

func NewRoleService(roleRepo *repository.RoleRepository) *RoleService {
	return &RoleService{roleRepo: roleRepo}
}

func (s *RoleService) GetList(page, size int, filters map[string]interface{}) ([]model.Role, int64, error) {
	return s.roleRepo.GetList(page, size, filters)
}

func (s *RoleService) Create(role *model.Role) error {
	return s.roleRepo.Create(role)
}

func (s *RoleService) Update(role *model.Role) error {
	return s.roleRepo.Update(role)
}

func (s *RoleService) Delete(id uint) error {
	return s.roleRepo.Delete(id)
}

func (s *RoleService) BatchDelete(ids []uint) error {
	return s.roleRepo.BatchDelete(ids)
}

func (s *RoleService) GetByID(id uint) (*model.Role, error) {
	return s.roleRepo.GetByID(id)
}

func (s *RoleService) GetPermissions(roleID uint) ([]model.Permission, error) {
	return s.roleRepo.GetPermissions(roleID)
}

func (s *RoleService) UpdatePermissions(roleID uint, permissionIds []uint) error {
	return s.roleRepo.UpdatePermissions(roleID, permissionIds)
}
