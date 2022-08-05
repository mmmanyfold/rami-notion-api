package rami

import (
	"github.com/jomei/notionapi"
)

type Transcript struct {
	UUID        string   `json:"uuid,omitempty"`
	ProjectUUID string   `json:"ProjectUUID,omitempty"`
	Blocks      []string `json:"blocks,omitempty"`
}

type HomePageAsset struct {
	Files []notionapi.File `json:"files,omitempty"`
	UUID  string           `json:"uuid,omitempty"`
	Type  string           `json:"type"`
	Urls  []string         `json:"urls"`
}

type Tag string

type Project struct {
	UUID           string        `json:"uuid,omitempty"`
	ID             string        `json:"id,omitempty"`
	Title          string        `json:"title" json:"title,omitempty"`
	Tags           []Tag         `json:"tags,omitempty"`
	Year           string        `json:"year,omitempty"`
	Thumbnail      string        `json:"thumbnail,omitempty"`
	Blocks         []string      `json:"blocks,omitempty"`
	Medium         string        `json:"medium,omitempty"`
	Description    string        `json:"description,omitempty"`
	HomePageAssets HomePageAsset `json:"homePageAssets,omitempty"`
	Transcript     Transcript    `json:"transcript,omitempty"`
}

type ProjectsResponse struct {
	LastRefreshed string    `json:"lastRefreshed"`
	AllProjects   []Project `json:"allProjects"`
}
