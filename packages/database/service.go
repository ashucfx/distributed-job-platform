package database

type DBService interface {
	CreateJob(job *Job) error
	GetJobByID(id string) (*Job, error)
	ListJobs(limit int) ([]Job, error)
	UpdateJobStatus(id string, status JobStatus, errStr string) error
	IncrementRetries(id string) error
	CreateJobLog(jobID string, message string, level string) error
}
