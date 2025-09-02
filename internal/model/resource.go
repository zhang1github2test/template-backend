package model

import "time"

type Resource struct {
	ID             int64     `json:"id" gorm:"primaryKey;autoIncrement;comment:资源ID，主键"`
	ResourceName   string    `json:"resource_name" gorm:"type:varchar(50);not null;comment:资源名称"`
	PermissionCode string    `json:"permission_code" gorm:"type:varchar(100);not null;uniqueIndex:uk_permission_code;comment:权限标识码，唯一标识"`
	Desc           *string   `json:"desc" gorm:"type:varchar(200);comment:资源描述"`
	Type           string    `json:"type" gorm:"type:varchar(20);not null;comment:资源类型：MENU-菜单，BUTTON-按钮，API-接口"`
	ResourcePath   *string   `json:"resource_path" gorm:"type:varchar(500);comment:资源路径（菜单路由或API路径）"`
	HTTPMethod     *string   `json:"http_method" gorm:"type:varchar(10);comment:HTTP方法（API类型使用）：GET,POST,PUT,DELETE等"`
	ParentID       *int64    `json:"parent_id" gorm:"comment:父权限ID（用于构建菜单树形结构）"`
	Sort           int       `json:"sort" gorm:"not null;default:0;comment:排序字段"`
	Status         int8      `json:"status" gorm:"not null;default:1;comment:权限状态：1-启用，0-禁用"`
	RequiresAuth   int8      `json:"requires_auth" gorm:"not null;default:1;comment:是否需要鉴权：1-需要，0-不需要"`
	Remark         *string   `json:"remark" gorm:"type:varchar(500);comment:备注信息"`
	CreatedBy      *int64    `json:"created_by" gorm:"comment:创建人ID"`
	CreatedAt      time.Time `json:"created_at" gorm:"not null;default:CURRENT_TIMESTAMP;comment:创建时间"`
	UpdatedBy      *int64    `json:"updated_by" gorm:"comment:更新人ID"`
	UpdatedAt      time.Time `json:"updated_at" gorm:"not null;default:CURRENT_TIMESTAMP;autoUpdateTime;comment:更新时间"`
}

func (Resource) TableName() string {
	return "resources"
}
