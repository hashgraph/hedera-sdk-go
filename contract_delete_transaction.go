package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

type ContractDeleteTransaction struct {
	Transaction
	contractID        ContractID
	transferContactID ContractID
	transferAccountID AccountID
}

func NewContractDeleteTransaction() *ContractDeleteTransaction {
	transaction := ContractDeleteTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func contractDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ContractDeleteTransaction {
	return ContractDeleteTransaction{
		Transaction:       transaction,
		contractID:        contractIDFromProtobuf(pb.GetContractDeleteInstance().GetContractID()),
		transferContactID: contractIDFromProtobuf(pb.GetContractDeleteInstance().GetTransferContractID()),
		transferAccountID: accountIDFromProtobuf(pb.GetContractDeleteInstance().GetTransferAccountID()),
	}
}

//Sets the contract ID which should be deleted.
func (transaction *ContractDeleteTransaction) SetContractID(id ContractID) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.contractID = id
	return transaction
}

func (transaction *ContractDeleteTransaction) GetContractID() ContractID {
	return transaction.contractID
}

//Sets the contract ID which will receive all remaining hbars.
func (transaction *ContractDeleteTransaction) SetTransferContractID(id ContractID) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.transferContactID = id

	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransferContractID() ContractID {
	return transaction.transferContactID
}

//Sets the account ID which will receive all remaining hbars.
func (transaction *ContractDeleteTransaction) SetTransferAccountID(accountID AccountID) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.transferAccountID = accountID

	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransferAccountID() AccountID {
	return transaction.transferAccountID
}

func (transaction *ContractDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = transaction.contractID.Validate(client)
	if err != nil {
		return err
	}
	err = transaction.transferContactID.Validate(client)
	if err != nil {
		return err
	}
	err = transaction.transferAccountID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *ContractDeleteTransaction) build() *proto.TransactionBody {
	body := &proto.ContractDeleteTransactionBody{}

	if !transaction.contractID.isZero() {
		body.ContractID = transaction.contractID.toProtobuf()
	}

	if !transaction.transferContactID.isZero() {
		body.Obtainers = &proto.ContractDeleteTransactionBody_TransferContractID{
			TransferContractID: transaction.transferContactID.toProtobuf(),
		}
	}

	if !transaction.transferAccountID.isZero() {
		body.Obtainers = &proto.ContractDeleteTransactionBody_TransferAccountID{
			TransferAccountID: transaction.transferAccountID.toProtobuf(),
		}
	}

	pb := proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: body,
		},
	}

	return &pb
}

func (transaction *ContractDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *ContractDeleteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.ContractDeleteTransactionBody{}

	if !transaction.contractID.isZero() {
		body.ContractID = transaction.contractID.toProtobuf()
	}

	if !transaction.transferContactID.isZero() {
		body.Obtainers = &proto.ContractDeleteTransactionBody_TransferContractID{
			TransferContractID: transaction.transferContactID.toProtobuf(),
		}
	}

	if !transaction.transferAccountID.isZero() {
		body.Obtainers = &proto.ContractDeleteTransactionBody_TransferAccountID{
			TransferAccountID: transaction.transferAccountID.toProtobuf(),
		}
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: body,
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func contractDeleteTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getContract().DeleteContract,
	}
}

func (transaction *ContractDeleteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
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
func (transaction *ContractDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ContractDeleteTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
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
		contractDeleteTransaction_getMethod,
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

func (transaction *ContractDeleteTransaction) Freeze() (*ContractDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ContractDeleteTransaction) FreezeWith(client *Client) (*ContractDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &ContractDeleteTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
}

func (transaction *ContractDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetMaxTransactionFee(fee Hbar) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionMemo(memo string) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ContractDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetTransactionID(transactionID TransactionID) *ContractDeleteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the node AccountID for this ContractDeleteTransaction.
func (transaction *ContractDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ContractDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *ContractDeleteTransaction) SetMaxRetry(count int) *ContractDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *ContractDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *ContractDeleteTransaction {
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
