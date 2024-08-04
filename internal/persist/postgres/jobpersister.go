package postgres

import (
	"errors"
	"sync"

	"github.com/calderwd/jobframework/api"
	ll "github.com/calderwd/jobframework/internal/logger"
	"github.com/google/uuid"
)

type InMemoryPersister struct {
}

const (
	host     = ""
	port     = 9432
	user     = ""
	password = ""
	dbname   = ""
)

var (
	logger               ll.Logger = ll.GetInstance()
	jobOnce              sync.Once
	jobPersisterInstance *PostgresJobPersister
)

type PostgresJobPersister struct {
	db *PGPersister
}

func (p *PostgresJobPersister) AddJob(js api.JobSummary, user string) error {
	res, err := p.db.db.Exec(`INSERT INTO jobs (uuid, name, description, state, status ) VALUES ($1, $2, $3, $4, $5) 
							ON CONFLICT (uuid) DO NOTHING`,
		js.Uuid, js.Name, js.State, js.Status, js.Description)

	if err != nil {
		return err
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("job with this UUID already exists")
	}

	return err
}

func (p *PostgresJobPersister) GetJob(uuid uuid.UUID, user string) (api.JobSummary, error) {
	var job api.JobSummary
	err := p.db.db.QueryRow(`SELECT uuid, name, description, state, status FROM jobs WHERE uuid = $1`, uuid).Scan(
		&job.Uuid, &job.Name, &job.Description, &job.State, &job.Status)
	return job, err
}

func (p *PostgresJobPersister) UpdateJob(js api.JobSummary, user string) error {
	res, err := p.db.db.Exec(`UPDATE jobs SET name = $2, description = $3, state = $4, status = $5 WHERE uuid = $1`,
		js.Uuid, js.Name, js.Description, js.State, js.Status)
	if err != nil {
		return err
	}
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("no rows updated, job not found")
	}
	return nil
}

func (p *PostgresJobPersister) ListJobs(filter api.JobFilter, user string) []api.JobSummary {
	rows, err := p.db.db.Query(`SELECT uuid, name, description, state, status FROM jobs`)
	if err != nil {
		logger.Error(err, "failed to retrieve jobs from postgres")
		return nil
	}
	defer rows.Close()

	var jobs []api.JobSummary
	for rows.Next() {
		var job api.JobSummary
		if err := rows.Scan(&job.Uuid, &job.Name, &job.Description, &job.State, &job.Status); err != nil {
			logger.Error(err, "failed to retrieve job row from postgres")
			return nil
		}
		jobs = append(jobs, job)
	}
	return jobs
}

func GetJobPersister() *PostgresJobPersister {
	jobOnce.Do(func() {
		jobPersisterInstance = &PostgresJobPersister{
			db: GetPGPersister(),
		}
	})

	return jobPersisterInstance
}

func (p *PostgresJobPersister) Shutdown(force bool) {
	p.db.db.Close()
}
