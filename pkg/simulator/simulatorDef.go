package simulator

import (
	"simulator/pkg/directory"
	"time"
)

type SimulatorInterface interface {
}

type Simulator struct {
	currTime       time.Time
	carbonEmission map[directory.AIModelDefinition]float64
	sloViolations  map[directory.AIModelDefinition]int

	incomingJobs         WorkloadQueue // Jobs as they enter
	queuedJobs           PriorityQueue // Jobs that are queued
	currentlyRunningJobs PriorityQueue // Jobs that are currently running
}
