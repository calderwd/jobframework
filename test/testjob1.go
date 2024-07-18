package test

import (
	"context"
	"fmt"

	"github.com/calderwd/jobframework/api"
)

type TestJob struct {
}

func (j TestJob) Process(js api.JobSummary, ctx context.Context) (bool, error) {
	fmt.Println("Test job running")
	return true, nil
}

type TestJobProfile struct {
}
