package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"

	"time"
)

type ContractDeleteTransaction struct {
	Transaction
	contractID        *ContractID
	transferContactID *ContractID
	transferAccountID *AccountID
}

func NewContractDeleteTransaction() *ContractDeleteTransaction {
	transaction := ContractDeleteTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _ContractDeleteTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) ContractDeleteTransaction {
	return ContractDeleteTransaction{
		Transaction:       transaction,
		contractID:        _ContractIDFromProtobuf(pb.GetContractDeleteInstance().GetContractID()),
		transferContactID: _ContractIDFromProtobuf(pb.GetContractDeleteInstance().GetTransferContractID()),
		transferAccountID: _AccountIDFromProtobuf(pb.GetContractDeleteInstance().GetTransferAccountID()),
	}
}

// Sets the contract ID which should be deleted.
func (transaction *ContractDeleteTransaction) SetContractID(contractID ContractID) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.contractID = &contractID
	return transaction
}

func (transaction *ContractDeleteTransaction) GetContractID() ContractID {
	if transaction.contractID == nil {
		return ContractID{}
	}

	return *transaction.contractID
}

// Sets the contract ID which will receive all remaining hbars.
func (transaction *ContractDeleteTransaction) SetTransferContractID(transferContactID ContractID) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.transferContactID = &transferContactID
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransferContractID() ContractID {
	if transaction.transferContactID == nil {
		return ContractID{}
	}

	return *transaction.transferContactID
}

// Sets the account ID which will receive all remaining hbars.
func (transaction *ContractDeleteTransaction) SetTransferAccountID(accountID AccountID) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.transferAccountID = &accountID

	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransferAccountID() AccountID {
	if transaction.transferAccountID == nil {
		return AccountID{}
	}

	return *transaction.transferAccountID
}

func (transaction *ContractDeleteTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.contractID != nil {
		if err := transaction.contractID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if transaction.transferContactID != nil {
		if err := transaction.transferContactID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if transaction.transferAccountID != nil {
		if err := transaction.transferAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *ContractDeleteTransaction) _Build() *services.TransactionBody {
	body := &services.ContractDeleteTransactionBody{}

	if transaction.contractID != nil {
		body.ContractID = transaction.contractID._ToProtobuf()
	}

	if transaction.transferContactID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferContractID{
			TransferContractID: transaction.transferContactID._ToProtobuf(),
		}
	}

	if transaction.transferAccountID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferAccountID{
			TransferAccountID: transaction.transferAccountID._ToProtobuf(),
		}
	}

	pb := services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: body,
		},
	}

	return &pb
}

func (transaction *ContractDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *ContractDeleteTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.ContractDeleteTransactionBody{}

	if transaction.contractID != nil {
		body.ContractID = transaction.contractID._ToProtobuf()
	}

	if transaction.transferContactID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferContractID{
			TransferContractID: transaction.transferContactID._ToProtobuf(),
		}
	}

	if transaction.transferAccountID != nil {
		body.Obtainers = &services.ContractDeleteTransactionBody_TransferAccountID{
			TransferAccountID: transaction.transferAccountID._ToProtobuf(),
		}
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: body,
		},
	}, nil
}

func _ContractDeleteTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetContract().DeleteContract,
	}
}

func (transaction *ContractDeleteTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ContractDeleteTransaction) Sign(
	privateKey PrivateKey,
) *ContractDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *ContractDeleteTransaction) SignWithOperator(
	client *Client,
) (*ContractDeleteTransaction, error) {
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
func (transaction *ContractDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractDeleteTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ContractDeleteTransaction) Execute(
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
		_ContractDeleteTransactionGetMethod,
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

func (transaction *ContractDeleteTransaction) Freeze() (*ContractDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ContractDeleteTransaction) FreezeWith(client *Client) (*ContractDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &ContractDeleteTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *ContractDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetMaxTransactionFee(fee Hbar) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionMemo(memo string) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionID(transactionID TransactionID) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractDeleteTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *ContractDeleteTransaction) SetMaxRetry(count int) *ContractDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *ContractDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractDeleteTransaction {
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

func (transaction *ContractDeleteTransaction) SetMaxBackoff(max time.Duration) *ContractDeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *ContractDeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *ContractDeleteTransaction) SetMinBackoff(min time.Duration) *ContractDeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *ContractDeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
