package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type FileInfoQuery struct {
	Query
	pb *proto.FileGetInfoQuery
}

func NewFileInfoQuery() *FileInfoQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.FileGetInfoQuery{Header: &header}
	query.pb.Query = &proto.Query_FileGetInfo{
		FileGetInfo: &pb,
	}

	return &FileInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

func (query *FileInfoQuery) SetFileID(id FileID) *FileInfoQuery {
	query.pb.FileID = id.toProtobuf()
	return query
}

func (query *FileInfoQuery) GetFileID(id FileID) FileID {
	return fileIDFromProtobuf(query.pb.GetFileID())
}

func (query *FileInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return Hbar{}, err
	}

	query.pbHeader.Payment = paymentTransaction
	query.pbHeader.ResponseType = proto.ResponseType_COST_ANSWER
	query.nodeIDs = client.getNodeAccountIDsForExecute()

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		query_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		accountInfoQuery_getMethod,
		accountInfoQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetFileGetInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	} else {
		return HbarFromTinybar(cost), nil
	}
}

func fileInfoQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetFileGetInfo().Header.NodeTransactionPrecheckCode)
}

func fileInfoQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getFile().GetFileInfo,
	}
}

func (query *FileInfoQuery) Execute(client *Client) (FileInfo, error) {
	if client == nil || client.operator == nil {
		return FileInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.getNodeAccountIDsForExecute())
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
		return FileInfo{}, err
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
		fileInfoQuery_getMethod,
		fileInfoQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return FileInfo{}, err
	}

	return fileInfoFromProtobuf(resp.query.GetFileGetInfo().FileInfo)
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *FileInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *FileInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *FileInfoQuery) SetQueryPayment(paymentAmount Hbar) *FileInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *FileInfoQuery) SetNodeAccountID(accountID []AccountID) *FileInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *FileInfoQuery) GetNodeAccountId() []AccountID {
	return query.Query.GetNodeAccountIDs()
}
