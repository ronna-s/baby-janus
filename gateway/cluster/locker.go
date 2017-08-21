package cluster

import (
	"sync"
)

type (
	mutexLocker struct {
		mutex sync.Mutex
	}
)

func newLocker() *mutexLocker {
	return &mutexLocker{mutex: sync.Mutex{}}
}

func (l *mutexLocker) run(f func()) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	f()
}
