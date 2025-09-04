package repository

import (
	"template-backend/internal/dto"
	"template-backend/internal/model"

	"gorm.io/gorm"
)

type SchoolAdmissionRepository interface {
	Create(info *model.SchoolAdmissionInfo) error
	GetByID(id int) (*model.SchoolAdmissionInfo, error)
	List(req *dto.SchoolAdmissionQueryRequest) ([]model.SchoolAdmissionInfo, int64, error)
	Update(info *model.SchoolAdmissionInfo) error
	Delete(id int) error
}

type schoolAdmissionRepository struct {
	db *gorm.DB
}

func NewSchoolAdmissionRepository(db *gorm.DB) SchoolAdmissionRepository {
	return &schoolAdmissionRepository{db: db}
}

func (r *schoolAdmissionRepository) Create(info *model.SchoolAdmissionInfo) error {
	return r.db.Create(info).Error
}

func (r *schoolAdmissionRepository) GetByID(id int) (*model.SchoolAdmissionInfo, error) {
	var info model.SchoolAdmissionInfo
	err := r.db.First(&info, id).Error
	return &info, err
}

// 分页 + 搜索查询
func (r *schoolAdmissionRepository) List(req *dto.SchoolAdmissionQueryRequest) ([]model.SchoolAdmissionInfo, int64, error) {
	var list []model.SchoolAdmissionInfo
	var total int64

	query := r.db.Model(&model.SchoolAdmissionInfo{})
	if req.SchoolName != "" {
		query = query.Where("school_name LIKE ?", "%"+req.SchoolName+"%")
	}
	if req.Year != 0 {
		query = query.Where("year = ?", req.Year)
	}
	if req.Category != "" {
		query = query.Where("category LIKE ?", "%"+req.Category+"%")
	}

	// 统计总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 分页查询
	offset := (req.Page - 1) * req.PageSize
	if err := query.Order("total_score DESC").Limit(req.PageSize).Offset(offset).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *schoolAdmissionRepository) Update(info *model.SchoolAdmissionInfo) error {
	return r.db.Save(info).Error
}

func (r *schoolAdmissionRepository) Delete(id int) error {
	return r.db.Delete(&model.SchoolAdmissionInfo{}, id).Error
}
