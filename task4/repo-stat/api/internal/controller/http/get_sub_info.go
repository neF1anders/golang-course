package http

import (
	"encoding/json"
	"log/slog"
	"net/http"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"

	_ "repo-stat/api/docs"
)

// SubRepoInfo godoc
// @Summary Get subscribed repositories info
// @Description Get subscribed GitHub repositories information
// @Tags repositories
// @Produce json
// @Success 200 {object} dto.Repos
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 503 {string} string
// @Router /subscriptions/info [get]
func NewGetSubInfoHandler(log *slog.Logger, fetch *usecase.Fetch) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		info, err := fetch.GetSubInfo(r.Context())
		if err != nil {
			writeGRPCError(w, err)
			log.Error("failed to execute sub-fetch response", "error", err)
			return
		}
		response := dto.Repos{}
		for _, el := range info {
			response.Repos = append(response.Repos, dto.Repo{
				Name:        el.Name,
				Description: el.Description,
				Stars:       el.Stars,
				Forks:       el.Forks,
				Date:        el.Date,
			})
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write sub-fetch response", "error", err)
			return
		}
	}
}
