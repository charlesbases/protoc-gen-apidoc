package protoc

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"

	"github.com/charlesbases/protobuf/types/httppb"
	"github.com/charlesbases/protobuf/types/servicepb"
	"github.com/charlesbases/protoc-gen-apidoc/conf"
	"github.com/charlesbases/protoc-gen-apidoc/logger"
	"github.com/charlesbases/protoc-gen-apidoc/types"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
	"google.golang.org/protobuf/types/pluginpb"
)

// Plugin .
func Plugin(fn func(p *types.Package) *pluginpb.CodeGeneratorResponse) {
	var buff = new(bytes.Buffer)
	if _, err := io.Copy(buff, os.Stdin); err != nil {
		logger.Fatal("read os.Stdin failed. ", err)
	}

	var req = new(pluginpb.CodeGeneratorRequest)
	if err := proto.Unmarshal(buff.Bytes(), req); err != nil {
		logger.Fatal("unmarshal os.Stdin failed. ", err)
	}
	if len(req.GetFileToGenerate()) == 0 {
		logger.Fatal("no file to generate")
	}

	// 解析配置参数
	conf.Parse(req.GetParameter())

	// proto 解析
	stdout(fn(parse(req)))
}

// stdout .
func stdout(rsp *pluginpb.CodeGeneratorResponse) {
	if data, err := proto.Marshal(rsp); err != nil {
		logger.Fatal(err)
	} else {
		os.Stdout.Write(data)
	}
}

// parse 解析 proto 文件
func parse(req *pluginpb.CodeGeneratorRequest) *types.Package {
	var p = newPackage(req.GetProtoFile()[0].GetPackage())

	var swg = sync.WaitGroup{}
	swg.Add(len(req.GetProtoFile()))

	for fidx := range req.GetProtoFile() {
		go func(file *descriptorpb.FileDescriptorProto) {
			if !strings.HasPrefix(file.GetPackage(), "google.protobuf") {
				// parse comment
				var cs = parseComments(file.SourceCodeInfo)

				// parse enum
				for idx, protoEnum := range file.GetEnumType() {
					p.AppendEnum(cs.parseEnum(protoEnum, COMMENT_PATH_ENUM, idx))
				}

				// parse message
				for midx, protoMessage := range file.GetMessageType() {
					var paths = []int{COMMENT_PATH_MESSAGE, midx}

					for eidx, protoEnum := range protoMessage.GetEnumType() {
						p.AppendEnum(cs.parseMessageEnum(protoEnum, protoMessage.GetName(), append(paths, COMMENT_PATH_MESSAGE_ENUM, eidx)...))
					}

					for nidx, protoNested := range protoMessage.GetNestedType() {
						p.AppendMessage(cs.parseMessageNested(protoNested, protoMessage.GetName(), append(paths, COMMENT_PATH_MESSAGE_MESSAGE, nidx)...))
					}

					p.AppendMessage(cs.parseMessage(protoMessage, paths...))
				}

				// parse service
				for idx, protoService := range file.GetService() {
					p.Services = append(p.Services, cs.parseService(protoService, COMMENT_PATH_SERVICE, idx))
				}
			}

			swg.Done()
		}(req.GetProtoFile()[fidx])
	}

	swg.Wait()

	return p.Sort()
}

// parseComments paarse comments in proto
func parseComments(infor *descriptorpb.SourceCodeInfo) comments {
	cs := make(map[string]*comment, 0)

	for _, location := range infor.GetLocation() {
		if location.GetLeadingComments() == "" && location.GetTrailingComments() == "" && len(location.GetLeadingDetachedComments()) == 0 {
			continue
		}

		detached := make([]string, 0)
		for _, val := range location.GetLeadingDetachedComments() {
			detached = append(detached, trim(val, "*", "\n"))
		}

		cs[fmt.Sprintf("%v", location.GetPath())] = &comment{
			leading:  trim(location.GetLeadingComments(), "*", "\n"),
			trailing: trim(location.GetTrailingComments(), "*", "\n"),
			detached: detached,
		}
	}
	return cs
}

// parseservice parse service in proto
func (cs comments) parseService(dsdp *descriptorpb.ServiceDescriptorProto, paths ...int) *types.Service {
	var service = newService(dsdp.GetName(), cs.comment(dsdp.GetName(), paths...))

	// descriptorpb.ServiceOptions
	// if opt := parseServiceOption(dsdp.GetOptions()); opt != nil {
	//
	// }

	for idx, protoRPC := range dsdp.GetMethod() {
		method := cs.parseMethod(protoRPC, append(paths, COMMENT_PATH_SERVICE_METHOD, idx)...)
		if len(method.Path) == 0 {
			method.Path = methodPath(service.Name, method.Name)
		}
		service.Methods = append(service.Methods, method)
	}
	return service
}

