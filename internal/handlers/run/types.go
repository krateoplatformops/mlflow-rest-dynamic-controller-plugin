package run

// Run:
// type: object
// properties:
//   info:
// 	$ref: '#/components/schemas/RunInfo'
//   data:
// 	$ref: '#/components/schemas/RunData'
//   inputs:
// 	$ref: '#/components/schemas/RunInputs'
// RunData:
// type: object
// properties:
//   metrics:
// 	type: array
// 	items:
// 	  $ref: '#/components/schemas/Metric'
//   params:
// 	type: array
// 	items:
// 	  $ref: '#/components/schemas/Param'
//   tags:
// 	type: array
// 	items:
// 	  $ref: '#/components/schemas/RunTag'
// RunInfo:
// type: object
// properties:
//   run_id:
// 	type: string
// 	description: Unique identifier for the run.
//   run_uuid:
// 	type: string
// 	description: "[Deprecated, use run_id instead] Unique identifier for the run. This field will be removed in a future MLflow version."
//   run_name:
// 	type: string
// 	description: The name of the run.
//   experiment_id:
// 	type: string
// 	description: The experiment ID.
//   user_id:
// 	type: string
// 	description: User who initiated the run. This field is deprecated as of MLflow 1.0, and will be removed in a future MLflow release. Use ‘mlflow.user’ tag instead.
//   status:
// 	type: string
// 	description: Current status of the run.
//   start_time:
// 	type: integer
// 	format: int64
// 	description: Unix timestamp of when the run started in milliseconds.
//   end_time:
// 	type: integer
// 	format: int64
// 	description: Unix timestamp of when the run ended in milliseconds.
//   artifact_uri:
// 	type: string
// 	description: URI of the directory where artifacts should be uploaded.
//   lifecycle_stage:
// 	type: string
// 	description: "Current life cycle stage of the experiment : OneOf(“active”, “deleted”)"
// RunInputs:
// type: object
// properties:
//   dataset_inputs:
// 	type: array
// 	items:
// 	  $ref: '#/components/schemas/DatasetInput'
// RunTag:
// type: object
// properties:
//   key:
// 	type: string
// 	description: The tag key.
//   value:
// 	type: string
// 	description: The tag value.
// Metric:
// type: object
// properties:
//   key:
// 	type: string
// 	description: The metric key.
//   value:
// 	type: number
// 	format: float
// 	description: The metric value.
//   timestamp:
// 	type: integer
// 	format: int64
// 	description: The timestamp of the metric.
// Param:
// type: object
// properties:
//   key:
// 	type: string
// 	description: The parameter key.
//   value:
// 	type: string
// 	description: The parameter value.
// DatasetInput:
// type: object
// properties:
//   dataset_id:
// 	type: string
// 	description: The dataset ID.
//   dataset_name:
// 	type: string
// 	description: The dataset name.

type RunResponse struct {
	Run Run `json:"run"`
}

type Run struct {
	RunId          string    `json:"run_id"`
	RunUuid        string    `json:"run_uuid"`
	RunName        string    `json:"run_name"`
	ExperimentId   string    `json:"experiment_id"`
	UserId         string    `json:"user_id"`
	Status         string    `json:"status"`
	StartTime      int64     `json:"start_time"`
	EndTime        int64     `json:"end_time"`
	ArtifactUri    string    `json:"artifact_uri"`
	LifecycleStage string    `json:"lifecycle_stage"`
	Info           RunInfo   `json:"info"`
	Data           RunData   `json:"data"`
	Inputs         RunInputs `json:"inputs"`
}

type RunData struct {
	Metrics []Metric `json:"metrics"`
	Params  []Param  `json:"params"`
	Tags    []RunTag `json:"tags"`
}

type RunInfo struct {
	RunId          string `json:"run_id"`
	RunUuid        string `json:"run_uuid"`
	RunName        string `json:"run_name"`
	ExperimentId   string `json:"experiment_id"`
	UserId         string `json:"user_id"`
	Status         string `json:"status"`
	StartTime      int64  `json:"start_time"`
	EndTime        int64  `json:"end_time"`
	ArtifactUri    string `json:"artifact_uri"`
	LifecycleStage string `json:"lifecycle_stage"`
}

type RunInputs struct {
	DatasetInputs []DatasetInput `json:"dataset_inputs"`
}

type RunTag struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type Metric struct {
	Key       string  `json:"key"`
	Value     float64 `json:"value"`
	Timestamp int64   `json:"timestamp"`
}

type Param struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

type DatasetInput struct {
	DatasetId   string `json:"dataset_id"`
	DatasetName string `json:"dataset_name"`
}
