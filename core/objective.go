package core

import (
	"github.com/mitas1/psa-core/config"
)

type objective interface {
	get(*Solution) int
	isProfitable(s *Solution, i, j int, spans ...int) bool
}

type spanTime struct{}
type totalTime struct{}
type totalTimeA struct{}

func NewObjective(opts config.Optimization) objective {
	switch {
	case "time" == opts.Objective && opts.Asymetric:
		return totalTimeA{}
	case "time" == opts.Objective:
		return totalTime{}
	case "span" == opts.Objective:
		return spanTime{}
	default:
		return spanTime{}
	}
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
