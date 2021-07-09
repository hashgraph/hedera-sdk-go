package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type TransferTransaction struct {
	Transaction
	pb             *proto.CryptoTransferTransactionBody
	tokenTransfers map[TokenID]map[AccountID]int64
	hbarTransfers  map[AccountID]Hbar
	nftTransfers   map[TokenID][]TokenNftTransfer
}

func NewTransferTransaction() *TransferTransaction {
	pb := &proto.CryptoTransferTransactionBody{
		Transfers: &proto.TransferList{
			AccountAmounts: []*proto.AccountAmount{},
		},
	}

	transaction := TransferTransaction{
		pb:             pb,
		Transaction:    newTransaction(),
		tokenTransfers: make(map[TokenID]map[AccountID]int64),
		hbarTransfers:  make(map[AccountID]Hbar),
	}

	transaction.SetMaxTransactionFee(NewHbar(1))

	return &transaction
}

func transferTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TransferTransaction {
	hbarTransfers := make(map[AccountID]Hbar)
	tokenTransfers := make(map[TokenID]map[AccountID]int64)

	for _, aa := range pb.GetCryptoTransfer().GetTransfers().AccountAmounts {
		accountID := accountIDFromProtobuf(aa.AccountID, nil)
		amount := HbarFromTinybar(aa.Amount)

		if value, ok := hbarTransfers[accountID]; ok {
			hbarTransfers[accountID] = HbarFromTinybar(amount.AsTinybar() + value.AsTinybar())
		} else {
			hbarTransfers[accountID] = amount
		}
	}

	for _, tokenTransfersList := range pb.GetCryptoTransfer().GetTokenTransfers() {
		tokenID := tokenIDFromProtobuf(tokenTransfersList.Token, nil)

		var currentTokenTransfers map[AccountID]int64

		if value, ok := tokenTransfers[tokenID]; ok {
			currentTokenTransfers = value
		} else {
			currentTokenTransfers = make(map[AccountID]int64)
		}

		for _, aa := range tokenTransfersList.GetTransfers() {
			accountID := accountIDFromProtobuf(aa.AccountID, nil)

			if value, ok := currentTokenTransfers[accountID]; ok {
				currentTokenTransfers[accountID] = aa.Amount + value
			} else {
				currentTokenTransfers[accountID] = aa.Amount
			}
		}

		tokenTransfers[tokenID] = currentTokenTransfers
	}

	return TransferTransaction{
		Transaction:    transaction,
		pb:             pb.GetCryptoTransfer(),
		hbarTransfers:  hbarTransfers,
		tokenTransfers: tokenTransfers,
	}
}

func (transaction *TransferTransaction) GetNftTransfers() map[TokenID][]TokenNftTransfer {
	nftTransferMap := make(map[TokenID][]TokenNftTransfer, len(transaction.pb.TokenTransfers))

	if len(transaction.pb.TokenTransfers) == 0 {
		return nftTransferMap
	}

	for _, tokenTransfer := range transaction.pb.TokenTransfers {
		for _, nftTransfer := range tokenTransfer.NftTransfers {
			token := tokenIDFromProtobuf(tokenTransfer.Token, nil)
			nftTransferMap[token] = append(nftTransferMap[token], nftTransferFromProtobuf(nftTransfer, nil))
		}
	}

	return nftTransferMap
}

func (transaction *TransferTransaction) GetTokenTransfers() map[TokenID][]TokenTransfer {
	transfers := make(map[TokenID][]TokenTransfer)

	for tokenID, tokenTransfers := range transaction.tokenTransfers {
		tokenTransfersList := make([]TokenTransfer, 0)

		for accountID, amount := range tokenTransfers {
			tokenTransfersList = append(tokenTransfersList, TokenTransfer{
				AccountID: accountID,
				Amount:    amount,
			})
		}

		transfers[tokenID] = tokenTransfersList
	}

	return transfers
}

func (transaction *TransferTransaction) GetHbarTransfers() map[AccountID]Hbar {
	return transaction.hbarTransfers
}

func (transaction *TransferTransaction) AddHbarTransfer(accountID AccountID, amount Hbar) *TransferTransaction {
	transaction.requireNotFrozen()

	if value, ok := transaction.hbarTransfers[accountID]; ok {
		transaction.hbarTransfers[accountID] = HbarFromTinybar(amount.AsTinybar() + value.AsTinybar())
	} else {
		transaction.hbarTransfers[accountID] = amount
	}

	return transaction
}

