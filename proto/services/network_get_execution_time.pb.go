// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.35.2
// 	protoc        v4.25.3
// source: network_get_execution_time.proto

package services

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

// *
// Gets the time in nanoseconds spent in <tt>handleTransaction</tt> for one or more
// TransactionIDs (assuming they have reached consensus "recently", since only a limited
// number of execution times are kept in-memory, depending on the value of the node-local
// property <tt>stats.executionTimesToTrack</tt>).
type NetworkGetExecutionTimeQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// standard info sent from client to node including the signed payment, and what kind of response
	// is requested (cost, state proof, both, or neither).
	Header *QueryHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// *
	// The id(s) of the transactions to get the execution time(s) of
	TransactionIds []*TransactionID `protobuf:"bytes,2,rep,name=transaction_ids,json=transactionIds,proto3" json:"transaction_ids,omitempty"`
}

func (x *NetworkGetExecutionTimeQuery) Reset() {
	*x = NetworkGetExecutionTimeQuery{}
	mi := &file_network_get_execution_time_proto_msgTypes[0]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NetworkGetExecutionTimeQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkGetExecutionTimeQuery) ProtoMessage() {}

func (x *NetworkGetExecutionTimeQuery) ProtoReflect() protoreflect.Message {
	mi := &file_network_get_execution_time_proto_msgTypes[0]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkGetExecutionTimeQuery.ProtoReflect.Descriptor instead.
func (*NetworkGetExecutionTimeQuery) Descriptor() ([]byte, []int) {
	return file_network_get_execution_time_proto_rawDescGZIP(), []int{0}
}

func (x *NetworkGetExecutionTimeQuery) GetHeader() *QueryHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *NetworkGetExecutionTimeQuery) GetTransactionIds() []*TransactionID {
	if x != nil {
		return x.TransactionIds
	}
	return nil
}

// *
// Response when the client sends the node NetworkGetExecutionTimeQuery; returns
// INVALID_TRANSACTION_ID if any of the given TransactionIDs do not have available
// execution times in the answering node.
type NetworkGetExecutionTimeResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	// *
	// Standard response from node to client, including the requested fields: cost, or state proof,
	// or both, or neither
	Header *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"`
	// *
	// The execution time(s) of the requested TransactionIDs, if available
	ExecutionTimes []uint64 `protobuf:"varint,2,rep,packed,name=execution_times,json=executionTimes,proto3" json:"execution_times,omitempty"`
}

func (x *NetworkGetExecutionTimeResponse) Reset() {
	*x = NetworkGetExecutionTimeResponse{}
	mi := &file_network_get_execution_time_proto_msgTypes[1]
	ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
	ms.StoreMessageInfo(mi)
}

func (x *NetworkGetExecutionTimeResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NetworkGetExecutionTimeResponse) ProtoMessage() {}

func (x *NetworkGetExecutionTimeResponse) ProtoReflect() protoreflect.Message {
	mi := &file_network_get_execution_time_proto_msgTypes[1]
	if x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NetworkGetExecutionTimeResponse.ProtoReflect.Descriptor instead.
func (*NetworkGetExecutionTimeResponse) Descriptor() ([]byte, []int) {
	return file_network_get_execution_time_proto_rawDescGZIP(), []int{1}
}

func (x *NetworkGetExecutionTimeResponse) GetHeader() *ResponseHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *NetworkGetExecutionTimeResponse) GetExecutionTimes() []uint64 {
	if x != nil {
		return x.ExecutionTimes
	}
	return nil
}

var File_network_get_execution_time_proto protoreflect.FileDescriptor

