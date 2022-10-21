package main

import (
	"os"
	"path/filepath"
	"sync"

	"github.com/charlesbases/protoc-gen-apidoc/conf"
	"github.com/charlesbases/protoc-gen-apidoc/generator/swagger"
	"github.com/charlesbases/protoc-gen-apidoc/logger"
	"github.com/charlesbases/protoc-gen-apidoc/protoc"
	"github.com/charlesbases/protoc-gen-apidoc/types"
)

func main() {
	protoc.Plugin(func(p *types.Package) {
		var swg = sync.WaitGroup{}
		swg.Add(len(conf.Get().Document))

		for _, dt := range conf.Get().Document {
			dt.File, _ = filepath.Abs(dt.File)
			var dir = filepath.Dir(dt.File)
			if _, err := os.Stat(dir); err != nil && os.IsNotExist(err) {
				os.MkdirAll(dir, 0755)
			}

			go func(dt *conf.Document) {
				var data []byte
				switch dt.Type {
				case types.DocumentType_HTML:
				case types.DocumentType_Markdown:
				case types.DocumentType_Swagger:
					data = swagger.New(p).Generate()
				case types.DocumentType_Postman:
				default:
					logger.Fatalf(`invalid type of "%s"`, dt.Type)
				}

				if len(data) != 0 {
					df, err := os.OpenFile(dt.File, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0755)
					if err != nil {
						logger.Fatalf(`open file "%s" failed. %v`, dt.File, err)
					}
					defer df.Close()

					df.Write(data)
				}

				swg.Done()
			}(dt)
		}

		swg.Wait()
	})
}