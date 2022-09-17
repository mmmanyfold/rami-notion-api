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
	UUID           string               `json:"uuid,omitempty"`
	ID             string               `json:"id,omitempty"`
	Title          string               `json:"title" json:"title,omitempty"`
	Tags           []Tag                `json:"tags,omitempty"`
	Year           string               `json:"year,omitempty"`
	Thumbnail      string               `json:"thumbnail,omitempty"`
	Blocks         []string             `json:"blocks,omitempty"`
	Medium         []notionapi.RichText `json:"medium,omitempty"`
	Description    []notionapi.RichText `json:"description,omitempty"`
	HomePageAssets HomePageAsset        `json:"homePageAssets,omitempty"`
	Transcript     Transcript           `json:"transcript,omitempty"`
	Slug           string               `json:"slug,omitempty"`
}

type ProjectsResponse struct {
	LastRefreshed string    `json:"lastRefreshed"`
	Rows          []Project `json:"rows"`
}

type Info struct {
	UUID        string               `json:"uuid,omitempty"`
	ProjectUUID string               `json:"ProjectUUID,omitempty"`
	Tag         string               `json:"tag,omitempty"`
	Line1       []notionapi.RichText `json:"line-1,omitempty"`
	Line2       []notionapi.RichText `json:"line-2,omitempty"`
	Line3       []notionapi.RichText `json:"line-3,omitempty"`
	Line4       []notionapi.RichText `json:"line-4,omitempty"`
	URL         string               `json:"url"`
	Download    []File               `json:"download,omitempty"`
}

type InfoResponse struct {
	LastRefreshed string `json:"lastRefreshed"`
	Rows          []Info `json:"rows"`
}

type CVAdditional struct {
	UUID        string               `json:"uuid,omitempty"`
	Tag         string               `json:"tag"`
	Title       []notionapi.RichText `json:"title"`
	Description []notionapi.RichText `json:"description"`
	Detail      []notionapi.RichText `json:"detail"`
	URL         string               `json:"url"`
	Download    []File               `json:"download,omitempty"`
}

type CVExhibitionsAndScreening struct {
	UUID        string               `json:"uuid,omitempty"`
	Title       []notionapi.RichText `json:"title"`
	Description []notionapi.RichText `json:"description"`
	Detail      []notionapi.RichText `json:"detail"`
	URL         string               `json:"url"`
	Download    []File               `json:"download,omitempty"`
	Year        string               `json:"year,omitempty"`
}

type CVExhibitionsAndScreeningResponse struct {
	LastRefreshed string                      `json:"lastRefreshed"`
	Rows          []CVExhibitionsAndScreening `json:"rows"`
}

type CVAdditionalResponse struct {
	LastRefreshed string         `json:"lastRefreshed"`
	Rows          []CVAdditional `json:"rows"`
}
