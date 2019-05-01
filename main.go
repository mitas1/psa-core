package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"time"

	_config "github.com/mitas1/psa-core/config"
	"github.com/mitas1/psa-core/core"
	"github.com/mitas1/psa-core/logging"
	"github.com/spf13/pflag"
)

var (
	log = logging.GetLogger()
)

const (
	SOLUTION_PATH = "_solutions"
)

func parseFlags() (config, logFile, instanceName, instancePath *string, iterations *int) {
	config = pflag.StringP(
		"config",
		"c",
		"config.yaml",
		"Path to a config file.",
	)
	logFile = pflag.StringP(
		"logfile",
		"f",
		"psa-core.log",
		"Path to a log file.",
	)
	instanceName = pflag.StringP(
		"instance-path",
		"i",
		"",
		"Path to specific instance file.",
	)
	instancePath = pflag.StringP(
		"instances-path",
		"p",
		"_instances/wan-rong-jih",
		"Path to instances dir.",
	)
	iterations = pflag.IntP(
		"iterations",
		"t",
		1,
		"Number of iterations.",
	)
	pflag.Parse()
	return
}

type solver struct {
	core *core.Core
}

func (s solver) solveInstance(_path, name string, maxIter int) (latexOut string) {
	instancePath := path.Join(_path, name)
	file, err := os.Open(instancePath)
	if err != nil {
		log.Fatal(err)
	}

	fileInfo, err := file.Stat()
	if err != nil {
		log.Fatal(err)
	}

	if fileInfo.IsDir() {
		log.Warningf("Skipping directory: %v", name)
		return
	}

	pdptw := core.ReadFromFile(_path, name)

	var totalDuration float64
	var totalObjective int

	for iteration := 1; iteration < maxIter+1; iteration++ {
		log.Infof("Solving instance: %v, iteration: %d", name, iteration)
		start := time.Now()
		sol, err := s.core.Process(pdptw)
		if err != nil {
			log.Error(err)
			return
		}
		duration := time.Since(start)
		log.Infof(`Instance solved!
		Name:			%v
		Tasks:			%v
		Objective:		%v
		Duration:		%.4f (s)
		Checks:			%v`, name, pdptw.NumberOfTasks(), sol.MakeSpan(),
			duration.Seconds(), sol.Check())

		totalObjective += sol.MakeSpan()
		totalDuration += duration.Seconds()
	}

	latexOut = fmt.Sprintf("%v	&	%v	&	%.4f\n", pdptw.NumberOfTasks(),
		totalObjective / maxIter, totalDuration / float64(maxIter))

	return
}

func main() {
	config, file, instanceName, instancesPath, iterations := parseFlags()

	log = logging.SetupLogger(file)

	c := _config.Config{}

	if err := c.LoadConfig(*config); err != nil {
		log.Fatal(err)
	}

	log = logging.SetupLogger(file)

	var latex string

	solver := solver{core: core.NewCore(&c)}

	if instanceName != nil && *instanceName != "" {
		latex += solver.solveInstance("", *instanceName, *iterations)
	} else {
		files, err := ioutil.ReadDir(*instancesPath)
		if err != nil {
			log.Fatal(err)
		}
		for _, file := range files {
			latex += solver.solveInstance(*instancesPath, file.Name(), *iterations)
		}
	}

	log.Infof("Well done copy paste the following into your awesome thesis")
	fmt.Print(latex)

	return
}
