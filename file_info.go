package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

// FileInfo contains the details about a file stored on Hedera
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
