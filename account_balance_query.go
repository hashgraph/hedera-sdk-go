package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

// AccountBalanceQuery gets the balance of a CryptoCurrency account. This returns only the balance, so it is a smaller
// and faster reply than AccountInfoQuery, which returns the balance plus additional information.
type AccountBalanceQuery struct {
	Query
	accountID  AccountID
	contractID ContractID
}

// NewAccountBalanceQuery creates an AccountBalanceQuery query which can be used to construct and execute
// an AccountBalanceQuery.
// It is recommended that you use this for creating new instances of an AccountBalanceQuery
// instead of manually creating an instance of the struct.
func NewAccountBalanceQuery() *AccountBalanceQuery {
	return &AccountBalanceQuery{
		Query: newQuery(false),
	}
}

// SetAccountID sets the AccountID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this method.
func (query *AccountBalanceQuery) SetAccountID(id AccountID) *AccountBalanceQuery {
	query.accountID = id

	return query
}

func (query *AccountBalanceQuery) GetAccountID() AccountID {
	return query.accountID
}

// SetContractID sets the ContractID for which you wish to query the balance.
//
// Note: you can only query an Account or Contract but not both -- if a Contract ID or Account ID has already been set,
// it will be overwritten by this method.
func (query *AccountBalanceQuery) SetContractID(id ContractID) *AccountBalanceQuery {
	query.contractID = id

	return query
}

func (query *AccountBalanceQuery) GetContractID() ContractID {
	return query.contractID
}

func (query *AccountBalanceQuery) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = query.accountID.Validate(client)
	if err != nil {
		return err
	}
	err = query.contractID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (query *AccountBalanceQuery) build() *proto.Query_CryptogetAccountBalance {
	pb := proto.CryptoGetAccountBalanceQuery{Header: &proto.QueryHeader{}}

	if !query.accountID.isZero() {
		pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_AccountID{
			AccountID: query.accountID.toProtobuf(),
		}
	}

	if !query.contractID.isZero() {
		pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_ContractID{
			ContractID: query.contractID.toProtobuf(),
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
		accountBalanceQuery_shouldRetry,
		protoReq,
		costQuery_advanceRequest,
		costQuery_getNodeAccountID,
		accountBalanceQuery_getMethod,
		accountBalanceQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetCryptogetAccountBalance().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func (query *AccountBalanceQuery) queryMakeRequest() protoRequest {
	pb := query.build()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.CryptogetAccountBalance.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.CryptogetAccountBalance.Header.ResponseType = proto.ResponseType_ANSWER_ONLY
	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *AccountBalanceQuery) costQueryMakeRequest(client *Client) (protoRequest, error) {
	pb := query.build()

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return protoRequest{}, err
	}

	pb.CryptogetAccountBalance.Header.Payment = paymentTransaction
	pb.CryptogetAccountBalance.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return protoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func accountBalanceQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetCryptogetAccountBalance().Header.NodeTransactionPrecheckCode))
}

func accountBalanceQuery_mapStatusError(_ request, response response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetCryptogetAccountBalance().Header.NodeTransactionPrecheckCode),
	}
}

func accountBalanceQuery_getMethod(_ request, channel *channel) method {
	return method{
		query: channel.getCrypto().CryptoGetBalance,
	}
}

func (query *AccountBalanceQuery) Execute(client *Client) (AccountBalance, error) {
	if client == nil {
		return AccountBalance{}, errNoClientProvided
	}

	query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return AccountBalance{}, err
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		accountBalanceQuery_shouldRetry,
		query.queryMakeRequest(),
		query_advanceRequest,
		query_getNodeAccountID,
		accountBalanceQuery_getMethod,
		accountBalanceQuery_mapStatusError,
		query_mapResponse,
	)

	if err != nil {
		return AccountBalance{}, err
	}

	return accountBalanceFromProtobuf(resp.query.GetCryptogetAccountBalance()), nil
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

// SetNodeAccountIDs sets the node AccountID for this AccountBalanceQuery.
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
