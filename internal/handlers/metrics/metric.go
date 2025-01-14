package metric

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/krateoplatformops/mlflow-rest-dynamic-controller-plugin/internal/handlers"
)

func GetMetric(opts handlers.HandlerOptions) handlers.Handler {
	return &handler{
		HandlerOptions: opts,
	}
}

var _ handlers.Handler = &handler{}

type handler struct {
	handlers.HandlerOptions
}

// @Summary Get metadata for an metric
// @Description Get metadata for an metric
// @ID get-metric
// @Param run_id query string true "ID of the associated metric"
// @Produce json
// @Success 200 {object} map[string]any
// @Router /2.0/mlflow/metrics/get-history [get]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	runID := r.URL.Query().Get("run_id")
	metricKey := r.URL.Query().Get("metric_key")
	if runID == "" || metricKey == "" {
		http.Error(w, "missing required parameter: run_id or metric_key", http.StatusBadRequest)
		return
	}

	log := h.Log.With(
		"operation", "/2.0/mlflow/metrics/get-history",
		"method", "GET",
		"run_id", runID,
		"metric_key", metricKey,
	)

	log.Debug("calling MLFlow metric history API")

	url := fmt.Sprintf("%s/2.0/mlflow/metrics/get-history?run_id=%s&metric_key=%s", h.Server.String(), runID, metricKey)

	log.Debug("url: ", slog.Any("url", url))
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

	log.Debug("calling MLFlow API", slog.Any("url", req.URL.String()))
	log.Debug("request headers", slog.Any("headers", req.Header))

	resp, err := h.Client.Do(req)
	if err != nil {
		log.Error("calling MLFlow API", slog.Any("error", err))
		http.Error(w, "failed to call MLFlow API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	log.Debug("response status", slog.Any("status", resp.StatusCode))
	log.Debug("response body", slog.Any("body", resp.Body))

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("reading response body", slog.Any("error", err))
		http.Error(w, "failed to read response", http.StatusInternalServerError)
		return
	}

	var metricResp map[string]any
	if resp.StatusCode != http.StatusOK {
		metricResp = map[string]any{
			"metrics": []any{},
		}

		if err := json.NewEncoder(w).Encode(metricResp); err != nil {
			log.Error("encoding response", slog.Any("error", err))
			http.Error(w, "failed to encode response", http.StatusInternalServerError)
			return
		}
		return
	}

	if err := json.Unmarshal(body, &metricResp); err != nil {
		log.Error("unmarshalling response", slog.Any("error", err))
		http.Error(w, "failed to parse response", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(metricResp); err != nil {
		log.Error("encoding response", slog.Any("error", err))
		http.Error(w, "failed to encode response", http.StatusInternalServerError)
		return
	}

	log.Debug("successfully processed request",
		"url", req.URL.String(),
		"status", http.StatusOK,
	)
}
