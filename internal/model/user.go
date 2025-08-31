package model

import "time"

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"uniqueIndex;size:64;not null" json:"username"`
	Nickname  string    `gorm:"size:64" json:"nickname"`
	Email     string    `gorm:"size:128" json:"email"`
	Phone     string    `gorm:"size:20" json:"phone"`
	Gender    string    `gorm:"size:10" json:"gender"`
	Status    int       `gorm:"default:1" json:"status"`
	Password  string    `gorm:"size:128" json:"-"` // 不返回密码
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
