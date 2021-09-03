package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TokenInfoQuery struct {
	Query
	tokenID *TokenID
}

// NewTopicInfoQuery creates a TopicInfoQuery query which can be used to construct and execute a
//  Get Topic Info Query.
func NewTokenInfoQuery() *TokenInfoQuery {
	return &TokenInfoQuery{
		Query: _NewQuery(true),
	}
}

// SetTopicID sets the topic to retrieve info about (the parameters and running state of).
func (query *TokenInfoQuery) SetTokenID(tokenID TokenID) *TokenInfoQuery {
	query.tokenID = &tokenID
	return query
}

func (query *TokenInfoQuery) GetTokenID() TokenID {
	if query.tokenID == nil {
		return TokenID{}
	}

	return *query.tokenID
}

func (query *TokenInfoQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.tokenID != nil {
		if err := query.tokenID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *TokenInfoQuery) _Build() *proto.Query_TokenGetInfo {
	body := &proto.TokenGetInfoQuery{
		Header: &proto.QueryHeader{},
	}
	if !query.tokenID._IsZero() {
		body.Token = query.tokenID._ToProtobuf()
	}

	return &proto.Query_TokenGetInfo{
		TokenGetInfo: body,
	}
}

func (query *TokenInfoQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.TokenGetInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.TokenGetInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *TokenInfoQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.TokenGetInfo.Header.Payment = paymentTransaction
	pb.TokenGetInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *TokenInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	query.nodeIDs = client.network._GetNodeAccountIDsForExecute()

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	protoReq, err := query._CostQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_TokenInfoQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_TokenInfoQueryGetMethod,
		_TokenInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetTokenGetInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	}
	return HbarFromTinybar(cost), nil
}

func _TokenInfoQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetTokenGetInfo().Header.NodeTransactionPrecheckCode))
}

func _TokenInfoQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetTokenGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func _TokenInfoQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetToken().GetTokenInfo,
	}
}

// Execute executes the TopicInfoQuery using the provided client
func (query *TokenInfoQuery) Execute(client *Client) (TokenInfo, error) {
	if client == nil || client.operator == nil {
		return TokenInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

	err := query._ValidateNetworkOnIDs(client)
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

	err = _QueryGeneratePayments(&query.Query, client, cost)
	if err != nil {
		return TokenInfo{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_TokenInfoQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TokenInfoQueryGetMethod,
		_TokenInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return TokenInfo{}, err
	}

	info := _TokenInfoFromProtobuf(resp.query.GetTokenGetInfo().TokenInfo)

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
