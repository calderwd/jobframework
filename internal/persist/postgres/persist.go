package postgres

import (
	"errors"
	"sync"

	"github.com/calderwd/jobframework/api"
)

var jobOnce sync.Once
var jobPersisterInstance PostgresJobPersister

var logOnce sync.Once
var logPersisterInstance PostgresLogPersister

var resultOnce sync.Once
var resultPersisterInstance PostgresResultPersister

type PostgresJobPersister struct {
}

func (p PostgresJobPersister) AddJob(js api.JobSummary, user string) error {
	return errors.New("not yet supported")
}

func (p PostgresJobPersister) UpdateJob(js api.JobSummary, user string) error {
	return errors.New("not yet supported")
}

type PostgresLogPersister struct {
}

type PostgresResultPersister struct {
}

func GetJobPersister() PostgresJobPersister {
	jobOnce.Do(func() {
		jobPersisterInstance = PostgresJobPersister{}
	})
	return jobPersisterInstance
}

func GetLogPersister() PostgresLogPersister {
	logOnce.Do(func() {
		logPersisterInstance = PostgresLogPersister{}
	})
	return logPersisterInstance
}

func GetResultPersister() PostgresResultPersister {
	resultOnce.Do(func() {
		resultPersisterInstance = PostgresResultPersister{}
	})
	return resultPersisterInstance
}
