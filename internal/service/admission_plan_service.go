package service

import (
	"go.uber.org/zap"
	"template-backend/internal/model"
	"template-backend/internal/repository"
	"template-backend/pkg/logger"
)

type AdmissionPlanService struct {
	repo *repository.AdmissionPlanRepo
}

func NewAdmissionPlanService(repo *repository.AdmissionPlanRepo) *AdmissionPlanService {
	return &AdmissionPlanService{repo: repo}
}

func (s *AdmissionPlanService) List(page, pageSize int, filters map[string]interface{}) ([]model.HighSchoolAdmissionPlan, int64, error) {
	logger.Logger().Info("List 服务层调用", zap.Int("page", page), zap.Int("pageSize", pageSize), zap.Any("filters", filters))
	return s.repo.List(page, pageSize, filters)
}

func (s *AdmissionPlanService) GetByID(id int) (*model.HighSchoolAdmissionPlan, error) {
	logger.Logger().Info("GetByID 服务层调用", zap.Int("id", id))
	return s.repo.GetByID(id)
}

func (s *AdmissionPlanService) Create(plan *model.HighSchoolAdmissionPlan) error {
	logger.Logger().Info("Create 服务层调用", zap.Any("plan", plan))
	return s.repo.Create(plan)
}

func (s *AdmissionPlanService) Update(id int, plan *model.HighSchoolAdmissionPlan) error {
	logger.Logger().Info("Update 服务层调用", zap.Int("id", id), zap.Any("plan", plan))
	return s.repo.Update(id, plan)
}

func (s *AdmissionPlanService) Delete(id int) error {
	logger.Logger().Info("Delete 服务层调用", zap.Int("id", id))
	return s.repo.Delete(id)
}
