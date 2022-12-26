package conf

import (
	"strings"

	"github.com/charlesbases/protoc-gen-apidoc/logger"
	"github.com/charlesbases/protoc-gen-apidoc/types"
)

const (
	argHost   arg = "host"
	argPort   arg = "port"
	argTitle  arg = "title"
	argHeader arg = "header"
	argOutput arg = "output"
)

// argsOptions .
type argsOptions struct {
	args string
}

// newArgsParser .
func newArgsParser(args string) parser {
	return &argsOptions{args: args}
}

// parse .
func (opts *argsOptions) parse() *configuration {
	var conf = &configuration{
		Header:   make([]types.Header, 0),
		Document: make([]*Document, 0),
	}

	if len(opts.args) != 0 {
		for _, param := range strings.Split(opts.args, ",") {
			var value string
			if i := strings.Index(param, "="); i >= 0 {
				value = param[i+1:]
				param = param[0:i]
			}

			switch arg(param) {
			case argHost:
				conf.Host = value
			case argPort:
				conf.Port = value
			case argTitle:
				conf.Title = value
			case argHeader:
				conf.Header = append(conf.Header, types.Header(value))
			case argOutput:
				switch types.DocumentType(value) {
				case types.DocumentType_Swagger:
					conf.Document = append(conf.Document, &Document{
						Type: types.DocumentType_Swagger,
						File: "swagger.json",
					})
				case types.DocumentType_Postman:
					conf.Document = append(conf.Document, &Document{
						Type: types.DocumentType_Postman,
						File: "postman.json",
					})
				case types.DocumentType_HTML:
					conf.Document = append(conf.Document, &Document{
						Type: types.DocumentType_Swagger,
						File: "apidoc.html",
					})
				case types.DocumentType_Markdown:
					conf.Document = append(conf.Document, &Document{
						Type: types.DocumentType_Swagger,
						File: "apidoc.md",
					})
				default:
					logger.Fatalf(`invalid type of "%s"`, value)
				}
			}
		}
	}

	return conf
}
