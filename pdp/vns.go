package pdp

// VNS - Variable Neighborhood Search
func VNS(tsp *PDP, iterationMax, levelMax int, localOpt bool) *Solution {
	var s *Solution

	tsp.preprocess()

	iteration := 0
	cons := Construction{}
	local2opt := Local2Opt{inc: make([]int, len(tsp.matrix))}

	// Generate feasible solution
	best := cons.Process(tsp)

	for iteration < iterationMax {
		iteration++

		// Generate feasible solution
		s = cons.Process(tsp)

		// Try to improve
		if localOpt {
			s = local2opt.Process(s)
		}

		if s.TotalDistance() < best.TotalDistance() {
			best = s
		}
	}

	return best
}
