package hedera

import (
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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
	return &AccountStakersQuery{
		Query: _NewQuery(true),
	}
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
		if err := query.accountID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *AccountStakersQuery) _Build() *proto.Query_CryptoGetProxyStakers {
	pb := proto.Query_CryptoGetProxyStakers{
		CryptoGetProxyStakers: &proto.CryptoGetStakersQuery{
			Header: &proto.QueryHeader{},
		},
	}

	if query.accountID != nil {
		pb.CryptoGetProxyStakers.AccountID = query.accountID._ToProtobuf()
	}

	return &pb
}

func (query *AccountStakersQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	_ = query._BuildAllPaymentTransactions()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.CryptoGetProxyStakers.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.CryptoGetProxyStakers.Header.ResponseType = proto.ResponseType_ANSWER_ONLY
	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *AccountStakersQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionIDGenerate(client.GetOperatorAccountID()), AccountID{}, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	paymentBytes, err := protobuf.Marshal(paymentTransaction)
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.CryptoGetProxyStakers.Header.Payment = &proto.Transaction{
		SignedTransactionBytes: paymentBytes,
	}
	pb.CryptoGetProxyStakers.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *AccountStakersQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.nodeIDs = client.network._GetNodeAccountIDsForExecute()
	}

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
		_AccountStakersQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_AccountStakersQueryGetMethod,
		_AccountStakersQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetCryptoGetProxyStakers().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _AccountStakersQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetCryptoGetProxyStakers().Header.NodeTransactionPrecheckCode))
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

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.nodeIDs = client.network._GetNodeAccountIDsForExecute()
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return []Transfer{}, err
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

	query.actualCost = cost

	if !query.IsFrozen() {
		_, err := query.FreezeWith(client)
		if err != nil {
			return []Transfer{}, err
		}
	}

	transactionID := query.paymentTransactionID

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		query.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_AccountStakersQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_AccountStakersQueryGetMethod,
		_AccountStakersQueryMapStatusError,
		_QueryMapResponse,
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

func (query *AccountStakersQuery) IsFrozen() bool {
	return query._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (query *AccountStakersQuery) Sign(
	privateKey PrivateKey,
) *AccountStakersQuery {
	return query.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (query *AccountStakersQuery) SignWithOperator(
	client *Client,
) (*AccountStakersQuery, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !query.IsFrozen() {
		_, err := query.FreezeWith(client)
		if err != nil {
			return query, err
		}
	}
	return query.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (query *AccountStakersQuery) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountStakersQuery {
	if !query._KeyAlreadySigned(publicKey) {
		query._SignWith(publicKey, signer)
	}

	return query
}

func (query *AccountStakersQuery) Freeze() (*AccountStakersQuery, error) {
	return query.FreezeWith(nil)
}

func (query *AccountStakersQuery) FreezeWith(client *Client) (*AccountStakersQuery, error) {
	if query.IsFrozen() {
		return query, nil
	}
	if query.actualCost.AsTinybar() == 0 {
		if query.queryPayment.tinybar != 0 {
			query.actualCost = query.queryPayment
		} else {
			if query.maxQueryPayment.tinybar == 0 {
				query.actualCost = client.maxQueryPayment
			} else {
				query.actualCost = query.maxQueryPayment
			}

			actualCost, err := query.GetCost(client)
			if err != nil {
				return &AccountStakersQuery{}, err
			}

			if query.actualCost.tinybar < actualCost.tinybar {
				return &AccountStakersQuery{}, ErrMaxQueryPaymentExceeded{
					QueryCost:       actualCost,
					MaxQueryPayment: query.actualCost,
					query:           "AccountStakersQuery",
				}
			}

			query.actualCost = actualCost
		}
	}
	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return &AccountStakersQuery{}, err
	}
	if err := query._InitPaymentTransactionID(client); err != nil {
		return query, err
	}

	return query, _QueryGeneratePayments(&query.Query, query.actualCost)
}
