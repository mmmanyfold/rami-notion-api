package api

import (
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jomei/notionapi"
	"github.com/mmmanyfold/rami-notion-api/pkg/notion"
	"github.com/mmmanyfold/rami-notion-api/pkg/rami"
	"log"
	"net/http"
	"time"
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

func (api *API) GetDB(w http.ResponseWriter, r *http.Request) {
	dbID := chi.URLParam(r, "id")
	if dbID == "" {
		http.Error(w, fmt.Sprintf("failed, db id URL param not found"), http.StatusNotFound)
		return
	}

	fmt.Printf("requesting db: %s\n", dbID)

	encoder := json.NewEncoder(w)

	// TODO: build mapping from db-id value to name of db
	switch dbID {
	case "78403af9f31145ce98c7a9ffa57931f8":
		rows, err := notion.GetCVAdditionalDB(api.notionClient)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("failed to retrieve GetCVAdditional DB from notion API"), http.StatusInternalServerError)
			return
		}
		encoder.SetEscapeHTML(false) // don't encode <, >, &
		if err := encoder.Encode(rami.CVAdditionalResponse{
			LastRefreshed: Timestamp(),
			Rows:          rows,
		}); err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("failed to json encode GetCVAdditionalDB response"), http.StatusInternalServerError)
			return
		}
	}

	w.WriteHeader(http.StatusOK)
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

	rows, err := notion.GetProjects(api.notionClient, assets, transcripts)
	if err != nil {
		http.Error(w, fmt.Sprintf("failed to retrieve Projects from notion API"), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	encoder := json.NewEncoder(w)
	encoder.SetEscapeHTML(false) // don't encode <, >, &
	if err := encoder.Encode(rami.ProjectsResponse{
		LastRefreshed: Timestamp(),
		Rows:          rows,
	}); err != nil {
		log.Println(err)
		http.Error(w, fmt.Sprintf("failed to json encode Projects response"), http.StatusInternalServerError)
		return
	}
}

func Timestamp() string {
	n := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		n.Year(), n.Month(), n.Day(),
		n.Hour(), n.Minute(), n.Second())
}
