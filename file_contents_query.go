package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

// FileContentsQuery retrieves the contents of a file.
type FileContentsQuery struct {
	QueryBuilder
	pb *proto.FileGetContentsQuery
}

// NewFileContentsQuery creates a FileContentsQuery transaction which can be used to construct and execute a
// File Get Contents Query.
func NewFileContentsQuery() *FileContentsQuery {
	pb := &proto.FileGetContentsQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_FileGetContents{FileGetContents: pb}

	return &FileContentsQuery{inner, pb}
}

// SetFileID sets the FileID of the file whose contents are requested.
func (transaction *FileContentsQuery) SetFileID(id FileID) *FileContentsQuery {
	transaction.pb.FileID = id.toProto()
	return transaction
}

// Execute executes the FileContentsQuery using the provided client. The returned byte slice will be empty if the file
// is empty.
func (transaction *FileContentsQuery) Execute(client *Client) ([]byte, error) {
	resp, err := transaction.execute(client)
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
func (transaction *FileContentsQuery) SetMaxQueryPayment(maxPayment Hbar) *FileContentsQuery {
	return &FileContentsQuery{*transaction.QueryBuilder.SetMaxQueryPayment(maxPayment), transaction.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (transaction *FileContentsQuery) SetQueryPayment(paymentAmount Hbar) *FileContentsQuery {
	return &FileContentsQuery{*transaction.QueryBuilder.SetQueryPayment(paymentAmount), transaction.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (transaction *FileContentsQuery) SetQueryPaymentTransaction(tx Transaction) *FileContentsQuery {
	return &FileContentsQuery{*transaction.QueryBuilder.SetQueryPaymentTransaction(tx), transaction.pb}
}
