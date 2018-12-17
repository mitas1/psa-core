package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"time"

	"github.com/mitas1/psa-core/pdptw"
)

func runAll(instancesPath, solutionsPath string) {
	files, err := ioutil.ReadDir(instancesPath)
	if err != nil {
		log.Fatal(err)
	}

	res := ""

	for _, file := range files {
		log.Printf("Solving: %s", file.Name())
		tsp := pdptw.ReadFromFile(instancesPath, file.Name())
		start := time.Now()
		s := pdptw.VNS(tsp, 1, 50)
		res += fmt.Sprintf("%v	&	%v	&	%v	&	%v	&\n", file.Name(), s.MakeSpan(), s.TotalDistance(), time.Since(start))
		log.Print(res)
		s.WriteToFile(solutionsPath, file.Name())
	}
	log.Print(res)
}

func main() {

	//instancesPath := flag.String("instances-path", "../_instances/dumas", "Path to instances")
	//solutionsPath := flag.String("solution-path", "../_solutions", "Path to solutions")

	instancesPath := *flag.String("instances-path", "", "Path to instances")
	solutionsPath := *flag.String("solution-path", "", "Path to solutions")

	flag.Parse()

	if instancesPath == "" {
		instancesPath = "../_instances/pdptw/psa"
	}

	if solutionsPath == "" {
		solutionsPath = "../_solutions/pdptw"
	}

	runAll(instancesPath, solutionsPath)

	/* 	if *instance != "" {
	   		tsp := tsppd.ReadFromFile(*instancesPath, *instance)
	   		s := tsppd.GVNS(tsp, 15, 8)
	   		s.Print()
	   		s.WriteToFile()

	   	} else {
	   		runAll(*instancesPath, *solutionsPath, nil)
	   	} */
}
