package service

import (
	"go.uber.org/zap"
	"template-backend/internal/global"
	"template-backend/internal/model"
	"template-backend/internal/repository"
	"template-backend/pkg/logger"
	"time"
)

type LogService interface {
	CreateLog(log *model.Log) error
	GetLogByID(id uint) (*model.Log, error)
	GetLogList(pageNum, pageSize int, conditions map[string]interface{}) ([]model.Log, int64, error)
	DeleteLog(id uint) error
	DeleteLogs(ids []uint) error
	CleanLogs() error
}

type logService struct {
	repo repository.LogRepository
}

func NewLogService(repo repository.LogRepository) LogService {
	service := &logService{repo: repo}
	global.InitLogChan(10000)
	service.startLogConsumer(100)
	return service
}
func (s *logService) CreateInBatches(logs []*model.Log) error {
	logger.Logger().Info("Creating log entry in batches", zap.Any("log", logs))
	err := s.repo.CreateInBatches(logs)
	if err != nil {
		logger.Logger().Error("Failed to create log entry", zap.Error(err))
		return err
	}
	logger.Logger().Info("Successfully created log entry")
	return nil
}

func (s *logService) CreateLog(log *model.Log) error {
	logger.Logger().Info("Creating log entry", zap.Any("log", log))
	err := s.repo.Create(log)
	if err != nil {
		logger.Logger().Error("Failed to create log entry", zap.Error(err))
		return err
	}
	logger.Logger().Info("Successfully created log entry", zap.Uint("id", log.ID))
	return nil
}

func (s *logService) GetLogByID(id uint) (*model.Log, error) {
	logger.Logger().Info("Fetching log by ID", zap.Uint("id", id))
	log, err := s.repo.GetByID(id)
	if err != nil {
		logger.Logger().Error("Failed to fetch log by ID", zap.Uint("id", id), zap.Error(err))
		return nil, err
	}
	logger.Logger().Info("Successfully fetched log by ID", zap.Uint("id", id))
	return log, nil
}

func (s *logService) GetLogList(pageNum, pageSize int, conditions map[string]interface{}) ([]model.Log, int64, error) {
	logger.Logger().Info("Fetching log list",
		zap.Int("pageNum", pageNum),
		zap.Int("pageSize", pageSize),
		zap.Any("conditions", conditions))

	logs, total, err := s.repo.List(pageNum, pageSize, conditions)
	if err != nil {
		logger.Logger().Error("Failed to fetch log list", zap.Error(err))
		return nil, 0, err
	}

	logger.Logger().Info("Successfully fetched log list",
		zap.Int("count", len(logs)),
		zap.Int64("total", total))
	return logs, total, nil
}

func (s *logService) DeleteLog(id uint) error {
	logger.Logger().Info("Deleting log", zap.Uint("id", id))
	err := s.repo.Delete(id)
	if err != nil {
		logger.Logger().Error("Failed to delete log", zap.Uint("id", id), zap.Error(err))
		return err
	}
	logger.Logger().Info("Successfully deleted log", zap.Uint("id", id))
	return nil
}

func (s *logService) DeleteLogs(ids []uint) error {
	logger.Logger().Info("Batch deleting logs", zap.Any("ids", ids))
	err := s.repo.DeleteBatch(ids)
	if err != nil {
		logger.Logger().Error("Failed to batch delete logs", zap.Any("ids", ids), zap.Error(err))
		return err
	}
	logger.Logger().Info("Successfully batch deleted logs", zap.Int("count", len(ids)))
	return nil
}

func (s *logService) CleanLogs() error {
	logger.Logger().Info("Cleaning all logs")
	err := s.repo.Clean()
	if err != nil {
		logger.Logger().Error("Failed to clean logs", zap.Error(err))
		return err
	}
	logger.Logger().Info("Successfully cleaned all logs")
	return nil
}

func (s *logService) startLogConsumer(batchSize int) {
	go func() {
		logBuffer := make([]*model.Log, 0, batchSize)
		ticker := time.NewTicker(5 * time.Second)
		defer ticker.Stop()

		for {
			select {
			case logEntry, ok := <-global.GetLogChan():
				if !ok {
					// 通道已关闭，处理剩余日志并退出
					if len(logBuffer) > 0 {
						s.CreateInBatches(logBuffer)
					}
					return
				}

				logBuffer = append(logBuffer, logEntry)

				// 达到批次大小时批量处理
				if len(logBuffer) >= batchSize {
					if err := s.CreateInBatches(logBuffer); err != nil {
						logger.Logger().Error("Failed to batch create logs", zap.Error(err))
					}
					logBuffer = make([]*model.Log, 0, batchSize)
				}

			case <-ticker.C:
				// 定时刷新缓冲区
				if len(logBuffer) > 0 {
					if err := s.CreateInBatches(logBuffer); err != nil {
						logger.Logger().Error("Failed to batch create logs", zap.Error(err))
					}
					logBuffer = make([]*model.Log, 0, 100)
				}
			}
		}
	}()
}
