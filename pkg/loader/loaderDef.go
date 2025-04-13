package loader

import "time"

type LoaderInterface interface {
	// Private Methods
	loadFromFile() error

	// Public Methods
	String() string
	PrintAllData() error
	GetIntentsityByDate(date time.Time) (float64, error)
	GetNumEntries() int
	GetStartDate() time.Time
	GetEndDate() time.Time
}

type Loader struct {
	filename   string
	data       []*DataPoint
	startDate  time.Time
	numEntries int
}

type DataPoint struct {
	StartDate       time.Time `csv:"start_date"`
	CarbonIntensity float64   `csv:"generated_rate_kg_per_mwh"`
}
