package policies

import (
	"fmt"
	"math"
	"math/rand"
	"simulator/pkg/directory"
	"simulator/pkg/workload"
	"sync"
	"time"
)

type ModelSelection struct {
	requiredAccuracy  float64
	currTotalAccuracy float64
	processedJobs     int
}

type ModelSelectionEstimate struct {
	CarbonEstimate float64
	Model          *directory.AIModelDefinition
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
	arrayLen := 0
	var wg sync.WaitGroup
	for _, model := range models {
		newAccuracy := (m.currTotalAccuracy + model.Accuracy) / float64(m.processedJobs+1)
		if newAccuracy >= m.requiredAccuracy {
			arrayLen++
		}
	}
	if arrayLen == 0 {
		return fmt.Errorf("no model found that meets the required accuracy, there is likely no model with an accuracy >= to the required accuracy")
	}
	wg.Add(arrayLen)
	array := make([]ModelSelectionEstimate, arrayLen)
	i := 0
	for _, model := range models {
		newAccuracy := (m.currTotalAccuracy + model.Accuracy) / float64(m.processedJobs+1)
		if newAccuracy >= m.requiredAccuracy {
			go func(index int, model *directory.AIModelDefinition) {
				defer wg.Done()
				carbonEstimate := FIFOCarbonEstimate(job, model)
				array[index] = ModelSelectionEstimate{
					CarbonEstimate: carbonEstimate,
					Model:          model,
				}
			}(i, &model)
			i += 1
		}
	}
	wg.Wait()
	bestCarbon := math.MaxFloat64
	bestAccuracy := math.MaxFloat64
	for _, estimate := range array {
		if estimate.CarbonEstimate <= bestCarbon && estimate.Model.Accuracy <= bestAccuracy {
			bestCarbon = estimate.CarbonEstimate
			bestAccuracy = estimate.Model.Accuracy
			selectedModel = estimate.Model
		}
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
