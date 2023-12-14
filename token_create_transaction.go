package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// TokenCreateTransaction
// Create a new token. After the token is created, the Token ID for it is in the receipt.
// The specified Treasury Account is receiving the initial supply of tokens as-well as the tokens
// from the Token Mint operation once executed. The balance of the treasury account is decreased
// when the Token Burn operation is executed.
//
// The initialSupply is the initial supply of the smallest parts of a token (like a
// tinybar, not an hbar). These are the smallest units of the token which may be transferred.
//
// The supply can change over time. If the total supply at some moment is S parts of tokens,
// and the token is using D decimals, then S must be less than or equal to
// 2<sup>63</sup>-1, which is 9,223,372,036,854,775,807. The number of whole tokens (not parts) will
// be S / 10<sup>D</sup>.
//
// If decimals is 8 or 11, then the number of whole tokens can be at most a few billions or
// millions, respectively. For example, it could match Bitcoin (21 million whole tokens with 8
// decimals) or hbars (50 billion whole tokens with 8 decimals). It could even match Bitcoin with
// milli-satoshis (21 million whole tokens with 11 decimals).
//
// Note that a created token is immutable if the adminKey is omitted. No property of
// an immutable token can ever change, with the sole exception of its expiry. Anyone can pay to
// extend the expiry time of an immutable token.
//
// A token can be either FUNGIBLE_COMMON or NON_FUNGIBLE_UNIQUE, based on its
// TokenType. If it has been omitted, FUNGIBLE_COMMON type is used.
//
// A token can have either INFINITE or FINITE supply type, based on its
// TokenType. If it has been omitted, INFINITE type is used.
//
// If a FUNGIBLE TokenType is used, initialSupply should explicitly be set to a
// non-negative. If not, the transaction will resolve to INVALID_TOKEN_INITIAL_SUPPLY.
//
// If a NON_FUNGIBLE_UNIQUE TokenType is used, initialSupply should explicitly be set
// to 0. If not, the transaction will resolve to INVALID_TOKEN_INITIAL_SUPPLY.
//
// If an INFINITE TokenSupplyType is used, maxSupply should explicitly be set to 0. If
// it is not 0, the transaction will resolve to INVALID_TOKEN_MAX_SUPPLY.
//
// If a FINITE TokenSupplyType is used, maxSupply should be explicitly set to a
// non-negative value. If it is not, the transaction will resolve to INVALID_TOKEN_MAX_SUPPLY.
type TokenCreateTransaction struct {
	Transaction
	treasuryAccountID  *AccountID
	autoRenewAccountID *AccountID
	customFees         []Fee
	tokenName          string
	memo               string
	tokenSymbol        string
	decimals           uint32
	tokenSupplyType    TokenSupplyType
	tokenType          TokenType
	maxSupply          int64
	adminKey           Key
	kycKey             Key
	freezeKey          Key
	wipeKey            Key
	scheduleKey        Key
	supplyKey          Key
	pauseKey           Key
	initialSupply      uint64
	freezeDefault      *bool
	expirationTime     *time.Time
	autoRenewPeriod    *time.Duration
}

