package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type FileInfoQuery struct {
	QueryBuilder
	pb *proto.FileGetInfoQuery
}

type FileInfo struct {
	FileID         FileID
	Size           int64
	ExpirationTime time.Time
	Deleted        bool
	// TODO: When KeyList is implemented
	// Keys []Keys
}

func NewFileInfoQuery() *FileInfoQuery {
	pb := &proto.FileGetInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_FileGetInfo{pb}

	return &FileInfoQuery{inner, pb}
}

func (builder *FileInfoQuery) SetFileID(id FileID) *FileInfoQuery {
	builder.pb.FileID = id.toProto()
	return builder
}

func (builder *FileInfoQuery) Execute(client *Client) (*FileInfo, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return nil, err
	}

	return &FileInfo{
		FileID:         fileIDFromProto(resp.GetFileGetInfo().FileInfo.FileID),
		Size:           resp.GetFileGetInfo().FileInfo.Size,
		ExpirationTime: timeFromProto(resp.GetFileGetInfo().FileInfo.ExpirationTime),
		Deleted:        resp.GetFileGetInfo().FileInfo.Deleted,
	}, nil
}
