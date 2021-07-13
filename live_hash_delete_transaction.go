package hedera

import (
	"errors"
	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

type LiveHashDeleteTransaction struct {
	Transaction
	pb        *services.CryptoDeleteLiveHashTransactionBody
	accountID AccountID
}

func NewLiveHashDeleteTransaction() *LiveHashDeleteTransaction {
	pb := &services.CryptoDeleteLiveHashTransactionBody{}

	transaction := LiveHashDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func liveHashDeleteTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) LiveHashDeleteTransaction {
	return LiveHashDeleteTransaction{
		Transaction: transaction,
		pb:          pb.GetCryptoDeleteLiveHash(),
		accountID:   accountIDFromProtobuf(pb.GetCryptoDeleteLiveHash().GetAccountOfLiveHash(), nil),
	}
}

func (transaction *LiveHashDeleteTransaction) SetHash(hash []byte) *LiveHashDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.LiveHashToDelete = hash
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetHash() []byte {
	return transaction.pb.GetLiveHashToDelete()
}

func (transaction *LiveHashDeleteTransaction) SetAccountID(id AccountID) *LiveHashDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.accountID = id
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetAccountID() AccountID {
	return transaction.accountID
}

func (transaction *LiveHashDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	err = transaction.accountID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *LiveHashDeleteTransaction) build() *LiveHashDeleteTransaction {
	if !transaction.accountID.isZero() {
		transaction.pb.AccountOfLiveHash = transaction.accountID.toProtobuf()
	}

	return transaction
}

func (transaction *LiveHashDeleteTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `LiveHashAddTransaction`")
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func liveHashDeleteTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getCrypto().DeleteLiveHash,
	}
}

func (transaction *LiveHashDeleteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *LiveHashDeleteTransaction) Sign(
	privateKey PrivateKey,
) *LiveHashDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *LiveHashDeleteTransaction) SignWithOperator(
	client *Client,
) (*LiveHashDeleteTransaction, error) {
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
func (transaction *LiveHashDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *LiveHashDeleteTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *LiveHashDeleteTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
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
		liveHashDeleteTransaction_getMethod,
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

func (transaction *LiveHashDeleteTransaction) onFreeze(
	pbBody *services.TransactionBody,
) bool {
	pbBody.Data = &services.TransactionBody_CryptoDeleteLiveHash{
		CryptoDeleteLiveHash: transaction.pb,
	}

	return true
}

func (transaction *LiveHashDeleteTransaction) Freeze() (*LiveHashDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *LiveHashDeleteTransaction) FreezeWith(client *Client) (*LiveHashDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &LiveHashDeleteTransaction{}, err
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

func (transaction *LiveHashDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetMaxTransactionFee(fee Hbar) *LiveHashDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetTransactionMemo(memo string) *LiveHashDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *LiveHashDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetTransactionID(transactionID TransactionID) *LiveHashDeleteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *LiveHashDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) SetMaxRetry(count int) *LiveHashDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *LiveHashDeleteTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
