package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"
)

type Query struct {
	paymentTransactionID        TransactionID
	nodeIDs                     []AccountID
	maxQueryPayment             Hbar
	queryPayment                Hbar
	nextPaymentTransactionIndex int
	nextTransactionIndex        int
	maxRetry                    int

	paymentTransactions []*proto.Transaction

	isPaymentRequired bool

	maxBackoff *time.Duration
	minBackoff *time.Duration
}

func newQuery(isPaymentRequired bool) Query {
	return Query{
		paymentTransactionID: TransactionID{},
		nextTransactionIndex: 0,
		maxRetry:             10,
		paymentTransactions:  make([]*proto.Transaction, 0),
		isPaymentRequired:    isPaymentRequired,
		maxQueryPayment:      NewHbar(0),
		queryPayment:         NewHbar(0),
	}
}

func (query *Query) SetNodeAccountIDs(accountID []AccountID) *Query {
	query.nodeIDs = append(query.nodeIDs, accountID...)
	return query
}

func (query *Query) GetNodeAccountIDs() []AccountID {
	return query.nodeIDs
}

func _QueryGetNodeAccountID(request _Request) AccountID {
	if len(request.query.nodeIDs) > 0 {
		return request.query.nodeIDs[request.query.nextPaymentTransactionIndex]
	}

	panic("Query _Node AccountID's not set before executing")
}

func _CostQueryGetNodeAccountID(request _Request) AccountID {
	return request.query.nodeIDs[request.query.nextPaymentTransactionIndex]
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

func _QueryAdvanceRequest(request _Request) {
	if request.query.isPaymentRequired && len(request.query.paymentTransactions) > 0 {
		request.query.nextPaymentTransactionIndex = (request.query.nextPaymentTransactionIndex + 1) % len(request.query.paymentTransactions)
	}
}

func _CostQueryAdvanceRequest(request _Request) {
	request.query.nextPaymentTransactionIndex = (request.query.nextPaymentTransactionIndex + 1) % len(request.query.nodeIDs)
}

func _QueryMapResponse(request _Request, response _Response, _ AccountID, protoRequest _ProtoRequest) (_IntermediateResponse, error) {
	return _IntermediateResponse{
		query: response.query,
	}, nil
}

func _QueryGeneratePayments(query *Query, client *Client, cost Hbar) error {
	for _, nodeID := range query.nodeIDs {
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

func _QueryMakePaymentTransaction(transactionID TransactionID, nodeAccountID AccountID, operator *_Operator, cost Hbar) (*proto.Transaction, error) {
	accountAmounts := make([]*proto.AccountAmount, 0)
	accountAmounts = append(accountAmounts, &proto.AccountAmount{
		AccountID: nodeAccountID.toProtobuf(),
		Amount:    cost.tinybar,
	})
	accountAmounts = append(accountAmounts, &proto.AccountAmount{
		AccountID: operator.accountID.toProtobuf(),
		Amount:    -cost.tinybar,
	})

	body := proto.TransactionBody{
		TransactionID:  transactionID.toProtobuf(),
		NodeAccountID:  nodeAccountID.toProtobuf(),
		TransactionFee: uint64(NewHbar(1).tinybar),
		TransactionValidDuration: &proto.Duration{
			Seconds: 120,
		},
		Data: &proto.TransactionBody_CryptoTransfer{
			CryptoTransfer: &proto.CryptoTransferTransactionBody{
				Transfers: &proto.TransferList{
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
	sigPairs := make([]*proto.SignaturePair, 0)
	sigPairs = append(sigPairs, operator.publicKey.toSignaturePairProtobuf(signature))

	return &proto.Transaction{
		BodyBytes: bodyBytes,
		SigMap: &proto.SignatureMap{
			SigPair: sigPairs,
		},
	}, nil
}
