package simulator

import "simulator/pkg/workload"

type PriorityQueue []*workload.Job

func (pq PriorityQueue) Len() int {
	return len(pq)
}

func (pq PriorityQueue) Less(i, j int) bool {
	// We want to return the lowest (starttime + duration) first
	return pq[i].StartTime.Add(*pq[i].Duration).Before(pq[j].StartTime.Add(*pq[j].Duration))
}

func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
}

func (pq *PriorityQueue) Push(x any) {
	*pq = append(*pq, x.(*workload.Job))
}

func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	x := old[n-1]
	*pq = old[0 : n-1]
	return x
}

func (pq *PriorityQueue) Peek() *workload.Job {
	if pq.Len() == 0 {
		return nil
	}
	return (*pq)[0]
}
