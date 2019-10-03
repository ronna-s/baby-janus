package gateway

import (
	"sync"
)

type (
	lock struct {
		mutex sync.Mutex
	}
)

func newLocker() *lock {
	return &lock{mutex: sync.Mutex{}}
}

func (l *lock) run(f func()) {
	l.mutex.Lock()
	defer l.mutex.Unlock()
	f()
}
