package core

import (
	"math/rand"

	"github.com/mitas1/psa-core/config"
	"github.com/mitas1/psa-core/utils"
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
	penalty  config.Penalty
}

func NewCons(opts config.Construction) *Construction {
	var strategy constructionStrategy

	switch opts.Strategy {
	case "random":
		strategy = random{}
	case "greedy":
		strategy = greedy{}
	case "sortBydueDate":
		strategy = sortBydueDate{}
	case "sortByTW":
		strategy = sortByTW{}
	default:
		strategy = random{}
	}

	return &Construction{levelMax: opts.LevelMax, strategy: strategy, penalty: opts.Penalty}
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

		if c.Penalty(x2) < c.Penalty(x) {
			for i := 1; i < len(x.route); i++ {
				x.route[i] = x2.route[i]
			}

			level = 1

			if c.Penalty(x) == 0 {
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

func (c *Construction) shifting(s *Solution, direction Direction, setType SetType) (penalty int, improved bool) {
	var newPenalty int
	// get feasible or unfeasible set
	set := s.getSet(setType)
	// calculate penalty
	penalty = c.Penalty(s)
	// set pointer to set length
	pointer := len(set)
	improvement := false

	for pointer > 0 {
		// select random node within set
		pos := rand.Intn(pointer)
		npos := set[pos]
		n := s.route[npos]

		improvement = false

		if direction == BACKWARD {
			for i := npos - 1; i > 0; i-- {
				if !s.tsp.arcs[n][s.route[i]] {
					break
				}

				s.exchange(npos, i)

				newPenalty = c.Penalty(s)

				if newPenalty < penalty {
					improvement = true
					break
				}
				s.exchange(i, npos)
			}
		} else {
			for i := npos + 1; i < len(s.route)-1; i++ {
				if !s.tsp.arcs[s.route[i]][n] {
					break
				}

				s.exchange(npos, i)

				newPenalty = c.Penalty(s)

				if newPenalty < penalty {
					improvement = true
					break
				}
				s.exchange(i, npos)
			}
		}

		if improvement {
			if newPenalty == 0 {
				return
			}
			penalty = newPenalty
			// reset after improvement
			set = s.getSet(setType)
			pointer = len(set)
			// improvement found
			improved = true
		} else {
			pointer--
			// move selected node to end
			set[pos], set[pointer] = set[pointer], set[pos]
		}
	}
	return
}

func (*Construction) disturb(s *Solution, level int) *Solution {
	var n1, n2 int
	newSolution := s.Copy()

	for i := 0; i < level; i++ {
		n1 = utils.Random(1, len(s.route))
		n2 = utils.Random(1, len(s.route))
		newSolution.exchange(n1, n2)
	}
	return newSolution
}

// Penalty is sum of all differences between the time to reach each customer
// and its due date
func (c Construction) Penalty(s *Solution) (penalty int) {
	traveled := s.tsp.traveled
	carrying := s.tsp.carrying
	hasNode := false

	p_tw := 0
	p_pd := 0
	p_c := 0

	for i := 1; i < len(s.route); i++ {

		hasNode = false
		traveled += s.tsp.matrix[s.route[i-1]][s.route[i]]
		carrying += s.tsp.demands[s.route[i-1]]

		// wait to ready to time
		if traveled < s.tsp.readyTime[s.route[i]] {
			traveled = s.tsp.readyTime[s.route[i]]
		}

		if carrying > s.tsp.capacity {
			p_c = p_c + (carrying - s.tsp.capacity)
		}

		if value, ok := s.tsp.precedence[s.route[i]]; ok {
			for j := i; j > 0; j-- {
				if value == s.route[j] {
					hasNode = true
					break
				}
			}
			if !hasNode {
				for j := i; j < s.tsp.numNodes; j++ {
					if value == s.route[j] {
						p_pd = p_pd + j
						break
					}
				}
			}
		}

		if s.tsp.dueDate[s.route[i]] != 0 && s.tsp.dueDate[s.route[i]] < traveled ||
			s.tsp.readyTime[s.route[i]] > traveled {
			p_tw = p_tw + traveled - s.tsp.dueDate[s.route[i]]
		}
	}

	penalty = c.penalty.TimeWindows*p_tw + c.penalty.PickupDelivery*p_pd + c.penalty.Capacity*p_c
	return
}
