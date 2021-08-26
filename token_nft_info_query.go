package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"github.com/pkg/errors"
	"time"
)

type TokenNftInfoQuery struct {
	Query
	byNft     bool
	byToken   bool
	byAccount bool
	tokenID   TokenID
	nftID     NftID
	accountID AccountID
	start     int64
	end       int64
}

func NewTokenNftInfoQuery() *TokenNftInfoQuery {
	return &TokenNftInfoQuery{
		Query:     newQuery(true),
		byNft:     false,
		byToken:   false,
		byAccount: false,
		tokenID:   TokenID{},
		nftID:     NftID{},
		accountID: AccountID{},
	}
}

func (query *TokenNftInfoQuery) SetNftID(id NftID) *TokenNftInfoQuery {
	query.nftID = id
	return query
}

func (query *TokenNftInfoQuery) GetNftID() NftID {
	return query.nftID
}

func (query *TokenNftInfoQuery) SetTokenID(id TokenID) *TokenNftInfoQuery {
	query.tokenID = id
	return query
}

func (query *TokenNftInfoQuery) GetTokenID() TokenID {
	return query.tokenID
}

func (query *TokenNftInfoQuery) SetAccountID(id AccountID) *TokenNftInfoQuery {
	query.accountID = id
	return query
}

func (query *TokenNftInfoQuery) GetAccountID() AccountID {
	return query.accountID
}

func (query *TokenNftInfoQuery) SetStart(start int64) *TokenNftInfoQuery {
	query.start = start
	return query
}

func (query *TokenNftInfoQuery) GetStart() int64 {
	return query.start
}

func (query *TokenNftInfoQuery) SetEnd(end int64) *TokenNftInfoQuery {
	query.end = end
	return query
}

func (query *TokenNftInfoQuery) GetEnd() int64 {
	return query.end
}

func (query *TokenNftInfoQuery) isByNft() bool {
	return query.byNft
}

func (query *TokenNftInfoQuery) isByToken() bool {
	return query.byToken
}

func (query *TokenNftInfoQuery) isByAccount() bool {
	return query.byAccount
}

func (query *TokenNftInfoQuery) ByNftID(id NftID) *TokenNftInfoQuery {
	query.byNft = true
	query.nftID = id

	return query
}

func (query *TokenNftInfoQuery) ByTokenID(id TokenID) *TokenNftInfoQuery {
	query.byToken = true
	query.tokenID = id

	return query
}

func (query *TokenNftInfoQuery) ByAccountID(id AccountID) *TokenNftInfoQuery {
	query.byAccount = true
	query.accountID = id

	return query
}

func (query *TokenNftInfoQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	if query.isByToken() {
		err = query.tokenID.Validate(client)
		if err != nil {
			return err
		}
	}
	if query.isByAccount() {
		err = query.accountID.Validate(client)
		if err != nil {
			return err
		}
	}
	if query.isByNft() {
		err = query.nftID.Validate(client)
		if err != nil {
			return err
		}
	}

	return nil
}

func (query *TokenNftInfoQuery) buildByToken() *proto.Query_TokenGetNftInfos {
	if query.isByToken() {
		body := &proto.TokenGetNftInfosQuery{
			Header: &proto.QueryHeader{},
		}
		if !query.tokenID.isZero() {
			body.TokenID = query.tokenID.toProtobuf()
		}
		body.Start = query.start
		body.End = query.end

		return &proto.Query_TokenGetNftInfos{
			TokenGetNftInfos: body,
		}
	}

	return &proto.Query_TokenGetNftInfos{}
}

func (query *TokenNftInfoQuery) buildByNft() *proto.Query_TokenGetNftInfo {
	if query.isByNft() {
		body := &proto.TokenGetNftInfoQuery{
			Header: &proto.QueryHeader{},
		}
		body.NftID = query.nftID.toProtobuf()

		return &proto.Query_TokenGetNftInfo{
			TokenGetNftInfo: body,
		}
	}

	return &proto.Query_TokenGetNftInfo{}
}

func (query *TokenNftInfoQuery) buildByAccount() *proto.Query_TokenGetAccountNftInfos {
	if query.isByAccount() {
		body := &proto.TokenGetAccountNftInfosQuery{Header: &proto.QueryHeader{}}
		if !query.accountID.isZero() {
			body.AccountID = query.accountID.toProtobuf()
		}
		body.Start = query.start
		body.End = query.end

		return &proto.Query_TokenGetAccountNftInfos{
			TokenGetAccountNftInfos: body,
		}
	}

	return &proto.Query_TokenGetAccountNftInfos{}
}

func (query *TokenNftInfoQuery) queryMakeRequest() protoRequest {
	if query.isByToken() {
		pb := query.buildByToken()
		if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
			pb.TokenGetNftInfos.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
		}
		pb.TokenGetNftInfos.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

		return protoRequest{
			query: &proto.Query{
				Query: pb,
			},
		}
	}
	if query.isByNft() {
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
	if query.isByAccount() {
		pb := query.buildByAccount()
		if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
			pb.TokenGetAccountNftInfos.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
		}
		pb.TokenGetAccountNftInfos.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

		return protoRequest{
			query: &proto.Query{
				Query: pb,
			},
		}
	}

	return protoRequest{}
}

