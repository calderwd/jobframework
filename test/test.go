package test

import (
	"fmt"
	"time"

	"github.com/calderwd/jobframework/api"
	jobframework "github.com/calderwd/jobframework/pkg"
)

func RunAddTest() {
	fmt.Println("Start")

	jf := jobframework.NewInstance()

	jobType := "my-job-type"
	myJob := TestJob{}
	myJobProfile := api.DefaultJobProfile()
	jf.GetJobRegistrar().RegisterJobType(jobType, myJob, myJobProfile, api.Standard)

	jobType2 := "my-job2-type"
	myJob2 := TestJob2{}
	myJobProfile2 := api.DefaultJobProfile()
	jf.GetJobRegistrar().RegisterJobType(jobType2, myJob2, myJobProfile2, api.Standard)

	jobSchedule := api.JobSchedule{
		Immediate: true,
	}

	jobContext := api.JobContext{}

	jf.AddJob(jobType, jobSchedule, jobContext, "myUser")
	jf.AddJob(jobType2, jobSchedule, jobContext, "myUser")

	time.Sleep(60 * time.Second)

	jf.Shutdown(false)
}
