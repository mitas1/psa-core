package pdptw

import (
	"math"
	"math/rand"
	"sort"
)

type random struct{}

// generate random solution
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

type sortByDuedate struct{}

func (sortByDuedate) getSolution(tsp *PDPTW) *Solution {
	var route []int

	for i := 0; i < tsp.numNodes; i++ {
		if i != tsp.startNode {
			route = append(route, i)
		}
	}

	// sort by duedate
	sort.Slice(route, func(i, j int) bool {
		return tsp.duedate[route[i]] < tsp.duedate[route[j]]
	})

	s := NewSolution(tsp, append([]int{tsp.startNode}, route...))

	return &s
}

type sortByTW struct{}

func (sortByTW) getSolution(tsp *PDPTW) *Solution {
	route := []int{}

	median := make(map[int]int)

	for i := 0; i < tsp.numNodes; i++ {
		if i != tsp.startNode {
			route = append(route, i)
			median[i] = tsp.duedate[i] - ((tsp.duedate[i] - tsp.readytime[i]) / 2)
		}
	}

	// sort by median
	sort.Slice(route, func(i, j int) bool {
		return median[route[i]] < median[route[j]]
	})

	s := NewSolution(tsp, append([]int{tsp.startNode}, route...))

	return &s
}

type greedy struct{}

// returns solution constructed by nearest neighborhood heuristic
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
		if value, ok := tsp.precedence[i]; ok {
			r1 = append(r1, value)
			r2 = append(r2, i)
			tmp[i] = true
			tmp[tsp.precedence[i]] = true
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
