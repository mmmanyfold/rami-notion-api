package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jomei/notionapi"
	"net/http"
)

type API struct {
	notionClient *notionapi.Client
	rateLimiter  *RateLimiter
}

func NewAPI(notionAPIKey string, ctx context.Context) (*API, error) {
	client := notionapi.NewClient(notionapi.Token(notionAPIKey))
	rateLimiter, err := NewRateLimiter(ctx, "projects", Rate, Limit, true)
	if err != nil {
		return nil, err
	}

	return &API{
		notionClient: client,
		rateLimiter:  rateLimiter,
	}, nil
}

func (api *API) Sync(w http.ResponseWriter, r *http.Request) {
	_, _, err := api.rateLimiter.Take()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transcripts, err := GetTranscripts(api.notionClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve HomePageAssets from notion API"), http.StatusInternalServerError)
		return
	}

	//time.Sleep(time.Second)

	assets, err := GetHomePageAssets(api.notionClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve HomePageAssets from notion API"), http.StatusInternalServerError)
		return
	}

	//time.Sleep(time.Second)

	projects, err := GetProjects(api.notionClient, assets, transcripts)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve Projects from notion API"), http.StatusInternalServerError)
		return
	}

	if err := api.rateLimiter.Close(); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false) // don't encode <, >, &
	if err := encoder.Encode(projects); err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve Projects from notion API"), http.StatusInternalServerError)
		return
	}
}
