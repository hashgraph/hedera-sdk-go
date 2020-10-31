package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type NetworkVersionInfoQuery struct {
	Query
	pb *proto.NetworkGetVersionInfoQuery
}

func NewNetworkVersionQuery() *NetworkVersionInfoQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.NetworkGetVersionInfoQuery{Header: &header}
	query.pb.Query = &proto.Query_NetworkGetVersionInfo{
		NetworkGetVersionInfo: &pb,
	}

	return &NetworkVersionInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

func networkVersionInfoQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetNetworkGetVersionInfo().Header.NodeTransactionPrecheckCode)
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
		return NetworkVersionInfo{}, err
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
		networkVersionInfoQuery_getMethod,
		networkVersionInfoQuery_mapResponseStatus,
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

func (query *NetworkVersionInfoQuery) GetNodeAccountIDs() []AccountID {
	return query.Query.GetNodeAccountIDs()
}
