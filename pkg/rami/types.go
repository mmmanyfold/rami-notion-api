package rami

import "github.com/jomei/notionapi"

type HomePageAsset struct {
	UUID  string           `json:"uuid,omitempty"`
	Files []notionapi.File `json:"files,omitempty"`
}

type Tag string

type Project struct {
	UUID           string        `json:"uuid,omitempty"`
	ID             string        `json:"id,omitempty"`
	Title          string        `json:"name" json:"title,omitempty"`
	Tags           []Tag         `json:"tags,omitempty"`
	Year           string        `json:"year,omitempty"`
	Thumbnail      string        `json:"thumbnail,omitempty"`
	Blocks         []string      `json:"blocks,omitempty"`
	Medium         string        `json:"medium,omitempty"`
	Description    string        `json:"description,omitempty"`
	HomePageAssets HomePageAsset `json:"homePageAssets,omitempty"`
}
