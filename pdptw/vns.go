package pdptw

import "log"

// VNS - Variable Neighborhood Search
func VNS(tsp *PDPTW, iterationMax, levelMax int) *Solution {
	var s *Solution

	tsp.preprocess()

	iteration := 0
	cons := Construction{}

	// Generate feasible solution
	best := cons.process(tsp)

	for iteration < iterationMax {
		iteration++

		//Generate feasible solution
		s = cons.process(tsp)

		// Try to improve

		// local2opt

		// Simmulated annealing with local 2 opt
		// sa := SA{}

		// s = sa.Process(s)

		local2opt := Local2Opt{inc: make([]int, len(tsp.matrix))}
		s = local2opt.Process(s)

		if s.MakeSpan() < best.MakeSpan() {
			best = s
			log.Print(s.MakeSpan())
		}
	}

	return best
}
