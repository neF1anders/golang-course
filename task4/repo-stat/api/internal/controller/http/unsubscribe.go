package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// Unsubscribe godoc
// @Summary Unsubscribe from a repository
// @Description Remove an existing subscription
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param owner query string true "GitHub repo owner"
// @Param repo query string true "GitHub repo name"
// @Success 204 {object} map[string]string
// @Failure 400 {string} string "invalid request body or missing fields"
// @Failure 404 {string} string "subscription not found"
// @Failure 500 {string} string "internal server error"
// @Router /subscriptions [delete]
func NewUnsubscribeHandler(log *slog.Logger, retriever *usecase.RetrieverUseCase) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		owner := r.URL.Query().Get("owner")
		repo := r.URL.Query().Get("repo")
		if owner == "" || repo == "" {
			http.Error(w, "owner and repo are required", http.StatusBadRequest)
			return
		}
		req := dto.Slug{
			Owner: owner,
			Repo:  repo,
		}
		defer func() {
			err := r.Body.Close()
			if err != nil {
				log.Printf("Error during closing responce body in fetcher: %v", err)
			}
		}()
		err := retriever.Unsubscribe(r.Context(), domain.Slug{
			Owner: req.Owner,
			Repo:  req.Repo,
		})
		if err != nil {
			writeGRPCError(w, err)
			log.Error("failed to unsubscribe", "owner", req.Owner, "repo", req.Repo, "error", err)
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusNoContent)
			return
		}
		err = json.NewEncoder(w).Encode(map[string]string{
			"status": "unsubscribed",
			"owner":  req.Owner,
			"repo":   req.Repo,
		})
		if err != nil {
			log.Error("failed to encode a response", "error", err)
		}
	}
}
