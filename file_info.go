package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

type FileInfo struct {
	FileID         FileID
	Size           int64
	ExpirationTime time.Time
	IsDeleted      bool
	Keys           KeyList
	FileMemo       string
}

func newFileInfo(fileID FileID, size int64, expirationTime time.Time, isDeleted bool, keys KeyList, fileMemo string) FileInfo {
	return FileInfo{
		FileID:         fileID,
		Size:           size,
		ExpirationTime: expirationTime,
		IsDeleted:      isDeleted,
		Keys:           keys,
		FileMemo:       fileMemo,
	}
}

func fileInfoFromProtobuf(fileInfo *proto.FileGetInfoResponse_FileInfo) (FileInfo, error) {
	if fileInfo == nil {
		return FileInfo{}, errParameterNull
	}
	var keys KeyList
	var err error
	if fileInfo.Keys != nil {
		keys, err = keyListFromProtobuf(fileInfo.Keys)
		if err != nil {
			return FileInfo{}, err
		}
	}

	return FileInfo{
		FileID:         fileIDFromProtobuf(fileInfo.FileID),
		Size:           fileInfo.Size,
		ExpirationTime: timeFromProtobuf(fileInfo.ExpirationTime),
		IsDeleted:      fileInfo.Deleted,
		Keys:           keys,
		FileMemo:       fileInfo.Memo,
	}, nil
}

func (fileInfo *FileInfo) toProtobuf() *proto.FileGetInfoResponse_FileInfo {
	return &proto.FileGetInfoResponse_FileInfo{
		FileID: fileInfo.FileID.toProtobuf(),
		Size:   fileInfo.Size,
		ExpirationTime: &proto.Timestamp{
			Seconds: int64(fileInfo.ExpirationTime.Second()),
			Nanos:   int32(fileInfo.ExpirationTime.Nanosecond()),
		},
		Deleted: fileInfo.IsDeleted,
		Keys:    fileInfo.Keys.toProtoKeyList(),
		Memo:    fileInfo.FileMemo,
	}
}

func (fileInfo FileInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(fileInfo.toProtobuf())
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

	info, err := fileInfoFromProtobuf(&pb)
	if err != nil {
		return FileInfo{}, err
	}

	return info, nil
}
