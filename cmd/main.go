package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"simulator/pkg/simulator"
	"simulator/pkg/simulator/policies"
	"simulator/pkg/workload"
	"strconv"
	"time"
)

func chooseRegion() string {
	switch os.Args[1] {
	case "CAISO":
		return "CAISO.csv"
	case "ERCOT":
		return "ERCOT.csv"
	case "MISO":
		return "MISO.csv"
	case "NYISO":
		return "NYISO.csv"
	default:
		panic("Invalid region specified. Please choose CAISO, ERCOT, MISO, or NYISO.")
	}
}

func chooseSLO() time.Duration {
	switch os.Args[2] {
	case "30min":
		return 30 * time.Minute
	case "1hr":
		return 1 * time.Hour
	case "6hr":
		return 6 * time.Hour
	case "12hr":
		return 12 * time.Hour
	case "24hr":
		return 24 * time.Hour
	default:
		panic("Invalid SLO specified. Please choose 1hr, 6hr, 12hr, or 24hr.")
	}
}

func chooseWorkload() string {
	switch os.Args[3] {
	case "random":
		return "random"
	case "uniform":
		return "uniform"
	case "morningSpike":
		return "morningSpike"
	case "afternoonSpike":
		return "afternoonSpike"
	case "eveningSpike":
		return "eveningSpike"
	case "nightSpike":
		return "nightSpike"
	case "weekdaySpike":
		return "weekdaySpike"
	case "weekendSpike":
		return "weekendSpike"
	default:
		panic("Invalid workload specified. Please choose random, uniform, or normal.")
	}
}

func choosePolicy() simulator.PolicyInterface {
	switch os.Args[4] {
	case "fifo":
		model, err := directory.FetchDirectory().GetModelDefinition(os.Args[5])
		if err != nil {
			panic("Error getting model definition")
		}
		return policies.NewFIFO(model)
	case "temporal":
		model, err := directory.FetchDirectory().GetModelDefinition(os.Args[5])
		if err != nil {
			panic("Error getting model definition")
		}
		return policies.NewTemporal(model, 0)
	case "modelSelection":
		accuracy, err := strconv.ParseFloat(os.Args[5], 64)
		if err != nil {
			panic("Error parsing accuracy")
		}
		return policies.NewModelSelection(accuracy)
	case "hybridSelection":
		accuracy, err := strconv.ParseFloat(os.Args[5], 64)
		if err != nil {
			panic("Error parsing accuracy")
		}
		return policies.NewHybridSelection(accuracy, 0)
	default:
		panic("Invalid policy specified. Please choose fifo, temporal, modelSelection, or hybridSelection.")
	}
}

func main() {
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	/*
		Load in carbon emission data
	*/
	dataPath := filepath.Join(currDir, "..", "data", "collected", chooseRegion())
	dataLoader := loader.NewLoader(dataPath)
	log.Println(dataLoader)

	/*
		Load in AI Model Definitions & Workload information
	*/
	modelPath := filepath.Join(currDir, "..", "cmd", "AIModels.json")
	_ = directory.NewDirectory(modelPath)

	/*
		Generate and load in workload information
	*/
	jobInfo := workload.NewJobInfo(chooseSLO(), 1000000, chooseWorkload())
	workload := workload.GetWorkload(jobInfo)

	/*
		Initialize the simulator
	*/
	simElement := simulator.NewSimulator(workload.Jobs, choosePolicy())
	if simElement == nil {
		log.Println("Simulator not initialized. Exiting.")
		return
	}
	log.Println(simElement)
	simElement.Begin()
	log.Println("Simulation complete.")
	log.Println(simElement)
}
