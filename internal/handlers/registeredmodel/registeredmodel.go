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

	log := h.Log.With(
		"Performing", "/2.0/mlflow/registered-models/get [GET]",
		"name", name)

	log.Debug("Calling MLFlow Experiment API")

	url := h.Server.String() + "/2.0/mlflow/registered-models/get?name=" + name

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

				var model RegisteredModelResponse
				err = json.Unmarshal(body, &model)
				if err != nil {
					h.Log.Error("unmarshalling response body", slog.Any("error", err))
					w.Write([]byte(fmt.Sprint("Error: ", err)))
				}

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				err = json.NewEncoder(w).Encode(model.RegisteredModel)
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
