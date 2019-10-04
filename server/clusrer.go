package main

import (
	"fmt"
	"math/rand"
)

type (
	Cluster struct {
		numInstances int
		numParts     int
		slicer       func(int) interface{}
		randomize    func([]string) []string
	}
)

const (
	numParts     = 454
	numInstances = 10
)

func NewCluster() *Cluster {
	return &Cluster{
		numInstances: numInstances,
		numParts:     numParts,
		slicer:       getPart,
		randomize:    randomize,
	}
}

func (c *Cluster) GetParts() []string {
	res := make([]string, c.numParts)
	for i := range res {
		res[i] = fmt.Sprintf("%v", c.slicer(i))
	}
	return res
}

func (c *Cluster) GetInstanceParts(instanceId int) []string {
	parts := c.randomize(c.GetParts())[0:c.numParts]
	lhs := instanceId * (c.numParts + 1) / c.numInstances
	rhs := (instanceId + 1) * (c.numParts + 1) / c.numInstances
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
