package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type LiveHashQuery struct {
	Query
	accountID *AccountID
	hash      []byte
}

func NewLiveHashQuery() *LiveHashQuery {
	header := services.QueryHeader{}
	return &LiveHashQuery{
		Query: _NewQuery(true, &header),
	}
}

func (query *LiveHashQuery) SetAccountID(accountID AccountID) *LiveHashQuery {
	query.accountID = &accountID
	return query
}

func (query *LiveHashQuery) GetAccountID() AccountID {
	if query.accountID == nil {
		return AccountID{}
	}

	return *query.accountID
}

func (query *LiveHashQuery) SetHash(hash []byte) *LiveHashQuery {
	query.hash = hash
	return query
}

func (query *LiveHashQuery) GetGetHash() []byte {
	return query.hash
}

func (query *LiveHashQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.accountID != nil {
		if err := query.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *LiveHashQuery) _Build() *services.Query_CryptoGetLiveHash {
	body := &services.CryptoGetLiveHashQuery{
		Header: &services.QueryHeader{},
	}
	if query.accountID != nil {
		body.AccountID = query.accountID._ToProtobuf()
	}

	if len(query.hash) > 0 {
		body.Hash = query.hash
	}

	return &services.Query_CryptoGetLiveHash{
		CryptoGetLiveHash: body,
	}
}

func (query *LiveHashQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error
	if len(query.Query.GetNodeAccountIDs()) == 0 {
		nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
		if err != nil {
			return Hbar{}, err
		}

		query.SetNodeAccountIDs(nodeAccountIDs)
	}

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	for range query.nodeAccountIDs {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.CryptoGetLiveHash.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_LiveHashQueryShouldRetry,
		_CostQueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
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

func _LiveHashQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetCryptoGetLiveHash().Header.NodeTransactionPrecheckCode))
}

func _LiveHashQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetCryptoGetLiveHash().Header.NodeTransactionPrecheckCode),
	}
}

func _LiveHashQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetLiveHash,
	}
}

func (query *LiveHashQuery) Execute(client *Client) (LiveHash, error) {
	if client == nil || client.operator == nil {
		return LiveHash{}, errNoClientProvided
	}

	var err error
	if len(query.Query.GetNodeAccountIDs()) == 0 {
		nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
		if err != nil {
			return LiveHash{}, err
		}

		query.SetNodeAccountIDs(nodeAccountIDs)
	}
	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return LiveHash{}, err
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

	query.nextPaymentTransactionIndex = 0
	query.paymentTransactions = make([]*services.Transaction, 0)

	err = _QueryGeneratePayments(&query.Query, client, cost)
	if err != nil {
		return LiveHash{}, err
	}

	pb := query._Build()
	pb.CryptoGetLiveHash.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_LiveHashQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_LiveHashQueryGetMethod,
		_LiveHashQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return LiveHash{}, err
	}

	liveHash, err := _LiveHashFromProtobuf(resp.query.GetCryptoGetLiveHash().LiveHash)
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
