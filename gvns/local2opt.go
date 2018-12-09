package gvns

import (
	"math/rand"
)

type Local2Opt struct {
	inc []int
}

func (l *Local2Opt) process(s *Solution) {
	var pos, npos, n1, n2, n3, n4, e1, e2, e3, e4 int
	improvement := false
	numNodes := s.tsp.numNodes

	l.calcInc(s)

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

			// check
			if !s.tsp.arcs[n3][n2] {
				break
			}

			n4 = 0
			if i < numNodes-1 {
				n4 = s.route[i+1]
			}

			e2 = s.tsp.matrix[n3][n4]

			e3 = s.tsp.matrix[n1][n3]
			e4 = s.tsp.matrix[n2][n4]

			if e1+e2 > e3+e4 {
				switch l.isFeasible(s, npos, i) {
				case 0:
					l.exchange(s, npos, i)
					improvement = true
					break
				case -2:
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
func (l *Local2Opt) exchange(s *Solution, iaux, jaux int) {
	start := iaux + 1 // start = n1 +1
	end := jaux       // end = n3

	aux := 0
	j := 0
	f := (end-start+1)/2 + start - 1

	for i := start; i <= f; i++ {
		j = end - (i - start)
		aux = s.route[i]
		s.route[i] = s.route[j]
		s.route[j] = aux
	}

	// calcInc

	numNodes := s.tsp.numNodes

	min := iaux
	max := jaux

	sum := 0

	if min > 1 {
		sum = l.inc[min-2]
	}

	ni := 0
	next := 0

	for i := min - 1; i < max; i++ {
		ni = s.route[i]
		next = s.route[i+1]

		if s.tsp.readytime[ni] > sum {
			sum = s.tsp.readytime[ni]
		}
		sum += s.tsp.matrix[ni][next]
		l.inc[i] = sum
	}

	for i := max; i < numNodes-1; i++ {
		ni = s.route[i]
		next = s.route[i+1]

		if s.tsp.readytime[ni] > sum {
			sum = s.tsp.readytime[ni]
		}
		sum += s.tsp.matrix[ni][next]
		if l.inc[i] == sum {
			break
		}
		l.inc[i] = sum
	}
}

// npos -> j
func (l *Local2Opt) isFeasible(s *Solution, iaux, jaux int) int {
	// ni -> next
	numNodes := s.tsp.numNodes
	sum := l.inc[iaux-1]

	ni := s.route[iaux]
	next := s.route[jaux]

	if s.tsp.readytime[ni] > sum {
		sum = s.tsp.readytime[ni]
	}

	sum += s.tsp.matrix[ni][next]
	if sum > s.tsp.duedate[next] {
		return -1
	}

	// for i = n2  i > n1 + 1 --
	for i := jaux; i > iaux+1; i-- {
		ni = s.route[i]
		next = s.route[i-1]

		if s.tsp.readytime[ni] > sum {
			sum = s.tsp.readytime[ni]
		}
		sum += s.tsp.matrix[ni][next]
		if sum > s.tsp.duedate[next] {
			return -2
		}
	}

	if jaux+1 < numNodes {
		ni = s.route[iaux+1]
		next = s.route[jaux+1]

		if s.tsp.readytime[ni] > sum {
			sum = s.tsp.readytime[ni]
		}
		sum += s.tsp.matrix[ni][next]

		if sum > s.tsp.duedate[next] {
			return -3
		}
	}

	for i := jaux + 1; i < numNodes-1; i++ {
		ni = s.route[i]
		next = s.route[i+1]

		if s.tsp.readytime[ni] > sum {
			sum = s.tsp.readytime[ni]
		}

		sum += s.tsp.matrix[ni][next]
		if sum > s.tsp.duedate[next] {
			return -4
		}
		if sum <= l.inc[i] {
			break
		}
	}

	return 0
}

func (l *Local2Opt) calcInc(s *Solution) {
	ni := 0
	next := 0
	sum := 0

	for i := 0; i < s.tsp.numNodes-1; i++ {
		ni = s.route[i]
		next = s.route[i+1]
		if s.tsp.readytime[ni] > sum {
			sum = s.tsp.readytime[ni]
		}
		sum += s.tsp.matrix[ni][next]
		l.inc[i] = sum
	}
}
