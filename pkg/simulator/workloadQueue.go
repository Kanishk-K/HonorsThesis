package simulator

import (
	"simulator/pkg/directory"
	"simulator/pkg/workload"
	"slices"
)

type WorkloadQueue []*workload.Job

func CreateWorkloadQueue() WorkloadQueue {
	// Get directory of all jobs in workloads
	directory := directory.FetchDirectory()
	if directory == nil {
		panic("Failed to create directory")
	}
	// Load in all AI Model Definitions
	workloadsSize := len(directory.GetModels())
	workloads := make([]workload.Workload, 0, workloadsSize)
	for _, value := range directory.GetModels() {
		workloads = append(workloads, workload.GetModelWorkload(&value))
	}
	var wq WorkloadQueue
	allocateSize := 0
	for _, wl := range workloads {
		allocateSize += len(wl.Jobs)
	}
	wq = make(WorkloadQueue, 0, allocateSize)
	for _, wl := range workloads {
		wq = append(wq, wl.Jobs...)
	}
	slices.SortFunc(wq, func(a, b *workload.Job) int {
		return a.StartTime.Compare(b.StartTime)
	})
	return wq
}

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
