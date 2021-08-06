package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type NetworkVersionInfoQuery struct {
	Query
}

func NewNetworkVersionQuery() *NetworkVersionInfoQuery {
	return &NetworkVersionInfoQuery{
		Query: newQuery(true),
	}
}

func (query *NetworkVersionInfoQuery) queryMakeRequest() protoRequest {
	pb := &proto.Query_NetworkGetVersionInfo{
		NetworkGetVersionInfo: &proto.NetworkGetVersionInfoQuery{
			Header: &proto.QueryHeader{},
		},
	}
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.NetworkGetVersionInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.NetworkGetVersionInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *NetworkVersionInfoQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	pb := &proto.Query_NetworkGetVersionInfo{
		NetworkGetVersionInfo: &proto.NetworkGetVersionInfoQuery{
			Header: &proto.QueryHeader{},
		},
	}

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	pb.NetworkGetVersionInfo.Header.Payment = paymentTransaction
	pb.NetworkGetVersionInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *NetworkVersionInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

	protoReq, err := query.costQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		networkVersionInfoQuery_shouldRetry,
		protoReq,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		networkVersionInfoQuery_getMethod,
		networkVersionInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetNetworkGetVersionInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	} else {
		return HbarFromTinybar(cost), nil
	}
}

func networkVersionInfoQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetNetworkGetVersionInfo().Header.NodeTransactionPrecheckCode))
}

func networkVersionInfoQuery_mapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetNetworkGetVersionInfo().Header.NodeTransactionPrecheckCode),
	}
}

func networkVersionInfoQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getNetwork().GetVersionInfo,
	}
}

func (query *NetworkVersionInfoQuery) Execute(client *Client) (NetworkVersionInfo, error) {
	if client == nil || client.operator == nil {
		return NetworkVersionInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
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

	err := query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return NetworkVersionInfo{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		networkVersionInfoQuery_shouldRetry,
		query.queryMakeRequest(),
		query_advanceRequest,
		query_getNodeAccountID,
		networkVersionInfoQuery_getMethod,
		networkVersionInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return NetworkVersionInfo{}, err
	}

	return networkVersionInfoFromProtobuf(resp.query.GetNetworkGetVersionInfo()), err
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