func (transaction *TransferTransaction) AddTokenTransfer(tokenID TokenID, accountID AccountID, value int64) *TransferTransaction {
	transaction.requireNotFrozen()

	var tokenTransfers map[AccountID]int64
	var amount int64

	if value, ok := transaction.tokenTransfers[tokenID]; ok {
		tokenTransfers = value
	} else {
		tokenTransfers = make(map[AccountID]int64)
	}

	if transfer, ok := tokenTransfers[accountID]; ok {
		amount = transfer + value
	} else {
		amount = value
	}

	tokenTransfers[accountID] = amount
	transaction.tokenTransfers[tokenID] = tokenTransfers

	return transaction
}

func (transaction *TransferTransaction) AddNftTransfer(nftID NftID, sender AccountID, receiver AccountID) *TransferTransaction {
	transaction.requireNotFrozen()

	if transaction.nftTransfers == nil {
		transaction.nftTransfers = make(map[TokenID][]TokenNftTransfer, 0)
	}

	if transaction.nftTransfers[nftID.TokenID] == nil {
		transaction.nftTransfers[nftID.TokenID] = make([]TokenNftTransfer, 0)
	}

	transaction.nftTransfers[nftID.TokenID] = append(transaction.nftTransfers[nftID.TokenID], TokenNftTransfer{
		SenderAccountID:   sender,
		ReceiverAccountID: receiver,
		SerialNumber:      nftID.SerialNumber,
	})

	return transaction
}

func (transaction *TransferTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	for tokenID, accountMap := range transaction.tokenTransfers {
		err = tokenID.Validate(client)
		for accountID, _ := range accountMap {
			err = accountID.Validate(client)
		}
	}
	for nftID, _ := range transaction.nftTransfers {
		err = nftID.Validate(client)
	}
	for accountID, _ := range transaction.hbarTransfers {
		err = accountID.Validate(client)
	}
	if err != nil {
		return err
	}

	return nil
}

func (transaction *TransferTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TransferTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	transaction.buildHbarTransfers()
	transaction.buildTokenTransfers()

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_CryptoTransfer{
			CryptoTransfer: &proto.CryptoTransferTransactionBody{
				Transfers:      transaction.pb.GetTransfers(),
				TokenTransfers: transaction.pb.GetTokenTransfers(),
			},
		},
	}, nil
}

func transferTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getCrypto().CryptoTransfer,
	}
}

func (transaction *TransferTransaction) AddSignature(publicKey PublicKey, signature []byte) *TransferTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
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
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TransferTransaction) Execute(
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
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		transferTransaction_getMethod,
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

func (transaction *TransferTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_CryptoTransfer{
		CryptoTransfer: transaction.pb,
	}

	return true
}

func (transaction *TransferTransaction) buildHbarTransfers() {
	transaction.pb.Transfers.AccountAmounts = make([]*proto.AccountAmount, 0)
	for accountID, amount := range transaction.hbarTransfers {
		transaction.pb.Transfers.AccountAmounts = append(transaction.pb.Transfers.AccountAmounts, &proto.AccountAmount{
			AccountID: accountID.toProtobuf(),
			Amount:    amount.AsTinybar(),
		})
	}
}

func (transaction *TransferTransaction) buildTokenTransfers() {
	if transaction.pb.TokenTransfers == nil {
		transaction.pb.TokenTransfers = make([]*proto.TokenTransferList, 0)
	}

	for tokenID, tokenTransfers := range transaction.tokenTransfers {
		transfers := make([]*proto.AccountAmount, 0)

		for accountID, amount := range tokenTransfers {
			transfers = append(transfers, &proto.AccountAmount{
				AccountID: accountID.toProtobuf(),
				Amount:    amount,
			})
		}

		transaction.pb.TokenTransfers = append(transaction.pb.TokenTransfers, &proto.TokenTransferList{
			Token:     tokenID.toProtobuf(),
			Transfers: transfers,
		})
	}
}

func (transaction *TransferTransaction) buildNftTransfers() {
	if transaction.pb.TokenTransfers == nil {
		transaction.pb.TokenTransfers = make([]*proto.TokenTransferList, 0)
	}

	for tokenID, nftTransfers := range transaction.nftTransfers {
		transfers := make([]*proto.NftTransfer, 0)

		for _, nftTransfer := range nftTransfers {
			transfers = append(transfers, nftTransfer.toProtobuf())
		}

		transaction.pb.TokenTransfers = append(transaction.pb.TokenTransfers, &proto.TokenTransferList{
			Token:        tokenID.toProtobuf(),
			NftTransfers: transfers,
		})
	}
}

func (transaction *TransferTransaction) Freeze() (*TransferTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TransferTransaction) FreezeWith(client *Client) (*TransferTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}

	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TransferTransaction{}, err
	}
	transaction.buildHbarTransfers()
	transaction.buildTokenTransfers()
	transaction.buildNftTransfers()

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
