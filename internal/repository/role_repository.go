// repository/role_repository.go
package repository

import (
	"template-backend/internal/model"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository(db *gorm.DB) *RoleRepository {
	return &RoleRepository{db: db}
}

func (r *RoleRepository) GetList(page, size int, filters map[string]interface{}) ([]model.Role, int64, error) {
	var roles []model.Role
	var total int64

	query := r.db.Model(&model.Role{})

	if roleName, ok := filters["roleName"]; ok {
		query = query.Where("role_name LIKE ?", "%"+roleName.(string)+"%")
	}
	if roleCode, ok := filters["roleCode"]; ok {
		query = query.Where("role_code LIKE ?", "%"+roleCode.(string)+"%")
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * size).Limit(size).Find(&roles).Error
	return roles, total, err
}

func (r *RoleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}

func (r *RoleRepository) BatchDelete(ids []uint) error {
	return r.db.Delete(&model.Role{}, ids).Error
}

func (r *RoleRepository) GetByID(id uint) (*model.Role, error) {
	var role model.Role
	if err := r.db.First(&role, id).Error; err != nil {
		return nil, err
	}
	return &role, nil
}

func (r *RoleRepository) GetPermissions(roleID uint) ([]model.Permission, error) {
	var role model.Role
	if err := r.db.Preload("Permissions").First(&role, roleID).Error; err != nil {
		return nil, err
	}
	return role.Permissions, nil
}

func (r *RoleRepository) UpdatePermissions(roleID uint, permissionIds []uint) error {
	var role model.Role
	if err := r.db.First(&role, roleID).Error; err != nil {
		return err
	}

	var permissions []model.Permission
	if err := r.db.Where("id IN ?", permissionIds).Find(&permissions).Error; err != nil {
		return err
	}

	return r.db.Model(&role).Association("Permissions").Replace(permissions)
}
