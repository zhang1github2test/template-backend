package service

import (
	"template-backend/internal/dto"
	"template-backend/internal/model"
	"template-backend/internal/repository"
)

type SchoolAdmissionService interface {
	Create(info *model.SchoolAdmissionInfo) error
	GetByID(id int) (*model.SchoolAdmissionInfo, error)
	List(req *dto.SchoolAdmissionQueryRequest) ([]model.SchoolAdmissionInfo, int64, error)
	Update(info *model.SchoolAdmissionInfo) error
	Delete(id int) error
}

type schoolAdmissionService struct {
	repo repository.SchoolAdmissionRepository
}

func NewSchoolAdmissionService(repo repository.SchoolAdmissionRepository) SchoolAdmissionService {
	return &schoolAdmissionService{repo: repo}
}

func (s *schoolAdmissionService) Create(info *model.SchoolAdmissionInfo) error {
	return s.repo.Create(info)
}

func (s *schoolAdmissionService) GetByID(id int) (*model.SchoolAdmissionInfo, error) {
	return s.repo.GetByID(id)
}

func (s *schoolAdmissionService) List(req *dto.SchoolAdmissionQueryRequest) ([]model.SchoolAdmissionInfo, int64, error) {
	return s.repo.List(req)
}

func (s *schoolAdmissionService) Update(info *model.SchoolAdmissionInfo) error {
	return s.repo.Update(info)
}

func (s *schoolAdmissionService) Delete(id int) error {
	return s.repo.Delete(id)
}
