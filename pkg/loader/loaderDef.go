package loader

import "time"

type LoaderInterface interface {
	// Private Methods
	loadFromFile() error

	// Public Methods
	String() string
	PrintAllData() error
	GetIndexByDate(date time.Time) (int, error)
	NumEntries() int
	StartDate() time.Time
	EndDate() time.Time
}

type Loader struct {
	filename   string
	startDate  time.Time
	numEntries int
	Data       []*DataPoint
}

type DataPoint struct {
	StartDate       time.Time `csv:"start_date"`
	CarbonIntensity float64   `csv:"generated_rate_kg_per_mwh"`
}
