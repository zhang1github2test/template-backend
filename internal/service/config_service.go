package service

import (
	"template-backend/internal/model"
	"template-backend/internal/repository"
)

type ConfigService interface {
	GetList(params map[string]interface{}) ([]model.Config, int64, error)
	GetByID(id int64) (*model.Config, error)
	Create(config *model.Config) error
	Update(config *model.Config) error
	Delete(id int64) error
}

type configService struct {
	repo repository.ConfigRepository
}

func NewConfigService(repo repository.ConfigRepository) ConfigService {
	return &configService{repo: repo}
}

func (s *configService) GetList(params map[string]interface{}) ([]model.Config, int64, error) {
	return s.repo.GetList(params)
}

func (s *configService) GetByID(id int64) (*model.Config, error) {
	return s.repo.GetByID(id)
}

func (s *configService) Create(config *model.Config) error {
	return s.repo.Create(config)
}

func (s *configService) Update(config *model.Config) error {
	return s.repo.Update(config)
}

func (s *configService) Delete(id int64) error {
	return s.repo.Delete(id)
}
