package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type FileInfoQuery struct {
	Query
	fileID *FileID
}

func NewFileInfoQuery() *FileInfoQuery {
	return &FileInfoQuery{
		Query: newQuery(true),
	}
}

func (query *FileInfoQuery) SetSetFileID(fileID FileID) *FileInfoQuery {
	query.fileID = &fileID
	return query
}

func (query *FileInfoQuery) GetFileID(id FileID) FileID {
	return query.fileID
}

func (query *FileInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.fileID != nil {
		if err := query.fileID.Validate(client); err != nil {
			return err
		}
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

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
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
		_FileInfoQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_FileInfoQueryGetMethod,
		_FileInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetFileGetInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	}
	return HbarFromTinybar(cost), nil
}

func _FileInfoQueryShouldRetry(_ request, response response) executionState {
	return _QueryShouldRetry(Status(response.query.GetFileGetInfo().Header.NodeTransactionPrecheckCode))
}

func _FileInfoQueryMapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetFileGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func _FileInfoQueryGetMethod(_ request, channel *channel) method {
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

	err = _QueryGeneratePayments(&query.Query, client, cost)
	if err != nil {
		return FileInfo{}, err
	}
	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		_FileInfoQueryShouldRetry,
		query.queryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_FileInfoQueryGetMethod,
		_FileInfoQueryMapStatusError,
		_QueryMapResponse,
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

func (query *FileInfoQuery) GetNodeAccountIDs() []AccountID {
	return query.Query.GetNodeAccountIDs()
}

func (query *FileInfoQuery) SetMaxRetry(count int) *FileInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *FileInfoQuery) SetMaxBackoff(max time.Duration) *FileInfoQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *FileInfoQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *FileInfoQuery) SetMinBackoff(min time.Duration) *FileInfoQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *FileInfoQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}
