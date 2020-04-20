package entities

// PreviewParams
type PreviewParams struct {
	Width  int    `json:"width"`
	Height int    `json:"height"`
	URL    string `json:"url,omitempty"`
}

// Preview
type Preview struct {
	Params *PreviewParams `json:"params"`
	Image  []byte         `json:"image,omitempty"`
}
