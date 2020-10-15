package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type CryptoTransferTransaction struct {
	Transaction
	pb *proto.CryptoTransferTransactionBody
}

func NewCryptoTransferTransaction() *CryptoTransferTransaction {
	pb := &proto.CryptoTransferTransactionBody{
		Transfers: &proto.TransferList{
			AccountAmounts: []*proto.AccountAmount{},
		},
	}

	transaction := CryptoTransferTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	return &transaction
}

func (transaction *CryptoTransferTransaction) AddTransfer(accountID AccountID, amount Hbar) *CryptoTransferTransaction {
	transaction.pb.Transfers.AccountAmounts = append(transaction.pb.Transfers.AccountAmounts, &proto.AccountAmount{AccountID: accountID.toProtobuf(), Amount: amount.tinybar})
	return transaction
}

func (transaction *CryptoTransferTransaction) AddSender(accountID AccountID, amount Hbar) *CryptoTransferTransaction {
	return transaction.AddTransfer(accountID, amount.negated())
}

func (transaction *CryptoTransferTransaction) AddRecipient(accountID AccountID, amount Hbar) *CryptoTransferTransaction {
	return transaction.AddTransfer(accountID, amount)
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func cryptoTransferTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getCrypto().CryptoTransfer,
	}
}

func (transaction *CryptoTransferTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *CryptoTransferTransaction) Sign(
	privateKey PrivateKey,
) *CryptoTransferTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *CryptoTransferTransaction) SignWithOperator(
	client *Client,
) (*CryptoTransferTransaction, error) {
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
func (transaction *CryptoTransferTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *CryptoTransferTransaction {
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
func (transaction *CryptoTransferTransaction) Execute(
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
		cryptoTransferTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
		query_makePaymentTransaction,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{TransactionID: transaction.id}, nil
}

func (transaction *CryptoTransferTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_CryptoTransfer{
		CryptoTransfer: transaction.pb,
	}

	return true
}

func (transaction *CryptoTransferTransaction) Freeze() (*CryptoTransferTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *CryptoTransferTransaction) FreezeWith(client *Client) (*CryptoTransferTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *CryptoTransferTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this CryptoTransferTransaction.
func (transaction *CryptoTransferTransaction) SetMaxTransactionFee(fee Hbar) *CryptoTransferTransaction {
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *CryptoTransferTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this CryptoTransferTransaction.
func (transaction *CryptoTransferTransaction) SetTransactionMemo(memo string) *CryptoTransferTransaction {
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *CryptoTransferTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this CryptoTransferTransaction.
func (transaction *CryptoTransferTransaction) SetTransactionValidDuration(duration time.Duration) *CryptoTransferTransaction {
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *CryptoTransferTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this CryptoTransferTransaction.
func (transaction *CryptoTransferTransaction) SetTransactionID(transactionID TransactionID) *CryptoTransferTransaction {
	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *CryptoTransferTransaction) GetNodeID() AccountID {
	return transaction.Transaction.GetNodeID()
}

// SetNodeID sets the node AccountID for this CryptoTransferTransaction.
func (transaction *CryptoTransferTransaction) SetNodeID(nodeID AccountID) *CryptoTransferTransaction {
	transaction.Transaction.SetNodeID(nodeID)
	return transaction
}
