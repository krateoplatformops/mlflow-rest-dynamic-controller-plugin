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

func CreateRun(opts handlers.HandlerOptions) handlers.Handler {
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
// @Router /2.0/mlflow/runs/create [post]
func (h *handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log := h.Log.With(
		"Performing", "/2.0/mlflow/runs/create [POST]")

	log.Debug("Calling MLFlow Run API")

	url := h.Server.String() + "/2.0/mlflow/runs/create"

	req, err := http.NewRequest("POST", url, r.Body)
	if err != nil {
		h.Log.Error("creating request", slog.Any("error", err))
		w.Write([]byte(fmt.Sprint("Error: ", err)))
	}
	if r.Header.Get("Authorization") != "" {
		req.Header.Set("Authorization", r.Header.Get("Authorization"))
	}

	// set application/json content type
	req.Header.Set("Content-Type", "application/json")
	resp, err := h.Client.Do(req)
	if err != nil {
		h.Log.Error("calling MLFlow Run POST API", slog.Any("error", err))
		w.Write([]byte(fmt.Sprint("Error: ", err)))
	}

	if resp != nil {
		// read response body
		if resp.StatusCode == http.StatusOK {
			if resp.Body != nil {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					h.Log.Error("reading response body", slog.Any("error", err))
					w.Write([]byte(fmt.Sprint("Error: ", err)))
				}

				var run run.RunResponse
				err = json.Unmarshal(body, &run)
				if err != nil {
					h.Log.Error("unmarshalling response body", slog.Any("error", err))
					w.Write([]byte(fmt.Sprint("Error: ", err)))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				run.Run.RunId = run.Run.Info.RunId
				run.Run.RunUuid = run.Run.Info.RunUuid
				run.Run.RunName = run.Run.Info.RunName
				run.Run.ExperimentId = run.Run.Info.ExperimentId
				run.Run.UserId = run.Run.Info.UserId
				run.Run.Status = run.Run.Info.Status
				err = json.NewEncoder(w).Encode(run.Run)
				if err != nil {
					h.Log.Error("encoding response body", slog.Any("error", err))
					w.Write([]byte(fmt.Sprint("Error: ", err)))
				}

				log.Debug("Successfully called", slog.Any("URL", req.URL))
				return
			}
		} else {
			w.WriteHeader(resp.StatusCode)
			w.Write([]byte(fmt.Sprint("Error: ", resp.Status)))
		}
	}

}
