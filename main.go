package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/mitas1/psa-core/config"
	"github.com/mitas1/psa-core/core"
	logging "github.com/op/go-logging"
	"github.com/spf13/pflag"
)

const SOLUTION_PATH = "_solutions"
const LOGGER_FORMAT = `%{color}%{time:2006/02/01-15:04:05.000} %{shortpkg}: %{shortfile} %{level:.4s} %{id:03x}%{color:reset} %{message}`

var logger = logging.MustGetLogger("main")

func parseFlags() (config *string) {
	config = pflag.StringP(
		"config",
		"c",
		"config.yaml",
		"Path to a config file.",
	)
	pflag.Parse()
	return
}

func runAll(instancesPath, solutionsPath string, config *config.Config) {
	files, err := ioutil.ReadDir(instancesPath)
	if err != nil {
		log.Fatal(err)
	}

	res := ""

	vns := core.NewCore(config)

	for _, file := range files {
		log.Printf("Solving: %s", file.Name())
		tsp := core.ReadFromFile(instancesPath, file.Name())
		start := time.Now()
		s, err := vns.Process(tsp)
		if err != nil {
			log.Print(err)
		}
		if s != nil {
			duration := time.Since(start)
			line := fmt.Sprintf("%v	&	%v	&	%v	&	%.4f	%v\n", file.Name(), s.MakeSpan(),
				s.TotalDistance(), duration.Seconds(), s.Check())
			res += line
			fmt.Print(line)
			s.WriteToFile(solutionsPath, file.Name())
		}
	}
	log.Print(res)
}

func main() {
	// Load configuration file
	c := config.Config{}

	config := parseFlags()

	err := c.LoadConfig(*config)
	if err != nil {
		os.Exit(1)
	}

	// Setup logger
	stdOutBackend := logging.NewLogBackend(os.Stdout, "", 0)
	logging.SetBackend(logging.NewBackendFormatter(
		stdOutBackend,
		logging.MustStringFormatter(LOGGER_FORMAT)))

	instancesPath := ""
	//instanceName := ""

	if instancesPath == "" {
		instancesPath = "_instances/wan-rong-jih/psa"
	}

	runAll(instancesPath, SOLUTION_PATH, &c)
	return

	log.Printf("Solving: %s", "a.psa")

	tsp := core.ReadFromFile(instancesPath, "01.psa")

	fmt.Printf("%#v", tsp)

	vns := core.NewCore(&c)

	s, err := vns.Process(tsp)
	if err != nil {
		log.Print(err)
	}
	s.Print()
	log.Print("done")
}
