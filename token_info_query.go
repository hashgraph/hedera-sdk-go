package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenInfoQuery struct {
	Query
	pb *proto.TokenGetInfoQuery
}

// NewTopicInfoQuery creates a TopicInfoQuery query which can be used to construct and execute a
//  Get Topic Info Query.
func NewTokenInfoQuery() *TokenInfoQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.TokenGetInfoQuery{Header: &header}
	query.pb.Query = &proto.Query_TokenGetInfo{
		TokenGetInfo: &pb,
	}

	return &TokenInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

// SetTopicID sets the topic to retrieve info about (the parameters and running state of).
func (query *TokenInfoQuery) SetTokenID(id TokenID) *TokenInfoQuery {
	query.pb.Token = id.toProtobuf()
	return query
}

func (query *TokenInfoQuery) GetTokenID() TokenID {
	return tokenIDFromProtobuf(query.pb.Token)
}

func (query *TokenInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return Hbar{}, err
	}

	query.pbHeader.Payment = paymentTransaction
	query.pbHeader.ResponseType = proto.ResponseType_COST_ANSWER
	query.nodeIDs = client.network.getNodeAccountIDsForExecute()

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

	cost := int64(resp.query.GetTokenGetInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	} else {
		return HbarFromTinybar(cost), nil
	}
}

func tokenInfoQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetTokenGetInfo().Header.NodeTransactionPrecheckCode)
}

func tokenInfoQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getToken().GetTokenInfo,
	}
}

// Execute executes the TopicInfoQuery using the provided client
func (query *TokenInfoQuery) Execute(client *Client) (TokenInfo, error) {
	if client == nil || client.operator == nil {
		return TokenInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	query.queryPayment = NewHbar(2)
	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		cost = client.maxQueryPayment

		actualCost, err := query.GetCost(client)
		if err != nil {
			return TokenInfo{}, err
		}

		if cost.tinybar > actualCost.tinybar {
			return TokenInfo{}, ErrMaxQueryPaymentExceeded{}
		}

		cost = actualCost
	}

	err := query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return TokenInfo{}, err
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
		tokenInfoQuery_getMethod,
		tokenInfoQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return TokenInfo{}, err
	}

	return tokenInfoFromProtobuf(resp.query.GetTokenGetInfo().TokenInfo), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *TokenInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TokenInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *TokenInfoQuery) SetQueryPayment(paymentAmount Hbar) *TokenInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *TokenInfoQuery) SetNodeAccountIDs(accountID []AccountID) *TokenInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *TokenInfoQuery) GetNodeAccountIDs() []AccountID {
	return query.Query.GetNodeAccountIDs()
}
