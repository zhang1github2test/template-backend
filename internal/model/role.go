// model/role.go
package model

import "time"

type Role struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	RoleName  string    `gorm:"size:64;not null" json:"roleName"`
	RoleCode  string    `gorm:"size:64;uniqueIndex;not null" json:"roleCode"`
	RoleDesc  string    `gorm:"size:255" json:"roleDesc"`
	Status    int       `gorm:"default:1" json:"status"`
	CreatedAt time.Time `json:"createTime"`
	UpdatedAt time.Time `json:"updateTime"`
	// 多对多关系
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
}
