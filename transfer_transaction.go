package hedera

import (
	"sort"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TransferTransaction struct {
	Transaction
	tokenTransfers   map[TokenID]map[AccountID]int64
	hbarTransfers    map[AccountID]Hbar
	nftTransfers     map[TokenID][]TokenNftTransfer
	expectedDecimals map[TokenID]uint32
}

func NewTransferTransaction() *TransferTransaction {
	transaction := TransferTransaction{
		Transaction:      _NewTransaction(),
		tokenTransfers:   make(map[TokenID]map[AccountID]int64),
		hbarTransfers:    make(map[AccountID]Hbar),
		nftTransfers:     make(map[TokenID][]TokenNftTransfer),
		expectedDecimals: make(map[TokenID]uint32),
		hbarApprovals:    make(map[AccountID]bool),
		tokenApprovals:   make(map[TokenID]map[AccountID]bool),
	}

	transaction.SetMaxTransactionFee(NewHbar(1))

	return &transaction
}

func _TransferTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TransferTransaction {
	hbarTransfers := make(map[AccountID]Hbar)
	tokenTransfers := make(map[TokenID]map[AccountID]int64)
	nftTransfers := make(map[TokenID][]TokenNftTransfer)

	for _, aa := range pb.GetCryptoTransfer().GetTransfers().AccountAmounts {
		accountID := _AccountIDFromProtobuf(aa.AccountID)
		amount := HbarFromTinybar(aa.Amount)

		if value, ok := hbarTransfers[*accountID]; ok {
			hbarTransfers[*accountID] = HbarFromTinybar(amount.AsTinybar() + value.AsTinybar())
		} else {
			hbarTransfers[*accountID] = amount
		}
	}

	for _, tokenTransfersList := range pb.GetCryptoTransfer().GetTokenTransfers() {
		if tokenID := _TokenIDFromProtobuf(tokenTransfersList.Token); tokenID != nil {
			var currentTokenTransfers map[AccountID]int64

			if value, ok := tokenTransfers[*tokenID]; ok {
				currentTokenTransfers = value
			} else {
				currentTokenTransfers = make(map[AccountID]int64)
			}

			for _, aa := range tokenTransfersList.GetTransfers() {
				if accountID := _AccountIDFromProtobuf(aa.AccountID); accountID != nil {
					if value, ok := currentTokenTransfers[*accountID]; ok {
						currentTokenTransfers[*accountID] = aa.Amount + value
					} else {
						currentTokenTransfers[*accountID] = aa.Amount
					}
				}
			}

			tokenTransfers[*tokenID] = currentTokenTransfers
		}
	}

	for _, tokenTransfersList := range pb.GetCryptoTransfer().GetTokenTransfers() {
		if tokenID := _TokenIDFromProtobuf(tokenTransfersList.Token); tokenID != nil {
			for _, aa := range tokenTransfersList.GetNftTransfers() {
				if nftTransfers[*tokenID] == nil {
					nftTransfers[*tokenID] = make([]TokenNftTransfer, 0)
				}
				nftTransfers[*tokenID] = append(nftTransfers[*tokenID], _NftTransferFromProtobuf(aa))
			}
		}
	}

	return &TransferTransaction{
		Transaction:    transaction,
		hbarTransfers:  hbarTransfers,
		tokenTransfers: tokenTransfers,
		nftTransfers:   nftTransfers,
	}
}

func (transaction *TransferTransaction) GetNftTransfers() map[TokenID][]TokenNftTransfer {
	return transaction.nftTransfers
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

		tempTokenTransferList := _TokenTransfers{tokenTransfersList}
		sort.Sort(tempTokenTransferList)

		transfers[tokenID] = tempTokenTransferList.transfers
	}

	return transfers
}

func (transaction *TransferTransaction) GetHbarTransfers() map[AccountID]Hbar {
	return transaction.hbarTransfers
}

func (transaction *TransferTransaction) AddHbarTransfer(accountID AccountID, amount Hbar) *TransferTransaction {
	transaction._RequireNotFrozen()

	if value, ok := transaction.hbarTransfers[accountID]; ok {
		transaction.hbarTransfers[accountID] = HbarFromTinybar(amount.AsTinybar() + value.AsTinybar())
	} else {
		transaction.hbarTransfers[accountID] = amount
	}

	return transaction
}

func (transaction *TransferTransaction) GetTokenIDDecimals() map[TokenID]uint32 {
	return transaction.expectedDecimals
}

