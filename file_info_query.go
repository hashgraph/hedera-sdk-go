package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type FileInfoQuery struct {
	Query
	fileID FileID
}

func NewFileInfoQuery() *FileInfoQuery {

	return &FileInfoQuery{
		Query: newQuery(true),
	}
}

func (query *FileInfoQuery) SetFileID(id FileID) *FileInfoQuery {
	query.fileID = id
	return query
}

func (query *FileInfoQuery) GetFileID(id FileID) FileID {
	return query.fileID
}

func (query *FileInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = query.fileID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *FileInfoQuery) build() *proto.Query_FileGetInfo {
	body := &proto.FileGetInfoQuery{
		Header: &proto.QueryHeader{},
	}
	if !query.fileID.isZero() {
		body.FileID = query.fileID.toProtobuf()
	}

	return &proto.Query_FileGetInfo{
		FileGetInfo: body,
	}
}

func (query *FileInfoQuery) queryMakeRequest() protoRequest {
	pb := query.build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.FileGetInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.FileGetInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *FileInfoQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	pb := query.build()

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	pb.FileGetInfo.Header.Payment = paymentTransaction
	pb.FileGetInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *FileInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	protoReq, err := query.costQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		fileInfoQuery_shouldRetry,
		protoReq,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		fileInfoQuery_getMethod,
		fileInfoQuery_mapStatusError,
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

func fileInfoQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetFileGetInfo().Header.NodeTransactionPrecheckCode))
}

func fileInfoQuery_mapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetFileGetInfo().Header.NodeTransactionPrecheckCode),
	}
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
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return FileInfo{}, err
	}

	query.build()

	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		if query.maxQueryPayment.tinybar == 0 {
			cost = client.maxQueryPayment
		} else {
			cost = query.maxQueryPayment
		}

		actualCost, err := query.GetCost(client)
		if err != nil {
			return FileInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return FileInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "FileInfoQuery",
			}
		}

		cost = actualCost
	}

	err = query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return FileInfo{}, err
	}
	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		fileInfoQuery_shouldRetry,
		query.queryMakeRequest(),
		query_advanceRequest,
		query_getNodeAccountID,
		fileInfoQuery_getMethod,
		fileInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return FileInfo{}, err
	}

	info, err := fileInfoFromProtobuf(resp.query.GetFileGetInfo().FileInfo)
	if err != nil {
		return FileInfo{}, err
	}

	return info, nil
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

func (query *FileInfoQuery) SetNodeAccountIDs(accountID []AccountID) *FileInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *FileInfoQuery) GetNodeAccountId() []AccountID {
	return query.Query.GetNodeAccountIDs()
}

func (query *FileInfoQuery) SetMaxRetry(count int) *FileInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}
