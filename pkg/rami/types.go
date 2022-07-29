package rami

type Asset struct {
}

type HomePageAsset struct {
	UUID   string `json:"UUID,omitempty"`
	ID     string `json:"id" json:"ID,omitempty"`
	Name   string `json:"name" json:"name,omitempty"`
	Assets Asset  `json:"assets" json:"assets"`
}
