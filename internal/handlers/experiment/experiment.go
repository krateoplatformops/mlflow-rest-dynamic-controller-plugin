package experiment

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/krateoplatformops/mlflow-rest-dynamic-controller-plugin/internal/handlers"
)

func GetExperiment(opts handlers.HandlerOptions) handlers.Handler {
	return &handler{
		HandlerOptions: opts,
	}
}

var _ handlers.Handler = &handler{}

type handler struct {
	handlers.HandlerOptions
}

// @Summary Get metadata for an experiment
// @Description Get metadata for an experiment
// @ID get-experiment
// @Param experiment_id query string true "ID of the associated experiment"
// @Produce json
// @Success 200 {object} Experiment
// @Router /2.0/mlflow/experiments/get [get]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	experimentID := r.URL.Query().Get("experiment_id")
	if experimentID == "" {
		http.Error(w, "missing required parameter: experiment_id", http.StatusBadRequest)
		return
	}

	log := h.Log.With(
		"operation", "/2.0/mlflow/experiments/get",
		"method", "GET",
		"experiment_id", experimentID,
	)

	log.Debug("calling MLFlow Experiment API")

	url := fmt.Sprintf("%s/2.0/mlflow/experiments/get?experiment_id=%s", h.Server.String(), experimentID)

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
			slog.Int("status", resp.StatusCode),
			slog.String("response", string(body)),
		)
		http.Error(w, fmt.Sprintf("MLFlow API error: %s", string(body)), resp.StatusCode)
		return
	}

	var experimentResp ExperimentResponse
	if err := json.Unmarshal(body, &experimentResp); err != nil {
		log.Error("unmarshalling response", slog.Any("error", err))
		http.Error(w, "failed to parse response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(experimentResp.Experiment); err != nil {
		log.Error("encoding response", slog.Any("error", err))
		return
	}

	log.Debug("successfully processed request",
		"url", req.URL.String(),
		"status", http.StatusOK,
	)
}
