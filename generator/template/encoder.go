package template

import (
	"fmt"
	"html/template"
	"strings"

	"github.com/charlesbases/protoc-gen-apidoc/protoc"
	"github.com/charlesbases/protoc-gen-apidoc/types"
	"google.golang.org/protobuf/types/descriptorpb"
)

// encoder .
type encoder struct {
	p *types.Package

	// nested 嵌套结构，防止序列化时堆栈溢出
	nesteds map[string]int
}

// indent .
func (e *encoder) indent(layer int) string {
	return strings.Repeat("  ", layer)
}

// marshal .
func (e *encoder) marshal(name string) template.HTML {
	e.nesteds = make(map[string]int, len(e.p.Messages))

	if mess, found := e.p.MessageDic[name]; found && len(mess.Fields) != 0 {
		var br strings.Builder
		br.WriteString("{\n")
		br.WriteString(e.encodeMessage(mess, 1))
		br.WriteString("\n}")
		return template.HTML(br.String())
	}
	return "null"
}

// encodeMessage .
func (e *encoder) encodeMessage(mess *types.Message, layer int) string {
	indent := e.indent(layer)

	var br strings.Builder
	for idx, field := range mess.Fields {
		switch field.JsonType {
		case types.JsonType_Object:
			switch field.ProtoType {
			case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
				value := e.encodeEnum(field.ProtoTypeName)

				switch field.JsonLabel {
				case types.JsonLabel_Repeated:
					br.WriteString(fmt.Sprintf(`%s"%s": [`, indent, field.JsonName))
					br.WriteString(fmt.Sprintf("\n%s%v\n", e.indent(layer+1), value))
					br.WriteString(fmt.Sprintf("%s]", indent))
				default:
					br.WriteString(fmt.Sprintf(`%s"%s": "%s"`, indent, field.JsonName, value))
				}
			case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
				// 预防同名结构体嵌套导致 goroutine 堆栈字节溢出
				if nesteds := fmt.Sprintf("%s.%s", field.MessageName, field.ProtoName); e.nesteds[nesteds] == 2 {
					br.WriteString(fmt.Sprintf(`%s"%s": null`, indent, field.JsonName))
					goto ending
				} else {
					e.nesteds[nesteds]++
				}

				if nesteds, found := e.p.MessageDic[field.ProtoTypeName]; found {
					switch field.JsonLabel {
					case types.JsonLabel_Repeated:
						br.WriteString(fmt.Sprintf(`%s"%s": [`, indent, field.JsonName))
						br.WriteString(fmt.Sprintf("\n%s{\n", e.indent(layer+1)))

						// 是否为 EntryMessage
						if protoc.IsEntry(field) && len(nesteds.Fields) == 2 {
							// nesteds.Fields[0]: key field
							// nesteds.Fields[1]: value field
							entryDemo1, entryDemo2 := new(types.MessageField), new(types.MessageField)
							*entryDemo1, *entryDemo2 = *nesteds.Fields[1], *nesteds.Fields[1]
							entryDemo1.JsonName, entryDemo2.JsonName = "key1", "key2"

							var entryMessage = &types.Message{Fields: []*types.MessageField{entryDemo1, entryDemo2}}
							br.WriteString(e.encodeMessage(entryMessage, layer+2))
						} else {
							br.WriteString(e.encodeMessage(nesteds, layer+2))
						}

						br.WriteString(fmt.Sprintf("\n%s}\n", e.indent(layer+1)))
						br.WriteString(fmt.Sprintf(`%s]`, indent))
					default:
						br.WriteString(fmt.Sprintf(`%s"%s": {`, indent, field.JsonName))
						br.WriteString("\n")

						// 是否为 EntryMessage
						if protoc.IsEntry(field) && len(nesteds.Fields) == 2 {
							// nesteds.Fields[0]: key field
							// nesteds.Fields[1]: value field
							entryDemo1, entryDemo2 := new(types.MessageField), new(types.MessageField)
							*entryDemo1, *entryDemo2 = *nesteds.Fields[1], *nesteds.Fields[1]
							entryDemo1.JsonName, entryDemo2.JsonName = "key1", "key2"

							var entryMessage = &types.Message{Fields: []*types.MessageField{entryDemo1, entryDemo2}}
							br.WriteString(e.encodeMessage(entryMessage, layer+1))
						} else {
							br.WriteString(e.encodeMessage(nesteds, layer+1))
						}

						br.WriteString("\n")
						br.WriteString(fmt.Sprintf(`%s}`, indent))
					}
				}
			}
		default:
			switch field.JsonLabel {
			case types.JsonLabel_Repeated:
				br.WriteString(fmt.Sprintf(`%s"%s": [`, indent, field.JsonName))
				br.WriteString(fmt.Sprintf("\n%s%v\n", e.indent(layer+1), field.JsonDefaultValue))
				br.WriteString(fmt.Sprintf(`%s]`, indent))
			default:
				br.WriteString(fmt.Sprintf(`%s"%s": %v`, indent, field.JsonName, field.JsonDefaultValue))
			}
		}
	ending:
		if idx != len(mess.Fields)-1 {
			br.WriteString(",\n")
		}
	}
	return br.String()
}

// encodeEnum .
func (e *encoder) encodeEnum(enumName string) string {
	if enum, found := e.p.EnumDic[enumName]; found && len(enum.Fields) != 0 {
		return enum.Fields[0].Name

	}
	return ""
}
