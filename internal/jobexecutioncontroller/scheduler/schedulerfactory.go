package scheduler

import (
	"fmt"

	"github.com/calderwd/jobframework/api"
	"github.com/calderwd/jobframework/internal/jobexecutioncontroller/scheduler/standard"
	ll "github.com/calderwd/jobframework/internal/logger"
)

var logger ll.Logger = ll.GetInstance()

type Scheduler interface {
	Start()
	ScheduleJob(js api.JobSummary, job api.IJob) error
	CancelJob(js api.JobSummary) bool
	Stop()
}

type schedulerFactory struct {
	schedulers map[api.ScheduleType]Scheduler
}

var sf = schedulerFactory{
	schedulers: map[api.ScheduleType]Scheduler{
		api.Standard: &(standard.StandardScheduler{}),
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
	sc.Start()

	return sc, nil
}

func Shutdown(force bool) {
	for k := range sf.schedulers {
		sch := sf.schedulers[k]
		sch.Stop()
	}
}
