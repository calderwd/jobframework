package jobexecutioncontroller

import (
	"github.com/calderwd/jobframework/api"
	"github.com/calderwd/jobframework/internal/jobexecutioncontroller/scheduler"
	ll "github.com/calderwd/jobframework/internal/logger"
	"github.com/calderwd/jobframework/internal/persist"
)

var logger ll.Logger = ll.GetInstance()

type JobExecutionController struct {
}

func (jec JobExecutionController) ScheduleJob(js api.JobSummary) error {

	jc, err := GetJobRegistrar().GetJobConfig(js.JobType)

	if err == nil {

		scheduler, err := scheduler.GetJobScheduler(jc.Scheduler)

		if err != nil {
			return err
		}

		scheduler.scheduleJob(js)

		if js.LastExecutionStart != nil {
			js.State = api.Rescheduled
		} else {
			js.State = api.Scheduled
		}

		persist.GetJobPersister().updateJob(js)
	} else {

		logger.Error(err)
	}

	return err
}
