package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type SystemUndeleteTransaction struct {
	Transaction
	pb *proto.SystemUndeleteTransactionBody
}

func NewSystemUndeleteTransaction() *SystemUndeleteTransaction {
	pb := &proto.SystemUndeleteTransactionBody{}

	transaction := SystemUndeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func (transaction *SystemUndeleteTransaction) SetContractID(contractID ContractID) *SystemUndeleteTransaction {
	transaction.pb.Id = &proto.SystemUndeleteTransactionBody_ContractID{ContractID: contractID.toProto()}
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetContract() ContractID {
	return contractIDFromProto(transaction.pb.GetContractID())
}

func (transaction *SystemUndeleteTransaction) SetFileID(fileID FileID) *SystemUndeleteTransaction {
	transaction.pb.Id = &proto.SystemUndeleteTransactionBody_FileID{FileID: fileID.toProto()}
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetFileID() FileID {
	return fileIDFromProto(transaction.pb.GetFileID())
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func systemUndeleteTransaction_getMethod(request request, channel *channel) method {
	switch request.transaction.pbBody.GetSystemUndelete().Id.(type) {
	case *proto.SystemUndeleteTransactionBody_ContractID:
		return method{
			transaction: channel.getContract().SystemUndelete,
		}
	default:
		return method{
			transaction: channel.getFile().SystemUndelete,
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

	if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
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
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.transactions); index++ {
		signature := signer(transaction.transactions[index].GetBodyBytes())

		transaction.signatures[index].SigPair = append(
			transaction.signatures[index].SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *SystemUndeleteTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	transactionID := transaction.id

	if !client.GetOperatorID().isZero() && client.GetOperatorID().equals(transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorKey(),
			client.operator.signer,
		)
	}

	_, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeId,
		systemUndeleteTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
		query_makePaymentTransaction,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{TransactionID: transaction.id}, nil
}

func (transaction *SystemUndeleteTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_SystemUndelete{
		SystemUndelete: transaction.pb,
	}

	return true
}

func (transaction *SystemUndeleteTransaction) Freeze() (*SystemUndeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *SystemUndeleteTransaction) FreezeWith(client *Client) (*SystemUndeleteTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *SystemUndeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetMaxTransactionFee(fee Hbar) *SystemUndeleteTransaction {
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetTransactionMemo(memo string) *SystemUndeleteTransaction {
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetTransactionValidDuration(duration time.Duration) *SystemUndeleteTransaction {
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetTransactionID(transactionID TransactionID) *SystemUndeleteTransaction {
	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetNodeID() AccountID {
	return transaction.Transaction.GetNodeID()
}

// SetNodeID sets the node AccountID for this SystemUndeleteTransaction.
func (transaction *SystemUndeleteTransaction) SetNodeID(nodeID AccountID) *SystemUndeleteTransaction {
	transaction.Transaction.SetNodeID(nodeID)
	return transaction
}
