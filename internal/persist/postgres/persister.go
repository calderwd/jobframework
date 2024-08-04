package postgres

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
)

var pgPersisterInstance *PGPersister

type PGPersister struct {
	db *sql.DB
}

func GetPGPersister() *PGPersister {
	jobOnce.Do(func() {
		psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)

		db, err := sql.Open("postgres", psqlInfo)
		if err != nil {
			panic(err)
		}

		err = db.Ping()
		if err != nil {
			panic(err)
		}

		pgPersisterInstance = &PGPersister{
			db: db,
		}
	})

	return pgPersisterInstance
}
