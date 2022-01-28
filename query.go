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

	paymentTransactionID        TransactionID
	nodeAccountIDs              []AccountID
	maxQueryPayment             Hbar
	queryPayment                Hbar
	nextPaymentTransactionIndex int
	nextNodeIndex               int
	nextTransactionIndex        int
	maxRetry                    int

	paymentTransactions []*services.Transaction

	isPaymentRequired bool

	maxBackoff *time.Duration
	minBackoff *time.Duration
}

func _NewQuery(isPaymentRequired bool, header *services.QueryHeader) Query {
	return Query{
		pb:                   &services.Query{},
		pbHeader:             header,
		paymentTransactionID: TransactionID{},
		nextTransactionIndex: 0,
		nextNodeIndex:        0,
		maxRetry:             10,
		paymentTransactions:  make([]*services.Transaction, 0),
		isPaymentRequired:    isPaymentRequired,
		maxQueryPayment:      NewHbar(0),
		queryPayment:         NewHbar(0),
	}
}

func (query *Query) SetNodeAccountIDs(nodeAccountIDs []AccountID) *Query {
	for _, nodeAccountID := range nodeAccountIDs {
		if nodeAccountID._IsZero() {
			panic("cannot set node account ID of 0.0.0")
		}
	}
	query.nodeAccountIDs = nodeAccountIDs
	return query
}

func (query *Query) GetNodeAccountIDs() []AccountID {
	return query.nodeAccountIDs
}

func _QueryGetNodeAccountID(request _Request) AccountID {
	if len(request.query.nodeAccountIDs) > 0 {
		return request.query.nodeAccountIDs[request.query.nextNodeIndex]
	}

	panic("Query node AccountID's not set before executing")
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *Query) SetMaxQueryPayment(maxPayment Hbar) *Query {
	query.maxQueryPayment = maxPayment
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *Query) SetQueryPayment(paymentAmount Hbar) *Query {
	query.queryPayment = paymentAmount
	return query
}

func (query *Query) GetMaxRetryCount() int {
	return query.maxRetry
}

func (query *Query) SetMaxRetry(count int) *Query {
	query.maxRetry = count
	return query
}

func _QueryShouldRetry(status Status) _ExecutionState {
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
		request.query.pbHeader.Payment = request.query.paymentTransactions[request.query.nextPaymentTransactionIndex]
	}
	request.query.pbHeader.ResponseType = services.ResponseType_ANSWER_ONLY

	return _ProtoRequest{
		query: request.query.pb,
	}
}

func _CostQueryMakeRequest(request _Request) _ProtoRequest {
	if request.query.isPaymentRequired && len(request.query.paymentTransactions) > 0 {
		request.query.pbHeader.Payment = request.query.paymentTransactions[request.query.nextPaymentTransactionIndex]
	}
	request.query.pbHeader.ResponseType = services.ResponseType_COST_ANSWER
	return _ProtoRequest{
		query: request.query.pb,
	}
}

func _QueryAdvanceRequest(request _Request) {
	if request.query.isPaymentRequired && len(request.query.paymentTransactions) > 0 {
		request.query.nextPaymentTransactionIndex = (request.query.nextPaymentTransactionIndex + 1) % len(request.query.paymentTransactions)
	}
	length := len(request.query.nodeAccountIDs)
	currentIndex := request.query.nextNodeIndex
	request.query.nextNodeIndex = (currentIndex + 1) % length
}

func _QueryMapResponse(request _Request, response _Response, _ AccountID, protoRequest _ProtoRequest) (_IntermediateResponse, error) {
	return _IntermediateResponse{
		query: response.query,
	}, nil
}

func _QueryGeneratePayments(query *Query, client *Client, cost Hbar) error {
	for _, nodeID := range query.nodeAccountIDs {
		transaction, err := _QueryMakePaymentTransaction(
			query.paymentTransactionID,
			nodeID,
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
