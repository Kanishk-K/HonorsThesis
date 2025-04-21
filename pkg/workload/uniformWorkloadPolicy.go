package workload

import (
	"fmt"
	"simulator/pkg/loader"
	"time"
)

type UniformWorkload struct{}

// GenerateWorkload generates a uniform workload for the given model using the loader's data.
func (uw *UniformWorkload) GenerateWorkload(jobInfo JobMetadata) ([]*Job, error) {
	loader := loader.GetLoader()
	if loader == nil {
		return nil, fmt.Errorf("loader is not initialized")
	}

	if loader.NumEntries() == 0 {
		return nil, fmt.Errorf("no data available in loader")
	}

	// Define the number of jobs to generate based on the model's NumberOfRuns
	jobList := make([]*Job, jobInfo.NumJobs)

	// Calculate the time interval between jobs
	startDate := loader.StartDate()
	endDate := loader.EndDate().Add(-jobInfo.DueTime)
	totalDuration := endDate.Sub(startDate)
	if totalDuration <= 0 {
		return nil, fmt.Errorf("invalid loader date range")
	}

	interval := totalDuration / time.Duration(jobInfo.NumJobs)

	// Generate jobs uniformly spaced within the loader's data range
	for index := range jobInfo.NumJobs {
		startTime := startDate.Add(time.Duration(index) * interval)

		job := &Job{
			Model:     nil,
			StartTime: startTime,
			DueTime:   startTime.Add(jobInfo.DueTime),
			EndTime:   startTime,
		}
		jobList[index] = job
	}

	return jobList, nil
}
