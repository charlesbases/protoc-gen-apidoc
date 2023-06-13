package conf

import (
	"fmt"
	"strings"

	"github.com/charlesbases/protoc-gen-apidoc/logger"
	"github.com/charlesbases/protoc-gen-apidoc/types"
)

const (
	argHost    arg = "host"
	argPort    arg = "port"
	argTitle   arg = "title"
	argHeader  arg = "header"
	argOutput  arg = "output"
	argschemes arg = "scheme"
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
			case argschemes:
				if len(conf.Schemes) == 0 {
					conf.Schemes = make([]string, 0, 2)
				}
				conf.Schemes = append(conf.Schemes, value)
			case argOutput:
				switch types.DocumentType(value) {
				case types.DocumentTypeSwagger:
					var filename = "swagger.json"
					if len(conf.Title) != 0 {
						filename = fmt.Sprintf("%s.%s", strings.ToLower(conf.Title), filename)
					}

					conf.Document = append(conf.Document, &Document{
						Type: types.DocumentTypeSwagger,
						File: filename,
					})
				case types.DocumentTypePostman:
					var filename = "postman.json"
					if len(conf.Title) != 0 {
						filename = fmt.Sprintf("%s.%s", strings.ToLower(conf.Title), filename)
					}

					conf.Document = append(conf.Document, &Document{
						Type: types.DocumentTypePostman,
						File: filename,
					})
				case types.DocumentTypeHTML:
					var filename = "apidoc.html"
					if len(conf.Title) != 0 {
						filename = fmt.Sprintf("%s.html", strings.ToLower(conf.Title))
					}

					conf.Document = append(conf.Document, &Document{
						Type: types.DocumentTypeHTML,
						File: filename,
					})
				case types.DocumentTypeMarkdown:
					var filename = "apidoc.md"
					if len(conf.Title) != 0 {
						filename = fmt.Sprintf("%s.md", strings.ToLower(conf.Title))
					}

					conf.Document = append(conf.Document, &Document{
						Type: types.DocumentTypeMarkdown,
						File: filename,
					})
				default:
					logger.Fatalf(`invalid type of "%s"`, value)
				}
			}
		}
	}

	// default document
	if len(conf.Document) == 0 {
		logger.Debugf(`use default document: "swagger"`)
		conf.Document = append(conf.Document, &Document{
			Type: types.DocumentTypeSwagger,
			File: "swagger.json",
		})
	}

	return conf
}
