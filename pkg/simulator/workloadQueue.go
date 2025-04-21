package simulator

import (
	"simulator/pkg/workload"
)

type WorkloadQueue []*workload.Job

func (wq WorkloadQueue) Len() int {
	return len(wq)
}

func (wq WorkloadQueue) Peak() *workload.Job {
	if len(wq) == 0 {
		return nil
	}
	return wq[0]
}

func (wq *WorkloadQueue) Pop() *workload.Job {
	if len(*wq) == 0 {
		return nil
	}
	job := (*wq)[0]
	*wq = (*wq)[1:]
	return job
}
