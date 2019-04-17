package core

import (
	"github.com/mitas1/psa-core/utils"
)

const (
	iterMax = 40
)

type localshifting struct {
	traveled   []int
	precedence map[int]int
	carrying   map[int]int
	objective  objective
}

func (local localshifting) process(x *Solution) {
	local.calcGlobals(x)

	for k := 0; k < iterMax; k++ {
		i := utils.Random(1, len(x.route)-1)

		for j := 1; j < len(x.route); j++ {
			if i != j {
				if local.isFeasible(x, i, j) == 0 {
					local.shift(x, i, j)
					break
				}
			}
		}
	}
	return
}

func (c *localshifting) calcGlobals(s *Solution) {
	var n1, n2 int

	traveled := s.tsp.traveled
	carrying := s.tsp.carrying

	c.traveled = make([]int, s.tsp.numNodes)
	c.precedence = make(map[int]int)
	c.carrying = make(map[int]int)

	for i := 0; i < len(s.route)-1; i++ {
		// traveled
		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readyTime[n1] > traveled {
			traveled = s.tsp.readyTime[n1]
		}

		traveled += s.tsp.matrix[n1][n2]

		c.traveled[i+1] = traveled

		// precedence
		if n, ok := s.tsp.precedence[n1]; ok {
			index := utils.IndexOf(n, s.route)
			c.precedence[index] = i
			c.precedence[i] = index
		} else if _, ok := c.precedence[i]; !ok {
			// ignore precedence of vertex
			c.precedence[i] = -1
		}

		carrying += s.tsp.demands[n1]

		c.carrying[i] = carrying
	}

	i := len(s.route) - 1

	n := s.tsp.precedence[s.route[i]]

	index := utils.IndexOf(n, s.route)
	c.precedence[i] = index
	c.precedence[index] = i

	return
}

func (c *localshifting) updateGlobals(s *Solution, from int) {
	var n1, n2 int

	traveled := c.traveled[from-1]
	carrying := c.carrying[from-1]

	for i := from - 1; i < len(s.route)-1; i++ {
		// traveled
		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readyTime[n1] > traveled {
			traveled = s.tsp.readyTime[n1]
		}

		traveled += s.tsp.matrix[n1][n2]

		c.traveled[i+1] = traveled

		// precedence
		if n, ok := s.tsp.precedence[n1]; ok {
			index := utils.IndexOf(n, s.route)
			c.precedence[index] = i
			c.precedence[i] = index
		} else if _, ok := c.precedence[i]; !ok {
			// ignore precedence of vertex
			c.precedence[i] = -1
		}

		carrying += s.tsp.demands[n2]

		c.carrying[i+1] = carrying
	}

	i := len(s.route) - 1

	n := s.tsp.precedence[s.route[i]]

	index := utils.IndexOf(n, s.route)
	c.precedence[i] = index
	c.precedence[index] = i

	return
}

func (c localshifting) isFeasible(s *Solution, pos, newPos int) int {
	var tail, traveled, carrying int

	predPos := c.precedence[pos]

	if newPos > pos {
		// FORWARD
		tail = newPos

		// precedence
		if predPos > pos {
			// move of A
			if newPos >= predPos {
				return -1
			}
		}

		// time window and capacity
		traveled = c.traveled[pos-1]
		carrying = c.carrying[pos-1]

		if !s.isFeasibleEdge(pos-1, pos+1, &traveled, &carrying) {
			return -1
		}

		if !s.isFeasibleRange(pos+1, newPos, &traveled, &carrying) {
			return -1
		}

		if !s.isFeasibleEdge(newPos, pos, &traveled, &carrying) {
			return -1
		}

		if traveled > c.traveled[newPos] {
			return -2
		}

		if newPos+1 < len(s.route) {
			if !s.isFeasibleEdge(pos, newPos+1, &traveled, &carrying) {
				return -1
			}

			if traveled > c.traveled[newPos+1] {
				return -2
			}
		}

	} else {
		// BACKWARD
		tail = pos

		// precedence
		if predPos < pos {
			// move of B
			if newPos <= predPos {
				return -1
			}
		}

		// time window and capacity
		traveled = c.traveled[newPos-1]
		carrying = c.carrying[newPos-1]

		if !s.isFeasibleEdge(newPos-1, pos, &traveled, &carrying) {
			return -1
		}

		if !s.isFeasibleEdge(pos, newPos, &traveled, &carrying) {
			return -1
		}

		if !s.isFeasibleRange(newPos, pos-1, &traveled, &carrying) {
			return -1
		}

		if traveled > c.traveled[pos] {
			return -2
		}

		if pos+1 < len(s.route) {
			if !s.isFeasibleEdge(pos-1, pos+1, &traveled, &carrying) {
				return -1
			}

			if traveled > c.traveled[pos+1] {
				return -2
			}
		}
	}

	if !s.isFeasibleRange(tail+1, len(s.route)-1, &traveled, &carrying) {
		return -1
	}

	return 0
}

func (local localshifting) shift(x *Solution, pos, newPos int) {
	var from int
	node := x.route[pos]
	if newPos > pos {
		from = pos
		// BACKWARD
		for i := pos; i < newPos; i++ {
			x.route[i] = x.route[i+1]
		}
	} else {
		// FORWARD
		from = newPos
		for i := pos; i > newPos; i-- {
			x.route[i] = x.route[i-1]
		}

	}
	x.route[newPos] = node
	local.updateGlobals(x, from)
}
