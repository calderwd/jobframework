package standard

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type jobWorker struct {
	name            string
	jobCancelStream chan string
}

type workerTerm struct {
	name string
}

type workerPool struct {
	maxWorkers        int
	workerPool        []jobWorker
	terminationStream chan workerTerm
}

const (
	inactivitySeconds int64 = (int64(time.Second) * 10)
)

// Define a new worker pool
//
//	maxWorkers: The maximum numer of workers in pool
//	terminationStream: Worker termination notification channel
func newWorkerPool(maxWorkers int, terminationStream chan workerTerm) workerPool {
	return workerPool{
		maxWorkers:        maxWorkers,
		workerPool:        make([]jobWorker, 0, maxWorkers),
		terminationStream: terminationStream,
	}
}

// Adds a new worker to the worker pool, pulling jobs from jobStream
func (wp *workerPool) addWorker(jobStream chan jobEntry) {
	jobWorker := jobWorker{
		name: fmt.Sprintf("Worker-%d", len(wp.workerPool)+1),
	}
	jobWorker.jobCancelStream = jobWorker.start(jobStream, wp.terminationStream)
	wp.workerPool = append(wp.workerPool, jobWorker)
}

// Removes a worker from the worker pool
func (wp *workerPool) removeWorker(name string) {
	for i := 0; i < wp.length(); i++ {
		wk := wp.workerPool[i]
		if wk.name == name {
			if i == 0 {
				wp.workerPool = wp.workerPool[1:]
				logger.Info(fmt.Sprintf("removing worker %s", name))
				return
			}
			wp.workerPool = append(wp.workerPool[0:i], wp.workerPool[i+1:]...)
			logger.Info(fmt.Sprintf("removing worker %s", name))
			return
		}
	}
}

// Starts a worker, receiving from jobStream and checking for inactivity
// Returns a job cancellation request channel for this worker
func (jw *jobWorker) start(jobStream <-chan jobEntry, terminationStream chan<- workerTerm) chan string {
	jobCancelStream := make(chan string)

	go func() {
		defer close(jobCancelStream)
		inactivityStream := time.NewTimer(time.Duration(inactivitySeconds))

		for {
			if !inactivityStream.Stop() {
				<-inactivityStream.C
			}
			inactivityStream.Reset(time.Duration(inactivitySeconds))
			select {
			case jqe, ok := <-jobStream:
				if ok {
					logger.Info(fmt.Sprintf("Running job in %s", jw.name))
					jqe.job.Process(jqe.jobSummary, jw.jobCancelStream)
				} else {
					logger.Info(fmt.Sprintf("closing worker %s", jw.name))
					terminationStream <- workerTerm{name: jw.name}
					return
				}
			case <-inactivityStream.C:
				logger.Info(fmt.Sprintf("Been waiting too long for a new job [%s]", jw.name))
				terminationStream <- workerTerm{name: jw.name}
				return
			}
		}
	}()

	return jobCancelStream
}

// Broadcast a job cancellation request to all workers for the specified job (uuid)
func (wp *workerPool) sendCancelRequest(uuid uuid.UUID) {
	for _, w := range wp.workerPool {
		w.jobCancelStream <- uuid.String()
	}
}

// Return the number of workers in pool
func (wp workerPool) length() int {
	return len(wp.workerPool)
}
