package core

import (
	"github.com/mitas1/psa-core/config"
)

// interface of local 2opt search
type localSearch interface {
	process(*Solution)
}

func getLocalSearch(local config.LocalSearch, objective objective) localSearch {
	localShift := localshifting{objective: objective}
	local2Opt := local2Opt{objective: objective}

	switch local {

	case config.VND:
		return vnd{
			objective:     objective,
			local2Opt:     local2Opt,
			localShifting: localShift}
	case config.Shifting:
		return localShift
	default:
		return local2Opt
	}
}
