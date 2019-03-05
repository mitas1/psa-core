package pdptw

import (
	"log"

	"github.com/mitas1/psa-core/config"
)

// VNS - Variable Neighborhood Search
func VNS(tsp *PDPTW, config *config.Config) (s *Solution) {
	iterationMax := config.VNS.IterMax
	iteration := 0

	// Preprocess incompatible arcs
	tsp.preprocess()

	// init structs
	cons := NewCons(config.Construction)
	opt := NewOptimization(config.Optimization, len(tsp.matrix))

	// Generate feasible solution
	best := cons.process(tsp)
	return best

	for iteration < iterationMax {
		iteration++

		//Generate feasible solution
		s = cons.process(tsp)

		// Try to improve

		// local2opt

		// Simmulated annealing with local 2 opt
		// sa := SA{}

		// s = sa.Process(s)

		s = opt.Process(s)

		log.Print(s.IsFeasibleLog())

		if s.MakeSpan() < best.MakeSpan() {
			best = s
			log.Print(s.MakeSpan())
		}
	}

	return best
}
