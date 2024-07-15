package main

import (
	"fmt"

	"github.com/calderwd/jobframework/test"
)

type TestJob struct {
}

func (j TestJob) Process() (bool, error) {
	fmt.Println("Test job running")
	return true, nil
}

type TestJobProfile struct {
}

func main() {
	fmt.Println("Start")

	test.RunAddTest()

}
