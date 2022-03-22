package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// AccountStakersQuery gets all of the accounts that are proxy staking to this account. For each of  them, the amount
// currently staked will be given. This is not yet implemented, but will be in a future version of the API.
type AccountStakersQuery struct {
	Query
	accountID *AccountID
}

// NewAccountStakersQuery creates an AccountStakersQuery query which can be used to construct and execute
// an AccountStakersQuery.
//
// It is recommended that you use this for creating new instances of an AccountStakersQuery
// instead of manually creating an instance of the struct.
func NewAccountStakersQuery() *AccountStakersQuery {
	header := services.QueryHeader{}
	return &AccountStakersQuery{
		Query: _NewQuery(true, &header),
	}
}

func (query *AccountStakersQuery) SetGrpcDeadline(deadline *time.Duration) *AccountStakersQuery {
	query.Query.SetGrpcDeadline(deadline)
	return query
}

// SetAccountID sets the Account ID for which the stakers should be retrieved
func (query *AccountStakersQuery) SetAccountID(accountID AccountID) *AccountStakersQuery {
	query.accountID = &accountID
	return query
}

func (query *AccountStakersQuery) GetAccountID() AccountID {
	if query.accountID == nil {
		return AccountID{}
	}

	return *query.accountID
}

func (query *AccountStakersQuery) _ValidateNetworkOnIDs(client *Client) error {
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

func (query *AccountStakersQuery) _Build() *services.Query_CryptoGetProxyStakers {
	pb := services.Query_CryptoGetProxyStakers{
		CryptoGetProxyStakers: &services.CryptoGetStakersQuery{
			Header: &services.QueryHeader{},
		},
	}

	if query.accountID != nil {
		pb.CryptoGetProxyStakers.AccountID = query.accountID._ToProtobuf()
	}

	return &pb
}

func (query *AccountStakersQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}
	if query.Query.nodeAccountIDs.locked {
		for range query.nodeAccountIDs.slice {
			paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
			if err != nil {
				return Hbar{}, err
			}
			query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
		}
	}

	pb := query._Build()
	pb.CryptoGetProxyStakers.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_AccountStakersQueryShouldRetry,
		_CostQueryMakeRequest,
		_CostQueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_AccountStakersQueryGetMethod,
		_AccountStakersQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetCryptoGetProxyStakers().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _AccountStakersQueryShouldRetry(logID string, _ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(logID, Status(response.query.GetCryptoGetProxyStakers().Header.NodeTransactionPrecheckCode))
}

func _AccountStakersQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetCryptoGetProxyStakers().Header.NodeTransactionPrecheckCode),
	}
}

func _AccountStakersQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetStakersByAccountID,
	}
}

func (query *AccountStakersQuery) Execute(client *Client) ([]Transfer, error) {
	if client == nil || client.operator == nil {
		return []Transfer{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return []Transfer{}, err
	}

	if !query.paymentTransactionIDs.locked {
		query.paymentTransactionIDs._Clear()._Push(TransactionIDGenerate(client.operator.accountID))
	}

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
			return []Transfer{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return []Transfer{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "AccountStakersQuery",
			}
		}

		cost = actualCost
	}

	query.paymentTransactions = make([]*services.Transaction, 0)

	if query.nodeAccountIDs.locked {
		err = _QueryGeneratePayments(&query.Query, client, cost)
		if err != nil {
			return []Transfer{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return []Transfer{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.CryptoGetProxyStakers.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_AccountStakersQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_AccountStakersQueryGetMethod,
		_AccountStakersQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
	)

	if err != nil {
		return []Transfer{}, err
	}

	var stakers = make([]Transfer, len(resp.query.GetCryptoGetProxyStakers().Stakers.ProxyStaker))

	// TODO: This is wrong, this _Method shold return `[]ProxyStaker` not `[]Transfer`
	for i, element := range resp.query.GetCryptoGetProxyStakers().Stakers.ProxyStaker {
		id := _AccountIDFromProtobuf(element.AccountID)
		accountID := AccountID{}

		if id == nil {
			accountID = *id
		}

		stakers[i] = Transfer{
			AccountID: accountID,
			Amount:    HbarFromTinybar(element.Amount),
		}
	}

	return stakers, err
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *AccountStakersQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountStakersQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *AccountStakersQuery) SetQueryPayment(paymentAmount Hbar) *AccountStakersQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountStakersQuery.
func (query *AccountStakersQuery) SetNodeAccountIDs(accountID []AccountID) *AccountStakersQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *AccountStakersQuery) SetMaxRetry(count int) *AccountStakersQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *AccountStakersQuery) SetMaxBackoff(max time.Duration) *AccountStakersQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *AccountStakersQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *AccountStakersQuery) SetMinBackoff(min time.Duration) *AccountStakersQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *AccountStakersQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *AccountStakersQuery) _GetLogID() string {
	timestamp := query.timestamp.UnixNano()
	if query.paymentTransactionIDs._Length() > 0 && query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("AccountStakersQuery:%d", timestamp)
}

func (query *AccountStakersQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountStakersQuery {
	query.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return query
}
