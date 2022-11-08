package template

const Markdown Template = `# {{$packagename := .Name -}}
Package {{$packagename}}

---
## 导航 <a name="top"> </a>
+ [服务](#srv)
+ [结构](#msg)
+ [枚举](#enu)
---

## 服务 <a name="srv"> </a>

{{range $serviceindex, $service := .Services -}}
+ ###### {{$service.Name}}  [{{$service.Description}}]
  {{range $apiindex, $method := $service.Methods -}}
  + [{{$method.Path}}](#{{$service.Name}}.{{$method.Name}}){{dynamic $method.Path}}[{{$method.Description}}]
  {{end}}
{{end}}
---

## 接口
{{range $serviceindex, $service := .Services -}}
{{range $apiindex, $method := $service.Methods -}}
#### {{$method.Path}} <a name="{{$service.Name}}.{{$method.Name}}"> </a> [服务](#srv) [结构](#msg) [枚举](#enu)
{{codeblock}}
描述: {{$method.Description}}
{{codeblock}}
+ 请求

{{$message := getMessage $method.RequestName -}}
| 字段 | 类型 | 标签 | 描述 |
| :----------------------: | :---------------------: | :----------------------: | :----------------------: |
{{range $fieldindex, $field := $message.Fields -}}
| {{$field.JsonName}} | [{{jsonType $field}}](#{{$field.ProtoTypeName}}) | {{$field.JsonLabel}} | {{$field.Description}} |
{{end}}
**示例**
{{codeblock "json"}}
{{jsonMarshal $message.Name}}
{{codeblock}}
+ 响应

{{$message := getMessage $method.ResponseName -}}
| 字段 | 类型 | 标签 | 描述 |
| :----------------------: | :---------------------: | :----------------------: | :----------------------: |
{{range $fieldindex, $field := $message.Fields -}}
| {{$field.JsonName}} | [{{jsonType $field}}](#{{$field.ProtoTypeName}}) | {{$field.JsonLabel}} | {{$field.Description}} |
{{end}}
**示例**
{{codeblock "json"}}
{{jsonMarshal $message.Name}}
{{codeblock}}
---
{{end}}
{{end}}

## 结构 <a name="msg"> </a>

| 类型 | 描述 |
| :----------------------: | :---------------------: |
{{range $messageindex, $message := .Messages -}}
| [{{$message.Name}}](#{{$message.Name}}) | {{$message.Description}} |
{{end}}
---
{{range $messageindex, $message := .Messages -}}
+ ##### {{$message.Name}} <a name="{{$message.Name}}"> </a> [服务](#srv) [结构](#msg) [枚举](#enu)
{{codeblock}}
描述: {{$message.Description}}
{{codeblock}}

| 字段 | 类型 | 标签 | 描述 |
| :----------------------: | :---------------------: | :----------------------: | :----------------------: |
{{range $fieldindex, $field := $message.Fields -}}
| {{$field.JsonName}} | [{{jsonType $field}}](#{{$field.ProtoTypeName}}) | {{$field.JsonLabel}} | {{$field.Description}} |
{{end}}
{{end}}

---
## 枚举 <a name="enu"> </a>

{{range $enumindex, $enum := .Enums -}}
+ ##### {{$enum.Name}} <a name="{{$enum.Name}}"> </a> [服务](#srv) [结构](#msg) [枚举](#enu)
| 键 | 值 | 描述 |
| :--------------------: | :--------------------: | :---------------------: |
{{range $fieldindex, $field := $enum.Fields -}}
| {{$field.Name}} | {{$field.Value}} | {{$enum.Description}}:    {{$field.Description}} |
{{end}}
{{end}}
---
`
