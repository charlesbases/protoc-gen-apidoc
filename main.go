package main

import (
	"github.com/charlesbases/protoc-gen-apidoc/conf"
	"github.com/charlesbases/protoc-gen-apidoc/generator"
	"github.com/charlesbases/protoc-gen-apidoc/generator/postman"
	"github.com/charlesbases/protoc-gen-apidoc/generator/swagger"
	"github.com/charlesbases/protoc-gen-apidoc/generator/template"
	"github.com/charlesbases/protoc-gen-apidoc/logger"
	"github.com/charlesbases/protoc-gen-apidoc/protoc"
	"github.com/charlesbases/protoc-gen-apidoc/types"
	"google.golang.org/protobuf/types/pluginpb"
)

func main() {
	protoc.Plugin(func(p *types.Package) *pluginpb.CodeGeneratorResponse {
		var rsp = new(pluginpb.CodeGeneratorResponse)

		for _, dt := range conf.Get().Document {
			var gen generator.Generator
			switch dt.Type {
			case types.DocumentTypeHTML:
				gen = template.NewGenerator(p, template.HTML)
			case types.DocumentTypeMarkdown:
				gen = template.NewGenerator(p, template.Markdown)
			case types.DocumentTypeSwagger:
				gen = swagger.NewGenerator(p)
			case types.DocumentTypePostman:
				gen = postman.NewGenerator(p)
			default:
				logger.Fatalf(`invalid type of "%s"`, dt.Type)
			}

			if data := gen.Generate(); len(data) != 0 {
				var content = string(data)
				rsp.File = append(rsp.File, &pluginpb.CodeGeneratorResponse_File{
					Name:    &dt.File,
					Content: &content,
				})
			}
		}

		return rsp
	})
}