// NewTokenCreateTransaction creates TokenCreateTransaction which creates a new token.
// After the token is created, the Token ID for it is in the receipt.
// The specified Treasury Account is receiving the initial supply of tokens as-well as the tokens
// from the Token Mint operation once executed. The balance of the treasury account is decreased
// when the Token Burn operation is executed.
//
// The initialSupply is the initial supply of the smallest parts of a token (like a
// tinybar, not an hbar). These are the smallest units of the token which may be transferred.
//
// The supply can change over time. If the total supply at some moment is S parts of tokens,
// and the token is using D decimals, then S must be less than or equal to
// 2<sup>63</sup>-1, which is 9,223,372,036,854,775,807. The number of whole tokens (not parts) will
// be S / 10<sup>D</sup>.
//
// If decimals is 8 or 11, then the number of whole tokens can be at most a few billions or
// millions, respectively. For example, it could match Bitcoin (21 million whole tokens with 8
// decimals) or hbars (50 billion whole tokens with 8 decimals). It could even match Bitcoin with
// milli-satoshis (21 million whole tokens with 11 decimals).
//
// Note that a created token is immutable if the adminKey is omitted. No property of
// an immutable token can ever change, with the sole exception of its expiry. Anyone can pay to
// extend the expiry time of an immutable token.
//
// A token can be either FUNGIBLE_COMMON or NON_FUNGIBLE_UNIQUE, based on its
// TokenType. If it has been omitted, FUNGIBLE_COMMON type is used.
//
// A token can have either INFINITE or FINITE supply type, based on its
// TokenType. If it has been omitted, INFINITE type is used.
//
// If a FUNGIBLE TokenType is used, initialSupply should explicitly be set to a
// non-negative. If not, the transaction will resolve to INVALID_TOKEN_INITIAL_SUPPLY.
//
// If a NON_FUNGIBLE_UNIQUE TokenType is used, initialSupply should explicitly be set
// to 0. If not, the transaction will resolve to INVALID_TOKEN_INITIAL_SUPPLY.
//
// If an INFINITE TokenSupplyType is used, maxSupply should explicitly be set to 0. If
// it is not 0, the transaction will resolve to INVALID_TOKEN_MAX_SUPPLY.
//
// If a FINITE TokenSupplyType is used, maxSupply should be explicitly set to a
// non-negative value. If it is not, the transaction will resolve to INVALID_TOKEN_MAX_SUPPLY.
func NewTokenCreateTransaction() *TokenCreateTransaction {
	tx := TokenCreateTransaction{
		Transaction: _NewTransaction(),
	}

	tx.SetAutoRenewPeriod(7890000 * time.Second)
	tx._SetDefaultMaxTransactionFee(NewHbar(40))
	tx.SetTokenType(TokenTypeFungibleCommon)

	return &tx
}

func _TokenCreateTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenCreateTransaction {
	customFees := make([]Fee, 0)

	for _, fee := range pb.GetTokenCreation().GetCustomFees() {
		customFees = append(customFees, _CustomFeeFromProtobuf(fee))
	}
	adminKey, _ := _KeyFromProtobuf(pb.GetTokenCreation().GetAdminKey())
	kycKey, _ := _KeyFromProtobuf(pb.GetTokenCreation().GetKycKey())
	freezeKey, _ := _KeyFromProtobuf(pb.GetTokenCreation().GetFreezeKey())
	wipeKey, _ := _KeyFromProtobuf(pb.GetTokenCreation().GetWipeKey())
	scheduleKey, _ := _KeyFromProtobuf(pb.GetTokenCreation().GetFeeScheduleKey())
	supplyKey, _ := _KeyFromProtobuf(pb.GetTokenCreation().GetSupplyKey())
	pauseKey, _ := _KeyFromProtobuf(pb.GetTokenCreation().GetPauseKey())

	freezeDefault := pb.GetTokenCreation().GetFreezeDefault()

	expirationTime := _TimeFromProtobuf(pb.GetTokenCreation().GetExpiry())
	autoRenew := _DurationFromProtobuf(pb.GetTokenCreation().GetAutoRenewPeriod())

	resultTx := &TokenCreateTransaction{
		Transaction:        tx,
		treasuryAccountID:  _AccountIDFromProtobuf(pb.GetTokenCreation().GetTreasury()),
		autoRenewAccountID: _AccountIDFromProtobuf(pb.GetTokenCreation().GetAutoRenewAccount()),
		customFees:         customFees,
		tokenName:          pb.GetTokenCreation().GetName(),
		memo:               pb.GetTokenCreation().GetMemo(),
		tokenSymbol:        pb.GetTokenCreation().GetSymbol(),
		decimals:           pb.GetTokenCreation().GetDecimals(),
		tokenSupplyType:    TokenSupplyType(pb.GetTokenCreation().GetSupplyType()),
		tokenType:          TokenType(pb.GetTokenCreation().GetTokenType()),
		maxSupply:          pb.GetTokenCreation().GetMaxSupply(),
		adminKey:           adminKey,
		kycKey:             kycKey,
		freezeKey:          freezeKey,
		wipeKey:            wipeKey,
		scheduleKey:        scheduleKey,
		supplyKey:          supplyKey,
		pauseKey:           pauseKey,
		initialSupply:      pb.GetTokenCreation().InitialSupply,
		freezeDefault:      &freezeDefault,
		expirationTime:     &expirationTime,
		autoRenewPeriod:    &autoRenew,
	}
	return resultTx
}

