package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type accountAmount struct {
	accountID AccountID
	amount    Hbar
}

type tokenTransfer struct {
	tokenID TokenID

	transfers   []TokenTransfer
	transferMap map[AccountID]int

	nftTransfers []TokenNftTransfer
}

type TransferTransaction struct {
	Transaction

	hbarTransfers   []accountAmount
	hbarTransferMap map[AccountID]int

	tokenTransfers   []tokenTransfer
	tokenTransferMap map[TokenID]int
}

func NewTransferTransaction() *TransferTransaction {
	transaction := TransferTransaction{
		Transaction: _NewTransaction(),

		tokenTransfers:   []tokenTransfer{},
		tokenTransferMap: map[TokenID]int{},

		hbarTransfers:   []accountAmount{},
		hbarTransferMap: map[AccountID]int{},
	}

	transaction.SetMaxTransactionFee(NewHbar(1))

	return &transaction
}

func _TransferTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TransferTransaction {
	tx := TransferTransaction{
		Transaction: transaction,

		tokenTransfers:   []tokenTransfer{},
		tokenTransferMap: map[TokenID]int{},

		hbarTransfers:   []accountAmount{},
		hbarTransferMap: map[AccountID]int{},
	}

	for _, aa := range pb.GetCryptoTransfer().GetTransfers().AccountAmounts {
		accountID := _AccountIDFromProtobuf(aa.AccountID)
		amount := HbarFromTinybar(aa.Amount)

		tx.AddHbarTransfer(*accountID, amount)
	}

	for _, tokenTransfersList := range pb.GetCryptoTransfer().GetTokenTransfers() {
		if tokenID := _TokenIDFromProtobuf(tokenTransfersList.Token); tokenID != nil {
			for _, tokenTransfer := range tokenTransfersList.GetTransfers() {
				if accountID := _AccountIDFromProtobuf(tokenTransfer.AccountID); accountID != nil {
					tx.AddTokenTransfer(*tokenID, *accountID, tokenTransfer.Amount)
				}
			}

			for _, nftTransfer := range tokenTransfersList.GetNftTransfers() {
				senderAccountID := _AccountIDFromProtobuf(nftTransfer.SenderAccountID)
				receiverAccountID := _AccountIDFromProtobuf(nftTransfer.ReceiverAccountID)

				if senderAccountID != nil && receiverAccountID != nil {
					tx.AddNftTransfer(NftID{
						TokenID:      *tokenID,
						SerialNumber: nftTransfer.SerialNumber,
					}, *senderAccountID, *receiverAccountID)
				}
			}
		}
	}

	return tx
}

func (transaction *TransferTransaction) GetNftTransfers() map[TokenID][]TokenNftTransfer {
	nftTransferMap := map[TokenID][]TokenNftTransfer{}

	for _, tokenTransfer := range transaction.tokenTransfers {
		if len(tokenTransfer.nftTransfers) > 0 {
			nftTransferMap[tokenTransfer.tokenID] = tokenTransfer.nftTransfers
		}
	}

	return nftTransferMap
}

func (transaction *TransferTransaction) GetTokenTransfers() map[TokenID][]TokenTransfer {
	tokenTransferMap := make(map[TokenID][]TokenTransfer)

	for _, tokenTransfer := range transaction.tokenTransfers {
		if len(tokenTransfer.transfers) > 0 {
			tokenTransferMap[tokenTransfer.tokenID] = tokenTransfer.transfers
		}
	}

	return tokenTransferMap
}

func (transaction *TransferTransaction) GetHbarTransfers() map[AccountID]Hbar {
	transferMap := make(map[AccountID]Hbar)

	for _, transfer := range transaction.hbarTransfers {
		transferMap[transfer.accountID] = transfer.amount
	}

	return transferMap
}

func (transaction *TransferTransaction) AddHbarTransfer(accountID AccountID, amount Hbar) *TransferTransaction {
	transaction._RequireNotFrozen()

	if index, ok := transaction.hbarTransferMap[accountID]; ok {
		currentAmount := transaction.hbarTransfers[index].amount

		transaction.hbarTransfers[index].amount = HbarFromTinybar(currentAmount.AsTinybar() + amount.AsTinybar())
	} else {
		index := len(transaction.hbarTransfers)

		transaction.hbarTransferMap[accountID] = index
		transaction.hbarTransfers = append(transaction.hbarTransfers, accountAmount{
			accountID: accountID,
			amount:    amount,
		})
	}

	return transaction
}

