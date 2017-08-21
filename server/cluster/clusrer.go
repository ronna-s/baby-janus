package cluster

import (
	"fmt"
	"math/rand"
)

type (
	Cluster interface {
		GetParts() []string
		GetInstanceParts(int) []string
	}

	cluster struct {
		numInstances int
		numParts     int
		slicer       func(int) interface{}
		randomize    func([]string) []string
	}
)

const (
	NUM_PARTS = 136
	NUM_INSTANCES = 10
)

func NewCluster() Cluster {
	return &cluster{
		numInstances: NUM_INSTANCES,
		numParts: NUM_PARTS,
		slicer:   getPart,
		randomize:    randomize,
	}
}

func (c *cluster) GetParts() []string {
	res := make([]string, c.numParts)
	for i := range res {
		res[i] = fmt.Sprintf("%v", c.slicer(i))
	}
	return res
}

func (c *cluster) GetInstanceParts(instanceId int) []string {
	parts := c.randomize(c.GetParts())[0:c.numParts]
	lhs := instanceId * (c.numParts + 1) / c.numInstances
	rhs := ((instanceId + 1) * (c.numParts + 1) / c.numInstances)
	if lhs > len(parts) {
		return []string{}
	}
	if rhs > len(parts) {
		return parts[lhs: ]
	}
	return parts[lhs: rhs]
}

func randomize(slice []string) []string {
	if len(slice) == 1 {
		return slice
	}
	for i := range slice {
		j := rand.Intn(len(slice) - 1)
		slice[i], slice[j] = slice[j], slice[i]
	}
	return slice
}

func getPart(pos int) interface{} {
	return fmt.Sprintf("/parts/%d.part", pos)
}
