package policies

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"simulator/pkg/workload"
	"time"
)

type Temporal struct {
	aiModel *directory.AIModelDefinition
}

func NewTemporal(aiModel *directory.AIModelDefinition) Temporal {
	return Temporal{
		aiModel: aiModel,
	}
}

func (t Temporal) HandleIncoming(job *workload.Job) error {
	loader := loader.GetLoader()
	if loader == nil {
		return fmt.Errorf("loader is not initialized")
	}
	if loader.NumEntries() == 0 {
		return fmt.Errorf("loader has no entries")
	}
	// Assign model job
	job.Model = t.aiModel
	currStartTime := job.StartTime

	bestTime := currStartTime
	minCarbon := math.MaxFloat64
	for currStartTime.Add(time.Duration(t.aiModel.MeanRunTime) * time.Second).Before(job.DueTime) {
		totalCarbon := 0.0
		currTime := currStartTime
		expectedEnd := currStartTime.Add(time.Duration(t.aiModel.MeanRunTime) * time.Second)
		// Calculate carbon emissions for the current time
		carbonIdx, err := loader.GetIndexByDate(currTime)
		if err != nil {
			return fmt.Errorf("error getting index by date: %w", err)
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
			modelRate := job.Model.EnergyUsage                         // in MW
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
			modelRate := job.Model.EnergyUsage                                                  // in MW
			carbon := timeDiff * modelRate * 3.6e-9 * 1e3 * carbonRate                          // in gCO2
			// Update the carbon emission
			totalCarbon += carbon
		}
		// If the carbon emissions are less than the minimum, update the best time and minimum carbon
		log.Printf("[TEMPORAL PREDICT] For start time %s, total carbon is predicted %f", currStartTime.Format(time.RFC3339), totalCarbon)
		if totalCarbon < minCarbon {
			minCarbon = totalCarbon
			bestTime = currStartTime
		}
		// Move to the next time slot which is the next carbon index
		nextTime, err := loader.GetIndexByDate(currStartTime)
		if err != nil {
			return fmt.Errorf("error getting index by date: %w", err)
		}
		if nextTime >= loader.NumEntries()-1 {
			break
		}
		currStartTime = loader.Data[nextTime+1].StartDate
	}
	job.StartTime = bestTime
	duration := max(t.aiModel.MeanRunTime+t.aiModel.StdDevRunTime*rand.NormFloat64(), 0)
	job.EndTime = job.StartTime.Add(time.Duration(duration) * time.Second)
	return nil
}

func (t Temporal) HandleQueued(job *workload.Job) error {
	return nil
}

func (t Temporal) HandleRunning(job *workload.Job) error {
	return nil
}

func (t Temporal) String() string {
	return fmt.Sprintf("Temporal with %s", t.aiModel.ModelName)
}
