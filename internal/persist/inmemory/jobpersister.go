package inmemory

import (
	"sync"

	"github.com/calderwd/jobframework/api"
	ll "github.com/calderwd/jobframework/internal/logger"
	"github.com/calderwd/jobframework/internal/persist/cache"
	"github.com/google/uuid"
)

type InMemoryPersister struct {
}

var (
	logger               ll.Logger = ll.GetInstance()
	jobOnce              sync.Once
	jobPersisterInstance InMemoryJobPersister
)

type InMemoryJobPersister struct {
	jobs  map[uuid.UUID]api.JobSummary
	cache cache.JobCache
}

func (p *InMemoryJobPersister) AddJob(js api.JobSummary, user string) error {
	p.jobs[js.Uuid] = js
	return nil
}

func (p *InMemoryJobPersister) GetJob(uuid uuid.UUID, user string) (api.JobSummary, error) {
	if js, found := p.cache.Get(uuid); found {
		return js, nil
	}
	return p.jobs[uuid], nil
}

func (p *InMemoryJobPersister) UpdateJob(js api.JobSummary, user string) error {
	if !p.cache.Update(js) {
		logger.Info("Persister " + js.Dump())
		p.jobs[js.Uuid] = js
	}
	return nil
}

func (p *InMemoryJobPersister) ListJobs(filter api.JobFilter, user string) []api.JobSummary {
	result := make([]api.JobSummary, len(p.jobs))
	for _, v := range p.jobs {
		if js, found := p.cache.Get(v.Uuid); found {
			result = append(result, js)
		} else {
			result = append(result, v)
		}
	}
	return result
}

func (p *InMemoryJobPersister) Shutdown(force bool) {
	// Nothing to do
}

func GetJobPersister() *InMemoryJobPersister {
	jobOnce.Do(func() {
		jobPersisterInstance = InMemoryJobPersister{
			jobs:  make(map[uuid.UUID]api.JobSummary),
			cache: cache.NewJobCache(),
		}
	})
	return &jobPersisterInstance
}
