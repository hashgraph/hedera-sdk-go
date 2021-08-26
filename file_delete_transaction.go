package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

type FileDeleteTransaction struct {
	Transaction
	fileID FileID
}

func NewFileDeleteTransaction() *FileDeleteTransaction {
	transaction := FileDeleteTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func fileDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) FileDeleteTransaction {
	return FileDeleteTransaction{
		Transaction: transaction,
		fileID:      fileIDFromProtobuf(pb.GetFileDelete().GetFileID()),
	}
}

func (transaction *FileDeleteTransaction) SetFileID(id FileID) *FileDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.fileID = id
	return transaction
}

func (transaction *FileDeleteTransaction) GetFileID() FileID {
	return transaction.fileID
}

func (transaction *FileDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = transaction.fileID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *FileDeleteTransaction) build() *proto.TransactionBody {
	body := &proto.FileDeleteTransactionBody{}
	if !transaction.fileID.isZero() {
		body.FileID = transaction.fileID.toProtobuf()
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_FileDelete{
			FileDelete: body,
		},
	}
}

func (transaction *FileDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *FileDeleteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.FileDeleteTransactionBody{}
	if !transaction.fileID.isZero() {
		body.FileID = transaction.fileID.toProtobuf()
	}
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_FileDelete{
			FileDelete: body,
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func fileDeleteTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getFile().DeleteFile,
	}
}

func (transaction *FileDeleteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *FileDeleteTransaction) Sign(
	privateKey PrivateKey,
) *FileDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *FileDeleteTransaction) SignWithOperator(
	client *Client,
) (*FileDeleteTransaction, error) {
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
func (transaction *FileDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileDeleteTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *FileDeleteTransaction) Execute(
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
		fileDeleteTransaction_getMethod,
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

func (transaction *FileDeleteTransaction) Freeze() (*FileDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *FileDeleteTransaction) FreezeWith(client *Client) (*FileDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}

	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &FileDeleteTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
}

func (transaction *FileDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this FileDeleteTransaction.
func (transaction *FileDeleteTransaction) SetMaxTransactionFee(fee Hbar) *FileDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *FileDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileDeleteTransaction.
func (transaction *FileDeleteTransaction) SetTransactionMemo(memo string) *FileDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *FileDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileDeleteTransaction.
func (transaction *FileDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *FileDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *FileDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileDeleteTransaction.
func (transaction *FileDeleteTransaction) SetTransactionID(transactionID TransactionID) *FileDeleteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this FileDeleteTransaction.
func (transaction *FileDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *FileDeleteTransaction) SetMaxRetry(count int) *FileDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *FileDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileDeleteTransaction {
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

func (transaction *FileDeleteTransaction) SetMaxBackoff(max time.Duration) *FileDeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *FileDeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *FileDeleteTransaction) SetMinBackoff(min time.Duration) *FileDeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *FileDeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