func (transaction *TransferTransaction) AddTokenTransfer(tokenID TokenID, accountID AccountID, value int64) *TransferTransaction {
	transaction._RequireNotFrozen()

	if tokenIndex, ok := transaction.tokenTransferMap[tokenID]; ok {
		tokenTransfer := transaction.tokenTransfers[tokenIndex]

		if transferIndex, ok := tokenTransfer.transferMap[accountID]; ok {
			tokenTransfer.transfers[transferIndex].Amount += value
		} else {
			transferIndex := len(tokenTransfer.transfers)

			tokenTransfer.transferMap[accountID] = transferIndex
			tokenTransfer.transfers = append(tokenTransfer.transfers, TokenTransfer{
				AccountID: accountID,
				Amount:    value,
			})
		}

		transaction.tokenTransfers[tokenIndex] = tokenTransfer
	} else {
		tokenIndex := len(transaction.tokenTransfers)

		transaction.tokenTransferMap[tokenID] = tokenIndex
		transaction.tokenTransfers = append(transaction.tokenTransfers, tokenTransfer{
			tokenID:      tokenID,
			nftTransfers: []TokenNftTransfer{},
			transferMap:  map[AccountID]int{accountID: 0},
			transfers: []TokenTransfer{{
				AccountID: accountID,
				Amount:    value,
			}},
		})
	}

	return transaction
}

func (transaction *TransferTransaction) AddNftTransfer(nftID NftID, sender AccountID, receiver AccountID) *TransferTransaction {
	transaction._RequireNotFrozen()

	nftTransfer := TokenNftTransfer{
		SenderAccountID:   sender,
		ReceiverAccountID: receiver,
		SerialNumber:      nftID.SerialNumber,
	}

	if tokenIndex, ok := transaction.tokenTransferMap[nftID.TokenID]; ok {
		tokenTransfer := transaction.tokenTransfers[tokenIndex]
		tokenTransfer.nftTransfers = append(tokenTransfer.nftTransfers, nftTransfer)
	} else {
		tokenIndex := len(transaction.tokenTransfers)

		transaction.tokenTransferMap[nftID.TokenID] = tokenIndex
		transaction.tokenTransfers = append(transaction.tokenTransfers, tokenTransfer{
			tokenID:      nftID.TokenID,
			nftTransfers: []TokenNftTransfer{nftTransfer},
			transferMap:  map[AccountID]int{},
			transfers:    []TokenTransfer{},
		})
	}

	return transaction
}

func (transaction *TransferTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	var err error

	for _, tokenTransfer := range transaction.tokenTransfers {
		err = tokenTransfer.tokenID.ValidateChecksum(client)
		if err != nil {
			return err
		}

		for _, transfer := range tokenTransfer.transfers {
			err = transfer.AccountID.ValidateChecksum(client)
			if err != nil {
				return err
			}
		}

		for _, nftTransfer := range tokenTransfer.nftTransfers {
			err = nftTransfer.ReceiverAccountID.ValidateChecksum(client)
			if err != nil {
				return err
			}

			err = nftTransfer.SenderAccountID.ValidateChecksum(client)
			if err != nil {
				return err
			}
		}
	}

	for _, transfer := range transaction.hbarTransfers {
		err = transfer.accountID.ValidateChecksum(client)
		if err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TransferTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TransferTransaction) _ConstructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.CryptoTransferTransactionBody{
		Transfers: &proto.TransferList{
			AccountAmounts: []*proto.AccountAmount{},
		},
		TokenTransfers: []*proto.TokenTransferList{},
	}

	if len(transaction.hbarTransfers) > 0 {
		body.Transfers.AccountAmounts = make([]*proto.AccountAmount, 0)

		for _, transfer := range transaction.hbarTransfers {
			body.Transfers.AccountAmounts = append(body.Transfers.AccountAmounts, &proto.AccountAmount{
				AccountID: transfer.accountID._ToProtobuf(),
				Amount:    transfer.amount.AsTinybar(),
			})
		}
	}

	if len(transaction.tokenTransfers) > 0 {
		if body.TokenTransfers == nil {
			body.TokenTransfers = make([]*proto.TokenTransferList, 0)
		}

		for _, tokenTransfer := range transaction.tokenTransfers {
			transfers := make([]*proto.AccountAmount, 0)
			nftTransfers := make([]*proto.NftTransfer, 0)

			for _, transfer := range tokenTransfer.transfers {
				transfers = append(transfers, &proto.AccountAmount{
					AccountID: transfer.AccountID._ToProtobuf(),
					Amount:    transfer.Amount,
				})
			}

			for _, nftTransfer := range tokenTransfer.nftTransfers {
				nftTransfers = append(nftTransfers, &proto.NftTransfer{
					SenderAccountID:   nftTransfer.SenderAccountID._ToProtobuf(),
					ReceiverAccountID: nftTransfer.ReceiverAccountID._ToProtobuf(),
					SerialNumber:      nftTransfer.SerialNumber,
				})
			}

			body.TokenTransfers = append(body.TokenTransfers, &proto.TokenTransferList{
				Token:        tokenTransfer.tokenID._ToProtobuf(),
				Transfers:    transfers,
				NftTransfers: nftTransfers,
			})
		}
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_CryptoTransfer{
			CryptoTransfer: body,
		},
	}, nil
}

