package conf

import (
	"strings"

	"github.com/charlesbases/protoc-gen-apidoc/types"
)

type arg string

var config *configuration

// configuration .
type configuration struct {
	Host     string         `yaml:"host"`
	Port     string         `yaml:"port"`
	Title    string         `yaml:"title"`
	Header   []types.Header `yaml:"header"`
	Schemes  []string       `yaml:"schemes"`
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
	// 配置文件解析
	// config = newFileParser(args).parse()
	// 输入参数解析
	config = newArgsParser(args).parse()

	if len(config.Host) != 0 {
		config.Host = strings.ToLower(config.Host)
	}
}

// Get .
func Get() *configuration {
	return config
}
