package pdptw

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
			//log.Print("CAPACITY OVERFLOW")
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
			//log.Print("TIME WINDOW OVERFLOW")
			return false
		}
	}

	// log.Printf("%v - %v - %v", carrying)
	return true
}

// Penalty is sum of all differences between the time to reach each customer
// and its due date
func (s *Solution) Penalty() (penalty int) {
	traveled := 0
	hasNode := false
	carrying := 0

	p_tw := 0
	p_pd := 0
	p_c := 0

	for i := 1; i < len(s.route); i++ {

		hasNode = false
		traveled += s.tsp.matrix[s.route[i-1]][s.route[i]]
		carrying += s.tsp.demands[s.route[i-1]]

		// wait to ready to time
		if traveled < s.tsp.readytime[s.route[i]] {
			traveled = s.tsp.readytime[s.route[i]]
		}

		if carrying > s.tsp.capacity {
			// log.Printf("penalty: %v - %v - %v -> %v", p_c, carrying, s.tsp.capacity, carrying-s.tsp.capacity)
			p_c = p_c + (carrying - s.tsp.capacity)
		}

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
						//log.Printf("%v fucking constraint", value)
						p_pd = p_pd + j
						break
					}
				}
			}
		}

		if s.tsp.duedate[s.route[i]] != 0 && s.tsp.duedate[s.route[i]] < traveled || s.tsp.readytime[s.route[i]] > traveled {
			p_tw = p_tw + traveled - s.tsp.duedate[s.route[i]]
		}
	}

	/* if p_tw == 0 {
		log.Printf("|| %v %v %v", p_pd, p_tw, p_c)
		_, un := s.calcSets()
		tmp := []int{}
		for i := 0; i < len(un); i++ {
			tmp = append(tmp, un[i])
		}
		log.Printf("%v", tmp)
	} */

	//log.Printf("|| %v %v %v", p_pd, p_tw, p_c)

	penalty = 10*p_tw + p_pd + p_c
	return
}

func (s *Solution) calcSets() (feasible []int, unfeasible []int) {
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

		if val, ok := s.tsp.precendense[s.route[i]]; ok {
			for j := i; j > 0; j-- {
				if val == s.route[j] {
					hasNode = true
					break
				}
			}
			if !hasNode {
				tmp = true
			}
		}

		if tmp || (s.tsp.duedate[s.route[i]] != 0 && s.tsp.duedate[s.route[i]] < traveled) || carrying > s.tsp.capacity {
			unfeasible = append(unfeasible, i)
		} else {
			feasible = append(feasible, i)
		}
	}

	//log.Printf("%v - %v", unfeasible, feasible)

	return feasible, unfeasible
}

func (s *Solution) change(i, j int) {
	aux := s.route[i]
	s.route[i] = s.route[j]
	s.route[j] = aux
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

func (s *Solution) MakeSpan() int {
	traveled := 0
	for i := 1; i < len(s.route); i++ {
		traveled += s.tsp.matrix[s.route[i-1]][s.route[i]]

		if traveled < s.tsp.readytime[s.route[i]] {
			traveled = s.tsp.readytime[s.route[i]]
		}

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
