package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
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
		Transaction: _NewTransaction(),
	}

	transaction.SetExpirationTime(time.Now().Add(7890000 * time.Second))
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _FileCreateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) FileCreateTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetFileCreate().GetKeys())
	expiration := _TimeFromProtobuf(pb.GetFileCreate().GetExpirationTime())

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
	transaction._RequireNotFrozen()
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
	transaction._RequireNotFrozen()
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
	transaction._RequireNotFrozen()
	transaction.contents = contents
	return transaction
}

func (transaction *FileCreateTransaction) GetContents() []byte {
	return transaction.contents
}

func (transaction *FileCreateTransaction) SetMemo(memo string) *FileCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo
	return transaction
}

func (transaction *FileCreateTransaction) GetMemo() string {
	return transaction.memo
}

func (transaction *FileCreateTransaction) _Build() *services.TransactionBody {
	body := &services.FileCreateTransactionBody{
		Memo: transaction.memo,
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.keys != nil {
		body.Keys = transaction.keys._ToProtoKeyList()
	}

	if transaction.contents != nil {
		body.Contents = transaction.contents
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_FileCreate{
			FileCreate: body,
		},
	}
}

func (transaction *FileCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *FileCreateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.FileCreateTransactionBody{
		Memo: transaction.memo,
	}

	if transaction.expirationTime != nil {
		body.ExpirationTime = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.keys != nil {
		body.Keys = transaction.keys._ToProtoKeyList()
	}

	if transaction.contents != nil {
		body.Contents = transaction.contents
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_FileCreate{
			FileCreate: body,
		},
	}, nil
}

func _FileCreateTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().CreateFile,
	}
}

func (transaction *FileCreateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
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
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

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
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
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

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := _Execute(
		client,
		_Request{
			transaction: &transaction.Transaction,
		},
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_FileCreateTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.transaction.NodeID,
		}, err
	}

	hash, err := transaction.GetTransactionHash()
	if err != nil {
		return TransactionResponse{}, err
	}

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
	transaction._InitFee(client)
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *FileCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this FileCreateTransaction.
func (transaction *FileCreateTransaction) SetMaxTransactionFee(fee Hbar) *FileCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *FileCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileCreateTransaction.
func (transaction *FileCreateTransaction) SetTransactionMemo(memo string) *FileCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *FileCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileCreateTransaction.
func (transaction *FileCreateTransaction) SetTransactionValidDuration(duration time.Duration) *FileCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *FileCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileCreateTransaction.
func (transaction *FileCreateTransaction) SetTransactionID(transactionID TransactionID) *FileCreateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this FileCreateTransaction.
func (transaction *FileCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *FileCreateTransaction) SetMaxRetry(count int) *FileCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *FileCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileCreateTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
		return transaction
	}

	if len(transaction.signedTransactions) == 0 {
		return transaction
	}

	transaction.transactions = make([]*services.Transaction, 0)
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)

	for index := 0; index < len(transaction.signedTransactions); index++ {
		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

func (transaction *FileCreateTransaction) SetMaxBackoff(max time.Duration) *FileCreateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *FileCreateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *FileCreateTransaction) SetMinBackoff(min time.Duration) *FileCreateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *FileCreateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