// parseMethod parse method in service
func (cs comments) parseMethod(dmdp *descriptorpb.MethodDescriptorProto, paths ...int) *types.ServiceMethod {
	var method = newServiceMethod(dmdp.GetName(), cs.comment(dmdp.GetName(), paths...))
	method.RequestName = split(dmdp.GetInputType())[1]
	method.ResponseName = split(dmdp.GetOutputType())[1]

	// descriptorpb.MethodOptions
	if opt := parseMethodOptions(dmdp.GetOptions()); opt != nil {
		switch opt.GetPattern().(type) {
		case *httppb.Http_Get:
			method.Path = opt.GetGet()
			method.Method = types.Method_Get
		case *httppb.Http_Put:
			method.Path = opt.GetPut()
			method.Method = types.Method_Put
		case *httppb.Http_Post:
			method.Path = opt.GetPost()
			method.Method = types.Method_Post
		case *httppb.Http_Delete:
			method.Path = opt.GetDelete()
			method.Method = types.Method_Delete
		}

		method.Consume = types.ContentType(opt.GetConsume())
		method.Produce = types.ContentType(opt.GetProduce())
	}
	if method.Produce == "" {
		method.Produce = types.ContentType_Json
	}
	if method.Consume == "" && method.Method != types.Method_Get {
		method.Consume = types.ContentType_Json
	}

	return method
}

// parseMessage parse message in proto
func (cs comments) parseMessage(protoMessage *descriptorpb.DescriptorProto, paths ...int) *types.Message {
	var message = newMessage(protoMessage.GetName(), cs.comment(protoMessage.GetName(), paths...))

	for idx, field := range protoMessage.GetField() {
		message.Fields = append(message.Fields, cs.parseMessageField(protoMessage, field, append(paths, COMMENT_PATH_MESSAGE_FIELD, idx)...))
	}
	return message
}

// parseMessageNested parse message nested in message
func (cs comments) parseMessageNested(nested *descriptorpb.DescriptorProto, parent string, paths ...int) *types.Message {
	name := nestedName(parent, nested.GetName())
	var message = newMessage(name, cs.comment(name, paths...))

	for idx, field := range nested.GetField() {
		message.Fields = append(message.Fields, cs.parseMessageField(nested, field, append(paths, COMMENT_PATH_MESSAGE_FIELD, idx)...))
	}
	return message
}

// parseMessageEnum parse enum in message
func (cs comments) parseMessageEnum(protoEnum *descriptorpb.EnumDescriptorProto, parent string, paths ...int) *types.Enum {
	name := nestedName(parent, protoEnum.GetName())
	var enum = newEnum(name, cs.comment(name, paths...))

	for idx, enumField := range protoEnum.GetValue() {
		enum.Fields = append(enum.Fields, cs.parseEnumField(enumField, append(paths, COMMENT_PATH_ENUM_VALUE, idx)...))
	}
	return enum
}

// parseMessageField parse field in message
func (cs comments) parseMessageField(protoMessage *descriptorpb.DescriptorProto, protoField *descriptorpb.FieldDescriptorProto, paths ...int) *types.MessageField {
	var field = &types.MessageField{MessageName: protoMessage.GetName(), Description: cs.comment(protoField.GetName(), paths...)}

	// Json
	field.JsonName = protoField.GetName()
	field.JsonLabel = types.Convert2JsonLabel(protoField.GetLabel())
	field.JsonType = types.Convert2JsonType(protoField.GetType())
	field.JsonDefaultValue = field.JsonType.DefaultValue()

	// Proto
	field.ProtoName = protoField.GetName()
	field.ProtoLaber = protoField.GetLabel()
	field.ProtoType = protoField.GetType()
	field.ProtoNumber = protoField.GetNumber()

	switch field.JsonType {
	case types.JsonType_Object:
		typename := split(protoField.GetTypeName())

		field.ProtoTypeName = typename[1]
		field.ProtoPackagePath = typename[0]
		field.ProtoFullName = protoField.GetTypeName()
	case types.JsonType_Number, types.JsonType_String, types.JsonType_Boolean:
		field.ProtoTypeName = descriptorpb.FieldDescriptorProto_Type_name[int32(field.ProtoType)]
	}

	return field
}

// parseEnum parse enum in proto
func (cs comments) parseEnum(protoEnum *descriptorpb.EnumDescriptorProto, paths ...int) *types.Enum {
	var enum = newEnum(protoEnum.GetName(), cs.comment(protoEnum.GetName(), paths...))

	for idx, enumField := range protoEnum.GetValue() {
		enum.Fields = append(enum.Fields, cs.parseEnumField(enumField, append(paths, COMMENT_PATH_ENUM_VALUE, idx)...))
	}
	return enum
}

// parseEnumField parse field in enum
func (cs comments) parseEnumField(protoEnumField *descriptorpb.EnumValueDescriptorProto, paths ...int) *types.EnumField {
	return &types.EnumField{
		Name:        protoEnumField.GetName(),
		Value:       protoEnumField.GetNumber(),
		Description: cs.comment(protoEnumField.GetName(), paths...),
	}
}

// parseMethodOptions .
func parseMethodOptions(opts *descriptorpb.MethodOptions) *httppb.Http {
	if opts != nil {
		if exp, ok := proto.GetExtension(opts, httppb.E_Http).(*httppb.Http); ok {
			return exp
		}
	}
	return nil
}

// parseServiceOption .
func parseServiceOption(opts *descriptorpb.ServiceOptions) *servicepb.Service {
	if opts != nil {
		if exp, ok := proto.GetExtension(opts, httppb.E_Http).(*servicepb.Service); ok {
			return exp
		}
	}
	return nil
}
