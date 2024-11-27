package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TokenUpdateTransaction
// At consensus, updates an already created token to the given values.
//
// If no value is given for a field, that field is left unchanged. For an immutable tokens (that is,
// a token without an admin key), only the expiry may be updated. Setting any other field in that
// case will cause the transaction status to resolve to TOKEN_IS_IMMUTABLE.
//
// --- Signing Requirements ---
//  1. Whether or not a token has an admin key, its expiry can be extended with only the transaction
//     payer's signature.
//  2. Updating any other field of a mutable token requires the admin key's signature.
//  3. If a new admin key is set, this new key must sign <b>unless</b> it is exactly an empty
//     <tt>KeyList</tt>. This special sentinel key removes the existing admin key and causes the
//     token to become immutable. (Other <tt>Key</tt> structures without a constituent
//     <tt>Ed25519</tt> key will be rejected with <tt>INVALID_ADMIN_KEY</tt>.)
//  4. If a new treasury is set, the new treasury account's key must sign the transaction.
//
// --- Nft Requirements ---
//  1. If a non fungible token has a positive treasury balance, the operation will abort with
//     CURRENT_TREASURY_STILL_OWNS_NFTS.
type TokenUpdateTransaction struct {
	*Transaction[*TokenUpdateTransaction]
	tokenID                  *TokenID
	treasuryAccountID        *AccountID
	autoRenewAccountID       *AccountID
	tokenName                string
	memo                     *string
	metadata                 []byte
	tokenSymbol              string
	adminKey                 Key
	kycKey                   Key
	freezeKey                Key
	wipeKey                  Key
	scheduleKey              Key
	supplyKey                Key
	pauseKey                 Key
	metadataKey              Key
	tokenKeyVerificationMode TokenKeyValidation
	expirationTime           *time.Time
	autoRenewPeriod          *time.Duration
}

// NewTokenUpdateTransaction creates TokenUpdateTransaction which at consensus,
// updates an already created token to the given values.
//
// If no value is given for a field, that field is left unchanged. For an immutable tokens (that is,
// a token without an admin key), only the expiry may be updated. Setting any other field in that
// case will cause the transaction status to resolve to TOKEN_IS_IMMUTABLE.
//
// --- Signing Requirements ---
//  1. Whether or not a token has an admin key, its expiry can be extended with only the transaction
//     payer's signature.
//  2. Updating any other field of a mutable token requires the admin key's signature.
//  3. If a new admin key is set, this new key must sign <b>unless</b> it is exactly an empty
//     <tt>KeyList</tt>. This special sentinel key removes the existing admin key and causes the
//     token to become immutable. (Other <tt>Key</tt> structures without a constituent
//     <tt>Ed25519</tt> key will be rejected with <tt>INVALID_ADMIN_KEY</tt>.)
//  4. If a new treasury is set, the new treasury account's key must sign the transaction.
//
// --- Nft Requirements ---
//  1. If a non fungible token has a positive treasury balance, the operation will abort with
//     CURRENT_TREASURY_STILL_OWNS_NFTS.
func NewTokenUpdateTransaction() *TokenUpdateTransaction {
	tx := TokenUpdateTransaction{
		memo: nil,
	}

	tx.Transaction = _NewTransaction(&tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return &tx
}

func _TokenUpdateTransactionFromProtobuf(tx Transaction[*TokenUpdateTransaction], pb *services.TransactionBody) TokenUpdateTransaction {
	adminKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetAdminKey())
	kycKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetKycKey())
	freezeKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetFreezeKey())
	wipeKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetWipeKey())
	scheduleKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetFeeScheduleKey())
	supplyKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetSupplyKey())
	pauseKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetPauseKey())
	metadataKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetMetadataKey())
	keyVerificationMode := pb.GetTokenUpdate().GetKeyVerificationMode()

	expirationTime := _TimeFromProtobuf(pb.GetTokenUpdate().GetExpiry())
	autoRenew := _DurationFromProtobuf(pb.GetTokenUpdate().GetAutoRenewPeriod())

	var memo *string
	if m := pb.GetTokenUpdate().GetMemo(); m != nil {
		memo = &m.Value
	}

	var metadata []byte
	if m := pb.GetTokenUpdate().GetMetadata(); m != nil {
		metadata = m.Value
	}

	tokenUpdateTransaction := TokenUpdateTransaction{
		tokenID:                  _TokenIDFromProtobuf(pb.GetTokenUpdate().GetToken()),
		treasuryAccountID:        _AccountIDFromProtobuf(pb.GetTokenUpdate().GetTreasury()),
		autoRenewAccountID:       _AccountIDFromProtobuf(pb.GetTokenUpdate().GetAutoRenewAccount()),
		tokenName:                pb.GetTokenUpdate().GetName(),
		memo:                     memo,
		metadata:                 metadata,
		tokenSymbol:              pb.GetTokenUpdate().GetSymbol(),
		adminKey:                 adminKey,
		kycKey:                   kycKey,
		freezeKey:                freezeKey,
		wipeKey:                  wipeKey,
		scheduleKey:              scheduleKey,
		supplyKey:                supplyKey,
		pauseKey:                 pauseKey,
		metadataKey:              metadataKey,
		tokenKeyVerificationMode: TokenKeyValidation(keyVerificationMode),
		expirationTime:           &expirationTime,
		autoRenewPeriod:          &autoRenew,
	}

	tx.childTransaction = &tokenUpdateTransaction
	tokenUpdateTransaction.Transaction = &tx
	return tokenUpdateTransaction
}

