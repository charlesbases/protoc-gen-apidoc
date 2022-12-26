package postman

import (
	"encoding/json"
	"strings"

	"github.com/charlesbases/protoc-gen-apidoc/conf"
	"github.com/charlesbases/protoc-gen-apidoc/encoder"
	"github.com/charlesbases/protoc-gen-apidoc/generator"
	"github.com/charlesbases/protoc-gen-apidoc/logger"
	"github.com/charlesbases/protoc-gen-apidoc/types"
	"github.com/google/uuid"
)

// newPostman .
func newPostman(p *types.Package) *Postman {
	return &Postman{
		p: p,
		host: func() *URL {
			var (
				url  = new(URL)
				host = conf.Get().Host
			)

			// http or https
			if strings.HasPrefix(host, "https") {
				url.Protocol = "https"
			} else {
				url.Protocol = "http"
			}

			if idx := strings.Index(host, "//"); idx != -1 {
				host = host[idx+2:]
			}

			url.Host = strings.Split(host, ".")
			url.Port = conf.Get().Port

			return url
		}(),
		header: func() []*Header {
			var h = make([]*Header, 0, len(conf.Get().Header))
			for _, header := range conf.Get().Header {
				h = append(h, &Header{
					Key:  header.String(),
					Type: "default",
				})
			}
			return h
		}(),
		Info: &Info{
			ID:     uuid.New().String(),
			Name:   conf.Get().Title,
			Schema: "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		},
		Item: make([]*Service, 0, len(p.Services)),
	}
}

// NewGenerator .
func NewGenerator(p *types.Package) generator.Generator {
	pt := newPostman(p)

	pt.parseServiceList()
	return pt
}

// Generate .
func (pt *Postman) Generate() []byte {
	data, err := json.MarshalIndent(pt, "", "  ")
	if err != nil {
		logger.Fatal(err)
	}

	return data
}

// parseServiceList .
func (pt *Postman) parseServiceList() {
	for _, srv := range pt.p.Services {
		pt.Item = append(pt.Item, pt.parseService(srv))
	}
}

// parseService .
func (pt *Postman) parseService(srv *types.Service) *Service {
	var ptService = &Service{
		Name: srv.Name,
		Item: make([]*API, 0, len(srv.Methods)),
	}

	for _, api := range srv.Methods {
		ptService.Item = append(ptService.Item, pt.parseServiceAPI(api))
	}

	return ptService
}

// parseServiceAPI .
func (pt *Postman) parseServiceAPI(api *types.ServiceMethod) *API {
	var ptAPI = &API{
		Name: api.Path,
		Request: &Request{
			Method: api.Method,
			Header: pt.header,
			URL: &URL{
				Raw:      conf.Get().Addr + api.Path,
				Protocol: pt.host.Protocol,
				Host:     pt.host.Host,
				Port:     pt.host.Port,
				Path:     strings.Split(strings.TrimPrefix(api.Path, "/"), "/"),
			},
		},
	}

	if mess, find := pt.p.MessageDic[api.RequestName]; find && len(mess.Fields) != 0 {
		switch api.Method {
		// Query
		case types.Method_Get:
			ptAPI.Request.URL.Query = make([]*Query, 0, len(mess.Fields))
			for _, field := range mess.Fields {
				ptAPI.Request.URL.Query = append(ptAPI.Request.URL.Query, &Query{
					Key:         field.JsonName,
					Description: field.Description,
				})
			}
		// Body
		default:
			switch api.Consume {
			case types.ContentType_Json:
				ptAPI.Request.Body = &Body{
					Mode: "raw",
					Raw:  encoder.NewEncoder(pt.p).EncodeJson(api.RequestName),
					Options: BodyOptions{
						Raw: struct {
							Language string `json:"language"`
						}(struct{ Language string }{Language: "json"}),
					},
				}
			case types.ContentType_Data:
				ptAPI.Request.Body = &Body{
					Mode: "formdata",
					Formdata: []*BodyFormData{
						{
							Key:  "data",
							Type: "file",
						},
					},
				}
			}
		}
	}

	return ptAPI
}
