package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "repo-stat/api/docs"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, ping *usecase.Ping, fetch *usecase.Fetch) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info", NewGetInfoHandler(log, fetch))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
}
