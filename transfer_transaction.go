package hedera

import (
	"fmt"
	"sort"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TransferTransaction struct {
	Transaction
	tokenTransfers map[TokenID]*_TokenTransfer
	hbarTransfers  []*_HbarTransfer
	nftTransfers   map[TokenID][]*TokenNftTransfer
}

func NewTransferTransaction() *TransferTransaction {
	transaction := TransferTransaction{
		Transaction:    _NewTransaction(),
		tokenTransfers: make(map[TokenID]*_TokenTransfer),
		hbarTransfers:  make([]*_HbarTransfer, 0),
		nftTransfers:   make(map[TokenID][]*TokenNftTransfer),
	}

	transaction.SetMaxTransactionFee(NewHbar(1))

	return &transaction
}

func _TransferTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TransferTransaction {
	tokenTransfers := make(map[TokenID]*_TokenTransfer)
	nftTransfers := make(map[TokenID][]*TokenNftTransfer)

	for _, tokenTransfersList := range pb.GetCryptoTransfer().GetTokenTransfers() {
		tok := _TokenIDFromProtobuf(tokenTransfersList.Token)
		tokenTransfers[*tok] = _TokenTransferPrivateFromProtobuf(tokenTransfersList)
	}

	for _, tokenTransfersList := range pb.GetCryptoTransfer().GetTokenTransfers() {
		if tokenID := _TokenIDFromProtobuf(tokenTransfersList.Token); tokenID != nil {
			for _, aa := range tokenTransfersList.GetNftTransfers() {
				if nftTransfers[*tokenID] == nil {
					nftTransfers[*tokenID] = make([]*TokenNftTransfer, 0)
				}
				nftTransfer := _NftTransferFromProtobuf(aa)
				nftTransfers[*tokenID] = append(nftTransfers[*tokenID], &nftTransfer)
			}
		}
	}

	return &TransferTransaction{
		Transaction:    transaction,
		hbarTransfers:  _HbarTransferFromProtobuf(pb.GetCryptoTransfer().GetTransfers().GetAccountAmounts()),
		tokenTransfers: tokenTransfers,
		nftTransfers:   nftTransfers,
	}
}

func (transaction *TransferTransaction) SetGrpcDeadline(deadline *time.Duration) *TransferTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

func (transaction *TransferTransaction) SetTokenTransferApproval(tokenID TokenID, accountID AccountID, approval bool) *TransferTransaction { //nolint
	for token, tokenTransfer := range transaction.tokenTransfers {
		if token.Compare(tokenID) == 0 {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID.Compare(accountID) == 0 {
					transfer.IsApproved = approval
				}
			}
		}
	}

	return transaction
}

func (transaction *TransferTransaction) SetHbarTransferApproval(spenderAccountID AccountID, approval bool) *TransferTransaction { //nolint
	for _, k := range transaction.hbarTransfers {
		if k.accountID.String() == spenderAccountID.String() {
			k.IsApproved = approval
		}
	}
	return transaction
}

func (transaction *TransferTransaction) SetNftTransferApproval(nftID NftID, approval bool) *TransferTransaction {
	for token, nftTransfers := range transaction.nftTransfers {
		if token.Compare(nftID.TokenID) == 0 {
			for _, nftTransfer := range nftTransfers {
				if nftTransfer.SerialNumber == nftID.SerialNumber {
					nftTransfer.IsApproved = approval
				}
			}
		}
	}
	return transaction
}

func (transaction *TransferTransaction) GetNftTransfers() map[TokenID][]TokenNftTransfer {
	nftResult := make(map[TokenID][]TokenNftTransfer)
	for token, nftTransfers := range transaction.nftTransfers {
		tempArray := make([]TokenNftTransfer, 0)
		for _, nftTransfer := range nftTransfers {
			tempArray = append(tempArray, *nftTransfer)
		}

		nftResult[token] = tempArray
	}

	return nftResult
}

