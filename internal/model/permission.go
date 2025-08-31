package model

type Permission struct {
	ID        uint         `gorm:"primaryKey" json:"id"`
	Name      string       `gorm:"size:64;not null" json:"name"`
	Code      string       `gorm:"size:64;uniqueIndex;not null" json:"code"`
	Type      int          `gorm:"not null" json:"type"` // 1: 菜单, 2: 按钮
	ParentID  uint         `json:"parentId"`
	Path      string       `gorm:"size:128" json:"path"`
	Component string       `gorm:"size:128" json:"component"`
	Icon      string       `gorm:"size:64" json:"icon"`
	Children  []Permission `gorm:"-" json:"children,omitempty"`
}
