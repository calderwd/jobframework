package scheduler

import (
	"fmt"

	"github.com/calderwd/jobframework/api"
	ll "github.com/calderwd/jobframework/internal/logger"
)

var logger ll.Logger = ll.GetInstance()

type Scheduler interface {
	ScheduleJob(js api.JobSummary, job api.IJob)
}

type schedulerFactory struct {
	schedulers map[api.ScheduleType]Scheduler
}

var sf = schedulerFactory{
	schedulers: map[api.ScheduleType]Scheduler{
		api.Standard: &StandardScheduler{},
		// api.Sequential: SequentialScheduler{},
		// api.Periodic: PeriodicScheduler{},
		// api.Cron: CronScheduler{}.
	},
}

func GetJobScheduler(s api.ScheduleType) (Scheduler, error) {
	switch s {
	case api.Standard:
	case api.Sequential:
	case api.Periodic:
	case api.Cron:
	default:
		err := fmt.Errorf("unknown scheduler type %d ", s)
		logger.Error(err)
		return nil, err
	}

	sc := sf.schedulers[s]
	return sc, nil
}
