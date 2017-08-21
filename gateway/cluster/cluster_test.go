package cluster

import (
	"fmt"
	"sync"
	"testing"
)

func TestIncrClusterId(t *testing.T) {
	var wg sync.WaitGroup
	c := NewCluster();
	iterations := 1000

	locker := newLocker()
	results := make([]bool, iterations)

	for i := 0; i < iterations; i++ {
		wg.Add(1)
		go func() {
			id := c.IncrClusterId()
			locker.run(func() {
				results[id] = true
			})
			defer wg.Done()
		}()
	}
	wg.Wait()
	for i, b := range results {
		if !b {
			t.Error(fmt.Sprintf("failed to return %d", i))
		}
	}
}
