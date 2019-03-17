package core

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mitas1/psa-core/config"
)

type optimization interface {
	process(*Solution) *Solution
}

type result struct {
	solution *Solution
	err      error
}

type core struct {
	common       config.Common
	cons         *Construction
	optimization optimization
	objective    objective
}

func NewCore(c *config.Config) core {
	cons := NewCons(c.Construction)

	objective := NewObjective(c.Optimization)

	var optimization optimization

	if c.Optimization.GVNS != (config.GVNS{}) {
		optimization = NewGVNS(c.Optimization.GVNS, objective)
	} else {
		optimization = NewSA(objective)
	}
	return core{cons: cons, optimization: optimization, objective: objective, common: c.Common}
}

// Process PDPTW instance
func (c core) Process(tsp *PDPTW) (*Solution, error) {
	iterationMax := c.common.IterMax
	iteration := 0

	// TODO: Check PDPTW instance

	// Preprocess incompatible arcs
	tsp.preprocess()

	// init structs

	var best, s *Solution

	channel := make(chan result)

	go func() {
		for {
			// set random seed
			rand.Seed(time.Now().UnixNano())

			//Generate feasible solution
			s = c.cons.process(tsp)
			channel <- result{solution: s, err: nil}

			// Try to improve

			// Simmulated annealing with local 2 optimization
			// sa := SA{}

			// s = sa.Process(s)

			s = c.optimization.process(s)
			channel <- result{solution: s, err: nil}
		}
	}()

	timeout := time.After(c.common.MaxTime * time.Second)

	for iteration < iterationMax*2 {
		iteration++
		select {
		case res := <-channel:
			if res.err != nil {
				return nil, res.err
			}
			if best == nil || c.objective.get(s) < c.objective.get(best) {
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
