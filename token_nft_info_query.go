package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenNftInfoQuery struct {
	Query
	nftID *NftID
}

func NewTokenNftInfoQuery() *TokenNftInfoQuery {
	return &TokenNftInfoQuery{
		Query: _NewQuery(true),
		nftID: nil,
	}
}

func (query *TokenNftInfoQuery) SetNftID(nftID NftID) *TokenNftInfoQuery {
	query.nftID = &nftID
	return query
}

func (query *TokenNftInfoQuery) GetNftID() NftID {
	if query.nftID == nil {
		return NftID{}
	}

	return *query.nftID
}

// Deprecated
func (query *TokenNftInfoQuery) SetTokenID(id TokenID) *TokenNftInfoQuery {
	return query
}

// Deprecated
func (query *TokenNftInfoQuery) GetTokenID() TokenID {
	return TokenID{}
}

// Deprecated
func (query *TokenNftInfoQuery) SetAccountID(id AccountID) *TokenNftInfoQuery {
	return query
}

// Deprecated
func (query *TokenNftInfoQuery) GetAccountID() AccountID {
	return AccountID{}
}

// Deprecated
func (query *TokenNftInfoQuery) SetStart(start int64) *TokenNftInfoQuery {
	return query
}

// Deprecated
func (query *TokenNftInfoQuery) GetStart() int64 {
	return 0
}

// Deprecated
func (query *TokenNftInfoQuery) SetEnd(end int64) *TokenNftInfoQuery {
	return query
}

// Deprecated
func (query *TokenNftInfoQuery) GetEnd() int64 {
	return 0
}

// Deprecated
func (query *TokenNftInfoQuery) ByNftID(id NftID) *TokenNftInfoQuery {
	query.nftID = &id
	return query
}

// Deprecated
func (query *TokenNftInfoQuery) ByTokenID(id TokenID) *TokenNftInfoQuery {
	return query
}

// Deprecated
func (query *TokenNftInfoQuery) ByAccountID(id AccountID) *TokenNftInfoQuery {
	return query
}

func (query *TokenNftInfoQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.nftID != nil {
		if err := query.nftID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *TokenNftInfoQuery) _BuildByNft() *proto.Query_TokenGetNftInfo {
	body := &proto.TokenGetNftInfoQuery{
		Header: &proto.QueryHeader{},
	}

	if query.nftID != nil {
		body.NftID = query.nftID._ToProtobuf()
	}

	return &proto.Query_TokenGetNftInfo{
		TokenGetNftInfo: body,
	}
}

func (query *TokenNftInfoQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._BuildByNft()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.TokenGetNftInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.TokenGetNftInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *TokenNftInfoQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb := query._BuildByNft()
	pb.TokenGetNftInfo.Header.Payment = paymentTransaction
	pb.TokenGetNftInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *TokenNftInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error
	nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
	if err != nil {
		return Hbar{}, err
	}
	query.SetNodeAccountIDs(nodeAccountIDs)

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	protoReq, err := query._CostQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	var resp _IntermediateResponse
	resp, err = _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_TokenNftInfoQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_TokenNftInfoQueryGetMethod,
		_TokenNftInfoQueryMapStatusError,
		_QueryMapResponse,
	)
	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetTokenGetNftInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	}
	return HbarFromTinybar(cost), nil
}

func _TokenNftInfoQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetTokenGetNftInfo().Header.NodeTransactionPrecheckCode))
}

func _TokenNftInfoQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetTokenGetNftInfo().Header.NodeTransactionPrecheckCode),
	}
}

func _TokenNftInfoQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetToken().GetTokenNftInfo,
	}
}

func (query *TokenNftInfoQuery) Execute(client *Client) ([]TokenNftInfo, error) {
	if client == nil || client.operator == nil {
		return []TokenNftInfo{}, errNoClientProvided
	}

	var err error

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
		if err != nil {
			return []TokenNftInfo{}, err
		}

		query.SetNodeAccountIDs(nodeAccountIDs)
	}
	err = query._ValidateNetworkOnIDs(client)
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

	err = _QueryGeneratePayments(&query.Query, client, cost)
	if err != nil {
		return []TokenNftInfo{}, err
	}

	var resp _IntermediateResponse
	tokenInfos := make([]TokenNftInfo, 0)
	resp, err = _Execute(
		client,
		_Request{

			query: &query.Query,
		},
		_TokenNftInfoQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TokenNftInfoQueryGetMethod,
		_TokenNftInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return []TokenNftInfo{}, err
	}

	tokenInfos = append(tokenInfos, _TokenNftInfoFromProtobuf(resp.query.GetTokenGetNftInfo().GetNft()))
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
