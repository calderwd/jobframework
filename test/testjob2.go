package test

import (
	"fmt"
	"time"

	"github.com/calderwd/jobframework/api"
)

type TestJob2 struct {
}

func (j TestJob2) Process(js api.JobSummary, jobCancelStream <-chan string) (bool, error) {
	name := js.Context.Context
	fmt.Printf("%s is running\n", name)
	for i := 0; i < 3; i++ {
		select {
		case uuid := <-jobCancelStream:
			fmt.Printf("Request to cancel job %s\n", uuid)
			return false, nil
		default:
		}
		time.Sleep(10 * time.Second)
	}
	fmt.Printf("%s has ended\n", name)
	return true, nil
}

type TestJobProfile2 struct {
}
