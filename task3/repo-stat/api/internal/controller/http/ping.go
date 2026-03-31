package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

func NewPingHandler(log *slog.Logger, ping *usecase.Ping) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		processor_status, subscriber_status := ping.Execute(r.Context())
		services := []dto.ServiceStatus{
			{
				Name:   "processor",
				Status: string(processor_status),
			},
			{
				Name:   "subscriber",
				Status: string(subscriber_status),
			},
		}
		var status string = "degraded"
		if subscriber_status == processor_status &&
			processor_status == domain.PingStatusUp {
			status = "ok"
		}
		response := dto.PingResponse{
			Status:   status,
			Services: services,
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
		}
	}
}
