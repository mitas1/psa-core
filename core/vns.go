package core

import (
	"github.com/mitas1/psa-core/config"
)

// NewVNS returns VNS optimization strategy
func NewVNS(opts config.VNS, obj objective) vns {
	localSearch := getLocalSearch(opts.LocalSearch, obj)
	return vns{
		search:    localSearch,
		levelMax:  opts.LevelMax,
		iterMax:   opts.IterMax,
		objective: obj}
}

type vns struct {
	search   localSearch
	levelMax int
	iterMax  int
	objective
}

// process search
func (local vns) process(x *Solution) *Solution {
	var x2 *Solution
	level := 1
	iterLevel := 0

	best := x.Copy()

	local.search.process(best)

	for level < local.levelMax {

		x2 = best.disturb(level)

		local.search.process(x2)

		if local.objective.get(x2) < local.objective.get(best) {
			iterLevel = 0
			level = 1
			best = x2
		} else {
			if iterLevel > local.iterMax {
				level++
				iterLevel = 0
			}
		}

		iterLevel++
	}

	return best
}
