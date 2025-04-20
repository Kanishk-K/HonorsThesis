package workload

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"time"
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

func (ws *WeekdaySpikeWorkload) GenerateWorkload(model *directory.AIModelDefinition) ([]*Job, error) {
	loader := loader.GetLoader()
	if loader == nil {
		panic("Loader is not initialized")
	}
	if loader.NumEntries() == 0 {
		panic("No entries in the loader")
	}
	numsOnSpike := int(math.Round(float64(model.NumberOfRuns) * ws.PercentOn))
	if numsOnSpike == 0 {
		return nil, fmt.Errorf("no jobs generated for the spike period")
	}
	if numsOnSpike > model.NumberOfRuns {
		numsOnSpike = model.NumberOfRuns
	}
	numOffSpike := model.NumberOfRuns - numsOnSpike

	jobList := make([]*Job, model.NumberOfRuns)
	startDate := loader.StartDate()
	endDate := loader.EndDate().Add(-time.Duration(model.SLOThreshold) * time.Second)

	log.Printf("Putting %d jobs in the weekday spike period and %d jobs in the off-peak period\n", numsOnSpike, numOffSpike)
	for i := range numsOnSpike {
		startTime := getWeekdayInRange(startDate, endDate)
		duration := time.Duration(rand.NormFloat64()*model.StdDevRunTime+model.MeanRunTime) * time.Second
		if duration < 0 {
			duration = time.Second
		}
		jobList[i] = &Job{
			Model:     model,
			StartTime: startTime,
			DueTime:   startTime.Add(time.Duration(model.SLOThreshold) * time.Second),
			Duration:  duration,
		}
	}
	for i := range numOffSpike {
		startTime := getWeekendInRange(startDate, endDate)
		duration := time.Duration(rand.NormFloat64()*model.StdDevRunTime+model.MeanRunTime) * time.Second
		if duration < 0 {
			duration = time.Second
		}
		jobList[numsOnSpike+i] = &Job{
			Model:     model,
			StartTime: startTime,
			DueTime:   startTime.Add(time.Duration(model.SLOThreshold) * time.Second),
			Duration:  duration,
		}
	}
	return jobList, nil
}
