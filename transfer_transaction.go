package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TransferTransaction
// Transfers cryptocurrency among two or more accounts by making the desired adjustments to their
// balances. Each transfer list can specify up to 10 adjustments. Each negative amount is withdrawn
// from the corresponding account (a sender), and each positive one is added to the corresponding
// account (a receiver). The amounts list must sum to zero. Each amount is a number of tinybars
// (there are 100,000,000 tinybars in one hbar).  If any sender account fails to have sufficient
// hbars, then the entire transaction fails, and none of those transfers occur, though the
// transaction fee is still charged. This transaction must be signed by the keys for all the sending
// accounts, and for any receiving accounts that have receiverSigRequired == true. The signatures
// are in the same order as the accounts, skipping those accounts that don't need a signature.
type TransferTransaction struct {
	*Transaction[*TransferTransaction]
	tokenTransfers map[TokenID]*_TokenTransfer
	hbarTransfers  []*_HbarTransfer
	nftTransfers   map[TokenID][]*_TokenNftTransfer
}

// NewTransferTransaction creates TransferTransaction which
// transfers cryptocurrency among two or more accounts by making the desired adjustments to their
// balances. Each transfer list can specify up to 10 adjustments. Each negative amount is withdrawn
// from the corresponding account (a sender), and each positive one is added to the corresponding
// account (a receiver). The amounts list must sum to zero. Each amount is a number of tinybars
// (there are 100,000,000 tinybars in one hbar).  If any sender account fails to have sufficient
// hbars, then the entire transaction fails, and none of those transfers occur, though the
// transaction fee is still charged. This transaction must be signed by the keys for all the sending
// accounts, and for any receiving accounts that have receiverSigRequired == true. The signatures
// are in the same order as the accounts, skipping those accounts that don't need a signature.
func NewTransferTransaction() *TransferTransaction {
	tx := &TransferTransaction{
		tokenTransfers: make(map[TokenID]*_TokenTransfer),
		hbarTransfers:  make([]*_HbarTransfer, 0),
		nftTransfers:   make(map[TokenID][]*_TokenNftTransfer),
	}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(1))

	return tx
}

func _TransferTransactionFromProtobuf(tx Transaction[*TransferTransaction], pb *services.TransactionBody) TransferTransaction {
	tokenTransfers := make(map[TokenID]*_TokenTransfer)
	nftTransfers := make(map[TokenID][]*_TokenNftTransfer)

	for _, tokenTransfersList := range pb.GetCryptoTransfer().GetTokenTransfers() {
		tok := _TokenIDFromProtobuf(tokenTransfersList.Token)
		tokenTransfers[*tok] = _TokenTransferPrivateFromProtobuf(tokenTransfersList)
	}

	for _, tokenTransfersList := range pb.GetCryptoTransfer().GetTokenTransfers() {
		if tokenID := _TokenIDFromProtobuf(tokenTransfersList.Token); tokenID != nil {
			for _, aa := range tokenTransfersList.GetNftTransfers() {
				if nftTransfers[*tokenID] == nil {
					nftTransfers[*tokenID] = make([]*_TokenNftTransfer, 0)
				}
				nftTransfer := _NftTransferFromProtobuf(aa)
				nftTransfers[*tokenID] = append(nftTransfers[*tokenID], &nftTransfer)
			}
		}
	}

	transferTransaction := TransferTransaction{
		hbarTransfers:  _HbarTransferFromProtobuf(pb.GetCryptoTransfer().GetTransfers().GetAccountAmounts()),
		tokenTransfers: tokenTransfers,
		nftTransfers:   nftTransfers,
	}

	tx.childTransaction = &transferTransaction
	transferTransaction.Transaction = &tx
	return transferTransaction
}

// SetTokenTransferApproval Sets the desired token unit balance adjustments
func (tx *TransferTransaction) SetTokenTransferApproval(tokenID TokenID, accountID AccountID, approval bool) *TransferTransaction { //nolint
	for token, tokenTransfer := range tx.tokenTransfers {
		if token.Compare(tokenID) == 0 {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID.Compare(accountID) == 0 {
					transfer.IsApproved = approval
				}
			}
		}
	}

	return tx
}

