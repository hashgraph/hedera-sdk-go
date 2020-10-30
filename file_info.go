package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type FileInfo struct {
	FileID         FileID
	Size           int64
	ExpirationTime time.Time
	IsDeleted      bool
	Keys           []Key
}

func newFileInfo(fileID FileID, size int64, expirationTime time.Time, isDeleted bool, keys []Key) FileInfo {
	return FileInfo{
		FileID:         fileID,
		Size:           size,
		ExpirationTime: expirationTime,
		IsDeleted:      isDeleted,
		Keys:           keys,
	}
}

func fileInfoFromProtobuf(fileInfo *proto.FileGetInfoResponse_FileInfo) (FileInfo, error) {
	var keys []Key
	if fileInfo.Keys != nil {
		keyList, err := keyListFromProtobuf(fileInfo.Keys)
		if err != nil {
			return FileInfo{}, err
		}

		keys = keyList.keys
	}

	return FileInfo{
		FileID:         fileIDFromProtobuf(fileInfo.FileID),
		Size:           fileInfo.Size,
		ExpirationTime: timeFromProtobuf(fileInfo.ExpirationTime),
		IsDeleted:      fileInfo.Deleted,
		Keys:           keys,
	}, nil
}

func (fileInfo *FileInfo) toProtobuf() *proto.FileGetInfoResponse_FileInfo {
	var keys = make([]*proto.Key, 0)
	for _, key := range fileInfo.Keys {
		keys = append(keys, key.toProtoKey())
	}

	return &proto.FileGetInfoResponse_FileInfo{
		FileID: fileInfo.FileID.toProtobuf(),
		Size:   fileInfo.Size,
		ExpirationTime: &proto.Timestamp{
			Seconds: int64(fileInfo.ExpirationTime.Second()),
			Nanos:   int32(fileInfo.ExpirationTime.Nanosecond()),
		},
		Deleted: fileInfo.IsDeleted,
		Keys: &proto.KeyList{
			Keys: keys,
		},
	}
}
