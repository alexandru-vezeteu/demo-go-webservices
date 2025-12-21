package http

type Link struct {
	Href   string `json:"href"`
	Rel    string `json:"rel,omitempty"`
	Method string `json:"method"`
	Title  string `json:"title,omitempty"`
}
