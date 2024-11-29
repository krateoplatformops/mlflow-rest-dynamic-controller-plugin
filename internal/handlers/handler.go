package handlers

import (
	"log/slog"
	"net/http"
	"net/url"
)

type HandlerOptions struct {
	Log    *slog.Logger
	Client *http.Client
	Server url.URL
}

type Handler interface {
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}
