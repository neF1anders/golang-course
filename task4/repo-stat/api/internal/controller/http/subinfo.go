package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	_ "repo-stat/api/docs"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// SubInfo godoc
// @Summary Get subscriptions list
// @Description Returns the list of all subscribed repositories
// @Tags subscriptions
// @Produce json
// @Success 200 {object} dto.Slugs
// @Failure 500 {string} string "internal server error"
// @Router /subscriptions [get]
func NewSubInfoHandler(log *slog.Logger, retriever *usecase.RetrieverUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		slugs, err := retriever.SubInfo(r.Context())
		if err != nil {
			log.Error("failed to list subscriptions", "error", err)
			writeGRPCError(w, err)
			return
		}
		var slugs_dto dto.Slugs
		for _, el := range slugs {
			slugs_dto.Slugs = append(slugs_dto.Slugs, dto.Slug{
				Owner: el.Owner,
				Repo:  el.Repo,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(slugs_dto); err != nil {
			log.Error("failed to write list response", "error", err)
		}
	}
}
