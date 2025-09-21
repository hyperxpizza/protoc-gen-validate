package module

import (
	"github.com/gogo/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

const (
	validateFieldNumber = 837932
)

var (
	E_Validate = &proto.ExtensionDesc{
		ExtendedType:  (*descriptorpb.FieldOptions)(nil),
		ExtensionType: (*string)(nil),
		Field:         validateFieldNumber,
		Name:          "validate.tag",
		Tag:           "varint,837932,opt,name=ignore",
	}
)
