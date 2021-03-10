package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

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
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func liveHashDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) LiveHashDeleteTransaction {
	return LiveHashDeleteTransaction{
		Transaction: transaction,
		pb:          pb.GetCryptoDeleteLiveHash(),
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

func (transaction *LiveHashDeleteTransaction) SetAccountID(accountID AccountID) *LiveHashDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AccountOfLiveHash = accountID.toProtobuf()
	return transaction
}

func (transaction *LiveHashDeleteTransaction) GetAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetAccountOfLiveHash())
}

func (transaction *LiveHashDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	txBytes, err := protobuf.Marshal(transaction.constructProtobuf())
	if err != nil {
		return &ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setTransactionBodyBytes(txBytes), nil
}

func (transaction *LiveHashDeleteTransaction) constructProtobuf() *proto.TransactionBody{
	return &proto.TransactionBody{
		TransactionID:            transaction.pbBody.GetTransactionID(),
		NodeAccountID:            transaction.pbBody.GetNodeAccountID(),
		TransactionFee:           transaction.pbBody.GetTransactionFee(),
		TransactionValidDuration: transaction.pbBody.GetTransactionValidDuration(),
		GenerateRecord:           transaction.pbBody.GetGenerateRecord(),
		Memo:                     transaction.pbBody.GetMemo(),
		Data: &proto.TransactionBody_CryptoDeleteLiveHash{
			CryptoDeleteLiveHash: &proto.CryptoDeleteLiveHashTransactionBody{
				AccountOfLiveHash: transaction.pb.GetAccountOfLiveHash(),
				LiveHashToDelete:  transaction.pb.GetLiveHashToDelete(),
			},
		},
	}
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
func (transaction *LiveHashDeleteTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil || client.operator == nil {
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
		transaction_mapResponseStatus,
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
	if transaction.IsFrozen() {
		return transaction, nil
	}
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
