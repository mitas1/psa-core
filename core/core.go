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

type Core struct {
	common       config.Common
	cons         *Construction
	optimization optimization
	objective    objective
}

func NewCore(c *config.Config) *Core {
	cons := NewCons(c.Construction)

	objective := NewObjective(c.Optimization)

	var optimization optimization

	if c.Optimization.VNS != (config.VNS{}) {
		optimization = NewVNS(c.Optimization.VNS, objective)
	} else {
		optimization = NewSA(c.Optimization.SA, objective)
	}
	return &Core{cons: cons, optimization: optimization, objective: objective, common: c.Common}
}

// Process PDPTW instance
func (c Core) Process(tsp *PDPTW) (*Solution, error) {
	iterationMax := c.common.IterMax
	i := 0
	iteration := 0

	// TODO: Check PDPTW instance

	// Preprocess incompatible arcs
	tsp.preprocess()

	// init structs

	var best, s *Solution

	channel := make(chan result)

	for i < iterationMax {
		go func() {
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

			s = c.optimization.process(s)
			channel <- result{solution: s, err: nil}

			s = c.optimization.process(s)
			channel <- result{solution: s, err: nil}
		}()
		i++
	}

	timeout := time.After(c.common.MaxTime * time.Second)

	for iteration < iterationMax*4 {
		iteration++
		select {
		case res := <-channel:
			if res.err != nil {
				return nil, res.err
			}

			if best == nil || c.objective.get(res.solution) < c.objective.get(best) {
				best = res.solution
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
