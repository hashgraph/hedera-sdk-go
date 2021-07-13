package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/pkg/errors"
	"time"
)

type ScheduleSignTransaction struct {
	Transaction
	pb         *services.ScheduleSignTransactionBody
	scheduleID ScheduleID
}

func NewScheduleSignTransaction() *ScheduleSignTransaction {
	pb := &services.ScheduleSignTransactionBody{}

	transaction := ScheduleSignTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func scheduleSignTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) ScheduleSignTransaction {
	return ScheduleSignTransaction{
		Transaction: transaction,
		pb:          pb.GetScheduleSign(),
		scheduleID:  scheduleIDFromProtobuf(pb.GetScheduleSign().GetScheduleID(), nil),
	}
}

func (transaction *ScheduleSignTransaction) SetScheduleID(id ScheduleID) *ScheduleSignTransaction {
	transaction.requireNotFrozen()
	transaction.scheduleID = id

	return transaction
}

func (transaction *ScheduleSignTransaction) GetScheduleID() ScheduleID {
	return transaction.scheduleID
}

func (transaction *ScheduleSignTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	err = transaction.scheduleID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *ScheduleSignTransaction) build() *ScheduleSignTransaction {
	if !transaction.scheduleID.isZero() {
		transaction.pb.ScheduleID = transaction.scheduleID.toProtobuf()
	}

	return transaction
}

func (transaction *ScheduleSignTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `ScheduleSignTransaction")
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func scheduleSignTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getSchedule().SignSchedule,
	}
}

func (transaction *ScheduleSignTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ScheduleSignTransaction) Sign(
	privateKey PrivateKey,
) *ScheduleSignTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *ScheduleSignTransaction) SignWithOperator(
	client *Client,
) (*ScheduleSignTransaction, error) {
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
func (transaction *ScheduleSignTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ScheduleSignTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ScheduleSignTransaction) Execute(
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
		scheduleSignTransaction_getMethod,
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
		TransactionID: transaction.GetTransactionID(),
		NodeID:        resp.transaction.NodeID,
		Hash:          hash,
	}, nil
}

func (transaction *ScheduleSignTransaction) onFreeze(
	pbBody *services.TransactionBody,
) bool {
	pbBody.Data = &services.TransactionBody_ScheduleSign{
		ScheduleSign: transaction.pb,
	}

	return true
}

func (transaction *ScheduleSignTransaction) Freeze() (*ScheduleSignTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ScheduleSignTransaction) FreezeWith(client *Client) (*ScheduleSignTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &ScheduleSignTransaction{}, err
	}
	transaction.build()

	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *ScheduleSignTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this ScheduleSignTransaction.
func (transaction *ScheduleSignTransaction) SetMaxTransactionFee(fee Hbar) *ScheduleSignTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ScheduleSignTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ScheduleSignTransaction.
func (transaction *ScheduleSignTransaction) SetTransactionMemo(memo string) *ScheduleSignTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ScheduleSignTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ScheduleSignTransaction.
func (transaction *ScheduleSignTransaction) SetTransactionValidDuration(duration time.Duration) *ScheduleSignTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ScheduleSignTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ScheduleSignTransaction.
func (transaction *ScheduleSignTransaction) SetTransactionID(transactionID TransactionID) *ScheduleSignTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this ScheduleSignTransaction.
func (transaction *ScheduleSignTransaction) SetNodeAccountIDs(nodeID []AccountID) *ScheduleSignTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *ScheduleSignTransaction) SetMaxRetry(count int) *ScheduleSignTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}
