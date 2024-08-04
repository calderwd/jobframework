package inmemory

import "sync"

var logOnce sync.Once
var logPersisterInstance InMemoryLogPersister

type InMemoryLogPersister struct {
}

func GetLogPersister() *InMemoryLogPersister {
	logOnce.Do(func() {
		logPersisterInstance = InMemoryLogPersister{}
	})
	return &logPersisterInstance
}
