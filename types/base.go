package types

import (
	"sort"
	"sync"

	"google.golang.org/protobuf/types/descriptorpb"
)

type (
	Package struct {
		enumLocker sync.RWMutex
		messLocker sync.RWMutex

		// Name Package.Name
		Name string
		// Version version
		Version string
		// Prefix uri prefix
		Prefix string
		// Services Service list
		Services []*Service
		// Enums Enum list
		Enums []*Enum
		// EnumDic Enum map
		EnumDic map[string]*Enum
		// Messages Message list
		Messages []*Message
		// MessageDic Message map
		MessageDic map[string]*Message
	}

	Service struct {
		Name        string
		Description string
		// Methods rpc list
		Methods []*ServiceMethod
	}

	// ServiceMethod service.rpc
	ServiceMethod struct {
		Name         string
		Path         string
		Method       Method
		Description  string
		Consume      ContentType
		Produce      ContentType
		RequestName  string
		ResponseName string
	}

	Enum struct {
		Name        string
		Description string
		Fields      []*EnumField
	}

	EnumField struct {
		Name        string
		Value       int32
		Description string
	}

	Message struct {
		Name        string
		Description string
		Fields      []*MessageField
	}

	MessageField struct {
		// MessageName Message.Name
		MessageName string
		// Description field description
		Description string

		ProtoName        string                                  // proto field name
		ProtoType        descriptorpb.FieldDescriptorProto_Type  // 隐式类型
		ProtoLaber       descriptorpb.FieldDescriptorProto_Label // proto 标签
		ProtoTypeName    string                                  // 显示类型
		ProtoFullName    string                                  // 包名.结构名
		ProtoPackagePath string                                  // 包路径
		ProtoNumber      int32                                   // 排序

		JsonName         string      // json field name
		JsonType         JsonType    // json 类型
		JsonLabel        JsonLabel   // json 标签
		JsonDefaultValue interface{} // json 数据默认值
	}
)

// ascending 升序
func ascending(l, r string) bool {
	var length int
	if len(l) < len(r) {
		length = len(l)
	} else {
		length = len(r)
	}

	for i := 0; i < length; i++ {
		if l[i] != r[i] {
			return l[i] < r[i]
		}
	}
	return len(l) < len(r)
}

// descending 降序
func descending(l, r string) bool {
	var length int
	if len(l) < len(r) {
		length = len(l)
	} else {
		length = len(r)
	}

	for i := 0; i < length; i++ {
		if l[i] != r[i] {
			return l[i] > r[i]
		}
	}
	return len(l) > len(r)
}

// AppendEnum .
func (p *Package) AppendEnum(def *Enum) {
	p.enumLocker.Lock()
	if _, found := p.EnumDic[def.Name]; !found {
		p.EnumDic[def.Name] = def
		p.Enums = append(p.Enums, def)
	}
	p.enumLocker.Unlock()
}

// AppendMessage .
func (p *Package) AppendMessage(def *Message) {
	p.messLocker.Lock()
	if _, found := p.MessageDic[def.Name]; !found {
		p.MessageDic[def.Name] = def
		p.Messages = append(p.Messages, def)
	}
	p.messLocker.Unlock()
}

// Sort .
func (p *Package) Sort() *Package {
	var swg = sync.WaitGroup{}
	swg.Add(3)

	// Services
	go func() {
		if len(p.Services) != 0 {
			sort.Slice(p.Services, func(i, j int) bool {
				return ascending(p.Services[i].Name, p.Services[j].Name)
			})
		}

		swg.Done()
	}()

	// Messages
	go func() {
		if len(p.Messages) != 0 {
			sort.Slice(p.Messages, func(i, j int) bool {
				return ascending(p.Messages[i].Name, p.Messages[j].Name)
			})
		}

		swg.Done()
	}()

	// Enums
	go func() {
		if len(p.Enums) != 0 {
			sort.Slice(p.Enums, func(i, j int) bool {
				return ascending(p.Enums[i].Name, p.Enums[j].Name)
			})
		}

		swg.Done()
	}()

	swg.Wait()
	return p
}
