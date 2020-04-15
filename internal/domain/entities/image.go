package entities

// Image data
type Image struct {
	Width   int           `json:"width"`
	Height  int           `json:"height"`
	URL     string        `json:"url,omitempty"`
	Headers []ImageHeader `json:"headers,omitempty"`
}

// ImageHeader holds key-value pair for HTTP header
type ImageHeader struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}
