package workload

import (
	"fmt"
	"math/rand"
	"simulator/pkg/loader"
	"slices"
	"time"
)

type RandomWorkload struct{}

// GenerateWorkload generates a random workload for the given model using the loader's data.
func (rw *RandomWorkload) GenerateWorkload(jobInfo JobMetadata) ([]*Job, error) {
	loader := loader.GetLoader()
	if loader == nil {
		return nil, fmt.Errorf("loader is not initialized")
	}

	// Define the number of jobs to generate based on the model's NumberOfRuns
	jobList := make([]*Job, jobInfo.NumJobs)

	// Generate jobs randomly within the loader's data range
	for index := range jobInfo.NumJobs {
		startTime, err := generateRandomTime(jobInfo.DueTime)
		if err != nil {
			return nil, err
		}
		job := &Job{
			Model:     nil,
			StartTime: startTime,
			DueTime:   startTime.Add(jobInfo.DueTime),
			EndTime:  nil,
		}
		jobList[index] = job
	}
	// Sort jobs by start time
	slices.SortFunc(jobList, func(a, b *Job) int {
		return a.StartTime.Compare(b.StartTime)
	})
	return jobList, nil
}

// generateRandomTime generates a random time within the loader's data range.
func generateRandomTime(SLOTime time.Duration) (time.Time, error) {
	loader := loader.GetLoader()
	if loader == nil {
		return time.Time{}, fmt.Errorf("loader is not initialized")
	}
	if loader.NumEntries() == 0 {
		return time.Time{}, fmt.Errorf("no data available in loader")
	}

	startDate := loader.StartDate()
	// The last entry in the loader's data is the end date minus the SLOTime
	endDate := loader.EndDate().Add(-SLOTime)

	// Ensure the end date is after the start date
	if endDate.Before(startDate) {
		return time.Time{}, fmt.Errorf("end date is before start date")
	}

	// Generate a random time between startDate and endDate
	randomDuration := time.Duration(rand.Int63n(endDate.Sub(startDate).Nanoseconds()))
	return startDate.Add(randomDuration), nil
}
