package inmemory

import (
	"sync"

	"github.com/calderwd/jobframework/api"
	"github.com/google/uuid"
)

var jobOnce sync.Once
var jobPersisterInstance InMemoryJobPersister

var logOnce sync.Once
var logPersisterInstance InMemoryLogPersister

var resultOnce sync.Once
var resultPersisterInstance InMemoryResultPersister

type InMemoryJobPersister struct {
	jobs map[uuid.UUID]api.JobSummary
}

func (p *InMemoryJobPersister) AddJob(js api.JobSummary, user string) error {
	p.jobs[js.Uuid] = js
	return nil
}

func (p *InMemoryJobPersister) UpdateJob(js api.JobSummary, user string) error {
	p.jobs[js.Uuid] = js
	return nil
}

type InMemoryLogPersister struct {
}

type InMemoryResultPersister struct {
}

func GetJobPersister() *InMemoryJobPersister {
	jobOnce.Do(func() {
		jobPersisterInstance = InMemoryJobPersister{
			jobs: make(map[uuid.UUID]api.JobSummary),
		}
	})
	return &jobPersisterInstance
}

func GetLogPersister() *InMemoryLogPersister {
	logOnce.Do(func() {
		logPersisterInstance = InMemoryLogPersister{}
	})
	return &logPersisterInstance
}

func GetResultPersister() *InMemoryResultPersister {
	resultOnce.Do(func() {
		resultPersisterInstance = InMemoryResultPersister{}
	})
	return &resultPersisterInstance
}
