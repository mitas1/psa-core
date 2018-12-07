package main

import (
	"flag"
	"io/ioutil"
	"log"

	"github.com/mitas1/psa-core/check"

	"github.com/mitas1/psa-core/gvns"
)

func runAll(instancesPath, solutionsPath string, checker *check.Check) {
	files, err := ioutil.ReadDir(instancesPath)
	if err != nil {
		log.Fatal(err)
	}

	errCounter := 0

	for _, file := range files {

		log.Printf("Solving: %s", file.Name())
		tsp := gvns.ReadFromFile(instancesPath, file.Name())
		s := gvns.GVNS(tsp, 15, 8)
		s.WriteToFile()

		if checker.CheckSolution(file.Name()) < 0 {
			log.Printf("Check fails for %s", file.Name())
			errCounter++
		}
		s = nil
		tsp = nil

	}
}

func main() {

	instancesPath := flag.String("instances-path", "../_instances/dumas", "Path to instances")
	solutionsPath := flag.String("solution-path", "../_solutions", "Path to solutions")

	instance := flag.String("instance", "", "Path to specific instance")

	checker := check.GetChecker(*instancesPath, *solutionsPath)

	flag.Parse()

	if *instance != "" {
		tsp := gvns.ReadFromFile(*instancesPath, *instance)
		s := gvns.GVNS(tsp, 15, 8)
		s.WriteToFile()

		if checker.CheckSolution(*instance) < 0 {
			log.Printf("Check fails for %s", *instance)
		}

	} else {
		runAll(*instancesPath, *solutionsPath, &checker)
	}
}
