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
	Keys 		   []PublicKey
}

func NewFileInfoQuery() *FileInfoQuery {
	pb := &proto.FileGetInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_FileGetInfo{FileGetInfo: pb}

	return &FileInfoQuery{inner, pb}
}

func (builder *FileInfoQuery) SetFileID(id FileID) *FileInfoQuery {
	builder.pb.FileID = id.toProto()
	return builder
}

func (builder *FileInfoQuery) Execute(client *Client) (FileInfo, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return FileInfo{}, err
	}

	pbKeyList := resp.GetFileGetInfo().FileInfo.Keys.Keys

	keyList := make([]PublicKey, len(pbKeyList))

	for i, key := range pbKeyList {

		// todo: support more than ed25519keys
		keyList[i], err = Ed25519PublicKeyFromBytes(key.GetEd25519())

		if err != nil {
			return FileInfo{}, err
		}
	}

	return FileInfo{
		FileID:         fileIDFromProto(resp.GetFileGetInfo().FileInfo.FileID),
		Size:           resp.GetFileGetInfo().FileInfo.Size,
		ExpirationTime: timeFromProto(resp.GetFileGetInfo().FileInfo.ExpirationTime),
		Deleted:        resp.GetFileGetInfo().FileInfo.Deleted,
		Keys:			keyList,
	}, nil
}

func (builder *FileInfoQuery) Cost(client *Client) (uint64, error) {
	// deleted files return a COST_ANSWER of zero which triggers `INSUFFICIENT_TX_FEE`
	// if you set that as the query payment; 25 tinybar seems to be enough to get
	// `FILE_DELETED` back instead.
	cost, err := builder.QueryBuilder.Cost(client)
	if err != nil {
		return 0, err
	}

	// math.Min requires float64 and returns float64
	if cost > 25 {
		return cost, nil
	}

	return 25, nil
}