func (transaction *TransferTransaction) AddTokenTransferWithDecimals(tokenID TokenID, accountID AccountID, value int64, decimal uint32) *TransferTransaction {
	transaction._RequireNotFrozen()

	var tokenTransfers map[AccountID]int64
	var amount int64

	if value, ok := transaction.tokenTransfers[tokenID]; ok {
		tokenTransfers = value
	} else {
		tokenTransfers = make(map[AccountID]int64)
	}

	transaction.expectedDecimals[tokenID] = decimal

	if transfer, ok := tokenTransfers[accountID]; ok {
		amount = transfer + value
	} else {
		amount = value
	}

	tokenTransfers[accountID] = amount
	transaction.tokenTransfers[tokenID] = tokenTransfers

	return transaction
}

func (transaction *TransferTransaction) AddTokenTransfer(tokenID TokenID, accountID AccountID, value int64) *TransferTransaction {
	transaction._RequireNotFrozen()

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
	transaction._RequireNotFrozen()

	if transaction.nftTransfers == nil {
		transaction.nftTransfers = make(map[TokenID][]TokenNftTransfer)
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

func (transaction *TransferTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	for tokenID, accountMap := range transaction.tokenTransfers {
		err = tokenID.ValidateChecksum(client)
		for accountID := range accountMap {
			err = accountID.ValidateChecksum(client)
		}
	}
	for nftID := range transaction.nftTransfers {
		err = nftID.ValidateChecksum(client)
	}
	for accountID := range transaction.hbarTransfers {
		err = accountID.ValidateChecksum(client)
	}
	if err != nil {
		return err
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

func (transaction *TransferTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.CryptoTransferTransactionBody{
		Transfers: &services.TransferList{
			AccountAmounts: []*services.AccountAmount{},
		},
		TokenTransfers: []*services.TokenTransferList{},
	}

	tempAccountIDarray := make([]AccountID, 0)
	for k := range transaction.hbarTransfers {
		tempAccountIDarray = append(tempAccountIDarray, k)
	}
	sort.Sort(_AccountIDs{accountIDs: tempAccountIDarray})

	if len(tempAccountIDarray) > 0 {
		body.Transfers.AccountAmounts = make([]*services.AccountAmount, 0)
		for _, accountID := range tempAccountIDarray {
			body.Transfers.AccountAmounts = append(body.Transfers.AccountAmounts, &services.AccountAmount{
				AccountID: accountID._ToProtobuf(),
				Amount:    transaction.hbarTransfers[accountID].AsTinybar(),
			})
		}
	}

	tempTokenIDarray := make([]TokenID, 0)
	for k := range transaction.tokenTransfers {
		tempTokenIDarray = append(tempTokenIDarray, k)
	}
	sort.Sort(_TokenIDs{tokenIDs: tempTokenIDarray})

	tempTokenTransfers := make(map[TokenID][]AccountID)
	for _, k := range tempTokenIDarray {
		initialAccountMap := transaction.tokenTransfers[k]

		tempAccountIDarray2 := make([]AccountID, 0)
		for k2 := range initialAccountMap {
			tempAccountIDarray2 = append(tempAccountIDarray2, k2)
		}
		sort.Sort(_AccountIDs{accountIDs: tempAccountIDarray2})

		tempTokenTransfers[k] = tempAccountIDarray2
	}

	if len(tempTokenIDarray) > 0 {
		if body.TokenTransfers == nil {
			body.TokenTransfers = make([]*services.TokenTransferList, 0)
		}

		for _, tokenID := range tempTokenIDarray {
			transfers := make([]*services.AccountAmount, 0)

			for _, accountID := range tempTokenTransfers[tokenID] {
				temp := transaction.tokenTransfers[tokenID]
				transfers = append(transfers, &services.AccountAmount{
					AccountID: accountID._ToProtobuf(),
					Amount:    temp[accountID],
				})
			}

			bod := &services.TokenTransferList{
				Token:     tokenID._ToProtobuf(),
				Transfers: transfers,
			}

			if decimal, ok := transaction.expectedDecimals[tokenID]; ok {
				bod.ExpectedDecimals = &wrapperspb.UInt32Value{Value: decimal}
			}

			body.TokenTransfers = append(body.TokenTransfers, bod)
		}
	}

	tempTokenIDarray = make([]TokenID, 0)
	for k := range transaction.nftTransfers {
		tempTokenIDarray = append(tempTokenIDarray, k)
	}
	sort.Sort(_TokenIDs{tokenIDs: tempTokenIDarray})

	tempNftTransfers := make(map[TokenID][]TokenNftTransfer)
	for _, k := range tempTokenIDarray {
		tempTokenNftTransfer := transaction.nftTransfers[k]

		sort.Sort(_TokenNftTransfers{tempTokenNftTransfer})

		tempNftTransfers[k] = tempTokenNftTransfer
	}

	if len(tempTokenIDarray) > 0 {
		if body.TokenTransfers == nil {
			body.TokenTransfers = make([]*services.TokenTransferList, 0)
		}

		for _, tokenID := range tempTokenIDarray {
			nftTransfers := make([]*services.NftTransfer, 0)

			for _, nftT := range tempNftTransfers[tokenID] {
				nftTransfers = append(nftTransfers, nftT._ToProtobuf())
			}

			body.TokenTransfers = append(body.TokenTransfers, &services.TokenTransferList{
				Token:        tokenID._ToProtobuf(),
				NftTransfers: nftTransfers,
			})
		}
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoTransfer{
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

	transaction.transactions = make([]*services.Transaction, 0)
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

func (transaction *TransferTransaction) _Build() *services.TransactionBody {
	body := &services.CryptoTransferTransactionBody{
		Transfers: &services.TransferList{
			AccountAmounts: []*services.AccountAmount{},
		},
		TokenTransfers: []*services.TokenTransferList{},
	}

	tempAccountIDarray := make([]AccountID, 0)
	for k := range transaction.hbarTransfers {
		tempAccountIDarray = append(tempAccountIDarray, k)
	}
	sort.Sort(_AccountIDs{accountIDs: tempAccountIDarray})

	if len(tempAccountIDarray) > 0 {
		body.Transfers.AccountAmounts = make([]*services.AccountAmount, 0)
		for _, accountID := range tempAccountIDarray {
			body.Transfers.AccountAmounts = append(body.Transfers.AccountAmounts, &services.AccountAmount{
				AccountID: accountID._ToProtobuf(),
				Amount:    transaction.hbarTransfers[accountID].AsTinybar(),
			})
		}
	}

	tempTokenIDarray := make([]TokenID, 0)
	for k := range transaction.tokenTransfers {
		tempTokenIDarray = append(tempTokenIDarray, k)
	}
	sort.Sort(_TokenIDs{tokenIDs: tempTokenIDarray})

	tempTokenTransfers := make(map[TokenID][]AccountID)
	for _, tokenID := range tempTokenIDarray {
		initialAccountMap := transaction.tokenTransfers[tokenID]

		tempAccountIDarray2 := make([]AccountID, 0)
		for k2 := range initialAccountMap {
			tempAccountIDarray2 = append(tempAccountIDarray2, k2)
		}
		sort.Sort(_AccountIDs{accountIDs: tempAccountIDarray2})

		tempTokenTransfers[tokenID] = tempAccountIDarray2
	}

	if len(tempTokenIDarray) > 0 {
		if body.TokenTransfers == nil {
			body.TokenTransfers = make([]*services.TokenTransferList, 0)
		}

		for _, tokenID := range tempTokenIDarray {
			transfers := make([]*services.AccountAmount, 0)

			for _, accountID := range tempTokenTransfers[tokenID] {
				temp := transaction.tokenTransfers[tokenID]
				transfers = append(transfers, &services.AccountAmount{
					AccountID: accountID._ToProtobuf(),
					Amount:    temp[accountID],
				})
			}

			bod := &services.TokenTransferList{
				Token:     tokenID._ToProtobuf(),
				Transfers: transfers,
			}

			if decimal, ok := transaction.expectedDecimals[tokenID]; ok {
				bod.ExpectedDecimals = &wrapperspb.UInt32Value{Value: decimal}
			}

			body.TokenTransfers = append(body.TokenTransfers, bod)
		}
	}

	tempTokenIDarray = make([]TokenID, 0)
	for k := range transaction.nftTransfers {
		tempTokenIDarray = append(tempTokenIDarray, k)
	}
	sort.Sort(_TokenIDs{tokenIDs: tempTokenIDarray})

	tempNftTransfers := make(map[TokenID][]TokenNftTransfer)
	for _, k := range tempTokenIDarray {
		tempTokenNftTransfer := transaction.nftTransfers[k]

		sort.Sort(_TokenNftTransfers{tempTokenNftTransfer})

		tempNftTransfers[k] = tempTokenNftTransfer
	}

	if len(tempTokenIDarray) > 0 {
		if body.TokenTransfers == nil {
			body.TokenTransfers = make([]*services.TokenTransferList, 0)
		}

		for _, tokenID := range tempTokenIDarray {
			nftTransfers := make([]*services.NftTransfer, 0)

			for _, nftT := range tempNftTransfers[tokenID] {
				nftTransfers = append(nftTransfers, nftT._ToProtobuf())
			}

			body.TokenTransfers = append(body.TokenTransfers, &services.TokenTransferList{
				Token:        tokenID._ToProtobuf(),
				NftTransfers: nftTransfers,
			})
		}
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoTransfer{
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
