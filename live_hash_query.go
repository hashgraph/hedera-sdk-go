package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type LiveHashQuery struct {
	Query
	accountID AccountID
	hash      []byte
}

func NewLiveHashQuery() *LiveHashQuery {
	return &LiveHashQuery{
		Query: newQuery(true),
	}
}

func (query *LiveHashQuery) SetAccountID(id AccountID) *LiveHashQuery {
	query.accountID = id
	return query
}

func (query *LiveHashQuery) GetAccountID() AccountID {
	return query.accountID
}

func (query *LiveHashQuery) SetHash(hash []byte) *LiveHashQuery {
	query.hash = hash
	return query
}

func (query *LiveHashQuery) GetGetHash() []byte {
	return query.hash
}

func (query *LiveHashQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if err := query.accountID.Validate(client); err != nil {
		return err
	}

	return nil
}

func (query *LiveHashQuery) build() *proto.Query_CryptoGetLiveHash {
	body := &proto.CryptoGetLiveHashQuery{
		Header: &proto.QueryHeader{},
	}
	if !query.accountID.isZero() {
		body.AccountID = query.accountID.toProtobuf()
	}

	if len(query.hash) > 0 {
		body.Hash = query.hash
	}

	return &proto.Query_CryptoGetLiveHash{
		CryptoGetLiveHash: body,
	}
}

func (query *LiveHashQuery) queryMakeRequest() protoRequest {
	pb := query.build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.CryptoGetLiveHash.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.CryptoGetLiveHash.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *LiveHashQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	pb := query.build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	pb.CryptoGetLiveHash.Header.Payment = paymentTransaction
	pb.CryptoGetLiveHash.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *LiveHashQuery) GetCost(client *Client) (Hbar, error) {
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
		_LiveHashQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_LiveHashQueryGetMethod,
		_LiveHashQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetCryptoGetLiveHash().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _LiveHashQueryShouldRetry(_ request, response response) executionState {
	return _QueryShouldRetry(Status(response.query.GetCryptoGetLiveHash().Header.NodeTransactionPrecheckCode))
}

func _LiveHashQueryMapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetCryptoGetLiveHash().Header.NodeTransactionPrecheckCode),
	}
}

func _LiveHashQueryGetMethod(_ request, channel *channel) method {
	return method{
		query: channel.getCrypto().GetLiveHash,
	}
}

func (query *LiveHashQuery) Execute(client *Client) (LiveHash, error) {
	if client == nil || client.operator == nil {
		return LiveHash{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return LiveHash{}, err
	}

	query.build()

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
			return LiveHash{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return LiveHash{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "LiveHashQuery",
			}
		}

		cost = actualCost
	}

	err = _QueryGeneratePayments(&query.Query, client, cost)
	if err != nil {
		return LiveHash{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		_LiveHashQueryShouldRetry,
		query.queryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_LiveHashQueryGetMethod,
		_LiveHashQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return LiveHash{}, err
	}

	liveHash, err := liveHashFromProtobuf(resp.query.GetCryptoGetLiveHash().LiveHash)
	if err != nil {
		return LiveHash{}, err
	}

	return liveHash, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *LiveHashQuery) SetMaxQueryPayment(maxPayment Hbar) *LiveHashQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *LiveHashQuery) SetQueryPayment(paymentAmount Hbar) *LiveHashQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *LiveHashQuery) SetNodeAccountIDs(accountID []AccountID) *LiveHashQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *LiveHashQuery) SetMaxBackoff(max time.Duration) *LiveHashQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *LiveHashQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *LiveHashQuery) SetMinBackoff(min time.Duration) *LiveHashQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *LiveHashQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}
