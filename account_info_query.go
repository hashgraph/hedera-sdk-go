package hedera

import (
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type AccountInfoQuery struct {
	Query
	accountID *AccountID
}

func NewAccountInfoQuery() *AccountInfoQuery {
	return &AccountInfoQuery{
		Query: _NewQuery(true),
	}
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
		if err := query.accountID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *AccountInfoQuery) _Build() *proto.Query_CryptoGetInfo {
	pb := proto.Query_CryptoGetInfo{
		CryptoGetInfo: &proto.CryptoGetInfoQuery{
			Header: &proto.QueryHeader{},
		},
	}

	if query.accountID != nil {
		pb.CryptoGetInfo.AccountID = query.accountID._ToProtobuf()
	}

	return &pb
}

func _AccountInfoQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetCryptoGetInfo().Header.NodeTransactionPrecheckCode))
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

func (query *AccountInfoQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	_ = query._BuildAllPaymentTransactions()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.CryptoGetInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.CryptoGetInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY
	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *AccountInfoQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionIDGenerate(client.GetOperatorAccountID()), AccountID{}, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	paymentBytes, err := protobuf.Marshal(paymentTransaction)
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.CryptoGetInfo.Header.Payment = &proto.Transaction{
		SignedTransactionBytes: paymentBytes,
	}
	pb.CryptoGetInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *AccountInfoQuery) GetCost(client *Client) (Hbar, error) {
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
		_AccountInfoQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_AccountInfoQueryGetMethod,
		_AccountInfoQueryMapStatusError,
		_QueryMapResponse,
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

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.nodeIDs = client.network._GetNodeAccountIDsForExecute()
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return AccountInfo{}, err
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

	query.actualCost = cost

	if !query.IsFrozen() {
		_, err := query.FreezeWith(client)
		if err != nil {
			return AccountInfo{}, err
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
		_AccountInfoQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_AccountInfoQueryGetMethod,
		_AccountInfoQueryMapStatusError,
		_QueryMapResponse,
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

func (query *AccountInfoQuery) IsFrozen() bool {
	return query._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (query *AccountInfoQuery) Sign(
	privateKey PrivateKey,
) *AccountInfoQuery {
	return query.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (query *AccountInfoQuery) SignWithOperator(
	client *Client,
) (*AccountInfoQuery, error) {
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
func (query *AccountInfoQuery) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountInfoQuery {
	if !query._KeyAlreadySigned(publicKey) {
		query._SignWith(publicKey, signer)
	}

	return query
}

func (query *AccountInfoQuery) Freeze() (*AccountInfoQuery, error) {
	return query.FreezeWith(nil)
}

func (query *AccountInfoQuery) FreezeWith(client *Client) (*AccountInfoQuery, error) {
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
				return &AccountInfoQuery{}, err
			}

			if query.actualCost.tinybar < actualCost.tinybar {
				return &AccountInfoQuery{}, ErrMaxQueryPaymentExceeded{
					QueryCost:       actualCost,
					MaxQueryPayment: query.actualCost,
					query:           "AccountInfoQuery",
				}
			}

			query.actualCost = actualCost
		}
	}
	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return &AccountInfoQuery{}, err
	}
	if err := query._InitPaymentTransactionID(client); err != nil {
		return query, err
	}

	return query, _QueryGeneratePayments(&query.Query, query.actualCost)
}
