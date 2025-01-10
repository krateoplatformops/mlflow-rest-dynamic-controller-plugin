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

func CreateModelVersion(opts handlers.HandlerOptions) handlers.Handler {
	return &handler{
		HandlerOptions: opts,
	}
}

var _ handlers.Handler = &handler{}

type handler struct {
	handlers.HandlerOptions
}

// @Summary Create a new model version
// @Description Create a new model version
// @ID create-model-version
// @Accept json
// @Produce json
// @Success 200 {object} map[string]any
// @Router /2.0/mlflow/model-versions/create [post]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := h.Log.With(
		"operation", "/2.0/mlflow/model-versions/create",
		"method", "POST",
	)

	log.Debug("calling MLFlow Model Version API")

	url := h.Server.String() + "/2.0/mlflow/model-versions/create"

	req, err := http.NewRequest("POST", url, r.Body)
	if err != nil {
		h.Log.Error("creating request", "error", err)
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	if r.Header.Get("Authorization") != "" {
		req.Header.Set("Authorization", r.Header.Get("Authorization"))
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := h.Client.Do(req)
	if err != nil {
		h.Log.Error("calling MLFlow Model Version API", "error", err)
		http.Error(w, "failed to call MLFlow Model Version API", http.StatusInternalServerError)
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
