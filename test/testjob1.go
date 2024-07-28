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
	for i := 0; i < 3; i++ {
		select {
		case uuid := <-jobCancelStream:
			fmt.Printf("Request to cancel job %s\n", uuid)
			return false, nil
		default:
		}
		time.Sleep(5 * time.Second)
	}
	fmt.Println("Test job end")
	return true, nil
}

type TestJobProfile struct {
}
