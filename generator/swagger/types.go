package swagger

import (
	"github.com/charlesbases/protoc-gen-apidoc/types"
	"google.golang.org/protobuf/types/descriptorpb"
)

var prototypes = map[descriptorpb.FieldDescriptorProto_Type]*Definition{
	descriptorpb.FieldDescriptorProto_TYPE_BYTES: {
		Type:   "string",
		Format: "bytes",
	},
	descriptorpb.FieldDescriptorProto_TYPE_STRING: {
		Type:   "string",
		Format: "string",
	},
	descriptorpb.FieldDescriptorProto_TYPE_FLOAT: {
		Type:   "number",
		Format: "float",
	},
	descriptorpb.FieldDescriptorProto_TYPE_DOUBLE: {
		Type:   "number",
		Format: "double",
	},
	descriptorpb.FieldDescriptorProto_TYPE_BOOL: {
		Type:   "boolean",
		Format: "boolean",
	},
	descriptorpb.FieldDescriptorProto_TYPE_INT32: {
		Type:   "integer",
		Format: "int32",
	},
	descriptorpb.FieldDescriptorProto_TYPE_INT64: {
		Type:   "integer",
		Format: "int64",
	},
	descriptorpb.FieldDescriptorProto_TYPE_UINT32: {
		Type:   "integer",
		Format: "uint32",
	},
	descriptorpb.FieldDescriptorProto_TYPE_UINT64: {
		Type:   "integer",
		Format: "uint64",
	},
}

type Position string

const (
	PositionFormData Position = "formData"
	PositionHeader   Position = "header"
	PositionQuery    Position = "query"
	PositionBody     Position = "body"
	PositionPath     Position = "path"
)

// SecurityType type
type SecurityType string

const (
	SecurityTypeBasic  SecurityType = "basic"
	SecurityTypeApiKey SecurityType = "apiKey"
	SecurityTypeOAuth2 SecurityType = "oauth2"
)

// Swagger .
type Swagger struct {
	p *types.Package `json:"-"`

	// Swagger version
	Swagger string `json:"swagger,omitempty"`
	// Info service info
	Info *Info `json:"info,omitempty"`
	// Host service host
	Host string `json:"host,omitempty"`
	// BasePath uri prefix
	BasePath string `json:"basePath,omitempty"`
	// Tags router group list
	Tags []*Tag `json:"tags,omitempty"`
	// Schemes scheme HTTP and HTTPS
	Schemes []string `json:"schemes,omitempty"`
	// Paths api list. map[uri][method]*API
	Paths map[string]map[string]*API `json:"paths,omitempty"`
	// Definitions model list
	Definitions map[string]*Definition `json:"definitions,omitempty"`
	// SecurityDefinitions .
	SecurityDefinitions map[string]*Security `json:"securityDefinitions,omitempty"`
}

// Info service info
type Info struct {
	// Title api title
	Title string `json:"title,omitempty"`
	// Version api version
	Version string `json:"version,omitempty"`
	// Description api description
	Description string `json:"description,omitempty"`
}

// Tag router group
type Tag struct {
	// Name group name
	Name string `json:"name,omitempty"`
	// Description tag description
	Description string `json:"description,omitempty"`
}

// Definition model
type Definition struct {
	// Name name
	Name string `json:"-"`
	// Type json type
	Type string `json:"type,omitempty"`
	// Description description
	Description string `json:"description,omitempty"`

	// Format data type
	Format string `json:"format,omitempty"`

	// Enum enum keys
	Enum []string `json:"enum,omitempty"`
	// Default enum default
	Default string `json:"default,omitempty"`

	// Reflex others Definition point
	Reflex string `json:"$ref,omitempty"`

	// Items array info
	Items *Definition `json:"items,omitempty"`

	// Entry proto entry type
	Entry *Definition `json:"additionalProperties,omitempty"`

	// Nesteds nested
	Nesteds map[string]*Definition `json:"properties,omitempty"`
}

// Security .
type Security struct {
	Type SecurityType `json:"type,omitempty"`
	Name types.Header `json:"name,omitempty"`
	In   Position     `json:"in,omitempty"`
}

// API .
type API struct {
	// Tags tag name list
	Tags []string `json:"tags,omitempty"`
	// Summary summary
	Summary string `json:"summary,omitempty"`
	// Description description
	Description string `json:"description,omitempty"`
	// OperationID operationId
	OperationID string `json:"operationId,omitempty"`
	// Consumes request ContentType
	Consumes []types.ContentType `json:"consumes,omitempty"`
	// Produces response ContentType
	Produces []types.ContentType `json:"produces,omitempty"`
	// Parameters request
	Parameters []*Parameter `json:"parameters,omitempty"`
	// Responses response
	Responses map[string]*Parameter `json:"responses,omitempty"`
}

// Parameter .
type Parameter struct {
	In       Position `json:"in,omitempty"`
	Name     string   `json:"name,omitempty"`
	Type     string   `json:"type,omitempty"`
	Required bool     `json:"required,omitempty"`
	// Enum enum keys
	Enum []string `json:"enum,omitempty"`
	// Default default value
	Default string `json:"default,omitempty"`
	// Description description
	Description string `json:"description,omitempty"`
	// Schema Definition path
	Schema *Definition `json:"schema,omitempty"`
	// Items array info
	Items *Definition `json:"items,omitempty"`
}
