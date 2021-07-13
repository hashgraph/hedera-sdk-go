package hedera

import "github.com/hashgraph/hedera-protobufs-go/services"

type NetworkVersionInfoQuery struct {
	Query
	pb *services.NetworkGetVersionInfoQuery
}

func NewNetworkVersionQuery() *NetworkVersionInfoQuery {
	header := services.QueryHeader{}
	query := newQuery(true, &header)
	pb := services.NetworkGetVersionInfoQuery{Header: &header}
	query.pb.Query = &services.Query_NetworkGetVersionInfo{
		NetworkGetVersionInfo: &pb,
	}

	return &NetworkVersionInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

func (query *NetworkVersionInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return Hbar{}, err
	}

	query.pbHeader.Payment = paymentTransaction
	query.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		networkVersionInfoQuery_shouldRetry,
		costQuery_makeRequest,
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

func networkVersionInfoQuery_mapStatusError(_ request, response response, _ *NetworkName) error {
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
		query_makeRequest,
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
