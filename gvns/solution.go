package gvns

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mitas1/psa-core/utils"
)

// Solution struct
type Solution struct {
	route []int
	nodes map[int]bool
	tsp   *TSP
}

// NewSolution returns a new instance of the Solution struct
func NewSolution(tsp *TSP, route []int) Solution {
	nodes := make(map[int]bool)
	s := Solution{route, nodes, tsp}
	for node := range route {
		s.nodes[node] = true
	}
	return s
}

func (s *Solution) strings() []string {
	return utils.MapIntToStr(s.route, func(x int) string {
		return fmt.Sprintf("%d", x)
	})
}

// Print the solution
func (s *Solution) Print() {
	fmt.Print(strings.Join(s.strings(), " -> "))
}

func (s *Solution) WriteToFile() {
	dir := "../_solutions"

	os.MkdirAll(dir, os.ModePerm)

	filePath := path.Join(dir, s.tsp.name)

	err := ioutil.WriteFile(filePath, []byte(strings.Join(s.strings()[1:len(s.route)-1], " ")), 0666)
	if err != nil {
		panic(err)
	}
}

// IsFeasible checks if solution is feasible
func (s *Solution) IsFeasible() bool {
	traveled := 0
	for i := 1; i < len(s.route); i++ {
		traveled += s.tsp.matrix[s.route[i-1]][s.route[i]]

		// wait to ready to time
		if traveled < s.tsp.readytime[s.route[i]] {
			traveled = s.tsp.readytime[s.route[i]]
		}

		if s.tsp.duedate[s.route[i]] < traveled {
			return false
		}
	}
	return true
}

// Penalty is sum of all differences between the time to reach each customer
// and its due date
func (s *Solution) Penalty() (penalty int) {
	traveled := 0
	for i := 1; i < len(s.route); i++ {
		traveled += s.tsp.matrix[s.route[i-1]][s.route[i]]

		// wait to ready to time
		if traveled < s.tsp.readytime[s.route[i]] {
			traveled = s.tsp.readytime[s.route[i]]
		}

		if s.tsp.duedate[s.route[i]] < traveled {
			penalty = penalty + traveled - s.tsp.duedate[s.route[i]]
		}
	}
	return penalty
}

func (s *Solution) calcSets() (feasible []int, unfeasible []int) {
	traveled := 0

	for i := 1; i <= len(s.route)-2; i++ {
		traveled += s.tsp.matrix[s.route[i-1]][s.route[i]]

		// wait to ready to time
		if traveled < s.tsp.readytime[s.route[i]] {
			traveled = s.tsp.readytime[s.route[i]]
		}

		if s.tsp.duedate[s.route[i]] < traveled ||
			s.tsp.readytime[s.route[i]] > traveled {
			unfeasible = append(unfeasible, i)
		} else {
			feasible = append(feasible, i)
		}
	}

	return feasible, unfeasible
}

func (s *Solution) exchange(p, paux int) {
	n := s.route[p]
	if paux > p {
		for i := p; i < paux; i++ {
			s.route[i] = s.route[i+1]
		}
	} else {
		for i := p; i > paux; i-- {
			s.route[i] = s.route[i-1]
		}
	}
	s.route[paux] = n
}

func (s *Solution) HasNode(node int) bool {
	if val, _ := s.nodes[node]; val {
		return true
	}
	return false
}

func (s *Solution) AddNode(node int) {
	s.route = append(s.route, node)
	s.nodes[node] = true
}

func (s *Solution) RemoveNode(node int) {
	s.route = s.route[:len(s.route)-1]
	s.nodes[node] = false
}

func (s *Solution) GetNode(index int) int {
	return s.route[index]
}

func (s *Solution) GetCurrent() int {
	return s.route[len(s.route)-1]
}

func (s *Solution) TotalDistance() int {
	total := 0
	for i := 1; i <= len(s.route)-1; i++ {
		total += s.tsp.matrix[s.route[i-1]][s.route[i]]
	}
	return total
}
