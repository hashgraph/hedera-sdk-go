package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

// FileContentsQuery retrieves the contents of a file.
type FileContentsQuery struct {
	Query
	pb *proto.FileGetContentsQuery
}

// NewFileContentsQuery creates a FileContentsQuery query which can be used to construct and execute a
// File Get Contents Query.
func NewFileContentsQuery() *FileContentsQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.FileGetContentsQuery{Header: &header}
	query.pb.Query = &proto.Query_FileGetContents{
		FileGetContents: &pb,
	}

	return &FileContentsQuery{
		Query: query,
		pb:    &pb,
	}
}

// SetFileID sets the FileID of the file whose contents are requested.
func (query *FileContentsQuery) SetFileID(id FileID) *FileContentsQuery {
	query.pb.FileID = id.toProtobuf()
	return query
}

func (query *FileContentsQuery) GetFileID(id FileID) FileID {
	return fileIDFromProtobuf(query.pb.GetFileID())
}

func fileContentsQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetFileGetContents().Header.NodeTransactionPrecheckCode)
}

func fileContentsQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getFile().GetFileContent,
	}
}

func (query *FileContentsQuery) Execute(client *Client) ([]byte, error) {
	if client == nil || client.operator == nil {
		return []byte{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.getNodeAccountIDsForTransaction())
	}

	query.queryPayment = NewHbar(2)
	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		cost = client.maxQueryPayment

		// actualCost := CostQuery()
	}

	err := query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return []byte{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		query_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeAccountID,
		fileContentsQuery_getMethod,
		fileContentsQuery_mapResponseStatus,
		query_mapResponse,
	)

	return resp.query.GetFileGetContents().FileContents.Contents, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *FileContentsQuery) SetMaxQueryPayment(maxPayment Hbar) *FileContentsQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *FileContentsQuery) SetQueryPayment(paymentAmount Hbar) *FileContentsQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *FileContentsQuery) SetNodeAccountIDs(accountID []AccountID) *FileContentsQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *FileContentsQuery) GetNodeAccountIDs() []AccountID {
	return query.Query.GetNodeAccountIDs()
}
