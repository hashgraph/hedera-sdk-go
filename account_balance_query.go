package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
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
	return &AccountBalanceQuery{
		Query: _NewQuery(false),
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

func (query *AccountBalanceQuery) _Build() *services.Query_CryptogetAccountBalance {
	pb := services.CryptoGetAccountBalanceQuery{Header: &services.QueryHeader{}}

	if query.accountID != nil {
		pb.BalanceSource = &services.CryptoGetAccountBalanceQuery_AccountID{
			AccountID: query.accountID._ToProtobuf(),
		}
	}

	if query.contractID != nil {
		pb.BalanceSource = &services.CryptoGetAccountBalanceQuery_ContractID{
			ContractID: query.contractID._ToProtobuf(),
		}
	}

	return &services.Query_CryptogetAccountBalance{
		CryptogetAccountBalance: &pb,
	}
}

func (query *AccountBalanceQuery) GetCost(client *Client) (Hbar, error) {
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
		_AccountBalanceQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
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

func (query *AccountBalanceQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.CryptogetAccountBalance.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.CryptogetAccountBalance.Header.ResponseType = services.ResponseType_ANSWER_ONLY
	return _ProtoRequest{
		query: &services.Query{
			Query: pb,
		},
	}
}

func (query *AccountBalanceQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.CryptogetAccountBalance.Header.Payment = paymentTransaction
	pb.CryptogetAccountBalance.Header.ResponseType = services.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &services.Query{
			Query: pb,
		},
	}, nil
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

	query.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return AccountBalance{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_AccountBalanceQueryShouldRetry,
		query._QueryMakeRequest(),
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
