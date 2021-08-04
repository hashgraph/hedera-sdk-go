package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// AccountBalanceQuery gets the balance of a CryptoCurrency account. This returns only the balance, so it is a smaller
// and faster reply than AccountInfoQuery, which returns the balance plus additional information.
type AccountBalanceQuery struct {
	Query
	pb         *proto.CryptoGetAccountBalanceQuery
	accountID  AccountID
	contractID ContractID
}

// NewAccountBalanceQuery creates an AccountBalanceQuery query which can be used to construct and execute
// an AccountBalanceQuery.
// It is recommended that you use this for creating new instances of an AccountBalanceQuery
// instead of manually creating an instance of the struct.
func NewAccountBalanceQuery() *AccountBalanceQuery {
	header := proto.QueryHeader{}
	query := newQuery(false, &header)
	pb := proto.CryptoGetAccountBalanceQuery{Header: &header}
	query.pb.Query = &proto.Query_CryptogetAccountBalance{
		CryptogetAccountBalance: &pb,
	}

	return &AccountBalanceQuery{
		Query: query,
		pb:    &pb,
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
	if !client.autoValidateChecksums {
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

func (query *AccountBalanceQuery) build() *AccountBalanceQuery {
	if !query.accountID.isZero() {
		query.pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_AccountID{
			AccountID: query.accountID.toProtobuf(),
		}
	}

	if !query.contractID.isZero() {
		query.pb.BalanceSource = &proto.CryptoGetAccountBalanceQuery_ContractID{
			ContractID: query.contractID.toProtobuf(),
		}
	}

	return query
}

func (query *AccountBalanceQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	paymentTransaction, err := query_makePaymentTransaction(TransactionID{}, AccountID{}, client.operator, Hbar{})
	if err != nil {
		return Hbar{}, err
	}

	query.pbHeader.Payment = paymentTransaction
	query.pbHeader.ResponseType = proto.ResponseType_COST_ANSWER
	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())
	}

	err = query.validateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	query.build()

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		accountBalanceQuery_shouldRetry,
		costQuery_makeRequest,
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

func accountBalanceQuery_shouldRetry(_ request, response response) executionState {
	return query_shouldRetry(Status(response.query.GetCryptogetAccountBalance().Header.NodeTransactionPrecheckCode))
}

func accountBalanceQuery_mapStatusError(_ request, response response, networkName *NetworkName) error {
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
	if client == nil || client.operator == nil {
		return AccountBalance{}, errNoClientProvided
	}

	query.SetNodeAccountIDs(client.network.getNodeAccountIDsForExecute())

	err := query.validateNetworkOnIDs(client)
	if err != nil {
		return AccountBalance{}, err
	}

	query.build()

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		accountBalanceQuery_shouldRetry,
		query_makeRequest,
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
