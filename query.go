package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

type Query struct {
	pb       *services.Query
	pbHeader *services.QueryHeader //nolint

	paymentTransactionIDs *_LockableSlice
	nodeAccountIDs        *_LockableSlice
	maxQueryPayment       Hbar
	queryPayment          Hbar
	maxRetry              int

	paymentTransactions []*services.Transaction

	isPaymentRequired bool

	maxBackoff   *time.Duration
	minBackoff   *time.Duration
	grpcDeadline *time.Duration
	timestamp    time.Time
}

func _NewQuery(isPaymentRequired bool, header *services.QueryHeader) Query {
	return Query{
		pb:                    &services.Query{},
		pbHeader:              header,
		paymentTransactionIDs: _NewLockableSlice(),
		maxRetry:              10,
		nodeAccountIDs:        _NewLockableSlice(),
		paymentTransactions:   make([]*services.Transaction, 0),
		isPaymentRequired:     isPaymentRequired,
		maxQueryPayment:       NewHbar(0),
		queryPayment:          NewHbar(0),
		timestamp:             time.Now(),
	}
}

func (this *Query) SetGrpcDeadline(deadline *time.Duration) *Query {
	this.grpcDeadline = deadline
	return this
}

func (this *Query) GetGrpcDeadline() *time.Duration {
	return this.grpcDeadline
}

func (this *Query) SetNodeAccountIDs(nodeAccountIDs []AccountID) *Query {
	for _, nodeAccountID := range nodeAccountIDs {
		this.nodeAccountIDs._Push(nodeAccountID)
	}
	this.nodeAccountIDs._SetLocked(true)
	return this
}

func (this *Query) GetNodeAccountIDs() (nodeAccountIDs []AccountID) {
	nodeAccountIDs = []AccountID{}

	for _, value := range this.nodeAccountIDs.slice {
		nodeAccountIDs = append(nodeAccountIDs, value.(AccountID))
	}

	return nodeAccountIDs
}

func _QueryGetNodeAccountID(request _Request) AccountID {
	return request.query.nodeAccountIDs._GetCurrent().(AccountID)
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (this *Query) SetMaxQueryPayment(maxPayment Hbar) *Query {
	this.maxQueryPayment = maxPayment
	return this
}

// SetQueryPayment sets the payment amount for this Query.
func (this *Query) SetQueryPayment(paymentAmount Hbar) *Query {
	this.queryPayment = paymentAmount
	return this
}

func (this *Query) GetMaxRetryCount() int {
	return this.maxRetry
}

func (this *Query) SetMaxRetry(count int) *Query {
	this.maxRetry = count
	return this
}

func _QueryShouldRetry(logID string, status Status) _ExecutionState {
	logCtx.Trace().Str("requestId", logID).Str("status", status.String()).Msg("query precheck status received")
	switch status {
	case StatusPlatformTransactionNotCreated, StatusBusy:
		return executionStateRetry
	case StatusOk:
		return executionStateFinished
	}

	return executionStateError
}

func _QueryMakeRequest(request _Request) _ProtoRequest {
	if request.query.isPaymentRequired && len(request.query.paymentTransactions) > 0 {
		request.query.pbHeader.Payment = request.query.paymentTransactions[request.query.paymentTransactionIDs.index]
	}
	request.query.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	return _ProtoRequest{
		query: request.query.pb,
	}
}

func _CostQueryMakeRequest(request _Request) _ProtoRequest {
	if request.query.isPaymentRequired && len(request.query.paymentTransactions) > 0 {
		request.query.pbHeader.Payment = request.query.paymentTransactions[request.query.paymentTransactionIDs.index]
	}
	request.query.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	return _ProtoRequest{
		query: request.query.pb,
	}
}

func _QueryAdvanceRequest(request _Request) {
	if request.query.isPaymentRequired && len(request.query.paymentTransactions) > 0 {
		request.query.paymentTransactionIDs._Advance()
	}
	request.query.nodeAccountIDs._Advance()
}

func _CostQueryAdvanceRequest(request _Request) {
	request.query.paymentTransactionIDs._Advance()
	request.query.nodeAccountIDs._Advance()
}

func _QueryMapResponse(request _Request, response _Response, _ AccountID, protoRequest _ProtoRequest) (_IntermediateResponse, error) {
	return _IntermediateResponse{
		query: response.query,
	}, nil
}

func _QueryGeneratePayments(query *Query, client *Client, cost Hbar) error {
	for _, nodeID := range query.nodeAccountIDs.slice {
		transaction, err := _QueryMakePaymentTransaction(
			query.paymentTransactionIDs._GetCurrent().(TransactionID),
			nodeID.(AccountID),
			client.operator,
			cost,
		)
		if err != nil {
			return err
		}

		query.paymentTransactions = append(query.paymentTransactions, transaction)
	}

	return nil
}

func _QueryMakePaymentTransaction(transactionID TransactionID, nodeAccountID AccountID, operator *_Operator, cost Hbar) (*services.Transaction, error) {
	accountAmounts := make([]*services.AccountAmount, 0)
	accountAmounts = append(accountAmounts, &services.AccountAmount{
		AccountID: nodeAccountID._ToProtobuf(),
		Amount:    cost.tinybar,
	})
	accountAmounts = append(accountAmounts, &services.AccountAmount{
		AccountID: operator.accountID._ToProtobuf(),
		Amount:    -cost.tinybar,
	})

	body := services.TransactionBody{
		TransactionID:  transactionID._ToProtobuf(),
		NodeAccountID:  nodeAccountID._ToProtobuf(),
		TransactionFee: uint64(NewHbar(1).tinybar),
		TransactionValidDuration: &services.Duration{
			Seconds: 120,
		},
		Data: &services.TransactionBody_CryptoTransfer{
			CryptoTransfer: &services.CryptoTransferTransactionBody{
				Transfers: &services.TransferList{
					AccountAmounts: accountAmounts,
				},
			},
		},
	}

	bodyBytes, err := protobuf.Marshal(&body)
	if err != nil {
		return nil, errors.Wrap(err, "error serializing query body")
	}

	signature := operator.signer(bodyBytes)
	sigPairs := make([]*services.SignaturePair, 0)
	sigPairs = append(sigPairs, operator.publicKey._ToSignaturePairProtobuf(signature))

	return &services.Transaction{
		BodyBytes: bodyBytes,
		SigMap: &services.SignatureMap{
			SigPair: sigPairs,
		},
	}, nil
}

func (this *Query) GetPaymentTransactionID() TransactionID {
	return this.paymentTransactionIDs._GetCurrent().(TransactionID)
}

func (this *Query) SetPaymentTransactionID(transactionID TransactionID) *Query {
	this.paymentTransactionIDs._Clear()._Push(transactionID)._SetLocked(true)
	return this
}
