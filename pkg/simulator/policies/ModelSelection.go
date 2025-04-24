package policies

import (
	"fmt"
	"math"
	"math/rand"
	"simulator/pkg/directory"
	"simulator/pkg/workload"
	"time"
)

type ModelSelection struct {
	requiredAccuracy  float64
	currTotalAccuracy float64
	processedJobs     int
}

func NewModelSelection(requiredAccuracy float64) *ModelSelection {
	return &ModelSelection{
		requiredAccuracy:  requiredAccuracy,
		currTotalAccuracy: 0,
		processedJobs:     0,
	}
}

func (m *ModelSelection) HandleIncoming(job *workload.Job) error {
	modelDirectory := directory.FetchDirectory()
	if modelDirectory == nil {
		return fmt.Errorf("model directory not initialized")
	}
	models := modelDirectory.GetModels()
	var selectedModel *directory.AIModelDefinition
	bestCarbon := math.MaxFloat64
	bestAccuracy := 0.0
	for _, model := range models {
		newAccuracy := (m.currTotalAccuracy + model.Accuracy) / float64(m.processedJobs+1)
		if newAccuracy >= m.requiredAccuracy {
			carbonEstimate, err := CarbonEstimate(job, model)
			if err != nil {
				return fmt.Errorf("error estimating carbon: %w", err)
			}
			if carbonEstimate <= bestCarbon && model.Accuracy >= bestAccuracy {
				bestCarbon = carbonEstimate
				bestAccuracy = model.Accuracy
				selectedModel = &model
			}
		}
	}
	if selectedModel == nil {
		return fmt.Errorf("no model found that meets the required accuracy, there is likely no model with an accuracy >= to the required accuracy")
	}
	job.Model = selectedModel
	duration := max(selectedModel.MeanRunTime+selectedModel.StdDevRunTime*rand.NormFloat64(), 0)
	job.EndTime = job.StartTime.Add(time.Duration(duration) * time.Second)

	m.currTotalAccuracy += selectedModel.Accuracy
	m.processedJobs++
	return nil
}

func (m ModelSelection) HandleQueued(job *workload.Job) error {
	return nil
}

func (m ModelSelection) HandleRunning(job *workload.Job) error {
	return nil
}

func (m ModelSelection) String() string {
	return fmt.Sprintf("Model Selection with required accuracy %f", m.requiredAccuracy)
}
