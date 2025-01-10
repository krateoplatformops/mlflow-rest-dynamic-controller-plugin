package modelversion

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/krateoplatformops/mlflow-rest-dynamic-controller-plugin/internal/handlers"
	"github.com/krateoplatformops/mlflow-rest-dynamic-controller-plugin/internal/handlers/modelversion"
)

func GetModelVersion(opts handlers.HandlerOptions) handlers.Handler {
	return &handler{
		HandlerOptions: opts,
	}
}

var _ handlers.Handler = &handler{}

type handler struct {
	handlers.HandlerOptions
}

// @Summary Get metadata for a model version
// @Description Get metadata for a model version
// @ID get-model-version
// @Param name query string true "Name of the registered model"
// @Param version query string true "Model version number"
// @Produce json
// @Success 200 {object} ModelVersion
// @Router /2.0/mlflow/model-versions/get [get]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	version := r.URL.Query().Get("version")

	// Validate required parameters
	if name == "" || version == "" {
		http.Error(w, "missing required parameters: name and version", http.StatusBadRequest)
		return
	}

	log := h.Log.With(
		"operation", "/2.0/mlflow/model-versions/get",
		"method", "GET",
		"name", name,
		"version", version,
	)

	log.Debug("calling MLFlow Model Version API")

	url := fmt.Sprintf("%s/2.0/mlflow/model-versions/get?name=%s&version=%s", h.Server.String(), name, version)

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

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("reading response body", slog.Any("error", err))
		http.Error(w, "failed to read response", http.StatusInternalServerError)
		return
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		log.Error("MLFlow API error",
			"status", resp.StatusCode,
			"response", string(body),
		)
		http.Error(w, fmt.Sprintf("MLFlow API error: %s", string(body)), resp.StatusCode)
		return
	}

	var model modelversion.ModelVersionResponse
	if err := json.Unmarshal(body, &model); err != nil {
		log.Error("unmarshalling response", slog.Any("error", err))
		http.Error(w, "failed to parse response", http.StatusInternalServerError)
		return
	}

	// Set response headers
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if model.ModelVersion != nil {
		if err := json.NewEncoder(w).Encode(model.ModelVersion); err != nil {
			log.Error("encoding response", slog.Any("error", err))
			// Cannot write error to client at this point as headers are already sent
			return
		}
	}

	if model.ModelVersionCamelCase != nil {
		if err := json.NewEncoder(w).Encode(model.ModelVersionCamelCase); err != nil {
			log.Error("encoding response", slog.Any("error", err))
			return
		}
	}

	log.Debug("successfully processed request",
		"url", req.URL.String(),
		"status", http.StatusOK,
	)
}
