package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	_ "repo-stat/api/docs"
	"repo-stat/api/internal/domain"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
)

// Subscribe godoc
// @Summary Subscribe to a repository
// @Description Add a new subscription to track a GitHub repository
// @Tags subscriptions
// @Accept json
// @Produce json
// @Param owner query string true "GitHub repo owner"
// @Param repo query string true "GitHub repo name"
// @Success 200 {object} map[string]string
// @Failure 400 {string} string "invalid request body or missing fields"
// @Failure 404 {string} string "repository not found on GitHub"
// @Failure 409 {string} string "subscription already exists"
// @Failure 500 {string} string "internal server error"
// @Router /subscriptions [post]
func NewSubscribeHandler(log *slog.Logger, retriever *usecase.RetrieverUseCase) http.HandlerFunc {
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
				log.Error("Error during closing responce body in fetcher", "error", err)
			}
		}()
		err := retriever.Subscribe(r.Context(), domain.Slug{
			Owner: req.Owner,
			Repo:  req.Repo,
		})
		if err != nil {
			writeGRPCError(w, err)
			log.Error("failed to subscribe", "owner", req.Owner, "repo", req.Repo, "error", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		err = json.NewEncoder(w).Encode(map[string]string{
			"status": "subscribed",
			"owner":  req.Owner,
			"repo":   req.Repo,
		})
		if err != nil {
			log.Error("failed to encode a response", "error", err)
		}
	}
}
