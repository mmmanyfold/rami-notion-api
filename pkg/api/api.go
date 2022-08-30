package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/go-chi/chi"
	"github.com/jomei/notionapi"
	"github.com/mmmanyfold/rami-notion-api/pkg/notion"
	"github.com/mmmanyfold/rami-notion-api/pkg/rami"
	"io/ioutil"
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

func (api *API) Sync(w http.ResponseWriter, r *http.Request) {
	name := chi.URLParam(r, "name")
	if name == "" {
		http.Error(w, fmt.Sprintf("failed, db id URL param not found"), http.StatusNotFound)
		return
	}

	fmt.Printf("requesting db: %s\n", name)

	encoder := json.NewEncoder(w)

	// TODO: build mapping from db-id value to name of db
	switch name {
	case "projects":
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
		payload := rami.ProjectsResponse{
			LastRefreshed: Timestamp(),
			Rows:          rows,
		}
		if err := writeToFile("projects.json", payload); err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("failed to persist json response to disk"), http.StatusInternalServerError)
			return
		}
		if err := encoder.Encode(payload); err != nil {
			http.Error(w, fmt.Sprintf("failed to json encode Projects response"), http.StatusInternalServerError)
			return
		}
	case "cv-additional":
		rows, err := notion.GetCVAdditionalDB(api.notionClient)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("failed to retrieve GetCVAdditional DB from notion API"), http.StatusInternalServerError)
			return
		}
		encoder.SetEscapeHTML(false) // don't encode <, >, &
		payload := rami.CVAdditionalResponse{
			LastRefreshed: Timestamp(),
			Rows:          rows,
		}
		if err := writeToFile("cv-additional.json", payload); err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("failed to persist json response to disk"), http.StatusInternalServerError)
			return
		}
		if err := encoder.Encode(payload); err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("failed to json encode GetCVAdditionalDB response"), http.StatusInternalServerError)
			return
		}
	case "info":
		rows, err := notion.GetInfoDB(api.notionClient)
		if err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("failed to retrieve info DB from notion API"), http.StatusInternalServerError)
			return
		}
		encoder.SetEscapeHTML(false) // don't encode <, >, &
		payload := rami.InfoResponse{
			LastRefreshed: Timestamp(),
			Rows:          rows,
		}
		if err := writeToFile("info.json", payload); err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("failed to persist json response to disk"), http.StatusInternalServerError)
			return
		}
		if err := encoder.Encode(payload); err != nil {
			log.Println(err)
			http.Error(w, fmt.Sprintf("failed to json encode info response"), http.StatusInternalServerError)
			return
		}
	default:
		http.Error(w, fmt.Sprintf("failed to retrieve DB: %s from notion API", name), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func Timestamp() string {
	n := time.Now()
	return fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
		n.Year(), n.Month(), n.Day(),
		n.Hour(), n.Minute(), n.Second())
}

func writeToFile(filename string, jsonData interface{}) error {
	buf := new(bytes.Buffer)
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(false)
	filepath := fmt.Sprintf("./public/%s", filename)

	if err := enc.Encode(&jsonData); err != nil {
		log.Println(err)
	}

	err := ioutil.WriteFile(filepath, buf.Bytes(), 0644)
	if err != nil {
		return err
	}

	return nil
}
