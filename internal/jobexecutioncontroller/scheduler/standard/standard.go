package standard

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/calderwd/jobframework/api"
	ll "github.com/calderwd/jobframework/internal/logger"
	"github.com/google/uuid"
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
	jobCancelStream      chan<- uuid.UUID
}

type jobWorker struct {
	name              string
	terminationStream <-chan struct{}
	jobCancelStream   chan string
}

type workerPool struct {
	maxWorkers int
	workerPool []jobWorker
}

var (
	logger ll.Logger = ll.GetInstance()
)

func newWorkerPool(maxWorkers int) workerPool {
	return workerPool{
		maxWorkers: maxWorkers,
		workerPool: make([]jobWorker, 0, maxWorkers),
	}
}

func (wp *workerPool) addWorker(jobStream chan jobEntry) {
	jobWorker := jobWorker{
		name: fmt.Sprintf("Worker-%d", len(wp.workerPool)+1),
	}

	jobCancelStream := make(chan string)
	jobWorker.jobCancelStream = jobCancelStream
	jobWorker.terminationStream = jobWorker.Start(jobStream)
	wp.workerPool = append(wp.workerPool, jobWorker)
}

func (wp *workerPool) removeWorker(name string) {
	for i, wk := range wp.workerPool {
		if wk.name == name {
			if i == 0 {
				wp.workerPool = wp.workerPool[1:]
				logger.Info(fmt.Sprintf("removing worker %s", name))
				return
			}
			wp.workerPool = append(wp.workerPool[0:i], wp.workerPool[i+1:]...)
			logger.Info(fmt.Sprintf("removing worker %s", name))
		}
	}
}

func (wp *workerPool) waitForShutdown() {
	logger.Info("closing dispatcher - closing workers")
	for _, worker := range wp.workerPool {
		_, ok := <-worker.terminationStream
		if !ok {
			defer close(worker.jobCancelStream)
			wp.removeWorker(worker.name)
			logger.Info(fmt.Sprintf("Received termination from %s\n", worker.name))
		}
	}
}

func (wp *workerPool) sweepTerminatedWorkers() {
	for _, worker := range wp.workerPool {
		select {
		case _, ok := <-worker.terminationStream:
			if !ok {
				defer close(worker.jobCancelStream)
				wp.removeWorker(worker.name)
				logger.Info(fmt.Sprintf("sweeper removing worker %s", worker.name))
			}
		default:
		}
	}
}

func (wp *workerPool) sendCancelRequest(uuid uuid.UUID) {
	for _, w := range wp.workerPool {
		w.jobCancelStream <- uuid.String()
	}
}

func (wp workerPool) length() int {
	return len(wp.workerPool)
}

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

func (sch *StandardScheduler) runDispatcher() (chan<- jobEntry, <-chan struct{}, chan<- uuid.UUID) {
	dispatcherStream := make(chan jobEntry)
	dispatcherTermStream := make(chan struct{})
	jobCancelStream := make(chan uuid.UUID)

	go func() {
		logger.Info("Dispatcher Started")
		jobStream := make(chan jobEntry)
		defer close(jobStream)
		defer close(dispatcherTermStream)
		defer close(jobCancelStream)

		workerPool := newWorkerPool(maxConcurrentJobs)

		for je := range dispatcherStream {

			if workerPool.length() <= maxConcurrentJobs {
				workerPool.addWorker(jobStream)
			}

			select {
			case jobStream <- je:
			case uuid := <-jobCancelStream:
				workerPool.sendCancelRequest(uuid)
			}

			workerPool.sweepTerminatedWorkers()
		}
		workerPool.waitForShutdown()

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
	sch.jobCancelStream <- js.Uuid
	return false
}

func (sch *StandardScheduler) Stop() {
	close(sch.dispatcherStream)
	logger.Info("Waiting for dispatcher to terminate")
	<-sch.dispatcherTermStream
	logger.Info("Dispatcher terminated")
}