func (transaction *TransferTransaction) GetTokenTransfers() map[TokenID][]TokenTransfer {
	transfers := make(map[TokenID][]TokenTransfer)
	for tokenID, tokenTransfers := range transaction.tokenTransfers {
		tokenTransfersList := make([]TokenTransfer, 0)

		for _, transfer := range tokenTransfers.Transfers {
			var acc AccountID
			if transfer.accountID != nil {
				acc = *transfer.accountID
			}
			tokenTransfersList = append(tokenTransfersList, TokenTransfer{
				AccountID:  acc,
				Amount:     transfer.Amount.AsTinybar(),
				IsApproved: transfer.IsApproved,
			})
		}

		tempTokenTransferList := _TokenTransfers{tokenTransfersList}
		sort.Sort(tempTokenTransferList)

		transfers[tokenID] = tempTokenTransferList.transfers
	}

	return transfers
}

func (transaction *TransferTransaction) GetHbarTransfers() map[AccountID]Hbar {
	result := make(map[AccountID]Hbar)
	for _, hbarTransfers := range transaction.hbarTransfers {
		result[*hbarTransfers.accountID] = hbarTransfers.Amount
	}
	return result
}

func (transaction *TransferTransaction) AddHbarTransfer(accountID AccountID, amount Hbar) *TransferTransaction {
	transaction._RequireNotFrozen()

	for _, transfer := range transaction.hbarTransfers {
		if transfer.accountID.Compare(accountID) == 0 {
			transfer.Amount = HbarFromTinybar(amount.AsTinybar() + transfer.Amount.AsTinybar())
			return transaction
		}
	}

	transaction.hbarTransfers = append(transaction.hbarTransfers, &_HbarTransfer{
		accountID:  &accountID,
		Amount:     amount,
		IsApproved: false,
	})

	return transaction
}

func (transaction *TransferTransaction) GetTokenIDDecimals() map[TokenID]uint32 {
	result := make(map[TokenID]uint32)
	for token, tokenTransfer := range transaction.tokenTransfers {
		if tokenTransfer.ExpectedDecimals != nil {
			result[token] = *tokenTransfer.ExpectedDecimals
		}
	}
	return result
}

func (transaction *TransferTransaction) AddTokenTransferWithDecimals(tokenID TokenID, accountID AccountID, value int64, decimal uint32) *TransferTransaction { //nolint
	transaction._RequireNotFrozen()

	for token, tokenTransfer := range transaction.tokenTransfers {
		if token.Compare(tokenID) == 0 {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID.Compare(accountID) == 0 {
					transfer.Amount = HbarFromTinybar(transfer.Amount.AsTinybar() + value)
					tokenTransfer.ExpectedDecimals = &decimal

					return transaction
				}
			}
		}
	}

	if v, ok := transaction.tokenTransfers[tokenID]; ok {
		v.Transfers = append(v.Transfers, &_HbarTransfer{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: false,
		})
		v.ExpectedDecimals = &decimal

		return transaction
	}

	transaction.tokenTransfers[tokenID] = &_TokenTransfer{
		Transfers: []*_HbarTransfer{{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: false,
		}},
		ExpectedDecimals: &decimal,
	}

	return transaction
}

func (transaction *TransferTransaction) AddTokenTransfer(tokenID TokenID, accountID AccountID, value int64) *TransferTransaction { //nolint
	transaction._RequireNotFrozen()

	for token, tokenTransfer := range transaction.tokenTransfers {
		if token.Compare(tokenID) == 0 {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID.Compare(accountID) == 0 {
					transfer.Amount = HbarFromTinybar(transfer.Amount.AsTinybar() + value)

					return transaction
				}
			}
		}
	}

	if v, ok := transaction.tokenTransfers[tokenID]; ok {
		v.Transfers = append(v.Transfers, &_HbarTransfer{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: false,
		})

		return transaction
	}

	transaction.tokenTransfers[tokenID] = &_TokenTransfer{
		Transfers: []*_HbarTransfer{{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: false,
		}},
	}

	return transaction
}

