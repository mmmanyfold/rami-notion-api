package api

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/jomei/notionapi"
	"github.com/mmmanyfold/rami-notion-api/repo/notion"
	"net/http"
	"os"
)

type API struct {
	notionClient *notionapi.Client
}

func New() *API {
	notionAPIKey := os.Getenv("NOTION_API_KEY")
	client := notionapi.NewClient(notionapi.Token(notionAPIKey))
	return &API{
		notionClient: client,
	}
}

func (api *API) Sync(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	rateLimiter, err := notion.NewRateLimiter(ctx, "projects", notion.Rate, true)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	_, _, err = rateLimiter.Take()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	transcripts, err := notion.GetTranscripts(api.notionClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve HomePageAssets from notion API"), http.StatusInternalServerError)
		return
	}

	//time.Sleep(time.Second)

	assets, err := notion.GetHomePageAssets(api.notionClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve HomePageAssets from notion API"), http.StatusInternalServerError)
		return
	}

	//time.Sleep(time.Second)

	projects, err := notion.GetProjects(api.notionClient, assets, transcripts)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve Projects from notion API"), http.StatusInternalServerError)
		return
	}

	if err := rateLimiter.Close(); err != nil {
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
