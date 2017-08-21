package cluster

import (
	"time"
)

type (
	Cluster interface {
		IncrClusterId() int
		GetSeed() int64
	}

	cluster struct {
		locker         *mutexLocker
		currInstanceId int
		seed           int64
	}
)

func NewCluster() *cluster {
	return &cluster{
		locker:         newLocker(),
		currInstanceId: -1,
		seed:           time.Now().UnixNano(),
	}
}

func (c *cluster) IncrClusterId() int {
	var res int;
	c.locker.run(func() {
		c.currInstanceId++
		res = c.currInstanceId
	})
	return res
}

func (c *cluster) GetSeed() int64 {
	return c.seed
}
