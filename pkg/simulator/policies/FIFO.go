package policies

import (
	"fmt"
	"log"
	"math/rand"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"simulator/pkg/workload"
	"time"
)

type FIFO struct {
	aiModel *directory.AIModelDefinition
}

func NewFIFO(aiModel *directory.AIModelDefinition) *FIFO {
	return &FIFO{
		aiModel: aiModel,
	}
}

func (f *FIFO) HandleIncoming(job *workload.Job) error {
	job.Model = f.aiModel
	_, err := CarbonEstimate(job, *f.aiModel)
	if err != nil {
		return fmt.Errorf("error estimating carbon: %w", err)
	}
	// Generate the duration of the job
	duration := max(f.aiModel.MeanRunTime+f.aiModel.StdDevRunTime*rand.NormFloat64(), 0)
	job.EndTime = job.StartTime.Add(time.Duration(duration) * time.Second)
	return nil
}

func (f *FIFO) HandleQueued(job *workload.Job) error {
	return nil
}

func (f *FIFO) HandleRunning(job *workload.Job) error {
	return nil
}

func (f *FIFO) String() string {
	return fmt.Sprintf("FIFO with %s", f.aiModel.ModelName)
}

func CarbonEstimate(job *workload.Job, aiModel directory.AIModelDefinition) (float64, error) {
	loader := loader.GetLoader()
	if loader == nil {
		return 0, fmt.Errorf("loader not initialized")
	}
	if loader.NumEntries() == 0 {
		return 0, fmt.Errorf("loader has no data")
	}
	totalCarbon := 0.0
	currTime := job.StartTime
	expectedEnd := job.StartTime.Add(time.Duration(aiModel.MeanRunTime) * time.Second)
	carbonIdx, err := loader.GetIndexByDate(currTime)
	if err != nil {
		return 0, fmt.Errorf("error getting index by date: %w", err)
	}
	for carbonIdx < loader.NumEntries()-1 && currTime.Before(expectedEnd) {
		// Find the next time, smaller of either the nextEntry or the newTime
		nextTime := loader.Data[carbonIdx+1].StartDate
		if nextTime.After(expectedEnd) {
			nextTime = expectedEnd
		}
		// Calculate the time difference
		timeDiff := nextTime.Sub(currTime).Seconds() // in seconds
		// Calculate the carbon emission
		carbonRate := loader.Data[carbonIdx].CarbonIntensity       // in kgCO2/MWh
		modelRate := aiModel.EnergyUsage                           // in MW
		carbon := timeDiff * modelRate * 3.6e-9 * 1e3 * carbonRate // in gCO2
		// Update the carbon emission
		totalCarbon += carbon
		currTime = nextTime
		carbonIdx++
	}
	if expectedEnd.After(loader.EndDate()) {
		// If the newTime is after all recorded data points, use the last entry as a heuristic
		timeDiff := expectedEnd.Sub(loader.Data[loader.NumEntries()-1].StartDate).Seconds() // in seconds
		carbonRate := loader.Data[loader.NumEntries()-1].CarbonIntensity                    // in kgCO2/MWh
		modelRate := aiModel.EnergyUsage                                                    // in MW
		carbon := timeDiff * modelRate * 3.6e-9 * 1e3 * carbonRate                          // in gCO2
		// Update the carbon emission
		totalCarbon += carbon
	}
	log.Printf("[FIFO PREDICT] For start time %s and model %s, total carbon is predicted %f gCO2", job.StartTime.Format(time.RFC3339), aiModel.ModelName, totalCarbon)
	return totalCarbon, nil
}
