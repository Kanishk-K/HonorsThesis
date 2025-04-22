package simulator

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"simulator/pkg/workload"
	"sync"
	"time"
)

var lock = &sync.Mutex{}
var singleton *Simulator

func NewSimulator(accuracyRequirement float64, workload []*workload.Job, schedulingPolicy PolicyInterface) *Simulator {
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
			queueJobHeap := make(AwaitingHeap, 0)
			heap.Init(&queueJobHeap)
			runningJobHeap := make(RunningHeap, 0)
			heap.Init(&runningJobHeap)
			singleton = &Simulator{
				currTime:               loader.StartDate(),
				carbonEmission:         make(map[directory.AIModelDefinition]float64),
				sloTimeouts:            make(map[directory.AIModelDefinition]int),
				sloAccuracyRequirement: accuracyRequirement,

				schedulingPolicy: schedulingPolicy,

				incomingJobs:         workload,
				queuedJobs:           queueJobHeap,
				currentlyRunningJobs: runningJobHeap,
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
	// Create a log file with the current date and time
	logFileName := fmt.Sprintf("simulator_log_%s.log", time.Now().Format("2006-01-02_15-04-05"))
	logFile, err := os.Create(logFileName)
	if err != nil {
		return fmt.Errorf("error creating log file: %w", err)
	}
	defer logFile.Close()
	log.SetFlags(log.Flags() &^ (log.Ldate | log.Ltime))
	log.SetOutput(logFile)
	// Run the simulator
	if err := s.run(); err != nil {
		return fmt.Errorf("error running simulator: %w", err)
	}
	// Set log output to stdout
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Flags() | log.Ldate | log.Ltime)
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
		// Fetch job from incoming jobs
		nextEvent := s.incomingJobs.Pop()
		nextTime := nextEvent.StartTime
		// Measure carbon emissions
		s.carbonMeasure(nextTime)
		// Policy assigns the job to be processed
		log.Printf("Removed job from incoming jobs at time %v with due date %v. ", nextTime.Format(time.RFC3339), nextEvent.DueTime.Format(time.RFC3339))
		s.schedulingPolicy.HandleIncoming(nextEvent)
		log.Printf("Model %s assigned with termination at time %v.\n", nextEvent.Model.ModelName, nextEvent.EndTime.Format(time.RFC3339))
		s.queuedJobs.Push(nextEvent)
	} else if origin == workload.QueuedJob {
		// Fetch job from queued jobs
		nextEvent := s.queuedJobs.Pop().(*workload.Job)
		nextTime := nextEvent.StartTime
		// Measure carbon emissions
		s.carbonMeasure(nextTime)
		// Policy is allowed to make modifications should it choose to
		log.Printf("Removed job from queued jobs at time %v with due date %v. ", nextTime.Format(time.RFC3339), nextEvent.DueTime.Format(time.RFC3339))
		s.schedulingPolicy.HandleQueued(nextEvent)
		// Add the job to the currently running jobs
		s.currentlyRunningJobs.Push(nextEvent)
	} else if origin == workload.RunningJob {
		// Fetch job from currently running jobs
		nextEvent := s.currentlyRunningJobs.Peek()
		nextTime := nextEvent.EndTime
		// Measure carbon emissions
		s.carbonMeasure(nextTime)
		nextEvent = s.currentlyRunningJobs.Pop().(*workload.Job)
		// Policy is allowed to make modifications should it choose to
		log.Printf("Removed job from currently running jobs at time %v with due date %v. ", nextTime.Format(time.RFC3339), nextEvent.DueTime.Format(time.RFC3339))
		s.schedulingPolicy.HandleRunning(nextEvent)
		// Add the job to the completed jobs
		s.completedJobs.Push(nextEvent)
		// Validate that the job hasn't violated the SLO
		if nextEvent.DueTime.Before(nextTime) {
			// SLO violation
			s.sloTimeouts[*nextEvent.Model]++
		}
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
	if IJJob != nil && (AJJob == nil || IJJob.StartTime.Before(AJJob.StartTime)) && (RJJob == nil || IJJob.StartTime.Before(RJJob.EndTime)) {
		nextJob = IJJob
		origin = workload.IncomingJob
	} else if AJJob != nil && (RJJob == nil || AJJob.StartTime.Before(RJJob.EndTime)) {
		nextJob = AJJob
		origin = workload.QueuedJob
	} else if RJJob != nil {
		nextJob = RJJob
		origin = workload.RunningJob
	}

	return nextJob, origin
}

func (s *Simulator) carbonMeasure(newTime time.Time) error {
	// newTime is always less than or equal to the end time of the runningQueue
	loader := loader.GetLoader()
	if loader == nil {
		return fmt.Errorf("loader not initialized")
	}
	if loader.NumEntries() == 0 {
		return fmt.Errorf("loader has no data")
	}
	totalCarbon := 0.0
	for _, job := range s.currentlyRunningJobs {
		carbonIdx, err := loader.GetIndexByDate(s.currTime)
		if err != nil {
			return fmt.Errorf("error getting index by date: %w", err)
		}
		currTime := s.currTime
		// Iterate until there are no more data points or the newTime is reached
		for carbonIdx < loader.NumEntries()-1 && currTime.Before(newTime) {
			// Find the next time, smaller of either the nextEntry or the newTime
			nextTime := loader.Data[carbonIdx+1].StartDate
			if nextTime.After(newTime) {
				nextTime = newTime
			}
			// Calculate the time difference
			timeDiff := nextTime.Sub(currTime).Seconds() // in seconds
			// Calculate the carbon emission
			carbonRate := loader.Data[carbonIdx].CarbonIntensity       // in kgCO2/MWh
			modelRate := job.Model.EnergyUsage                         // in MW
			carbon := timeDiff * modelRate * 3.6e-9 * 1e3 * carbonRate // in gCO2
			// Update the carbon emission
			s.carbonEmission[*job.Model] += carbon
			totalCarbon += carbon
			currTime = nextTime
			carbonIdx++
		}
		if newTime.After(loader.EndDate()) {
			// If the newTime is after all recorded data points, use the last entry as a heuristic
			timeDiff := newTime.Sub(loader.Data[loader.NumEntries()-1].StartDate).Seconds() // in seconds
			carbonRate := loader.Data[loader.NumEntries()-1].CarbonIntensity                // in kgCO2/MWh
			modelRate := job.Model.EnergyUsage                                              // in MW
			carbon := timeDiff * modelRate * 3.6e-9 * 1e3 * carbonRate                      // in gCO2
			// Update the carbon emission
			s.carbonEmission[*job.Model] += carbon
			totalCarbon += carbon
			currTime = newTime
		}
	}
	s.currTime = newTime
	log.Printf("Total carbon emission from %v to %v: %.2f gCO2\n", s.currTime.Format(time.RFC3339), newTime.Format(time.RFC3339), totalCarbon)
	return nil
}
