package api

import (
	"fmt"
	"github.com/jomei/notionapi"
	"github.com/mmmanyfold/rami-notion-api/repo/notion"
	"net/http"
	"os"
	"sync"
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

func (a *API) Sync(w http.ResponseWriter, r *http.Request) {
	var wg sync.WaitGroup
	wg.Add(1)
	err := notion.GetDenormalizedProjects(a.notionClient)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve HomePage projects from notion API"), http.StatusInternalServerError)
		return
	}

	wg.Wait()

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("processed successfully"))
}
