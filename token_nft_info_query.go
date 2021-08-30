package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

type TokenNftInfoQuery struct {
	Query
	nftID NftID
}

func NewTokenNftInfoQuery() *TokenNftInfoQuery {
	return &TokenNftInfoQuery{
		Query: newQuery(true),
		nftID: NftID{},
	}
}

func (query *TokenNftInfoQuery) SetNftID(id NftID) *TokenNftInfoQuery {
	query.nftID = id
	return query
}

func (query *TokenNftInfoQuery) GetNftID() NftID {
	return query.nftID
}

//Deprecated
func (query *TokenNftInfoQuery) SetTokenID(id TokenID) *TokenNftInfoQuery {
	return query
}

//Deprecated
func (query *TokenNftInfoQuery) GetTokenID() TokenID {
	return TokenID{}
}

//Deprecated
func (query *TokenNftInfoQuery) SetAccountID(id AccountID) *TokenNftInfoQuery {
	return query
}

//Deprecated
func (query *TokenNftInfoQuery) GetAccountID() AccountID {
	return AccountID{}
}

//Deprecated
func (query *TokenNftInfoQuery) SetStart(start int64) *TokenNftInfoQuery {
	return query
}

//Deprecated
func (query *TokenNftInfoQuery) GetStart() int64 {
	return 0
}

//Deprecated
func (query *TokenNftInfoQuery) SetEnd(end int64) *TokenNftInfoQuery {
	return query
}

//Deprecated
func (query *TokenNftInfoQuery) GetEnd() int64 {
	return 0
}

//Deprecated
func (query *TokenNftInfoQuery) ByNftID(id NftID) *TokenNftInfoQuery {
	query.nftID = id

	return query
}

//Deprecated
func (query *TokenNftInfoQuery) ByTokenID(id TokenID) *TokenNftInfoQuery {
	return query
}

//Deprecated
func (query *TokenNftInfoQuery) ByAccountID(id AccountID) *TokenNftInfoQuery {
	return query
}

func (query *TokenNftInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = query.nftID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *TokenNftInfoQuery) buildByNft() *proto.Query_TokenGetNftInfo {
	body := &proto.TokenGetNftInfoQuery{
		Header: &proto.QueryHeader{},
	}
	body.NftID = query.nftID.toProtobuf()

	return &proto.Query_TokenGetNftInfo{
		TokenGetNftInfo: body,
	}
}

func (query *TokenNftInfoQuery) queryMakeRequest() protoRequest {

	pb := query.buildByNft()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.TokenGetNftInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.TokenGetNftInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *TokenNftInfoQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	pb := query.buildByNft()
	pb.TokenGetNftInfo.Header.Payment = paymentTransaction
	pb.TokenGetNftInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *TokenNftInfoQuery) GetCost(client *Client) (Hbar, error) {
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

	var resp intermediateResponse
	resp, err = execute(
		client,
		request{
			query: &query.Query,
		},
		tokenNftInfoQuery_shouldRetry,
		protoReq,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		tokenNftInfoQuery_getMethod,
		tokenNftInfoQuery_mapStatusError,
		query_mapResponse,
	)
	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetTokenGetNftInfo().Header.Cost)
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

func accountNftInfoQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetTokenGetAccountNftInfos().Header.NodeTransactionPrecheckCode))
}

func accountNftInfoQuery_mapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetTokenGetAccountNftInfos().Header.NodeTransactionPrecheckCode),
	}
}

func accountNftInfoQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getToken().GetAccountNftInfos,
	}
}

func (query *TokenNftInfoQuery) Execute(client *Client) ([]TokenNftInfo, error) {
	if client == nil || client.operator == nil {
		return []TokenNftInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return []TokenNftInfo{}, err
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
				query:           "TokenNftInfoQuery",
			}
		}

		cost = actualCost
	}

	err = query_generatePayments(&query.Query, client, cost)
	if err != nil {
		return []TokenNftInfo{}, err
	}

	var resp intermediateResponse
	tokenInfos := make([]TokenNftInfo, 0)
	resp, err = execute(
		client,
		request{

			query: &query.Query,
		},
		tokenNftInfoQuery_shouldRetry,
		query.queryMakeRequest(),
		query_advanceRequest,
		query_getNodeAccountID,
		tokenNftInfoQuery_getMethod,
		tokenNftInfoQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return []TokenNftInfo{}, err
	}

	tokenInfos = append(tokenInfos, tokenNftInfoFromProtobuf(resp.query.GetTokenGetNftInfo().GetNft()))
	return tokenInfos, nil
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

func (query *TokenNftInfoQuery) SetMaxBackoff(max time.Duration) *TokenNftInfoQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *TokenNftInfoQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *TokenNftInfoQuery) SetMinBackoff(min time.Duration) *TokenNftInfoQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *TokenNftInfoQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}