// SetTokenName Sets the publicly visible name of the token, specified as a string of only ASCII characters
func (tx *TokenCreateTransaction) SetTokenName(name string) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.tokenName = name
	return tx
}

// GetTokenName returns the token name
func (tx *TokenCreateTransaction) GetTokenName() string {
	return tx.tokenName
}

// SetTokenSymbol Sets the publicly visible token symbol. It is UTF-8 capitalized alphabetical string identifying the token
func (tx *TokenCreateTransaction) SetTokenSymbol(symbol string) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.tokenSymbol = symbol
	return tx
}

// SetTokenMemo Sets the publicly visible token memo. It is max 100 bytes.
func (tx *TokenCreateTransaction) SetTokenMemo(memo string) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.memo = memo
	return tx
}

// GetTokenMemo returns the token memo
func (tx *TokenCreateTransaction) GetTokenMemo() string {
	return tx.memo
}

// GetTokenSymbol returns the token symbol
func (tx *TokenCreateTransaction) GetTokenSymbol() string {
	return tx.tokenSymbol
}

// SetDecimals Sets the number of decimal places a token is divisible by. This field can never be changed!
func (tx *TokenCreateTransaction) SetDecimals(decimals uint) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.decimals = uint32(decimals)
	return tx
}

// GetDecimals returns the number of decimal places a token is divisible by
func (tx *TokenCreateTransaction) GetDecimals() uint {
	return uint(tx.decimals)
}

// SetTokenType Specifies the token type. Defaults to FUNGIBLE_COMMON
func (tx *TokenCreateTransaction) SetTokenType(t TokenType) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.tokenType = t
	return tx
}

// GetTokenType returns the token type
func (tx *TokenCreateTransaction) GetTokenType() TokenType {
	return tx.tokenType
}

// SetSupplyType Specifies the token supply type. Defaults to INFINITE
func (tx *TokenCreateTransaction) SetSupplyType(tokenSupply TokenSupplyType) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.tokenSupplyType = tokenSupply
	return tx
}

// GetSupplyType returns the token supply type
func (tx *TokenCreateTransaction) GetSupplyType() TokenSupplyType {
	return tx.tokenSupplyType
}

// SetMaxSupply Depends on TokenSupplyType. For tokens of type FUNGIBLE_COMMON - sets the
// maximum number of tokens that can be in circulation. For tokens of type NON_FUNGIBLE_UNIQUE -
// sets the maximum number of NFTs (serial numbers) that can be minted. This field can never be
// changed!
func (tx *TokenCreateTransaction) SetMaxSupply(maxSupply int64) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.maxSupply = maxSupply
	return tx
}

// GetMaxSupply returns the max supply
func (tx *TokenCreateTransaction) GetMaxSupply() int64 {
	return tx.maxSupply
}

// SetTreasuryAccountID Sets the account which will act as a treasury for the token. This account will receive the specified initial supply
func (tx *TokenCreateTransaction) SetTreasuryAccountID(treasuryAccountID AccountID) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.treasuryAccountID = &treasuryAccountID
	return tx
}

// GetTreasuryAccountID returns the treasury account ID
func (tx *TokenCreateTransaction) GetTreasuryAccountID() AccountID {
	if tx.treasuryAccountID == nil {
		return AccountID{}
	}

	return *tx.treasuryAccountID
}

// SetAdminKey Sets the key which can perform update/delete operations on the token. If empty, the token can be perceived as immutable (not being able to be updated/deleted)
func (tx *TokenCreateTransaction) SetAdminKey(publicKey Key) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.adminKey = publicKey
	return tx
}

