package entities

const (
	ProxyingRequest int = iota
	SourceResponse
)

var messageTypeText = map[int]string{
	ProxyingRequest: "Proxying HTTP request",
	SourceResponse:  "Response from image source",
}

// MessageTypeToString convert from integer value
func MessageTypeToString(t int) string {
	return messageTypeText[t]
}

// StringToMessageType convert from string
func StringToMessageType(s string) (key int, ok bool) {
	for k, v := range messageTypeText {
		if v == s {
			key = k
			ok = true
			return
		}
	}
	return
}

// Request data
type Request struct {
	Width   int                 `json:"width"`
	Height  int                 `json:"height"`
	URL     string              `json:"url,omitempty"`
	Headers map[string][]string `json:"headers,omitempty"`
}

// Type implementation for HTTP request proxying
func (r Request) Type() int {
	return ProxyingRequest
}

// Response data
type Response struct {
	Preview  []byte `json:"preview,omitempty"`
	Filename string `json:"filename,omitempty"`
	Status   Status `json:"status"`
}

// Type implementation for image source response
func (r Response) Type() int {
	return SourceResponse
}

// StatusError holds together HTTP response status code and text
type Status struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
