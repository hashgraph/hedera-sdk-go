package hedera

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type FileUpdateTransaction struct {
	Transaction
	pb     *proto.FileUpdateTransactionBody
	fileID FileID
}

func NewFileUpdateTransaction() *FileUpdateTransaction {
	pb := &proto.FileUpdateTransactionBody{}

	transaction := FileUpdateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func fileUpdateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) FileUpdateTransaction {
	return FileUpdateTransaction{
		Transaction: transaction,
		pb:          pb.GetFileUpdate(),
		fileID:      fileIDFromProtobuf(pb.GetFileUpdate().GetFileID(), nil),
	}
}

func (transaction *FileUpdateTransaction) SetFileID(id FileID) *FileUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.fileID = id
	return transaction
}

func (transaction *FileUpdateTransaction) GetFileID() FileID {
	return transaction.fileID
}

func (transaction *FileUpdateTransaction) SetKeys(keys ...Key) *FileUpdateTransaction {
	transaction.requireNotFrozen()
	if transaction.pb.Keys == nil {
		transaction.pb.Keys = &proto.KeyList{Keys: []*proto.Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	transaction.pb.Keys = keyList.toProtoKeyList()

	return transaction
}

func (transaction *FileUpdateTransaction) GetKeys() KeyList {
	keys := transaction.pb.GetKeys()
	if keys != nil {
		keyList, err := keyListFromProtobuf(keys, nil)
		if err != nil {
			return KeyList{}
		}

		return keyList
	} else {
		return KeyList{}
	}
}

func (transaction *FileUpdateTransaction) SetExpirationTime(expiration time.Time) *FileUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ExpirationTime = timeToProtobuf(expiration)
	return transaction
}

func (transaction *FileUpdateTransaction) GetExpirationTime() time.Time {
	return timeFromProtobuf(transaction.pb.ExpirationTime)
}

func (transaction *FileUpdateTransaction) SetContents(contents []byte) *FileUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Contents = contents
	return transaction
}

func (transaction *FileUpdateTransaction) GetContents() []byte {
	return transaction.pb.Contents
}

func (transaction *FileUpdateTransaction) SetFileMemo(memo string) *FileUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Memo = &wrappers.StringValue{Value: memo}

	return transaction
}

func (transaction *FileUpdateTransaction) GeFileMemo() string {
	if transaction.pb.Memo != nil {
		return transaction.pb.Memo.GetValue()
	}

	return ""
}

func (transaction *FileUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	err = transaction.fileID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *FileUpdateTransaction) build() *FileUpdateTransaction {
	if !transaction.fileID.isZero() {
		transaction.pb.FileID = transaction.fileID.toProtobuf()
	}

	return transaction
}

func (transaction *FileUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *FileUpdateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	transaction.build()
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_FileUpdate{
			FileUpdate: &proto.FileUpdateTransactionBody{
				FileID:         transaction.pb.GetFileID(),
				ExpirationTime: transaction.pb.GetExpirationTime(),
				Keys:           transaction.pb.GetKeys(),
				Contents:       transaction.pb.GetContents(),
				Memo:           transaction.pb.GetMemo(),
			},
		},
	}, nil
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
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
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
		fileUpdateTransaction_getMethod,
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
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &FileUpdateTransaction{}, err
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

func (transaction *FileUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetMaxTransactionFee(fee Hbar) *FileUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *FileUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetTransactionMemo(memo string) *FileUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *FileUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *FileUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *FileUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetTransactionID(transactionID TransactionID) *FileUpdateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this FileUpdateTransaction.
func (transaction *FileUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *FileUpdateTransaction) SetMaxRetry(count int) *FileUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *FileUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileUpdateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
