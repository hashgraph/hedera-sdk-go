package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TokenMintTransaction
// Mints tokens from the Token's treasury Account. If no Supply Key is defined, the transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
// The operation decreases the Total Supply of the Token. Total supply cannot go below zero.
// The amount provided must be in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to mint 100 tokens, one must provide amount of 10000. In order
// to mint 100.55 tokens, one must provide amount of 10055.
type TokenMintTransaction struct {
	*Transaction[*TokenMintTransaction]
	tokenID *TokenID
	amount  uint64
	meta    [][]byte
}

// NewTokenMintTransaction creates TokenMintTransaction which
// mints tokens from the Token's treasury Account. If no Supply Key is defined, the transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
// The operation decreases the Total Supply of the Token. Total supply cannot go below zero.
// The amount provided must be in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to mint 100 tokens, one must provide amount of 10000. In order
// to mint 100.55 tokens, one must provide amount of 10055.
func NewTokenMintTransaction() *TokenMintTransaction {
	tx := &TokenMintTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return tx
}

func _TokenMintTransactionFromProtobuf(tx Transaction[*TokenMintTransaction], pb *services.TransactionBody) TokenMintTransaction {
	tokenMintTransaction := TokenMintTransaction{
		tokenID: _TokenIDFromProtobuf(pb.GetTokenMint().GetToken()),
		amount:  pb.GetTokenMint().GetAmount(),
		meta:    pb.GetTokenMint().GetMetadata(),
	}

	tx.childTransaction = &tokenMintTransaction
	tokenMintTransaction.Transaction = &tx
	return tokenMintTransaction
}

// SetTokenID Sets the token for which to mint tokens. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (tx *TokenMintTransaction) SetTokenID(tokenID TokenID) *TokenMintTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the TokenID for this TokenMintTransaction
func (tx *TokenMintTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetAmount Sets the amount to mint from the Treasury Account. Amount must be a positive non-zero number, not
// bigger than the token balance of the treasury account (0; balance], represented in the lowest
// denomination.
func (tx *TokenMintTransaction) SetAmount(amount uint64) *TokenMintTransaction {
	tx._RequireNotFrozen()
	tx.amount = amount
	return tx
}

// GetAmount returns the amount to mint from the Treasury Account
func (tx *TokenMintTransaction) GetAmount() uint64 {
	return tx.amount
}

// SetMetadatas
// Applicable to tokens of type NON_FUNGIBLE_UNIQUE. A list of metadata that are being created.
// Maximum allowed size of each metadata is 100 bytes
func (tx *TokenMintTransaction) SetMetadatas(meta [][]byte) *TokenMintTransaction {
	tx._RequireNotFrozen()
	tx.meta = meta
	return tx
}

// SetMetadata
// Applicable to tokens of type NON_FUNGIBLE_UNIQUE. A list of metadata that are being created.
// Maximum allowed size of each metadata is 100 bytes
func (tx *TokenMintTransaction) SetMetadata(meta []byte) *TokenMintTransaction {
	tx._RequireNotFrozen()
	if tx.meta == nil {
		tx.meta = make([][]byte, 0)
	}
	tx.meta = append(tx.meta, [][]byte{meta}...)
	return tx
}

// GetMetadatas returns the metadata that are being created.
func (tx *TokenMintTransaction) GetMetadatas() [][]byte {
	return tx.meta
}

// ----------- Overridden functions ----------------

func (tx TokenMintTransaction) getName() string {
	return "TokenMintTransaction"
}

func (tx TokenMintTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.tokenID != nil {
		if err := tx.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx TokenMintTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenMint{
			TokenMint: tx.buildProtoBody(),
		},
	}
}

func (tx TokenMintTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenMint{
			TokenMint: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenMintTransaction) buildProtoBody() *services.TokenMintTransactionBody {
	body := &services.TokenMintTransactionBody{
		Amount: tx.amount,
	}

	if tx.meta != nil {
		body.Metadata = tx.meta
	}

	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	return body
}

func (tx TokenMintTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().MintToken,
	}
}

func (tx TokenMintTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenMintTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
