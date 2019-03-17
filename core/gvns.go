package core

import (
	"github.com/mitas1/psa-core/config"
)

// NewGVNS returns GVNS optimization strategy
func NewGVNS(opts config.GVNS, obj objective) gvns {
	return gvns{
		levelMax:  opts.LevelMax,
		iterMax:   opts.IterMax,
		objective: obj,
		search:    cons2Opt{objective: obj}}
}

type gvns struct {
	search   localSearch
	levelMax int
	iterMax  int
	objective
}

// process search
func (local gvns) process(x *Solution) *Solution {
	level := 1
	iterLevel := 0

	local.search.process(x)

	for level < local.levelMax {
		x2 := x.disturb(level)

		local.search.process(x2)

		if local.objective.get(x2) < local.objective.get(x) {
			iterLevel = 0
			level = 1
			x = x2
		} else {
			if iterLevel > local.iterMax {
				level++
				iterLevel = 0
			}
		}
		iterLevel++
	}

	return x
}
