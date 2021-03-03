package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// FileCreateTransaction creates a new file, containing the given contents.  It is referenced by its FileID, and does
// not have a filename, so it is important to get and hold onto the FileID. After the file is created, the FileID for
// it can be found in the receipt, or retrieved with a GetByKey query, or by asking for a Record of the transaction to
// be created, and retrieving that.
//
// See FileInfoQuery for more information about files.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0, with
// a null key. Future versions of the API will support multiple realms and multiple shards.
type FileCreateTransaction struct {
	Transaction
	pb *proto.FileCreateTransactionBody
}

// NewFileCreateTransaction creates a FileCreateTransaction transaction which can be
// used to construct and execute a File Create Transaction.
func NewFileCreateTransaction() *FileCreateTransaction {
	pb := &proto.FileCreateTransactionBody{}

	transaction := FileCreateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetExpirationTime(time.Now().Add(7890000 * time.Second))
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func fileCreateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) FileCreateTransaction {
	return FileCreateTransaction{
		Transaction: transaction,
		pb:          pb.GetFileCreate(),
	}
}

// AddKey adds a key to the internal list of keys associated with the file. All of the keys on the list must sign to
// create or modify a file, but only one of them needs to sign in order to delete the file. Each of those "keys" may
// itself be threshold key containing other keys (including other threshold keys). In other words, the behavior is an
// AND for create/modify, OR for delete. This is useful for acting as a revocation server. If it is desired to have the
// behavior be AND for all 3 operations (or OR for all 3), then the list should have only a single Key, which is a
// threshold key, with N=1 for OR, N=M for AND.
//
// If a file is created without adding ANY keys, the file is immutable and ONLY the
// expirationTime of the file can be changed using FileUpdateTransaction. The file contents or its keys will not be
// mutable.
func (transaction *FileCreateTransaction) SetKeys(keys ...Key) *FileCreateTransaction {
	transaction.requireNotFrozen()
	if transaction.pb.Keys == nil {
		transaction.pb.Keys = &proto.KeyList{Keys: []*proto.Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	transaction.pb.Keys = keyList.toProtoKeyList()

	return transaction
}

func (transaction *FileCreateTransaction) GetKeys() KeyList {
	keys := transaction.pb.GetKeys()
	if keys != nil {
		keyList, err := keyListFromProtobuf(keys)
		if err != nil {
			return KeyList{}
		}

		return keyList
	} else {
		return KeyList{}
	}
}

// SetExpirationTime sets the time at which this file should expire (unless FileUpdateTransaction is used before then to
// extend its life). The file will automatically disappear at the fileExpirationTime, unless its expiration is extended
// by another transaction before that time. If the file is deleted, then its contents will become empty and it will be
// marked as deleted until it expires, and then it will cease to exist.
func (transaction *FileCreateTransaction) SetExpirationTime(expiration time.Time) *FileCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ExpirationTime = timeToProtobuf(expiration)
	return transaction
}

func (transaction *FileCreateTransaction) GetExpirationTime() time.Time {
	return timeFromProtobuf(transaction.pb.GetExpirationTime())
}

// SetContents sets the bytes that are the contents of the file (which can be empty). If the size of the file and other
// fields in the transaction exceed the max transaction size then FileAppendTransaction can be used to continue
// uploading the file.
func (transaction *FileCreateTransaction) SetContents(contents []byte) *FileCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Contents = contents
	return transaction
}

func (transaction *FileCreateTransaction) GetContents() []byte {
	return transaction.pb.GetContents()
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func fileCreateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getFile().CreateFile,
	}
}

func (transaction *FileCreateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *FileCreateTransaction) Sign(
	privateKey PrivateKey,
) *FileCreateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *FileCreateTransaction) SignWithOperator(
	client *Client,
) (*FileCreateTransaction, error) {
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
func (transaction *FileCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileCreateTransaction {
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
func (transaction *FileCreateTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil || client.operator == nil {
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
		fileCreateTransaction_getMethod,
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

func (transaction *FileCreateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_FileCreate{
		FileCreate: transaction.pb,
	}

	return true
}

func (transaction *FileCreateTransaction) Freeze() (*FileCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *FileCreateTransaction) FreezeWith(client *Client) (*FileCreateTransaction, error) {
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

func (transaction *FileCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this FileCreateTransaction.
func (transaction *FileCreateTransaction) SetMaxTransactionFee(fee Hbar) *FileCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *FileCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileCreateTransaction.
func (transaction *FileCreateTransaction) SetTransactionMemo(memo string) *FileCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *FileCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileCreateTransaction.
func (transaction *FileCreateTransaction) SetTransactionValidDuration(duration time.Duration) *FileCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *FileCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileCreateTransaction.
func (transaction *FileCreateTransaction) SetTransactionID(transactionID TransactionID) *FileCreateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this FileCreateTransaction.
func (transaction *FileCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *FileCreateTransaction) SetMaxRetry(count int) *FileCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *FileCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileCreateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
