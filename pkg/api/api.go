package api

import (
	"encoding/json"
	"fmt"
	"github.com/jomei/notionapi"
	"github.com/mmmanyfold/rami-notion-api/pkg/notion"
	"log"
	"net/http"
)

type API struct {
	notionClient *notionapi.Client
}

func NewAPI(notionAPIKey string) (*API, error) {
	client := notionapi.NewClient(notionapi.Token(notionAPIKey))
	return &API{
		notionClient: client,
	}, nil
}

func (api *API) Sync(w http.ResponseWriter, r *http.Request) {
	transcripts, err := notion.GetTranscripts(api.notionClient)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("failed to retrieve Transcripts from notion API"), http.StatusInternalServerError)
		return
	}

	assets, err := notion.GetHomePageAssets(api.notionClient)
	if err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("failed to retrieve HomePageAssets from notion API"), http.StatusInternalServerError)
		return
	}

	projects, err := notion.GetProjects(api.notionClient, assets, transcripts)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve Projects from notion API"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false) // don't encode <, >, &
	if err := encoder.Encode(projects); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("failed to retrieve Projects from notion API"), http.StatusInternalServerError)
		return
	}
}
