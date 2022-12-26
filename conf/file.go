package conf

import (
	"os"
	"path/filepath"

	"gopkg.in/yaml.v3"
	_ "gopkg.in/yaml.v3"

	"github.com/charlesbases/protoc-gen-apidoc/logger"
)

const _defaultConfigfile = "apidoc.yaml"

// fileOptions .
type fileOptions struct {
	configfile string
}

// fileParser .
func fileParser(configfile string) parser {
	return &fileOptions{configfile: configfile}
}

// parse .
func (opts *fileOptions) parse() *configuration {
	abspath, err := filepath.Abs(opts.configfile)
	if err != nil {
		logger.Fatal(err)
	}

	configfile, err := os.Open(abspath)
	if err != nil {
		logger.Fatal(err)
	}
	defer configfile.Close()

	var config = new(configuration)
	if err := yaml.NewDecoder(configfile).Decode(config); err != nil {
		logger.Fatal(err)
	}

	return config
}
