package directory

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
)

var lock = &sync.Mutex{}
var singleton *Directory

func NewDirectory(filename string) *Directory {
	if singleton == nil {
		lock.Lock()
		defer lock.Unlock()
		if singleton == nil {
			log.Println("Initializing Directory with filename:", filename)
			singleton = &Directory{
				filename: filename,
				models:   make(map[string]AIModelDefinition),
			}
			err := singleton.loadFromFile()
			if err != nil {
				log.Printf("Error loading from file: %v", err)
				return nil
			}
		} else {
			log.Println("Directory already initialized")
		}
	} else {
		log.Println("Directory already initialized")
	}
	return singleton
}

func (d *Directory) loadFromFile() error {
	// Load models from a JSON file
	jsonFile, err := os.Open(d.filename)
	if err != nil {
		return err
	}
	defer jsonFile.Close()
	data, err := io.ReadAll(jsonFile)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, &d.models)
	if err != nil {
		log.Printf("Error unmarshalling JSON: %v", err)
		return err
	}
	log.Println("Loaded models from file.")
	return nil
}

func (d *Directory) String() string {
	return "Directory for filename: " + d.filename
}

func (d *Directory) GetModelDefinition(modelName string) (*AIModelDefinition, error) {
	if model, exists := d.models[modelName]; exists {
		return &model, nil
	} else {
		return nil, fmt.Errorf("model %s not found", modelName)
	}
}
