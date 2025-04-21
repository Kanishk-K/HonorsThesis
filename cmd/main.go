package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
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
	dataPath := filepath.Join(currDir, "..", "data", "collected", "MISO.csv")
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
	jobInfo := workload.NewJobInfo(5*time.Minute, 100, "random")
	workload := workload.GetWorkload(jobInfo)
	for _, job := range workload.Jobs {
		fmt.Printf("StartTime: %s, DueTime: %s\n", job.StartTime.Format(time.RFC3339), job.DueTime.Format(time.RFC3339))
	}
}
