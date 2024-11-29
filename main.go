package main

import (
	"context"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"sync/atomic"
	"syscall"
	"time"

	_ "github.com/krateoplatformops/github-rest-dynamic-controller-plugin/docs"
	"github.com/krateoplatformops/github-rest-dynamic-controller-plugin/internal/env"
	"github.com/krateoplatformops/github-rest-dynamic-controller-plugin/internal/handlers"
	"github.com/krateoplatformops/github-rest-dynamic-controller-plugin/internal/handlers/experiment"
	httpSwagger "github.com/swaggo/http-swagger"
)

var (
	serviceName = "mlflow-rest-dynamic-controller-plugin"
)

// @title           MLFlow Plugin API for Krateo Operator Generator (KOG)
// @version         1.0
// @description     Simple wrapper around MLFlow API for Krateo Operator Generator (KOG)
// @termsOfService  http://swagger.io/terms/

// @contact.name   Krateo Support
// @contact.url    https://krateo.io
// @contact.email  contact@krateoplatformops.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host			localhost:8080
// @BasePath		/
// @schemes 	 	http

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	debugOn := flag.Bool("debug", env.Bool("PLUGIN_DEBUG", true), "dump verbose output")
	port := flag.Int("port", env.Int("PLUGIN_PORT", 8081), "port to listen on")
	serverStr := flag.String("server", env.String("MLFLOW_SERVER", "http://localhost:5000"), "MLFlow server URL")

	flag.Parse()

	mux := http.NewServeMux()

	lopts := &slog.HandlerOptions{Level: slog.LevelInfo}
	if *debugOn {
		lopts = &slog.HandlerOptions{Level: slog.LevelDebug}
	}

	// log := *slog.New(slog.NewJSONHandler(os.Stdout,
	// 	lopts))
	log := slog.New(slog.NewTextHandler(os.Stdout, lopts))

	log = log.With("service", serviceName)

	serverUrl, err := url.Parse(*serverStr)
	if err != nil {
		log.Error("parsing server URL", slog.Any("error", err))
	}

	opts := handlers.HandlerOptions{
		Log:    log,
		Client: http.DefaultClient,
		Server: *serverUrl,
	}

	healthy := int32(0)

	mux.Handle("GET /2.0/mlflow/experiments/get", experiment.GetExperiment(opts))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", *port),
		Handler:      mux,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 50 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	ctx, stop := signal.NotifyContext(context.Background(), []os.Signal{
		os.Interrupt,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGKILL,
		syscall.SIGHUP,
		syscall.SIGQUIT,
	}...)
	defer stop()

	go func() {
		atomic.StoreInt32(&healthy, 1)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("could not listen on %s - %v", server.Addr, err)
		}
	}()

	// Listen for the interrupt signal.
	log.Info("server is ready to handle requests", slog.Any("port", *port))
	<-ctx.Done()

	// Restore default behavior on the interrupt signal and notify user of shutdown.
	stop()
	log.Info("server is shutting down gracefully, press Ctrl+C again to force")
	atomic.StoreInt32(&healthy, 0)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	server.SetKeepAlivesEnabled(false)
	if err := server.Shutdown(ctx); err != nil {
		log.Error("server forced to shutdown", slog.Any("error", err))
	}

	log.Info("server gracefully stopped")
}
