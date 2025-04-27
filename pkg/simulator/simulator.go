package simulator

import (
	"container/heap"
	"fmt"
	"log"
	"os"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"simulator/pkg/simulator/policies"
	"simulator/pkg/workload"
	"sync"
	"time"
)

var lock = &sync.Mutex{}
var singleton *Simulator

func NewSimulator(workload []*workload.Job, schedulingPolicy PolicyInterface) *Simulator {
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
				currTime:       loader.StartDate(),
				carbonEmission: make(map[directory.AIModelDefinition]float64),
				sloTimeouts:    make(map[directory.AIModelDefinition]int),

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
			"\tScheduling Policy: %s\n"+
			"\tIncoming Jobs Length: %d\n"+
			"\tQueued Jobs: %v\n"+
			"\tCurrently Running Jobs: %v\n"+
			"\tCompleted Jobs Length: %d\n",
		s.currTime,
		s.carbonEmission,
		s.sloTimeouts,
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
		s.currTime = nextEvent.StartTime
		// Policy assigns the job to be processed
		log.Printf("[INCOMING] Process requested at time %v with due date %v. ", s.currTime.Format(time.ANSIC), nextEvent.DueTime.Format(time.ANSIC))
		s.schedulingPolicy.HandleIncoming(nextEvent)
		log.Printf("[POLICY] Model %s assigned with start at time %v, true end %v.\n", nextEvent.Model.ModelName, nextEvent.StartTime.Format(time.ANSIC), nextEvent.EndTime.Format(time.ANSIC))
		heap.Push(&s.queuedJobs, nextEvent)
	} else if origin == workload.QueuedJob {
		// Fetch job from queued jobs
		nextEvent := heap.Pop(&s.queuedJobs).(*workload.Job)
		s.currTime = nextEvent.StartTime
		// Policy is allowed to make modifications should it choose to
		log.Printf("[AWAITING] Job begins processing at time %v, will complete by %v ", s.currTime.Format(time.ANSIC), nextEvent.EndTime.Format(time.ANSIC))
		s.schedulingPolicy.HandleQueued(nextEvent)
		// Add the job to the currently running jobs
		heap.Push(&s.currentlyRunningJobs, nextEvent)
	} else if origin == workload.RunningJob {
		// Fetch job from currently running jobs
		nextEvent := heap.Pop(&s.currentlyRunningJobs).(*workload.Job)
		s.currTime = nextEvent.EndTime
		// Measure carbon emissions
		s.carbonMeasure(nextEvent)
		// Policy is allowed to make modifications should it choose to
		log.Printf("[COMPLETE] Job completed at %v", s.currTime.Format(time.ANSIC))
		s.schedulingPolicy.HandleRunning(nextEvent)
		// Add the job to the completed jobs
		s.completedJobs.Push(nextEvent)
		// Validate that the job hasn't violated the SLO
		if nextEvent.DueTime.Before(nextJob.EndTime) {
			// SLO violation
			s.sloTimeouts[*nextEvent.Model]++
			log.Printf("[SLO VIOLATION] Job %s with start time %v and end time %v. SLO violated. ", nextEvent.Model.ModelName, nextEvent.StartTime.Format(time.ANSIC), nextEvent.EndTime.Format(time.ANSIC))
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

func (s *Simulator) carbonMeasure(job *workload.Job) error {
	totalCarbon := policies.CarbonCalculate(job.StartTime, job.EndTime, job.Model)
	log.Printf("[EMISSION] Job %s with start time %v and end time %v. Carbon released %f gCO2. ", job.Model.ModelName, job.StartTime.Format(time.ANSIC), job.EndTime.Format(time.ANSIC), totalCarbon)
	s.carbonEmission[*job.Model] += totalCarbon
	return nil
}
