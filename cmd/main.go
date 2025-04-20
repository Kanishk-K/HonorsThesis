package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
	"simulator/pkg/simulator"
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
	wlq := simulator.CreateWorkloadQueue()
	if wlq == nil {
		log.Println("Failed to create workload queue")
		return
	}
	for _, job := range wlq {
		log.Println(*job)
	}
}
