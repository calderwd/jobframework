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

	uuid, error := jf.AddJob(jobType, jobSchedule, jobContext, "myUser")
	uuid, error = jf.AddJob(jobType2, jobSchedule, jobContext, "myUser")

	time.Sleep(10 * time.Second)
	fmt.Printf("uuid = %s error = %s", uuid, error)

}
