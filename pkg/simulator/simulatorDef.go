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

	incomingJobs         WorkloadQueue
	readyQueue           PriorityQueue
	currentlyRunningJobs PriorityQueue
}
