package pdptw

import (
	"log"
	"math/rand"

	"github.com/mitas1/psa-core/config"
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
				travelled:  make([]int, numNodes),
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

type objective interface {
	get(*Solution) int
	isProfitable(s *Solution, i, j int, spans ...int) bool
}

type spanTime struct{}
type totalTime struct{}
type totalTimeA struct{}

func (spanTime) get(s *Solution) int {
	traveled := 0
	for i := 0; i < len(s.route)-1; i++ {
		if traveled < s.tsp.readytime[s.route[i]] {
			traveled = s.tsp.readytime[s.route[i]]
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

	if s.tsp.readytime[n1] > sum {
		sum = s.tsp.readytime[n1]
	}

	sum += s.tsp.matrix[n1][n2]

	for k := j; k > i+1; k-- {
		n1 = s.route[k]
		n2 = s.route[k-1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]
	}

	n1 = s.route[i+1]
	n2 = s.route[j+1]

	if s.tsp.readytime[n1] > sum {
		sum = s.tsp.readytime[n1]
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
	travelled  []int
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
			if c.objective.isProfitable(s, i, j, c.travelled[j+1], c.travelled[i]) {
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

	sum := c.travelled[iaux]
	carrying := c.carrying[iaux-1]

	// update the reversed path
	for i := iaux; i < jaux; i++ {
		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]
		carrying += s.tsp.demands[n1]

		c.travelled[i+1] = sum
		c.carrying[i] = carrying
	}

	// update the rest
	for i := jaux; i < len(s.route)-1; i++ {
		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]
		carrying += s.tsp.demands[n1]

		c.travelled[i+1] = sum
		c.carrying[i] = carrying
	}

	return
}

// exchange (i,i+1), (j,j+1) ===> (i,j), (i+1,j+1)
func (c cons2Opt) isFeasible(s *Solution, i, j int) bool {

	sum := c.travelled[i]

	carrying := c.carrying[i]

	n1 := s.route[i]
	n2 := s.route[j]

	if s.tsp.readytime[n1] > sum {
		sum = s.tsp.readytime[n1]
	}

	sum += s.tsp.matrix[n1][n2]

	if sum > s.tsp.duedate[n2] {
		return false
	}

	if carrying > s.tsp.capacity {
		return false
	}

	for k := j; k > i+1; k-- {

		n1 = s.route[k]
		n2 = s.route[k-1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]

		if sum > s.tsp.duedate[n2] {
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

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]
		carrying += s.tsp.demands[n1]

		if sum > s.tsp.duedate[n2] {
			return false
		}

		if carrying > s.tsp.capacity {
			return false
		}
	}

	for k := j + 1; k < len(s.route)-1; k++ {
		n1 = s.route[k]
		n2 = s.route[k+1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]
		carrying += s.tsp.demands[n1]

		if sum > s.tsp.duedate[n2] {
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

	if s.tsp.readytime[s.route[len(s.route)-2]] > sum {
		sum = s.tsp.readytime[s.route[len(s.route)-2]]
	}

	return true
}

func (c *cons2Opt) calcGlobals(s *Solution) {
	var n1, n2 int

	sum := s.tsp.travelled
	carrying := s.tsp.carrying

	c.precedence = make(map[int]int)

	for i := 0; i < len(s.route)-1; i++ {
		// travelled
		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]

		c.travelled[i+1] = sum

		// precedence
		if n, ok := s.tsp.precedence[n1]; ok {
			index := indexOf(n, s.route)
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

	index := indexOf(n, s.route)
	c.precedence[i] = index
	c.precedence[index] = i

	return
}

func (cons2Opt) isProfitable(s *Solution, i, j int, spans ...int) bool {
	var n1, n2 int

	sum := spans[1]
	n1 = s.route[i]
	n2 = s.route[j]

	if s.tsp.readytime[n1] > sum {
		sum = s.tsp.readytime[n1]
	}

	sum += s.tsp.matrix[n1][n2]

	for k := j; k > i+1; k-- {
		n1 = s.route[k]
		n2 = s.route[k-1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]
	}

	n1 = s.route[i+1]
	n2 = s.route[j+1]

	if s.tsp.readytime[n1] > sum {
		sum = s.tsp.readytime[n1]
	}

	sum += s.tsp.matrix[n1][n2]

	return spans[0] > sum
}

type Local2Opt struct {
	inc []int
}

func (local Local2Opt) objective(x *Solution) int {
	return 0
}

func (local *Local2Opt) process(s *Solution) {
	var pos, npos, n1, n2, n3, n4, e1, e2, e3, e4 int
	improvement := false
	numNodes := s.tsp.numNodes

	// create auxiliary set
	pointer := numNodes - 2
	set := make([]int, pointer)

	for i := 0; i < pointer; i++ {
		set[i] = i + 1
	}

	// n1 -> n2 -> n3 -> n4
	// n1 -> n3 -> n2 -> n4
	for pointer > 0 {
		improvement = false

		pos = rand.Int() % pointer

		npos = set[pos]

		n1 = s.route[npos]
		n2 = s.route[npos+1]
		e1 = s.tsp.matrix[n1][n2]

		for i := npos + 2; i < numNodes; i++ {
			n3 = s.route[i]

			n4 = 0
			if i < numNodes-1 {
				n4 = s.route[i+1]
			}

			e2 = s.tsp.matrix[n3][n4]

			e3 = s.tsp.matrix[n1][n3]
			e4 = s.tsp.matrix[n2][n4]

			if e1+e2 > e3+e4 {
				s, improvement = local.exchange(s, npos, i)
				if improvement {
					break
				}
			}
		}

		if improvement {
			pointer = numNodes - 2
		} else {
			aux := set[pointer-1]
			set[pointer-1] = npos
			set[pos] = aux
			pointer--
		}
	}
	return
}

// npos, i
// n1, n3
func (*Local2Opt) exchange(s *Solution, iaux, jaux int) (*Solution, bool) {
	if iaux > jaux {
		return s, false
	}
	start := iaux + 1 // start = n1 +1
	end := jaux       // end = n3

	prev := s.MakeSpan()

	x := s.Copy()

	aux := 0
	j := 0

	// median
	f := (end-start+1)/2 + start - 1

	for i := start; i <= f; i++ {
		j = end - (i - start)
		aux = x.route[i]
		x.route[i] = x.route[j]
		x.route[j] = aux
	}

	if x.IsFeasible() && prev > x.MakeSpan() {
		return x, true
	}
	// log.Print("SKIP")
	return s, false
}

func (*Local2Opt) exchangeDisturb(s *Solution, iaux, jaux int) (*Solution, bool) {
	if iaux > jaux {
		return s, false
	}
	start := iaux + 1 // start = n1 +1
	end := jaux       // end = n3

	x := s.Copy()

	aux := 0
	j := 0

	// median
	f := (end-start+1)/2 + start - 1

	for i := start; i <= f; i++ {
		j = end - (i - start)
		aux = x.route[i]
		x.route[i] = x.route[j]
		x.route[j] = aux
	}

	if x.IsFeasible() {
		return x, true
	}
	return s, false
}

func (local *Local2Opt) disturb(s *Solution, level int) *Solution {
	levelMax := level
	var j int
	imp := false

	for j < levelMax {

		// take a random node
		n1 := rand.Int()%len(s.route) + 1

		for i := n1; i < len(s.route); i++ {
			if i != n1 {
				s, imp = local.exchange(s, n1, i)
				if imp {
					break
				}
			}
		}
		j++
	}
	return s
}

// exchange (i,i+1), (j,j+1) ===> (i,j), (i+1,j+1)
func (local *Local2Opt) isTWfeasible(s *Solution, i, j int, inc []int) bool {

	sum := inc[i]

	n1 := s.route[i]
	n2 := s.route[j]

	if s.tsp.readytime[n1] > sum {
		sum = s.tsp.readytime[n1]
	}

	sum += s.tsp.matrix[n1][n2]

	if sum > s.tsp.duedate[n2] {
		return false
	}

	for k := j; k > i+1; k-- {

		n1 = s.route[k]
		n2 = s.route[k-1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]

		if sum > s.tsp.duedate[n2] {
			return false
		}
	}

	if j+1 < len(s.route) {
		n1 = s.route[i+1]
		n2 = s.route[j+1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}
		sum += s.tsp.matrix[n1][n2]
		if sum > s.tsp.duedate[n2] {
			return false
		}
	}

	for k := j + 1; k < len(s.route)-1; k++ {
		n1 = s.route[k]
		n2 = s.route[k+1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}
		sum += s.tsp.matrix[n1][n2]
		if sum > s.tsp.duedate[n2] {
			return false
		}
		/* if (sum <= this->inc[i])
		   {
		   	break;
		   } */
	}

	return true
}

func (Local2Opt) isProfitable(s *Solution, i, j int) bool {
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

func (local *Local2Opt) local2OptLexical(s *Solution) {

	log.Print(s.IsFeasible())

	inc := local.calcGlobalTSP(s)

	i := 0

	traveled := 0

	mark := make(map[int]int)

	for x := 0; x < len(s.route); x++ {
		mark[x] = 0
		traveled += s.tsp.matrix[s.route[i]][s.route[i+1]]
	}

	tmp := false

	// outer loop
	for {

		if i >= len(s.route)-4 {
			break
		}

		// inner loop
		for j := i + 2; j < len(s.route)-1; j++ {

			// log.Printf("%v - %v", i+1, j+1)

			if !local.isTWfeasible(s, i, j, inc) {
				break
			}

			if local.isProfitable(s, i, j) {
				// log.Print(s.MakeSpan())
				log.Printf("EXCHANGE: %v - %v - %v", i, j, local.isTWfeasible(s, i, j, inc))
				local.clearExchange(s, i, j)
				inc = local.calcGlobalTSP(s)
				// s.Print()
				// log.Print(s.MakeSpan())
				break
			}

			// log.Printf("PROFIT: %v", local.isProfitable(s, traveled, i, j))

			feas := true

			if Contains(s.route[i+1:j], s.tsp.pred[s.route[j]]) {
				feas = false
			}

			for k := i + 1; k < j; k++ {
				if mark[s.route[k]] != 0 {
					feas = false
					log.Print(s.route[k])
					break
				}
			}

			if !feas {
				i++
				//log.Print("NOT FEAS")
				break
			} else {
				// log.Printf("%v - %v", i+1, j+1)
			}

			// j + 1 examination
			if Contains(s.route[i+1:j+1], s.tsp.pred[s.route[j+1]]) {
				mark[s.route[j+1]] = 1

				log.Print(s.tsp.pred[s.route[j+1]])

				if s.tsp.pred[s.route[j+1]] == s.route[j] {
					log.Print("pred(j+1) = j")

					for k := i + 1; k <= j; k++ {
						if s.tsp.pred[s.route[k]] < 0 {
							mark[-s.tsp.pred[s.route[k]]] = 0
						}
					}

					if s.tsp.pred[s.route[j+2]] < 0 {
						mark[-s.tsp.pred[s.route[j+2]]] = 1
					}

					i = j
					tmp = true
				} else {

					log.Print("pred(j+1) < j")

					for k := i + 1; k < s.tsp.pred[s.route[j+1]]; k++ {
						if s.tsp.pred[s.route[k]] < 0 {
							mark[-s.tsp.pred[s.route[k]]] = 0
						}
					}

					i = indexOf(s.tsp.pred[s.route[j+1]], s.route)
					tmp = true
				}
				break
			}

			if s.tsp.pred[s.route[j+1]] < 0 {
				// log.Printf("MARKING %v %v", s.route[j+1], s.tsp.pred[s.route[j+1]])
				mark[-s.tsp.pred[s.route[j+1]]] = 1
				// log.Printf("%v", mark)
			}

		}

		if !tmp {
			i++
		} else {
			tmp = false
		}

		if i >= len(s.route)-4 {
			break
		}
	}
}

func Contains(a []int, x int) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func indexOf(element int, data []int) int {
	for k, v := range data {
		if element == v {
			return k
		}
	}
	return -1 //not found.
}

func (*Local2Opt) clearExchange(s *Solution, iaux, jaux int) bool {
	if iaux > jaux {
		return false
	}
	start := iaux + 1 // start = n1 +1
	end := jaux       // end = n3

	aux := 0
	j := 0

	// median
	f := (end-start+1)/2 + start - 1

	for i := start; i <= f; i++ {
		j = end - (i - start)
		aux = s.route[i]
		s.route[i] = s.route[j]
		s.route[j] = aux
	}

	return true
}

func (*Local2Opt) calcGlobalTSP(s *Solution) []int {
	sum := 0
	n1 := 0
	n2 := 0

	inc := make([]int, len(s.route))

	for i := 0; i < len(s.route)-1; i++ {
		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]

		inc[i] = sum
	}
	return inc
}

func (*Local2Opt) calcGlobals(s *Solution) ([]int, map[int]int, map[int]int) {
	sum := 0
	carrying := 0
	n1 := 0
	n2 := 0

	travelled := make([]int, len(s.route))
	precendence := make(map[int]int)
	capacity := make(map[int]int)

	for i := 0; i < len(s.route)-1; i++ {
		// travelled
		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]

		travelled[i] = sum

		// precedence

		n3 := s.tsp.pred[n1]
		if n3 < 0 {
			n3 = -n3
		}

		precendence[i] = indexOf(n3, s.route)

		// capacity

		carrying += s.tsp.demands[n1]

		capacity[i] = carrying
	}

	precendence[len(s.route)-1] = indexOf(s.tsp.pred[s.route[len(s.route)-1]], s.route)

	log.Print(precendence)

	return travelled, precendence, capacity
}
