package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

// FileInfo contains the details about a file stored on Hiero
type FileInfo struct {
	FileID         FileID
	Size           int64
	ExpirationTime time.Time
	IsDeleted      bool
	Keys           KeyList
	FileMemo       string
	LedgerID       LedgerID
}

func _FileInfoFromProtobuf(fileInfo *services.FileGetInfoResponse_FileInfo) (FileInfo, error) {
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
		LedgerID:       LedgerID{fileInfo.LedgerId},
	}, nil
}

func (fileInfo *FileInfo) _ToProtobuf() *services.FileGetInfoResponse_FileInfo {
	return &services.FileGetInfoResponse_FileInfo{
		FileID: fileInfo.FileID._ToProtobuf(),
		Size:   fileInfo.Size,
		ExpirationTime: &services.Timestamp{
			Seconds: int64(fileInfo.ExpirationTime.Second()),
			Nanos:   int32(fileInfo.ExpirationTime.Nanosecond()),
		},
		Deleted:  fileInfo.IsDeleted,
		Keys:     fileInfo.Keys._ToProtoKeyList(),
		Memo:     fileInfo.FileMemo,
		LedgerId: fileInfo.LedgerID.ToBytes(),
	}
}

// ToBytes returns the byte representation of the FileInfo
func (fileInfo FileInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(fileInfo._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// FileInfoFromBytes returns a FileInfo object from a raw byte array
func FileInfoFromBytes(data []byte) (FileInfo, error) {
	if data == nil {
		return FileInfo{}, errByteArrayNull
	}
	pb := services.FileGetInfoResponse_FileInfo{}
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
