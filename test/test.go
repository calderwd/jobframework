package test

import (
	"fmt"

	"github.com/calderwd/jobframework/api"
	jobframework "github.com/calderwd/jobframework/pkg"
)

func RunAddTest() {
	fmt.Println("Start")

	jf := jobframework.NewInstance()

	jobType := "my-job-type"
	myJob := TestJob{}
	myJobProfile := TestJobProfile{}

	jf.GetJobRegistrar().RegisterJobType(jobType, myJob, myJobProfile, api.Standard)

	jobSchedule := api.JobSchedule{
		Immediate: true,
	}

	jobContext := api.JobContext{}

	uuid, error := jf.AddJob(jobType, jobSchedule, jobContext, "myUser")

	fmt.Printf("uuid = %s error = %s", uuid, error)

}
