package workload

import (
	"fmt"
	"log"
	"math"
	"simulator/pkg/loader"
	"slices"
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
func (sp *SpikeWorkload) GenerateWorkload(jobInfo JobMetadata) ([]*Job, error) {
	loader := loader.GetLoader()
	if loader == nil {
		return nil, fmt.Errorf("loader is not initialized")
	}

	if loader.NumEntries() == 0 {
		return nil, fmt.Errorf("no data available in loader")
	}

	numJobs := jobInfo.NumJobs
	numOnSpike := int(math.Ceil(float64(numJobs) * sp.PercentOn))
	if numOnSpike == 0 {
		return nil, fmt.Errorf("no jobs to generate")
	}
	if numOnSpike > numJobs {
		numOnSpike = numJobs
	}
	numOffSpike := numJobs - numOnSpike

	jobList := make([]*Job, numJobs)

	startDate := loader.StartDate()
	endDate := loader.EndDate().Add(-jobInfo.DueTime)

	log.Printf("Putting %d jobs in the morning spike period and %d jobs in the off-peak period\n", numOnSpike, numOffSpike)
	// Assign jobs during the spike period (5 AM to 12 PM)
	for i := range numOnSpike {
		startTime := getStartTimeInRange(startDate, endDate, sp.TimeOn, sp.TimeOff)

		job := &Job{
			Model:     nil,
			StartTime: startTime,
			DueTime:   startTime.Add(jobInfo.DueTime),
			Duration:  nil,
		}
		jobList[i] = job
	}

	// Assign jobs during the off-peak period (12 PM to 5 AM)
	for i := range numOffSpike {
		startTime := getStartTimeInRange(startDate, endDate, sp.TimeOff, sp.TimeOn)
		job := &Job{
			Model:     nil,
			StartTime: startTime,
			DueTime:   startTime.Add(jobInfo.DueTime),
			Duration:  nil,
		}
		jobList[numOnSpike+i] = job
	}
	slices.SortFunc(jobList, func(a, b *Job) int {
		return a.StartTime.Compare(b.StartTime)
	})
	return jobList, nil
}
