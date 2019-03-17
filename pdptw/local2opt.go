package pdptw

import (
	"math/rand"

	"github.com/mitas1/psa-core/config"
	"github.com/mitas1/psa-core/utils"
)

// NewOptimization returns local2opt optimization strategy
func NewOptimization(opts config.Optimization, numNodes int) local2Opt {
	var obj objective
	switch {
	case "time" == opts.Objective && opts.Asymetric:
		obj = totalTimeA{}
	case "time" == opts.Objective:
		obj = totalTime{}
	case "span" == opts.Objective:
		obj = spanTime{}
	default:
		obj = spanTime{}
	}

	switch opts.Strategy {
	default:
		return local2Opt{
			levelMax:  opts.LevelMax,
			iterMax:   opts.IterMax,
			objective: obj,
			strategy: cons2Opt{
				traveled:   make([]int, numNodes),
				precedence: make(map[int]int),
				carrying:   make(map[int]int),
				objective:  obj}}
	}

	return local2Opt{}
}

// interface of local 2opt search
type localSearch interface {
	process(*Solution)
}

type local2Opt struct {
	strategy localSearch
	levelMax int
	iterMax  int
	objective
}

func (spanTime) get(s *Solution) int {
	traveled := 0
	for i := 0; i < len(s.route)-1; i++ {
		if traveled < s.tsp.readyTime[s.route[i]] {
			traveled = s.tsp.readyTime[s.route[i]]
		}
		traveled += s.tsp.matrix[s.route[i]][s.route[i+1]]
	}
	return traveled
}

func (spanTime) isProfitable(s *Solution, i, j int, spans ...int) bool {
	var n1, n2 int

	sum := spans[1]
	n1 = s.route[i]
	n2 = s.route[j]

	if s.tsp.readyTime[n1] > sum {
		sum = s.tsp.readyTime[n1]
	}

	sum += s.tsp.matrix[n1][n2]

	for k := j; k > i+1; k-- {
		n1 = s.route[k]
		n2 = s.route[k-1]

		if s.tsp.readyTime[n1] > sum {
			sum = s.tsp.readyTime[n1]
		}

		sum += s.tsp.matrix[n1][n2]
	}

	n1 = s.route[i+1]
	n2 = s.route[j+1]

	if s.tsp.readyTime[n1] > sum {
		sum = s.tsp.readyTime[n1]
	}

	sum += s.tsp.matrix[n1][n2]

	return spans[0] > sum
}

func (totalTime) get(s *Solution) int {
	traveled := 0
	for i := 0; i < len(s.route)-1; i++ {
		traveled += s.tsp.matrix[s.route[i]][s.route[i+1]]
	}
	return traveled
}

func (totalTime) isProfitable(s *Solution, i, j int, spans ...int) bool {
	n1 := s.route[i]
	n2 := s.route[i+1]
	n3 := s.route[j]
	n4 := 0

	if j < s.tsp.numNodes-1 {
		n4 = s.route[j+1]
	}

	e1 := s.tsp.matrix[n1][n2]
	e2 := s.tsp.matrix[n3][n4]

	e3 := s.tsp.matrix[n1][n3]
	e4 := s.tsp.matrix[n2][n4]

	return e1+e2 > e3+e4
}

func (totalTimeA) get(s *Solution) int {
	traveled := 0
	for i := 0; i < len(s.route)-1; i++ {
		traveled += s.tsp.matrix[s.route[i]][s.route[i+1]]
	}
	return traveled
}

func (totalTimeA) isProfitable(s *Solution, i, j int, spans ...int) bool {
	n1 := s.route[i]
	n2 := s.route[i+1]
	n3 := s.route[j]
	n4 := 0

	if j < s.tsp.numNodes-1 {
		n4 = s.route[j+1]
	}

	e1 := s.tsp.matrix[n1][n2]
	e2 := s.tsp.matrix[n3][n4]

	e3 := s.tsp.matrix[n1][n3]
	e4 := s.tsp.matrix[n2][n4]

	// TODO handle asymetric

	return e1+e2 > e3+e4
}

// process strategy
func (local local2Opt) process(x *Solution) *Solution {
	level := 1
	iterLevel := 0

	local.strategy.process(x)

	for level < local.levelMax {
		x2 := x.disturb(level)

		local.strategy.process(x2)

		if local.objective.get(x2) < local.objective.get(x) {
			iterLevel = 0
			level = 1
			x = x2
		} else {
			if iterLevel > local.iterMax {
				level++
				iterLevel = 0
			}
		}
		iterLevel++
	}

	return x
}

// constrained 2 opt
type cons2Opt struct {
	traveled   []int
	precedence map[int]int
	carrying   map[int]int
	objective
}

func (c cons2Opt) process(s *Solution) {
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

func (c cons2Opt) exchangeGlobalUpdate(s *Solution, iaux, jaux int) {
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
func (c cons2Opt) isFeasible(s *Solution, i, j int) bool {

	sum := c.traveled[i]

	carrying := c.carrying[i]

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
		/* if (sum <= this->inc[i])
		   {
		   	break;
		   } */

		if carrying > s.tsp.capacity {
			return false
		}
	}

	if s.tsp.readyTime[s.route[len(s.route)-2]] > sum {
		sum = s.tsp.readyTime[s.route[len(s.route)-2]]
	}

	return true
}

func (c *cons2Opt) calcGlobals(s *Solution) {
	var n1, n2 int

	sum := s.tsp.traveled
	carrying := s.tsp.carrying

	c.precedence = make(map[int]int)

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

func (cons2Opt) isProfitable(s *Solution, i, j int, spans ...int) bool {
	var n1, n2 int

	sum := spans[1]
	n1 = s.route[i]
	n2 = s.route[j]

	if s.tsp.readyTime[n1] > sum {
		sum = s.tsp.readyTime[n1]
	}

	sum += s.tsp.matrix[n1][n2]

	for k := j; k > i+1; k-- {
		n1 = s.route[k]
		n2 = s.route[k-1]

		if s.tsp.readyTime[n1] > sum {
			sum = s.tsp.readyTime[n1]
		}

		sum += s.tsp.matrix[n1][n2]
	}

	n1 = s.route[i+1]
	n2 = s.route[j+1]

	if s.tsp.readyTime[n1] > sum {
		sum = s.tsp.readyTime[n1]
	}

	sum += s.tsp.matrix[n1][n2]

	return spans[0] > sum
}
