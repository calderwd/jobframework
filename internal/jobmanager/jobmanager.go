package jobmanager

import (
	"sync"

	"github.com/calderwd/jobframework/api"
	jec "github.com/calderwd/jobframework/internal/jobexecutioncontroller"
	"github.com/calderwd/jobframework/internal/persist"
	"github.com/google/uuid"
)

var instance *jobManager

var once sync.Once

func GetInstance() api.JobService {

	once.Do(func() {
		instance = &jobManager{
			jobExec: jec.JobExecutionController{},
		}
	})

	return instance
}

type jobManager struct {
	jobExec jec.JobExecutionController
}

func (jm *jobManager) GetJobRegistrar() api.JobRegistrar {
	return jec.GetJobRegistrar()
}

func (jm *jobManager) AddJob(jobType string, jobSchedule api.JobSchedule, jobContext api.JobContext, user string) (uuid.UUID, error) {

	js := api.JobSummary{}.Build(jobType, jobSchedule, jobContext)

	jc, err := jec.GetJobRegistrar().GetJobConfig(jobType)

	if err != nil {
		return uuid.Nil, err
	}

	if jc.Profile.CanAdd() {
		persist.GetJobPersister().AddJob(js, user)
		jm.jobExec.ScheduleJob(js)

	}
	return js.Uuid, nil
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
	return []api.JobSummary{}
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