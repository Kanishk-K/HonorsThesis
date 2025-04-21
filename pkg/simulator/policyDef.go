package simulator

import (
	"simulator/pkg/workload"
)

type PolicyInterface interface {
	HandleIncoming(job *workload.Job) error // Assign model and duration for job
	HandleQueued(job *workload.Job) error
	HandleRunning(job *workload.Job) error

	String() string
}
