package pdptw

import (
	"math/rand"
	"time"

	"github.com/mitas1/psa-core/config"
)

// VNS - Variable Neighborhood Search
func VNS(tsp *PDPTW, config *config.Config) (s *Solution) {
	// set random seed
	rand.Seed(time.Now().UnixNano())

	iterationMax := config.VNS.IterMax
	iteration := 0

	// Preprocess incompatible arcs
	tsp.preprocess()

	// init structs
	cons := NewCons(config.Construction)

	opt := NewOptimization(config.Optimization, len(tsp.matrix))

	// Generate feasible solution
	best := cons.process(tsp)

	for iteration < iterationMax {
		iteration++

		//Generate feasible solution
		s = cons.process(tsp)

		// Try to improve

		// Simmulated annealing with local 2 opt
		// sa := SA{}

		// s = sa.Process(s)

		// local2opt
		s = opt.process(s)

		if opt.objective.get(s) < opt.objective.get(best) {
			best = s
		}
	}
	return best
}
