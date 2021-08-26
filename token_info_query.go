package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

type TokenInfoQuery struct {
	Query
	tokenID TokenID
}

// NewTopicInfoQuery creates a TopicInfoQuery query which can be used to construct and execute a
//  Get Topic Info Query.
func NewTokenInfoQuery() *TokenInfoQuery {
	return &TokenInfoQuery{
		Query: newQuery(true),
	}
}

// SetTopicID sets the topic to retrieve info about (the parameters and running state of).
func (query *TokenInfoQuery) SetTokenID(id TokenID) *TokenInfoQuery {
	query.tokenID = id
	return query
}

func (query *TokenInfoQuery) GetTokenID() TokenID {
	return query.tokenID
}

func (query *TokenInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = query.tokenID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *TokenInfoQuery) build() *proto.Query_TokenGetInfo {
	body := &proto.TokenGetInfoQuery{
		Header: &proto.QueryHeader{},
	}
	if !query.tokenID.isZero() {
		body.Token = query.tokenID.toProtobuf()
	}

	return &proto.Query_TokenGetInfo{
		TokenGetInfo: body,
	}
}

func (query *TokenInfoQuery) queryMakeRequest() protoRequest {
	pb := query.build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.TokenGetInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.TokenGetInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *TokenInfoQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	pb := query.build()

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	pb.TokenGetInfo.Header.Payment = paymentTransaction
	pb.TokenGetInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *TokenInfoQuery) GetCost(client *Client) (Hbar, error) {
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
		tokenInfoQuery_shouldRetry,
		protoReq,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		tokenInfoQuery_getMethod,
		tokenInfoQuery_mapStatusError,
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

func tokenInfoQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetTokenGetInfo().Header.NodeTransactionPrecheckCode))
}

func tokenInfoQuery_mapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetTokenGetInfo().Header.NodeTransactionPrecheckCode),
	}
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

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return TokenInfo{}, err
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
			return TokenInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return TokenInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "TokenInfoQuery",
			}
		}

		cost = actualCost
	}

	err = query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return TokenInfo{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		tokenInfoQuery_shouldRetry,
		query.queryMakeRequest(),
		query_advanceRequest,
		query_getNodeAccountID,
		tokenInfoQuery_getMethod,
		tokenInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return TokenInfo{}, err
	}

	info := tokenInfoFromProtobuf(resp.query.GetTokenGetInfo().TokenInfo)

	return info, nil
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

func (query *TokenInfoQuery) SetMaxRetry(count int) *TokenInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *TokenInfoQuery) SetMaxBackoff(max time.Duration) *TokenInfoQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *TokenInfoQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *TokenInfoQuery) SetMinBackoff(min time.Duration) *TokenInfoQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *TokenInfoQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}
