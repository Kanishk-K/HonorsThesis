package policies

import (
	"fmt"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"time"
)

func CarbonCalculate(start time.Time, end time.Time, model *directory.AIModelDefinition) float64 {
	// newTime is always less than or equal to the end time of the runningQueue
	loader := loader.GetLoader()
	if loader == nil {
		panic("loader not initialized")
	}
	if loader.NumEntries() == 0 {
		panic("loader has no data")
	}
	totalCarbon := 0.0
	currTime := start
	// Iterate until we run out of carbon data or the job end time
	for currTime.Before(end) && currTime.Before(loader.EndDate()) {
		carbonIdx, err := loader.GetIndexByDate(currTime)
		if err != nil {
			panic(fmt.Sprintf("error getting index by date: %v", err))
		}
		var nextTime time.Time
		if carbonIdx < loader.NumEntries()-1 && loader.Data[carbonIdx+1].StartDate.Before(end) {
			nextTime = loader.Data[carbonIdx+1].StartDate
		} else {
			// If carbonIdx is the last entry, this will find the heuristic from ending outside boundary
			// Otherwise, if the next entry is after the end time, we will use the end time
			// to calculate the carbon emission
			nextTime = end
		}
		// Calculate the time difference
		timeDelta := nextTime.Sub(currTime).Seconds() // in seconds
		if timeDelta <= 0 {
			panic(fmt.Sprintf("non-positive timedelta: %f", timeDelta))
		}
		// Calculate the carbon emission
		carbonRate := loader.Data[carbonIdx].CarbonIntensity // in kgCO2/MWh
		modelRate := model.EnergyUsage                       // in MW
		carbon := timeDelta * modelRate * 3.6e-9 * 1e3 * carbonRate
		totalCarbon += carbon // in gCO2

		currTime = nextTime
	}
	return totalCarbon
}
