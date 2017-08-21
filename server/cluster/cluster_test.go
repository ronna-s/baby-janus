package cluster

import (
	"fmt"
	"strconv"
	"testing"
)

func TestGetParts(t *testing.T) {
	c := NewCluster()
	c.numParts = 0
	if len(c.GetParts()) != 0 {
		t.Error("should return empty array")
	}
	c.numParts = 10
	parts := c.GetParts()
	if len(parts) != 10 {
		t.Error("should return 10 items")
	}
	if parts[5] != fmt.Sprintf("/parts/5.part") {
		t.Error("returned incorrect path")
	}
}


func TestGetInstanceParts(t *testing.T) {
	c := NewCluster()
	c.randomize = func(s []string) []string { return s }
	c.slicer = func(pos int) interface{} { return pos }

	c.numParts = 0
	c.numInstances = 10
	for i := 0; i < c.numInstances; i++ {
		equals(t, len(c.GetInstanceParts(i)), 0)
	}

	c.numInstances = 10
	c.numParts = NUM_PARTS
	resStr := []string{}

	for i := 0; i < c.numInstances; i++ {
		resStr = append(resStr, c.GetInstanceParts(i)...)
	}
	equals(t, len(resStr), c.numParts)
	for i := 0; i < c.numParts; i++ {
		resInt, _ := strconv.Atoi(resStr[i])
		equals(t, resInt, i)
	}

}

func equals(t *testing.T, a interface{}, b interface{}) {
	if a != b {
		t.Fatal(fmt.Sprintf("expected %v to equal %v", a, b))
	}
}
