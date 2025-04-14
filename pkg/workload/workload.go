package workload

import (
	"math"
	"math/rand"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"time"
)

func GetModelWorkload(model *directory.AIModelDefinition) Workload {
	loader := loader.GetLoader()
	if loader == nil {
		panic("Loader is not initialized")
	}
	policy := model.WorkloadPolicy
	var workload PolicyInterface
	switch policy {
	case "random":
		workload = &RandomWorkload{}
	case "uniform":
		workload = &UniformWorkload{}
	case "morningSpike":
		// Morning spike from 5 AM to 12 PM
		workload = NewMorningSpikeWorkload(0.8, 0.2)
	case "afternoonSpike":
		// Afternoon spike from 12 PM to 5 PM
		workload = NewAfternoonSpikeWorkload(0.8, 0.2)
	case "eveningSpike":
		// Evening spike from 5 PM to 12 AM
		workload = NewEveningSpikeWorkload(0.8, 0.2)
	case "nightSpike":
		// Night spike from 12 AM to 5 AM
		workload = NewNightSpikeWorkload(0.8, 0.2)
	case "weekdaySpike":
		// Assigns jobs to weekdays
		workload = NewWeekdaySpikeWorkload(0.8, 0.2)
	case "weekendSpike":
		// Assigns jobs to weekends
		workload = NewWeekendSpikeWorkload(1, 0)
	default:
		panic("Unknown workload policy: " + policy)
	}

	jobs, err := workload.GenerateWorkload(model)
	if err != nil {
		panic("Failed to generate workload: " + err.Error())
	}
	return Workload{
		Model:  model,
		Policy: workload,
		Jobs:   jobs,
	}
}

func getStartTimeInRange(startDate, endDate time.Time, startHour, endHour int) time.Time {
	// Generate a random time within the specified hour range
	duration := endDate.Sub(startDate)
	totalDays := int(duration.Hours() / 24)
	if endHour < startHour {
		endHour += 24
	}

	randomDay := rand.Intn(totalDays + 1)
	chosenDate := startDate.AddDate(0, 0, randomDay)

	randomHour := (rand.Intn(endHour-startHour) + startHour) % 24
	randomMinute := rand.Intn(60)
	randomSecond := rand.Intn(60)

	retVal := time.Date(chosenDate.Year(), chosenDate.Month(), chosenDate.Day(),
		randomHour, randomMinute, randomSecond, 0, chosenDate.Location())

	// Ensure the generated time is within the specified range
	if retVal.Before(startDate) {
		// If the generated time is before the start date, add days to make it valid
		daysToAdd := int(math.Ceil(startDate.Sub(retVal).Hours() / 24))
		retVal = retVal.AddDate(0, 0, daysToAdd)
	}
	if retVal.After(endDate) {
		// If the generated time is after the end date, subtract days to make it valid
		daysToSubtract := int(math.Ceil(retVal.Sub(endDate).Hours() / 24))
		retVal = retVal.AddDate(0, 0, -daysToSubtract)
	}
	return retVal
}

func getWeekdayInRange(startDate, endDate time.Time) time.Time {
	if endDate.Before(startDate) {
		panic("endDate must be after startDate")
	}

	// Calculate the duration between the two times
	duration := endDate.Sub(startDate)

	for {
		// Generate a random duration within the range
		randomOffset := time.Duration(rand.Int63n(int64(duration)))
		candidate := startDate.Add(randomOffset)

		// Check if it's a weekday (Monday to Friday)
		if candidate.Weekday() >= time.Monday && candidate.Weekday() <= time.Friday {
			return candidate
		}
	}
}

func getWeekendInRange(startDate, endDate time.Time) time.Time {
	if endDate.Before(startDate) {
		panic("endDate must be after startDate")
	}

	// Calculate the duration between the two times
	duration := endDate.Sub(startDate)

	for {
		// Generate a random duration within the range
		randomOffset := time.Duration(rand.Int63n(int64(duration)))
		candidate := startDate.Add(randomOffset)

		// Check if it's a weekend (Saturday or Sunday)
		if candidate.Weekday() == time.Saturday || candidate.Weekday() == time.Sunday {
			return candidate
		}
	}
}