func _TransferTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().CryptoTransfer,
	}
}

func (transaction *TransferTransaction) AddSignature(publicKey PublicKey, signature []byte) *TransferTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
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
			publicKey._ToSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

func (transaction *TransferTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
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
func (transaction *TransferTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TransferTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
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

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := _Execute(
		client,
		_Request{
			transaction: &transaction.Transaction,
		},
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_TransferTransactionGetMethod,
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

func (transaction *TransferTransaction) _Build() *proto.TransactionBody {
	body := &proto.CryptoTransferTransactionBody{
		Transfers: &proto.TransferList{
			AccountAmounts: []*proto.AccountAmount{},
		},
		TokenTransfers: []*proto.TokenTransferList{},
	}

	if len(transaction.hbarTransfers) > 0 {
		body.Transfers.AccountAmounts = make([]*proto.AccountAmount, 0)

		for _, transfer := range transaction.hbarTransfers {
			body.Transfers.AccountAmounts = append(body.Transfers.AccountAmounts, &proto.AccountAmount{
				AccountID: transfer.accountID._ToProtobuf(),
				Amount:    transfer.amount.AsTinybar(),
			})
		}
	}

	if len(transaction.tokenTransfers) > 0 {
		if body.TokenTransfers == nil {
			body.TokenTransfers = make([]*proto.TokenTransferList, 0)
		}

		for _, tokenTransfer := range transaction.tokenTransfers {
			transfers := make([]*proto.AccountAmount, 0)
			nftTransfers := make([]*proto.NftTransfer, 0)

			for _, transfer := range tokenTransfer.transfers {
				transfers = append(transfers, &proto.AccountAmount{
					AccountID: transfer.AccountID._ToProtobuf(),
					Amount:    transfer.Amount,
				})
			}

			for _, nftTransfer := range tokenTransfer.nftTransfers {
				nftTransfers = append(nftTransfers, &proto.NftTransfer{
					SenderAccountID:   nftTransfer.SenderAccountID._ToProtobuf(),
					ReceiverAccountID: nftTransfer.ReceiverAccountID._ToProtobuf(),
					SerialNumber:      nftTransfer.SerialNumber,
				})
			}

			body.TokenTransfers = append(body.TokenTransfers, &proto.TokenTransferList{
				Token:        tokenTransfer.tokenID._ToProtobuf(),
				Transfers:    transfers,
				NftTransfers: nftTransfers,
			})
		}
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &proto.TransactionBody_CryptoTransfer{
			CryptoTransfer: body,
		},
	}
}

func (transaction *TransferTransaction) Freeze() (*TransferTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TransferTransaction) FreezeWith(client *Client) (*TransferTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}

	transaction._InitFee(client)
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TransferTransaction{}, err
	}
	transaction._InitFee(client)
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TransferTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenUpdateTransaction.
func (transaction *TransferTransaction) SetMaxTransactionFee(fee Hbar) *TransferTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TransferTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenUpdateTransaction.
func (transaction *TransferTransaction) SetTransactionMemo(memo string) *TransferTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TransferTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenUpdateTransaction.
func (transaction *TransferTransaction) SetTransactionValidDuration(duration time.Duration) *TransferTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TransferTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenUpdateTransaction.
func (transaction *TransferTransaction) SetTransactionID(transactionID TransactionID) *TransferTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenUpdateTransaction.
func (transaction *TransferTransaction) SetNodeAccountIDs(nodeID []AccountID) *TransferTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TransferTransaction) SetMaxRetry(count int) *TransferTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TransferTransaction) SetMaxBackoff(max time.Duration) *TransferTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TransferTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TransferTransaction) SetMinBackoff(min time.Duration) *TransferTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TransferTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
