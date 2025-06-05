package links_client

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"
)

var (
	baseURL    = env("LINKS_STORE_URL", "http://localhost:8080")
	httpClient = &http.Client{Timeout: 8 * time.Second}
)

func env(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

type Link struct {
	ID    int64  `json:"id"`
	URL   string `json:"url"`
	Label string `json:"label"`
}

func List() ([]Link, error) {
	resp, err := httpClient.Get(baseURL + "/api/v1/links")
	if err != nil {
		return nil, fmt.Errorf("links_client: list: %w", err)
	}
	defer resp.Body.Close()

	var links []Link
	if err := json.NewDecoder(resp.Body).Decode(&links); err != nil {
		return nil, fmt.Errorf("links_client: decode list: %w", err)
	}
	return links, nil
}

func Get(id int64) (Link, error) {
	resp, err := httpClient.Get(fmt.Sprintf("%s/api/v1/links/%d", baseURL, id))
	if err != nil {
		return Link{}, fmt.Errorf("links_client: get: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return Link{}, fmt.Errorf("links_client: %d not found", id)
	}

	var l Link
	if err := json.NewDecoder(resp.Body).Decode(&l); err != nil {
		return Link{}, fmt.Errorf("links_client: decode link: %w", err)
	}
	return l, nil
}
