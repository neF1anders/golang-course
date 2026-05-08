package http

import (
	"log/slog"
	"net/http"
	"repo-stat/api/internal/usecase"

	httpSwagger "github.com/swaggo/http-swagger"

	_ "repo-stat/api/docs"
)

func AddRoutes(mux *http.ServeMux, log *slog.Logger, ping *usecase.Ping, fetch *usecase.Fetch, retriever *usecase.RetrieverUseCase) {
	mux.Handle("GET /api/ping", NewPingHandler(log, ping))
	mux.Handle("GET /api/repositories/info", NewGetInfoHandler(log, fetch))
	mux.Handle("/swagger/", httpSwagger.WrapHandler)
	mux.Handle("GET /subscriptions", NewSubInfoHandler(log, retriever))
	mux.Handle("POST /subscriptions", NewSubscribeHandler(log, retriever))
	mux.Handle("DELETE /subscriptions", NewUnsubscribeHandler(log, retriever))
	mux.Handle("GET /subscriptions/info", NewGetSubInfoHandler(log, fetch))
}
