package simulator

import (
	"simulator/pkg/directory"
	"time"
)

type SimulatorInterface interface {
	// Private Methods
	run() error
	update() error
	carbonMeasure(newTime time.Time) error

	// Public Methods
	String() string
	Begin() error
}

type Simulator struct {
	currTime               time.Time
	carbonEmission         map[directory.AIModelDefinition]float64
	sloTimeouts            map[directory.AIModelDefinition]int
	sloAccuracyRequirement float64

	schedulingPolicy string

	incomingJobs         WorkloadQueue // Jobs as they enter
	queuedJobs           AwaitingHeap  // Jobs that are queued
	currentlyRunningJobs RunningHeap   // Jobs that are currently running
	completedJobs        WorkloadQueue // Jobs that have completed
}
