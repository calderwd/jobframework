package postgres

import "sync"

var logOnce sync.Once
var logPersisterInstance PostgresLogPersister

type PostgresLogPersister struct {
	db *PGPersister
}

func GetLogPersister() *PostgresLogPersister {
	logOnce.Do(func() {
		logPersisterInstance = PostgresLogPersister{
			db: GetPGPersister(),
		}
	})
	return &logPersisterInstance
}

func (p *PostgresLogPersister) Shutdown(force bool) {

}
