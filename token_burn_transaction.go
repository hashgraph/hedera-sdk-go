package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TokenBurnTransaction Burns tokens from the Token's treasury Account.
// If no Supply Key is defined, the transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
// The operation decreases the Total Supply of the Token. Total supply cannot go below
// zero.
// The amount provided must be in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to burn 100 tokens, one must provide amount of 10000. In order
// to burn 100.55 tokens, one must provide amount of 10055.
type TokenBurnTransaction struct {
	*Transaction[*TokenBurnTransaction]
	tokenID *TokenID
	amount  uint64
	serial  []int64
}

// NewTokenBurnTransaction creates TokenBurnTransaction which burns tokens from the Token's treasury Account.
// If no Supply Key is defined, the transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
// The operation decreases the Total Supply of the Token. Total supply cannot go below
// zero.
// The amount provided must be in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to burn 100 tokens, one must provide amount of 10000. In order
// to burn 100.55 tokens, one must provide amount of 10055.
func NewTokenBurnTransaction() *TokenBurnTransaction {
	tx := &TokenBurnTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _TokenBurnTransactionFromProtobuf(tx Transaction[*TokenBurnTransaction], pb *services.TransactionBody) TokenBurnTransaction {
	tokenBurnTransaction := TokenBurnTransaction{
		tokenID: _TokenIDFromProtobuf(pb.GetTokenBurn().Token),
		amount:  pb.GetTokenBurn().GetAmount(),
		serial:  pb.GetTokenBurn().GetSerialNumbers(),
	}

	tx.childTransaction = &tokenBurnTransaction
	tokenBurnTransaction.Transaction = &tx
	return tokenBurnTransaction
}

// SetTokenID Sets the token for which to burn tokens. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (tx *TokenBurnTransaction) SetTokenID(tokenID TokenID) *TokenBurnTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the TokenID for the token which will be burned.
func (tx *TokenBurnTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetAmount Sets the amount to burn from the Treasury Account. Amount must be a positive non-zero number, not
// bigger than the token balance of the treasury account (0; balance], represented in the lowest
// denomination.
func (tx *TokenBurnTransaction) SetAmount(amount uint64) *TokenBurnTransaction {
	tx._RequireNotFrozen()
	tx.amount = amount
	return tx
}

// Deprecated: Use TokenBurnTransaction.GetAmount() instead.
func (tx *TokenBurnTransaction) GetAmmount() uint64 {
	return tx.amount
}

func (tx *TokenBurnTransaction) GetAmount() uint64 {
	return tx.amount
}

// SetSerialNumber
// Applicable to tokens of type NON_FUNGIBLE_UNIQUE.
// The list of serial numbers to be burned.
func (tx *TokenBurnTransaction) SetSerialNumber(serial int64) *TokenBurnTransaction {
	tx._RequireNotFrozen()
	if tx.serial == nil {
		tx.serial = make([]int64, 0)
	}
	tx.serial = append(tx.serial, serial)
	return tx
}

// SetSerialNumbers sets the list of serial numbers to be burned.
func (tx *TokenBurnTransaction) SetSerialNumbers(serial []int64) *TokenBurnTransaction {
	tx._RequireNotFrozen()
	tx.serial = serial
	return tx
}

// GetSerialNumbers returns the list of serial numbers to be burned.
func (tx *TokenBurnTransaction) GetSerialNumbers() []int64 {
	return tx.serial
}

// ----------- Overridden functions ----------------

func (tx TokenBurnTransaction) getName() string {
	return "TokenBurnTransaction"
}

func (tx TokenBurnTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx TokenBurnTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenBurn{
			TokenBurn: tx.buildProtoBody(),
		},
	}
}

func (tx TokenBurnTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenBurn{
			TokenBurn: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenBurnTransaction) buildProtoBody() *services.TokenBurnTransactionBody {
	body := &services.TokenBurnTransactionBody{
		Amount: tx.amount,
	}

	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	if tx.serial != nil {
		body.SerialNumbers = tx.serial
	}

	return body
}

func (tx TokenBurnTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().BurnToken,
	}
}

func (tx TokenBurnTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenBurnTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
