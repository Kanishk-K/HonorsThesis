package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"simulator/pkg/directory"
	"simulator/pkg/loader"
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
	log.Println(modelDirectory.GetModelDefinition("gpt-3.5-turbo"))
}
