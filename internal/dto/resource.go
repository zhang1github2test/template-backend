package dto

import "time"

type CreateResourceRequest struct {
	ResourceName   string  `json:"resource_name" binding:"required" validate:"max=50"`
	PermissionCode string  `json:"permission_code" binding:"required" validate:"max=100"`
	Desc           *string `json:"desc" validate:"max=200"`
	Type           string  `json:"type" binding:"required" validate:"max=20"`
	ResourcePath   *string `json:"resource_path" validate:"max=500"`
	HTTPMethod     *string `json:"http_method" validate:"max=10"`
	ParentID       *int64  `json:"parent_id"`
	Sort           int     `json:"sort"`
	RequiresAuth   int8    `json:"requires_auth"`
	Remark         *string `json:"remark" validate:"max=500"`
	CreatedBy      *int64  `json:"created_by"`
}

type UpdateResourceRequest struct {
	ResourceName   *string `json:"resource_name" validate:"max=50"`
	PermissionCode *string `json:"permission_code" validate:"max=100"`
	Desc           *string `json:"desc" validate:"max=200"`
	Type           *string `json:"type" validate:"max=20"`
	ResourcePath   *string `json:"resource_path" validate:"max=500"`
	HTTPMethod     *string `json:"http_method" validate:"max=10"`
	ParentID       *int64  `json:"parent_id"`
	Sort           *int    `json:"sort"`
	Status         *int8   `json:"status"`
	RequiresAuth   *int8   `json:"requires_auth"`
	Remark         *string `json:"remark" validate:"max=500"`
	UpdatedBy      *int64  `json:"updated_by"`
}

type ResourceQueryRequest struct {
	ID             int64  `form:"id"`
	ResourceName   string `form:"resource_name"`
	PermissionCode string `form:"permission_code"`
	Type           string `form:"type"`
	ResourcePath   string `form:"resource_path"`
	HTTPMethod     string `form:"http_method"`
	ParentID       *int64 `form:"parent_id"`
	Status         *int8  `form:"status"`
	RequiresAuth   *int8  `form:"requires_auth"`
	Page           int    `form:"page" `
	PageSize       int    `form:"page_size"`
}

type ResourceResponse struct {
	ID             int64              `json:"id"`
	ResourceName   string             `json:"resource_name"`
	PermissionCode string             `json:"permission_code"`
	Desc           *string            `json:"desc"`
	Type           string             `json:"type"`
	ResourcePath   *string            `json:"resource_path"`
	HTTPMethod     *string            `json:"http_method"`
	ParentID       *int64             `json:"parent_id"`
	Sort           int                `json:"sort"`
	Status         int8               `json:"status"`
	RequiresAuth   int8               `json:"requires_auth"`
	Remark         *string            `json:"remark"`
	CreatedBy      *int64             `json:"created_by"`
	CreatedAt      time.Time          `json:"created_at"`
	UpdatedBy      *int64             `json:"updated_by"`
	UpdatedAt      time.Time          `json:"updated_at"`
	HasChildren    bool               `json:"has_children"`
	Children       []ResourceResponse `json:"children"`
}

type PagedResponse struct {
	Data       []ResourceResponse `json:"data"`
	Total      int64              `json:"total"`
	Page       int                `json:"page"`
	PageSize   int                `json:"page_size"`
	TotalPages int                `json:"total_pages"`
}
