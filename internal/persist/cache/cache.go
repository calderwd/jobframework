package cache

import (
	"sync"

	"github.com/calderwd/jobframework/api"
	ll "github.com/calderwd/jobframework/internal/logger"
	"github.com/google/uuid"
)

type JobCache struct {
	mu    sync.RWMutex
	cache map[uuid.UUID]api.JobSummary
}

type Persister interface {
}

var (
	logger ll.Logger = ll.GetInstance()
)

func (jc *JobCache) Update(js api.JobSummary) bool {

	if js.State == api.Running {
		logger.Info("Cache " + js.Dump())
		jc.mu.Lock()
		defer jc.mu.Unlock()

		jc.cache[js.Uuid] = js
		return true
	}
	return false
}

func (jc *JobCache) Get(uuid uuid.UUID) (api.JobSummary, bool) {

	jc.mu.RLock()
	defer jc.mu.RUnlock()

	if j, ok := jc.cache[uuid]; !ok {
		return j, true
	}
	return api.JobSummary{}, false
}

func NewJobCache() JobCache {
	return JobCache{
		cache: make(map[uuid.UUID]api.JobSummary),
	}
}