type TokenKeyValidation int32

const (
	// FULL_VALIDATION performs all token key validations.
	FULL_VALIDATION TokenKeyValidation = iota
	// NO_VALIDATION performs no validations at all for all passed token keys.
	NO_VALIDATION
)

// SetTokenID Sets the Token to be updated
func (tx *TokenUpdateTransaction) SetTokenID(tokenID TokenID) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the TokenID for this TokenUpdateTransaction
func (tx *TokenUpdateTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetTokenSymbol Sets the new Symbol of the Token.
// Must be UTF-8 capitalized alphabetical string identifying the token.
func (tx *TokenUpdateTransaction) SetTokenSymbol(symbol string) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.tokenSymbol = symbol
	return tx
}

func (tx *TokenUpdateTransaction) GetTokenSymbol() string {
	return tx.tokenSymbol
}

// SetTokenName Sets the new Name of the Token. Must be a string of ASCII characters.
func (tx *TokenUpdateTransaction) SetTokenName(name string) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.tokenName = name
	return tx
}

// GetTokenName returns the TokenName for this TokenUpdateTransaction
func (tx *TokenUpdateTransaction) GetTokenName() string {
	return tx.tokenName
}

// SetTreasuryAccountID sets thehe new Treasury account of the Token.
// If the provided treasury account is not existing or deleted,
// the _Response will be INVALID_TREASURY_ACCOUNT_FOR_TOKEN. If successful, the Token
// balance held in the previous Treasury Account is transferred to the new one.
func (tx *TokenUpdateTransaction) SetTreasuryAccountID(treasuryAccountID AccountID) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.treasuryAccountID = &treasuryAccountID
	return tx
}

func (tx *TokenUpdateTransaction) GetTreasuryAccountID() AccountID {
	if tx.treasuryAccountID == nil {
		return AccountID{}
	}

	return *tx.treasuryAccountID
}

// SetAdminKey Sets the new Admin key of the Token.
// If Token is immutable, transaction will resolve to TOKEN_IS_IMMUTABlE.
func (tx *TokenUpdateTransaction) SetAdminKey(publicKey Key) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = publicKey
	return tx
}

func (tx *TokenUpdateTransaction) GetAdminKey() Key {
	return tx.adminKey
}

// SetPauseKey Sets the Key which can pause and unpause the Token. If the Token does not currently have a pause key,
// transaction will resolve to TOKEN_HAS_NO_PAUSE_KEY
func (tx *TokenUpdateTransaction) SetPauseKey(publicKey Key) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.pauseKey = publicKey
	return tx
}

