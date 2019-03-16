package pdptw

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

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
	hasNode := false
	traveled := 0
	carrying := 0

	for i := 1; i < s.tsp.numNodes; i++ {
		hasNode = false
		traveled += s.tsp.matrix[s.route[i-1]][s.route[i]]
		carrying += s.tsp.demands[s.route[i-1]]

		//log.Printf("%v - %v - %v", s.route[i-1], s.tsp.demands[s.route[i-1]], carrying)

		// wait to ready to time
		if traveled < s.tsp.readytime[s.route[i]] {
			traveled = s.tsp.readytime[s.route[i]]
		}

		if carrying > s.tsp.capacity {
			// log.Print("CAPACITY OVERFLOW")
			return false
		}

		if value, ok := s.tsp.precendense[s.route[i]]; ok {
			for j := i; j > 0; j-- {
				if value == s.route[j] {
					hasNode = true
					break
				}
			}
			if !hasNode {
				// log.Printf("%v - %v", s.route[i], i)
				// log.Print("PRECEDENCE OVERFLOW")
				return false
			}
		}

		if s.tsp.duedate[s.route[i]] != 0 && s.tsp.duedate[s.route[i]] < traveled {
			// log.Printf("%v TIME WINDOW OVERFLOW", i)
			return false
		}
	}

	// log.Printf("%v - %v - %v", carrying)
	return true
}

func (s *Solution) getSet(setType SetType) (set []int) {
	traveled := 0
	carrying := 0
	hasNode := false
	tmp := false

	for i := 1; i < len(s.route); i++ {
		traveled += s.tsp.matrix[s.route[i-1]][s.route[i]]
		carrying += s.tsp.demands[s.route[i-1]]
		hasNode = false
		tmp = false

		// wait to ready to time
		if traveled < s.tsp.readytime[s.route[i]] {
			traveled = s.tsp.readytime[s.route[i]]
		}

		if val, ok := s.tsp.pred[s.route[i]]; ok {
			for j := i; j < len(s.route); j++ {
				if -val == s.route[j] {
					hasNode = true
					break
				}
			}
			/* for j := i; j > 0; j-- {
				if val == s.route[j] {
					hasNode = true
					break
				}
			} */
			if !hasNode {
				tmp = true
			}
		}

		isFeasible := tmp || (s.tsp.duedate[s.route[i]] != 0 && s.tsp.duedate[s.route[i]] < traveled) || carrying > s.tsp.capacity

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

func (s *Solution) MakeSpan() int {
	traveled := 0
	for i := 0; i < len(s.route)-1; i++ {
		if traveled < s.tsp.readytime[s.route[i]] {
			traveled = s.tsp.readytime[s.route[i]]
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
	for j := 0; j < level; j++ {

		n1 := utils.Random(0, len(s.route)-2)
		n2 := utils.Random(n1+2, len(s.route))

		x = s.kExchange(n1, n2)
	}
	return x
}

func (s *Solution) kExchange(iaux, jaux int) *Solution {
	if iaux > jaux {
		return s
	}
	start := iaux + 1 // start = n1 +1
	end := jaux       // end = n3

	x := s.Copy()
	j := 0

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

//// ------- TESTING -----------------------------------------------------------

// IsFeasible checks if solution is feasible
func (s *Solution) IsFeasiblePrecendence() bool {
	hasNode := false
	traveled := 0
	carrying := 0

	for i := 1; i < s.tsp.numNodes; i++ {
		hasNode = false
		traveled += s.tsp.matrix[s.route[i-1]][s.route[i]]
		carrying += s.tsp.demands[s.route[i-1]]

		//log.Printf("%v - %v - %v", s.route[i-1], s.tsp.demands[s.route[i-1]], carrying)

		// wait to ready to time
		if traveled < s.tsp.readytime[s.route[i]] {
			traveled = s.tsp.readytime[s.route[i]]
		}

		if value, ok := s.tsp.precendense[s.route[i]]; ok {
			for j := i; j > 0; j-- {
				if value == s.route[j] {
					hasNode = true
					break
				}
			}
			if !hasNode {
				// log.Printf("%v - %v", s.route[i], i)
				// log.Print("PRECEDENCE OVERFLOW")
				return false
			}
		}
	}

	// log.Printf("%v - %v - %v", carrying)
	return true
}

// IsFeasible checks if solution is feasible
func (s *Solution) IsFeasibleLog() bool {
	hasNode := false
	traveled := 0
	carrying := 0

	for i := 1; i < s.tsp.numNodes; i++ {
		hasNode = false
		traveled += s.tsp.matrix[s.route[i-1]][s.route[i]]
		carrying += s.tsp.demands[s.route[i-1]]

		//log.Printf("%v - %v - %v", s.route[i-1], s.tsp.demands[s.route[i-1]], carrying)

		// wait to ready to time
		if traveled < s.tsp.readytime[s.route[i]] {
			traveled = s.tsp.readytime[s.route[i]]
		}

		if carrying > s.tsp.capacity {
			log.Print("CAPACITY OVERFLOW")
			return false
		}

		if value, ok := s.tsp.precendense[s.route[i]]; ok {
			for j := i; j > 0; j-- {
				if value == s.route[j] {
					hasNode = true
					break
				}
			}
			if !hasNode {
				// log.Printf("%v - %v", s.route[i], i)
				log.Print("PRECEDENCE OVERFLOW")
				return false
			}
		}

		if s.tsp.duedate[s.route[i]] != 0 && s.tsp.duedate[s.route[i]] < traveled {
			log.Printf("%v TIME WINDOW OVERFLOW", i)
			return false
		}
	}

	// log.Printf("%v - %v - %v", carrying)
	return true
}
