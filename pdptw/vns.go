package pdptw

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mitas1/psa-core/config"
)

type result struct {
	solution *Solution
	err      error
}

type vns struct {
	config *config.Config
}

func NewVNS(config *config.Config) vns {
	return vns{config}
}

// VNS - Variable Neighborhood Search
func (v vns) Process(tsp *PDPTW) (*Solution, error) {
	iterationMax := v.config.VNS.IterMax
	iteration := 0

	// TODO: Check TSP instance

	// Preprocess incompatible arcs
	tsp.preprocess()

	// init structs
	cons := NewCons(v.config.Construction)

	opt := NewOptimization(v.config.Optimization, len(tsp.matrix))

	var best, s *Solution

	channel := make(chan result)

	go func() {
		for {
			// set random seed
			rand.Seed(time.Now().UnixNano())

			//Generate feasible solution
			s = cons.process(tsp)
			channel <- result{solution: s, err: nil}

			// Try to improve

			// Simmulated annealing with local 2 opt
			// sa := SA{}

			// s = sa.Process(s)

			s = opt.process(s)
			channel <- result{solution: s, err: nil}
		}
	}()

	timeout := time.After(v.config.VNS.MaxTime * time.Second)

	for iteration < iterationMax*2 {
		iteration++
		select {
		case res := <-channel:
			if res.err != nil {
				return nil, res.err
			}
			if best == nil || opt.objective.get(s) < opt.objective.get(best) {
				best = s
			}
		case <-timeout:
			if best == nil {
				// raise timeout error
				return nil, fmt.Errorf("timeout: Unable to find solution")
			} else {
				return best, fmt.Errorf("timeout: Only partial solution found")
			}
		}
	}

	return best, nil
}
