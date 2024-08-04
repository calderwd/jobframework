package inmemory

import "sync"

var resultOnce sync.Once
var resultPersisterInstance InMemoryResultPersister

type InMemoryResultPersister struct {
}

func GetResultPersister() *InMemoryResultPersister {
	resultOnce.Do(func() {
		resultPersisterInstance = InMemoryResultPersister{}
	})
	return &resultPersisterInstance
}
