package pdp

import (
	"math/rand"
)

func GetRandom(tsp *PDP) *Solution {
	list := rand.Perm(tsp.numNodes - 2)

	for i := range list {
		list[i] += 2
	}

	list = append(list, 1)

	route := []int{0}
	route = append(route, list...)

	s := NewSolution(tsp, route)

	return &s
}

func GetWellRandom(tsp *PDP) *Solution {

	x := GetRandom(tsp)

	return x
}
