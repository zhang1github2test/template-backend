// pkg/repository/log.go
package repository

import (
	"gorm.io/gorm"
	"template-backend/internal/model"
	"time"
)

type LogRepository interface {
	Create(log *model.Log) error
	CreateInBatches(logs []*model.Log) error
	GetByID(id uint) (*model.Log, error)
	List(pageNum, pageSize int, conditions map[string]interface{}) ([]*model.Log, int64, error)
	Delete(id uint) error
	DeleteBatch(ids []uint) error
	Clean() error
}

type logRepository struct {
	db *gorm.DB
}

func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{db: db}
}

// BatchCreate 批量创建日志记录
func (r *logRepository) CreateInBatches(logs []*model.Log) error {
	// 如果日志列表为空，直接返回
	if len(logs) == 0 {
		return nil
	}
	return r.db.CreateInBatches(logs, 100).Error
}

func (r *logRepository) Create(log *model.Log) error {
	return r.db.Create(log).Error
}

func (r *logRepository) GetByID(id uint) (*model.Log, error) {
	var log model.Log
	err := r.db.Where("id = ?", id).First(&log).Error
	if err != nil {
		return nil, err
	}
	return &log, nil
}

func (r *logRepository) List(pageNum, pageSize int, conditions map[string]interface{}) ([]*model.Log, int64, error) {
	var logs []*model.Log
	var total int64

	db := r.db.Model(&model.Log{})

	// 添加查询条件
	if method, ok := conditions["method"]; ok && method != "" {
		db = db.Where("method = ?", method)
	}
	if path, ok := conditions["path"]; ok && path != "" {
		db = db.Where("path LIKE ?", "%"+path.(string)+"%")
	}
	if status, ok := conditions["status"]; ok && status != 0 {
		db = db.Where("status = ?", status)
	}
	if ip, ok := conditions["ip"]; ok && ip != "" {
		db = db.Where("ip = ?", ip)
	}
	if handler, ok := conditions["handler"]; ok && handler != "" {
		db = db.Where("handler LIKE ?", "%"+handler.(string)+"%")
	}
	if timestampRange, ok := conditions["timestamp"]; ok {
		if rangeArr, valid := timestampRange.([]time.Time); valid && len(rangeArr) == 2 {
			db = db.Where("timestamp BETWEEN ? AND ?", rangeArr[0], rangeArr[1])
		}
	}

	// 获取总数
	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (pageNum - 1) * pageSize
	err := db.Offset(offset).Limit(pageSize).Order("timestamp DESC").Find(&logs).Error
	if err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *logRepository) Delete(id uint) error {
	return r.db.Delete(&model.Log{}, id).Error
}

func (r *logRepository) DeleteBatch(ids []uint) error {
	return r.db.Where("id IN ?", ids).Delete(&model.Log{}).Error
}

func (r *logRepository) Clean() error {
	return r.db.Session(&gorm.Session{AllowGlobalUpdate: true}).Delete(&model.Log{}).Error
}
