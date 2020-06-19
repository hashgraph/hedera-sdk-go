package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type FileInfoQuery struct {
	QueryBuilder
	pb *proto.FileGetInfoQuery
}

type FileInfo struct {
	FileID         FileID
	Size           int64
	ExpirationTime time.Time
	IsDeleted      bool
	Keys           []PublicKey
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

	pbKeys := resp.GetFileGetInfo().FileInfo.Keys
	keys, err := publicKeyListFromProto(pbKeys)

	if err != nil {
		return FileInfo{}, err
	}

	return FileInfo{
		FileID:         fileIDFromProto(resp.GetFileGetInfo().FileInfo.FileID),
		Size:           resp.GetFileGetInfo().FileInfo.Size,
		ExpirationTime: timeFromProto(resp.GetFileGetInfo().FileInfo.ExpirationTime),
		IsDeleted:      resp.GetFileGetInfo().FileInfo.Deleted,
		Keys:           keys,
	}, nil
}

func (builder *FileInfoQuery) Cost(client *Client) (Hbar, error) {
	// deleted files return a COST_ANSWER of zero which triggers `INSUFFICIENT_TX_FEE`
	// if you set that as the query payment; 25 tinybar seems to be enough to get
	// `FILE_DELETED` back instead.
	cost, err := builder.QueryBuilder.GetCost(client)
	if err != nil {
		return ZeroHbar, err
	}

	// math.Min requires float64 and returns float64
	if cost.AsTinybar() > 25 {
		return cost, nil
	}

	return HbarFromTinybar(25), nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (builder *FileInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *FileInfoQuery {
	return &FileInfoQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (builder *FileInfoQuery) SetQueryPayment(paymentAmount Hbar) *FileInfoQuery {
	return &FileInfoQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (builder *FileInfoQuery) SetQueryPaymentTransaction(tx Transaction) *FileInfoQuery {
	return &FileInfoQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
