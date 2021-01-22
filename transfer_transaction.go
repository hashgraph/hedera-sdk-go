package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TransferTransaction struct {
	Transaction
	pb           *proto.CryptoTransferTransactionBody
	tokenIndexes map[TokenID]int
}

func NewTransferTransaction() *TransferTransaction {
	pb := &proto.CryptoTransferTransactionBody{
		Transfers: &proto.TransferList{
			AccountAmounts: []*proto.AccountAmount{},
		},
	}

	transaction := TransferTransaction{
		pb:           pb,
		Transaction:  newTransaction(),
		tokenIndexes: make(map[TokenID]int),
	}

	transaction.SetMaxTransactionFee(NewHbar(1))

	return &transaction
}

func transferTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TransferTransaction {
	return TransferTransaction{
		Transaction: transaction,
		pb:          pb.GetCryptoTransfer(),
	}
}

func (transaction *TransferTransaction) GetTokenTransfers() map[TokenID][]TokenTransfer {
	tokenTransferMap := make(map[TokenID][]TokenTransfer, len(transaction.pb.TokenTransfers))

	if len(transaction.pb.TokenTransfers) == 0 {
		return tokenTransferMap
	}

	for _, tokenTransfer := range transaction.pb.TokenTransfers {
		for _, accountAmount := range tokenTransfer.Transfers {
			token := tokenIDFromProtobuf(tokenTransfer.Token)
			tokenTransferMap[token] = append(tokenTransferMap[token], tokenTransferFromProtobuf(accountAmount))
		}
	}

	return tokenTransferMap
}

func (transaction *TransferTransaction) GetHbarTransfers() map[AccountID]Hbar {
	hbarTransferMap := make(map[AccountID]Hbar, len(transaction.pb.Transfers.AccountAmounts))

	if len(transaction.pb.Transfers.AccountAmounts) == 0 {
		return hbarTransferMap
	}

	for _, hbarTransfer := range transaction.pb.Transfers.AccountAmounts {
		hbarTransferMap[accountIDFromProtobuf(hbarTransfer.AccountID)] = HbarFromTinybar(hbarTransfer.Amount)
	}

	return hbarTransferMap
}

func (transaction *TransferTransaction) AddHbarTransfer(accountID AccountID, amount Hbar) *TransferTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Transfers.AccountAmounts = append(transaction.pb.Transfers.AccountAmounts, &proto.AccountAmount{AccountID: accountID.toProtobuf(), Amount: amount.tinybar})
	return transaction
}

func (transaction *TransferTransaction) AddTokenTransfer(tokenID TokenID, accountID AccountID, value int64) *TransferTransaction {
	transaction.requireNotFrozen()

	accountAmount := proto.AccountAmount{
		AccountID: accountID.toProtobuf(),
		Amount:    value,
	}

	if index, ok := transaction.tokenIndexes[tokenID]; ok {
		transaction.pb.TokenTransfers[index].Transfers = append(
			transaction.pb.TokenTransfers[index].Transfers,
			&accountAmount,
		)
	} else {
		transaction.tokenIndexes[tokenID] = len(transaction.pb.TokenTransfers)
		transaction.pb.TokenTransfers = append(transaction.pb.TokenTransfers, &proto.TokenTransferList{
			Token:     tokenID.toProtobuf(),
			Transfers: []*proto.AccountAmount{&accountAmount},
		})
	}

	return transaction
}

func transferTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getCrypto().CryptoTransfer,
	}
}

func (transaction *TransferTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TransferTransaction) Sign(
	privateKey PrivateKey,
) *TransferTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TransferTransaction) SignWithOperator(
	client *Client,
) (*TransferTransaction, error) {
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
func (transaction *TransferTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TransferTransaction {
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
func (transaction *TransferTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil || client.operator == nil {
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
		transferTransaction_getMethod,
		transaction_mapResponseStatus,
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

func (transaction *TransferTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_CryptoTransfer{
		CryptoTransfer: transaction.pb,
	}

	return true
}

func (transaction *TransferTransaction) Freeze() (*TransferTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TransferTransaction) FreezeWith(client *Client) (*TransferTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *TransferTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenUpdateTransaction.
func (transaction *TransferTransaction) SetMaxTransactionFee(fee Hbar) *TransferTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TransferTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenUpdateTransaction.
func (transaction *TransferTransaction) SetTransactionMemo(memo string) *TransferTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TransferTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenUpdateTransaction.
func (transaction *TransferTransaction) SetTransactionValidDuration(duration time.Duration) *TransferTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TransferTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenUpdateTransaction.
func (transaction *TransferTransaction) SetTransactionID(transactionID TransactionID) *TransferTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenUpdateTransaction.
func (transaction *TransferTransaction) SetNodeAccountIDs(nodeID []AccountID) *TransferTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TransferTransaction) SetMaxRetry(count int) *TransferTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}
