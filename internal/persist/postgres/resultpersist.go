package postgres

import (
	"sync"
)

var resultOnce sync.Once
var resultPersisterInstance PostgresResultPersister

type PostgresResultPersister struct {
	db *PGPersister
}

func GetResultPersister() PostgresResultPersister {
	resultOnce.Do(func() {
		resultPersisterInstance = PostgresResultPersister{
			db: GetPGPersister(),
		}
	})
	return resultPersisterInstance
}
