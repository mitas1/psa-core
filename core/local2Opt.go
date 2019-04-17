package core

import (
	"math/rand"

	"github.com/mitas1/psa-core/utils"
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

	c.calcGlobals(s)

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

	sum := c.traveled[i]

	carrying := c.carrying[i]

	// i -> j
	n1 := s.route[i]
	n2 := s.route[j]

	if s.tsp.readyTime[n1] > sum {
		sum = s.tsp.readyTime[n1]
	}

	sum += s.tsp.matrix[n1][n2]

	if sum > s.tsp.dueDate[n2] {
		return false
	}

	if carrying > s.tsp.capacity {
		return false
	}

	// j -> i + 1

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
			// log.Printf("EXCHANGE: i=%v j=%v k=%v %v| %v", i, j, k, c.precedence[k], c.precedence)
			return false
		}

		// capacity

		carrying += s.tsp.demands[n1]

		if carrying > s.tsp.capacity {
			return false
		}
	}
	// i+1 -> j+1
	if j+1 < len(s.route) {
		n1 = s.route[i+1]
		n2 = s.route[j+1]

		if s.tsp.readyTime[n1] > sum {
			sum = s.tsp.readyTime[n1]
		}

		sum += s.tsp.matrix[n1][n2]
		carrying += s.tsp.demands[n1]

		if sum > s.tsp.dueDate[n2] {
			return false
		}

		if carrying > s.tsp.capacity {
			return false
		}
	}

	// j+1 ->

	for k := j + 1; k < len(s.route)-1; k++ {
		n1 = s.route[k]
		n2 = s.route[k+1]

		if s.tsp.readyTime[n1] > sum {
			sum = s.tsp.readyTime[n1]
		}

		sum += s.tsp.matrix[n1][n2]
		carrying += s.tsp.demands[n1]

		if sum > s.tsp.dueDate[n2] {
			return false
		}

		if carrying > s.tsp.capacity {
			return false
		}
	}

	if s.tsp.readyTime[s.route[len(s.route)-2]] > sum {
		sum = s.tsp.readyTime[s.route[len(s.route)-2]]
	}

	return true
}

func (c *local2Opt) calcGlobals(s *Solution) {
	var n1, n2 int

	sum := s.tsp.traveled
	carrying := s.tsp.carrying

	c.traveled = make([]int, s.tsp.numNodes)
	c.precedence = make(map[int]int)
	c.carrying = make(map[int]int)

	for i := 0; i < len(s.route)-1; i++ {
		// traveled
		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readyTime[n1] > sum {
			sum = s.tsp.readyTime[n1]
		}

		sum += s.tsp.matrix[n1][n2]

		c.traveled[i+1] = sum

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
