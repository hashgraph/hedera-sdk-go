package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type FileInfo struct {
	FileID         FileID
	Size           int64
	ExpirationTime time.Time
	IsDeleted      bool
	Keys           KeyList
	FileMemo       string
}

func _FileInfoFromProtobuf(fileInfo *proto.FileGetInfoResponse_FileInfo) (FileInfo, error) {
	if fileInfo == nil {
		return FileInfo{}, errParameterNull
	}
	var keys KeyList
	var err error
	if fileInfo.Keys != nil {
		keys, err = _KeyListFromProtobuf(fileInfo.Keys)
		if err != nil {
			return FileInfo{}, err
		}
	}

	fileID := FileID{}
	if fileInfo.FileID != nil {
		fileID = *_FileIDFromProtobuf(fileInfo.FileID)
	}

	return FileInfo{
		FileID:         fileID,
		Size:           fileInfo.Size,
		ExpirationTime: _TimeFromProtobuf(fileInfo.ExpirationTime),
		IsDeleted:      fileInfo.Deleted,
		Keys:           keys,
		FileMemo:       fileInfo.Memo,
	}, nil
}

func (fileInfo *FileInfo) _ToProtobuf() *proto.FileGetInfoResponse_FileInfo {
	return &proto.FileGetInfoResponse_FileInfo{
		FileID: fileInfo.FileID._ToProtobuf(),
		Size:   fileInfo.Size,
		ExpirationTime: &proto.Timestamp{
			Seconds: int64(fileInfo.ExpirationTime.Second()),
			Nanos:   int32(fileInfo.ExpirationTime.Nanosecond()),
		},
		Deleted: fileInfo.IsDeleted,
		Keys:    fileInfo.Keys._ToProtoKeyList(),
		Memo:    fileInfo.FileMemo,
	}
}

func (fileInfo FileInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(fileInfo._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func FileInfoFromBytes(data []byte) (FileInfo, error) {
	if data == nil {
		return FileInfo{}, errByteArrayNull
	}
	pb := proto.FileGetInfoResponse_FileInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return FileInfo{}, err
	}

	info, err := _FileInfoFromProtobuf(&pb)
	if err != nil {
		return FileInfo{}, err
	}

	return info, nil
}
