package workload

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"slices"
	"time"
)

type SpikeWorkload struct {
	PercentOn  float64 // Percent of jobs that are during the spike
	PercentOff float64 // Percent of jobs that are during the off-peak
	TimeOn     int     // Start time of the spike period in hours (0-23)
	TimeOff    int     // End time of the spike period in hours (0-23)
}

func NewMorningSpikeWorkload(percentOn, percentOff float64) *SpikeWorkload {
	return &SpikeWorkload{
		PercentOn:  percentOn,
		PercentOff: percentOff,
		TimeOn:     5,
		TimeOff:    12,
	}
}

func NewAfternoonSpikeWorkload(percentOn, percentOff float64) *SpikeWorkload {
	return &SpikeWorkload{
		PercentOn:  percentOn,
		PercentOff: percentOff,
		TimeOn:     12,
		TimeOff:    17,
	}
}

func NewEveningSpikeWorkload(percentOn, percentOff float64) *SpikeWorkload {
	return &SpikeWorkload{
		PercentOn:  percentOn,
		PercentOff: percentOff,
		TimeOn:     17,
		TimeOff:    24,
	}
}

func NewNightSpikeWorkload(percentOn, percentOff float64) *SpikeWorkload {
	return &SpikeWorkload{
		PercentOn:  percentOn,
		PercentOff: percentOff,
		TimeOn:     0,
		TimeOff:    5,
	}
}

// GenerateWorkload generates a workload with spikes during the specified time range.
func (sp *SpikeWorkload) GenerateWorkload(model *directory.AIModelDefinition) ([]*Job, error) {
	loader := loader.GetLoader()
	if loader == nil {
		return nil, fmt.Errorf("loader is not initialized")
	}

	if loader.NumEntries() == 0 {
		return nil, fmt.Errorf("no data available in loader")
	}

	numOnSpike := int(math.Ceil(float64(model.NumberOfRuns) * sp.PercentOn))
	if numOnSpike == 0 {
		return nil, fmt.Errorf("no jobs to generate")
	}
	if numOnSpike > model.NumberOfRuns {
		numOnSpike = model.NumberOfRuns
	}
	numOffSpike := model.NumberOfRuns - numOnSpike

	numJobs := model.NumberOfRuns
	jobList := make([]*Job, numJobs)

	startDate := loader.StartDate()
	endDate := loader.EndDate().Add(-time.Duration(model.SLOThreshold) * time.Second)

	log.Printf("Putting %d jobs in the morning spike period and %d jobs in the off-peak period\n", numOnSpike, numOffSpike)
	// Assign jobs during the spike period (5 AM to 12 PM)
	for i := range numOnSpike {
		startTime := getStartTimeInRange(startDate, endDate, sp.TimeOn, sp.TimeOff)
		duration := time.Duration(rand.NormFloat64()*model.StdDevRunTime+model.MeanRunTime) * time.Second
		if duration < 0 {
			duration = time.Second // Ensure duration is positive
		}

		job := &Job{
			Model:     model,
			StartTime: startTime,
			DueTime:   startTime.Add(time.Duration(model.SLOThreshold) * time.Second),
			Duration:  duration,
		}
		jobList[i] = job
	}

	// Assign jobs during the off-peak period (12 PM to 5 AM)
	for i := range numOffSpike {
		startTime := getStartTimeInRange(startDate, endDate, sp.TimeOff, sp.TimeOn)
		duration := time.Duration(rand.NormFloat64()*model.StdDevRunTime+model.MeanRunTime) * time.Second
		if duration < 0 {
			duration = time.Second // Ensure duration is positive
		}
		job := &Job{
			Model:     model,
			StartTime: startTime,
			DueTime:   startTime.Add(time.Duration(model.SLOThreshold) * time.Second),
			Duration:  duration,
		}
		jobList[numOnSpike+i] = job
	}
	slices.SortFunc(jobList, func(a, b *Job) int {
		return a.StartTime.Compare(b.StartTime)
	})
	return jobList, nil
}
