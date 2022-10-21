package conf

import (
	"strings"

	"github.com/charlesbases/protoc-gen-apidoc/types"
)

type arg string

const (
	// argConfigfile 配置文件路径
	argConfigfile arg = "configfile"

	// defaultAPIHost .
	defaultAPIHost = "127.0.0.1"
)

var config *configuration

// configuration .
type configuration struct {
	Host     string         `yaml:"host"`
	Title    string         `yaml:"title"`
	Header   []types.Header `yaml:"header"`
	Document []*Document    `yaml:"document"`
}

// Document .
type Document struct {
	Type types.DocumentType `yaml:"type"`
	File string             `yaml:"file"`
}

// parser 配置解析器
type parser interface {
	parse() *configuration
}

// Parse .
func Parse(args string) {
	var configfile = defaultConfigfile

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
				configfile = value
			}
		}
	}

	// 配置文件解析
	config = fileParser(configfile).parse()

	// Default
	if len(config.Host) == 0 {
		config.Host = defaultAPIHost
	}
}

// Get .
func Get() *configuration {
	return config
}
