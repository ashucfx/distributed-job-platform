package database

import "gorm.io/gorm"

type GormDBService struct {
	db *gorm.DB
}

func NewGormDBService(db *gorm.DB) DBService {
	return &GormDBService{db: db}
}

func (s *GormDBService) CreateJob(job *Job) error {
	return s.db.Create(job).Error
}

func (s *GormDBService) GetJobByID(id string) (*Job, error) {
	var job Job
	err := s.db.First(&job, "id = ?", id).Error
	return &job, err
}

func (s *GormDBService) UpdateJobStatus(id string, status JobStatus, errStr string) error {
	return s.db.Model(&Job{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":     status,
		"last_error": errStr,
	}).Error
}

func (s *GormDBService) IncrementRetries(id string) error {
	return s.db.Model(&Job{}).Where("id = ?", id).UpdateColumn("retries", gorm.Expr("retries + ?", 1)).Error
}

func (s *GormDBService) CreateJobLog(jobID string, message string, level string) error {
	return s.db.Create(&JobLog{
		JobID:   jobID,
		Message: message,
		Level:   level,
	}).Error
}
