package pdptw

import (
	"math"
	"math/rand"
)

type SA struct {
}

func (local *SA) Process(state *Solution) *Solution {

	iterMax := 30.0
	iter := 0.0
	T := 0.0
	fraction := 0.0
	new_cost := 0.0

	cost := float64(state.MakeSpan())
	var new_state *Solution

	for iter < iterMax {
		iter++

		fraction = iter / iterMax
		T = local.temperature(fraction)
		// log.Printf("temperature: %v", T)

		// log.Printf("Fraction - %v", int(fraction*100))
		new_state = local.disturb(state, int(fraction*100))
		//disturb := new_state.MakeSpan()
		new_state = local.local2Opt(new_state)

		//log.Printf("%v - %v", disturb, new_state.MakeSpan())

		new_cost = float64(new_state.MakeSpan())

		if local.probability(cost, new_cost, T) > rand.Float64() {
			//log.Printf("%v - %v - update", cost, new_cost)
			state = new_state
			cost = new_cost
		}
	}

	return state
}

func (local *SA) temperature(fraction float64) float64 {
	return math.Max(0.01, math.Min(1, 1-fraction))
}

func (local *SA) probability(cost, new_cost, T float64) float64 {
	if new_cost < cost {
		return 1
	}

	return math.Exp(-(new_cost - cost) / T)
}

func (local *SA) local2Opt(s *Solution) *Solution {
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
	return s
}

// npos, i
// n1, n3
func (*SA) exchange(s *Solution, iaux, jaux int) (*Solution, bool) {
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
		s = x
		return s, true
	}
	return s, false
}

func (*SA) exchangeDisturb(s *Solution, iaux, jaux int) (*Solution, bool) {
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
		s = x
		return s, true
	}
	return s, false
}

func (local *SA) disturb(s *Solution, level int) *Solution {
	levelMax := level * 2
	var j int
	//imp := false

	for j < levelMax {

		// take a random node
		n1 := rand.Int()%(len(s.route)-1) + 1

		for i := n1; i < len(s.route); i++ {
			if i != n1 {
				s, _ = local.exchangeDisturb(s, n1, i)
			}
		}

		for i := n1; i > 0; i-- {
			if i != n1 {
				s, _ = local.exchangeDisturb(s, n1, i)
			}
		}
		j++

	}

	return s
}
