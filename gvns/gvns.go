package gvns

// GVNS - General Variable Neighborhood Search
func GVNS(tsp *TSP, iterationMax, levelMax int) *Solution {
	var s *Solution

	tsp.preprocess()

	iteration := 0
	cons := Construction{}

	// Generate feasible solution
	best := cons.process(tsp)

	for iteration < iterationMax {
		iteration++

		// Generate feasible solution
		s = cons.process(tsp)

		// Try to improve
		local2opt := Local2Opt{inc: make([]int, len(tsp.matrix))}
		local2opt.process(s)

		if s.TotalDistance() < best.TotalDistance() {
			best = s
		}
	}

	return best
}
