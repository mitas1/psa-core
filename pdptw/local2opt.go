package pdptw

import (
	"log"
	"math/rand"
	"reflect"

	"github.com/mitas1/psa-core/config"
	"github.com/mitas1/psa-core/utils"
)

// interface of local 2opt search
type local2Opt interface {
	Process(*Solution) *Solution
}

type local2OptBase struct{}

// TODO - Rewrite needed
func (local local2OptBase) disturb(s *Solution, level int) *Solution {

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

func (local2OptBase) exchange(s *Solution, iaux, jaux int) (*Solution, bool) {
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
	return s, false
}

func NewOptimization(opts config.Optimization, numNodes int) (opt local2Opt) {
	switch opts.Strategy {
	case "2opt":
		opt = Local2Opt{inc: make([]int, numNodes)}
	case "cons2opt":
		opt = cons2Opt{inc: make([]int, numNodes)}
	default:
		opt = cons2Opt{inc: make([]int, numNodes)}
	}

	return opt
}

// constrained 2 opt
type cons2Opt struct {
	local2OptBase
	inc []int
}

func (local cons2Opt) Process(x *Solution) *Solution {
	level := 1
	levelMax := 20
	iterMax := 2

	iterLevel := 0

	local.local2Opt(x)

	for level < levelMax {
		x2 := local.disturb(x, level)

		local.local2Opt(x2)

		if x2.MakeSpan() < x.MakeSpan() {
			iterLevel = 0
			level = 1
			x = x2
		} else {
			if iterLevel > iterMax {
				level++
				iterLevel = 0
			}
		}
		iterLevel++
	}

	return x
}

func (local cons2Opt) local2Opt(s *Solution) {
	var pos, npos, n1, n2, n3, n4, e1, e2, e3, e4 int
	improvement := false
	numNodes := s.tsp.numNodes

	// create auxiliary set
	setSize := numNodes - 2
	set := make([]int, setSize)

	for i := 0; i < setSize; i++ {
		set[i] = i + 1
	}

	inc, pred, capacity := local.calcGlobals(s)

	// n1 -> n2 -> n3 -> n4
	// n1 -> n3 -> n2 -> n4
	for setSize > 0 {
		improvement = false

		pos = rand.Int() % setSize

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
				if local.isFeasible(s, npos, i, inc, pred, capacity) {
					inc, pred = local.exchangeGlobalUpdate(s, npos, i, inc, pred, capacity)

					inc2, pred2, _ := local.calcGlobals(s)

					if !utils.Equal(inc2, inc) {
						log.Printf("\n%v - %v\n%v\n%v\nFUCKING\n\n", npos, i, inc, inc2)
					}

					if !reflect.DeepEqual(pred, pred2) {
						s.Print()
						log.Printf("\n%v - %v\n%v\n%v\nFUCKING\n\n", npos, i, pred, pred2)
					}

					improvement = true
					break
				}
			}
		}

		if improvement {
			setSize = numNodes - 2
		} else {
			aux := set[setSize-1]
			set[setSize-1] = npos
			set[pos] = aux
			setSize--
		}
	}
	return
}

func (cons2Opt) exchangeGlobalUpdate(
	s *Solution, iaux, jaux int, inc []int, precedence map[int]int, capacity map[int]int,
) ([]int, map[int]int) {
	if iaux > jaux {
		return nil, nil
	}

	start := iaux + 1 // start = n1 +1
	end := jaux       // end = n3

	aux := 0
	j := 0

	var pred_i, pred_j int

	// median
	f := (end-start+1)/2 + start - 1

	// log.Printf("BEFORE: %v, %v \n%v\n", iaux, jaux, precedence)

	for i := start; i <= f; i++ {

		j = end - (i - start)

		// log.Printf("BEFORE: %v, %v ", i, j)

		aux = s.route[i]
		s.route[i] = s.route[j]
		s.route[j] = aux

		// update precedence
		pred_i = precedence[i]
		pred_j = precedence[j]

		precedence[pred_i] = j
		precedence[pred_j] = i

		precedence[j] = pred_i
		precedence[i] = pred_j
	}

	// log.Printf("AFTER: \n%v\n\n\n\n", precedence)

	var n1, n2 int

	sum := inc[iaux-1]

	// update the reversed path

	for i := iaux; i < jaux; i++ {

		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]

		inc[i] = sum
	}

	// update the rest

	for i := jaux; i < len(s.route)-1; i++ {
		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readytime[n1] > sum {
			sum = s.tsp.readytime[n1]
		}

		sum += s.tsp.matrix[n1][n2]

		inc[i] = sum
	}

	return inc, precedence
}

// exchange (i,i+1), (j,j+1) ===> (i,j), (i+1,j+1)
func (cons2Opt) isFeasible(
	s *Solution, i, j int, inc []int, precedence map[int]int, capacity map[int]int,
) bool {

	sum := inc[i]

	carrying := capacity[i]

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

		if precedence[k] > i && precedence[k] < j {
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

	return true
}

func (cons2Opt) calcGlobals(s *Solution) ([]int, map[int]int, map[int]int) {
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

	return travelled, precendence, capacity
}

type Local2Opt struct {
	inc []int
}

func (local Local2Opt) Process(x *Solution) *Solution {
	level := 1
	levelMax := 20
	iterMax := 2

	iterLevel := 0

	local.local2Opt(x)

	for level < levelMax {
		x2 := local.disturb(x, level)

		local.local2Opt(x2)

		if x2.MakeSpan() < x.MakeSpan() {
			iterLevel = 0
			level = 1
			x = x2
		} else {
			if iterLevel > iterMax {
				level++
				iterLevel = 0
			}
		}
		iterLevel++
	}

	return x
}

func (local *Local2Opt) local2Opt(s *Solution) {
	var pos, npos, n1, n2, n3, n4, e1, e2, e3, e4 int
	improvement := false
	numNodes := s.tsp.numNodes

	// create auxiliary set
	setSize := numNodes - 2
	set := make([]int, setSize)

	for i := 0; i < setSize; i++ {
		set[i] = i + 1
	}

	// n1 -> n2 -> n3 -> n4
	// n1 -> n3 -> n2 -> n4
	for setSize > 0 {
		improvement = false

		pos = rand.Int() % setSize

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
			setSize = numNodes - 2
		} else {
			aux := set[setSize-1]
			set[setSize-1] = npos
			set[pos] = aux
			setSize--
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

func (local *Local2Opt) calcProfit(s *Solution, travelled, i, j int) int {
	n1 := s.route[i]
	n2 := s.route[i+1]
	e1 := s.tsp.matrix[n1][n2]

	n3 := s.route[j]
	n4 := s.route[j+1]

	e2 := s.tsp.matrix[n3][n4]

	e3 := s.tsp.matrix[n1][n3]
	e4 := s.tsp.matrix[n2][n4]

	return (e3 + e4) - (e1 + e2)
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

			if local.calcProfit(s, inc[i], i, j) > 0 {
				// log.Print(s.MakeSpan())
				log.Printf("EXCHANGE: %v - %v - %v", i, j, local.isTWfeasible(s, i, j, inc))
				local.clearExchange(s, i, j)
				inc = local.calcGlobalTSP(s)
				// s.Print()
				// log.Print(s.MakeSpan())
				break
			}

			// log.Printf("PROFIT: %v", local.calcProfit(s, traveled, i, j))

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

	s.Print()

	log.Printf("IS FEASIBLE: %v\nIS FEASIBLE PRECEDENCE: %v\n\n", s.IsFeasible(), s.IsFeasiblePrecendence())
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

	return travelled, precendence, capacity
}
