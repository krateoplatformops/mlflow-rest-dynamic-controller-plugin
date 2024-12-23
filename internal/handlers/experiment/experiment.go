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
	experiment_id := r.URL.Query().Get("experiment_id")

	log := h.Log.With(
		"Performing", "/2.0/mlflow/experiments/get [get]",
		"experiment_id", experiment_id)

	log.Debug("Calling MLFlow Experiment API")

	url := h.Server.String() + "/2.0/mlflow/experiments/get?experiment_id=" + experiment_id

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		h.Log.Error("creating request", slog.Any("error", err))
		w.Write([]byte(fmt.Sprint("Error: ", err)))
	}

	if r.Header.Get("Authorization") != "" {
		req.Header.Set("Authorization", r.Header.Get("Authorization"))
	}

	resp, err := h.Client.Do(req)
	if err != nil {
		h.Log.Error("calling MLFlow Experiment GET API", slog.Any("error", err))
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

				var experiment ExperimentResponse
				err = json.Unmarshal(body, &experiment)
				if err != nil {
					h.Log.Error("unmarshalling response body", slog.Any("error", err))
					fmt.Println("body", string(body))
					w.Write([]byte(fmt.Sprint("Error: ", err)))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				err = json.NewEncoder(w).Encode(experiment.Experiment)
				if err != nil {
					h.Log.Error("encoding response body", slog.Any("error", err))
					w.Write([]byte(fmt.Sprint("Error: ", err)))
				}

				log.Debug("Successfully called", slog.Any("URL", req.URL))
				return
			}
		}
	}

}
