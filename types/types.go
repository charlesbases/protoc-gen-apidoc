package types

import (
	"google.golang.org/protobuf/types/descriptorpb"
)

// DocumentType 文档类型
type DocumentType string

const (
	DocumentType_HTML     DocumentType = "html"
	DocumentType_Markdown DocumentType = "markdown"
	DocumentType_Postman  DocumentType = "postman"
	DocumentType_Swagger  DocumentType = "swagger"
)

type ContentType string

const (
	ContentType_Json ContentType = "application/json"
	ContentType_Data ContentType = "multipart/form-data"
)

type Method string

const (
	Method_Get    Method = "GET"
	Method_Put    Method = "PUT"
	Method_Post   Method = "POST"
	Method_Delete Method = "DELETE"
)

// String .
func (m Method) String() string {
	return string(m)
}

// LowerCase .
func (m Method) LowerCase() string {
	switch m {
	case Method_Get:
		return "get"
	case Method_Put:
		return "put"
	case Method_Post:
		return "post"
	case Method_Delete:
		return "delete"
	default:
		return ""
	}
}

type Header string

// String .
func (h Header) String() string {
	return string(h)
}

// Desc .
func (h Header) Desc() string {
	return string(h) + " In Header"
}

type JsonType string

const (
	JsonType_Object  JsonType = "Object"
	JsonType_Number  JsonType = "Number"
	JsonType_String  JsonType = "String"
	JsonType_Boolean JsonType = "Boolean"
)

type JsonLabel string

const (
	JsonLabel_Optional JsonLabel = "可选"
	JsonLabel_Required JsonLabel = "必须"
	JsonLabel_Repeated JsonLabel = "重复"
)

// DefaultValue .
func (jt JsonType) DefaultValue() interface{} {
	switch jt {
	case JsonType_Number:
		return 0
	case JsonType_String:
		return `"string"`
	case JsonType_Boolean:
		return false
	default:
		return nil
	}
}

// Convert2JsonType descriptorpb.FieldDescriptorProto_Type to JsonType
func Convert2JsonType(pt descriptorpb.FieldDescriptorProto_Type) JsonType {
	switch pt {
	case
		descriptorpb.FieldDescriptorProto_TYPE_DOUBLE,
		descriptorpb.FieldDescriptorProto_TYPE_FLOAT,
		descriptorpb.FieldDescriptorProto_TYPE_INT64,
		descriptorpb.FieldDescriptorProto_TYPE_UINT64,
		descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_FIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_UINT32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED32,
		descriptorpb.FieldDescriptorProto_TYPE_SFIXED64,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT64:
		return JsonType_Number
	case
		descriptorpb.FieldDescriptorProto_TYPE_GROUP,
		descriptorpb.FieldDescriptorProto_TYPE_MESSAGE,
		descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		return JsonType_Object
	case
		descriptorpb.FieldDescriptorProto_TYPE_STRING,
		descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		return JsonType_String
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return JsonType_Boolean
	default:
		return ""
	}
}

// Convert2JsonLabel descriptorpb.FieldDescriptorProto_Label to JsonLabel
func Convert2JsonLabel(pl descriptorpb.FieldDescriptorProto_Label) JsonLabel {
	switch pl {
	case descriptorpb.FieldDescriptorProto_LABEL_OPTIONAL:
		return JsonLabel_Optional
	case descriptorpb.FieldDescriptorProto_LABEL_REQUIRED:
		return JsonLabel_Required
	case descriptorpb.FieldDescriptorProto_LABEL_REPEATED:
		return JsonLabel_Repeated
	default:
		return ""
	}
}
