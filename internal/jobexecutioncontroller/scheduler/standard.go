package scheduler

import (
	"context"
	"errors"

	"github.com/calderwd/jobframework/api"
)

type jobQueueEntry struct {
	jobSummary api.JobSummary
	job        api.IJob
}

type jobUpdateEntry struct {
	jobCancel  context.CancelFunc
	jobSummary api.JobSummary
}

type StandardScheduler struct {
	running         bool
	poolSize        uint
	jobCount        uint
	jobQueue        chan jobQueueEntry
	schedulerTerm   chan struct{}
	jobUpdateStream chan jobUpdateEntry
	schedulerCtx    context.Context
	stop            context.CancelFunc
}

func (sch *StandardScheduler) updateJobSummary(ju jobUpdateEntry) {
	logger.Info("update for job received")
}

func (sch *StandardScheduler) Start() {

	if sch.running {
		return
	}
	sch.running = true
	sch.poolSize = 10
	sch.jobQueue = make(chan jobQueueEntry)
	sch.schedulerTerm = make(chan struct{})
	sch.jobUpdateStream = make(chan jobUpdateEntry)
	sch.schedulerCtx, sch.stop = context.WithCancel(context.Background())

	go func() error {
		defer close(sch.jobUpdateStream)

		for {
			select {
			case <-sch.schedulerCtx.Done():
				return sch.schedulerCtx.Err()
			case ju := <-sch.jobUpdateStream:
				sch.updateJobSummary(ju)
			}
		}
	}()

	go func() error {
		defer close(sch.jobQueue)
		defer close(sch.schedulerTerm)

		for {
			select {
			case <-sch.schedulerCtx.Done():
				return sch.schedulerCtx.Err()
			case jqe := <-sch.jobQueue:
				ctx, jobCancel := context.WithCancel(sch.schedulerCtx)
				go jqe.job.Process(jqe.jobSummary, ctx)
				sch.jobUpdateStream <- jobUpdateEntry{
					jobCancel:  jobCancel,
					jobSummary: jqe.jobSummary,
				}
			}
		}
	}()
	//<-sch.schedulerTerm
}

func (sch *StandardScheduler) ScheduleJob(js api.JobSummary, job api.IJob) error {
	sch.jobCount++
	if sch.jobCount >= sch.poolSize {
		return errors.New("schedule thread pool exceeded")
	}

	sch.jobQueue <- jobQueueEntry{
		jobSummary: js,
		job:        job,
	}
	return nil
}

func (sch *StandardScheduler) CancelJob(js api.JobSummary) bool {
	return true
}

func (sch *StandardScheduler) Stop() {
	sch.stop()
}
