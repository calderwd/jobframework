package test

import "fmt"

type TestJob struct {
}

func (j TestJob) Process() (bool, error) {
	fmt.Println("Test job running")
	return true, nil
}

type TestJobProfile struct {
}
