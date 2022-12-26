package postman

import "github.com/charlesbases/protoc-gen-apidoc/types"

// Postman .
type Postman struct {
	p      *types.Package `json:"-"`
	host   *URL           `json:"-"`
	header []*Header      `json:"-"`

	Info *Info      `json:"info"`
	Item []*Service `json:"item"`
}

// Info .
type Info struct {
	ID     string `json:"_postman_id"`
	Name   string `json:"name"`
	Schema string `json:"schema"`
}

// Service .
type Service struct {
	Name string `json:"name"`
	Item []*API `json:"item"`
}

// API .
type API struct {
	Name     string    `json:"name"`
	Request  *Request  `json:"request,omitempty"`
	Response *Response `json:"response,omitempty"`
}

// Request .
type Request struct {
	Method types.Method `json:"method"`
	Header []*Header    `json:"header"`
	Body   *Body        `json:"body"`
	URL    *URL         `json:"url"`
}

// Response .
type Response struct{}

// Header .
type Header struct {
	Key   string `json:"key"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// Body .
type Body struct {
	Mode     string          `json:"mode"`
	Raw      string          `json:"raw,omitempty"`
	Options  BodyOptions     `json:"options,omitempty"`
	Formdata []*BodyFormData `json:"formdata,omitempty"`
}

// BodyOptions .
type BodyOptions struct {
	Raw struct {
		Language string `json:"language"`
	} `json:"raw"`
}

// BodyFormData .
type BodyFormData struct {
	Key  string `json:"key"`
	Type string `json:"type"`
	Src  []byte `json:"src"`
}

// Query .
type Query struct {
	Key      string `json:"key"`
	Value    string `json:"value"`
	Disabled bool   `json:"disabled"`
}

// URL .
type URL struct {
	Raw      string   `json:"raw"`
	Protocol string   `json:"protocol"`
	Host     []string `json:"host"`
	Port     string   `json:"port"`
	Path     []string `json:"path"`
	Query    []*Query `json:"query"`
}
