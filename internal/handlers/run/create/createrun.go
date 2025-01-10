package run

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/krateoplatformops/mlflow-rest-dynamic-controller-plugin/internal/handlers"
	"github.com/krateoplatformops/mlflow-rest-dynamic-controller-plugin/internal/handlers/run"
)

func CreateRun(opts handlers.HandlerOptions) handlers.Handler {
	return &handler{
		HandlerOptions: opts,
	}
}

var _ handlers.Handler = &handler{}

type handler struct {
	handlers.HandlerOptions
}

// @Summary Create a new run
// @Description Create a new run
// @ID create-run
// @Accept json
// @Produce json
// @Success 200 {object} Run
// @Router /2.0/mlflow/runs/create [post]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := h.Log.With(
		"operation", "createRun",
		"method", "POST",
	)

	// Read and validate request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Error("reading request body", slog.Any("error", err))
		http.Error(w, "failed to read request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Create API request
	url := h.Server.String() + "/2.0/mlflow/runs/create"
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(body))
	if err != nil {
		log.Error("creating request", slog.Any("error", err))
		http.Error(w, "failed to create request", http.StatusInternalServerError)
		return
	}

	// Forward authorization header if present
	if authHeader := r.Header.Get("Authorization"); authHeader != "" {
		req.Header.Set("Authorization", authHeader)
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute request
	resp, err := h.Client.Do(req)
	if err != nil {
		log.Error("calling MLFlow API", slog.Any("error", err))
		http.Error(w, "failed to call MLFlow API", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	// Read response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Error("reading response body", slog.Any("error", err))
		http.Error(w, "failed to read response", http.StatusInternalServerError)
		return
	}

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		log.Error("MLFlow API error",
			"status", resp.StatusCode,
			"response", string(respBody),
		)
		http.Error(w, fmt.Sprintf("MLFlow API error: %s", string(respBody)), resp.StatusCode)
		return
	}

	// Parse response
	var runResp run.RunResponse
	if err := json.Unmarshal(respBody, &runResp); err != nil {
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

	// Write response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(runResp.Run); err != nil {
		log.Error("encoding response", slog.Any("error", err))
		// Can't write error to client as headers are already sent
		return
	}

	log.Debug("successfully created run",
		"url", req.URL.String(),
		"status", http.StatusOK,
	)
}
