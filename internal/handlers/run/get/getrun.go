package run

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/krateoplatformops/mlflow-rest-dynamic-controller-plugin/internal/handlers"
	"github.com/krateoplatformops/mlflow-rest-dynamic-controller-plugin/internal/handlers/run"
)

func GetRun(opts handlers.HandlerOptions) handlers.Handler {
	return &handler{
		HandlerOptions: opts,
	}
}

var _ handlers.Handler = &handler{}

type handler struct {
	handlers.HandlerOptions
}

// @Summary Get metadata for a run
// @Description Get metadata for a run
// @ID get-run
// @Param run_id query string true "ID of the associated run"
// @Produce json
// @Success 200 {object} Run
// @Failure 400 {string} string "Bad Request"
// @Failure 404 {string} string "Not Found"
// @Failure 500 {string} string "Internal Server Error"
// @Router /2.0/mlflow/runs/get [get]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	runID := r.URL.Query().Get("run_id")
	if runID == "" {
		http.Error(w, "missing required parameter: run_id", http.StatusBadRequest)
		return
	}

	log := h.Log.With(
		"operation", "/2.0/mlflow/runs/get",
		"method", "GET",
		"run_id", runID,
	)

	log.Debug("calling MLFlow Run API")

	url := fmt.Sprintf("%s/2.0/mlflow/runs/get?run_id=%s", h.Server.String(), runID)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		log.Error("creating request", slog.Any("error", err))
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

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

	var runResp run.RunResponse
	if err := json.Unmarshal(body, &runResp); err != nil {
		log.Error("unmarshalling response", slog.Any("error", err))
		http.Error(w, "failed to parse response", http.StatusInternalServerError)
		return
	}

	// Map response fields
	runResp.Run.RunId = runResp.Run.Info.RunId
	runResp.Run.RunUuid = runResp.Run.Info.RunUuid
	runResp.Run.RunName = runResp.Run.Info.RunName
	runResp.Run.ExperimentId = runResp.Run.Info.ExperimentId
	runResp.Run.UserId = runResp.Run.Info.UserId
	runResp.Run.Status = runResp.Run.Info.Status

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(runResp.Run); err != nil {
		log.Error("encoding response", slog.Any("error", err))
		return
	}

	log.Debug("successfully processed request",
		"url", req.URL.String(),
		"status", http.StatusOK,
	)
}
