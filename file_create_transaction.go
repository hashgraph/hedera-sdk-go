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
	keys           *KeyList
	expirationTime *time.Time
	contents       []byte
	memo           string
}

// NewFileCreateTransaction creates a FileCreateTransaction transaction which can be
// used to construct and execute a File Create Transaction.
func NewFileCreateTransaction() *FileCreateTransaction {
	transaction := FileCreateTransaction{
		Transaction: newTransaction(),
	}

	transaction.SetExpirationTime(time.Now().Add(7890000 * time.Second))
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func fileCreateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) FileCreateTransaction {
	keys, _ := keyListFromProtobuf(pb.GetFileCreate().GetKeys())
	expiration := timeFromProtobuf(pb.GetFileCreate().GetExpirationTime())

	return FileCreateTransaction{
		Transaction:    transaction,
		keys:           &keys,
		expirationTime: &expiration,
		contents:       pb.GetFileCreate().GetContents(),
		memo:           pb.GetMemo(),
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
	if transaction.keys == nil {
		transaction.keys = &KeyList{keys: []Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	transaction.keys = keyList

	return transaction
}

func (transaction *FileCreateTransaction) GetKeys() KeyList {
	if transaction.keys != nil {
		return *transaction.keys
	}

	return KeyList{}
}

// SetExpirationTime sets the time at which this file should expire (unless FileUpdateTransaction is used before then to
// extend its life). The file will automatically disappear at the fileExpirationTime, unless its expiration is extended
// by another transaction before that time. If the file is deleted, then its contents will become empty and it will be
// marked as deleted until it expires, and then it will cease to exist.
func (transaction *FileCreateTransaction) SetExpirationTime(expiration time.Time) *FileCreateTransaction {
	transaction.requireNotFrozen()
	transaction.expirationTime = &expiration
	return transaction
}

func (transaction *FileCreateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}

	return time.Time{}
}

// SetContents sets the bytes that are the contents of the file (which can be empty). If the size of the file and other
// fields in the transaction exceed the max transaction size then FileAppendTransaction can be used to continue
// uploading the file.
func (transaction *FileCreateTransaction) SetContents(contents []byte) *FileCreateTransaction {
	transaction.requireNotFrozen()
	transaction.contents = contents
	return transaction
}

func (transaction *FileCreateTransaction) GetContents() []byte {
	return transaction.contents
}

func (transaction *FileCreateTransaction) SetMemo(memo string) *FileCreateTransaction {
	transaction.requireNotFrozen()
	transaction.memo = memo
	return transaction
}

func (transaction *FileCreateTransaction) GetMemo() string {
	return transaction.memo
}

func (transaction *FileCreateTransaction) build() *proto.TransactionBody {
	body := &proto.FileCreateTransactionBody{
		Memo: transaction.memo,
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = timeToProtobuf(*transaction.expirationTime)
	}

	if transaction.keys != nil {
		body.Keys = transaction.keys.toProtoKeyList()
	}

	if transaction.contents != nil {
		body.Contents = transaction.contents
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_FileCreate{
			FileCreate: body,
		},
	}
}

func (transaction *FileCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *FileCreateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.FileCreateTransactionBody{
		Memo: transaction.memo,
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = timeToProtobuf(*transaction.expirationTime)
	}

	if transaction.keys != nil {
		body.Keys = transaction.keys.toProtoKeyList()
	}

	if transaction.contents != nil {
		body.Contents = transaction.contents
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_FileCreate{
			FileCreate: body,
		},
	}, nil
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
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *FileCreateTransaction) Execute(
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
		fileCreateTransaction_getMethod,
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
	body := transaction.build()

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
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
	transaction.requireOneNodeAccountID()

	if !transaction.isFrozen() {
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	if len(transaction.signedTransactions) == 0 {
		return transaction
	}

	transaction.transactions = make([]*proto.Transaction, 0)
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)

	for index := 0; index < len(transaction.signedTransactions); index++ {
		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	//transaction.signedTransactions[0].SigMap.SigPair = append(transaction.signedTransactions[0].SigMap.SigPair, publicKey.toSignaturePairProtobuf(signature))
	return transaction
}
