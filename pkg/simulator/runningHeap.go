package simulator

import "simulator/pkg/workload"

type RunningHeap []*workload.Job

func (pq RunningHeap) Len() int {
	return len(pq)
}

func (pq RunningHeap) Less(i, j int) bool {
	// We want to return the lowest endTime first
	return pq[i].EndTime.Before(pq[j].EndTime)
}

func (pq RunningHeap) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *RunningHeap) Push(x any) {
	*pq = append(*pq, x.(*workload.Job))
}

func (pq *RunningHeap) Pop() any {
	old := *pq
	n := len(old)
	x := old[0]
	*pq = old[1:n]
	return x
}

func (pq *RunningHeap) Peek() *workload.Job {
	if pq.Len() == 0 {
		return nil
	}
	return (*pq)[0]
}