func (query *TokenNftInfoQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	if query.isByToken() {
		pb := query.buildByToken()
		pb.TokenGetNftInfos.Header.Payment = paymentTransaction
		pb.TokenGetNftInfos.Header.ResponseType = proto.ResponseType_COST_ANSWER

		return protoRequest{
			query: &proto.Query{
				Query: pb,
			},
		}, nil
	}

	if query.isByNft() {
		pb := query.buildByNft()
		pb.TokenGetNftInfo.Header.Payment = paymentTransaction
		pb.TokenGetNftInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

		return protoRequest{
			query: &proto.Query{
				Query: pb,
			},
		}, nil
	}

	if query.isByToken() {
		pb := query.buildByAccount()
		pb.TokenGetAccountNftInfos.Header.Payment = paymentTransaction
		pb.TokenGetAccountNftInfos.Header.ResponseType = proto.ResponseType_COST_ANSWER

		return protoRequest{
			query: &proto.Query{
				Query: pb,
			},
		}, nil
	}

	return protoRequest{}, nil
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

	enabled := 0
	if query.isByNft() {
		enabled += 1
	}
	if query.isByAccount() {
		enabled += 1
	}
	if query.isByToken() {
		enabled += 1
	}

	if enabled > 1 {
		return Hbar{}, errors.New("TokenNftInfoQuery must be one of ByNftId, ByTokenId, or ByAccountId, but multiple of these modes have been selected")
	} else if enabled == 0 {
		return Hbar{}, errors.New("TokenNftInfoQuery must be one of ByNftId, ByTokenId, or ByAccountId, but none of these modes have been selected")
	}

	protoReq, err := query.costQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	var resp intermediateResponse
	if query.isByNft() {
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
	} else if query.isByToken() {
		resp, err = execute(
			client,
			request{
				query: &query.Query,
			},
			tokenNftInfosQuery_shouldRetry,
			protoReq,
			costQuery_advanceRequest,
			costQuery_getNodeAccountID,
			tokenNftInfosQuery_getMethod,
			tokenNftInfosQuery_mapStatusError,
			query_mapResponse,
		)

		if err != nil {
			return Hbar{}, err
		}

		cost := int64(resp.query.GetTokenGetNftInfos().Header.Cost)
		if cost < 25 {
			return HbarFromTinybar(25), nil
		} else {
			return HbarFromTinybar(cost), nil
		}
	} else {
		resp, err = execute(
			client,
			request{
				query: &query.Query,
			},
			accountNftInfoQuery_shouldRetry,
			protoReq,
			costQuery_advanceRequest,
			costQuery_getNodeAccountID,
			accountNftInfoQuery_getMethod,
			accountNftInfoQuery_mapStatusError,
			query_mapResponse,
		)

		if err != nil {
			return Hbar{}, err
		}

		cost := int64(resp.query.GetTokenGetAccountNftInfos().Header.Cost)
		if cost < 25 {
			return HbarFromTinybar(25), nil
		} else {
			return HbarFromTinybar(cost), nil
		}
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

	enabled := 0
	if query.isByNft() {
		enabled += 1
	}
	if query.isByAccount() {
		enabled += 1
	}
	if query.isByToken() {
		enabled += 1
	}

	if enabled > 1 {
		return []TokenNftInfo{}, errors.New("TokenNftInfoQuery must be one of ByNftId, ByTokenId, or ByAccountId, but multiple of these modes have been selected")
	} else if enabled == 0 {
		return []TokenNftInfo{}, errors.New("TokenNftInfoQuery must be one of ByNftId, ByTokenId, or ByAccountId, but none of these modes have been selected")
	}

	var resp intermediateResponse
	tokenInfos := make([]TokenNftInfo, 0)
	if query.isByNft() {
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

	} else if query.isByToken() {
		resp, err = execute(
			client,
			request{
				query: &query.Query,
			},
			tokenNftInfosQuery_shouldRetry,
			query.queryMakeRequest(),
			query_advanceRequest,
			query_getNodeAccountID,
			tokenNftInfosQuery_getMethod,
			tokenNftInfosQuery_mapStatusError,
			query_mapResponse,
		)

		if err != nil {
			return []TokenNftInfo{}, err
		}

		nfts := resp.query.GetTokenGetNftInfos().GetNfts()
		for _, tokenInfo := range nfts {
			tokenInfos = append(tokenInfos, tokenNftInfoFromProtobuf(tokenInfo))
		}

	} else {
		resp, err = execute(
			client,
			request{
				query: &query.Query,
			},
			accountNftInfoQuery_shouldRetry,
			query.queryMakeRequest(),
			query_advanceRequest,
			query_getNodeAccountID,
			accountNftInfoQuery_getMethod,
			accountNftInfoQuery_mapStatusError,
			query_mapResponse,
		)

		if err != nil {
			return []TokenNftInfo{}, err
		}

		nfts := resp.query.GetTokenGetAccountNftInfos().GetNfts()
		for _, tokenInfo := range nfts {
			tokenInfos = append(tokenInfos, tokenNftInfoFromProtobuf(tokenInfo))
		}
	}

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
