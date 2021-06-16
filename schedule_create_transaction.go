package hedera

import (
	"errors"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type ScheduleCreateTransaction struct {
	Transaction
	pb *proto.ScheduleCreateTransactionBody
}

func NewScheduleCreateTransaction() *ScheduleCreateTransaction {
	pb := &proto.ScheduleCreateTransactionBody{}

	transaction := ScheduleCreateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func scheduleCreateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ScheduleCreateTransaction {
	return ScheduleCreateTransaction{
		Transaction: transaction,
		pb:          pb.GetScheduleCreate(),
	}
}

func (transaction *ScheduleCreateTransaction) SetPayerAccountID(id AccountID) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.PayerAccountID = id.toProtobuf()

	return transaction
}

func (transaction *ScheduleCreateTransaction) GetPayerAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.PayerAccountID)
}

func (transaction *ScheduleCreateTransaction) SetAdminKey(key Key) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AdminKey = key.toProtoKey()

	return transaction
}

func (transaction *ScheduleCreateTransaction) setSchedulableTransactionBody(txBody *proto.SchedulableTransactionBody) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ScheduledTransactionBody = txBody

	return transaction
}

func (transaction *ScheduleCreateTransaction) GetAdminKey() *Key {
	key, err := keyFromProtobuf(transaction.pb.GetAdminKey())
	if err != nil {
		return nil
	}
	return &key
}

func (transaction *ScheduleCreateTransaction) SetScheduleMemo(memo string) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Memo = memo

	return transaction
}

func (transaction *ScheduleCreateTransaction) GetScheduleMemo() string {
	return transaction.pb.GetMemo()
}

func (transaction *ScheduleCreateTransaction) SetScheduledTransaction(tx ITransaction) (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := tx.constructScheduleProtobuf()
	if err != nil {
		return transaction, err
	}

	transaction.pb.ScheduledTransactionBody = scheduled
	return transaction, nil
}

func (transaction *ScheduleCreateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `ScheduleCreateTransaction`")
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func scheduleCreateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getSchedule().CreateSchedule,
	}
}

func (transaction *ScheduleCreateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ScheduleCreateTransaction) Sign(
	privateKey PrivateKey,
) *ScheduleCreateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *ScheduleCreateTransaction) SignWithOperator(
	client *Client,
) (*ScheduleCreateTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
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
func (transaction *ScheduleCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ScheduleCreateTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ScheduleCreateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	transactionID := transaction.GetTransactionID()

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(*transactionID.AccountID) {
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
		scheduleCreateTransaction_getMethod,
		transaction_mapStatusError,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.transaction.NodeID,
		}, err
	}

	hash, err := transaction.GetTransactionHash()

	return TransactionResponse{
		TransactionID:          transaction.GetTransactionID(),
		NodeID:                 resp.transaction.NodeID,
		Hash:                   hash,
		ScheduledTransactionId: transaction.GetTransactionID(),
	}, nil
}

func (transaction *ScheduleCreateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_ScheduleCreate{
		ScheduleCreate: transaction.pb,
	}

	return true
}

func (transaction *ScheduleCreateTransaction) Freeze() (*ScheduleCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ScheduleCreateTransaction) FreezeWith(client *Client) (*ScheduleCreateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	//transaction.transactionIDs[0] = transaction.transactionIDs[0].SetScheduled(true)

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *ScheduleCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this ScheduleCreateTransaction.
func (transaction *ScheduleCreateTransaction) SetMaxTransactionFee(fee Hbar) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ScheduleCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ScheduleCreateTransaction.
func (transaction *ScheduleCreateTransaction) SetTransactionMemo(memo string) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ScheduleCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ScheduleCreateTransaction.
func (transaction *ScheduleCreateTransaction) SetTransactionValidDuration(duration time.Duration) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ScheduleCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ScheduleCreateTransaction.
func (transaction *ScheduleCreateTransaction) SetTransactionID(transactionID TransactionID) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this ScheduleCreateTransaction.
func (transaction *ScheduleCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *ScheduleCreateTransaction) SetMaxRetry(count int) *ScheduleCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}
