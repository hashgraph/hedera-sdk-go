package hedera

import (
	"errors"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type ScheduleCreateTransaction struct {
	Transaction
	payerAccountID  AccountID
	adminKey        Key
	schedulableBody *proto.SchedulableTransactionBody
	memo            string
}

func NewScheduleCreateTransaction() *ScheduleCreateTransaction {
	transaction := ScheduleCreateTransaction{
		Transaction: newTransaction(),
	}

	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func scheduleCreateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ScheduleCreateTransaction {
	key, _ := keyFromProtobuf(pb.GetScheduleCreate().GetAdminKey())

	return ScheduleCreateTransaction{
		Transaction:     transaction,
		payerAccountID:  accountIDFromProtobuf(pb.GetScheduleCreate().GetPayerAccountID()),
		adminKey:        key,
		schedulableBody: pb.GetScheduleCreate().GetScheduledTransactionBody(),
		memo:            pb.GetScheduleCreate().GetMemo(),
	}
}

func (transaction *ScheduleCreateTransaction) SetPayerAccountID(id AccountID) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()
	transaction.payerAccountID = id

	return transaction
}

func (transaction *ScheduleCreateTransaction) GetPayerAccountID() AccountID {
	return transaction.payerAccountID
}

func (transaction *ScheduleCreateTransaction) SetAdminKey(key Key) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()
	transaction.adminKey = key

	return transaction
}

func (transaction *ScheduleCreateTransaction) setSchedulableTransactionBody(txBody *proto.SchedulableTransactionBody) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()
	transaction.schedulableBody = txBody

	return transaction
}

func (transaction *ScheduleCreateTransaction) GetAdminKey() *Key {
	if transaction.adminKey == nil {
		return nil
	}
	return &transaction.adminKey
}

func (transaction *ScheduleCreateTransaction) SetScheduleMemo(memo string) *ScheduleCreateTransaction {
	transaction.requireNotFrozen()
	transaction.memo = memo

	return transaction
}

func (transaction *ScheduleCreateTransaction) GetScheduleMemo() string {
	return transaction.memo
}

func (transaction *ScheduleCreateTransaction) SetScheduledTransaction(tx ITransaction) (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := tx.constructScheduleProtobuf()
	if err != nil {
		return transaction, err
	}

	transaction.schedulableBody = scheduled
	return transaction, nil
}

func (transaction *ScheduleCreateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = transaction.payerAccountID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *ScheduleCreateTransaction) build() *proto.TransactionBody {
	body := &proto.ScheduleCreateTransactionBody{
		Memo: transaction.memo,
	}

	if !transaction.payerAccountID.isZero() {
		body.PayerAccountID = transaction.payerAccountID.toProtobuf()
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey.toProtoKey()
	}

	if transaction.schedulableBody != nil {
		body.ScheduledTransactionBody = transaction.schedulableBody
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_ScheduleCreate{
			ScheduleCreate: body,
		},
	}
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
		transaction_makeRequest(request{
			transaction: &transaction.Transaction,
		}),
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

func (transaction *ScheduleCreateTransaction) Freeze() (*ScheduleCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ScheduleCreateTransaction) FreezeWith(client *Client) (*ScheduleCreateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &ScheduleCreateTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	//transaction.transactionIDs[0] = transaction.transactionIDs[0].SetScheduled(true)

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
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
