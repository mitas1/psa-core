package core

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

func CreateInstance(
	startNode int,
	vehicleCapacity int,
	traveled int,
	carrying int,
	readyTime []int,
	dueDate []int,
	demands map[int]int,
	precedence map[int]int,
	matrix [][]int,
) PDPTW {
	return PDPTW{
		name:       "instance",
		startNode:  startNode,
		capacity:   vehicleCapacity,
		traveled:   traveled,
		carrying:   carrying,
		numNodes:   len(matrix),
		readyTime:  readyTime,
		dueDate:    dueDate,
		demands:    demands,
		precedence: precedence,
		matrix:     matrix,
	}
}

type PDPTW struct {
	name       string
	startNode  int
	capacity   int
	numNodes   int
	traveled   int
	carrying   int
	matrix     [][]int
	readyTime  []int
	dueDate    []int
	demands    map[int]int
	precedence map[int]int
	pred       map[int]int
	arcs       map[int]map[int]bool
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
	tsp.precedence = make(map[int]int)
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
				// number of nodes
				tsp.numNodes = elems[0]
				// capacity of vehicle
				tsp.capacity = elems[1]
				// start node
				tsp.startNode = elems[2]

				// init traveled and carrying if instance contains
				if len(elems) > 3 {
					tsp.traveled = elems[3]
					tsp.carrying = elems[4]
				}
				tsp.dueDate = make([]int, tsp.numNodes)
				tsp.readyTime = make([]int, tsp.numNodes)
			} else if i <= tsp.numNodes {
				tsp.matrix = append(tsp.matrix, elems)
			} else {
				if len(elems) == 7 {
					tsp.precedence[elems[1]] = elems[0]

					// TODO: Remove unused pred
					tsp.pred[elems[1]] = elems[0]
					tsp.pred[elems[0]] = -elems[1]

					tsp.demands[elems[0]] = elems[2]
					tsp.demands[elems[1]] = -elems[2]
					tsp.readyTime[elems[0]] = elems[3]
					tsp.dueDate[elems[0]] = elems[4]
					tsp.readyTime[elems[1]] = elems[5]
					tsp.dueDate[elems[1]] = elems[6]
				} else if len(elems) == 4 {
					tsp.precedence[elems[1]] = -1
					tsp.demands[elems[0]] = elems[1]
					tsp.readyTime[elems[0]] = elems[2]
					tsp.dueDate[elems[0]] = elems[3]
				} else {
					log.Printf("Wrong task format")
				}
			}
			i++
		}

	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &tsp
}

func (tsp *PDPTW) NumberOfTasks() int {
	return tsp.numNodes / 2
}

func (tsp *PDPTW) preprocess() {
	tsp.arcs = make(map[int]map[int]bool)
	for i, _ := range tsp.matrix {
		tsp.arcs[i] = make(map[int]bool)
		for j, _ := range tsp.matrix[i] {
			if i != j {
				if tsp.readyTime[i]+tsp.matrix[i][j] > tsp.dueDate[j] {
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
readyTime:          %3v
dueDate:            %3v
Demands:          	%v
Precendeces:		%v
`, tsp.name, tsp.numNodes, tsp.readyTime, tsp.dueDate, tsp.demands, tsp.precedence)
	for _, line := range tsp.matrix {
		fmt.Printf("%2v\n", line)
	}
}
