package pkg

import (
	"github.com/calderwd/jobframework/api"
	"github.com/calderwd/jobframework/internal/jobmanager"
)

func NewInstance() api.JobService {
	return jobmanager.GetInstance()
}
