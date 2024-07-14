package jobmanager

import (
	"sync"

	"github.com/calderwd/jobframework/api"
	"github.com/google/uuid"
)

var instance *jobManager

var once sync.Once

func NewInstance() api.JobService {

	once.Do(func() {
		instance = &jobManager{}
	})

	return instance
}

type jobManager struct {
}

func (jm *jobManager) GetJobRegistrar() api.JobRegistrar {
	return nil
}

func (jm *jobManager) AddJob(jobType string, jobSchedule api.JobSchedule, jobContext api.JobContext, user string) (uuid.UUID, error) {
	return uuid.New(), nil
}

func (jm *jobManager) GetJob(uuid uuid.UUID, user string) (api.JobSummary, error) {
	return api.JobSummary{}, nil
}

func (jm *jobManager) CancelJob(uuid uuid.UUID, user string) error {
	return nil
}

func (jm *jobManager) DeleteJob(uuid uuid.UUID, user string) error {
	return nil
}

func (jm *jobManager) ListJobs(filter api.JobFilter, user string) []api.JobSummary {
	return []api.JobSummary{api.JobSummary{}}
}

func (jm *jobManager) GetEvaluation(uuid uuid.UUID, evaluationId uint64, user string) (string, error) {
	return "", nil
}

func (jm *jobManager) GetJobLogPageCount(uuid uuid.UUID, evaluationId uint64, user string) (uint, error) {
	return 0, nil
}

func (jm *jobManager) GetJobLogPage(uuid uuid.UUID, evaluationId uint64, pageNumber uint, user string) ([]string, error) {
	return []string{}, nil
}

func (jm *jobManager) GetJobHistory(uuid uuid.UUID, user string) ([]api.JobSummary, error) {
	return []api.JobSummary{}, nil
}
