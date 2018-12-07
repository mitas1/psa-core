package gvns

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

type TSP struct {
	name      string
	numNodes  int
	matrix    [][]int
	readytime []int
	duedate   []int
	arcs      map[int]map[int]bool
}

// ReadFromFile reads the given tsptw instance from file
func ReadFromFile(_path string, name string) *TSP {
	tsp := TSP{}
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

	for scanner.Scan() {
		line = strings.Trim(scanner.Text(), " ")
		elems := utils.Map(strings.Fields(line), func(str string) int {
			num, err := strconv.Atoi(str)
			if err != nil {
				log.Fatal(err)
			}
			return num
		})
		if i == 0 {
			tsp.numNodes = elems[0]
		} else if i <= tsp.numNodes {
			tsp.matrix = append(tsp.matrix, elems)
		} else {
			tsp.readytime = append(tsp.readytime, elems[0])
			tsp.duedate = append(tsp.duedate, elems[1])
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &tsp
}

func (tsp *TSP) preprocess() {
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
func (tsp *TSP) Print() {
	fmt.Printf(`=============================TSPTW==============================
Instance name:      %s
Number of vertices: %v
Readytime:          %3v
Duedate:            %3v
`, tsp.name, tsp.numNodes, tsp.readytime, tsp.duedate)
	for _, line := range tsp.matrix {
		fmt.Printf("%2v\n", line)
	}
}
