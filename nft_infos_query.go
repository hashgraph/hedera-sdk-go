package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenNftInfos struct {
	TokenID TokenID
	Ntfs    []TokenNftInfo
}

type TokenNftInfosQuery struct {
	Query
	pb *proto.TokenGetNftInfosQuery
}

func NewTokenNftInfosQuery() *TokenNftInfosQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.TokenGetNftInfosQuery{Header: &header}
	query.pb.Query = &proto.Query_TokenGetNftInfos{
		TokenGetNftInfos: &pb,
	}

	return &TokenNftInfosQuery{
		Query: query,
		pb:    &pb,
	}
}

func (query *TokenNftInfosQuery) SetTokenID(id TokenID) *TokenNftInfosQuery {
	query.pb.TokenID = id.toProtobuf()
	return query
}

func (query *TokenNftInfosQuery) GetTokenID() TokenID {
	return tokenIDFromProtobuf(query.pb.TokenID)
}

func (query *TokenNftInfosQuery) SetStart(start int64) *TokenNftInfosQuery {
	query.pb.Start = start
	return query
}

func (query *TokenNftInfosQuery) GetStart() int64 {
	return query.pb.Start
}

func (query *TokenNftInfosQuery) SetEnd(end int64) *TokenNftInfosQuery {
	query.pb.End = end
	return query
}

func (query *TokenNftInfosQuery) GetEnd() int64 {
	return query.pb.End
}

func (query *TokenNftInfosQuery) GetCost(client *Client) (Hbar, error) {
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
		tokenNftInfosQuery_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		tokenNftInfosQuery_getMethod,
		tokenNftInfosQuery_mapStatusError,
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

func tokenNftInfosQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetTokenGetNftInfos().Header.NodeTransactionPrecheckCode))
}

func tokenNftInfosQuery_mapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetTokenGetNftInfos().Header.NodeTransactionPrecheckCode),
	}
}

func tokenNftInfosQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getToken().GetTokenNftInfos,
	}
}

func (query *TokenNftInfosQuery) Execute(client *Client) (TokenNftInfos, error) {
	if client == nil || client.operator == nil {
		return TokenNftInfos{}, errNoClientProvided
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
			return TokenNftInfos{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return TokenNftInfos{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "TokenNftInfosQuery",
			}
		}

		cost = actualCost
	}

	err := query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return TokenNftInfos{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		tokenNftInfosQuery_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeAccountID,
		tokenNftInfosQuery_getMethod,
		tokenNftInfosQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return TokenNftInfos{}, err
	}

	nfts := resp.query.GetTokenGetNftInfos().Nfts

	nftsConverted := make([]TokenNftInfo, len(nfts))

	for i, nft := range nfts {
		nftsConverted[i] = tokenNftInfoFromProtobuf(nft)
	}

	return TokenNftInfos{
		TokenID: tokenIDFromProtobuf(resp.query.GetTokenGetNftInfos().TokenID),
		Ntfs:    nftsConverted,
	}, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *TokenNftInfosQuery) SetMaxQueryPayment(maxPayment Hbar) *TokenNftInfosQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *TokenNftInfosQuery) SetQueryPayment(paymentAmount Hbar) *TokenNftInfosQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *TokenNftInfosQuery) SetNodeAccountIDs(accountID []AccountID) *TokenNftInfosQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *TokenNftInfosQuery) SetMaxRetry(count int) *TokenNftInfosQuery {
	query.Query.SetMaxRetry(count)
	return query
}
