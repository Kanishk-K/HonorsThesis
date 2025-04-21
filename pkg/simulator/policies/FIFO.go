package policies

import (
	"fmt"
	"math/rand"
	"simulator/pkg/directory"
	"simulator/pkg/workload"
	"time"
)

type FIFO struct {
	aiModel *directory.AIModelDefinition
}

func NewFIFO(aiModel *directory.AIModelDefinition) FIFO {
	return FIFO{
		aiModel: aiModel,
	}
}

func (f FIFO) HandleIncoming(job *workload.Job) error {
	job.Model = f.aiModel
	// Generate the duration of the job
	duration := max(f.aiModel.MeanRunTime+f.aiModel.StdDevRunTime*rand.NormFloat64(), 0)
	job.EndTime = job.StartTime.Add(time.Duration(duration) * time.Second)
	return nil
}

func (f FIFO) HandleQueued(job *workload.Job) error {
	return nil
}

func (f FIFO) HandleRunning(job *workload.Job) error {
	return nil
}

func (f FIFO) String() string {
	return fmt.Sprintf("FIFO with %s", f.aiModel.ModelName)
}
