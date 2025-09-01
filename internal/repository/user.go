package repository

import (
	"template-backend/internal/model"

	"gorm.io/gorm"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// 查询列表（带分页和筛选）
func (d *UserRepository) GetList(page, pageSize int, filters map[string]interface{}) ([]model.User, int64, error) {
	var users []model.User
	var total int64
	query := d.db.Model(&model.User{})

	// 动态条件
	if username, ok := filters["username"]; ok {
		query = query.Where("username LIKE ?", "%"+username.(string)+"%")
	}
	if nickname, ok := filters["nickname"]; ok {
		query = query.Where("nickname LIKE ?", "%"+nickname.(string)+"%")
	}
	if email, ok := filters["email"]; ok {
		query = query.Where("email LIKE ?", "%"+email.(string)+"%")
	}
	if status, ok := filters["status"]; ok {
		query = query.Where("status = ?", status)
	}

	err := query.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&users).Error
	return users, total, err
}

func (d *UserRepository) GetByID(id uint) (*model.User, error) {
	var user model.User
	if err := d.db.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (d *UserRepository) Create(user *model.User) error {
	return d.db.Create(user).Error
}

func (d *UserRepository) Update(user *model.User) error {
	return d.db.Save(user).Error
}

func (d *UserRepository) Delete(id uint) error {
	return d.db.Delete(&model.User{}, id).Error
}

// 为用户分配角色
func (d *UserRepository) AssignRoles(userID uint, roleIDs []uint) error {
	// 先删除用户现有的所有角色
	if err := d.db.Where("user_id = ?", userID).Delete(&model.UserRole{}).Error; err != nil {
		return err
	}

	// 添加新的角色关联
	if len(roleIDs) > 0 {
		var userRoles []model.UserRole
		for _, roleID := range roleIDs {
			userRoles = append(userRoles, model.UserRole{
				UserID: userID,
				RoleID: roleID,
			})
		}
		return d.db.Create(&userRoles).Error
	}

	return nil
}

// 获取用户的角色
func (d *UserRepository) GetUserRoles(userID uint) ([]model.Role, error) {
	var user model.User
	if err := d.db.Preload("Roles").First(&user, userID).Error; err != nil {
		return nil, err
	}
	return user.Roles, nil
}
