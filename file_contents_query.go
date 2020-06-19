package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

// FileContentsQuery retrieves the contents of a file.
type FileContentsQuery struct {
	QueryBuilder
	pb *proto.FileGetContentsQuery
}

// NewFileContentsQuery creates a FileContentsQuery builder which can be used to construct and execute a
// File Get Contents Query.
func NewFileContentsQuery() *FileContentsQuery {
	pb := &proto.FileGetContentsQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_FileGetContents{FileGetContents: pb}

	return &FileContentsQuery{inner, pb}
}

// SetFileID sets the FileID of the file whose contents are requested.
func (builder *FileContentsQuery) SetFileID(id FileID) *FileContentsQuery {
	builder.pb.FileID = id.toProto()
	return builder
}

// Execute executes the FileContentsQuery using the provided client. The returned byte slice will be empty if the file
// is empty.
func (builder *FileContentsQuery) Execute(client *Client) ([]byte, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return []byte{}, err
	}

	return resp.GetFileGetContents().FileContents.Contents, nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (builder *FileContentsQuery) SetMaxQueryPayment(maxPayment Hbar) *FileContentsQuery {
	return &FileContentsQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (builder *FileContentsQuery) SetQueryPayment(paymentAmount Hbar) *FileContentsQuery {
	return &FileContentsQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (builder *FileContentsQuery) SetQueryPaymentTransaction(tx Transaction) *FileContentsQuery {
	return &FileContentsQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
