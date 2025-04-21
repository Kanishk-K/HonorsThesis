package loader

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocarina/gocsv"
)

var lock = &sync.Mutex{}
var singleton *Loader

func NewLoader(filename string) *Loader {
	if singleton == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleton == nil {
			singleton = &Loader{
				filename: filename,
			}
			err := singleton.loadFromFile()
			if err != nil {
				log.Println("Error loading from file:", err)
				singleton = nil
				return nil
			}
		} else {
			log.Println("Loader already initialized")
		}
	} else {
		log.Println("Loader already initialized")
	}
	return singleton
}

func GetLoader() *Loader {
	if singleton == nil {
		log.Println("Loader not initialized")
		return nil
	}
	return singleton
}

func (l *Loader) loadFromFile() error {
	// Implement your file loading logic here
	log.Println("Loading from file:", l.filename)
	dataFile, err := os.OpenFile(l.filename, os.O_RDONLY, 0644)
	if err != nil {
		log.Println("Error opening file:", err)
		return err
	}
	defer dataFile.Close()
	if err := gocsv.UnmarshalFile(dataFile, &l.Data); err != nil {
		log.Println("Error unmarshalling file:", err)
		return err
	}
	l.numEntries = len(l.Data)
	if l.numEntries == 0 {
		log.Println("No data found in file")
		return fmt.Errorf("no data found in file")
	}
	l.startDate = l.Data[0].StartDate
	log.Println("Data loaded successfully, number of entries:", l.numEntries)
	return nil
}

func (l *Loader) String() string {
	return "Loader for filename: " + l.filename
}

func (l *Loader) PrintAllData() error {
	if l.numEntries == 0 {
		return fmt.Errorf("no data to print")
	}
	for _, dataPoint := range l.Data {
		fmt.Printf("StartDate: %s, CarbonIntensity: %f\n", dataPoint.StartDate.Format("2006-01-02 03:04 PM"), dataPoint.CarbonIntensity)
	}
	return nil
}

func (l *Loader) GetIndexByDate(date time.Time) (int, error) {
	// Check if the loader is initialized
	if l.numEntries == 0 {
		return -1, fmt.Errorf("no data available")
	}
	// If the date is before the start date or after the last date in the data, return an error
	// The last date is the start date of the last entry + 5 minutes
	if date.Before(l.startDate) || date.After(l.Data[l.numEntries-1].StartDate.Add(5*time.Minute)) {
		return -1, fmt.Errorf("date out of range")
	}
	// Entries are sorted by date we can binary search for the date
	// We are NOT guarenteed that every date is present in the data
	// So we need to find the closest date
	// If the exact date is not found, use the closest PREVIOUS date
	left, right := 0, l.numEntries-1
	for left <= right {
		mid := (left + right) / 2
		if l.Data[mid].StartDate.Equal(date) {
			return mid, nil
		} else if l.Data[mid].StartDate.Before(date) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	// If we reach here, the date is not found, return the closest previous date
	if right < 0 {
		return -1, fmt.Errorf("no data found for date")
	}
	return right, nil
}

func (l *Loader) NumEntries() int {
	return l.numEntries
}

func (l *Loader) StartDate() time.Time {
	return l.startDate
}

func (l *Loader) EndDate() time.Time {
	if l.numEntries == 0 {
		return time.Time{}
	}
	return l.Data[l.numEntries-1].StartDate.Add(5 * time.Minute)
}
