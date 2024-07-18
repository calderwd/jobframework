package scheduler

import "github.com/calderwd/jobframework/api"

type StandardScheduler struct {
}

func (sch *StandardScheduler) ScheduleJob(js api.JobSummary, job api.IJob) {
	go job.Process()
}