func (transaction *TransferTransaction) AddNftTransfer(nftID NftID, sender AccountID, receiver AccountID) *TransferTransaction {
	transaction._RequireNotFrozen()

	if transaction.nftTransfers == nil {
		transaction.nftTransfers = make(map[TokenID][]*TokenNftTransfer)
	}

	if transaction.nftTransfers[nftID.TokenID] == nil {
		transaction.nftTransfers[nftID.TokenID] = make([]*TokenNftTransfer, 0)
	}

	transaction.nftTransfers[nftID.TokenID] = append(transaction.nftTransfers[nftID.TokenID], &TokenNftTransfer{
		SenderAccountID:   sender,
		ReceiverAccountID: receiver,
		SerialNumber:      nftID.SerialNumber,
	})

	return transaction
}

func (transaction *TransferTransaction) AddApprovedHbarTransfer(accountID AccountID, amount Hbar, approve bool) *TransferTransaction {
	transaction._RequireNotFrozen()

	for _, transfer := range transaction.hbarTransfers {
		if transfer.accountID.Compare(accountID) == 0 {
			transfer.Amount = HbarFromTinybar(amount.AsTinybar() + transfer.Amount.AsTinybar())
			transfer.IsApproved = approve
			return transaction
		}
	}

	transaction.hbarTransfers = append(transaction.hbarTransfers, &_HbarTransfer{
		accountID:  &accountID,
		Amount:     amount,
		IsApproved: approve,
	})

	return transaction
}

func (transaction *TransferTransaction) AddApprovedTokenTransferWithDecimals(tokenID TokenID, accountID AccountID, value int64, decimal uint32, approve bool) *TransferTransaction { //nolint
	transaction._RequireNotFrozen()

	for token, tokenTransfer := range transaction.tokenTransfers {
		if token.Compare(tokenID) == 0 {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID.Compare(accountID) == 0 {
					transfer.Amount = HbarFromTinybar(transfer.Amount.AsTinybar() + value)
					tokenTransfer.ExpectedDecimals = &decimal
					for _, transfer := range tokenTransfer.Transfers {
						transfer.IsApproved = approve
					}

					return transaction
				}
			}
		}
	}

	if v, ok := transaction.tokenTransfers[tokenID]; ok {
		v.Transfers = append(v.Transfers, &_HbarTransfer{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: approve,
		})
		v.ExpectedDecimals = &decimal

		return transaction
	}

	transaction.tokenTransfers[tokenID] = &_TokenTransfer{
		Transfers: []*_HbarTransfer{{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: approve,
		}},
		ExpectedDecimals: &decimal,
	}

	return transaction
}

func (transaction *TransferTransaction) AddApprovedTokenTransfer(tokenID TokenID, accountID AccountID, value int64, approve bool) *TransferTransaction { //nolint
	transaction._RequireNotFrozen()

	for token, tokenTransfer := range transaction.tokenTransfers {
		if token.Compare(tokenID) == 0 {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID.Compare(accountID) == 0 {
					transfer.Amount = HbarFromTinybar(transfer.Amount.AsTinybar() + value)
					transfer.IsApproved = approve

					return transaction
				}
			}
		}
	}

	if v, ok := transaction.tokenTransfers[tokenID]; ok {
		v.Transfers = append(v.Transfers, &_HbarTransfer{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: approve,
		})

		return transaction
	}

	transaction.tokenTransfers[tokenID] = &_TokenTransfer{
		Transfers: []*_HbarTransfer{{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: approve,
		}},
	}

	return transaction
}

