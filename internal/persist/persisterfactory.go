package persist

import (
	"github.com/calderwd/jobframework/api"
	"github.com/calderwd/jobframework/internal/persist/inmemory"
	"github.com/calderwd/jobframework/internal/persist/postgres"
)

const (
	InMemory = iota
	Postgres
)

type JobPersister interface {
	AddJob(js api.JobSummary, user string)
}

type LogPersister interface {
}

type ResultPersister interface {
}

var persisterType = InMemory

func GetJobPersister() JobPersister {

	if persisterType == InMemory {
		return inmemory.GetJobPersister()
	}
	return postgres.GetJobPersister()
}

func GetLogPersister() LogPersister {
	if persisterType == InMemory {
		return inmemory.GetLogPersister()
	}
	return postgres.GetLogPersister()
}

func GetResultPersister() ResultPersister {
	if persisterType == InMemory {
		return inmemory.GetResultPersister()
	}
	return postgres.GetResultPersister()

}
