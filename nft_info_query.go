package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenNftInfoQuery struct {
	Query
	pb *proto.TokenGetNftInfoQuery
}

func NewTokenNftInfoQuery() *TokenNftInfoQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.TokenGetNftInfoQuery{Header: &header}
	query.pb.Query = &proto.Query_TokenGetNftInfo{
		TokenGetNftInfo: &pb,
	}

	return &TokenNftInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

func (query *TokenNftInfoQuery) SetNftID(id NftID) *TokenNftInfoQuery {
	query.pb.NftID = id.toProtobuf()
	return query
}

func (query *TokenNftInfoQuery) GetNftID() NftID {
	return nftIDFromProtobuf(query.pb.NftID)
}

func (query *TokenNftInfoQuery) GetCost(client *Client) (Hbar, error) {
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
		tokenNftInfoQuery_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		tokenNftInfoQuery_getMethod,
		tokenNftInfoQuery_mapStatusError,
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

func tokenNftInfoQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetTokenGetNftInfo().Header.NodeTransactionPrecheckCode))
}

func tokenNftInfoQuery_mapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetTokenGetNftInfo().Header.NodeTransactionPrecheckCode),
	}
}

func tokenNftInfoQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getToken().GetTokenNftInfo,
	}
}

func (query *TokenNftInfoQuery) Execute(client *Client) (TokenNftInfo, error) {
	if client == nil || client.operator == nil {
		return TokenNftInfo{}, errNoClientProvided
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
			return TokenNftInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return TokenNftInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "TokenNftInfoQuery",
			}
		}

		cost = actualCost
	}

	err := query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return TokenNftInfo{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		tokenNftInfoQuery_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeAccountID,
		tokenNftInfoQuery_getMethod,
		tokenNftInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return TokenNftInfo{}, err
	}

	return tokenNftInfoFromProtobuf(resp.query.GetTokenGetNftInfo().GetNft()), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *TokenNftInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TokenNftInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *TokenNftInfoQuery) SetQueryPayment(paymentAmount Hbar) *TokenNftInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *TokenNftInfoQuery) SetNodeAccountIDs(accountID []AccountID) *TokenNftInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *TokenNftInfoQuery) SetMaxRetry(count int) *TokenNftInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}
