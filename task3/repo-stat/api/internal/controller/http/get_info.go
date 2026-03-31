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
)

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
		}
	}
}

func ParseGitHubRepo(rawURL string) (owner, repo string, err error) {
	u, err := url.Parse(rawURL)
	if err != nil {
		return "", "", err
	}
	// Ожидаем путь вида /owner/repo (возможно с trailing slash)
	parts := strings.Split(strings.Trim(u.Path, "/"), "/")
	if len(parts) < 2 {
		return "", "", errors.New("invalid GitHub repo URL: missing owner or repo")
	}
	return parts[0], parts[1], nil
}