// GetAdminKey returns the admin key
func (tx *TokenCreateTransaction) GetAdminKey() Key {
	return tx.adminKey
}

// SetKycKey Sets the key which can grant or revoke KYC of an account for the token's transactions. If empty, KYC is not required, and KYC grant or revoke operations are not possible.
func (tx *TokenCreateTransaction) SetKycKey(publicKey Key) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.kycKey = publicKey
	return tx
}

func (tx *TokenCreateTransaction) GetKycKey() Key {
	return tx.kycKey
}

// SetFreezeKey Sets the key which can sign to freeze or unfreeze an account for token transactions. If empty, freezing is not possible
func (tx *TokenCreateTransaction) SetFreezeKey(publicKey Key) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.freezeKey = publicKey
	return tx
}

// GetFreezeKey returns the freeze key
func (tx *TokenCreateTransaction) GetFreezeKey() Key {
	return tx.freezeKey
}

// SetWipeKey Sets the key which can wipe the token balance of an account. If empty, wipe is not possible
func (tx *TokenCreateTransaction) SetWipeKey(publicKey Key) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.wipeKey = publicKey
	return tx
}

// GetWipeKey returns the wipe key
func (tx *TokenCreateTransaction) GetWipeKey() Key {
	return tx.wipeKey
}

// SetFeeScheduleKey Set the key which can change the token's custom fee schedule; must sign a TokenFeeScheduleUpdate
// transaction
func (tx *TokenCreateTransaction) SetFeeScheduleKey(key Key) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.scheduleKey = key
	return tx
}

// GetFeeScheduleKey returns the fee schedule key
func (tx *TokenCreateTransaction) GetFeeScheduleKey() Key {
	return tx.scheduleKey
}

// SetPauseKey Set the Key which can pause and unpause the Token.
// If Empty the token pause status defaults to PauseNotApplicable, otherwise Unpaused.
func (tx *TokenCreateTransaction) SetPauseKey(key Key) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.pauseKey = key
	return tx
}

// GetPauseKey returns the pause key
func (tx *TokenCreateTransaction) GetPauseKey() Key {
	return tx.pauseKey
}

// SetCustomFees Set the custom fees to be assessed during a CryptoTransfer that transfers units of this token
func (tx *TokenCreateTransaction) SetCustomFees(customFee []Fee) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.customFees = customFee
	return tx
}

// GetCustomFees returns the custom fees
func (tx *TokenCreateTransaction) GetCustomFees() []Fee {
	return tx.customFees
}

// The key which can change the supply of a token. The key is used to sign Token Mint/Burn operations
// SetInitialBalance sets the initial number of Hbar to put into the token
func (tx *TokenCreateTransaction) SetSupplyKey(publicKey Key) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.supplyKey = publicKey
	return tx
}

func (tx *TokenCreateTransaction) GetSupplyKey() Key {
	return tx.supplyKey
}

// Specifies the initial supply of tokens to be put in circulation. The initial supply is sent to the Treasury Account. The supply is in the lowest denomination possible.
func (tx *TokenCreateTransaction) SetInitialSupply(initialSupply uint64) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.initialSupply = initialSupply
	return tx
}

func (tx *TokenCreateTransaction) GetInitialSupply() uint64 {
	return tx.initialSupply
}

// The default Freeze status (frozen or unfrozen) of Hedera accounts relative to this token. If true, an account must be unfrozen before it can receive the token
func (tx *TokenCreateTransaction) SetFreezeDefault(freezeDefault bool) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.freezeDefault = &freezeDefault
	return tx
}

// GetFreezeDefault returns the freeze default
func (tx *TokenCreateTransaction) GetFreezeDefault() bool {
	return *tx.freezeDefault
}

// The epoch second at which the token should expire; if an auto-renew account and period are specified, this is coerced to the current epoch second plus the autoRenewPeriod
func (tx *TokenCreateTransaction) SetExpirationTime(expirationTime time.Time) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewPeriod = nil
	tx.expirationTime = &expirationTime

	return tx
}

