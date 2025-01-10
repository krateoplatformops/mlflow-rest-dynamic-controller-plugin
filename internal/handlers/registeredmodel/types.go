package registeredmodel

type RegisteredModelResponse struct {
	RegisteredModel          map[string]any `json:"registered_model" description:"Registered model metadata."`
	RegisteredModelCamelCase map[string]any `json:"registeredModel" description:"Registered model metadata."`
}
