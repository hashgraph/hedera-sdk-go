package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type LiveHashDeleteTransaction struct {
	Transaction
	pb *proto.CryptoDeleteLiveHashTransactionBody
}

func NewLiveHashDeleteTransaction() *LiveHashDeleteTransaction {
	pb := &proto.CryptoDeleteLiveHashTransactionBody{}

	transaction := LiveHashDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	return &transaction
}

func (transaction *LiveHashDeleteTransaction) SetHash(hash []byte) *LiveHashDeleteTransaction {
	transaction.pb.LiveHashToDelete = hash
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetHash() []byte {
	return transaction.pb.GetLiveHashToDelete()
}

func (transaction *LiveHashDeleteTransaction) SetAccountID(accountID AccountID) *LiveHashDeleteTransaction {
	transaction.pb.AccountOfLiveHash = accountID.toProtobuf()
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetAccountID() AccountID {
	return accountIDFromProto(transaction.pb.GetAccountOfLiveHash())
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

	if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
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
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.transactions); index++ {
		signature := signer(transaction.transactions[index].GetBodyBytes())

		transaction.signatures[index].SigPair = append(
			transaction.signatures[index].SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *LiveHashDeleteTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	transactionID := transaction.id

	if !client.GetOperatorID().isZero() && client.GetOperatorID().equals(transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorKey(),
			client.operator.signer,
		)
	}

	_, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeId,
		liveHashDeleteTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
		query_makePaymentTransaction,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{TransactionID: transaction.id}, nil
}

func (transaction *LiveHashDeleteTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_CryptoDeleteLiveHash{
		CryptoDeleteLiveHash: transaction.pb,
	}

	return true
}

func (transaction *LiveHashDeleteTransaction) Freeze() (*LiveHashDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *LiveHashDeleteTransaction) FreezeWith(client *Client) (*LiveHashDeleteTransaction, error) {
	transaction.initFee(client)
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
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetTransactionMemo(memo string) *LiveHashDeleteTransaction {
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *LiveHashDeleteTransaction {
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetTransactionID(transactionID TransactionID) *LiveHashDeleteTransaction {
	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetNodeID() AccountID {
	return transaction.Transaction.GetNodeID()
}

// SetNodeID sets the node AccountID for this LiveHashDeleteTransaction.
func (transaction *LiveHashDeleteTransaction) SetNodeID(nodeID AccountID) *LiveHashDeleteTransaction {
	transaction.Transaction.SetNodeID(nodeID)
	return transaction
}
