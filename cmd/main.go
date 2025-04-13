package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"simulator/pkg/workload"
)

func main() {
	currDir, err := os.Getwd()
	if err != nil {
		fmt.Println("Error getting current directory:", err)
		return
	}
	dataPath := filepath.Join(currDir, "..", "data", "collected", "MISO.csv")
	dataLoader := loader.NewLoader(dataPath)
	log.Println(dataLoader)

	modelPath := filepath.Join(currDir, "..", "cmd", "AIModels.json")
	modelDirectory := directory.NewDirectory(modelPath)

	modelDef, err := modelDirectory.GetModelDefinition("gpt-3.5-turbo")
	if err != nil {
		fmt.Println("Error getting model definition:", err)
		return
	}

	workload := workload.GetModelWorkload(modelDef)
	for _, job := range workload.Jobs {
		fmt.Printf("Job: StartTime=%s, DueTime=%s, Duration=%s\n", job.StartTime, job.DueTime, job.Duration)
	}
}
