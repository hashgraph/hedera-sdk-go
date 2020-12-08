package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"github.com/pkg/errors"
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

func systemUndeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) SystemUndeleteTransaction {
	return SystemUndeleteTransaction{
		Transaction: transaction,
		pb:          pb.GetSystemUndelete(),
	}
}

func (transaction *SystemUndeleteTransaction) SetContractID(contractID ContractID) *SystemUndeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Id = &proto.SystemUndeleteTransactionBody_ContractID{ContractID: contractID.toProtobuf()}
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetContract() ContractID {
	return contractIDFromProtobuf(transaction.pb.GetContractID())
}

func (transaction *SystemUndeleteTransaction) SetFileID(fileID FileID) *SystemUndeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Id = &proto.SystemUndeleteTransactionBody_FileID{FileID: fileID.toProtobuf()}
	return transaction
}

func (transaction *SystemUndeleteTransaction) GetFileID() FileID {
	return fileIDFromProtobuf(transaction.pb.GetFileID())
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

	if client == nil {
		return nil, errors.Wrap(errNoClientProvided, "for SignWithOperator")
	} else if client.operator == nil {
		return nil, errors.Wrap(errClientOperatorSigning, "for SignWithOperator")
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return transaction, errors.Wrap(err, "FreezeWith in SignWithOperator")
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
func (transaction *SystemUndeleteTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, errors.Wrap(err, "FreezeWith in Execute")
		}
	}

	transactionID := transaction.transactionIDs[0]

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(transactionID.AccountID) {
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
		systemUndeleteTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, errors.Wrap(err, "execution error")
	}

	return TransactionResponse{
		TransactionID: transaction.transactionIDs[0],
		NodeID:        resp.transaction.NodeID,
	}, nil
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
		return transaction, errors.Wrap(err, "initTransactionID in SystemUndeleteTransaction.FreezeWith")
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

func (transaction *SystemUndeleteTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
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
