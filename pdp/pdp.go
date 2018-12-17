package pdp

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

type PDP struct {
	name        string
	numNodes    int
	matrix      [][]int
	arcs        map[int]map[int]bool
	precendense map[int]int
}

// ReadFromFile reads the given tsptw instance from file
func ReadFromFile(_path string, name string) *PDP {
	tsp := PDP{}
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
			tsp.precendense[elems[1]] = elems[0]
		}
		i++
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}
	return &tsp
}

// Print the instance in human readable form
func (tsp *PDP) Print() {
	fmt.Printf(`=============================PDPTW==============================
Instance name:      %s
Number of vertices: %v
Precedense:         %v
`, tsp.name, tsp.numNodes, tsp.precendense)
	for _, line := range tsp.matrix {
		fmt.Printf("%2v\n", line)
	}
}