// GetPauseKey returns the Key which can pause and unpause the Token
func (tx *TokenUpdateTransaction) GetPauseKey() Key {
	return tx.pauseKey
}

// SetMetadataKey Set the Key which can update the metadata.
func (tx *TokenUpdateTransaction) SetMetadataKey(key Key) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.metadataKey = key
	return tx
}

// GetMetadataKey returns the metadata key
func (tx *TokenUpdateTransaction) GetMetadataKey() Key {
	return tx.metadataKey
}

// SetKycKey Sets the new KYC key of the Token. If Token does not have currently a KYC key, transaction will
// resolve to TOKEN_HAS_NO_KYC_KEY.
func (tx *TokenUpdateTransaction) SetKycKey(publicKey Key) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.kycKey = publicKey
	return tx
}

// GetKycKey returns the new KYC key of the Token
func (tx *TokenUpdateTransaction) GetKycKey() Key {
	return tx.kycKey
}

// SetFreezeKey Sets the new Freeze key of the Token. If the Token does not have currently a Freeze key, transaction
// will resolve to TOKEN_HAS_NO_FREEZE_KEY.
func (tx *TokenUpdateTransaction) SetFreezeKey(publicKey Key) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.freezeKey = publicKey
	return tx
}

// GetFreezeKey returns the new Freeze key of the Token
func (tx *TokenUpdateTransaction) GetFreezeKey() Key {
	return tx.freezeKey
}

// SetWipeKey Sets the new Wipe key of the Token. If the Token does not have currently a Wipe key, transaction
// will resolve to TOKEN_HAS_NO_WIPE_KEY.
func (tx *TokenUpdateTransaction) SetWipeKey(publicKey Key) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.wipeKey = publicKey
	return tx
}

// GetWipeKey returns the new Wipe key of the Token
func (tx *TokenUpdateTransaction) GetWipeKey() Key {
	return tx.wipeKey
}

// SetSupplyKey Sets the new Supply key of the Token. If the Token does not have currently a Supply key, transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
func (tx *TokenUpdateTransaction) SetSupplyKey(publicKey Key) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.supplyKey = publicKey
	return tx
}

// GetSupplyKey returns the new Supply key of the Token
func (tx *TokenUpdateTransaction) GetSupplyKey() Key {
	return tx.supplyKey
}

// SetFeeScheduleKey
// If set, the new key to use to update the token's custom fee schedule; if the token does not
// currently have this key, transaction will resolve to TOKEN_HAS_NO_FEE_SCHEDULE_KEY
func (tx *TokenUpdateTransaction) SetFeeScheduleKey(key Key) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.scheduleKey = key
	return tx
}

// GetFeeScheduleKey returns the new key to use to update the token's custom fee schedule
func (tx *TokenUpdateTransaction) GetFeeScheduleKey() Key {
	return tx.scheduleKey
}

// SetAutoRenewAccount Sets the new account which will be automatically charged to renew the token's expiration, at
// autoRenewPeriod interval.
func (tx *TokenUpdateTransaction) SetAutoRenewAccount(autoRenewAccountID AccountID) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewAccountID = &autoRenewAccountID
	return tx
}

func (tx *TokenUpdateTransaction) GetAutoRenewAccount() AccountID {
	if tx.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *tx.autoRenewAccountID
}

// SetAutoRenewPeriod Sets the new interval at which the auto-renew account will be charged to extend the token's expiry.
func (tx *TokenUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewPeriod = &autoRenewPeriod
	return tx
}

// GetAutoRenewPeriod returns the new interval at which the auto-renew account will be charged to extend the token's expiry.
func (tx *TokenUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if tx.autoRenewPeriod != nil {
		return time.Duration(int64(tx.autoRenewPeriod.Seconds()) * time.Second.Nanoseconds())
	}

	return time.Duration(0)
}

// SetExpirationTime Sets the new expiry time of the token. Expiry can be updated even if admin key is not set.
// If the provided expiry is earlier than the current token expiry, transaction wil resolve to
// INVALID_EXPIRATION_TIME
func (tx *TokenUpdateTransaction) SetExpirationTime(expirationTime time.Time) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.expirationTime = &expirationTime
	return tx
}

