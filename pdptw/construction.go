package pdptw

import (
	"math/rand"
)

type Construction struct{}

func (c *Construction) process(tsp *PDPTW) *Solution {
	level := 1
	levelMax := 25

	// generate random solution
	x := GetRandom(tsp)

	c.local1shift(x)

	for !x.IsFeasible() {
		x2 := c.disturb(x, level)

		/* local2opt := Local2Opt{inc: make([]int, len(tsp.matrix))}
		x2 = local2opt.Process(x2) */

		x2 = c.local1shift(x2)

		//log.Printf("%v - %v", x2.Penalty(), x2.IsFeasible())

		if x2.Penalty() < x.Penalty() {
			for i := 1; i < len(x.route); i++ {
				x.route[i] = x2.route[i]
			}

			level = 1

			if x.Penalty() == 0 {
				return x
			}
		} else {
			level++

			if levelMax < level {
				level = 1
				x = GetRandom(tsp)
			}
		}
	}

	return x
}

func (*Construction) local1shift(s *Solution) *Solution {
	penalty := s.Penalty()

	feasibleSet, unfeasibleSet := s.calcSets()
	unfeasibleSize := len(unfeasibleSet)
	feasibleSize := len(feasibleSet)

	improvement := false

	for penalty > 0 && (unfeasibleSize > 0 || feasibleSize > 0) {

		feasibleSet, unfeasibleSet = s.calcSets()
		unfeasibleSize = len(unfeasibleSet)
		feasibleSize = len(feasibleSet)

		//ufAuxSize := unfeasibleSize
		//fAuxSize := feasibleSize

		for unfeasibleSize > 0 {
			pos := rand.Intn(len(unfeasibleSet))
			npos := unfeasibleSet[pos]
			//n := s.route[npos]

			improvement = false

			for i := npos - 1; i > 0; i-- {
				/* if !s.tsp.arcs[n][s.route[i]] {
					break
				} */

				s.exchange(npos, i)

				penaltyAux := s.Penalty()

				if penaltyAux < penalty {
					penalty = penaltyAux
					improvement = true
					break
				} else {
					s.exchange(i, npos)
				}
			}

			if improvement {
				if penalty == 0 {
					return s
				}

				feasibleSet, unfeasibleSet = s.calcSets()
				unfeasibleSize = len(unfeasibleSet)
				feasibleSize = len(feasibleSet)
			} else {
				unfeasibleSize--
				aux := unfeasibleSet[pos]
				unfeasibleSet[pos] = unfeasibleSet[unfeasibleSize]
				unfeasibleSet[unfeasibleSize] = aux
			}

		}

		for feasibleSize > 0 {
			pos := rand.Intn(len(feasibleSet))
			npos := feasibleSet[pos]
			n := s.route[npos]

			improvement = false

			for i := npos + 1; i < s.tsp.numNodes-1; i++ {
				if !s.tsp.arcs[s.route[i]][n] {
					break
				}

				s.exchange(npos, i)

				penaltyAux := s.Penalty()

				if penaltyAux < penalty {
					penalty = penaltyAux
					improvement = true
					break
				} else {
					s.exchange(i, npos)
				}
			}

			if improvement {
				if penalty == 0 {
					return s
				}
				feasibleSet, unfeasibleSet = s.calcSets()
				unfeasibleSize = len(unfeasibleSet)
				feasibleSize = len(feasibleSet)
			} else {
				feasibleSize--
				aux := feasibleSet[pos]
				feasibleSet[pos] = feasibleSet[feasibleSize]
				feasibleSet[feasibleSize] = aux
			}
		}

		//feasibleSize = fAuxSize
		// unfeasibleSize = ufAuxSize

		for unfeasibleSize > 0 {
			pos := rand.Intn(unfeasibleSize)
			npos := unfeasibleSet[pos]
			n := s.route[npos]

			improvement = false

			for i := npos + 1; i < len(s.route)-1; i++ {
				if !s.tsp.arcs[s.route[i]][n] {
					break
				}

				s.exchange(npos, i)

				penaltyAux := s.Penalty()

				if penaltyAux < penalty {
					penalty = penaltyAux
					improvement = true
					break
				} else {
					s.exchange(i, npos)
				}
			}

			if improvement {
				if penalty == 0 {
					return s
				}

				feasibleSet, unfeasibleSet = s.calcSets()
				unfeasibleSize = len(unfeasibleSet)
				feasibleSize = len(feasibleSet)
			} else {
				unfeasibleSize--
				aux := unfeasibleSet[pos]
				unfeasibleSet[pos] = unfeasibleSet[unfeasibleSize]
				unfeasibleSet[unfeasibleSize] = aux
			}
		}

		for feasibleSize > 0 {
			pos := rand.Intn(len(feasibleSet))
			npos := feasibleSet[pos]
			n := s.route[npos]

			improvement = false

			for i := npos - 1; i > 0; i-- {
				if !s.tsp.arcs[n][s.route[i]] {
					break
				}

				s.exchange(npos, i)

				penaltyAux := s.Penalty()

				if penaltyAux < penalty {
					penalty = penaltyAux
					improvement = true
					break
				} else {
					s.exchange(i, npos)
				}
			}

			if improvement {
				if penalty == 0 {
					return s
				}
				feasibleSet, unfeasibleSet = s.calcSets()
				unfeasibleSize = len(unfeasibleSet)
				feasibleSize = len(feasibleSet)
			} else {
				feasibleSize--
				aux := feasibleSet[pos]
				feasibleSet[pos] = feasibleSet[feasibleSize]
				feasibleSet[feasibleSize] = aux
			}
		}
	}
	return s
}
func indexOf(element int, data []int) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func (*Construction) disturb(s *Solution, level int) *Solution {
	newSolution := NewSolution(s.tsp, s.route)

	feasibleSet, unfeasibleSet := s.calcSets()

	feasibleSize, unfeasibleSize := len(feasibleSet), len(unfeasibleSet)

	for i := 0; i <= level; i++ {
		if feasibleSize > 0 && unfeasibleSize > 0 {
			newSolution.exchange(feasibleSet[rand.Intn(feasibleSize)],
				unfeasibleSet[rand.Intn(unfeasibleSize)])
		}
	}
	return &newSolution
}
