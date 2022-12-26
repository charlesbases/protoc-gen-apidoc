package conf

import (
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
	_ "gopkg.in/yaml.v3"

	"github.com/charlesbases/protoc-gen-apidoc/logger"
)

const (
	argConfigfile arg = "configfile"

	defaultConfigfile = "apidoc.yaml"
)

// fileOptions .
type fileOptions struct {
	configfile string
}

// newFileParser .
func newFileParser(args string) parser {
	var opts = &fileOptions{configfile: defaultConfigfile}

	if len(args) != 0 {
		for _, param := range strings.Split(args, ",") {
			var value string
			if i := strings.Index(param, "="); i >= 0 {
				value = param[i+1:]
				param = param[0:i]
			}

			switch arg(param) {
			// 配置文件
			case argConfigfile:
				opts.configfile = value
			}
		}
	}

	return opts
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
