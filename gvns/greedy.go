package gvns

import (
	"math"
)

func getMinIndex(array []int) int {
	minIndex := 0
	var min int = array[minIndex]
	for index, value := range array {
		if value < min {
			min = value
			minIndex = index
		}
	}
	return minIndex
}

func Greedy(tsp *TSP) *Solution {
	best := NewSolution(tsp, []int{0})
	for i := 0; i < tsp.numNodes-1; i++ {
		current := best.GetNode(i)
		minIndex := 0
		var min = math.MaxInt64

		for index, value := range tsp.matrix[current] {
			if !best.HasNode(index) && value < min {
				min = value
				minIndex = index
			}
		}

		best.AddNode(minIndex)
	}
	best.AddNode(0)
	return &best
}
