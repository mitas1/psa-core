package pdptw

import (
	"math"
	"math/rand"
)

// returns random solution
type random struct{}

func (random) getSolution(tsp *PDPTW) *Solution {
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

	s := NewSolution(tsp, route)

	return &s
}

// returns solution constructed by nearest neighborhood heuristic
type greedy struct{}

func (greedy) getSolution(tsp *PDPTW) *Solution {
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

func GetRandomPD(tsp *PDPTW) *Solution {
	r1 := []int{0}
	r2 := []int{}

	tmp := make(map[int]bool, tsp.numNodes)

	for i := 0; i < tsp.numNodes; i++ {
		if value, ok := tsp.precendense[i]; ok {
			r1 = append(r1, value)
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

	return &s
}
