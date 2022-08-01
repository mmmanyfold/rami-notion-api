package rami

type Asset struct {
	UUID string `json:"uuid,omitempty"`
}

type Tag string

type Project struct {
	UUID           string   `json:"uuid,omitempty"`
	ID             string   `json:"id,omitempty"`
	Title          string   `json:"name" json:"title,omitempty"`
	Assets         Asset    `json:"assets,omitempty"`
	Tags           []Tag    `json:"tags,omitempty"`
	Year           string   `json:"year,omitempty"`
	Thumbnail      string   `json:"thumbnail,omitempty"`
	Blocks         []string `json:"blocks,omitempty"`
	Medium         string   `json:"medium,omitempty"`
	Description    string   `json:"description,omitempty"`
	HomePageAssets []Asset  `json:"homePageAssets,omitempty"`
}
