package jobexecutioncontroller

import (
	"github.com/calderwd/jobframework/api"
	sch "github.com/calderwd/jobframework/internal/jobexecutioncontroller/scheduler"
	ll "github.com/calderwd/jobframework/internal/logger"
	"github.com/calderwd/jobframework/internal/persist"
)

var logger ll.Logger = ll.GetInstance()

type JobExecutionController struct {
}

func (jec JobExecutionController) ScheduleJob(js api.JobSummary, user string) error {

	jc, err := GetJobRegistrar().GetJobConfig(js.JobType)

	if err == nil {

		scheduler, err := sch.GetJobScheduler(jc.Scheduler)

		if err != nil {
			return err
		}

		scheduler.ScheduleJob(js, jc.Job)

		if js.LastExecutionStart != nil {
			js.State = api.Rescheduled
		} else {
			js.State = api.Scheduled
		}

		persist.GetJobPersister().UpdateJob(js, user)
	} else {

		logger.Error(err)
	}

	return err
}

func (jec JobExecutionController) Shutdown(force bool) {
	sch.Shutdown(force)
}
