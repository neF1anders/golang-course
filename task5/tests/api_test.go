package api_test

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

const address = "http://localhost:28080"

var client = http.Client{
	Timeout: 30 * time.Second,
}

type PingService struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}

type PingResponse struct {
	Status   string        `json:"status"`
	Services []PingService `json:"services"`
}

type Repo struct {
	FullName    string `json:"full_name"`
	Description string `json:"description"`
	Stars       int64  `json:"stars"`
	Forks       int64  `json:"forks"`
	CreatedAt   string `json:"created_at"`
}

type ReposResponse struct {
	Repositories []Repo `json:"repositories"`
}

type Slug struct {
	Owner string `json:"owner"`
	Repo  string `json:"repo"`
}

type SlugsResponse struct {
	Subscriptions []Slug `json:"subscriptions"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func waitForAPI(t *testing.T) {
	t.Helper()

	require.Eventually(t, func() bool {
		resp, err := client.Get(address + "/api/ping")
		if err != nil {
			return false
		}
		defer resp.Body.Close()

		return resp.StatusCode == http.StatusOK || resp.StatusCode == http.StatusServiceUnavailable
	}, 20*time.Second, 500*time.Millisecond, "api did not become ready")
}

func serviceMap(services []PingService) map[string]string {
	res := make(map[string]string, len(services))
	for _, svc := range services {
		res[svc.Name] = svc.Status
	}
	return res
}

func subscribe(t *testing.T, owner, repo string) {
	t.Helper()
	url := fmt.Sprintf("%s/subscriptions?owner=%s&repo=%s", address, owner, repo)
	resp, err := client.Post(url, "", nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusConflict {
		t.Fatalf("unexpected status code %d from subscribe", resp.StatusCode)
	}
}

func unsubscribe(t *testing.T, owner, repo string) {
	t.Helper()
	url := fmt.Sprintf("%s/subscriptions?owner=%s&repo=%s", address, owner, repo)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusNoContent && resp.StatusCode != http.StatusNotFound {
		t.Fatalf("unexpected status code %d from unsubscribe", resp.StatusCode)
	}
}

func getSubscriptions(t *testing.T) []Slug {
	t.Helper()
	resp, err := client.Get(address + "/subscriptions")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var body SlugsResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	return body.Subscriptions
}

func getSubscriptionsInfo(t *testing.T) []Repo {
	t.Helper()
	resp, err := client.Get(address + "/subscriptions/info")
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusOK, resp.StatusCode)
	var body ReposResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body))
	return body.Repositories
}
func TestPreflight(t *testing.T) {
	require.Equal(t, true, true)
}

func TestPing(t *testing.T) {
	waitForAPI(t)

	resp, err := client.Get(address + "/api/ping")
	require.NoError(t, err, "cannot ping api")
	defer resp.Body.Close()

	var body PingResponse
	require.NoError(t, json.NewDecoder(resp.Body).Decode(&body), "cannot decode ping response")

	require.Equal(t, http.StatusOK, resp.StatusCode, "wrong status code")
	require.Equal(t, "ok", body.Status, "wrong overall status")

	services := serviceMap(body.Services)
	require.Equal(t, "up", services["processor"], "processor should be up")
	require.Equal(t, "up", services["subscriber"], "subscriber should be up")
}

func TestPingHelpfulFailureMessage(t *testing.T) {
	waitForAPI(t)

	resp, err := client.Get(address + "/api/ping")
	require.NoError(t, err)
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusServiceUnavailable {
		var body PingResponse
		_ = json.NewDecoder(resp.Body).Decode(&body)

		services := serviceMap(body.Services)
		t.Fatalf("api is degraded: processor=%s subscriber=%s", services["processor"], services["subscriber"])
	}

	require.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestSubscriptionLifecycle(t *testing.T) {
	waitForAPI(t)
	owner := "golang"
	repo := "go"
	unsubscribe(t, owner, repo)
	subs := getSubscriptions(t)
	for _, s := range subs {
		if s.Owner == owner && s.Repo == repo {
			t.Fatal("subscription already exists before test")
		}
	}
	subscribe(t, owner, repo)
	subs = getSubscriptions(t)
	found := false
	for _, s := range subs {
		if s.Owner == owner && s.Repo == repo {
			found = true
			break
		}
	}
	require.True(t, found, "subscription not found after subscribe")
	require.Eventually(t, func() bool {
		repos := getSubscriptionsInfo(t)
		for _, r := range repos {
			if r.FullName == owner+"/"+repo {
				if r.Stars > 0 && r.Description != "" {
					return true
				}
			}
		}
		return false
	}, 20*time.Second, 1*time.Second, "info for subscribed repo did not appear or has zero stars")
	unsubscribe(t, owner, repo)
	subs = getSubscriptions(t)
	for _, s := range subs {
		if s.Owner == owner && s.Repo == repo {
			t.Fatal("subscription still exists after unsubscribe")
		}
	}
}

func TestSubscribeDuplicate(t *testing.T) {
	waitForAPI(t)

	owner := "golang"
	repo := "go"
	unsubscribe(t, owner, repo)
	subscribe(t, owner, repo)
	url := fmt.Sprintf("%s/subscriptions?owner=%s&repo=%s", address, owner, repo)
	resp, err := client.Post(url, "", nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusConflict, resp.StatusCode, "duplicate subscription should return 409")
	unsubscribe(t, owner, repo)
}

func TestSubscribeInvalidRepo(t *testing.T) {
	waitForAPI(t)

	owner := "nonexistent"
	repo := "nonexistent"
	url := fmt.Sprintf("%s/subscriptions?owner=%s&repo=%s", address, owner, repo)
	resp, err := client.Post(url, "", nil)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNotFound, resp.StatusCode, "invalid repo should return 404")
}

func TestUnsubscribeNotFound(t *testing.T) {
	waitForAPI(t)

	owner := "golang"
	repo := "go"
	unsubscribe(t, owner, repo)
	url := fmt.Sprintf("%s/subscriptions?owner=%s&repo=%s", address, owner, repo)
	req, err := http.NewRequest(http.MethodDelete, url, nil)
	require.NoError(t, err)
	resp, err := client.Do(req)
	require.NoError(t, err)
	defer resp.Body.Close()
	require.Equal(t, http.StatusNotFound, resp.StatusCode, "unsubscribe non-existent should return 404")
}

func TestSubscriptionsInfoMultiple(t *testing.T) {
	waitForAPI(t)
	repos := []struct{ owner, repo string }{
		{"golang", "go"},
		{"kubernetes", "kubernetes"},
	}
	for _, r := range repos {
		unsubscribe(t, r.owner, r.repo)
	}
	for _, r := range repos {
		subscribe(t, r.owner, r.repo)
	}
	subs := getSubscriptions(t)
	found := make(map[string]bool)
	for _, s := range subs {
		key := s.Owner + "/" + s.Repo
		found[key] = true
	}
	for _, r := range repos {
		key := r.owner + "/" + r.repo
		require.True(t, found[key], "missing subscription for %s", key)
	}
	require.Eventually(t, func() bool {
		infos := getSubscriptionsInfo(t)
		infoMap := make(map[string]Repo)
		for _, repo := range infos {
			infoMap[repo.FullName] = repo
		}
		for _, r := range repos {
			repo, ok := infoMap[r.owner+"/"+r.repo]
			if !ok || repo.Stars == 0 {
				return false
			}
		}
		return true
	}, 20*time.Second, 1*time.Second, "not all repositories had info with stars")

	// Очистка
	for _, r := range repos {
		unsubscribe(t, r.owner, r.repo)
	}
}
