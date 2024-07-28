package standard

import (
	"context"
	"errors"
	"sync/atomic"

	"github.com/calderwd/jobframework/api"
	ll "github.com/calderwd/jobframework/internal/logger"
	"github.com/google/uuid"
)

var (
	logger ll.Logger = ll.GetInstance()
)

const (
	maxConcurrentJobs int = 10
	maxScheduledJobs  int = 100
)

type jobEntry struct {
	jobSummary api.JobSummary
	job        api.IJob
}

type StandardScheduler struct {
	running              atomic.Bool
	ctx                  context.Context
	cancel               context.CancelFunc
	dispatcherStream     chan<- jobEntry  // Jobs to be scheduled through dispatcher
	dispatcherTermStream <-chan struct{}  // Indicates termination of dispatcher
	jobCancelStream      chan<- uuid.UUID // Fan-out entry channel for job cancellation request
}

// Start the dispatcher thread if not already active
func (sch *StandardScheduler) Start() {
	if sch.running.CompareAndSwap(false, true) {
		sch.ctx, sch.cancel = context.WithCancel(context.Background())
		sch.dispatcherStream, sch.dispatcherTermStream, sch.jobCancelStream = sch.runDispatcher(sch.ctx)
	}
}

// Create dispatcher thread, which will manage the worker pool, sending jobs for execution to pool
// Termination of dispatcher will only complete once all workers report inactivity, i.e. jobs complete
func (sch *StandardScheduler) runDispatcher(ctx context.Context) (chan<- jobEntry, <-chan struct{}, chan<- uuid.UUID) {
	dispatcherStream := make(chan jobEntry, maxScheduledJobs)
	dispatcherTermStream := make(chan struct{})
	jobCancelStream := make(chan uuid.UUID)
	workerTermStream := make(chan workerTerm)

	go func() {
		logger.Info("Dispatcher Started")
		jobStream := make(chan jobEntry)
		defer close(jobStream)
		defer close(dispatcherTermStream)
		defer close(jobCancelStream)
		defer close(workerTermStream)

		workerPool := newWorkerPool(maxConcurrentJobs, workerTermStream)

		done := ctx.Done()

	outer:
		for {
			select {
			case je, ok := <-dispatcherStream:

				// If dispatcher channel closed then exit dispatcher (only reachable if channel empty)
				if !ok {
					break outer
				}

				// If worker(s) busy add worker, and post job
				select {
				case jobStream <- je:
				default:
					if workerPool.length() <= maxConcurrentJobs {
						workerPool.addWorker(jobStream)
					}
					jobStream <- je
				}

			case <-done:
				// Done received, so drain dispatcher stream and wait for worker pool to complete
				for range dispatcherStream {
				}

				done = nil
				dispatcherStream = nil

				// Wait for worker pool to empty
				for workerPool.length() > 0 {
					wt := <-workerPool.terminationStream
					workerPool.removeWorker(wt.name)
				}
				break outer

			case uuid := <-jobCancelStream:
				// Broadcast cancellation request to works for designated job uuid
				workerPool.sendCancelRequest(uuid)

			case wt := <-workerPool.terminationStream:
				// Remove worker from pool on worker termination indication
				workerPool.removeWorker(wt.name)
			}
		}
	}()

	return dispatcherStream, dispatcherTermStream, jobCancelStream
}

// Post a new job to the dispatcher for execution
func (sch *StandardScheduler) ScheduleJob(js api.JobSummary, job api.IJob) error {
	select {
	case sch.dispatcherStream <- jobEntry{
		jobSummary: js,
		job:        job,
	}:
	default:
		logger.Info("Maximum number of jobs reached - job discarded")
		return errors.New("maximum number of jobs reached - job discarded")
	}
	return nil
}

// Post a new job cancellation request to the dispatcher
func (sch *StandardScheduler) CancelJob(js api.JobSummary) bool {
	sch.jobCancelStream <- js.Uuid
	return false
}

// Shutdown the scheduler by closing the dispatcher channel and waiting for indication of shutdown
func (sch *StandardScheduler) Stop() {
	close(sch.dispatcherStream)
	sch.cancel()
	logger.Info("Waiting for dispatcher to terminate")
	<-sch.dispatcherTermStream
	logger.Info("Dispatcher terminated")
}
