package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// AccountBalanceQuery gets the balance of a CryptoCurrency account. This returns only the balance, so it is a smaller
// and faster reply than AccountInfoQuery, which returns the balance plus additional information.
type AccountBalanceQuery struct {
	Query
	accountID  *AccountID
	contractID *ContractID
}

// NewAccountBalanceQuery creates an AccountBalanceQuery query which can be used to construct and execute
// an AccountBalanceQuery.
// It is recommended that you use this for creating new instances of an AccountBalanceQuery
// instead of manually creating an instance of the struct.
func NewAccountBalanceQuery() *AccountBalanceQuery {
	header := proto.QueryHeader{}
	return &AccountBalanceQuery{
		Query: _NewQuery(false, &header),
	}
}

// SetAccountID sets the AccountID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this _Method.
func (query *AccountBalanceQuery) SetAccountID(accountID AccountID) *AccountBalanceQuery {
	query.accountID = &accountID
	return query
}

func (query *AccountBalanceQuery) GetAccountID() AccountID {
	if query.accountID == nil {
		return AccountID{}
	}

	return *query.accountID
}

// SetContractID sets the ContractID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this _Method.
func (query *AccountBalanceQuery) SetContractID(contractID ContractID) *AccountBalanceQuery {
	query.contractID = &contractID
	return query
}

func (query *AccountBalanceQuery) GetContractID() ContractID {
	if query.contractID == nil {
		return ContractID{}
	}

	return *query.contractID
}

func (query *AccountBalanceQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.accountID != nil {
		if err := query.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if query.contractID != nil {
		if err := query.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *AccountBalanceQuery) _Build() *proto.Query_CryptogetAccountBalance {
	pb := proto.CryptoGetAccountBalanceQuery{Header: &proto.QueryHeader{}}

	if query.accountID != nil {
		pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_AccountID{
			AccountID: query.accountID._ToProtobuf(),
		}
	}

	if query.contractID != nil {
		pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_ContractID{
			ContractID: query.contractID._ToProtobuf(),
		}
	}

	return &proto.Query_CryptogetAccountBalance{
		CryptogetAccountBalance: &pb,
	}
}

func (query *AccountBalanceQuery) GetCost(client *Client) (Hbar, error) {
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
	} else if len(query.Query.GetNodeAccountIDs()) == 1 {
		nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
		if err != nil {
			return Hbar{}, err
		}

		query.nodeAccountIDs = append(query.nodeAccountIDs, nodeAccountIDs[0])
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
	pb.CryptogetAccountBalance.Header = query.pbHeader

	query.pb = &proto.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_AccountBalanceQueryShouldRetry,
		_CostQueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_AccountBalanceQueryGetMethod,
		_AccountBalanceQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetCryptogetAccountBalance().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _AccountBalanceQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetCryptogetAccountBalance().Header.NodeTransactionPrecheckCode))
}

func _AccountBalanceQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetCryptogetAccountBalance().Header.NodeTransactionPrecheckCode),
	}
}

func _AccountBalanceQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().CryptoGetBalance,
	}
}

func (query *AccountBalanceQuery) Execute(client *Client) (AccountBalance, error) {
	if client == nil {
		return AccountBalance{}, errNoClientProvided
	}

	var err error
	if len(query.Query.GetNodeAccountIDs()) == 0 {
		nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
		if err != nil {
			return AccountBalance{}, err
		}

		query.SetNodeAccountIDs(nodeAccountIDs)
	} else if len(query.Query.GetNodeAccountIDs()) == 1 {
		nodeAccountIDs, err := client.network._GetNodeAccountIDsForExecute()
		if err != nil {
			return AccountBalance{}, err
		}

		query.nodeAccountIDs = append(query.nodeAccountIDs, nodeAccountIDs[0])
	}

	err = query._ValidateNetworkOnIDs(client)
	if err != nil {
		return AccountBalance{}, err
	}

	query.nextPaymentTransactionIndex = 0
	query.paymentTransactions = make([]*proto.Transaction, 0)

	pb := query._Build()
	pb.CryptogetAccountBalance.Header = query.pbHeader
	query.pb = &proto.Query{
		Query: pb,
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_AccountBalanceQueryShouldRetry,
		_QueryMakeRequest,
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_AccountBalanceQueryGetMethod,
		_AccountBalanceQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return AccountBalance{}, err
	}

	return _AccountBalanceFromProtobuf(resp.query.GetCryptogetAccountBalance()), nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *AccountBalanceQuery) SetMaxQueryPayment(maxPayment Hbar) *AccountBalanceQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *AccountBalanceQuery) SetQueryPayment(paymentAmount Hbar) *AccountBalanceQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountBalanceQuery.
func (query *AccountBalanceQuery) SetNodeAccountIDs(accountID []AccountID) *AccountBalanceQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *AccountBalanceQuery) SetMaxRetry(count int) *AccountBalanceQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *AccountBalanceQuery) SetMaxBackoff(max time.Duration) *AccountBalanceQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *AccountBalanceQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *AccountBalanceQuery) SetMinBackoff(min time.Duration) *AccountBalanceQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *AccountBalanceQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}
