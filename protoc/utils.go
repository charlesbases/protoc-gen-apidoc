package protoc

import (
	"strings"
	"time"

	"github.com/charlesbases/protoc-gen-apidoc/logger"
	"github.com/charlesbases/protoc-gen-apidoc/types"
	"google.golang.org/protobuf/types/descriptorpb"
)

// version .
func version() string {
	return time.Now().Format("20060102150405")
}

// nestedName message nested name
func nestedName(v ...string) string {
	return strings.Join(v, "_")
}

// methodPath .
func methodPath(v ...string) string {
	return "/" + strings.Join(v, "/")
}

// split split by "." and return package and type name
func split(typename string) [2]string {
	list := strings.Split(typename, ".")
	if len(list) < 3 {
		logger.Fatal("split type failed. ", typename)
	}
	return [2]string{list[1], strings.Join(list[2:], "_")}
}

// trim  prefix and suffix TODO 可优化
func trim(source string, cutsets ...string) string {
	source = strings.TrimSpace(source)

	for _, cutset := range cutsets {
		for {
			if strings.HasPrefix(source, cutset) {
				source = strings.TrimPrefix(source, cutset)
				continue
			}
			break
		}
		for {
			if strings.HasSuffix(source, cutset) {
				source = strings.TrimSuffix(source, cutset)
				continue
			}
			break
		}
	}
	return source
}

// IsEntry 是否为 proto 自动创建的 entry message. 例：map<string, string>
func IsEntry(mf *types.MessageField) bool {
	if mf.ProtoType == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
		var bs strings.Builder
		bs.Grow(len(mf.MessageName) + len(mf.ProtoName) + 6)

		bs.WriteString(mf.MessageName)
		bs.WriteString("_")

		for _, item := range strings.Split(mf.ProtoName, "_") {
			if len(item) > 0 {
				if c := item[0]; 'a' <= c && c <= 'z' {
					bs.WriteByte(c - ('a' - 'A'))
				} else {
					bs.WriteByte(c)
				}
				if len(item) > 1 {
					bs.WriteString(item[1:])
				}
			}
		}
		bs.WriteString("Entry")
		return mf.ProtoTypeName == bs.String()
	}
	return false
}
