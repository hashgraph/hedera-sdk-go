package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenAccountNftInfoQuery struct {
	Query
	pb *proto.TokenGetAccountNftInfoQuery
}

func NewTokenAccountNftInfoQuery() *TokenAccountNftInfoQuery {
	header := proto.QueryHeader{}
	query := newQuery(true, &header)
	pb := proto.TokenGetAccountNftInfoQuery{Header: &header}
	query.pb.Query = &proto.Query_TokenGetAccountNftInfo{
		TokenGetAccountNftInfo: &pb,
	}

	return &TokenAccountNftInfoQuery{
		Query: query,
		pb:    &pb,
	}
}

func (query *TokenAccountNftInfoQuery) SetAccountID(id AccountID) *TokenAccountNftInfoQuery {
	query.pb.AccountID = id.toProtobuf()
	return query
}

func (query *TokenAccountNftInfoQuery) GetAccountID() AccountID {
	return accountIDFromProtobuf(query.pb.AccountID)
}

func (query *TokenAccountNftInfoQuery) SetStart(start int64) *TokenAccountNftInfoQuery {
	query.pb.Start = start
	return query
}

func (query *TokenAccountNftInfoQuery) GetStart() int64 {
	return query.pb.Start
}

func (query *TokenAccountNftInfoQuery) SetEnd(end int64) *TokenAccountNftInfoQuery {
	query.pb.End = end
	return query
}

func (query *TokenAccountNftInfoQuery) GetEnd() int64 {
	return query.pb.End
}

func (query *TokenAccountNftInfoQuery) GetCost(client *Client) (Hbar, error) {
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
		tokenAccountNftInfoQuery_shouldRetry,
		costQuery_makeRequest,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		tokenAccountNftInfoQuery_getMethod,
		tokenAccountNftInfoQuery_mapStatusError,
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

func tokenAccountNftInfoQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetTokenGetAccountNftInfo().Header.NodeTransactionPrecheckCode))
}

func tokenAccountNftInfoQuery_mapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetTokenGetAccountNftInfo().Header.NodeTransactionPrecheckCode),
	}
}

func tokenAccountNftInfoQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getToken().GetAccountNftInfo,
	}
}

func (query *TokenAccountNftInfoQuery) Execute(client *Client) ([]TokenNftInfo, error) {
	if client == nil || client.operator == nil {
		return []TokenNftInfo{}, errNoClientProvided
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
			return []TokenNftInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return []TokenNftInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "TokenAccountNftInfoQuery",
			}
		}

		cost = actualCost
	}

	err := query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return []TokenNftInfo{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		tokenAccountNftInfoQuery_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeAccountID,
		tokenAccountNftInfoQuery_getMethod,
		tokenAccountNftInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return []TokenNftInfo{}, err
	}

	nfts := resp.query.GetTokenGetNftInfos().Nfts
	nftsConverted := make([]TokenNftInfo, len(nfts))

	for i, nft := range nfts {
		nftsConverted[i] = tokenNftInfoFromProtobuf(nft)
	}

	return nftsConverted, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *TokenAccountNftInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *TokenAccountNftInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *TokenAccountNftInfoQuery) SetQueryPayment(paymentAmount Hbar) *TokenAccountNftInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *TokenAccountNftInfoQuery) SetNodeAccountIDs(accountID []AccountID) *TokenAccountNftInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *TokenAccountNftInfoQuery) SetMaxRetry(count int) *TokenAccountNftInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}
