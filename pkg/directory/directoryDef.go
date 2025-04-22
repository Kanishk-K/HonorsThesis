package directory

type DirectoryInterface interface {
	// Private methods
	loadFromFile() error
	// Public methods
	String() string
	GetModelDefinition(modelName string) (*AIModelDefinition, error)
	GetModels() map[string]AIModelDefinition
}

type Directory struct {
	filename string
	models   map[string]AIModelDefinition
}

type AIModelDefinition struct {
	ModelName     string  `json:"model_name"`
	MeanRunTime   float64 `json:"mean_run_time"`    // in seconds
	StdDevRunTime float64 `json:"std_dev_run_time"` // in seconds
	EnergyUsage   float64 `json:"energy_usage"`     // in MW
	SLOThreshold  float64 `json:"slo_threshold"`    // in seconds
}
