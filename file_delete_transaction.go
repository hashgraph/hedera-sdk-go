package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

type FileDeleteTransaction struct {
	Transaction
	fileID *FileID
}

func NewFileDeleteTransaction() *FileDeleteTransaction {
	transaction := FileDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func _FileDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) FileDeleteTransaction {
	return FileDeleteTransaction{
		Transaction: transaction,
		fileID:      _FileIDFromProtobuf(pb.GetFileDelete().GetFileID()),
	}
}

func (transaction *FileDeleteTransaction) SetFileID(fileID FileID) *FileDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.fileID = &fileID
	return transaction
}

func (transaction *FileDeleteTransaction) GetFileID() FileID {
	if transaction.fileID == nil {
		return FileID{}
	}

	return *transaction.fileID
}

func (transaction *FileDeleteTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.fileID != nil {
		if err := transaction.fileID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *FileDeleteTransaction) _Build() *proto.TransactionBody {
	body := &proto.FileDeleteTransactionBody{}
	if transaction.fileID != nil {
		body.FileID = transaction.fileID._ToProtobuf()
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &proto.TransactionBody_FileDelete{
			FileDelete: body,
		},
	}
}

func (transaction *FileDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *FileDeleteTransaction) _ConstructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.FileDeleteTransactionBody{}
	if transaction.fileID != nil {
		body.FileID = transaction.fileID._ToProtobuf()
	}
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_FileDelete{
			FileDelete: body,
		},
	}, nil
}

func _FileDeleteTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetFile().DeleteFile,
	}
}

func (transaction *FileDeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
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
func (transaction *FileDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileDeleteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
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
		_FileDeleteTransactionGetMethod,
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

func (transaction *FileDeleteTransaction) Freeze() (*FileDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *FileDeleteTransaction) FreezeWith(client *Client) (*FileDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}

	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &FileDeleteTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *FileDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this FileDeleteTransaction.
func (transaction *FileDeleteTransaction) SetMaxTransactionFee(fee Hbar) *FileDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *FileDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FileDeleteTransaction.
func (transaction *FileDeleteTransaction) SetTransactionMemo(memo string) *FileDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *FileDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FileDeleteTransaction.
func (transaction *FileDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *FileDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *FileDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FileDeleteTransaction.
func (transaction *FileDeleteTransaction) SetTransactionID(transactionID TransactionID) *FileDeleteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this FileDeleteTransaction.
func (transaction *FileDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *FileDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *FileDeleteTransaction) SetMaxRetry(count int) *FileDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *FileDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *FileDeleteTransaction {
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
