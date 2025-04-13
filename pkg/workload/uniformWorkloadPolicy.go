package workload

import (
	"fmt"
	"math/rand"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"time"
)

type UniformWorkload struct{}

// GenerateWorkload generates a uniform workload for the given model using the loader's data.
func (uw *UniformWorkload) GenerateWorkload(model *directory.AIModelDefinition) ([]Job, error) {
	loader := loader.GetLoader()
	if loader == nil {
		return nil, fmt.Errorf("loader is not initialized")
	}

	if loader.NumEntries() == 0 {
		return nil, fmt.Errorf("no data available in loader")
	}

	// Define the number of jobs to generate based on the model's NumberOfRuns
	numJobs := model.NumberOfRuns
	jobList := make([]Job, numJobs)

	// Calculate the time interval between jobs
	startDate := loader.StartDate()
	endDate := loader.EndDate().Add(-time.Duration(model.SLOThreshold) * time.Second)
	totalDuration := endDate.Sub(startDate)
	if totalDuration <= 0 {
		return nil, fmt.Errorf("invalid loader date range")
	}

	interval := totalDuration / time.Duration(numJobs)

	// Generate jobs uniformly spaced within the loader's data range
	for index := range numJobs {
		startTime := startDate.Add(time.Duration(index) * interval)
		duration := time.Duration(rand.NormFloat64()*model.StdDevRunTime+model.MeanRunTime) * time.Second

		job := Job{
			Model:     model,
			StartTime: startTime,
			DueTime:   startTime.Add(time.Duration(model.SLOThreshold) * time.Second),
			Duration:  duration,
		}
		jobList[index] = job
	}

	return jobList, nil
}
