package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"time"

	"github.com/mitas1/psa-core/config"
	"github.com/mitas1/psa-core/pdptw"
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

	for _, file := range files {
		log.Printf("Solving: %s", file.Name())
		tsp := pdptw.ReadFromFile(instancesPath, file.Name())
		start := time.Now()
		s := pdptw.VNS(tsp, config)
		duration := time.Since(start)
		res += fmt.Sprintf("%v	&	%v	&	%v	&	%.4f\n", file.Name(), s.MakeSpan(),
			s.TotalDistance(), duration.Seconds())
		log.Print(res)
		s.WriteToFile(solutionsPath, file.Name())
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
		instancesPath = "_instances/hosny/psa"
	}

	runAll(instancesPath, SOLUTION_PATH, &c)
	return

	log.Printf("Solving: %s", "a.psa")

	tsp := pdptw.ReadFromFile(instancesPath, "01.psa")

	fmt.Printf("%#v", tsp)

	s := pdptw.VNS(tsp, &c)
	s.Print()
	log.Print("done")
}
