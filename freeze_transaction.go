package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

type FreezeTransaction struct {
	Transaction
	pb *proto.FreezeTransactionBody
}

func NewFreezeTransaction() *FreezeTransaction {
	pb := &proto.FreezeTransactionBody{}

	transaction := FreezeTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func freezeTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) FreezeTransaction {
	return FreezeTransaction{
		Transaction: transaction,
		pb:          pb.GetFreeze(),
	}
}

func (transaction *FreezeTransaction) SetStartTime(startTime time.Time) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.pb.StartHour = int32(startTime.Hour())
	transaction.pb.StartMin = int32(startTime.Minute())
	return transaction
}

func (transaction *FreezeTransaction) GetStartTime() time.Time {
	t1 := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		int(transaction.pb.StartHour), int(transaction.pb.StartMin),
		0, time.Now().Nanosecond(), time.Now().Location(),
	)
	return t1
}

func (transaction *FreezeTransaction) SetEndTime(endTime time.Time) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.pb.StartHour = int32(endTime.Hour())
	transaction.pb.StartMin = int32(endTime.Minute())
	return transaction
}

func (transaction *FreezeTransaction) GetEndTime() time.Time {
	t1 := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		int(transaction.pb.EndHour), int(transaction.pb.EndMin),
		0, time.Now().Nanosecond(), time.Now().Location(),
	)
	return t1
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func freezeTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getFreeze().Freeze,
	}
}

func (transaction *FreezeTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *FreezeTransaction) Sign(
	privateKey PrivateKey,
) *FreezeTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *FreezeTransaction) SignWithOperator(
	client *Client,
) (*FreezeTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return transaction, err
		}
	}
	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *FreezeTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FreezeTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	} else {
		transaction.transactions = make([]*proto.Transaction, 0)
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.signedTransactions); index++ {
		signature := signer(transaction.signedTransactions[index].GetBodyBytes())

		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *FreezeTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	transactionID := transaction.transactionIDs[0]

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		freezeTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{
		TransactionID: transaction.transactionIDs[0],
		NodeID:        resp.transaction.NodeID,
	}, nil
}

func (transaction *FreezeTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_Freeze{
		Freeze: transaction.pb,
	}

	return true
}

func (transaction *FreezeTransaction) Freeze() (*FreezeTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *FreezeTransaction) FreezeWith(client *Client) (*FreezeTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *FreezeTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this FreezeTransaction.
func (transaction *FreezeTransaction) SetMaxTransactionFee(fee Hbar) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *FreezeTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FreezeTransaction.
func (transaction *FreezeTransaction) SetTransactionMemo(memo string) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *FreezeTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FreezeTransaction.
func (transaction *FreezeTransaction) SetTransactionValidDuration(duration time.Duration) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *FreezeTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FreezeTransaction.
func (transaction *FreezeTransaction) SetTransactionID(transactionID TransactionID) *FreezeTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *FreezeTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
}

// SetNodeAccountID sets the node AccountID for this FreezeTransaction.
func (transaction *FreezeTransaction) SetNodeAccountIDs(nodeID []AccountID) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}
