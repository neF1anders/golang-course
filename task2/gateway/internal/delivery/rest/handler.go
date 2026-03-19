package rest

import (
	"encoding/json"
	"net/http"

	"gateway/internal/external/grpc"
)

type Handler struct {
	collectorClient *grpc.CollectorClient
}

func NewHandler(collectorClient *grpc.CollectorClient) *Handler {
	return &Handler{
		collectorClient: collectorClient,
	}
}
func (h *Handler) GetRepoInfo(w http.ResponseWriter, r *http.Request) {
	owner := r.URL.Query().Get("owner")
	repo := r.URL.Query().Get("repo")
	if owner == "" || repo == "" {
		http.Error(w, "owner and repo are required", http.StatusBadRequest)
		return
	}
	repoInfo, err := h.collectorClient.GetRepoInfo(owner, repo)
	if err != nil {
		http.Error(w, "failed to get repo info: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(repoInfo)
}
