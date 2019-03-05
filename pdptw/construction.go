package pdptw

import (
	"math/rand"

	"github.com/mitas1/psa-core/config"
)

// Direction for shifting
type Direction int32

// FORWARD AND BACKWARD direction
const (
	FORWARD  Direction = 0
	BACKWARD Direction = 1
)

type constructionStrategy interface {
	getSolution(tsp *PDPTW) *Solution
}

// Construction struct
type Construction struct {
	levelMax int
	strategy constructionStrategy
}

func NewCons(opts config.Construction) *Construction {
	var strategy constructionStrategy

	switch opts.Strategy {
	case "random":
		strategy = random{}
	case "greedy":
		strategy = greedy{}
	default:
		strategy = random{}
	}

	return &Construction{levelMax: opts.LevelMax, strategy: strategy}
}

func (c *Construction) process(tsp *PDPTW) *Solution {
	level := 1

	// generate first solution, using configured strategy
	x := c.strategy.getSolution(tsp)

	// c.localSearch(x)

	for !x.IsFeasible() {
		x2 := c.disturb(x, level)

		c.localSearch(x2)

		if x2.IsFeasible() {
			return x2
		}

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

			if c.levelMax < level {
				level = 1
				x = c.strategy.getSolution(tsp)
			}
		}
	}

	return x
}

func (c *Construction) localSearch(s *Solution) {
	penalty := 1

	// variables represent whether the improvement was found
	i1, i2, i3, i4 := true, true, true, true

	// break only if no shifting found improvement or penalty is 0
	for i1 || i2 || i3 || i4 {
		if penalty, i1 = c.shifting(s, BACKWARD, UNFEASIBLE_SET); penalty == 0 {
			break
		}
		if penalty, i2 = c.shifting(s, FORWARD, FEASIBLE_SET); penalty == 0 {
			break
		}
		if penalty, i3 = c.shifting(s, FORWARD, UNFEASIBLE_SET); penalty == 0 {
			break
		}
		if penalty, i4 = c.shifting(s, BACKWARD, FEASIBLE_SET); penalty == 0 {
			break
		}
	}
	return
}

func (*Construction) shifting(s *Solution, direction Direction, setType SetType) (penalty int, improved bool) {
	// get feasible or unfeasible set
	set := s.getSet(setType)
	// calculate penalty
	penalty = s.Penalty()
	// set pointer to set length
	pointer := len(set)
	improvement := false

	for pointer > 0 {
		// select random node within set
		pos := rand.Intn(pointer)
		npos := set[pos]
		n := s.route[npos]

		improvement = false

		// TODO: Rewrite needed, merge this code
		if direction == BACKWARD {
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
		} else {
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
		}

		if improvement {
			if penalty == 0 {
				return
			}
			// reset after improvement
			set = s.getSet(setType)
			pointer = len(set)
			// improvement found
			improved = true
		} else {
			pointer--
			// move selected node to end
			aux := set[pos]
			set[pos] = set[pointer]
			set[pointer] = aux
		}
	}
	return
}

// TODO: rewrite needed
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
