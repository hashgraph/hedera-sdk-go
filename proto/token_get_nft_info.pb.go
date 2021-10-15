// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.27.1
// 	protoc        v3.17.3
// source: proto/token_get_nft_info.proto

package proto

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

// Represents an NFT on the Ledger
type NftID struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	TokenID      *TokenID `protobuf:"bytes,1,opt,name=tokenID,proto3" json:"tokenID,omitempty"`            // The (non-fungible) token of which this NFT is an instance
	SerialNumber int64    `protobuf:"varint,2,opt,name=serialNumber,proto3" json:"serialNumber,omitempty"` // The unique identifier of this instance
}

func (x *NftID) Reset() {
	*x = NftID{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_token_get_nft_info_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *NftID) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*NftID) ProtoMessage() {}

func (x *NftID) ProtoReflect() protoreflect.Message {
	mi := &file_proto_token_get_nft_info_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use NftID.ProtoReflect.Descriptor instead.
func (*NftID) Descriptor() ([]byte, []int) {
	return file_proto_token_get_nft_info_proto_rawDescGZIP(), []int{0}
}

func (x *NftID) GetTokenID() *TokenID {
	if x != nil {
		return x.TokenID
	}
	return nil
}

func (x *NftID) GetSerialNumber() int64 {
	if x != nil {
		return x.SerialNumber
	}
	return 0
}

// Applicable only to tokens of type NON_FUNGIBLE_UNIQUE. Gets info on a NFT for a given TokenID (of type NON_FUNGIBLE_UNIQUE) and serial number
type TokenGetNftInfoQuery struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header *QueryHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"` // Standard info sent from client to node, including the signed payment, and what kind of response is requested (cost, state proof, both, or neither).
	NftID  *NftID       `protobuf:"bytes,2,opt,name=nftID,proto3" json:"nftID,omitempty"`   // The ID of the NFT
}

func (x *TokenGetNftInfoQuery) Reset() {
	*x = TokenGetNftInfoQuery{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_token_get_nft_info_proto_msgTypes[1]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenGetNftInfoQuery) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenGetNftInfoQuery) ProtoMessage() {}

func (x *TokenGetNftInfoQuery) ProtoReflect() protoreflect.Message {
	mi := &file_proto_token_get_nft_info_proto_msgTypes[1]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenGetNftInfoQuery.ProtoReflect.Descriptor instead.
func (*TokenGetNftInfoQuery) Descriptor() ([]byte, []int) {
	return file_proto_token_get_nft_info_proto_rawDescGZIP(), []int{1}
}

func (x *TokenGetNftInfoQuery) GetHeader() *QueryHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *TokenGetNftInfoQuery) GetNftID() *NftID {
	if x != nil {
		return x.NftID
	}
	return nil
}

type TokenNftInfo struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	NftID        *NftID     `protobuf:"bytes,1,opt,name=nftID,proto3" json:"nftID,omitempty"`               // The ID of the NFT
	AccountID    *AccountID `protobuf:"bytes,2,opt,name=accountID,proto3" json:"accountID,omitempty"`       // The current owner of the NFT
	CreationTime *Timestamp `protobuf:"bytes,3,opt,name=creationTime,proto3" json:"creationTime,omitempty"` // The effective consensus timestamp at which the NFT was minted
	Metadata     []byte     `protobuf:"bytes,4,opt,name=metadata,proto3" json:"metadata,omitempty"`         // Represents the unique metadata of the NFT
}

func (x *TokenNftInfo) Reset() {
	*x = TokenNftInfo{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_token_get_nft_info_proto_msgTypes[2]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenNftInfo) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenNftInfo) ProtoMessage() {}

func (x *TokenNftInfo) ProtoReflect() protoreflect.Message {
	mi := &file_proto_token_get_nft_info_proto_msgTypes[2]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenNftInfo.ProtoReflect.Descriptor instead.
func (*TokenNftInfo) Descriptor() ([]byte, []int) {
	return file_proto_token_get_nft_info_proto_rawDescGZIP(), []int{2}
}

func (x *TokenNftInfo) GetNftID() *NftID {
	if x != nil {
		return x.NftID
	}
	return nil
}

func (x *TokenNftInfo) GetAccountID() *AccountID {
	if x != nil {
		return x.AccountID
	}
	return nil
}

func (x *TokenNftInfo) GetCreationTime() *Timestamp {
	if x != nil {
		return x.CreationTime
	}
	return nil
}

func (x *TokenNftInfo) GetMetadata() []byte {
	if x != nil {
		return x.Metadata
	}
	return nil
}

type TokenGetNftInfoResponse struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Header *ResponseHeader `protobuf:"bytes,1,opt,name=header,proto3" json:"header,omitempty"` // Standard response from node to client, including the requested fields: cost, or state proof, or both, or neither
	Nft    *TokenNftInfo   `protobuf:"bytes,2,opt,name=nft,proto3" json:"nft,omitempty"`       // The information about this NFT
}

