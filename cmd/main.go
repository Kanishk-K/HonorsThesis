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
	"time"
)

func main() {
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	/*
		Load in carbon emission data
	*/
	dataPath := filepath.Join(currDir, "..", "data", "collected", "CAISO.csv")
	dataLoader := loader.NewLoader(dataPath)
	log.Println(dataLoader)

	/*
		Load in AI Model Definitions & Workload information
	*/
	modelPath := filepath.Join(currDir, "..", "cmd", "AIModels.json")
	modelDirectory := directory.NewDirectory(modelPath)

	/*
		Generate and load in workload information
	*/
	jobInfo := workload.NewJobInfo(6*time.Hour, 1000000, "random")
	workload := workload.GetWorkload(jobInfo)

	/*
		Initialize the simulator
	*/
	model, err := modelDirectory.GetModelDefinition("small")
	if err != nil {
		log.Println("Error getting model definition:", err)
		return
	}
	schedulingPolicy := policies.NewTemporal(model)
	simElement := simulator.NewSimulator(workload.Jobs, schedulingPolicy)
	if simElement == nil {
		log.Println("Simulator not initialized. Exiting.")
		return
	}
	log.Println(simElement)
	simElement.Begin()
	log.Println("Simulation complete.")
	log.Println(simElement)
}
