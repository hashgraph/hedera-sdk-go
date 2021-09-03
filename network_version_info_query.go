package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type NetworkVersionInfoQuery struct {
	Query
}

func NewNetworkVersionQuery() *NetworkVersionInfoQuery {
	return &NetworkVersionInfoQuery{
		Query: _NewQuery(true),
	}
}

func (query *NetworkVersionInfoQuery) _QueryMakeRequest() _ProtoRequest {
	pb := &proto.Query_NetworkGetVersionInfo{
		NetworkGetVersionInfo: &proto.NetworkGetVersionInfoQuery{
			Header: &proto.QueryHeader{},
		},
	}
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.NetworkGetVersionInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.NetworkGetVersionInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *NetworkVersionInfoQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := &proto.Query_NetworkGetVersionInfo{
		NetworkGetVersionInfo: &proto.NetworkGetVersionInfoQuery{
			Header: &proto.QueryHeader{},
		},
	}

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.NetworkGetVersionInfo.Header.Payment = paymentTransaction
	pb.NetworkGetVersionInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *NetworkVersionInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	query.nodeIDs = client.network._GetNodeAccountIDsForExecute()

	protoReq, err := query._CostQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_NetworkVersionInfoQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_NetworkVersionInfoQueryGetMethod,
		_NetworkVersionInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetNetworkGetVersionInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	}
	return HbarFromTinybar(cost), nil
}

func _NetworkVersionInfoQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetNetworkGetVersionInfo().Header.NodeTransactionPrecheckCode))
}

func _NetworkVersionInfoQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetNetworkGetVersionInfo().Header.NodeTransactionPrecheckCode),
	}
}

func _NetworkVersionInfoQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetNetwork().GetVersionInfo,
	}
}

func (query *NetworkVersionInfoQuery) Execute(client *Client) (NetworkVersionInfo, error) {
	if client == nil || client.operator == nil {
		return NetworkVersionInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

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
			return NetworkVersionInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return NetworkVersionInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "NetworkVersionInfoQuery",
			}
		}

		cost = actualCost
	}

	err := _QueryGeneratePayments(&query.Query, client, cost)
	if err != nil {
		return NetworkVersionInfo{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_NetworkVersionInfoQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_NetworkVersionInfoQueryGetMethod,
		_NetworkVersionInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return NetworkVersionInfo{}, err
	}

	return _NetworkVersionInfoFromProtobuf(resp.query.GetNetworkGetVersionInfo()), err
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *NetworkVersionInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *NetworkVersionInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *NetworkVersionInfoQuery) SetQueryPayment(paymentAmount Hbar) *NetworkVersionInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *NetworkVersionInfoQuery) SetNodeAccountIDs(accountID []AccountID) *NetworkVersionInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *NetworkVersionInfoQuery) SetMaxRetry(count int) *NetworkVersionInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *NetworkVersionInfoQuery) SetMaxBackoff(max time.Duration) *NetworkVersionInfoQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *NetworkVersionInfoQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *NetworkVersionInfoQuery) SetMinBackoff(min time.Duration) *NetworkVersionInfoQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *NetworkVersionInfoQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}
