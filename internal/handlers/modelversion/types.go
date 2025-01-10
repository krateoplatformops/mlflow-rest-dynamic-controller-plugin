package modelversion

type ModelVersionResponse struct {
	ModelVersion ModelVersion `json:"model_version" description:"Model version metadata"`
}

// ModelVersionTag represents a tag associated with a model version
type ModelVersionTag struct {
	Key   string `json:"key" description:"The tag key"`
	Value string `json:"value" description:"The tag value"`
}

// ModelVersionStatus represents the current status of a model version
type ModelVersionStatus string

const (
	// PendingRegistration indicates registration request is pending
	PendingRegistration ModelVersionStatus = "PENDING_REGISTRATION"
	// FailedRegistration indicates registration request has failed
	FailedRegistration ModelVersionStatus = "FAILED_REGISTRATION"
	// Ready indicates model version is ready for use
	Ready ModelVersionStatus = "READY"
)

// ModelVersion represents metadata about a specific version of a model
type ModelVersion struct {
	Name                 string             `json:"name" description:"Unique name of the model"`
	Version              string             `json:"version" description:"Model's version number"`
	CreationTimestamp    int64              `json:"creation_timestamp" description:"Timestamp recorded when this model_version was created"`
	LastUpdatedTimestamp int64              `json:"last_updated_timestamp" description:"Timestamp recorded when metadata for this model_version was last updated"`
	UserID               string             `json:"user_id" description:"User that created this model_version"`
	CurrentStage         string             `json:"current_stage" description:"Current stage for this model_version"`
	Description          string             `json:"description" description:"Description of this model_version"`
	Source               string             `json:"source" description:"URI indicating the location of the source model artifacts"`
	RunID                string             `json:"run_id" description:"MLflow run ID used when creating model_version"`
	Status               ModelVersionStatus `json:"status" description:"Current status of the model version"`
	StatusMessage        string             `json:"status_message" description:"Details on current status, if pending or failed"`
	Tags                 []ModelVersionTag  `json:"tags" description:"Additional metadata key-value pairs"`
	RunLink              string             `json:"run_link" description:"Direct link to the run that generated this version"`
	Aliases              []string           `json:"aliases" description:"Aliases pointing to this model_version"`
}
