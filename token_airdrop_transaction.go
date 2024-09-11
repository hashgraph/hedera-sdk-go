package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"sort"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type TokenAirdropTransaction struct {
	Transaction
	tokenTransfers map[TokenID]*_TokenTransfer
	nftTransfers   map[TokenID][]*_TokenNftTransfer
}

func NewTokenAirdropTransaction() *TokenAirdropTransaction {
	tx := TokenAirdropTransaction{
		Transaction:    _NewTransaction(),
		tokenTransfers: make(map[TokenID]*_TokenTransfer),
		nftTransfers:   make(map[TokenID][]*_TokenNftTransfer),
	}

	tx._SetDefaultMaxTransactionFee(NewHbar(1))

	return &tx
}

func _TokenAirdropTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenAirdropTransaction {
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

	return &TokenAirdropTransaction{
		Transaction:    tx,
		tokenTransfers: tokenTransfers,
		nftTransfers:   nftTransfers,
	}
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
		sort.Sort(tempTokenTransferList)

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

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenAirdropTransaction) Sign(privateKey PrivateKey) *TokenAirdropTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenAirdropTransaction) SignWithOperator(client *Client) (*TokenAirdropTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenAirdropTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenAirdropTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenAirdropTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenAirdropTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenAirdropTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenAirdropTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenAirdropTransaction) Freeze() (*TokenAirdropTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenAirdropTransaction) FreezeWith(client *Client) (*TokenAirdropTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenAirdropTransaction.
func (tx *TokenAirdropTransaction) SetMaxTransactionFee(fee Hbar) *TokenAirdropTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenAirdropTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenAirdropTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenAirdropTransaction.
func (tx *TokenAirdropTransaction) SetTransactionMemo(memo string) *TokenAirdropTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenAirdropTransaction.
func (tx *TokenAirdropTransaction) SetTransactionValidDuration(duration time.Duration) *TokenAirdropTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *TokenAirdropTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this TokenAirdropTransaction.
func (tx *TokenAirdropTransaction) SetTransactionID(transactionID TransactionID) *TokenAirdropTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenAirdropTransaction.
func (tx *TokenAirdropTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenAirdropTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenAirdropTransaction) SetMaxRetry(count int) *TokenAirdropTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenAirdropTransaction) SetMaxBackoff(max time.Duration) *TokenAirdropTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenAirdropTransaction) SetMinBackoff(min time.Duration) *TokenAirdropTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenAirdropTransaction) SetLogLevel(level LogLevel) *TokenAirdropTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenAirdropTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenAirdropTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenAirdropTransaction) getName() string {
	return "TokenAirdropTransaction"
}

func (tx *TokenAirdropTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TokenAirdropTransaction) build() *services.TransactionBody {
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

func (tx *TokenAirdropTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenAirdrop{
			TokenAirdrop: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenAirdropTransaction) buildProtoBody() *services.TokenAirdropTransactionBody {
	body := &services.TokenAirdropTransactionBody{
		TokenTransfers: []*services.TokenTransferList{},
	}

	tempTokenIDarray := make([]TokenID, 0)
	for transfer := range tx.tokenTransfers {
		tempTokenIDarray = append(tempTokenIDarray, transfer)
	}
	sort.Sort(_TokenIDs{tokenIDs: tempTokenIDarray})

	for _, k := range tempTokenIDarray {
		sort.Sort(&_HbarTransfers{transfers: tx.tokenTransfers[k].Transfers})
	}

	if len(tempTokenIDarray) > 0 {
		if body.TokenTransfers == nil {
			body.TokenTransfers = make([]*services.TokenTransferList, 0)
		}

		for _, tokenID := range tempTokenIDarray {
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

	tempTokenIDarray = make([]TokenID, 0)
	for k := range tx.nftTransfers {
		tempTokenIDarray = append(tempTokenIDarray, k)
	}
	sort.Sort(_TokenIDs{tokenIDs: tempTokenIDarray})

	tempNftTransfers := make(map[TokenID][]*_TokenNftTransfer)
	for _, k := range tempTokenIDarray {
		tempTokenNftTransfer := tx.nftTransfers[k]

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

	return body
}

func (tx *TokenAirdropTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().AirdropTokens,
	}
}

func (this *TokenAirdropTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return this.buildScheduled()
}
