package workload

import (
	"simulator/pkg/directory"
	"time"
)

type PolicyInterface interface {
	GenerateWorkload(jobInfo JobMetadata) ([]*Job, error)
}

// Defines the various workloads that can be used in the simulation.
type Workload struct {
	// The associated jobs for this workload.
	Jobs []*Job
	// The policy used to generate the workload.
	Policy PolicyInterface
}

type Job struct {
	Model     *directory.AIModelDefinition // The model this job is associated with.
	StartTime time.Time                    // When the job is queued
	DueTime   time.Time                    // When the job is due before SLO violation
	EndTime   *time.Time                   // How long the job will take to run
}

type JobMetadata struct {
	// The time the job is due before SLO violation
	DueTime        time.Duration
	NumJobs        int
	WorkloadPolicy string
}

type JobOrigin int

const (
	IncomingJob JobOrigin = iota
	QueuedJob
	RunningJob
)
