package jobexecutioncontroller

// type JobRegistrar interface {
// 	RegisterJobType(jobType string, job IJob, profile JobProfile, scheduleType ScheduleType)
// }

import (
	"errors"
	"sync"

	"github.com/calderwd/jobframework/api"
)

var (
	onceRegistry sync.Once
	instance     *jobRegistrar

	errUnknownJob error = errors.New("unknown job type specified")
)

type jobRegistrar struct {
	rj map[string]api.JobConfig
}

func GetJobRegistrar() api.JobRegistrar {
	onceRegistry.Do(func() {
		instance = &jobRegistrar{
			rj: make(map[string]api.JobConfig),
		}
	})

	return instance
}

func (jr *jobRegistrar) RegisterJobType(jobType string, job api.IJob, profile api.JobProfile, scheduleType api.ScheduleType) {
	jr.rj[jobType] = api.JobConfig{
		JobType:   jobType,
		Job:       job,
		Scheduler: scheduleType,
		Profile:   profile,
	}
}

func (jr *jobRegistrar) GetJobConfig(jobType string) (api.JobConfig, error) {

	if jc, ok := jr.rj[jobType]; ok {
		return jc, nil
	}

	return api.JobConfig{}, errUnknownJob
}
