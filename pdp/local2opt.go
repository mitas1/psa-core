package pdp

import (
	"math/rand"
)

type Local2Opt struct {
	inc []int
}

func (local *Local2Opt) Process(x *Solution) *Solution {
	level := 1
	levelMax := 8
	iterMax := 5

	iterLevel := 0

	local.local2Opt(x)

	for level < levelMax {
		x2 := local.disturb(x, level)

		local.local2Opt(x2)

		if x2.TotalDistance() < x.TotalDistance() {
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
	start := iaux + 1 // start = n1 +1
	end := jaux       // end = n3

	prev := s.TotalDistance()

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

	if prev > x.TotalDistance() && x.IsFeasible() {
		return x, true
	}
	return s, false
}

func (local *Local2Opt) disturb(s *Solution, level int) *Solution {
	levelMax := level * 2
	var j int
	imp := false

	for j < levelMax {

		// take a random node
		n1 := rand.Int()%len(s.route) - 1 + 1

		for i := n1; i < len(s.route)-1; i++ {
			if i != n1 {
				s, imp = local.exchange(s, n1, i)
				if imp {
					break
				}
			}
		}

		for i := n1; i > 1; i-- {
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
