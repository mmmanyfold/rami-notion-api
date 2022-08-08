package rami

import (
	"github.com/jomei/notionapi"
)

type Transcript struct {
	UUID        string            `json:"uuid,omitempty"`
	ProjectUUID string            `json:"ProjectUUID,omitempty"`
	Blocks      []notionapi.Block `json:"blocks,omitempty"`
}

type File struct {
	Url    string `json:"url,omitempty"`
	Width  uint64 `json:"width,omitempty"`
	Height uint64 `json:"height,omitempty"`
}

type HomePageAsset struct {
	NotionFiles []notionapi.File `json:"notionFiles,omitempty"`
	UUID        string           `json:"uuid,omitempty"`
	Type        string           `json:"type,omitempty"`
	Files       []File           `json:"files,omitempty"`
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
	Slug           string        `json:"slug,omitempty"`
}

type ProjectsResponse struct {
	LastRefreshed string    `json:"lastRefreshed"`
	Rows          []Project `json:"allProjects"`
}

type Info struct {
	UUID        string `json:"uuid,omitempty"`
	ProjectUUID string `json:"ProjectUUID,omitempty"`
}

type InfoResponse struct {
	LastRefreshed string    `json:"lastRefreshed"`
	AllInfo       []Project `json:"allProjects"`
}

type CVAdditional struct {
	UUID        string `json:"uuid,omitempty"`
	Tag         string `json:"tag"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Detail      string `json:"detail"`
	URL         string `json:"url"`
}

type CVAdditionalResponse struct {
	LastRefreshed string         `json:"lastRefreshed"`
	Rows          []CVAdditional `json:"rows"`
}
