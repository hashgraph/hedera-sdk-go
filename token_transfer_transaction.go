package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenTransferTransaction struct {
	Transaction
	pb *proto.TokenTransfersTransactionBody
}

func NewTokenTransferTransaction() *TokenTransferTransaction {
	pb := &proto.TokenTransfersTransactionBody{}

	transaction := TokenTransferTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func (transaction *TokenTransferTransaction) GetTransfers() map[TokenID][]TokenTransfer {
	tokenTransferMap := make(map[TokenID][]TokenTransfer, len(transaction.pb.TokenTransfers))
	for _, tokenTransfer := range transaction.pb.TokenTransfers {
		for _, accountAmount := range tokenTransfer.Transfers {
			token := tokenIDFromProtobuf(tokenTransfer.Token)
			tokenTransferMap[token] = append(tokenTransferMap[token], tokenTransferFromProtobuf(accountAmount))
		}
	}

	return tokenTransferMap
}

func (transaction *TokenTransferTransaction) AddSender(tokenID TokenID, accountID AccountID, value int64) *TokenTransferTransaction {
	return transaction.AddTransfer(tokenID, accountID, -value)
}

func (transaction *TokenTransferTransaction) AddRecipient(tokenID TokenID, accountID AccountID, value int64) *TokenTransferTransaction {
	return transaction.AddTransfer(tokenID, accountID, value)
}

func (transaction *TokenTransferTransaction) AddTransfer(tokenID TokenID, accountID AccountID, value int64) *TokenTransferTransaction {
	transaction.requireNotFrozen()

	println("value", value)

	accountAmount := proto.AccountAmount{
		AccountID: accountID.toProtobuf(),
		Amount:    value,
	}

	accountAmountArray := make([]*proto.AccountAmount, 1)
	accountAmountArray[0] = &accountAmount

	println("value", accountAmountArray[0].String())

	tokenTransfers := &proto.TokenTransferList{
		Token:     tokenID.toProtobuf(),
		Transfers: accountAmountArray,
	}

	transaction.pb.TokenTransfers = append(transaction.pb.TokenTransfers, tokenTransfers)

	return transaction
}

func tokenTransferTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().TransferTokens,
	}
}

func (transaction *TokenTransferTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenTransferTransaction) Sign(
	privateKey PrivateKey,
) *TokenTransferTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenTransferTransaction) SignWithOperator(
	client *Client,
) (*TokenTransferTransaction, error) {
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
func (transaction *TokenTransferTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenTransferTransaction {
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
func (transaction *TokenTransferTransaction) Execute(
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

	resp, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeId,
		tokenTransferTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{
		TransactionID: transaction.id,
		NodeID:        resp.transaction.NodeID,
	}, nil
}

func (transaction *TokenTransferTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenTransfers{
		TokenTransfers: transaction.pb,
	}

	return true
}

func (transaction *TokenTransferTransaction) Freeze() (*TokenTransferTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenTransferTransaction) FreezeWith(client *Client) (*TokenTransferTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *TokenTransferTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenUpdateTransaction.
func (transaction *TokenTransferTransaction) SetMaxTransactionFee(fee Hbar) *TokenTransferTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenTransferTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenUpdateTransaction.
func (transaction *TokenTransferTransaction) SetTransactionMemo(memo string) *TokenTransferTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenTransferTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenUpdateTransaction.
func (transaction *TokenTransferTransaction) SetTransactionValidDuration(duration time.Duration) *TokenTransferTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenTransferTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenUpdateTransaction.
func (transaction *TokenTransferTransaction) SetTransactionID(transactionID TransactionID) *TokenTransferTransaction {
	transaction.requireNotFrozen()
	transaction.id = transactionID
	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *TokenTransferTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
}

// SetNodeTokenID sets the node TokenID for this TokenUpdateTransaction.
func (transaction *TokenTransferTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenTransferTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}
