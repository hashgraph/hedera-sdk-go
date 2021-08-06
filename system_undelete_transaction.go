package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

type SystemUndeleteTransaction struct {
	Transaction
	contractID ContractID
	fileID     FileID
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

func (transaction *SystemUndeleteTransaction) SetContractID(id ContractID) *SystemUndeleteTransaction {
	transaction.requireNotFrozen()
	transaction.contractID = id
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetContract() ContractID {
	return transaction.contractID
}

func (transaction *SystemUndeleteTransaction) SetFileID(id FileID) *SystemUndeleteTransaction {
	transaction.requireNotFrozen()
	transaction.fileID = id
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetFileID() FileID {
	return transaction.fileID
}

func (transaction *SystemUndeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = transaction.contractID.Validate(client)
	if err != nil {
		return err
	}
	err = transaction.fileID.Validate(client)
	if err != nil {
		return err
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

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func systemUndeleteTransaction_getMethod(request request, channel *channel) method {
	if channel.getContract() == nil {
		return method{
			transaction: channel.getFile().SystemUndelete,
		}
	} else {
		return method{
			transaction: channel.getContract().SystemUndelete,
		}
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
func (transaction *SystemUndeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *SystemUndeleteTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

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
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest(request{
			transaction: &transaction.Transaction,
		}),
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		systemUndeleteTransaction_getMethod,
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

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
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

// SetNodeAccountID sets the node AccountID for this SystemUndeleteTransaction.
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
