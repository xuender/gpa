package gpa

import (
	"google.golang.org/protobuf/reflect/protoreflect"
)

type Document interface {
	GetId() uint64
	ProtoReflect() protoreflect.Message
}
