package repository

import (
	"gorm.io/gorm"
	"template-backend/internal/model"
)

type ResourceRepository interface {
	Create(resource *model.Resource) error
	GetByID(id int64) (*model.Resource, error)
	GetByPermissionCode(code string) (*model.Resource, error)
	Update(id int64, updates map[string]interface{}) error
	Delete(id int64) error
	List(query *ResourceQuery) ([]model.Resource, int64, error)
	ExistsByPermissionCode(code string, excludeID int64) bool
	QueryAll(query *ResourceQuery) ([]model.Resource, error)
}

type ResourceQuery struct {
	ID             int64
	ResourceName   string
	PermissionCode string
	Type           string
	ResourcePath   string
	HTTPMethod     string
	ParentID       *int64
	Status         *int8
	RequiresAuth   *int8
	Page           int
	PageSize       int
}

type resourceRepository struct {
	db *gorm.DB
}

func NewResourceRepository(db *gorm.DB) ResourceRepository {
	return &resourceRepository{db: db}
}

func (r *resourceRepository) Create(resource *model.Resource) error {
	return r.db.Create(resource).Error
}

func (r *resourceRepository) GetByID(id int64) (*model.Resource, error) {
	var resource model.Resource
	err := r.db.Where("id = ?", id).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *resourceRepository) GetByPermissionCode(code string) (*model.Resource, error) {
	var resource model.Resource
	err := r.db.Where("permission_code = ?", code).First(&resource).Error
	if err != nil {
		return nil, err
	}
	return &resource, nil
}

func (r *resourceRepository) Update(id int64, updates map[string]interface{}) error {
	return r.db.Model(&model.Resource{}).Where("id = ?", id).Updates(updates).Error
}

func (r *resourceRepository) Delete(id int64) error {
	return r.db.Where("id = ?", id).Delete(&model.Resource{}).Error
}

func (r *resourceRepository) List(query *ResourceQuery) ([]model.Resource, int64, error) {
	var resources []model.Resource
	var total int64

	db := r.db.Model(&model.Resource{})

	// 构建查询条件
	if query.ID != 0 {
		db = db.Where("id = ?", query.ID)
	}
	if query.ResourceName != "" {
		db = db.Where("resource_name LIKE ?", "%"+query.ResourceName+"%")
	}
	if query.PermissionCode != "" {
		db = db.Where("permission_code LIKE ?", "%"+query.PermissionCode+"%")
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.ResourcePath != "" {
		db = db.Where("resource_path LIKE ?", "%"+query.ResourcePath+"%")
	}
	if query.HTTPMethod != "" {
		db = db.Where("http_method = ?", query.HTTPMethod)
	}
	if query.ParentID != nil {
		db = db.Where("parent_id = ?", *query.ParentID)
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.RequiresAuth != nil {
		db = db.Where("requires_auth = ?", *query.RequiresAuth)
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (query.Page - 1) * query.PageSize
	err := db.Order("sort ASC, id DESC").
		Offset(offset).
		Limit(query.PageSize).
		Find(&resources).Error

	return resources, total, err
}

func (r *resourceRepository) QueryAll(query *ResourceQuery) ([]model.Resource, error) {
	var resources []model.Resource

	db := r.db.Model(&model.Resource{})

	// 构建查询条件
	if query.ID != 0 {
		db = db.Where("id = ?", query.ID)
	}
	if query.ResourceName != "" {
		db = db.Where("resource_name LIKE ?", "%"+query.ResourceName+"%")
	}
	if query.PermissionCode != "" {
		db = db.Where("permission_code LIKE ?", "%"+query.PermissionCode+"%")
	}
	if query.Type != "" {
		db = db.Where("type = ?", query.Type)
	}
	if query.ResourcePath != "" {
		db = db.Where("resource_path LIKE ?", "%"+query.ResourcePath+"%")
	}
	if query.HTTPMethod != "" {
		db = db.Where("http_method = ?", query.HTTPMethod)
	}
	if query.ParentID != nil {
		db = db.Where("parent_id = ?", *query.ParentID)
	} else {
		db = db.Where("parent_id IS NULL")
	}
	if query.Status != nil {
		db = db.Where("status = ?", *query.Status)
	}
	if query.RequiresAuth != nil {
		db = db.Where("requires_auth = ?", *query.RequiresAuth)
	}

	err := db.Order("sort ASC, id DESC").
		Find(&resources).Error

	return resources, err
}

func (r *resourceRepository) ExistsByPermissionCode(code string, excludeID int64) bool {
	var count int64
	query := r.db.Model(&model.Resource{}).Where("permission_code = ?", code)
	if excludeID > 0 {
		query = query.Where("id != ?", excludeID)
	}
	query.Count(&count)
	return count > 0
}
