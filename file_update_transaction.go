package hedera

import (
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type FileUpdateTransaction struct {
	Transaction
	fileID         *FileID
	keys           *KeyList
	expirationTime *time.Time
	contents       []byte
	memo           string
}

func NewFileUpdateTransaction() *FileUpdateTransaction {
	transaction := FileUpdateTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _FileUpdateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) FileUpdateTransaction {
	keys, _ := _KeyListFromProtobuf(pb.GetFileUpdate().GetKeys())
	expiration := _TimeFromProtobuf(pb.GetFileUpdate().GetExpirationTime())

	return FileUpdateTransaction{
		Transaction:    transaction,
		fileID:         _FileIDFromProtobuf(pb.GetFileUpdate().GetFileID()),
		keys:           &keys,
		expirationTime: &expiration,
		contents:       pb.GetFileUpdate().GetContents(),
		memo:           pb.GetFileUpdate().GetMemo().Value,
	}
}

func (transaction *FileUpdateTransaction) SetFileID(fileID FileID) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.fileID = &fileID
	return transaction
}

func (transaction *FileUpdateTransaction) GetFileID() FileID {
	if transaction.fileID == nil {
		return FileID{}
	}

	return *transaction.fileID
}

func (transaction *FileUpdateTransaction) SetKeys(keys ...Key) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	if transaction.keys == nil {
		transaction.keys = &KeyList{keys: []Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	transaction.keys = keyList

	return transaction
}

func (transaction *FileUpdateTransaction) GetKeys() KeyList {
	if transaction.keys != nil {
		return *transaction.keys
	}

	return KeyList{}
}

func (transaction *FileUpdateTransaction) SetExpirationTime(expiration time.Time) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.expirationTime = &expiration
	return transaction
}

func (transaction *FileUpdateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}

	return time.Time{}
}

func (transaction *FileUpdateTransaction) SetContents(contents []byte) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.contents = contents
	return transaction
}

func (transaction *FileUpdateTransaction) GetContents() []byte {
	return transaction.contents
}

func (transaction *FileUpdateTransaction) SetFileMemo(memo string) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo

	return transaction
}

func (transaction *FileUpdateTransaction) GeFileMemo() string {
	return transaction.memo
}

func (transaction *FileUpdateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil {
		return nil
	}

	if transaction.fileID != nil {
		if err := transaction.fileID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *FileUpdateTransaction) _Build() *proto.TransactionBody {
	body := &proto.FileUpdateTransactionBody{
		Memo: &wrappers.StringValue{Value: transaction.memo},
	}
	if transaction.fileID != nil {
		body.FileID = transaction.fileID._ToProtobuf()
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

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &proto.TransactionBody_FileUpdate{
			FileUpdate: body,
		},
	}
}

func (transaction *FileUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *FileUpdateTransaction) _ConstructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.FileUpdateTransactionBody{
		Memo: &wrappers.StringValue{Value: transaction.memo},
	}
	if transaction.fileID != nil {
		body.FileID = transaction.fileID._ToProtobuf()
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

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_FileUpdate{
			FileUpdate: body,
		},
	}, nil
}

func _FileUpdateTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().UpdateFile,
	}
}

func (transaction *FileUpdateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
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
func (transaction *FileUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileUpdateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *FileUpdateTransaction) Execute(
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
		_TransactionMakeRequest(_Request{
			transaction: &transaction.Transaction,
		}),
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_FileUpdateTransactionGetMethod,
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

func (transaction *FileUpdateTransaction) Freeze() (*FileUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *FileUpdateTransaction) FreezeWith(client *Client) (*FileUpdateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &FileUpdateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *FileUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetMaxTransactionFee(fee Hbar) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *FileUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetTransactionMemo(memo string) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *FileUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *FileUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetTransactionID(transactionID TransactionID) *FileUpdateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *FileUpdateTransaction) SetMaxRetry(count int) *FileUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *FileUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileUpdateTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
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
			publicKey._ToSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

func (transaction *FileUpdateTransaction) SetMaxBackoff(max time.Duration) *FileUpdateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *FileUpdateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *FileUpdateTransaction) SetMinBackoff(min time.Duration) *FileUpdateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *FileUpdateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
