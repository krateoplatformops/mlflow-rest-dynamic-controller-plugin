package registeredmodel

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/krateoplatformops/mlflow-rest-dynamic-controller-plugin/internal/handlers"
)

func GetRegisteredModel(opts handlers.HandlerOptions) handlers.Handler {
	return &handler{
		HandlerOptions: opts,
	}
}

var _ handlers.Handler = &handler{}

type handler struct {
	handlers.HandlerOptions
}

// @Summary Get metadata for a registered model
// @Description Get metadata for a registered model
// @ID get-registered-model
// @Param name query string true "Registered model unique name identifier"
// @Produce json
// @Success 200 {object} RegisteredModel
// @Router /2.0/mlflow/registered-models/get [get]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	if name == "" {
		http.Error(w, "missing required parameter: name", http.StatusBadRequest)
		return
	}

	log := h.Log.With(
		"operation", "/2.0/mlflow/registered-models/get",
		"method", "GET",
		"name", name,
	)

	log.Debug("calling MLFlow Registered Model API")

	url := fmt.Sprintf("%s/2.0/mlflow/registered-models/get?name=%s", h.Server.String(), name)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error("creating request", slog.Any("error", err))
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	// Forward authorization header if present
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		log.Error("calling MLFlow API", slog.Any("error", err))
		http.Error(w, "failed to call MLFlow API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("reading response body", slog.Any("error", err))
		http.Error(w, "failed to read response", http.StatusInternalServerError)
		return
	}

	if resp.StatusCode != http.StatusOK {
		log.Error("MLFlow API error",
			"status", resp.StatusCode,
			"response", string(body),
		)
		http.Error(w, fmt.Sprintf("MLFlow API error: %s", string(body)), resp.StatusCode)
		return
	}

	var model RegisteredModelResponse
	if err := json.Unmarshal(body, &model); err != nil {
		log.Error("unmarshalling response", slog.Any("error", err))
		http.Error(w, "failed to parse response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(model.RegisteredModel); err != nil {
		log.Error("encoding response", slog.Any("error", err))
		return
	}

	log.Debug("successfully processed request",
		"url", req.URL.String(),
		"status", http.StatusOK,
	)
}
