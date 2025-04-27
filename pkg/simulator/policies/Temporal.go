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

type Temporal struct {
	aiModel     *directory.AIModelDefinition
	safeguardSD float64
}

func NewTemporal(aiModel *directory.AIModelDefinition, safeguardSD float64) *Temporal {
	return &Temporal{
		aiModel:     aiModel,
		safeguardSD: safeguardSD,
	}
}

func (t *Temporal) HandleIncoming(job *workload.Job) error {
	// Assign model job
	job.Model = t.aiModel
	bestTime, carbonPredict, _ := TemporalCarbonEstimate(job, t.aiModel, t.safeguardSD)
	if !bestTime.Equal(job.StartTime) {
		estimatedEnd := bestTime.Add(time.Duration(t.aiModel.MeanRunTime) * time.Second).Add(time.Duration(t.aiModel.StdDevRunTime*t.safeguardSD) * time.Second)
		log.Printf("[TEMPORAL SHIFT PREDICT] For start time %s, estimated end %s, and model %s, total carbon is predicted %f gCO2", bestTime.Format(time.ANSIC), estimatedEnd.Format(time.ANSIC), t.aiModel.ModelName, carbonPredict)
	} else {
		estimatedEnd := job.StartTime.Add(time.Duration(t.aiModel.MeanRunTime) * time.Second).Add(time.Duration(t.aiModel.StdDevRunTime*t.safeguardSD) * time.Second)
		log.Printf("[TEMPORAL NO CHANGE PREDICT] For start time %s, estimated end %s, and model %s, total carbon is predicted %f gCO2", job.StartTime.Format(time.ANSIC), estimatedEnd.Format(time.ANSIC), t.aiModel.ModelName, carbonPredict)
	}
	job.StartTime = bestTime
	duration := max(t.aiModel.MeanRunTime+t.aiModel.StdDevRunTime*rand.NormFloat64(), 0)
	job.EndTime = job.StartTime.Add(time.Duration(duration) * time.Second)
	return nil
}

func (t *Temporal) HandleQueued(job *workload.Job) error {
	return nil
}

func (t *Temporal) HandleRunning(job *workload.Job) error {
	return nil
}

func (t *Temporal) String() string {
	return fmt.Sprintf("Temporal with %s and Standard Deviation Guard: %f", t.aiModel.ModelName, t.safeguardSD)
}

func TemporalCarbonEstimate(job *workload.Job, aiModel *directory.AIModelDefinition, safeguardSD float64) (time.Time, float64, error) {
	loader := loader.GetLoader()
	if loader == nil {
		return time.Time{}, 0, fmt.Errorf("loader not initialized")
	}
	if loader.NumEntries() == 0 {
		return time.Time{}, 0, fmt.Errorf("loader has no data")
	}
	// Default values should there not be space to temporally shift
	bestTime := job.StartTime
	currTime := job.StartTime
	currEnd := job.StartTime.Add(time.Duration(aiModel.MeanRunTime) * time.Second).Add(time.Duration(aiModel.StdDevRunTime*safeguardSD) * time.Second)
	minCarbon := CarbonCalculate(job.StartTime, currEnd, aiModel)

	for currEnd.Before(loader.EndDate()) && currEnd.Before(job.DueTime) {
		carbonIdx, err := loader.GetIndexByDate(currTime)
		if err != nil {
			return time.Time{}, 0, fmt.Errorf("error getting index by date: %v", err)
		}
		if carbonIdx >= loader.NumEntries()-1 {
			// We can't shift the job to a later time
			return bestTime, minCarbon, nil
		}
		carbon := CarbonCalculate(currTime, currEnd, aiModel)
		if carbon < minCarbon {
			// log.Printf("[TEMPORAL PREDICT] For start time %s, estimated end %s, and model %s, total carbon is predicted %f gCO2", currTime.Format(time.ANSIC), currEnd.Format(time.ANSIC), aiModel.ModelName, carbon)
			minCarbon = carbon
			bestTime = currTime
		}
		currTime = loader.Data[carbonIdx+1].StartDate
		currEnd = currTime.Add(time.Duration(aiModel.MeanRunTime) * time.Second).Add(time.Duration(aiModel.StdDevRunTime*safeguardSD) * time.Second)
	}

	return bestTime, minCarbon, nil
}
