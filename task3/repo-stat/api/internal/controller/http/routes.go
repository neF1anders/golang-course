package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, ping *usecase.Ping, fetch *usecase.Fetch) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info?url=<github_repo_url>", NewGetInfoHandler(log, fetch))
}
