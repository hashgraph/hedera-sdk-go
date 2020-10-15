package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type FileUpdateTransaction struct {
	Transaction
	pb *proto.FileUpdateTransactionBody
}

func NewFileUpdateTransaction() *FileUpdateTransaction {
	pb := &proto.FileUpdateTransactionBody{}

	transaction := FileUpdateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func (transaction *FileUpdateTransaction) SetFileID(ID FileID) *FileUpdateTransaction {
	transaction.pb.FileID = ID.toProto()
	return transaction
}

func (transaction *FileUpdateTransaction) GetFileID() FileID {
	return fileIDFromProto(transaction.pb.GetFileID())
}

func (transaction *FileUpdateTransaction) SetKeys(keys ...Key) *FileUpdateTransaction {
	if transaction.pb.Keys == nil {
		transaction.pb.Keys = &proto.KeyList{Keys: []*proto.Key{}}
	}
	keyList := KeyList{keys: []*proto.Key{}}
	keyList.AddAll(keys)

	transaction.pb.Keys = keyList.toProtoKeyList()

	return transaction
}

func (transaction *FileUpdateTransaction) GetKeys() KeyList {
	return keyListFromProto(transaction.pb.Keys)
}

func (transaction *FileUpdateTransaction) SetExpirationTime(expiration time.Time) *FileUpdateTransaction {
	transaction.pb.ExpirationTime = timeToProto(expiration)
	return transaction
}

func (transaction *FileUpdateTransaction) GetExpirationTime() time.Time {
	return timeFromProto(transaction.pb.ExpirationTime)
}

func (transaction *FileUpdateTransaction) SetContents(contents []byte) *FileUpdateTransaction {
	transaction.pb.Contents = contents
	return transaction
}

func (transaction *FileUpdateTransaction) GetContents() []byte {
	return transaction.pb.Contents
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func fileUpdateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getFile().UpdateFile,
	}
}

func (transaction *FileUpdateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *FileUpdateTransaction) Sign(
	privateKey PrivateKey,
) *FileUpdateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *FileUpdateTransaction) SignWithOperator(
	client *Client,
) (*FileUpdateTransaction, error) {
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
func (transaction *FileUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileUpdateTransaction {
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
func (transaction *FileUpdateTransaction) Execute(
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
		fileUpdateTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{TransactionID: transaction.id}, nil
}

func (transaction *FileUpdateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_FileUpdate{
		FileUpdate: transaction.pb,
	}

	return true
}

func (transaction *FileUpdateTransaction) Freeze() (*FileUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *FileUpdateTransaction) FreezeWith(client *Client) (*FileUpdateTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *FileUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetMaxTransactionFee(fee Hbar) *FileUpdateTransaction {
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *FileUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetTransactionMemo(memo string) *FileUpdateTransaction {
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *FileUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *FileUpdateTransaction {
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *FileUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetTransactionID(transactionID TransactionID) *FileUpdateTransaction {
	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *FileUpdateTransaction) GetNodeID() AccountID {
	return transaction.Transaction.GetNodeID()
}

// SetNodeID sets the node AccountID for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetNodeID(nodeID AccountID) *FileUpdateTransaction {
	transaction.Transaction.SetNodeID(nodeID)
	return transaction
}
