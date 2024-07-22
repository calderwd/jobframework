package test

import (
	"fmt"
	"time"

	"github.com/calderwd/jobframework/api"
)

type TestJob struct {
}

func (j TestJob) Process(js api.JobSummary, jobCancelStream <-chan string) (bool, error) {
	fmt.Println("Test job running")
	time.Sleep(10 * time.Second)
	fmt.Println("Test job end")
	return true, nil
}

type TestJobProfile struct {
}
