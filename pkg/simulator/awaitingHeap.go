package simulator

import "simulator/pkg/workload"

type AwaitingHeap []*workload.Job

func (pq AwaitingHeap) Len() int {
	return len(pq)
}

func (pq AwaitingHeap) Less(i, j int) bool {
	// We want to return the lowest starttime first
	return pq[i].StartTime.Before(pq[j].StartTime)
}

func (pq AwaitingHeap) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *AwaitingHeap) Push(x any) {
	*pq = append(*pq, x.(*workload.Job))
}

func (pq *AwaitingHeap) Pop() any {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func (pq *AwaitingHeap) Peek() *workload.Job {
	if pq.Len() == 0 {
		return nil
	}
	return (*pq)[0]
}
