package entities

const (
	ProxyingRequest int = iota
	PreviewerResponse
)

var messageTypeText = map[int]string{
	ProxyingRequest:   "Proxying HTTP request",
	PreviewerResponse: "Response from previewer",
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
	Params  *PreviewParams      `json:"params"`
	Headers map[string][]string `json:"headers,omitempty"`
}

// Type implementation for HTTP request proxying
func (r Request) Type() int {
	return ProxyingRequest
}

// Response data
type Response struct {
	Preview *Preview `json:"preview,omitempty"`
	Status  Status   `json:"status"`
}

// Type implementation for image source response
func (r Response) Type() int {
	return PreviewerResponse
}

// StatusError holds together HTTP response status code and text
type Status struct {
	Code int    `json:"code"`
	Text string `json:"text"`
}
