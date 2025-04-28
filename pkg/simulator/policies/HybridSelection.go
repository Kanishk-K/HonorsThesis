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

type HybridSelection struct {
	requiredAccuracy  float64
	safeguardSD       float64
	currTotalAccuracy float64
	processedJobs     int
}

type HybridSelectionEstimate struct {
	BestStartTime  time.Time
	CarbonEstimate float64
	Model          *directory.AIModelDefinition
}

func NewHybridSelection(requiredAccuracy float64, safeguardSD float64) *HybridSelection {
	return &HybridSelection{
		requiredAccuracy:  requiredAccuracy,
		safeguardSD:       safeguardSD,
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
	arrayLen := 0
	var wg sync.WaitGroup
	for _, model := range models {
		newAccuracy := (h.currTotalAccuracy + model.Accuracy) / float64(h.processedJobs+1)
		if newAccuracy >= h.requiredAccuracy {
			arrayLen++
		}
	}
	if arrayLen == 0 {
		return fmt.Errorf("no model found that meets the required accuracy, there is likely no model with an accuracy >= to the required accuracy")
	}
	wg.Add(arrayLen)
	array := make([]HybridSelectionEstimate, arrayLen)
	i := 0
	for _, model := range models {
		newAccuracy := (h.currTotalAccuracy + model.Accuracy) / float64(h.processedJobs+1)
		if newAccuracy >= h.requiredAccuracy {
			go func(index int, model *directory.AIModelDefinition) {
				defer wg.Done()
				bestTime, carbonEstimate, err := TemporalCarbonEstimate(job, model, h.safeguardSD)
				if err != nil {
					panic(err)
				}
				array[index] = HybridSelectionEstimate{
					BestStartTime:  bestTime,
					CarbonEstimate: carbonEstimate,
					Model:          model,
				}
			}(i, &model)
			i += 1
		}
	}
	wg.Wait()
	bestStartTime := job.StartTime
	bestCarbon := math.MaxFloat64
	bestAccuracy := math.MaxFloat64
	for _, estimate := range array {
		if estimate.CarbonEstimate <= bestCarbon && estimate.Model.Accuracy <= bestAccuracy {
			bestStartTime = estimate.BestStartTime
			bestCarbon = estimate.CarbonEstimate
			bestAccuracy = estimate.Model.Accuracy
			selectedModel = estimate.Model
		}
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
