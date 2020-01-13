package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type FileContentsQuery struct {
	QueryBuilder
	pb *proto.FileGetContentsQuery
}

func NewFileContentsQuery() *FileContentsQuery {
	pb := &proto.FileGetContentsQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_FileGetContents{FileGetContents: pb}

	return &FileContentsQuery{inner, pb}
}

func (builder *FileContentsQuery) SetFileID(id FileID) *FileContentsQuery {
	builder.pb.FileID = id.toProto()
	return builder
}

func (builder *FileContentsQuery) Execute(client *Client) ([]byte, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return []byte{}, err
	}

	return resp.GetFileGetContents().FileContents.Contents, nil
}
