package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TokenGrantKycTransaction
// Grants KYC to the account for the given token. Must be signed by the Token's kycKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no KYC Key is defined, the transaction will resolve to TOKEN_HAS_NO_KYC_KEY.
// Once executed the Account is marked as KYC Granted.
type TokenGrantKycTransaction struct {
	*Transaction[*TokenGrantKycTransaction]
	tokenID   *TokenID
	accountID *AccountID
}

// NewTokenGrantKycTransaction creates TokenGrantKycTransaction which
// grants KYC to the account for the given token. Must be signed by the Token's kycKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no KYC Key is defined, the transaction will resolve to TOKEN_HAS_NO_KYC_KEY.
// Once executed the Account is marked as KYC Granted.
func NewTokenGrantKycTransaction() *TokenGrantKycTransaction {
	tx := &TokenGrantKycTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return tx
}

func _TokenGrantKycTransactionFromProtobuf(tx Transaction[*TokenGrantKycTransaction], pb *services.TransactionBody) TokenGrantKycTransaction {
	tokenGrandKycTransaction := TokenGrantKycTransaction{
		tokenID:   _TokenIDFromProtobuf(pb.GetTokenGrantKyc().GetToken()),
		accountID: _AccountIDFromProtobuf(pb.GetTokenGrantKyc().GetAccount()),
	}

	tx.childTransaction = &tokenGrandKycTransaction
	tokenGrandKycTransaction.Transaction = &tx
	return tokenGrandKycTransaction
}

// SetTokenID Sets the token for which this account will be granted KYC.
// If token does not exist, transaction results in INVALID_TOKEN_ID
func (tx *TokenGrantKycTransaction) SetTokenID(tokenID TokenID) *TokenGrantKycTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the token for which this account will be granted KYC.
func (tx *TokenGrantKycTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetAccountID Sets the account to be KYCed
func (tx *TokenGrantKycTransaction) SetAccountID(accountID AccountID) *TokenGrantKycTransaction {
	tx._RequireNotFrozen()
	tx.accountID = &accountID
	return tx
}

// GetAccountID returns the AccountID that is being KYCed
func (tx *TokenGrantKycTransaction) GetAccountID() AccountID {
	if tx.accountID == nil {
		return AccountID{}
	}

	return *tx.accountID
}

// ----------- Overridden functions ----------------

func (tx TokenGrantKycTransaction) getName() string {
	return "TokenGrantKycTransaction"
}

func (tx TokenGrantKycTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.tokenID != nil {
		if err := tx.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.accountID != nil {
		if err := tx.accountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx TokenGrantKycTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenGrantKyc{
			TokenGrantKyc: tx.buildProtoBody(),
		},
	}
}

func (tx TokenGrantKycTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenGrantKyc{
			TokenGrantKyc: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenGrantKycTransaction) buildProtoBody() *services.TokenGrantKycTransactionBody {
	body := &services.TokenGrantKycTransactionBody{}
	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	if tx.accountID != nil {
		body.Account = tx.accountID._ToProtobuf()
	}

	return body
}

func (tx TokenGrantKycTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().GrantKycToTokenAccount,
	}
}

func (tx TokenGrantKycTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenGrantKycTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