// SetHbarTransferApproval Sets the desired hbar balance adjustments
func (tx *TransferTransaction) SetHbarTransferApproval(spenderAccountID AccountID, approval bool) *TransferTransaction { //nolint
	for _, k := range tx.hbarTransfers {
		if k.accountID.String() == spenderAccountID.String() {
			k.IsApproved = approval
		}
	}
	return tx
}

// SetNftTransferApproval Sets the desired nft token unit balance adjustments
func (tx *TransferTransaction) SetNftTransferApproval(nftID NftID, approval bool) *TransferTransaction {
	for token, nftTransfers := range tx.nftTransfers {
		if token.Compare(nftID.TokenID) == 0 {
			for _, nftTransfer := range nftTransfers {
				if nftTransfer.SerialNumber == nftID.SerialNumber {
					nftTransfer.IsApproved = approval
				}
			}
		}
	}
	return tx
}

// GetNftTransfers returns the nft transfers
func (tx *TransferTransaction) GetNftTransfers() map[TokenID][]_TokenNftTransfer {
	nftResult := make(map[TokenID][]_TokenNftTransfer)
	for token, nftTransfers := range tx.nftTransfers {
		tempArray := make([]_TokenNftTransfer, 0)
		for _, nftTransfer := range nftTransfers {
			tempArray = append(tempArray, *nftTransfer)
		}

		nftResult[token] = tempArray
	}

	return nftResult
}

// GetTokenTransfers returns the token transfers
func (tx *TransferTransaction) GetTokenTransfers() map[TokenID][]TokenTransfer {
	transfers := make(map[TokenID][]TokenTransfer)
	for tokenID, tokenTransfers := range tx.tokenTransfers {
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

		transfers[tokenID] = tempTokenTransferList.transfers
	}

	return transfers
}

// GetHbarTransfers returns the hbar transfers
func (tx *TransferTransaction) GetHbarTransfers() map[AccountID]Hbar {
	result := make(map[AccountID]Hbar)
	for _, hbarTransfers := range tx.hbarTransfers {
		result[*hbarTransfers.accountID] = hbarTransfers.Amount
	}
	return result
}

// AddHbarTransfer Sets The desired hbar balance adjustments
func (tx *TransferTransaction) AddHbarTransfer(accountID AccountID, amount Hbar) *TransferTransaction {
	tx._RequireNotFrozen()

	for _, transfer := range tx.hbarTransfers {
		if transfer.accountID.Compare(accountID) == 0 {
			transfer.Amount = HbarFromTinybar(amount.AsTinybar() + transfer.Amount.AsTinybar())
			return tx
		}
	}

	tx.hbarTransfers = append(tx.hbarTransfers, &_HbarTransfer{
		accountID:  &accountID,
		Amount:     amount,
		IsApproved: false,
	})

	return tx
}

// GetTokenIDDecimals returns the token decimals
func (tx *TransferTransaction) GetTokenIDDecimals() map[TokenID]uint32 {
	result := make(map[TokenID]uint32)
	for token, tokenTransfer := range tx.tokenTransfers {
		if tokenTransfer.ExpectedDecimals != nil {
			result[token] = *tokenTransfer.ExpectedDecimals
		}
	}
	return result
}

// AddTokenTransferWithDecimals Sets the desired token unit balance adjustments with decimals
func (tx *TransferTransaction) AddTokenTransferWithDecimals(tokenID TokenID, accountID AccountID, value int64, decimal uint32) *TransferTransaction { //nolint
	tx._RequireNotFrozen()

	for token, tokenTransfer := range tx.tokenTransfers {
		if token.Compare(tokenID) == 0 {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID.Compare(accountID) == 0 {
					transfer.Amount = HbarFromTinybar(transfer.Amount.AsTinybar() + value)
					tokenTransfer.ExpectedDecimals = &decimal

					return tx
				}
			}
		}
	}

	if v, ok := tx.tokenTransfers[tokenID]; ok {
		v.Transfers = append(v.Transfers, &_HbarTransfer{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: false,
		})
		v.ExpectedDecimals = &decimal

		return tx
	}

	tx.tokenTransfers[tokenID] = &_TokenTransfer{
		Transfers: []*_HbarTransfer{{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: false,
		}},
		ExpectedDecimals: &decimal,
	}

	return tx
}

