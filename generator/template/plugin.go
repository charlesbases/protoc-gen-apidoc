package template

import (
	"bytes"
	"fmt"
	"html/template"
	"strings"

	"github.com/charlesbases/protoc-gen-apidoc/encoder"
	"github.com/charlesbases/protoc-gen-apidoc/generator"
	"github.com/charlesbases/protoc-gen-apidoc/logger"
	"github.com/charlesbases/protoc-gen-apidoc/types"
)

const width = 66

type Template string

// Generator .
type Generator struct {
	p *types.Package
	t Template
}

// NewGenerator .
func NewGenerator(p *types.Package, t Template) generator.Generator {
	return &Generator{
		p: p,
		t: t,
	}
}

// Generate code generater
func (g *Generator) Generate() []byte {
	temp := template.New(string(g.t))

	temp.Funcs(template.FuncMap{
		"dynamic":     dynamic,
		"codeblock":   codeblock,
		"getMessage":  g.getMessage,
		"jsonType":    g.jsonType,
		"jsonMarshal": g.jsonMarshal,
		"increasing":  g.increasing,
		"polling":     g.polling,
	})

	html, err := temp.Parse(string(g.t))
	if err != nil {
		logger.Fatal(err)
	}

	var buffer bytes.Buffer
	if err := html.Execute(&buffer, g.p); err != nil {
		logger.Fatal(err)
	}

	return buffer.Bytes()
}

// dynamic 动态返回一定长度字符
func dynamic(v string) string {
	return strings.Repeat("·", width-len(v))
}

// codeblock markdown 代码块
func codeblock(languages ...string) template.HTML {
	for x := range languages {
		return template.HTML(fmt.Sprintf("``` %s", languages[x]))
	}
	return template.HTML("```")
}

// polling 判断 index 的奇偶, 偶数返回true, 奇数返回false, 制作条纹表格时需要
func (g *Generator) polling(index int) bool {
	return index%2 == 0
}

// increasing 递增
func (g *Generator) increasing(index int) int {
	return index + 1
}

// getMessage get message information by message name
func (g *Generator) getMessage(name string) *types.Message {
	if mess, found := g.p.MessageDic[name]; found {
		return mess
	}
	return &types.Message{}
}

// jsonType .
func (g *Generator) jsonType(field *types.MessageField) template.HTML {
	switch field.JsonType {
	case types.JsonType_Object:
		return template.HTML(field.ProtoTypeName)
	default:
		return template.HTML(field.JsonType)
	}
}

// jsonMarshal json parse for message
func (g *Generator) jsonMarshal(messageName string) template.HTML {
	if data := encoder.NewEncoder(g.p).EncodeJson(messageName); len(data) != 0 {
		return template.HTML(data)
	}
	return "null"
}
