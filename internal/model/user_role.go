// internal/model/user_role.go
package model

// 用户和角色的关联表
type UserRole struct {
	UserID uint `gorm:"primaryKey"`
	RoleID uint `gorm:"primaryKey"`
}
