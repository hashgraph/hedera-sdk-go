package hedera

import (
	"bytes"
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
	actualCost                  Hbar
	nextPaymentTransactionIndex int
	maxRetry                    int

	paymentTransactions       []*proto.Transaction
	signedPaymentTransactions []*proto.SignedTransaction
	paymentTransactionIDs     []TransactionID

	publicKeys         []PublicKey
	transactionSigners []TransactionSigner

	isPaymentRequired bool

	maxBackoff *time.Duration
	minBackoff *time.Duration
}

func _NewQuery(isPaymentRequired bool) Query {
	return Query{
		paymentTransactionID:        TransactionID{},
		nextPaymentTransactionIndex: 0,
		maxRetry:                    10,
		actualCost:                  Hbar{},
		paymentTransactions:         make([]*proto.Transaction, 0),
		signedPaymentTransactions:   make([]*proto.SignedTransaction, 0),
		paymentTransactionIDs:       make([]TransactionID, 0),
		nodeIDs:                     make([]AccountID, 0),
		isPaymentRequired:           isPaymentRequired,
		maxQueryPayment:             NewHbar(0),
		queryPayment:                NewHbar(0),
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

func _QueryGeneratePayments(query *Query, cost Hbar) error {
	for _, nodeID := range query.nodeIDs {
		transaction, err := _QueryMakePaymentTransaction(
			query.paymentTransactionID,
			nodeID,
			cost,
		)
		if err != nil {
			return err
		}

		query.signedPaymentTransactions = append(query.signedPaymentTransactions, transaction)
	}

	return nil
}

func _QueryMakePaymentTransaction(transactionID TransactionID, nodeAccountID AccountID, cost Hbar) (*proto.SignedTransaction, error) {
	accountAmounts := make([]*proto.AccountAmount, 0)
	accountAmounts = append(accountAmounts, &proto.AccountAmount{
		AccountID: nodeAccountID._ToProtobuf(),
		Amount:    cost.tinybar,
	})
	accountAmounts = append(accountAmounts, &proto.AccountAmount{
		AccountID: transactionID.AccountID._ToProtobuf(),
		Amount:    -cost.tinybar,
	})

	body := proto.TransactionBody{
		TransactionID:  transactionID._ToProtobuf(),
		NodeAccountID:  nodeAccountID._ToProtobuf(),
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

	tx := &proto.SignedTransaction{
		BodyBytes: bodyBytes,
		SigMap: &proto.SignatureMap{
			SigPair: make([]*proto.SignaturePair, 0),
		},
	}

	return tx, nil
}

func (query *Query) _IsFrozen() bool {
	return len(query.signedPaymentTransactions) > 0
}

func (query *Query) _SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) {
	query.publicKeys = append(query.publicKeys, publicKey)
	query.transactionSigners = append(query.transactionSigners, signer)
}

func (query *Query) _KeyAlreadySigned(
	pk PublicKey,
) bool {
	for _, key := range query.publicKeys {
		if key.String() == pk.String() {
			return true
		}
	}

	return false
}

func (query *Query) _InitPaymentTransactionID(client *Client) error {
	if len(query.paymentTransactionIDs) == 0 {
		if client != nil {
			if client.operator != nil {
				query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)
			} else {
				return errNoClientOrTransactionID
			}
		} else {
			return errNoClientOrTransactionID
		}
	}

	return nil
}

func (query *Query) _BuildAllPaymentTransactions() error {
	for i := 0; i < len(query.signedPaymentTransactions); i++ {
		err := query._BuildPaymentTransaction(i)
		if err != nil {
			return err
		}
	}

	return nil
}

func (query *Query) _BuildPaymentTransaction(index int) error {
	if len(query.paymentTransactions) < index {
		for i := len(query.paymentTransactions); i < index; i++ {
			query.paymentTransactions = append(query.paymentTransactions, nil)
		}
	} else if len(query.paymentTransactions) > index &&
		query.paymentTransactions[index] != nil &&
		query.paymentTransactions[index].SignedTransactionBytes != nil {
		return nil
	}

	query._SignPaymentTransaction(index)

	data, err := protobuf.Marshal(query.signedPaymentTransactions[index])
	if err != nil {
		return errors.Wrap(err, "failed to serialize transactions for building")
	}

	query.paymentTransactions = append(query.paymentTransactions, &proto.Transaction{
		SignedTransactionBytes: data,
	})

	return nil
}

func (query *Query) _SignPaymentTransaction(index int) {
	if len(query.signedPaymentTransactions[index].SigMap.SigPair) != 0 {
		for i, key := range query.publicKeys {
			if query.transactionSigners[i] != nil && bytes.Equal(query.signedPaymentTransactions[index].SigMap.SigPair[0].PubKeyPrefix, key.keyData) {
				return
			}
		}
	}

	bodyBytes := query.signedPaymentTransactions[index].GetBodyBytes()

	for i := 0; i < len(query.publicKeys); i++ {
		publicKey := query.publicKeys[i]
		signer := query.transactionSigners[i]

		if signer == nil {
			continue
		}

		query.signedPaymentTransactions[index].SigMap.SigPair = append(query.signedPaymentTransactions[index].SigMap.SigPair, publicKey._ToSignaturePairProtobuf(signer(bodyBytes)))
	}
}
