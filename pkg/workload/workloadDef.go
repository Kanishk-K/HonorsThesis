package workload

import (
	"simulator/pkg/directory"
	"time"
)

type PolicyInterface interface {
	GenerateWorkload(model *directory.AIModelDefinition) ([]Job, error)
}

// Defines the various workloads that can be used in the simulation.
type Workload struct {
	// The model this workload is associated with.
	Model *directory.AIModelDefinition
	// The associated jobs for this workload.
	Jobs []Job
	// The policy used to generate the workload.
	Policy PolicyInterface
}

type Job struct {
	Model     *directory.AIModelDefinition // The model this job is associated with.
	StartTime time.Time                    // When the job is queued
	DueTime   time.Time                    // When the job is due before SLO violation
	Duration  time.Duration                // How long the job will take to run
}
