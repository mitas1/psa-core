package core

import (
	"math/rand"
)

// constrained 2 opt
type local2Opt struct {
	traveled   []int
	precedence map[int]int
	carrying   map[int]int
	objective
}

func (c local2Opt) process(s *Solution) {
	var pos, i int

	numNodes := s.tsp.numNodes
	// create auxiliary set
	pointer := numNodes - 2
	set := make([]int, pointer)

	for i := 0; i < pointer; i++ {
		set[i] = i
	}

	c.setGlobals(s.calcGlobals())

	// outerloop
	for pointer > 0 {
	outer:
		// generate random outerloop i
		pos = rand.Intn(pointer)
		i = set[pos]

		// iner loop
		for j := i + 2; j < numNodes-1; j++ {
			if c.objective.isProfitable(s, i, j, c.traveled[j+1], c.traveled[i]) {
				if c.isFeasible(s, i, j) {
					c.exchangeGlobalUpdate(s, i, j)
					pointer = numNodes - 2
					goto outer
				}
				break
			}
		}
		pointer--
		set[pointer], set[pos] = i, set[pointer]
	}
	return
}

func (c local2Opt) exchangeGlobalUpdate(s *Solution, iaux, jaux int) {
	start := iaux + 1
	end := jaux

	j := 0

	median := (end-start+1)/2 + start - 1

	for i := start; i <= median; i++ {
		j = end - (i - start)
		// exchange
		s.route[i], s.route[j] = s.route[j], s.route[i]

		// update precedence
		c.precedence[c.precedence[i]], c.precedence[c.precedence[j]],
			c.precedence[j], c.precedence[i] = j, i, c.precedence[i], c.precedence[j]
	}

	var n1, n2 int

	sum := c.traveled[iaux]
	carrying := c.carrying[iaux-1]

	// update the reversed path
	for i := iaux; i < jaux; i++ {
		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readyTime[n1] > sum {
			sum = s.tsp.readyTime[n1]
		}

		sum += s.tsp.matrix[n1][n2]
		carrying += s.tsp.demands[n1]

		c.traveled[i+1] = sum
		c.carrying[i] = carrying
	}

	// update the rest
	for i := jaux; i < len(s.route)-1; i++ {
		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readyTime[n1] > sum {
			sum = s.tsp.readyTime[n1]
		}

		sum += s.tsp.matrix[n1][n2]
		carrying += s.tsp.demands[n1]

		c.traveled[i+1] = sum
		c.carrying[i] = carrying
	}

	return
}

// exchange (i,i+1), (j,j+1) ===> (i,j), (i+1,j+1)
func (c local2Opt) isFeasible(s *Solution, i, j int) bool {
	var n1, n2 int
	sum := c.traveled[i]
	carrying := c.carrying[i]

	if !s.isFeasibleEdge(i, j, &sum, &carrying) {
		return false
	}

	for k := j; k > i+1; k-- {

		n1 = s.route[k]
		n2 = s.route[k-1]

		if s.tsp.readyTime[n1] > sum {
			sum = s.tsp.readyTime[n1]
		}

		sum += s.tsp.matrix[n1][n2]

		if sum > s.tsp.dueDate[n2] {
			return false
		}

		// precendence

		if c.precedence[k] > i && c.precedence[k] < j {
			return false
		}

		// capacity

		carrying += s.tsp.demands[n2]

		if carrying > s.tsp.capacity {
			return false
		}
	}

	if j+1 < len(s.route) {
		if !s.isFeasibleEdge(i+1, j+1, &sum, &carrying) {
			return false
		}
	}

	if !s.isFeasibleRange(j+1, len(s.route)-1, &sum, &carrying) {
		return false
	}

	return true
}

func (c *local2Opt) setGlobals(traveled []int, carrying, precedence map[int]int) {
	c.traveled = traveled
	c.precedence = precedence
	c.carrying = carrying
}
