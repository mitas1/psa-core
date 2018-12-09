package check

import (
	"bytes"
	"log"
	"os/exec"
	"path"
)

type Check struct {
	pathInstances string
	pathSolutions string
}

func GetChecker(pathInstances string, pathSolutions string) Check {
	return Check{pathInstances, pathSolutions}
}

func (ch *Check) CheckSolution(name string) int {
	cmd := exec.Command("./check_solution", path.Join(ch.pathInstances, name),
		path.Join(ch.pathSolutions, name))

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	err := cmd.Run()
	if err != nil {
		log.Print(err.Error())
		log.Printf(stderr.String())
		log.Printf("%s\n", out.String())
		return -1
	}

	log.Printf("%s\n", out.String())

	return 1
}
