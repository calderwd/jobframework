package inmemory

import "github.com/calderwd/jobframework/api"

type InMemoryJobPersister struct {
}

func (p InMemoryJobPersister) AddJob(js api.JobSummary, user string) {

}

type InMemoryLogPersister struct {
}

type InMemoryResultPersister struct {
}

func GetJobPersister() InMemoryJobPersister {
	return InMemoryJobPersister{}
}

func GetLogPersister() InMemoryLogPersister {
	return InMemoryLogPersister{}
}

func GetResultPersister() InMemoryResultPersister {
	return InMemoryResultPersister{}
}
