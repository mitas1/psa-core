package core

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/mitas1/psa-core/logging"
	"github.com/mitas1/psa-core/utils"
)

type SetType int32

const (
	FEASIBLE_SET   SetType = 0
	UNFEASIBLE_SET SetType = 1
)

// Solution struct
type Solution struct {
	route []int
	nodes map[int]bool
	tsp   *PDPTW
}

// NewSolution returns a new instance of the Solution struct
func NewSolution(tsp *PDPTW, route []int) Solution {
	nodes := make(map[int]bool)
	s := Solution{route: route, nodes: nodes, tsp: tsp}
	for node := range route {
		s.nodes[node] = true
	}
	return s
}

// GetRoute returns array of calculated route
func (s *Solution) GetRoute() []int {
	return s.route
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

func (s *Solution) WriteToFile(dir, name string) {
	os.MkdirAll(dir, os.ModePerm)

	filePath := path.Join(dir, name)

	err := ioutil.WriteFile(filePath, []byte(strings.Join(s.strings(), " ")), 0666)
	if err != nil {
		panic(err)
	}
}

// IsFeasible checks if solution is feasible
func (s *Solution) IsFeasible() bool {
	traveled := s.tsp.traveled
	carrying := s.tsp.carrying

	for i := 1; i < s.tsp.numNodes; i++ {
		traveled += s.tsp.matrix[s.route[i-1]][s.route[i]]
		carrying += s.tsp.demands[s.route[i-1]]

		// wait to ready to time
		if traveled < s.tsp.readyTime[s.route[i]] {
			traveled = s.tsp.readyTime[s.route[i]]
		}

		if carrying > s.tsp.capacity {
			return false
		}

		if value, ok := s.tsp.precedence[s.route[i]]; ok {
			if i <= utils.IndexOf(value, s.route) {
				return false
			}
		}

		if s.tsp.dueDate[s.route[i]] != 0 && s.tsp.dueDate[s.route[i]] < traveled {
			return false
		}
	}

	i := len(s.route) - 1

	if carrying+s.tsp.demands[s.route[i]] != 0 {
		return false
	}

	return true
}

func (s *Solution) isFeasibleEdge(i, j int, sum, carrying *int) bool {
	n1 := s.route[i]
	n2 := s.route[j]

	if s.tsp.readyTime[n1] > *sum {
		*sum = s.tsp.readyTime[n1]
	}

	*sum += s.tsp.matrix[n1][n2]
	*carrying += s.tsp.demands[n2]

	if *sum > s.tsp.dueDate[n2] || *carrying > s.tsp.capacity {
		return false
	}
	return true
}

func (s *Solution) isFeasibleRange(start, end int, sum, carrying *int) bool {
	var n1, n2 int
	for i := start; i < end; i++ {
		n1 = s.route[i]
		n2 = s.route[i+1]
		if s.tsp.readyTime[n1] > *sum {
			*sum = s.tsp.readyTime[n1]
		}
		*sum += s.tsp.matrix[n1][n2]
		*carrying += s.tsp.demands[n2]
		if *sum > s.tsp.dueDate[n2] || *carrying > s.tsp.capacity {
			return false
		}
	}
	return true
}

func (s *Solution) getSet(setType SetType) (set []int) {
	traveled := s.tsp.traveled
	carrying := s.tsp.carrying
	predViolation := false

	for i := 1; i < len(s.route); i++ {
		traveled += s.tsp.matrix[s.route[i-1]][s.route[i]]
		carrying += s.tsp.demands[s.route[i-1]]
		predViolation = false

		// wait to ready to time
		if traveled < s.tsp.readyTime[s.route[i]] {
			traveled = s.tsp.readyTime[s.route[i]]
		}

		if value, ok := s.tsp.precedence[s.route[i]]; ok {
			if i >= utils.IndexOf(value, s.route) {
				predViolation = true
			}
		}

		isFeasible := !predViolation || (s.tsp.dueDate[s.route[i]] != 0 &&
			s.tsp.dueDate[s.route[i]] < traveled) || carrying > s.tsp.capacity

		if setType == FEASIBLE_SET && isFeasible {
			set = append(set, i)
		} else if isFeasible {
			set = append(set, i)
		}
	}

	return
}

// exchange move node in position pos to new position newPos
func (s *Solution) exchange(pos, newPos int) {
	node := s.route[pos]
	if newPos > pos {
		for i := pos; i < newPos; i++ {
			s.route[i] = s.route[i+1]
		}
	} else {
		for i := pos; i > newPos; i-- {
			s.route[i] = s.route[i-1]
		}
	}
	s.route[newPos] = node
}

func (s *Solution) change(i, j int) {
	s.route[i], s.route[j] = s.route[j], s.route[i]
}

func (s *Solution) hasNode(node int) bool {
	if val, _ := s.nodes[node]; val {
		return true
	}
	return false
}

func (s *Solution) addNode(node int) {
	s.route = append(s.route, node)
	s.nodes[node] = true
}

func (s *Solution) removeNode(node int) {
	s.route = s.route[:len(s.route)-1]
	s.nodes[node] = false
}

func (s *Solution) getNode(index int) int {
	return s.route[index]
}

func (s *Solution) getCurrent() int {
	return s.route[len(s.route)-1]
}

func (s *Solution) TotalDistance() int {
	total := 0
	for i := 1; i <= len(s.route)-1; i++ {
		total += s.tsp.matrix[s.route[i-1]][s.route[i]]
	}
	return total
}

func (s *Solution) MakeSpan() int {
	traveled := s.tsp.traveled
	for i := 0; i < len(s.route)-1; i++ {
		if traveled < s.tsp.readyTime[s.route[i]] {
			traveled = s.tsp.readyTime[s.route[i]]
		}
		traveled += s.tsp.matrix[s.route[i]][s.route[i+1]]
	}
	return traveled
}

// Copy make a copy of Solution
func (s Solution) Copy() *Solution {
	route := []int{}

	for i := 0; i < s.tsp.numNodes; i++ {
		route = append(route, s.route[i])
	}

	return &Solution{route: route, tsp: s.tsp}
}

func (s Solution) disturb(level int) (x *Solution) {
	x = s.Copy()

	for j := 0; j < level; j++ {
		n1 := utils.Random(0, len(s.route)-2)
		n2 := utils.Random(n1+2, len(s.route))
		x = x.kExchange(n1, n2)
	}
	return x
}

func (s *Solution) kExchange(iaux, jaux int) *Solution {
	if iaux > jaux {
		return s
	}
	start := iaux + 1 // start = n1 +1
	end := jaux       // end = n3

	j := 0
	x := s.Copy()

	// median
	median := (end-start+1)/2 + start - 1

	for i := start; i <= median; i++ {
		j = end - (i - start)
		x.route[i], x.route[j] = x.route[j], x.route[i]
	}

	if x.IsFeasible() {
		return x
	}
	return s
}

func (s *Solution) calcGlobals() (traveled []int, carrying, precedence map[int]int) {
	var n1, n2 int

	_traveled := s.tsp.traveled
	_carrying := s.tsp.carrying

	traveled = make([]int, s.tsp.numNodes)
	carrying = make(map[int]int)
	precedence = make(map[int]int)

	for i := 0; i < len(s.route)-1; i++ {
		// traveled
		n1 = s.route[i]
		n2 = s.route[i+1]

		if s.tsp.readyTime[n1] > _traveled {
			_traveled = s.tsp.readyTime[n1]
		}

		_traveled += s.tsp.matrix[n1][n2]

		traveled[i+1] = _traveled

		// precedence
		if n, ok := s.tsp.precedence[n1]; ok {
			index := utils.IndexOf(n, s.route)
			precedence[index] = i
			precedence[i] = index
		} else if _, ok := precedence[i]; !ok {
			// ignore precedence of vertex
			precedence[i] = -1
		}

		_carrying += s.tsp.demands[n1]

		carrying[i] = _carrying
	}

	i := len(s.route) - 1

	n := s.tsp.precedence[s.route[i]]

	index := utils.IndexOf(n, s.route)
	precedence[i] = index
	precedence[index] = i

	return
}

//// ------- TESTING -----------------------------------------------------------

var (
	log = logging.GetLogger()
)

func (s *Solution) Check() bool {

	if s.tsp.startNode != s.route[0] {
		log.Errorf("%v: %v", "Wrong startnode!", s.route)
		return false
	}

	if s.tsp.numNodes != len(s.route) {
		log.Errorf("%v: %v", "numNodes are not equal to route", s.route)
		return false
	}

	set := make(map[int]bool)

	for _, v := range s.route {
		set[v] = true
	}

	if len(set) != len(s.route) {
		log.Errorf("%v: %v", "Some nodes are duplicated", s.route)
		return false
	}

	if !s.IsFeasible() {
		log.Errorf("%v: %v", "Solution is not FEASIBLE!", s.route)
		return false
	}

	return true
}
