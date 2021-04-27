package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type SystemDeleteTransaction struct {
	Transaction
	pb *proto.SystemDeleteTransactionBody
}

func NewSystemDeleteTransaction() *SystemDeleteTransaction {
	pb := &proto.SystemDeleteTransactionBody{}

	transaction := SystemDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func systemDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) SystemDeleteTransaction {
	return SystemDeleteTransaction{
		Transaction: transaction,
		pb:          pb.GetSystemDelete(),
	}
}

func (transaction *SystemDeleteTransaction) SetExpirationTime(expiration time.Time) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.ExpirationTime = &proto.TimestampSeconds{
		Seconds: expiration.Unix(),
	}
	return transaction
}

func (transaction *SystemDeleteTransaction) GetExpirationTime() int64 {
	return transaction.pb.GetExpirationTime().Seconds
}

func (transaction *SystemDeleteTransaction) SetContractID(contractID ContractID) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Id = &proto.SystemDeleteTransactionBody_ContractID{ContractID: contractID.toProtobuf()}
	return transaction
}

func (transaction *SystemDeleteTransaction) GetContract() ContractID {
	return contractIDFromProtobuf(transaction.pb.GetContractID())
}

func (transaction *SystemDeleteTransaction) SetFileID(fileID FileID) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Id = &proto.SystemDeleteTransactionBody_FileID{FileID: fileID.toProtobuf()}
	return transaction
}

func (transaction *SystemDeleteTransaction) GetFileID() FileID {
	return fileIDFromProtobuf(transaction.pb.GetFileID())
}

func (transaction *SystemDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *SystemDeleteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_SystemDelete{
			SystemDelete: &proto.SystemDeleteTransactionBody{
				Id:             nil,
				ExpirationTime: transaction.pb.GetExpirationTime(),
			},
		},
	}

	switch transaction.pb.GetId().(type) {
	case *proto.SystemDeleteTransactionBody_ContractID:
		body.GetSystemUndelete().Id = &proto.SystemUndeleteTransactionBody_ContractID{ContractID: transaction.pb.GetContractID()}
	case *proto.SystemDeleteTransactionBody_FileID:
		body.GetSystemUndelete().Id = &proto.SystemUndeleteTransactionBody_FileID{FileID: transaction.pb.GetFileID()}
	}

	return body, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func systemDeleteTransaction_getMethod(request request, channel *channel) method {
	switch request.transaction.pbBody.GetSystemDelete().Id.(type) {
	case *proto.SystemDeleteTransactionBody_ContractID:
		return method{
			transaction: channel.getContract().SystemDelete,
		}
	default:
		return method{
			transaction: channel.getFile().SystemDelete,
		}
	}
}

func (transaction *SystemDeleteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
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
func (transaction *SystemDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *SystemDeleteTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	} else {
		transaction.transactions = make([]*proto.Transaction, 0)
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.signedTransactions); index++ {
		signature := signer(transaction.signedTransactions[index].GetBodyBytes())

		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
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
		systemDeleteTransaction_getMethod,
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

func (transaction *SystemDeleteTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_SystemDelete{
		SystemDelete: transaction.pb,
	}

	return true
}

func (transaction *SystemDeleteTransaction) Freeze() (*SystemDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *SystemDeleteTransaction) FreezeWith(client *Client) (*SystemDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *SystemDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetMaxTransactionFee(fee Hbar) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *SystemDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetTransactionMemo(memo string) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *SystemDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *SystemDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetTransactionID(transactionID TransactionID) *SystemDeleteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *SystemDeleteTransaction) SetMaxRetry(count int) *SystemDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *SystemDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *SystemDeleteTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
