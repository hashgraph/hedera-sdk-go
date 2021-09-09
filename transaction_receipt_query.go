package hedera

import (
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TransactionReceiptQuery struct {
	Query
	transactionID *TransactionID
	duplicates    *bool
}

func NewTransactionReceiptQuery() *TransactionReceiptQuery {
	return &TransactionReceiptQuery{
		Query: _NewQuery(false),
	}
}

func (query *TransactionReceiptQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if err := query.transactionID.AccountID.Validate(client); err != nil {
		return err
	}

	return nil
}

func (query *TransactionReceiptQuery) _Build() *proto.Query_TransactionGetReceipt {
	body := &proto.TransactionGetReceiptQuery{
		Header: &proto.QueryHeader{},
	}

	if query.transactionID.AccountID != nil {
		body.TransactionID = query.transactionID._ToProtobuf()
	}

	if query.duplicates != nil {
		body.IncludeDuplicates = *query.duplicates
	}

	return &proto.Query_TransactionGetReceipt{
		TransactionGetReceipt: body,
	}
}

func (query *TransactionReceiptQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	_ = query._BuildAllPaymentTransactions()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.TransactionGetReceipt.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.TransactionGetReceipt.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *TransactionReceiptQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionIDGenerate(client.GetOperatorAccountID()), AccountID{}, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}
	paymentBytes, err := protobuf.Marshal(paymentTransaction)
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.TransactionGetReceipt.Header.Payment = &proto.Transaction{
		SignedTransactionBytes: paymentBytes,
	}
	pb.TransactionGetReceipt.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *TransactionReceiptQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil {
		return Hbar{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
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
		_TransactionReceiptQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_TransactionReceiptQueryGetMethod,
		_TransactionReceiptQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetTransactionGetReceipt().Header.Cost)
	return HbarFromTinybar(cost), nil
}

func _TransactionReceiptQueryShouldRetry(request _Request, response _Response) _ExecutionState {
	switch Status(response.query.GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode()) {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusReceiptNotFound, StatusRecordNotFound:
		return executionStateRetry
	case StatusOk:
		break
	default:
		return executionStateError
	}

	switch Status(response.query.GetTransactionGetReceipt().GetReceipt().GetStatus()) {
	case StatusBusy, StatusUnknown, StatusOk, StatusReceiptNotFound, StatusRecordNotFound:
		return executionStateRetry
	case StatusSuccess:
		return executionStateFinished
	default:
		return executionStateError
	}
}

func _TransactionReceiptQueryMapStatusError(request _Request, response _Response) error {
	switch Status(response.query.GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode()) {
	case StatusPlatformTransactionNotCreated, StatusBusy, StatusUnknown, StatusReceiptNotFound, StatusRecordNotFound, StatusOk:
		break
	default:
		return ErrHederaPreCheckStatus{
			Status: Status(response.query.GetTransactionGetReceipt().GetHeader().GetNodeTransactionPrecheckCode()),
		}
	}

	return ErrHederaReceiptStatus{
		Status:  Status(response.query.GetTransactionGetReceipt().GetReceipt().GetStatus()),
		Receipt: _TransactionReceiptFromProtobuf(response.query.GetTransactionGetReceipt().GetReceipt()),
	}
}

func _TransactionReceiptQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetCrypto().GetTransactionReceipts,
	}
}

func (query *TransactionReceiptQuery) SetTransactionID(transactionID TransactionID) *TransactionReceiptQuery {
	query.transactionID = &transactionID
	return query
}

func (query *TransactionReceiptQuery) GetTransactionID() TransactionID {
	if query.transactionID == nil {
		return TransactionID{}
	}

	return *query.transactionID
}

func (query *TransactionReceiptQuery) SetNodeAccountIDs(accountID []AccountID) *TransactionReceiptQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *TransactionReceiptQuery) SetQueryPayment(queryPayment Hbar) *TransactionReceiptQuery {
	query.queryPayment = queryPayment
	return query
}

func (query *TransactionReceiptQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionReceiptQuery {
	query.maxQueryPayment = queryMaxPayment
	return query
}

func (query *TransactionReceiptQuery) SetMaxRetry(count int) *TransactionReceiptQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *TransactionReceiptQuery) SetMaxBackoff(max time.Duration) *TransactionReceiptQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *TransactionReceiptQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *TransactionReceiptQuery) SetMinBackoff(min time.Duration) *TransactionReceiptQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *TransactionReceiptQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *TransactionReceiptQuery) Execute(client *Client) (TransactionReceipt, error) {
	if client == nil {
		return TransactionReceipt{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.SetNodeAccountIDs(client.network._GetNodeAccountIDsForExecute())
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return TransactionReceipt{}, err
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
			return TransactionReceipt{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return TransactionReceipt{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "TransactionReceiptQuery",
			}
		}

		cost = actualCost
	}

	query.actualCost = cost
	if !query.IsFrozen() {
		_, err := query.FreezeWith(client)
		if err != nil {
			return TransactionReceipt{}, err
		}
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_TransactionReceiptQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_TransactionReceiptQueryGetMethod,
		_TransactionReceiptQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		if precheckErr, ok := err.(ErrHederaPreCheckStatus); ok {
			return TransactionReceipt{}, _NewErrHederaReceiptStatus(precheckErr.TxID, precheckErr.Status)
		}
		return TransactionReceipt{}, err
	}

	return _TransactionReceiptFromProtobuf(resp.query.GetTransactionGetReceipt().GetReceipt()), nil
}

func (query *TransactionReceiptQuery) IsFrozen() bool {
	return query._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (query *TransactionReceiptQuery) Sign(
	privateKey PrivateKey,
) *TransactionReceiptQuery {
	return query.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (query *TransactionReceiptQuery) SignWithOperator(
	client *Client,
) (*TransactionReceiptQuery, error) {
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
func (query *TransactionReceiptQuery) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TransactionReceiptQuery {
	if !query._KeyAlreadySigned(publicKey) {
		query._SignWith(publicKey, signer)
	}

	return query
}

func (query *TransactionReceiptQuery) Freeze() (*TransactionReceiptQuery, error) {
	return query.FreezeWith(nil)
}

func (query *TransactionReceiptQuery) FreezeWith(client *Client) (*TransactionReceiptQuery, error) {
	if query.IsFrozen() {
		return query, nil
	}
	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TransactionReceiptQuery{}, err
	}

	return query, nil
}
