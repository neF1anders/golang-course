package http

import (
	"encoding/json"
	"errors"
	"fmt"
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
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]string{"error": "missing url parameter"})
			if err != nil {
				log.Error("failed to encode error json", "error", err)
			}
			return
		}
		owner, repo, err := ParseGitHubRepo(repoURL)
		if owner == "" || repo == "" || err != nil {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusBadRequest)
			err := json.NewEncoder(w).Encode(map[string]string{"error": "invalid url parameters"})
			if err != nil {
				log.Error("failed to encode error json", "error", err)
			}
			fmt.Printf("Error in parsing url: %v\n", err)
			return
		}
		info, err := fetch.GetInfo(r.Context(), owner, repo)
		if err != nil {
			writeGRPCError(w, err)
			log.Error("failed to execute fetch response", "error", err)
			return
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