var file_network_get_execution_time_proto_rawDesc = []byte{
	0x0a, 0x20, 0x6e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x5f, 0x67, 0x65, 0x74, 0x5f, 0x65, 0x78,
	0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x2e, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x11, 0x62, 0x61, 0x73, 0x69, 0x63,
	0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x12, 0x71, 0x75,
	0x65, 0x72, 0x79, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x15, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x89, 0x01, 0x0a, 0x1c, 0x4e, 0x65, 0x74, 0x77,
	0x6f, 0x72, 0x6b, 0x47, 0x65, 0x74, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x54,
	0x69, 0x6d, 0x65, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x2a, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64,
	0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x12, 0x3d, 0x0a, 0x0f, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74,
	0x69, 0x6f, 0x6e, 0x5f, 0x69, 0x64, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x14, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f,
	0x6e, 0x49, 0x44, 0x52, 0x0e, 0x74, 0x72, 0x61, 0x6e, 0x73, 0x61, 0x63, 0x74, 0x69, 0x6f, 0x6e,
	0x49, 0x64, 0x73, 0x22, 0x79, 0x0a, 0x1f, 0x4e, 0x65, 0x74, 0x77, 0x6f, 0x72, 0x6b, 0x47, 0x65,
	0x74, 0x45, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x52, 0x65,
	0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12, 0x2d, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72,
	0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52,
	0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68,
	0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x27, 0x0a, 0x0f, 0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69,
	0x6f, 0x6e, 0x5f, 0x74, 0x69, 0x6d, 0x65, 0x73, 0x18, 0x02, 0x20, 0x03, 0x28, 0x04, 0x52, 0x0e,
	0x65, 0x78, 0x65, 0x63, 0x75, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x42, 0x26,
	0x0a, 0x22, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64, 0x65, 0x72, 0x61, 0x68, 0x61, 0x73, 0x68,
	0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x61, 0x70, 0x69, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e,
	0x6a, 0x61, 0x76, 0x61, 0x50, 0x01, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_network_get_execution_time_proto_rawDescOnce sync.Once
	file_network_get_execution_time_proto_rawDescData = file_network_get_execution_time_proto_rawDesc
)

func file_network_get_execution_time_proto_rawDescGZIP() []byte {
	file_network_get_execution_time_proto_rawDescOnce.Do(func() {
		file_network_get_execution_time_proto_rawDescData = protoimpl.X.CompressGZIP(file_network_get_execution_time_proto_rawDescData)
	})
	return file_network_get_execution_time_proto_rawDescData
}

var file_network_get_execution_time_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_network_get_execution_time_proto_goTypes = []any{
	(*NetworkGetExecutionTimeQuery)(nil),    // 0: proto.NetworkGetExecutionTimeQuery
	(*NetworkGetExecutionTimeResponse)(nil), // 1: proto.NetworkGetExecutionTimeResponse
	(*QueryHeader)(nil),                     // 2: proto.QueryHeader
	(*TransactionID)(nil),                   // 3: proto.TransactionID
	(*ResponseHeader)(nil),                  // 4: proto.ResponseHeader
}
var file_network_get_execution_time_proto_depIdxs = []int32{
	2, // 0: proto.NetworkGetExecutionTimeQuery.header:type_name -> proto.QueryHeader
	3, // 1: proto.NetworkGetExecutionTimeQuery.transaction_ids:type_name -> proto.TransactionID
	4, // 2: proto.NetworkGetExecutionTimeResponse.header:type_name -> proto.ResponseHeader
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_network_get_execution_time_proto_init() }
func file_network_get_execution_time_proto_init() {
	if File_network_get_execution_time_proto != nil {
		return
	}
	file_basic_types_proto_init()
	file_query_header_proto_init()
	file_response_header_proto_init()
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_network_get_execution_time_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_network_get_execution_time_proto_goTypes,
		DependencyIndexes: file_network_get_execution_time_proto_depIdxs,
		MessageInfos:      file_network_get_execution_time_proto_msgTypes,
	}.Build()
	File_network_get_execution_time_proto = out.File
	file_network_get_execution_time_proto_rawDesc = nil
	file_network_get_execution_time_proto_goTypes = nil
	file_network_get_execution_time_proto_depIdxs = nil
}
