package policies

import (
	"fmt"
	"math"
	"math/rand"
	"simulator/pkg/directory"
	"simulator/pkg/workload"
	"time"
)

type HybridSelection struct {
	requiredAccuracy  float64
	currTotalAccuracy float64
	processedJobs     int
}

func NewHybridSelection(requiredAccuracy float64) *HybridSelection {
	return &HybridSelection{
		requiredAccuracy:  requiredAccuracy,
		currTotalAccuracy: 0,
		processedJobs:     0,
	}
}

func (h *HybridSelection) HandleIncoming(job *workload.Job) error {
	modelDirectory := directory.FetchDirectory()
	if modelDirectory == nil {
		return fmt.Errorf("model directory not initialized")
	}
	models := modelDirectory.GetModels()
	var selectedModel *directory.AIModelDefinition
	bestStartTime := job.StartTime
	bestCarbon := math.MaxFloat64
	bestAccuracy := 0.0
	for _, model := range models {
		newAccuracy := (h.currTotalAccuracy + model.Accuracy) / float64(h.processedJobs+1)
		if newAccuracy >= h.requiredAccuracy {
			startTime, carbonEstimate, err := TemporalCarbonEstimate(job, &model)
			if err != nil {
				return fmt.Errorf("error estimating carbon: %w", err)
			}
			if carbonEstimate <= bestCarbon && model.Accuracy >= bestAccuracy {
				bestCarbon = carbonEstimate
				bestAccuracy = model.Accuracy
				bestStartTime = startTime
				selectedModel = &model
			}
		}
	}
	if selectedModel == nil {
		return fmt.Errorf("no model found that meets the required accuracy, there is likely no model with an accuracy >= to the required accuracy")
	}
	job.Model = selectedModel
	job.StartTime = bestStartTime
	duration := max(selectedModel.MeanRunTime+selectedModel.StdDevRunTime*rand.NormFloat64(), 0)
	job.EndTime = job.StartTime.Add(time.Duration(duration) * time.Second)

	h.currTotalAccuracy += selectedModel.Accuracy
	h.processedJobs++
	return nil
}

func (h HybridSelection) HandleQueued(job *workload.Job) error {
	return nil
}

func (h HybridSelection) HandleRunning(job *workload.Job) error {
	return nil
}

func (h HybridSelection) String() string {
	return fmt.Sprintf("HybridSelection with required accuracy: %f", h.requiredAccuracy)
}
