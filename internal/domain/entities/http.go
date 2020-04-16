package entities

// Request data
type Request struct {
	Width   int                 `json:"width"`
	Height  int                 `json:"height"`
	URL     string              `json:"url,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
}

// Response data
type Response struct {
	Preview []byte `json:"preview,omitempty"`
	Status  Status `json:"status"`
}

// StatusError holds together HTTP response status code and text
type Status struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
