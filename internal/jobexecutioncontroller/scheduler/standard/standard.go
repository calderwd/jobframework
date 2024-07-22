package standard

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/calderwd/jobframework/api"
	ll "github.com/calderwd/jobframework/internal/logger"
)

const (
	maxConcurrentJobs int           = 10
	maxScheduledJobs  int           = 100
	inactivitySeconds time.Duration = (time.Second * 10)
)

type jobEntry struct {
	jobSummary api.JobSummary
	job        api.IJob
}

type StandardScheduler struct {
	running              atomic.Bool
	dispatcherStream     chan<- jobEntry
	dispatcherTermStream <-chan struct{}
	jobCancelStream      chan<- string
}

type jobWorker struct {
	name              string
	terminationStream <-chan struct{}
	jobCancelStream   <-chan string
}

var (
	logger ll.Logger = ll.GetInstance()
)

func (jw *jobWorker) Start(jobStream <-chan jobEntry) <-chan struct{} {
	terminationStream := make(chan struct{})
	go func() {
		defer close(terminationStream)
		for {
			select {
			case jqe, ok := <-jobStream:
				if ok {
					logger.Info(fmt.Sprintf("Running job in %s\n", jw.name))
					jqe.job.Process(jqe.jobSummary, jw.jobCancelStream)
				} else {
					logger.Info(fmt.Sprintf("closing worker %s\n", jw.name))
					return
				}
			case <-time.After(inactivitySeconds):
				logger.Info("Been waiting too long for a new job")
				return
			}
		}
	}()
	return terminationStream
}

func (sch *StandardScheduler) Start() {
	if sch.running.CompareAndSwap(false, true) {
		sch.dispatcherStream, sch.dispatcherTermStream, sch.jobCancelStream = sch.runDispatcher()
	}
}

func (sch *StandardScheduler) runDispatcher() (chan<- jobEntry, <-chan struct{}, chan<- string) {
	dispatcherStream := make(chan jobEntry)
	dispatcherTermStream := make(chan struct{})
	jobCancelStream := make(chan string)

	go func() {
		logger.Info("Dispatcher Started")
		jobStream := make(chan jobEntry)
		defer close(jobStream)
		defer close(dispatcherTermStream)

		workerPool := make([]jobWorker, 0, maxConcurrentJobs)

		for je := range dispatcherStream {

			if len(workerPool) <= maxConcurrentJobs {
				jobWorker := jobWorker{
					name: fmt.Sprintf("Worker-%d", len(workerPool)+1),
				}
				jobCancelStream := make(chan string)
				jobWorker.terminationStream = jobWorker.Start(jobStream)
				jobWorker.jobCancelStream = jobCancelStream
				workerPool = append(workerPool, jobWorker)
			}

			jobStream <- je
		}
		logger.Info("closing dispatcher - closing workers")
		for _, worker := range workerPool {
			<-worker.terminationStream
			logger.Info(fmt.Sprintf("Received termination from %s\n", worker.name))
		}

	}()

	return dispatcherStream, dispatcherTermStream, jobCancelStream
}

func (sch *StandardScheduler) ScheduleJob(js api.JobSummary, job api.IJob) error {
	sch.dispatcherStream <- jobEntry{
		jobSummary: js,
		job:        job,
	}
	return nil
}

func (sch *StandardScheduler) CancelJob(js api.JobSummary) bool {
	sch.jobCancelStream <- js.Uuid.String()
	return false
}

func (sch *StandardScheduler) Stop() {
	close(sch.dispatcherStream)
	logger.Info("Waiting for dispatcher to terminate")
	<-sch.dispatcherTermStream
	logger.Info("Dispatcher terminated")
}
