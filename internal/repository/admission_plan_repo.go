package repository

import (
	"go.uber.org/zap"
	"gorm.io/gorm"
	"template-backend/internal/model"
	"template-backend/pkg/logger"
)

type AdmissionPlanRepo struct {
	db *gorm.DB
}

func NewAdmissionPlanRepo(db *gorm.DB) *AdmissionPlanRepo {
	return &AdmissionPlanRepo{db: db}
}

func (r *AdmissionPlanRepo) List(page, pageSize int, filters map[string]interface{}) ([]model.HighSchoolAdmissionPlan, int64, error) {
	var plans []model.HighSchoolAdmissionPlan
	var total int64

	query := r.db.Model(&model.HighSchoolAdmissionPlan{})

	for key, value := range filters {
		if value != "" {
			switch key {
			case "school_name":
				query = query.Where("school_name LIKE ?", "%"+value.(string)+"%")
			case "year":
				query = query.Where("year = ?", value)
			case "district_type":
				query = query.Where("district_type = ?", value)
			default:
				query = query.Where(key+" = ?", value)
			}
		}
	}

	if err := query.Count(&total).Error; err != nil {
		logger.Logger().Error("List Count 查询失败", zap.Error(err))
		return nil, 0, err
	}

	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&plans).Error; err != nil {
		logger.Logger().Error("List Find 查询失败", zap.Error(err))
		return nil, 0, err
	}

	logger.Logger().Info("List 查询成功", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.Int64("total", total))
	return plans, total, nil
}

func (r *AdmissionPlanRepo) GetByID(id int) (*model.HighSchoolAdmissionPlan, error) {
	var plan model.HighSchoolAdmissionPlan
	err := r.db.First(&plan, id).Error
	if err != nil {
		logger.Logger().Error("GetByID 查询失败", zap.Error(err), zap.Int("id", id))
		return nil, err
	}
	logger.Logger().Info("GetByID 查询成功", zap.Int("id", id))
	return &plan, nil
}

func (r *AdmissionPlanRepo) Create(plan *model.HighSchoolAdmissionPlan) error {
	err := r.db.Create(plan).Error
	if err != nil {
		logger.Logger().Error("Create 创建失败", zap.Error(err), zap.Any("plan", plan))
		return err
	}
	logger.Logger().Info("Create 创建成功", zap.Int("id", plan.ID))
	return nil
}

func (r *AdmissionPlanRepo) Update(id int, plan *model.HighSchoolAdmissionPlan) error {
	err := r.db.Model(&model.HighSchoolAdmissionPlan{}).Where("id = ?", id).Updates(plan).Error
	if err != nil {
		logger.Logger().Error("Update 更新失败", zap.Error(err), zap.Int("id", id))
		return err
	}
	logger.Logger().Info("Update 更新成功", zap.Int("id", id))
	return nil
}

func (r *AdmissionPlanRepo) Delete(id int) error {
	err := r.db.Delete(&model.HighSchoolAdmissionPlan{}, id).Error
	if err != nil {
		logger.Logger().Error("Delete 删除失败", zap.Error(err), zap.Int("id", id))
		return err
	}
	logger.Logger().Info("Delete 删除成功", zap.Int("id", id))
	return nil
}
