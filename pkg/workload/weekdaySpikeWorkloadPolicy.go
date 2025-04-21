package workload

import (
	"fmt"
	"log"
	"math"
	"simulator/pkg/loader"
	"slices"
)

type WeekdaySpikeWorkload struct {
	PercentOn  float64 // Percent of jobs that are during the spike
	PercentOff float64 // Percent of jobs that are during the off-peak
}

func NewWeekdaySpikeWorkload(percentOn, percentOff float64) *WeekdaySpikeWorkload {
	return &WeekdaySpikeWorkload{
		PercentOn:  percentOn,
		PercentOff: percentOff,
	}
}

func (ws *WeekdaySpikeWorkload) GenerateWorkload(jobInfo JobMetadata) ([]*Job, error) {
	loader := loader.GetLoader()
	if loader == nil {
		panic("Loader is not initialized")
	}
	if loader.NumEntries() == 0 {
		panic("No entries in the loader")
	}
	numJobs := jobInfo.NumJobs
	numsOnSpike := int(math.Round(float64(numJobs) * ws.PercentOn))
	if numsOnSpike == 0 {
		return nil, fmt.Errorf("no jobs generated for the spike period")
	}
	if numsOnSpike > numJobs {
		numsOnSpike = numJobs
	}
	numOffSpike := numJobs - numsOnSpike

	jobList := make([]*Job, numJobs)
	startDate := loader.StartDate()
	endDate := loader.EndDate().Add(-jobInfo.DueTime)

	log.Printf("Putting %d jobs in the weekday spike period and %d jobs in the off-peak period\n", numsOnSpike, numOffSpike)
	for i := range numsOnSpike {
		startTime := getWeekdayInRange(startDate, endDate)
		jobList[i] = &Job{
			Model:     nil,
			StartTime: startTime,
			DueTime:   startTime.Add(jobInfo.DueTime),
			EndTime:   nil,
		}
	}
	for i := range numOffSpike {
		startTime := getWeekendInRange(startDate, endDate)
		jobList[numsOnSpike+i] = &Job{
			Model:     nil,
			StartTime: startTime,
			DueTime:   startTime.Add(jobInfo.DueTime),
			EndTime:   nil,
		}
	}
	slices.SortFunc(jobList, func(a, b *Job) int {
		return a.StartTime.Compare(b.StartTime)
	})
	return jobList, nil
}
