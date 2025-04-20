package workload

import (
	"fmt"
	"math/rand"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"time"
)

type RandomWorkload struct{}

// GenerateWorkload generates a random workload for the given model using the loader's data.
func (rw *RandomWorkload) GenerateWorkload(model *directory.AIModelDefinition) ([]*Job, error) {
	loader := loader.GetLoader()
	if loader == nil {
		return nil, fmt.Errorf("loader is not initialized")
	}

	// Define the number of jobs to generate based on the model's NumberOfRuns
	numJobs := model.NumberOfRuns
	jobList := make([]*Job, numJobs)

	// Generate jobs randomly within the loader's data range
	for index := range numJobs {
		startTime, err := generateRandomTime(loader, model.SLOThreshold)
		if err != nil {
			return nil, err
		}

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
		jobList[index] = job
	}
	// Sort jobs by start time
	return jobList, nil
}

// generateRandomTime generates a random time within the loader's data range.
func generateRandomTime(loader *loader.Loader, SLOTime float64) (time.Time, error) {
	if loader.NumEntries() == 0 {
		return time.Time{}, fmt.Errorf("no data available in loader")
	}

	startDate := loader.StartDate()
	// The last entry in the loader's data is the end date minus the SLOTime
	endDate := loader.EndDate().Add(-time.Duration(SLOTime) * time.Second)

	// Ensure the end date is after the start date
	if endDate.Before(startDate) {
		return time.Time{}, fmt.Errorf("end date is before start date")
	}

	// Generate a random time between startDate and endDate
	randomDuration := time.Duration(rand.Int63n(endDate.Sub(startDate).Nanoseconds()))
	return startDate.Add(randomDuration), nil
}
