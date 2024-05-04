// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.33.0
// 	protoc        v5.26.1
// source: streaming/enum.proto

package pb

import (
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ResponseStatus int32

const (
	ResponseStatus_STATUS_UNDEFINED            ResponseStatus = 0
	ResponseStatus_STATUS_INIT                 ResponseStatus = 1
	ResponseStatus_STATUS_FRAMER_PROCESSING    ResponseStatus = 2
	ResponseStatus_STATUS_FRAMER_ERROR         ResponseStatus = 3
	ResponseStatus_STATUS_DETECTION_PROCESSING ResponseStatus = 4
	ResponseStatus_STATUS_DETECTION_ERROR      ResponseStatus = 5
	ResponseStatus_STATUS_SUCCESS              ResponseStatus = 6
	ResponseStatus_STATUS_ERROR                ResponseStatus = 7
	ResponseStatus_STATUS_CANCELED             ResponseStatus = 8
	ResponseStatus_STATUS_RESPONSER_PROCESSING ResponseStatus = 9
)

// Enum value maps for ResponseStatus.
var (
	ResponseStatus_name = map[int32]string{
		0: "STATUS_UNDEFINED",
		1: "STATUS_INIT",
		2: "STATUS_FRAMER_PROCESSING",
		3: "STATUS_FRAMER_ERROR",
		4: "STATUS_DETECTION_PROCESSING",
		5: "STATUS_DETECTION_ERROR",
		6: "STATUS_SUCCESS",
		7: "STATUS_ERROR",
		8: "STATUS_CANCELED",
		9: "STATUS_RESPONSER_PROCESSING",
	}
	ResponseStatus_value = map[string]int32{
		"STATUS_UNDEFINED":            0,
		"STATUS_INIT":                 1,
		"STATUS_FRAMER_PROCESSING":    2,
		"STATUS_FRAMER_ERROR":         3,
		"STATUS_DETECTION_PROCESSING": 4,
		"STATUS_DETECTION_ERROR":      5,
		"STATUS_SUCCESS":              6,
		"STATUS_ERROR":                7,
		"STATUS_CANCELED":             8,
		"STATUS_RESPONSER_PROCESSING": 9,
	}
)

func (x ResponseStatus) Enum() *ResponseStatus {
	p := new(ResponseStatus)
	*p = x
	return p
}

func (x ResponseStatus) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (ResponseStatus) Descriptor() protoreflect.EnumDescriptor {
	return file_streaming_enum_proto_enumTypes[0].Descriptor()
}

func (ResponseStatus) Type() protoreflect.EnumType {
	return &file_streaming_enum_proto_enumTypes[0]
}

func (x ResponseStatus) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use ResponseStatus.Descriptor instead.
func (ResponseStatus) EnumDescriptor() ([]byte, []int) {
	return file_streaming_enum_proto_rawDescGZIP(), []int{0}
}

type QueryType int32

const (
	QueryType_TYPE_UNDEFINED QueryType = 0
	QueryType_TYPE_VIDEO     QueryType = 1
	QueryType_TYPE_STREAM    QueryType = 2
)

// Enum value maps for QueryType.
var (
	QueryType_name = map[int32]string{
		0: "TYPE_UNDEFINED",
		1: "TYPE_VIDEO",
		2: "TYPE_STREAM",
	}
	QueryType_value = map[string]int32{
		"TYPE_UNDEFINED": 0,
		"TYPE_VIDEO":     1,
		"TYPE_STREAM":    2,
	}
)

func (x QueryType) Enum() *QueryType {
	p := new(QueryType)
	*p = x
	return p
}

func (x QueryType) String() string {
	return protoimpl.X.EnumStringOf(x.Descriptor(), protoreflect.EnumNumber(x))
}

func (QueryType) Descriptor() protoreflect.EnumDescriptor {
	return file_streaming_enum_proto_enumTypes[1].Descriptor()
}

func (QueryType) Type() protoreflect.EnumType {
	return &file_streaming_enum_proto_enumTypes[1]
}

func (x QueryType) Number() protoreflect.EnumNumber {
	return protoreflect.EnumNumber(x)
}

// Deprecated: Use QueryType.Descriptor instead.
func (QueryType) EnumDescriptor() ([]byte, []int) {
	return file_streaming_enum_proto_rawDescGZIP(), []int{1}
}

var File_streaming_enum_proto protoreflect.FileDescriptor

var file_streaming_enum_proto_rawDesc = []byte{
	0x0a, 0x14, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e, 0x67, 0x2f, 0x65, 0x6e, 0x75, 0x6d,
	0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12, 0x09, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d, 0x69, 0x6e,
	0x67, 0x2a, 0x87, 0x02, 0x0a, 0x0e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x53, 0x74,
	0x61, 0x74, 0x75, 0x73, 0x12, 0x14, 0x0a, 0x10, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x55,
	0x4e, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0f, 0x0a, 0x0b, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x49, 0x4e, 0x49, 0x54, 0x10, 0x01, 0x12, 0x1c, 0x0a, 0x18, 0x53,
	0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x46, 0x52, 0x41, 0x4d, 0x45, 0x52, 0x5f, 0x50, 0x52, 0x4f,
	0x43, 0x45, 0x53, 0x53, 0x49, 0x4e, 0x47, 0x10, 0x02, 0x12, 0x17, 0x0a, 0x13, 0x53, 0x54, 0x41,
	0x54, 0x55, 0x53, 0x5f, 0x46, 0x52, 0x41, 0x4d, 0x45, 0x52, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52,
	0x10, 0x03, 0x12, 0x1f, 0x0a, 0x1b, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x45, 0x54,
	0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x50, 0x52, 0x4f, 0x43, 0x45, 0x53, 0x53, 0x49, 0x4e,
	0x47, 0x10, 0x04, 0x12, 0x1a, 0x0a, 0x16, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x44, 0x45,
	0x54, 0x45, 0x43, 0x54, 0x49, 0x4f, 0x4e, 0x5f, 0x45, 0x52, 0x52, 0x4f, 0x52, 0x10, 0x05, 0x12,
	0x12, 0x0a, 0x0e, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x53, 0x55, 0x43, 0x43, 0x45, 0x53,
	0x53, 0x10, 0x06, 0x12, 0x10, 0x0a, 0x0c, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f, 0x45, 0x52,
	0x52, 0x4f, 0x52, 0x10, 0x07, 0x12, 0x13, 0x0a, 0x0f, 0x53, 0x54, 0x41, 0x54, 0x55, 0x53, 0x5f,
	0x43, 0x41, 0x4e, 0x43, 0x45, 0x4c, 0x45, 0x44, 0x10, 0x08, 0x12, 0x1f, 0x0a, 0x1b, 0x53, 0x54,
	0x41, 0x54, 0x55, 0x53, 0x5f, 0x52, 0x45, 0x53, 0x50, 0x4f, 0x4e, 0x53, 0x45, 0x52, 0x5f, 0x50,
	0x52, 0x4f, 0x43, 0x45, 0x53, 0x53, 0x49, 0x4e, 0x47, 0x10, 0x09, 0x2a, 0x40, 0x0a, 0x09, 0x51,
	0x75, 0x65, 0x72, 0x79, 0x54, 0x79, 0x70, 0x65, 0x12, 0x12, 0x0a, 0x0e, 0x54, 0x59, 0x50, 0x45,
	0x5f, 0x55, 0x4e, 0x44, 0x45, 0x46, 0x49, 0x4e, 0x45, 0x44, 0x10, 0x00, 0x12, 0x0e, 0x0a, 0x0a,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x56, 0x49, 0x44, 0x45, 0x4f, 0x10, 0x01, 0x12, 0x0f, 0x0a, 0x0b,
	0x54, 0x59, 0x50, 0x45, 0x5f, 0x53, 0x54, 0x52, 0x45, 0x41, 0x4d, 0x10, 0x02, 0x42, 0x17, 0x5a,
	0x15, 0x69, 0x6e, 0x74, 0x65, 0x72, 0x6e, 0x61, 0x6c, 0x2f, 0x73, 0x74, 0x72, 0x65, 0x61, 0x6d,
	0x69, 0x6e, 0x67, 0x2f, 0x70, 0x62, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_streaming_enum_proto_rawDescOnce sync.Once
	file_streaming_enum_proto_rawDescData = file_streaming_enum_proto_rawDesc
)

func file_streaming_enum_proto_rawDescGZIP() []byte {
	file_streaming_enum_proto_rawDescOnce.Do(func() {
		file_streaming_enum_proto_rawDescData = protoimpl.X.CompressGZIP(file_streaming_enum_proto_rawDescData)
	})
	return file_streaming_enum_proto_rawDescData
}

var file_streaming_enum_proto_enumTypes = make([]protoimpl.EnumInfo, 2)
var file_streaming_enum_proto_goTypes = []interface{}{
	(ResponseStatus)(0), // 0: streaming.ResponseStatus
	(QueryType)(0),      // 1: streaming.QueryType
}
var file_streaming_enum_proto_depIdxs = []int32{
	0, // [0:0] is the sub-list for method output_type
	0, // [0:0] is the sub-list for method input_type
	0, // [0:0] is the sub-list for extension type_name
	0, // [0:0] is the sub-list for extension extendee
	0, // [0:0] is the sub-list for field type_name
}

func init() { file_streaming_enum_proto_init() }
func file_streaming_enum_proto_init() {
	if File_streaming_enum_proto != nil {
		return
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_streaming_enum_proto_rawDesc,
			NumEnums:      2,
			NumMessages:   0,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_streaming_enum_proto_goTypes,
		DependencyIndexes: file_streaming_enum_proto_depIdxs,
		EnumInfos:         file_streaming_enum_proto_enumTypes,
	}.Build()
	File_streaming_enum_proto = out.File
	file_streaming_enum_proto_rawDesc = nil
	file_streaming_enum_proto_goTypes = nil
	file_streaming_enum_proto_depIdxs = nil
}