package core

import (
	"math"
	"math/rand"
)

type SA struct {
	search    localSearch
	objective objective
}

func NewSA(obj objective) SA {
	return SA{objective: obj, search: cons2Opt{objective: obj}}
}

func (local SA) process(state *Solution) *Solution {
	iterMax := 30.0
	iter := 0.0
	T := 0.0
	fraction := 0.0
	newCost := 0.0

	cost := float64(local.objective.get(state))
	var newState *Solution

	for iter < iterMax {
		iter++

		fraction = iter / iterMax
		T = local.temperature(fraction)
		// log.Printf("temperature: %v", T)

		// log.Printf("Fraction - %v", int(fraction*100))
		newState = state.disturb(int(fraction * 100))
		//disturb := newState.MakeSpan()
		local.search.process(newState)

		newCost = float64(local.objective.get(newState))

		if local.probability(cost, newCost, T) > rand.Float64() {
			//log.Printf("%v - %v - update", cost, newCost)
			state = newState
			cost = newCost
		}
	}

	return state
}

func (local *SA) temperature(fraction float64) float64 {
	return math.Max(0.01, math.Min(1, 1-fraction))
}

func (local *SA) probability(cost, newCost, T float64) float64 {
	if newCost < cost {
		return 1
	}

	return math.Exp(-(newCost - cost) / T)
}
