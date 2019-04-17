package core

import (
	"math"
	"math/rand"

	"github.com/mitas1/psa-core/config"
)

type SA struct {
	iterMax   float64
	search    localSearch
	objective objective
}

func NewSA(opts config.SA, obj objective) SA {
	localSearch := getLocalSearch(opts.LocalSearch, obj)
	return SA{objective: obj, search: localSearch, iterMax: opts.IterMax}
}

func (local SA) process(state *Solution) *Solution {
	iter := 0.0
	T := 0.0
	fraction := 0.0
	newCost := 0.0

	cost := float64(local.objective.get(state))
	var newState *Solution

	for iter < local.iterMax {
		iter++

		fraction = iter / iterMax
		T = local.temperature(fraction)

		newState = state.disturb(int(fraction * 100))
		local.search.process(newState)

		newCost = float64(local.objective.get(newState))

		if local.probability(cost, newCost, T) > rand.Float64() {
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
