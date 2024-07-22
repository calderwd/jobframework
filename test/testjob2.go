package test

import (
	"fmt"
	"time"

	"github.com/calderwd/jobframework/api"
)

type TestJob2 struct {
}

func (j TestJob2) Process(js api.JobSummary, jobCancelStream <-chan string) (bool, error) {
	fmt.Println("Test job 2 running")
	time.Sleep(10 * time.Second)
	fmt.Println("Test job 2 end")
	return true, nil
}

type TestJobProfile2 struct {
}
