package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

type SystemUndeleteTransaction struct {
	Transaction
	contractID *ContractID
	fileID     *FileID
}

func NewSystemUndeleteTransaction() *SystemUndeleteTransaction {
	transaction := SystemUndeleteTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func systemUndeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{
		Transaction: transaction,
		contractID:  contractIDFromProtobuf(pb.GetSystemUndelete().GetContractID()),
		fileID:      fileIDFromProtobuf(pb.GetSystemUndelete().GetFileID()),
	}
}

func (transaction *SystemUndeleteTransaction) SetContractID(contractID ContractID) *SystemUndeleteTransaction {
	transaction.requireNotFrozen()
	transaction.contractID = &contractID
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetContractID() ContractID {
	if transaction.contractID == nil {
		return ContractID{}
	}

	return *transaction.contractID
}

func (transaction *SystemUndeleteTransaction) SetFileID(fileID FileID) *SystemUndeleteTransaction {
	transaction.requireNotFrozen()
	transaction.fileID = &fileID
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetFileID() FileID {
	if transaction.fileID == nil {
		return FileID{}
	}

	return *transaction.fileID
}

func (transaction *SystemUndeleteTransaction) validateNetworkOnIDs(client *Client) error {
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

func (transaction *SystemUndeleteTransaction) build() *proto.TransactionBody {
	body := &proto.SystemUndeleteTransactionBody{}
	if !transaction.contractID.isZero() {
		body.Id = &proto.SystemUndeleteTransactionBody_ContractID{
			ContractID: transaction.contractID.toProtobuf(),
		}
	}

	if !transaction.fileID.isZero() {
		body.Id = &proto.SystemUndeleteTransactionBody_FileID{
			FileID: transaction.fileID.toProtobuf(),
		}
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_SystemUndelete{
			SystemUndelete: body,
		},
	}
}

func (transaction *SystemUndeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *SystemUndeleteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.SystemUndeleteTransactionBody{}
	if !transaction.contractID.isZero() {
		body.Id = &proto.SystemUndeleteTransactionBody_ContractID{
			ContractID: transaction.contractID.toProtobuf(),
		}
	}

	if !transaction.fileID.isZero() {
		body.Id = &proto.SystemUndeleteTransactionBody_FileID{
			FileID: transaction.fileID.toProtobuf(),
		}
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_SystemUndelete{
			SystemUndelete: body,
		},
	}, nil
}

func _SystemUndeleteTransactionGetMethod(request _Request, channel *_Channel) _Method {
	if channel.getContract() == nil {
		return _Method{
			transaction: channel.getFile().SystemUndelete,
		}
	}

	return _Method{
		transaction: channel.getContract().SystemUndelete,
	}
}

func (transaction *SystemUndeleteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *SystemUndeleteTransaction) Sign(
	privateKey PrivateKey,
) *SystemUndeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *SystemUndeleteTransaction) SignWithOperator(
	client *Client,
) (*SystemUndeleteTransaction, error) {
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
func (transaction *SystemUndeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *SystemUndeleteTransaction {
	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *SystemUndeleteTransaction) Execute(
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
		_Request{
			transaction: &transaction.Transaction,
		},
		_TransactionShouldRetry,
		_TransactionMakeRequest(_Request{
			transaction: &transaction.Transaction,
		}),
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_SystemUndeleteTransactionGetMethod,
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

func (transaction *SystemUndeleteTransaction) Freeze() (*SystemUndeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *SystemUndeleteTransaction) FreezeWith(client *Client) (*SystemUndeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &SystemUndeleteTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *SystemUndeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetMaxTransactionFee(fee Hbar) *SystemUndeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetTransactionMemo(memo string) *SystemUndeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetTransactionValidDuration(duration time.Duration) *SystemUndeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetTransactionID(transactionID TransactionID) *SystemUndeleteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *SystemUndeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *SystemUndeleteTransaction) SetMaxRetry(count int) *SystemUndeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *SystemUndeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *SystemUndeleteTransaction {
	transaction.requireOneNodeAccountID()

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

	return transaction
}

func (transaction *SystemUndeleteTransaction) SetMaxBackoff(max time.Duration) *SystemUndeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *SystemUndeleteTransaction) SetMinBackoff(min time.Duration) *SystemUndeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
