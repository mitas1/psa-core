package pdp

import (
	"math/rand"
)

type Construction struct{}

func (c *Construction) Process(tsp *PDP) *Solution {
	level := 1
	levelMax := 8

	// generate random solution
	x := GetRandom(tsp)
	x.PenaltySort()

	c.local1shift(x)

	for !x.IsFeasible() {
		x2 := c.disturb(x, level)

		c.local1shift(x2)

		if x2.Penalty() < x.Penalty() {
			for i := 1; i < len(x.route); i++ {
				x.route[i] = x2.route[i]
			}
			level = 1

			if x.Penalty() == x2.Penalty() {
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

func (*Construction) local1shift(s *Solution) {
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
					return
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
			//n := s.route[npos]

			improvement = false

			for i := npos + 1; i < s.tsp.numNodes-1; i++ {

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
					return
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
			//n := s.route[npos]

			improvement = false

			for i := npos + 1; i < len(s.route)-1; i++ {

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
					return
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
			//n := s.route[npos]

			improvement = false

			for i := npos - 1; i > 0; i-- {

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
					return
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