func (tx *TokenCreateTransaction) GetExpirationTime() time.Time {
	if tx.expirationTime != nil {
		return *tx.expirationTime
	}

	return time.Time{}
}

// An account which will be automatically charged to renew the token's expiration, at autoRenewPeriod interval
func (tx *TokenCreateTransaction) SetAutoRenewAccount(autoRenewAccountID AccountID) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewAccountID = &autoRenewAccountID
	return tx
}

func (tx *TokenCreateTransaction) GetAutoRenewAccount() AccountID {
	if tx.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *tx.autoRenewAccountID
}

// The interval at which the auto-renew account will be charged to extend the token's expiry
func (tx *TokenCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *TokenCreateTransaction {
	tx._RequireNotFrozen()
	tx.autoRenewPeriod = &autoRenewPeriod
	return tx
}

func (tx *TokenCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if tx.autoRenewPeriod != nil {
		return time.Duration(int64(tx.autoRenewPeriod.Seconds()) * time.Second.Nanoseconds())
	}

	return time.Duration(0)
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenCreateTransaction) Sign(privateKey PrivateKey) *TokenCreateTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenCreateTransaction) SignWithOperator(client *Client) (*TokenCreateTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the transaction's signature map
// with the publicKey as the map key.
func (tx *TokenCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenCreateTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenCreateTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenCreateTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenCreateTransaction) Freeze() (*TokenCreateTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenCreateTransaction) FreezeWith(client *Client) (*TokenCreateTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenCreateTransaction.
func (tx *TokenCreateTransaction) SetMaxTransactionFee(fee Hbar) *TokenCreateTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenCreateTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenCreateTransaction.
func (tx *TokenCreateTransaction) SetTransactionMemo(memo string) *TokenCreateTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenCreateTransaction.
func (tx *TokenCreateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenCreateTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// SetTransactionID sets the TransactionID for this TokenCreateTransaction.
func (tx *TokenCreateTransaction) SetTransactionID(transactionID TransactionID) *TokenCreateTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenCreateTransaction.
func (tx *TokenCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenCreateTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenCreateTransaction) SetMaxRetry(count int) *TokenCreateTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenCreateTransaction) SetMaxBackoff(max time.Duration) *TokenCreateTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenCreateTransaction) SetMinBackoff(min time.Duration) *TokenCreateTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenCreateTransaction) SetLogLevel(level LogLevel) *TokenCreateTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenCreateTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenCreateTransaction) getName() string {
	return "TokenCreateTransaction"
}

func (tx *TokenCreateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
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

	for _, customFee := range tx.customFees {
		if err := customFee.validateNetworkOnIDs(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *TokenCreateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenCreation{
			TokenCreation: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenCreateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenCreation{
			TokenCreation: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenCreateTransaction) buildProtoBody() *services.TokenCreateTransactionBody {
	body := &services.TokenCreateTransactionBody{
		Name:          tx.tokenName,
		Symbol:        tx.tokenSymbol,
		Memo:          tx.memo,
		Decimals:      tx.decimals,
		TokenType:     services.TokenType(tx.tokenType),
		SupplyType:    services.TokenSupplyType(tx.tokenSupplyType),
		MaxSupply:     tx.maxSupply,
		InitialSupply: tx.initialSupply,
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

	if body.CustomFees == nil {
		body.CustomFees = make([]*services.CustomFee, 0)
	}
	for _, customFee := range tx.customFees {
		body.CustomFees = append(body.CustomFees, customFee._ToProtobuf())
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

	if tx.freezeDefault != nil {
		body.FreezeDefault = *tx.freezeDefault
	}

	return body
}

func (tx *TokenCreateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().CreateToken,
	}
}

func (tx *TokenCreateTransaction) preFreezeWith(client *Client) {
	if tx.autoRenewAccountID == nil && tx.autoRenewPeriod != nil && client != nil && !client.GetOperatorAccountID()._IsZero() {
		tx.SetAutoRenewAccount(client.GetOperatorAccountID())
	}
}

func (tx *TokenCreateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
