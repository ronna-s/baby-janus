package cluster

import (
	"fmt"
	"sync"
	"testing"
)

func TestLock(t *testing.T) {

	var wg sync.WaitGroup
	i := 0
	locker := newLocker();
	f := func() {
		i++
		defer wg.Done()
	}

	for i := 0; i < 100; i++ {
		wg.Add(1)
		go locker.run(f)
	}

	wg.Wait()

	if i != 100 {
		t.Error(fmt.Sprintf("race condition in locker execution %d", i))
	}
}