// AddTokenTransfer Sets the desired token unit balance adjustments
// Applicable to tokens of type FUNGIBLE_COMMON.
func (tx *TransferTransaction) AddTokenTransfer(tokenID TokenID, accountID AccountID, value int64) *TransferTransaction { //nolint
	tx._RequireNotFrozen()

	for token, tokenTransfer := range tx.tokenTransfers {
		if token.Compare(tokenID) == 0 {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID.Compare(accountID) == 0 {
					transfer.Amount = HbarFromTinybar(transfer.Amount.AsTinybar() + value)

					return tx
				}
			}
		}
	}

	if v, ok := tx.tokenTransfers[tokenID]; ok {
		v.Transfers = append(v.Transfers, &_HbarTransfer{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: false,
		})

		return tx
	}

	tx.tokenTransfers[tokenID] = &_TokenTransfer{
		Transfers: []*_HbarTransfer{{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: false,
		}},
	}

	return tx
}

// AddNftTransfer Sets the desired nft token unit balance adjustments
// Applicable to tokens of type NON_FUNGIBLE_UNIQUE.
func (tx *TransferTransaction) AddNftTransfer(nftID NftID, sender AccountID, receiver AccountID) *TransferTransaction {
	tx._RequireNotFrozen()

	if tx.nftTransfers == nil {
		tx.nftTransfers = make(map[TokenID][]*_TokenNftTransfer)
	}

	if tx.nftTransfers[nftID.TokenID] == nil {
		tx.nftTransfers[nftID.TokenID] = make([]*_TokenNftTransfer, 0)
	}

	tx.nftTransfers[nftID.TokenID] = append(tx.nftTransfers[nftID.TokenID], &_TokenNftTransfer{
		SenderAccountID:   sender,
		ReceiverAccountID: receiver,
		SerialNumber:      nftID.SerialNumber,
	})

	return tx
}

// AddHbarTransferWithDecimals adds an approved hbar transfer
func (tx *TransferTransaction) AddApprovedHbarTransfer(accountID AccountID, amount Hbar, approve bool) *TransferTransaction {
	tx._RequireNotFrozen()

	for _, transfer := range tx.hbarTransfers {
		if transfer.accountID.Compare(accountID) == 0 {
			transfer.Amount = HbarFromTinybar(amount.AsTinybar() + transfer.Amount.AsTinybar())
			transfer.IsApproved = approve
			return tx
		}
	}

	tx.hbarTransfers = append(tx.hbarTransfers, &_HbarTransfer{
		accountID:  &accountID,
		Amount:     amount,
		IsApproved: approve,
	})

	return tx
}

// AddHbarTransfer adds an approved hbar transfer with decimals
func (tx *TransferTransaction) AddApprovedTokenTransferWithDecimals(tokenID TokenID, accountID AccountID, value int64, decimal uint32, approve bool) *TransferTransaction { //nolint
	tx._RequireNotFrozen()

	for token, tokenTransfer := range tx.tokenTransfers {
		if token.Compare(tokenID) == 0 {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID.Compare(accountID) == 0 {
					transfer.Amount = HbarFromTinybar(transfer.Amount.AsTinybar() + value)
					tokenTransfer.ExpectedDecimals = &decimal
					for _, transfer := range tokenTransfer.Transfers {
						transfer.IsApproved = approve
					}

					return tx
				}
			}
		}
	}

	if v, ok := tx.tokenTransfers[tokenID]; ok {
		v.Transfers = append(v.Transfers, &_HbarTransfer{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: approve,
		})
		v.ExpectedDecimals = &decimal

		return tx
	}

	tx.tokenTransfers[tokenID] = &_TokenTransfer{
		Transfers: []*_HbarTransfer{{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: approve,
		}},
		ExpectedDecimals: &decimal,
	}

	return tx
}

// AddHbarTransfer adds an approved hbar transfer
func (tx *TransferTransaction) AddApprovedTokenTransfer(tokenID TokenID, accountID AccountID, value int64, approve bool) *TransferTransaction { //nolint
	tx._RequireNotFrozen()

	for token, tokenTransfer := range tx.tokenTransfers {
		if token.Compare(tokenID) == 0 {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID.Compare(accountID) == 0 {
					transfer.Amount = HbarFromTinybar(transfer.Amount.AsTinybar() + value)
					transfer.IsApproved = approve

					return tx
				}
			}
		}
	}

	if v, ok := tx.tokenTransfers[tokenID]; ok {
		v.Transfers = append(v.Transfers, &_HbarTransfer{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: approve,
		})

		return tx
	}

	tx.tokenTransfers[tokenID] = &_TokenTransfer{
		Transfers: []*_HbarTransfer{{
			accountID:  &accountID,
			Amount:     HbarFromTinybar(value),
			IsApproved: approve,
		}},
	}

	return tx
}

