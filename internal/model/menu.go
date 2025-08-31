// internal/model/menu.go
package model

import "encoding/json"

type Menu struct {
	ID         uint    `gorm:"primaryKey" json:"id"`
	Name       string  `json:"name"`
	Path       string  `json:"path"`
	Component  *string `json:"component"`
	ParentID   *uint   `json:"parentId"`
	Type       int     `json:"type"` // 1:目录 2:菜单 3:按钮
	Redirect   string  `json:"redirect"`
	Permission *string `json:"permission"`
	Visible    *bool   `json:"visible"`
	Sort       *int    `json:"sort"`
	Meta       Meta    `json:"meta" gorm:"-"`
	MetaJSON   string  `gorm:"type:json" json:"-"` // 存储 meta 结构体的 JSON

	Children []*Menu `gorm:"-" json:"children,omitempty"`
}

func (m *Menu) MarshalMeta() {
	// 将 Meta 转换为 JSON 字符串
	metaBytes, err := json.Marshal(m.Meta)
	if err != nil {
		// 处理错误
	}
	m.MetaJSON = string(metaBytes)
}
func (m *Menu) UnMarshalMeta() {
	var meta Meta
	// 将 Meta 转换为 JSON 字符串
	err := json.Unmarshal([]byte(m.MetaJSON), &meta)
	if err != nil {
		// 处理错误
	}
	m.Meta = meta
}

type Meta struct {
	Title        string   `json:"title"`
	RequiresAuth bool     `json:"requiresAuth"`
	Hidden       bool     `json:"hidden"`
	Icon         string   `json:"icon"`
	Permissions  []string `json:"permissions"`
	Roles        []string `json:"roles"`
	KeepAlive    bool     `json:"keepAlive"`
	Breadcrumb   bool     `json:"breadcrumb"`
	ActiveMenu   string   `json:"activeMenu"`
}
