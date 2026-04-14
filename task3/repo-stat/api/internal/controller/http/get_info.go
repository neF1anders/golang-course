package http

import (
	"encoding/json"
	"errors"
	"log/slog"
	"net/http"
	"net/url"
	"repo-stat/api/internal/dto"
	"repo-stat/api/internal/usecase"
	"strings"

	_ "repo-stat/api/docs"
)

// RepoInfo godoc
// @Summary Get repository info
// @Description Get GitHub repository information
// @Tags repositories
// @Accept json
// @Produce json
// @Param url query string true "GitHub repo URL"
// @Success 200 {object} dto.Repo
// @Failure 400 {string} string
// @Failure 404 {string} string
// @Failure 503 {string} string
// @Router /api/repositories/info [get]
func NewGetInfoHandler(log *slog.Logger, fetch *usecase.Fetch) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		repoURL := r.URL.Query().Get("url")
		if repoURL == "" {
			http.Error(w, "missing url parameter", http.StatusBadRequest)
			return
		}
		owner, repo, err := ParseGitHubRepo(repoURL)
		if owner == "" || repo == "" {
			http.Error(w, "owner and repo are required", http.StatusBadRequest)
			return
		}
		info, err := fetch.Execute(r.Context(), owner, repo)
		if err != nil {
			writeGRPCError(w, err)
			log.Error("failed to execute fetch response", "error", err)
		}
		response := dto.Repo{
			Name:        info.Name,
			Description: info.Description,
			Stars:       info.Stars,
			Forks:       info.Forks,
			Date:        info.Date,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			log.Error("failed to write ping response", "error", err)
			return
		}
	}
}

func ParseGitHubRepo(rawURL string) (owner, repo string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}
	// Для SSH вида git@github.com:owner/repo.git путь будет "owner/repo.git"
	path := strings.Trim(u.Path, "/")
	if path == "" && u.Opaque != "" && strings.Contains(u.Opaque, ":") {
		// Возможно, это SSH: git@github.com:owner/repo.git
		parts := strings.SplitN(u.Opaque, ":", 2)
		if len(parts) == 2 {
			path = parts[1]
		}
	}
	parts := strings.Split(path, "/")
	if len(parts) < 2 {
		return "", "", errors.New("invalid GitHub repo URL: missing owner or repo")
	}
	owner = parts[0]
	repo = strings.TrimSuffix(parts[1], ".git")
	return owner, repo, nil
}
