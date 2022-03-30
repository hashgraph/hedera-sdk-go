package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type AccountInfoQuery struct {
	Query
	accountID *AccountID
}

func NewAccountInfoQuery() *AccountInfoQuery {
	header := services.QueryHeader{}
	return &AccountInfoQuery{
		Query: _NewQuery(true, &header),
	}
}

func (query *AccountInfoQuery) SetGrpcDeadline(deadline *time.Duration) *AccountInfoQuery {
	query.Query.SetGrpcDeadline(deadline)
	return query
}

// SetAccountID sets the AccountID for this AccountInfoQuery.
func (query *AccountInfoQuery) SetAccountID(accountID AccountID) *AccountInfoQuery {
	query.accountID = &accountID
	return query
}

func (query *AccountInfoQuery) GetAccountID() AccountID {
	if query.accountID == nil {
		return AccountID{}
	}

	return *query.accountID
}

func (query *AccountInfoQuery) _ValidateNetworkOnIDs(client *Client) error {
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

func (query *AccountInfoQuery) _Build() *services.Query_CryptoGetInfo {
	pb := services.Query_CryptoGetInfo{
		CryptoGetInfo: &services.CryptoGetInfoQuery{
			Header: &services.QueryHeader{},
		},
	}

	if query.accountID != nil {
		pb.CryptoGetInfo.AccountID = query.accountID._ToProtobuf()
	}

	return &pb
}

func _AccountInfoQueryShouldRetry(logID string, _ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(logID, Status(response.query.GetCryptoGetInfo().Header.NodeTransactionPrecheckCode))
}

func _AccountInfoQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetCryptoGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func _AccountInfoQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetAccountInfo,
	}
}

func (query *AccountInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	if query.nodeAccountIDs.locked {
		for range query.nodeAccountIDs.slice {
			paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
			if err != nil {
				return Hbar{}, err
			}
			query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return Hbar{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.CryptoGetInfo.Header = query.pbHeader

	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_AccountInfoQueryShouldRetry,
		_CostQueryMakeRequest,
		_CostQueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_AccountInfoQueryGetMethod,
		_AccountInfoQueryMapStatusError,
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

	cost := int64(resp.query.GetCryptoGetInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	}

	return HbarFromTinybar(cost), nil
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountInfoQuery.
func (query *AccountInfoQuery) SetNodeAccountIDs(accountID []AccountID) *AccountInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

// SetQueryPayment sets the Hbar payment to pay the _Node a fee for handling this query
func (query *AccountInfoQuery) SetQueryPayment(queryPayment Hbar) *AccountInfoQuery {
	query.queryPayment = queryPayment
	return query
}

// SetMaxQueryPayment sets the maximum payment allowable for this query.
func (query *AccountInfoQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *AccountInfoQuery {
	query.maxQueryPayment = queryMaxPayment
	return query
}

func (query *AccountInfoQuery) SetMaxRetry(count int) *AccountInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *AccountInfoQuery) Execute(client *Client) (AccountInfo, error) {
	if client == nil || client.operator == nil {
		return AccountInfo{}, errNoClientProvided
	}

	var err error

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return AccountInfo{}, err
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
			return AccountInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return AccountInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "AccountInfoQuery",
			}
		}

		cost = actualCost
	}

	query.paymentTransactions = make([]*services.Transaction, 0)
	if query.nodeAccountIDs.locked {
		err = _QueryGeneratePayments(&query.Query, client, cost)
		if err != nil {
			return AccountInfo{}, err
		}
	} else {
		paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
		if err != nil {
			return AccountInfo{}, err
		}
		query.paymentTransactions = append(query.paymentTransactions, paymentTransaction)
	}

	pb := query._Build()
	pb.CryptoGetInfo.Header = query.pbHeader
	query.pb = &services.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_AccountInfoQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_AccountInfoQueryGetMethod,
		_AccountInfoQueryMapStatusError,
		_QueryMapResponse,
		query._GetLogID(),
		query.grpcDeadline,
		query.maxBackoff,
		query.minBackoff,
		query.maxRetry,
	)

	if err != nil {
		return AccountInfo{}, err
	}

	return _AccountInfoFromProtobuf(resp.query.GetCryptoGetInfo().AccountInfo)
}

func (query *AccountInfoQuery) SetMaxBackoff(max time.Duration) *AccountInfoQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *AccountInfoQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *AccountInfoQuery) SetMinBackoff(min time.Duration) *AccountInfoQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *AccountInfoQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *AccountInfoQuery) _GetLogID() string {
	timestamp := query.timestamp.UnixNano()
	if query.paymentTransactionIDs._Length() > 0 && query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart != nil {
		timestamp = query.paymentTransactionIDs._GetCurrent().(TransactionID).ValidStart.UnixNano()
	}
	return fmt.Sprintf("AccountInfoQuery:%d", timestamp)
}

func (query *AccountInfoQuery) SetPaymentTransactionID(transactionID TransactionID) *AccountInfoQuery {
	query.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return query
}