// AddNftTransfer adds an approved nft transfer
func (tx *TransferTransaction) AddApprovedNftTransfer(nftID NftID, sender AccountID, receiver AccountID, approve bool) *TransferTransaction {
	tx._RequireNotFrozen()

	if tx.nftTransfers == nil {
		tx.nftTransfers = make(map[TokenID][]*_TokenNftTransfer)
	}

	if tx.nftTransfers[nftID.TokenID] == nil {
		tx.nftTransfers[nftID.TokenID] = make([]*_TokenNftTransfer, 0)
	}

	tx.nftTransfers[nftID.TokenID] = append(tx.nftTransfers[nftID.TokenID], &_TokenNftTransfer{
		SenderAccountID:   sender,
		ReceiverAccountID: receiver,
		SerialNumber:      nftID.SerialNumber,
		IsApproved:        approve,
	})

	return tx
}

// ----------- Overridden functions ----------------

func (tx TransferTransaction) getName() string {
	return "TransferTransaction"
}

func (tx TransferTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	for token, tokenTransfer := range tx.tokenTransfers {
		err = token.ValidateChecksum(client)
		if err != nil {
			return err
		}
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
	for token, nftTransfers := range tx.nftTransfers {
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
	for _, hbarTransfer := range tx.hbarTransfers {
		err = hbarTransfer.accountID.ValidateChecksum(client)
		if err != nil {
			return err
		}
	}

	return nil
}

func (tx TransferTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_CryptoTransfer{
			CryptoTransfer: tx.buildProtoBody(),
		},
	}
}

func (tx TransferTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoTransfer{
			CryptoTransfer: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TransferTransaction) buildProtoBody() *services.CryptoTransferTransactionBody {
	body := &services.CryptoTransferTransactionBody{
		Transfers: &services.TransferList{
			AccountAmounts: []*services.AccountAmount{},
		},
		TokenTransfers: []*services.TokenTransferList{},
	}

	if len(tx.hbarTransfers) > 0 {
		body.Transfers.AccountAmounts = make([]*services.AccountAmount, 0)
		for _, hbarTransfer := range tx.hbarTransfers {
			body.Transfers.AccountAmounts = append(body.Transfers.AccountAmounts, &services.AccountAmount{
				AccountID:  hbarTransfer.accountID._ToProtobuf(),
				Amount:     hbarTransfer.Amount.AsTinybar(),
				IsApproval: hbarTransfer.IsApproved,
			})
		}
	}

	if len(tx.tokenTransfers) > 0 {
		if body.TokenTransfers == nil {
			body.TokenTransfers = make([]*services.TokenTransferList, 0)
		}

		for tokenID := range tx.tokenTransfers {
			transfers := tx.tokenTransfers[tokenID]._ToProtobuf()

			bod := &services.TokenTransferList{
				Token:     tokenID._ToProtobuf(),
				Transfers: transfers,
			}

			if tx.tokenTransfers[tokenID].ExpectedDecimals != nil {
				bod.ExpectedDecimals = &wrapperspb.UInt32Value{Value: *tx.tokenTransfers[tokenID].ExpectedDecimals}
			}

			body.TokenTransfers = append(body.TokenTransfers, bod)
		}
	}

	if len(tx.nftTransfers) > 0 {
		if body.TokenTransfers == nil {
			body.TokenTransfers = make([]*services.TokenTransferList, 0)
		}

		for tokenID, nftTransferList := range tx.nftTransfers {
			nftTransfers := make([]*services.NftTransfer, 0)

			for _, nftT := range nftTransferList {
				nftTransfers = append(nftTransfers, nftT._ToProtobuf())
			}

			body.TokenTransfers = append(body.TokenTransfers, &services.TokenTransferList{
				Token:        tokenID._ToProtobuf(),
				NftTransfers: nftTransfers,
			})
		}
	}

	return body
}

func (tx TransferTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().CryptoTransfer,
	}
}

func (tx TransferTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TransferTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, tx)
}
