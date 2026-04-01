package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/adapter/processor"
	"repo-stat/api/internal/adapter/subscriber"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, subClient *subscriber.Client, procClient *processor.Client) {
	mux.Handle("GET /api/ping", NewPingHandler(log, subClient, procClient))
	mux.Handle("GET /api/repositories/info", NewRepoInfoHandler(log, procClient))
}
