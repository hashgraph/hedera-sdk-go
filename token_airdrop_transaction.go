package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

type TokenAirdropTransaction struct {
	*Transaction[*TokenAirdropTransaction]
	tokenTransfers map[TokenID]*_TokenTransfer
	nftTransfers   map[TokenID][]*_TokenNftTransfer
}

func NewTokenAirdropTransaction() *TokenAirdropTransaction {
	tx := &TokenAirdropTransaction{
		tokenTransfers: make(map[TokenID]*_TokenTransfer),
		nftTransfers:   make(map[TokenID][]*_TokenNftTransfer),
	}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(1))

	return tx
}

func _TokenAirdropTransactionFromProtobuf(tx Transaction[*TokenAirdropTransaction], pb *services.TransactionBody) TokenAirdropTransaction {
	tokenTransfers := make(map[TokenID]*_TokenTransfer)
	nftTransfers := make(map[TokenID][]*_TokenNftTransfer)

	for _, tokenTransfersList := range pb.GetTokenAirdrop().GetTokenTransfers() {
		tok := _TokenIDFromProtobuf(tokenTransfersList.Token)
		tokenTransfers[*tok] = _TokenTransferPrivateFromProtobuf(tokenTransfersList)
	}

	for _, tokenTransfersList := range pb.GetTokenAirdrop().GetTokenTransfers() {
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

	tokenAirdropTransaction := TokenAirdropTransaction{
		tokenTransfers: tokenTransfers,
		nftTransfers:   nftTransfers,
	}

	tx.childTransaction = &tokenAirdropTransaction
	tokenAirdropTransaction.Transaction = &tx
	return tokenAirdropTransaction
}

// SetTokenTransferApproval Sets the desired token unit balance adjustments
func (tx *TokenAirdropTransaction) SetTokenTransferApproval(tokenID TokenID, accountID AccountID, approval bool) *TokenAirdropTransaction { //nolint
	for token, tokenTransfer := range tx.tokenTransfers {
		if token.equals(tokenID) {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID._Equals(accountID) {
					transfer.IsApproved = approval
				}
			}
		}
	}

	return tx
}

// SetNftTransferApproval Sets the desired nft token unit balance adjustments
func (tx *TokenAirdropTransaction) SetNftTransferApproval(nftID NftID, approval bool) *TokenAirdropTransaction {
	for token, nftTransfers := range tx.nftTransfers {
		if token.equals(nftID.TokenID) {
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
func (tx *TokenAirdropTransaction) GetNftTransfers() map[TokenID][]_TokenNftTransfer {
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
func (tx *TokenAirdropTransaction) GetTokenTransfers() map[TokenID][]TokenTransfer {
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

// GetTokenIDDecimals returns the token decimals
func (tx *TokenAirdropTransaction) GetTokenIDDecimals() map[TokenID]uint32 {
	result := make(map[TokenID]uint32)
	for token, tokenTransfer := range tx.tokenTransfers {
		if tokenTransfer.ExpectedDecimals != nil {
			result[token] = *tokenTransfer.ExpectedDecimals
		}
	}
	return result
}

// AddTokenTransferWithDecimals Sets the desired token unit balance adjustments with decimals
func (tx *TokenAirdropTransaction) AddTokenTransferWithDecimals(tokenID TokenID, accountID AccountID, value int64, decimal uint32) *TokenAirdropTransaction { //nolint
	tx._RequireNotFrozen()

	for token, tokenTransfer := range tx.tokenTransfers {
		if token.equals(tokenID) {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID._Equals(accountID) {
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
func (tx *TokenAirdropTransaction) AddTokenTransfer(tokenID TokenID, accountID AccountID, value int64) *TokenAirdropTransaction { //nolint
	tx._RequireNotFrozen()

	for token, tokenTransfer := range tx.tokenTransfers {
		if token.equals(tokenID) {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID._Equals(accountID) {
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
func (tx *TokenAirdropTransaction) AddNftTransfer(nftID NftID, sender AccountID, receiver AccountID) *TokenAirdropTransaction {
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

// AddApprovedTokenTransferWithDecimals adds an approved token transfer with decimals
func (tx *TokenAirdropTransaction) AddApprovedTokenTransferWithDecimals(tokenID TokenID, accountID AccountID, value int64, decimal uint32, approve bool) *TokenAirdropTransaction { //nolint
	tx._RequireNotFrozen()

	for token, tokenTransfer := range tx.tokenTransfers {
		if token.equals(tokenID) {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID._Equals(accountID) {
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

// AddApprovedTokenTransfer adds an approved token transfer
func (tx *TokenAirdropTransaction) AddApprovedTokenTransfer(tokenID TokenID, accountID AccountID, value int64, approve bool) *TokenAirdropTransaction { //nolint
	tx._RequireNotFrozen()

	for token, tokenTransfer := range tx.tokenTransfers {
		if token.equals(tokenID) {
			for _, transfer := range tokenTransfer.Transfers {
				if transfer.accountID._Equals(accountID) {
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

// AddApprovedNftTransfer adds an approved nft transfer
func (tx *TokenAirdropTransaction) AddApprovedNftTransfer(nftID NftID, sender AccountID, receiver AccountID, approve bool) *TokenAirdropTransaction {
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

func (tx TokenAirdropTransaction) getName() string {
	return "TokenAirdropTransaction"
}

func (tx TokenAirdropTransaction) validateNetworkOnIDs(client *Client) error {
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

	return nil
}

func (tx TokenAirdropTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenAirdrop{
			TokenAirdrop: tx.buildProtoBody(),
		},
	}
}

func (tx TokenAirdropTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenAirdrop{
			TokenAirdrop: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenAirdropTransaction) buildProtoBody() *services.TokenAirdropTransactionBody {
	body := &services.TokenAirdropTransactionBody{
		TokenTransfers: []*services.TokenTransferList{},
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

func (tx TokenAirdropTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().AirdropTokens,
	}
}

func (tx TokenAirdropTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenAirdropTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
