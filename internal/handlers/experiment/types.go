package experiment

// Experiment:
// type: object
// properties:
//
//	  experiment_id:
//		type: string
//		description: Unique identifier for the experiment.
//	  name:
//		type: string
//		description: Human readable name that identifies the experiment.
//	  artifact_location:
//		type: string
//		description: Location where artifacts for the experiment are stored.
//	  lifecycle_stage:
//		type: string
//		description: "Current life cycle stage of the experiment: “active” or “deleted”."
//	  last_update_time:
//		type: integer
//		format: int64
//		description: Last update time.
//	  creation_time:
//		type: integer
//		format: int64
//		description: Creation time.
//	  tags:
//		type: array
//		items:
//		  $ref: '#/components/schemas/ExperimentTag'
//
// ExperimentTag:
// type: object
// properties:
//
//	  key:
//		type: string
//		description: The tag key.
//	  value:
//		type: string
//		description: The tag value.
type ExperimentResponse struct {
	Experiment Experiment `json:"experiment"`
}

type Experiment struct {
	ExperimentID     string          `json:"experiment_id"`
	Name             string          `json:"name"`
	ArtifactLocation string          `json:"artifact_location"`
	LifecycleStage   string          `json:"lifecycle_stage"`
	LastUpdateTime   int64           `json:"last_update_time"`
	CreationTime     int64           `json:"creation_time"`
	Tags             []ExperimentTag `json:"tags"`
}

type ExperimentTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
