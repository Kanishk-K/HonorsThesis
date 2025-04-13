package workload

import (
	"simulator/pkg/directory"
	"simulator/pkg/loader"
)

func GetModelWorkload(model *directory.AIModelDefinition) Workload {
	loader := loader.GetLoader()
	if loader == nil {
		panic("Loader is not initialized")
	}
	policy := model.WorkloadPolicy
	var workload PolicyInterface
	switch policy {
	case "random":
		workload = &RandomWorkload{}
	case "uniform":
		workload = &UniformWorkload{}
	default:
		panic("Unknown workload policy: " + policy)
	}

	jobs, err := workload.GenerateWorkload(model)
	if err != nil {
		panic("Failed to generate workload: " + err.Error())
	}
	return Workload{
		Model:  model,
		Policy: workload,
		Jobs:   jobs,
	}
}