func (transaction *TransferTransaction) AddApprovedNftTransfer(nftID NftID, sender AccountID, receiver AccountID, approve bool) *TransferTransaction {
	transaction._RequireNotFrozen()

	if transaction.nftTransfers == nil {
		transaction.nftTransfers = make(map[TokenID][]*TokenNftTransfer)
	}

	if transaction.nftTransfers[nftID.TokenID] == nil {
		transaction.nftTransfers[nftID.TokenID] = make([]*TokenNftTransfer, 0)
	}

	transaction.nftTransfers[nftID.TokenID] = append(transaction.nftTransfers[nftID.TokenID], &TokenNftTransfer{
		SenderAccountID:   sender,
		ReceiverAccountID: receiver,
		SerialNumber:      nftID.SerialNumber,
		IsApproved:        approve,
	})

	return transaction
}

func (transaction *TransferTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	for token, tokenTransfer := range transaction.tokenTransfers {
		err = token.ValidateChecksum(client)
		for _, transfer := range tokenTransfer.Transfers {
			err = transfer.accountID.ValidateChecksum(client)
			if err != nil {
				return err
			}
		}
		if err != nil {
			return err
		}
	}
	for token, nftTransfers := range transaction.nftTransfers {
		err = token.ValidateChecksum(client)
		if err != nil {
			return err
		}
		for _, nftTransfer := range nftTransfers {
			err = nftTransfer.SenderAccountID.ValidateChecksum(client)
			if err != nil {
				return err
			}
			err = nftTransfer.ReceiverAccountID.ValidateChecksum(client)
			if err != nil {
				return err
			}
		}
	}
	for _, hbarTransfer := range transaction.hbarTransfers {
		err = hbarTransfer.accountID.ValidateChecksum(client)
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

func (transaction *TransferTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.CryptoTransferTransactionBody{
		Transfers: &services.TransferList{
			AccountAmounts: []*services.AccountAmount{},
		},
		TokenTransfers: []*services.TokenTransferList{},
	}

	sort.Sort(&_HbarTransfers{transaction.hbarTransfers})

	if len(transaction.hbarTransfers) > 0 {
		body.Transfers.AccountAmounts = make([]*services.AccountAmount, 0)
		for _, hbarTransfer := range transaction.hbarTransfers {
			body.Transfers.AccountAmounts = append(body.Transfers.AccountAmounts, &services.AccountAmount{
				AccountID:  hbarTransfer.accountID._ToProtobuf(),
				Amount:     hbarTransfer.Amount.AsTinybar(),
				IsApproval: hbarTransfer.IsApproved,
			})
		}
	}

	tempTokenIDarray := make([]TokenID, 0)
	for k := range transaction.tokenTransfers {
		tempTokenIDarray = append(tempTokenIDarray, k)
	}
	sort.Sort(_TokenIDs{tokenIDs: tempTokenIDarray})

	for _, k := range tempTokenIDarray {
		sort.Sort(&_HbarTransfers{transfers: transaction.tokenTransfers[k].Transfers})
	}

	if len(tempTokenIDarray) > 0 {
		if body.TokenTransfers == nil {
			body.TokenTransfers = make([]*services.TokenTransferList, 0)
		}

		for _, tokenID := range tempTokenIDarray {
			transfers := transaction.tokenTransfers[tokenID]._ToProtobuf()

			bod := &services.TokenTransferList{
				Token:     tokenID._ToProtobuf(),
				Transfers: transfers,
			}

			if transaction.tokenTransfers[tokenID].ExpectedDecimals != nil {
				bod.ExpectedDecimals = &wrapperspb.UInt32Value{Value: *transaction.tokenTransfers[tokenID].ExpectedDecimals}
			}

			body.TokenTransfers = append(body.TokenTransfers, bod)
		}
	}

	tempTokenIDarray = make([]TokenID, 0)
	for k := range transaction.nftTransfers {
		tempTokenIDarray = append(tempTokenIDarray, k)
	}
	sort.Sort(_TokenIDs{tokenIDs: tempTokenIDarray})

	tempNftTransfers := make(map[TokenID][]*TokenNftTransfer)
	for _, k := range tempTokenIDarray {
		tempTokenNftTransfer := transaction.nftTransfers[k]

		sort.Sort(&_TokenNftTransfers{tempTokenNftTransfer})

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

	if transaction.signedTransactions._Length() == 0 {
		return transaction
	}

	transaction.transactions = _NewLockableSlice()
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)
	transaction.transactionIDs.locked = true

	for index := 0; index < transaction.signedTransactions._Length(); index++ {
		var temp *services.SignedTransaction
		switch t := transaction.signedTransactions._Get(index).(type) { //nolint
		case *services.SignedTransaction:
			temp = t
		}
		temp.SigMap.SigPair = append(
			temp.SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
		transaction.signedTransactions._Set(index, temp)
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

	transactionID := transaction.transactionIDs._GetCurrent().(TransactionID)

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
		transaction._GetLogID(),
		transaction.grpcDeadline,
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

	sort.Sort(&_HbarTransfers{transaction.hbarTransfers})

	if len(transaction.hbarTransfers) > 0 {
		body.Transfers.AccountAmounts = make([]*services.AccountAmount, 0)
		for _, hbarTransfer := range transaction.hbarTransfers {
			body.Transfers.AccountAmounts = append(body.Transfers.AccountAmounts, &services.AccountAmount{
				AccountID:  hbarTransfer.accountID._ToProtobuf(),
				Amount:     hbarTransfer.Amount.AsTinybar(),
				IsApproval: hbarTransfer.IsApproved,
			})
		}
	}

	tempTokenIDarray := make([]TokenID, 0)
	for k := range transaction.tokenTransfers {
		tempTokenIDarray = append(tempTokenIDarray, k)
	}
	sort.Sort(_TokenIDs{tokenIDs: tempTokenIDarray})

	for _, k := range tempTokenIDarray {
		sort.Sort(&_HbarTransfers{transfers: transaction.tokenTransfers[k].Transfers})
	}

	if len(tempTokenIDarray) > 0 {
		if body.TokenTransfers == nil {
			body.TokenTransfers = make([]*services.TokenTransferList, 0)
		}

		for _, tokenID := range tempTokenIDarray {
			transfers := transaction.tokenTransfers[tokenID]._ToProtobuf()

			bod := &services.TokenTransferList{
				Token:     tokenID._ToProtobuf(),
				Transfers: transfers,
			}

			if transaction.tokenTransfers[tokenID].ExpectedDecimals != nil {
				bod.ExpectedDecimals = &wrapperspb.UInt32Value{Value: *transaction.tokenTransfers[tokenID].ExpectedDecimals}
			}

			body.TokenTransfers = append(body.TokenTransfers, bod)
		}
	}

	tempTokenIDarray = make([]TokenID, 0)
	for k := range transaction.nftTransfers {
		tempTokenIDarray = append(tempTokenIDarray, k)
	}
	sort.Sort(_TokenIDs{tokenIDs: tempTokenIDarray})

	tempNftTransfers := make(map[TokenID][]*TokenNftTransfer)
	for _, k := range tempTokenIDarray {
		tempTokenNftTransfer := transaction.nftTransfers[k]

		sort.Sort(&_TokenNftTransfers{tempTokenNftTransfer})

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

// SetMaxTransactionFee sets the max transaction fee for this TransferTransaction.
func (transaction *TransferTransaction) SetMaxTransactionFee(fee Hbar) *TransferTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *TransferTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TransferTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *TransferTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

func (transaction *TransferTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TransferTransaction.
func (transaction *TransferTransaction) SetTransactionMemo(memo string) *TransferTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TransferTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TransferTransaction.
func (transaction *TransferTransaction) SetTransactionValidDuration(duration time.Duration) *TransferTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TransferTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TransferTransaction.
func (transaction *TransferTransaction) SetTransactionID(transactionID TransactionID) *TransferTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TransferTransaction.
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

func (transaction *TransferTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("TransferTransaction:%d", timestamp.UnixNano())
}
