package simulator

import (
	"fmt"
	"log"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"simulator/pkg/workload"
	"sync"
	"time"
)

var lock = &sync.Mutex{}
var singleton *Simulator

func NewSimulator(accuracyRequirement float64, workload []*workload.Job, schedulingPolicy string) *Simulator {
	if singleton == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleton == nil {
			loader := loader.GetLoader()
			if loader == nil {
				log.Println("Loader not initialized. Cannot create Simulator.")
				return nil
			}
			if loader.NumEntries() == 0 {
				log.Println("Loader has no data. Cannot create Simulator.")
				return nil
			}
			singleton = &Simulator{
				currTime:               loader.StartDate(),
				carbonEmission:         make(map[directory.AIModelDefinition]float64),
				sloTimeouts:            make(map[directory.AIModelDefinition]int),
				sloAccuracyRequirement: accuracyRequirement,

				schedulingPolicy: schedulingPolicy,

				incomingJobs:         workload,
				queuedJobs:           make(AwaitingHeap, 0),
				currentlyRunningJobs: make(RunningHeap, 0),
				completedJobs:        make(WorkloadQueue, 0),
			}
		} else {
			log.Println("Simulator already initialized")
		}
	} else {
		log.Println("Simulator already initialized")
	}
	return singleton
}

func (s *Simulator) String() string {
	return fmt.Sprintf(
		"\nSimulator State:\n"+
			"\tCurrent Time: %v\n"+
			"\tCarbon Emission: %v\n"+
			"\tSLO Timeouts: %v\n"+
			"\tSLO Accuracy Requirement: %.2f\n"+
			"\tScheduling Policy: %s\n"+
			"\tIncoming Jobs Length: %d\n"+
			"\tQueued Jobs: %v\n"+
			"\tCurrently Running Jobs: %v\n"+
			"\tCompleted Jobs Length: %d\n",
		s.currTime,
		s.carbonEmission,
		s.sloTimeouts,
		s.sloAccuracyRequirement,
		s.schedulingPolicy,
		len(s.incomingJobs),
		s.queuedJobs,
		s.currentlyRunningJobs,
		len(s.completedJobs),
	)
}

func (s *Simulator) Begin() error {
	if err := s.run(); err != nil {
		return fmt.Errorf("error running simulator: %w", err)
	}
	return nil
}

func (s *Simulator) run() error {
	for len(s.incomingJobs) > 0 || s.queuedJobs.Len() > 0 || s.currentlyRunningJobs.Len() > 0 {
		// Run until all jobs are completed
		err := s.update()
		if err != nil {
			return fmt.Errorf("error updating simulator: %w", err)
		}
	}
	return nil
}

func (s *Simulator) update() error {
	// Get the next job to run
	nextJob, origin := s.pickNextEvent()
	if nextJob == nil {
		// Should never happen
		return fmt.Errorf("no next job to run")
	}
	if origin == workload.IncomingJob {
		nextEvent := s.incomingJobs.Pop()
		nextTime := nextEvent.StartTime
		log.Printf("Popped incoming job: %v\n", nextEvent)
		s.carbonMeasure(nextTime)
		nextEvent.EndTime = &nextEvent.StartTime
		s.queuedJobs.Push(nextEvent) // TODO: ADD POLICY TO QUEUE JOBS
	} else if origin == workload.QueuedJob {
		nextEvent := s.queuedJobs.Pop().(*workload.Job)
		nextTime := nextEvent.StartTime
		log.Printf("Popped queued job: %v\n", nextEvent)
		s.carbonMeasure(nextTime)
		s.currentlyRunningJobs.Push(nextEvent)
	} else if origin == workload.RunningJob {
		nextEvent := s.currentlyRunningJobs.Pop().(*workload.Job)
		nextTime := nextEvent.StartTime
		log.Printf("Popped running job: %v\n", nextEvent)
		s.carbonMeasure(nextTime)
		s.completedJobs.Push(nextEvent)
	} else {
		return fmt.Errorf("unknown job origin: %v", origin)
	}
	return nil
}

func (s *Simulator) pickNextEvent() (*workload.Job, workload.JobOrigin) {
	var nextJob *workload.Job
	var origin workload.JobOrigin

	IJJob := s.incomingJobs.Peek()
	AJJob := s.queuedJobs.Peek()
	RJJob := s.currentlyRunningJobs.Peek()

	// Compare the jobs and pick the one with the smallest time
	if IJJob != nil && (AJJob == nil || IJJob.StartTime.Before(AJJob.StartTime)) && (RJJob == nil || IJJob.StartTime.Before(*RJJob.EndTime)) {
		nextJob = IJJob
		origin = workload.IncomingJob
	} else if AJJob != nil && (RJJob == nil || AJJob.StartTime.Before(*RJJob.EndTime)) {
		nextJob = AJJob
		origin = workload.QueuedJob
	} else if RJJob != nil {
		nextJob = RJJob
		origin = workload.RunningJob
	}

	return nextJob, origin
}

func (s *Simulator) carbonMeasure(newTime time.Time) error {
	// While loop until all complete jobs are processed
	
	// Loop over remaining incomplete jobs and calculate carbon emissions

	return nil
}
