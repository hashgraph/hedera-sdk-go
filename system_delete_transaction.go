package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type SystemDeleteTransaction struct {
	Transaction
	contractID     *ContractID
	fileID         *FileID
	expirationTime *time.Time
}

func NewSystemDeleteTransaction() *SystemDeleteTransaction {
	transaction := SystemDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _SystemDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) SystemDeleteTransaction {
	expiration := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute(),
		int(pb.GetSystemDelete().ExpirationTime.Seconds), time.Now().Nanosecond(), time.Now().Location(),
	)
	return SystemDeleteTransaction{
		Transaction:    transaction,
		contractID:     _ContractIDFromProtobuf(pb.GetSystemDelete().GetContractID()),
		fileID:         _FileIDFromProtobuf(pb.GetSystemDelete().GetFileID()),
		expirationTime: &expiration,
	}
}

func (transaction *SystemDeleteTransaction) SetExpirationTime(expiration time.Time) *SystemDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.expirationTime = &expiration
	return transaction
}

func (transaction *SystemDeleteTransaction) GetExpirationTime() int64 {
	if transaction.expirationTime != nil {
		return transaction.expirationTime.Unix()
	}

	return 0
}

func (transaction *SystemDeleteTransaction) SetContractID(contractID ContractID) *SystemDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.contractID = &contractID
	return transaction
}

func (transaction *SystemDeleteTransaction) GetContractID() ContractID {
	if transaction.contractID == nil {
		return ContractID{}
	}

	return *transaction.contractID
}

func (transaction *SystemDeleteTransaction) SetFileID(fileID FileID) *SystemDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.fileID = &fileID
	return transaction
}

func (transaction *SystemDeleteTransaction) GetFileID() FileID {
	if transaction.fileID == nil {
		return FileID{}
	}

	return *transaction.fileID
}

func (transaction *SystemDeleteTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.contractID != nil {
		if err := transaction.contractID.Validate(client); err != nil {
			return err
		}
	}

	if transaction.fileID != nil {
		if err := transaction.fileID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *SystemDeleteTransaction) _Build() *proto.TransactionBody {
	body := &proto.SystemDeleteTransactionBody{}

	if transaction.expirationTime != nil {
		body.ExpirationTime = &proto.TimestampSeconds{
			Seconds: transaction.expirationTime.Unix(),
		}
	}

	if !transaction.contractID._IsZero() {
		body.Id = &proto.SystemDeleteTransactionBody_ContractID{
			ContractID: transaction.contractID._ToProtobuf(),
		}
	}

	if !transaction.fileID._IsZero() {
		body.Id = &proto.SystemDeleteTransactionBody_FileID{
			FileID: transaction.fileID._ToProtobuf(),
		}
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &proto.TransactionBody_SystemDelete{
			SystemDelete: body,
		},
	}
}

func (transaction *SystemDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *SystemDeleteTransaction) _ConstructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.SystemDeleteTransactionBody{}

	if transaction.expirationTime != nil {
		body.ExpirationTime = &proto.TimestampSeconds{
			Seconds: transaction.expirationTime.Unix(),
		}
	}

	if !transaction.contractID._IsZero() {
		body.Id = &proto.SystemDeleteTransactionBody_ContractID{
			ContractID: transaction.contractID._ToProtobuf(),
		}
	}

	if !transaction.fileID._IsZero() {
		body.Id = &proto.SystemDeleteTransactionBody_FileID{
			FileID: transaction.fileID._ToProtobuf(),
		}
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_SystemDelete{
			SystemDelete: body,
		},
	}, nil
}

func _SystemDeleteTransactionGetMethod(request _Request, channel *_Channel) _Method {
	// switch os := runtime.GOOS; os {
	// case "darwin":
	//	fmt.Println("OS X.")
	//}
	if channel._GetContract() == nil {
		return _Method{
			transaction: channel._GetFile().SystemDelete,
		}
	}

	return _Method{
		transaction: channel._GetContract().SystemDelete,
	}
}

func (transaction *SystemDeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *SystemDeleteTransaction) Sign(
	privateKey PrivateKey,
) *SystemDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *SystemDeleteTransaction) SignWithOperator(
	client *Client,
) (*SystemDeleteTransaction, error) {
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
func (transaction *SystemDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *SystemDeleteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *SystemDeleteTransaction) Execute(
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
		_SystemDeleteTransactionGetMethod,
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

func (transaction *SystemDeleteTransaction) Freeze() (*SystemDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *SystemDeleteTransaction) FreezeWith(client *Client) (*SystemDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &SystemDeleteTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *SystemDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetMaxTransactionFee(fee Hbar) *SystemDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *SystemDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetTransactionMemo(memo string) *SystemDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *SystemDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *SystemDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *SystemDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetTransactionID(transactionID TransactionID) *SystemDeleteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *SystemDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *SystemDeleteTransaction) SetMaxRetry(count int) *SystemDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *SystemDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *SystemDeleteTransaction {
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

func (transaction *SystemDeleteTransaction) SetMaxBackoff(max time.Duration) *SystemDeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *SystemDeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *SystemDeleteTransaction) SetMinBackoff(min time.Duration) *SystemDeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *SystemDeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
