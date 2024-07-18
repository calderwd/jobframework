package test

import (
	"context"
	"fmt"

	"github.com/calderwd/jobframework/api"
)

type TestJob2 struct {
}

func (j TestJob2) Process(js api.JobSummary, ctx context.Context) (bool, error) {
	fmt.Println("Test job 2 running")
	return true, nil
}

type TestJobProfile2 struct {
}
