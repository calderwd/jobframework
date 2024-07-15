package postgres

import "github.com/calderwd/jobframework/api"

type PostgresJobPersister struct {
}

func (p PostgresJobPersister) AddJob(js api.JobSummary, user string) {

}

type PostgresLogPersister struct {
}

type PostgresResultPersister struct {
}

func GetJobPersister() PostgresJobPersister {
	return PostgresJobPersister{}
}

func GetLogPersister() PostgresLogPersister {
	return PostgresLogPersister{}
}

func GetResultPersister() PostgresResultPersister {
	return PostgresResultPersister{}
}
