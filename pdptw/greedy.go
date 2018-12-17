package pdptw

import (
	"log"
	"math"
	"math/rand"
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

func Greedy(tsp *PDPTW) *Solution {
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

func GetRandom(tsp *PDPTW) *Solution {
	//log.Printf("Start node: %v", tsp.startNode)
	route := []int{tsp.startNode}
	for i := 0; i < tsp.numNodes; i++ {
		if i != tsp.startNode {
			route = append(route, i)
		}
	}

	for i := 1; i < len(route); i++ {
		j := rand.Intn(i) + 1
		route[i], route[j] = route[j], route[i]
	}

	// log.Printf("%v", route)

	s := NewSolution(tsp, route)

	return &s
}

func GetRandomPD(tsp *PDPTW) *Solution {
	r1 := []int{0}
	r2 := []int{}

	tmp := make(map[int]bool, tsp.numNodes)

	for i := 0; i < tsp.numNodes; i++ {
		if value, ok := tsp.precendense[i]; ok {
			r1 = append(r1, value)
			log.Print(value)
			r2 = append(r2, i)
			tmp[i] = true
			tmp[tsp.precendense[i]] = true
		}
	}

	/* 	for i := 1; i < tsp.numNodes; i++ {
		if _, ok := tmp[i]; !ok {
			r1 = append(r1, i)
		}
	} */

	for i := 1; i < len(r1); i++ {
		j := rand.Intn(i) + 1
		r1[i], r1[j] = r1[j], r1[i]
	}

	for i := 1; i < len(r2); i++ {
		j := rand.Intn(i) + 1
		r2[i], r2[j] = r2[j], r2[i]
	}

	r1 = append(r1, r2...)

	s := NewSolution(tsp, r1)

	s.Print()

	return &s
}
