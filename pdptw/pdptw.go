package pdptw

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path"
	"strconv"
	"strings"

	"github.com/mitas1/psa-core/utils"
)

type TSP interface {
	ReadFromFile(_path string, name string) interface{}
	Print()
	preprocess()
}

type PDPTW struct {
	name        string
	startNode   int
	capacity    int
	numNodes    int
	matrix      [][]int
	readytime   []int
	demands     map[int]int
	duedate     []int
	precendense map[int]int
	pred        map[int]int
	arcs        map[int]map[int]bool
}

// ReadFromFile reads the given tsptw instance from file
func ReadFromFile(_path string, name string) *PDPTW {
	tsp := PDPTW{}
	file, err := os.Open(path.Join(_path, name))

	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()

	tsp.name = name

	scanner := bufio.NewScanner(file)
	var line string
	i := 0

	tsp.matrix = make([][]int, 0)
	tsp.precendense = make(map[int]int)
	tsp.pred = make(map[int]int)
	tsp.demands = make(map[int]int)

	for scanner.Scan() {
		line = strings.Trim(scanner.Text(), " ")
		if line[0] != '#' {
			elems := utils.Map(strings.Fields(line), func(str string) int {
				num, err := strconv.Atoi(str)
				if err != nil {
					log.Fatal(err)
				}
				return num
			})
			if i == 0 {
				tsp.numNodes = elems[0]
				tsp.capacity = elems[1]
				tsp.startNode = 0
				tsp.duedate = make([]int, tsp.numNodes)
				tsp.readytime = make([]int, tsp.numNodes)
			} else if i <= tsp.numNodes {
				tsp.matrix = append(tsp.matrix, elems)
			} else {
				tsp.precendense[elems[1]] = elems[0]

				tsp.pred[elems[1]] = elems[0]

				tsp.pred[elems[0]] = -elems[1]

				tsp.demands[elems[0]] = elems[2]
				tsp.demands[elems[1]] = -elems[2]
				tsp.readytime[elems[0]] = elems[3]
				tsp.duedate[elems[0]] = elems[4]
				tsp.readytime[elems[1]] = elems[5]
				tsp.duedate[elems[1]] = elems[6]
			}
			i++
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &tsp
}

func (tsp *PDPTW) preprocess() {
	tsp.arcs = make(map[int]map[int]bool)
	for i, _ := range tsp.matrix {
		tsp.arcs[i] = make(map[int]bool)
		for j, _ := range tsp.matrix[i] {
			if i != j {
				if tsp.readytime[i]+tsp.matrix[i][j] > tsp.duedate[j] {
					tsp.arcs[i][j] = false
				} else {
					tsp.arcs[i][j] = true
				}
			}
		}
	}
}

// Print the instance in human readable form
func (tsp *PDPTW) Print() {
	fmt.Printf(`=============================PDPTWTW==============================
Instance name:      %s
Number of vertices: %v
Readytime:          %3v
Duedate:            %3v
Demands:          	%v
Precendeces:		%v
`, tsp.name, tsp.numNodes, tsp.readytime, tsp.duedate, tsp.demands, tsp.precendense)
	for _, line := range tsp.matrix {
		fmt.Printf("%2v\n", line)
	}
}
