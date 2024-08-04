package api

import (
	"fmt"
	"time"

	"github.com/google/uuid"
)

type JobSchedule struct {
	SimplePeriodic string
	Immediate      bool
	Cron           string
}

type JobContext struct {
	Context string
}

type JobState int

func (s JobState) String() string {
	switch s {
	case 0:
		return "Idle"
	case 1:
		return "Scheduled"
	case 2:
		return "Rescheduler"
	case 3:
		return "Running"
	case 4:
		return "Deleting"
	case 5:
		return "Complete"
	case 6:
		return "Cancelled"
	}
	return "Unknown"
}

const (
	Idle        = iota // Has been added to the database
	Scheduled          // Waiting for first execution
	Rescheduled        // Has Completed execution and been rescheduled
	Running            // Currently executing
	Deleting           // Deleted from the system
	Complete           // Single execution has complete (not rescheduled)
	Cancelled          // // Job has been cancelled
)

type JobStatus int

func (s JobStatus) String() string {
	switch s {
	case 0:
		return "Calculating"
	case 1:
		return "Success"
	case 2:
		return "Failure"
	}
	return "Unknown"
}

const (
	Calculating = iota // Not yet complete
	Success            // Complete and successful
	Failure            // Complete but failed
)

type JobSummary struct {
	Uuid               uuid.UUID
	Name               string
	Description        string
	State              JobState
	Status             JobStatus
	JobType            string
	Progress           uint8
	Context            JobContext
	EvaluationId       uint64
	Schedule           JobSchedule
	LastExecutionStart *time.Time
	LastExeuctionEnd   *time.Time
	NextExecutionTime  *time.Time
}

func (js JobSummary) Build(jobType string, jobSchedule JobSchedule, jobContext JobContext) JobSummary {
	js.Uuid = uuid.New()
	js.Name = ""
	js.Description = ""
	js.State = Complete
	js.Status = Failure
	js.JobType = jobType
	js.Progress = 0
	js.EvaluationId = 0
	js.Schedule = jobSchedule
	js.LastExecutionStart = nil
	js.LastExeuctionEnd = nil
	js.NextExecutionTime = nil
	js.Context = jobContext

	return js
}

func (js JobSummary) Dump() string {

	return fmt.Sprintf(" { uuid : %s, name : %s, Desc : %s, State : %s, Status : %s, JobType : %s, Progress %d, Eval : %d", js.Uuid, js.Name, js.Description, js.State.String(), js.Status.String(), js.JobType, js.Progress, js.EvaluationId)
}

type JobFilter struct {
}

type JobFilterEntry struct {
}

type JobProfile interface {
	CanAdd() bool
}

type ScheduleType int

const (
	Standard   = iota // Standard scheduling of background jobs
	Sequential        // Only one such job executes at a time
	Periodic          // Runs periodically
	Cron              // Scheduled at a specific time
)

type IJob interface {
	Process(js JobSummary, jobCancelStream <-chan string) (bool, error)
}

type JobConfig struct {
	JobType   string
	Job       IJob
	Scheduler ScheduleType
	Profile   JobProfile
}

type JobRegistrar interface {
	RegisterJobType(jobType string, job IJob, profile JobProfile, scheduleType ScheduleType)
	GetJobConfig(jobType string) (JobConfig, error)
}

type JobService interface {
	GetJobRegistrar() JobRegistrar
	AddJob(jobType string, jobSchedule JobSchedule, jobContext JobContext, user string) (uuid.UUID, error)
	GetJob(uuid uuid.UUID, user string) (JobSummary, error)
	CancelJob(uuid uuid.UUID, user string) bool
	DeleteJob(uuid uuid.UUID, user string) error
	ListJobs(filter JobFilter, user string) []JobSummary
	GetEvaluation(uuid uuid.UUID, evaluationId uint64, user string) (string, error)
	GetJobLogPageCount(uuid uuid.UUID, evaluationId uint64, user string) (uint, error)
	GetJobLogPage(uuid uuid.UUID, evaluationId uint64, pageNumber uint, user string) ([]string, error)
	GetJobHistory(uuid uuid.UUID, user string) ([]JobSummary, error)
	Shutdown(force bool)
}