func (tx *TokenUpdateTransaction) GetExpirationTime() time.Time {
	if tx.expirationTime != nil {
		return *tx.expirationTime
	}
	return time.Time{}
}

// SetTokenMemo
// If set, the new memo to be associated with the token (UTF-8 encoding max 100 bytes)
func (tx *TokenUpdateTransaction) SetTokenMemo(memo string) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.memo = &memo

	return tx
}

func (tx *TokenUpdateTransaction) GetTokenMemo() string {
	return *tx.memo
}

// SetTokenMetadata sets the token metadata
func (tx *TokenUpdateTransaction) SetTokenMetadata(metadata []byte) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.metadata = metadata

	return tx
}

// GetTokenMetadata returns the token metadata
func (tx *TokenUpdateTransaction) GetTokenMetadata() []byte {
	return tx.metadata
}

// SetKeyVerificationMode sets the token key verification mode
func (tx *TokenUpdateTransaction) SetKeyVerificationMode(verificationMode TokenKeyValidation) *TokenUpdateTransaction {
	tx._RequireNotFrozen()
	tx.tokenKeyVerificationMode = verificationMode
	return tx
}

// GetKeyVerificationMode returns the token metadata
func (tx *TokenUpdateTransaction) GetKeyVerificationMode() TokenKeyValidation {
	return tx.tokenKeyVerificationMode
}

// ----------- Overridden functions ----------------

func (tx TokenUpdateTransaction) getName() string {
	return "TokenUpdateTransaction"
}

func (tx TokenUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.tokenID != nil {
		if err := tx.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.treasuryAccountID != nil {
		if err := tx.treasuryAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if tx.autoRenewAccountID != nil {
		if err := tx.autoRenewAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx TokenUpdateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenUpdate{
			TokenUpdate: tx.buildProtoBody(),
		},
	}
}

func (tx TokenUpdateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenUpdate{
			TokenUpdate: tx.buildProtoBody(),
		},
	}, nil
}

func (tx TokenUpdateTransaction) buildProtoBody() *services.TokenUpdateTransactionBody {
	body := &services.TokenUpdateTransactionBody{
		Name:                tx.tokenName,
		Symbol:              tx.tokenSymbol,
		KeyVerificationMode: services.TokenKeyValidation(tx.tokenKeyVerificationMode),
	}

	if tx.memo != nil {
		body.Memo = &wrapperspb.StringValue{Value: *tx.memo}
	}

	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	if tx.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*tx.autoRenewPeriod)
	}

	if tx.expirationTime != nil {
		body.Expiry = _TimeToProtobuf(*tx.expirationTime)
	}

	if tx.treasuryAccountID != nil {
		body.Treasury = tx.treasuryAccountID._ToProtobuf()
	}

	if tx.autoRenewAccountID != nil {
		body.AutoRenewAccount = tx.autoRenewAccountID._ToProtobuf()
	}

	if tx.adminKey != nil {
		body.AdminKey = tx.adminKey._ToProtoKey()
	}

	if tx.freezeKey != nil {
		body.FreezeKey = tx.freezeKey._ToProtoKey()
	}

	if tx.scheduleKey != nil {
		body.FeeScheduleKey = tx.scheduleKey._ToProtoKey()
	}

	if tx.kycKey != nil {
		body.KycKey = tx.kycKey._ToProtoKey()
	}

	if tx.wipeKey != nil {
		body.WipeKey = tx.wipeKey._ToProtoKey()
	}

	if tx.supplyKey != nil {
		body.SupplyKey = tx.supplyKey._ToProtoKey()
	}

	if tx.pauseKey != nil {
		body.PauseKey = tx.pauseKey._ToProtoKey()
	}

	if tx.metadataKey != nil {
		body.MetadataKey = tx.metadataKey._ToProtoKey()
	}

	if tx.metadata != nil {
		body.Metadata = &wrapperspb.BytesValue{Value: tx.metadata}
	}

	return body
}

func (tx TokenUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().UpdateToken,
	}
}

func (tx TokenUpdateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx TokenUpdateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
