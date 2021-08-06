package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"github.com/pkg/errors"
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

func query_getNodeAccountID(request request) AccountID {
	if len(request.query.nodeIDs) > 0 {
		return request.query.nodeIDs[request.query.nextPaymentTransactionIndex]
	} else {
		panic("Query node AccountID's not set before executing")
	}
}

func costQuery_getNodeAccountID(request request) AccountID {
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

func (query *Query) getTransactionID(paymentAmount Hbar) TransactionID {
	return query.paymentTransactionID
}

func (query *Query) getIsPaymentRequired() bool {
	return true
}

func query_shouldRetry(status Status) executionState {
	switch status {
	case StatusPlatformTransactionNotCreated, StatusBusy:
		return executionStateRetry
	case StatusOk:
		return executionStateFinished
	}

	return executionStateError
}

func query_advanceRequest(request request) {
	if request.query.isPaymentRequired && len(request.query.paymentTransactions) > 0 {
		request.query.nextPaymentTransactionIndex = (request.query.nextPaymentTransactionIndex + 1) % len(request.query.paymentTransactions)
	}
}

func costQuery_advanceRequest(request request) {
	request.query.nextPaymentTransactionIndex = (request.query.nextPaymentTransactionIndex + 1) % len(request.query.nodeIDs)
}

func query_mapResponse(request request, response response, _ AccountID, protoRequest protoRequest) (intermediateResponse, error) {
	return intermediateResponse{
		query: response.query,
	}, nil
}

func query_generatePayments(query *Query, client *Client, cost Hbar) error {
	for _, nodeID := range query.nodeIDs {
		transaction, err := query_makePaymentTransaction(
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

func query_makePaymentTransaction(transactionID TransactionID, nodeAccountID AccountID, operator *operator, cost Hbar) (*proto.Transaction, error) {
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