func (x *TokenGetNftInfoResponse) Reset() {
	*x = TokenGetNftInfoResponse{}
	if protoimpl.UnsafeEnabled {
		mi := &file_proto_token_get_nft_info_proto_msgTypes[3]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *TokenGetNftInfoResponse) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*TokenGetNftInfoResponse) ProtoMessage() {}

func (x *TokenGetNftInfoResponse) ProtoReflect() protoreflect.Message {
	mi := &file_proto_token_get_nft_info_proto_msgTypes[3]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use TokenGetNftInfoResponse.ProtoReflect.Descriptor instead.
func (*TokenGetNftInfoResponse) Descriptor() ([]byte, []int) {
	return file_proto_token_get_nft_info_proto_rawDescGZIP(), []int{3}
}

func (x *TokenGetNftInfoResponse) GetHeader() *ResponseHeader {
	if x != nil {
		return x.Header
	}
	return nil
}

func (x *TokenGetNftInfoResponse) GetNft() *TokenNftInfo {
	if x != nil {
		return x.Nft
	}
	return nil
}

var File_proto_token_get_nft_info_proto protoreflect.FileDescriptor

var file_proto_token_get_nft_info_proto_rawDesc = []byte{
	0x0a, 0x1e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x5f, 0x67, 0x65,
	0x74, 0x5f, 0x6e, 0x66, 0x74, 0x5f, 0x69, 0x6e, 0x66, 0x6f, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x12, 0x05, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x17, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x62,
	0x61, 0x73, 0x69, 0x63, 0x5f, 0x74, 0x79, 0x70, 0x65, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x1a, 0x18, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x71, 0x75, 0x65, 0x72, 0x79, 0x5f, 0x68, 0x65,
	0x61, 0x64, 0x65, 0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x1b, 0x70, 0x72, 0x6f, 0x74,
	0x6f, 0x2f, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x5f, 0x68, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x15, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2f, 0x74,
	0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0x55,
	0x0a, 0x05, 0x4e, 0x66, 0x74, 0x49, 0x44, 0x12, 0x28, 0x0a, 0x07, 0x74, 0x6f, 0x6b, 0x65, 0x6e,
	0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0e, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x49, 0x44, 0x52, 0x07, 0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x49,
	0x44, 0x12, 0x22, 0x0a, 0x0c, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e, 0x75, 0x6d, 0x62, 0x65,
	0x72, 0x18, 0x02, 0x20, 0x01, 0x28, 0x03, 0x52, 0x0c, 0x73, 0x65, 0x72, 0x69, 0x61, 0x6c, 0x4e,
	0x75, 0x6d, 0x62, 0x65, 0x72, 0x22, 0x66, 0x0a, 0x14, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x47, 0x65,
	0x74, 0x4e, 0x66, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x51, 0x75, 0x65, 0x72, 0x79, 0x12, 0x2a, 0x0a,
	0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x12, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x51, 0x75, 0x65, 0x72, 0x79, 0x48, 0x65, 0x61, 0x64, 0x65,
	0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x22, 0x0a, 0x05, 0x6e, 0x66, 0x74,
	0x49, 0x44, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x4e, 0x66, 0x74, 0x49, 0x44, 0x52, 0x05, 0x6e, 0x66, 0x74, 0x49, 0x44, 0x22, 0xb4, 0x01,
	0x0a, 0x0c, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x4e, 0x66, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x12, 0x22,
	0x0a, 0x05, 0x6e, 0x66, 0x74, 0x49, 0x44, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x0c, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x4e, 0x66, 0x74, 0x49, 0x44, 0x52, 0x05, 0x6e, 0x66, 0x74,
	0x49, 0x44, 0x12, 0x2e, 0x0a, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x18,
	0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x41, 0x63,
	0x63, 0x6f, 0x75, 0x6e, 0x74, 0x49, 0x44, 0x52, 0x09, 0x61, 0x63, 0x63, 0x6f, 0x75, 0x6e, 0x74,
	0x49, 0x44, 0x12, 0x34, 0x0a, 0x0c, 0x63, 0x72, 0x65, 0x61, 0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69,
	0x6d, 0x65, 0x18, 0x03, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x10, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f,
	0x2e, 0x54, 0x69, 0x6d, 0x65, 0x73, 0x74, 0x61, 0x6d, 0x70, 0x52, 0x0c, 0x63, 0x72, 0x65, 0x61,
	0x74, 0x69, 0x6f, 0x6e, 0x54, 0x69, 0x6d, 0x65, 0x12, 0x1a, 0x0a, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x18, 0x04, 0x20, 0x01, 0x28, 0x0c, 0x52, 0x08, 0x6d, 0x65, 0x74, 0x61,
	0x64, 0x61, 0x74, 0x61, 0x22, 0x6f, 0x0a, 0x17, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x47, 0x65, 0x74,
	0x4e, 0x66, 0x74, 0x49, 0x6e, 0x66, 0x6f, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x12,
	0x2d, 0x0a, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x18, 0x01, 0x20, 0x01, 0x28, 0x0b, 0x32,
	0x15, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x2e, 0x52, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65,
	0x48, 0x65, 0x61, 0x64, 0x65, 0x72, 0x52, 0x06, 0x68, 0x65, 0x61, 0x64, 0x65, 0x72, 0x12, 0x25,
	0x0a, 0x03, 0x6e, 0x66, 0x74, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x2e, 0x54, 0x6f, 0x6b, 0x65, 0x6e, 0x4e, 0x66, 0x74, 0x49, 0x6e, 0x66, 0x6f,
	0x52, 0x03, 0x6e, 0x66, 0x74, 0x42, 0x4b, 0x0a, 0x1a, 0x63, 0x6f, 0x6d, 0x2e, 0x68, 0x65, 0x64,
	0x65, 0x72, 0x61, 0x2e, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2e, 0x70, 0x72,
	0x6f, 0x74, 0x6f, 0x50, 0x01, 0x5a, 0x2b, 0x67, 0x69, 0x74, 0x68, 0x75, 0x62, 0x2e, 0x63, 0x6f,
	0x6d, 0x2f, 0x68, 0x61, 0x73, 0x68, 0x67, 0x72, 0x61, 0x70, 0x68, 0x2f, 0x68, 0x65, 0x64, 0x65,
	0x72, 0x61, 0x2d, 0x73, 0x64, 0x6b, 0x2d, 0x67, 0x6f, 0x2f, 0x76, 0x32, 0x2f, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x62, 0x06, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x33,
}

var (
	file_proto_token_get_nft_info_proto_rawDescOnce sync.Once
	file_proto_token_get_nft_info_proto_rawDescData = file_proto_token_get_nft_info_proto_rawDesc
)

func file_proto_token_get_nft_info_proto_rawDescGZIP() []byte {
	file_proto_token_get_nft_info_proto_rawDescOnce.Do(func() {
		file_proto_token_get_nft_info_proto_rawDescData = protoimpl.X.CompressGZIP(file_proto_token_get_nft_info_proto_rawDescData)
	})
	return file_proto_token_get_nft_info_proto_rawDescData
}

var file_proto_token_get_nft_info_proto_msgTypes = make([]protoimpl.MessageInfo, 4)
var file_proto_token_get_nft_info_proto_goTypes = []interface{}{
	(*NftID)(nil),                   // 0: proto.NftID
	(*TokenGetNftInfoQuery)(nil),    // 1: proto.TokenGetNftInfoQuery
	(*TokenNftInfo)(nil),            // 2: proto.TokenNftInfo
	(*TokenGetNftInfoResponse)(nil), // 3: proto.TokenGetNftInfoResponse
	(*TokenID)(nil),                 // 4: proto.TokenID
	(*QueryHeader)(nil),             // 5: proto.QueryHeader
	(*AccountID)(nil),               // 6: proto.AccountID
	(*Timestamp)(nil),               // 7: proto.Timestamp
	(*ResponseHeader)(nil),          // 8: proto.ResponseHeader
}
var file_proto_token_get_nft_info_proto_depIdxs = []int32{
	4, // 0: proto.NftID.tokenID:type_name -> proto.TokenID
	5, // 1: proto.TokenGetNftInfoQuery.header:type_name -> proto.QueryHeader
	0, // 2: proto.TokenGetNftInfoQuery.nftID:type_name -> proto.NftID
	0, // 3: proto.TokenNftInfo.nftID:type_name -> proto.NftID
	6, // 4: proto.TokenNftInfo.accountID:type_name -> proto.AccountID
	7, // 5: proto.TokenNftInfo.creationTime:type_name -> proto.Timestamp
	8, // 6: proto.TokenGetNftInfoResponse.header:type_name -> proto.ResponseHeader
	2, // 7: proto.TokenGetNftInfoResponse.nft:type_name -> proto.TokenNftInfo
	8, // [8:8] is the sub-list for method output_type
	8, // [8:8] is the sub-list for method input_type
	8, // [8:8] is the sub-list for extension type_name
	8, // [8:8] is the sub-list for extension extendee
	0, // [0:8] is the sub-list for field type_name
}

func init() { file_proto_token_get_nft_info_proto_init() }
func file_proto_token_get_nft_info_proto_init() {
	if File_proto_token_get_nft_info_proto != nil {
		return
	}
	file_proto_basic_types_proto_init()
	file_proto_query_header_proto_init()
	file_proto_response_header_proto_init()
	file_proto_timestamp_proto_init()
	if !protoimpl.UnsafeEnabled {
		file_proto_token_get_nft_info_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*NftID); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_token_get_nft_info_proto_msgTypes[1].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenGetNftInfoQuery); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_token_get_nft_info_proto_msgTypes[2].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenNftInfo); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
		file_proto_token_get_nft_info_proto_msgTypes[3].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*TokenGetNftInfoResponse); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_proto_token_get_nft_info_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   4,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_proto_token_get_nft_info_proto_goTypes,
		DependencyIndexes: file_proto_token_get_nft_info_proto_depIdxs,
		MessageInfos:      file_proto_token_get_nft_info_proto_msgTypes,
	}.Build()
	File_proto_token_get_nft_info_proto = out.File
	file_proto_token_get_nft_info_proto_rawDesc = nil
	file_proto_token_get_nft_info_proto_goTypes = nil
	file_proto_token_get_nft_info_proto_depIdxs = nil
}
