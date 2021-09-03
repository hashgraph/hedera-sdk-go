package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

// ContractExecuteTransaction calls a function of the given smart contract instance, giving it ContractFuncionParams as
// its inputs. it can use the given amount of gas, and any unspent gas will be refunded to the paying account.
//
// If this function stores information, it is charged gas to store it. There is a fee in hbars to maintain that storage
// until the expiration time, and that fee is added as part of the transaction fee.
//
// For a cheaper but more limited _Method to call functions, see ContractCallQuery.
type ContractExecuteTransaction struct {
	Transaction
	contractID *ContractID
	gas        int64
	amount     int64
	parameters []byte
}

// NewContractExecuteTransaction creates a ContractExecuteTransaction transaction which can be
// used to construct and execute a Contract Call Transaction.
func NewContractExecuteTransaction() *ContractExecuteTransaction {
	transaction := ContractExecuteTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func contractExecuteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ContractExecuteTransaction {
	return ContractExecuteTransaction{
		Transaction: transaction,
		contractID:  contractIDFromProtobuf(pb.GetContractCall().GetContractID()),
		gas:         pb.GetContractCall().GetGas(),
		amount:      pb.GetContractCall().GetAmount(),
		parameters:  pb.GetContractCall().GetFunctionParameters(),
	}
}

// SetContractID sets the contract instance to call.
func (transaction *ContractExecuteTransaction) SetContractID(contractID ContractID) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.contractID = &contractID
	return transaction
}

func (transaction *ContractExecuteTransaction) GetContractID() ContractID {
	if transaction.contractID == nil {
		return ContractID{}
	}

	return *transaction.contractID
}

// SetGas sets the maximum amount of gas to use for the call.
func (transaction *ContractExecuteTransaction) SetGas(gas uint64) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.gas = int64(gas)
	return transaction
}

func (transaction *ContractExecuteTransaction) GetGas() uint64 {
	return uint64(transaction.gas)
}

// SetPayableAmount sets the amount of Hbar sent (the function must be payable if this is nonzero)
func (transaction *ContractExecuteTransaction) SetPayableAmount(amount Hbar) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.amount = amount.AsTinybar()
	return transaction
}

func (transaction ContractExecuteTransaction) GetPayableAmount() Hbar {
	return HbarFromTinybar(transaction.amount)
}

// Sets the function parameters
func (transaction *ContractExecuteTransaction) SetFunctionParameters(params []byte) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.parameters = params
	return transaction
}

func (transaction *ContractExecuteTransaction) GetFunctionParameters() []byte {
	return transaction.parameters
}

// SetFunction sets which function to call, and the ContractFunctionParams to pass to the function
func (transaction *ContractExecuteTransaction) SetFunction(name string, params *ContractFunctionParameters) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	if params == nil {
		params = NewContractFunctionParameters()
	}

	transaction.parameters = params.build(&name)
	return transaction
}

func (transaction *ContractExecuteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.contractID != nil {
		if err := transaction.contractID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *ContractExecuteTransaction) build() *proto.TransactionBody {
	body := proto.ContractCallTransactionBody{
		Gas:                transaction.gas,
		Amount:             transaction.amount,
		FunctionParameters: transaction.parameters,
	}

	if !transaction.contractID.isZero() {
		body.ContractID = transaction.contractID.toProtobuf()
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_ContractCall{
			ContractCall: &body,
		},
	}
}

func (transaction *ContractExecuteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *ContractExecuteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := proto.ContractCallTransactionBody{
		Gas:                transaction.gas,
		Amount:             transaction.amount,
		FunctionParameters: transaction.parameters,
	}

	if !transaction.contractID.isZero() {
		body.ContractID = transaction.contractID.toProtobuf()
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_ContractCall{
			ContractCall: &body,
		},
	}, nil
}

func _ContractExecuteTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel.getContract().ContractCallMethod,
	}
}

func (transaction *ContractExecuteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ContractExecuteTransaction) Sign(
	privateKey PrivateKey,
) *ContractExecuteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *ContractExecuteTransaction) SignWithOperator(
	client *Client,
) (*ContractExecuteTransaction, error) {
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
func (transaction *ContractExecuteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractExecuteTransaction {
	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ContractExecuteTransaction) Execute(
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
		_ContractExecuteTransactionGetMethod,
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

func (transaction *ContractExecuteTransaction) Freeze() (*ContractExecuteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ContractExecuteTransaction) FreezeWith(client *Client) (*ContractExecuteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &ContractExecuteTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *ContractExecuteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetMaxTransactionFee(fee Hbar) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ContractExecuteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetTransactionMemo(memo string) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ContractExecuteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ContractExecuteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetTransactionID(transactionID TransactionID) *ContractExecuteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this ContractExecuteTransaction.
func (transaction *ContractExecuteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractExecuteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *ContractExecuteTransaction) SetMaxRetry(count int) *ContractExecuteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *ContractExecuteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractExecuteTransaction {
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

func (transaction *ContractExecuteTransaction) SetMaxBackoff(max time.Duration) *ContractExecuteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *ContractExecuteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *ContractExecuteTransaction) SetMinBackoff(min time.Duration) *ContractExecuteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *ContractExecuteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
