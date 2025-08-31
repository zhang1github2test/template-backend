package repository

import (
	"template-backend/internal/model"

	"gorm.io/gorm"
)

type ConfigRepository interface {
	GetList(params map[string]interface{}) ([]model.Config, int64, error)
	GetByID(id int64) (*model.Config, error)
	Create(config *model.Config) error
	Update(config *model.Config) error
	Delete(id int64) error
}

type configRepository struct {
	db *gorm.DB
}

func NewConfigRepository(db *gorm.DB) ConfigRepository {
	return &configRepository{db: db}
}

func (r *configRepository) GetList(params map[string]interface{}) ([]model.Config, int64, error) {
	var configs []model.Config
	var total int64

	query := r.db.Model(&model.Config{})

	// 动态条件
	if v, ok := params["configKey"].(string); ok && v != "" {
		query = query.Where("config_key LIKE ?", "%"+v+"%")
	}
	if v, ok := params["configName"].(string); ok && v != "" {
		query = query.Where("config_name LIKE ?", "%"+v+"%")
	}

	// 统计总数
	query.Count(&total)

	// 分页
	pageNum, _ := params["pageNum"].(int)
	pageSize, _ := params["pageSize"].(int)
	if pageNum > 0 && pageSize > 0 {
		offset := (pageNum - 1) * pageSize
		query = query.Offset(offset).Limit(pageSize)
	}

	err := query.Find(&configs).Error
	return configs, total, err
}

func (r *configRepository) GetByID(id int64) (*model.Config, error) {
	var config model.Config
	if err := r.db.First(&config, id).Error; err != nil {
		return nil, err
	}
	return &config, nil
}

func (r *configRepository) Create(config *model.Config) error {
	return r.db.Create(config).Error
}

func (r *configRepository) Update(config *model.Config) error {
	return r.db.Save(config).Error
}

func (r *configRepository) Delete(id int64) error {
	return r.db.Delete(&model.Config{}, id).Error
}
