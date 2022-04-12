package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type AccountDetailsQuery struct {
	Query
	accountID *AccountID
}

func NewAccountDetailsQuery() *AccountDetailsQuery {
	header := services.QueryHeader{}
	return &AccountDetailsQuery{
		Query: _NewQuery(true, &header),
	}
}

func (query *AccountDetailsQuery) SetAccountID(id AccountID) *AccountDetailsQuery {
	query.accountID = &id
	return query
}

func (query *AccountDetailsQuery) GetAccountID() AccountID {
	if query.accountID != nil {
		return *query.accountID
	}

	return AccountID{}
}

func (query *AccountDetailsQuery) _Build() *services.Query_AccountDetails {
	pb := services.GetAccountDetailsQuery{Header: &services.QueryHeader{}}

	if query.accountID != nil {
		pb.AccountId = query.accountID._ToProtobuf()
	}

	return &services.Query_AccountDetails{
		AccountDetails: &pb,
	}
}

func (query *AccountDetailsQuery) _ValidateNetworkOnIDs(client *Client) error {
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

func (query *AccountDetailsQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	pb := services.Query_AccountDetails{
		AccountDetails: &services.GetAccountDetailsQuery{},
	}
	pb.AccountDetails.Header = query.pbHeader

	query.pb = &services.Query{
		Query: &pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_AccountDetailsQueryShouldRetry,
		_CostQueryMakeRequest,
		_CostQueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_AccountDetailsQueryGetMethod,
		_AccountDetailsQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.(*services.Response).GetAccountDetails().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	}
	return HbarFromTinybar(cost), nil
}

func _AccountDetailsQueryShouldRetry(logID string, _ interface{}, response interface{}) _ExecutionState {
	return _QueryShouldRetry(logID, Status(response.(*services.Response).GetAccountDetails().Header.NodeTransactionPrecheckCode))
}

func _AccountDetailsQueryMapStatusError(_ interface{}, response interface{}) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.(*services.Response).GetAccountDetails().Header.NodeTransactionPrecheckCode),
	}
}

func _AccountDetailsQueryGetMethod(_ interface{}, channel *_Channel) _Method {
	return _Method{
		query: channel._GetNetwork().GetAccountDetails,
	}
}

func (query *AccountDetailsQuery) Execute(client *Client) (AccountDetails, error) {
	if client == nil || client.operator == nil {
		return AccountDetails{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return AccountDetails{}, err
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
			return AccountDetails{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return AccountDetails{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "NetworkVersionInfo",
			}
		}

		cost = actualCost
	}

	query.paymentTransactions = make([]*services.Transaction, 0)

	if query.nodeAccountIDs.locked {
		err = _QueryGeneratePayments(&query.Query, client, cost)
		if err != nil {
			return AccountDetails{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return AccountDetails{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.AccountDetails.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		&query.Query,
		_AccountDetailsQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_AccountDetailsQueryGetMethod,
		_AccountDetailsQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err != nil {
		return AccountDetails{}, err
	}

	return _AccountDetailsFromProtobuf(resp.(*services.Response).GetAccountDetails().AccountDetails)
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *AccountDetailsQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountDetailsQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *AccountDetailsQuery) SetQueryPayment(paymentAmount Hbar) *AccountDetailsQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *AccountDetailsQuery) SetNodeAccountIDs(accountID []AccountID) *AccountDetailsQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *AccountDetailsQuery) SetMaxRetry(count int) *AccountDetailsQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *AccountDetailsQuery) SetMaxBackoff(max time.Duration) *AccountDetailsQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *AccountDetailsQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *AccountDetailsQuery) SetMinBackoff(min time.Duration) *AccountDetailsQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *AccountDetailsQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *AccountDetailsQuery) _GetLogID() string {
	timestamp := query.timestamp.UnixNano()
	if query.paymentTransactionIDs._Length() > 0 && query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("AccountDetailsQuery:%d", timestamp)
}

func (query *AccountDetailsQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountDetailsQuery {
	query.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return query
}
