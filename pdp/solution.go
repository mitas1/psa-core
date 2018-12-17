package pdp

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"sort"
	"strings"

	"github.com/mitas1/psa-core/utils"
)

// Solution struct
type Solution struct {
	route []int
	nodes map[int]bool
	tsp   *PDP
}

// NewSolution returns a new instance of the Solution struct
func NewSolution(tsp *PDP, route []int) Solution {
	nodes := make(map[int]bool)
	s := Solution{route, nodes, tsp}
	for node := range route {
		s.nodes[node] = true
	}
	return s
}

// Copy make a copy of Solution
func (s Solution) Copy() *Solution {
	route := []int{}

	for i := 0; i < s.tsp.numNodes; i++ {
		route = append(route, s.route[i])
	}

	return &Solution{route: route, tsp: s.tsp}
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

func (s *Solution) WriteToFile(dir string) {
	os.MkdirAll(dir, os.ModePerm)

	filePath := path.Join(dir, s.tsp.name)

	err := ioutil.WriteFile(filePath, []byte(strings.Join(s.strings(), " ")), 0666)
	if err != nil {
		panic(err)
	}
}

// IsFeasible checks if solution is feasible
func (s *Solution) IsFeasible() bool {
	hasNode := false

	// do not include +0 and -0
	for i := 1; i < s.tsp.numNodes-1; i++ {
		hasNode = false

		if value, ok := s.tsp.precendense[s.route[i]]; ok {
			for j := i; j > 0; j-- {
				if value == s.route[j] {
					hasNode = true
					break
				}
			}
			if !hasNode {
				return false
			}
		}
	}

	return true
}

// Penalty TODO rewrite using hashing
func (s *Solution) Penalty() (penalty int) {
	hasNode := false

	// do not include +0 and -0
	for i := 1; i < s.tsp.numNodes-1; i++ {
		hasNode = false

		if value, ok := s.tsp.precendense[s.route[i]]; ok {
			for j := i; j > 0; j-- {
				if value == s.route[j] {
					hasNode = true
					break
				}
			}
			if !hasNode {
				for j := i; j < s.tsp.numNodes; j++ {
					if value == s.route[j] {
						penalty += j
						break
					}
				}
			}
		}
	}
	return penalty
}

func (s *Solution) calcSets() (feasible []int, unfeasible []int) {
	hasNode := false

	// do not include +0 and -0
	for i := 1; i < s.tsp.numNodes-1; i++ {
		hasNode = false

		if val, ok := s.tsp.precendense[s.route[i]]; ok {
			for j := i; j >= 0; j-- {
				if val == s.route[j] {
					hasNode = true
					break
				}
			}
			if hasNode {
				feasible = append(feasible, i)
			} else {
				unfeasible = append(feasible, i)
			}
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
	s.route = s.route[:s.tsp.numNodes-1]
	s.nodes[node] = false
}

func (s *Solution) GetNode(index int) int {
	return s.route[index]
}

func (s *Solution) GetCurrent() int {
	return s.route[s.tsp.numNodes-1]
}

// TotalDistance returns total distance traveled
func (s Solution) TotalDistance() int {
	total := 0
	for i := 1; i < s.tsp.numNodes; i++ {
		total += s.tsp.matrix[s.route[i-1]][s.route[i]]
	}
	return total
}

// Incompatible arcs, eliminate the search space
func (tsp *PDP) preprocess() {
	tsp.arcs = make(map[int]map[int]bool)
	for i := range tsp.matrix {
		tsp.arcs[i] = make(map[int]bool)

		for j := range tsp.matrix {
			tsp.arcs[i][j] = true
		}

		if value, ok := tsp.precendense[i]; ok {
			tsp.arcs[value][i] = false
		}
	}
}

func (s *Solution) NodePenalty(i int) (penalty int) {
	if value, ok := s.tsp.precendense[s.route[i]]; ok {
		for j := i; j < s.tsp.numNodes; j++ {
			if value == s.route[j] {
				penalty += j
				break
			}
		}
	}
	return
}

type Node struct {
	penalty int
	node    int
}

type ByPenalty []Node

func (a ByPenalty) Len() int           { return len(a) }
func (a ByPenalty) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPenalty) Less(i, j int) bool { return a[i].penalty < a[j].penalty }

func (s *Solution) PenaltySort() {
	penalties := []Node{}

	for node := range s.route {
		penalties = append(penalties, Node{penalty: s.NodePenalty(node), node: node})
	}
	sort.Sort(ByPenalty(penalties))
	for i, pen := range penalties {
		s.route[i] = pen.node
	}
	return
}
