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
	"fmt"
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
	transaction := TokenCreateTransaction{
		Transaction: _NewTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction._SetDefaultMaxTransactionFee(NewHbar(40))
	transaction.SetTokenType(TokenTypeFungibleCommon)

	return &transaction
}

func _TokenCreateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TokenCreateTransaction {
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

	return &TokenCreateTransaction{
		Transaction:        transaction,
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
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (transaction *TokenCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenCreateTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// SetTokenName Sets the publicly visible name of the token, specified as a string of only ASCII characters
func (transaction *TokenCreateTransaction) SetTokenName(name string) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenName = name
	return transaction
}

// GetTokenName returns the token name
func (transaction *TokenCreateTransaction) GetTokenName() string {
	return transaction.tokenName
}

// SetTokenSymbol Sets the publicly visible token symbol. It is UTF-8 capitalized alphabetical string identifying the token
func (transaction *TokenCreateTransaction) SetTokenSymbol(symbol string) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenSymbol = symbol
	return transaction
}

// SetTokenMemo Sets the publicly visible token memo. It is max 100 bytes.
func (transaction *TokenCreateTransaction) SetTokenMemo(memo string) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo
	return transaction
}

// GetTokenMemo returns the token memo
func (transaction *TokenCreateTransaction) GetTokenMemo() string {
	return transaction.memo
}

// GetTokenSymbol returns the token symbol
func (transaction *TokenCreateTransaction) GetTokenSymbol() string {
	return transaction.tokenSymbol
}

// SetDecimals Sets the number of decimal places a token is divisible by. This field can never be changed!
func (transaction *TokenCreateTransaction) SetDecimals(decimals uint) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.decimals = uint32(decimals)
	return transaction
}

// GetDecimals returns the number of decimal places a token is divisible by
func (transaction *TokenCreateTransaction) GetDecimals() uint {
	return uint(transaction.decimals)
}

// SetTokenType Specifies the token type. Defaults to FUNGIBLE_COMMON
func (transaction *TokenCreateTransaction) SetTokenType(t TokenType) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenType = t
	return transaction
}

// GetTokenType returns the token type
func (transaction *TokenCreateTransaction) GetTokenType() TokenType {
	return transaction.tokenType
}

// SetSupplyType Specifies the token supply type. Defaults to INFINITE
func (transaction *TokenCreateTransaction) SetSupplyType(tokenSupply TokenSupplyType) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenSupplyType = tokenSupply
	return transaction
}

// GetSupplyType returns the token supply type
func (transaction *TokenCreateTransaction) GetSupplyType() TokenSupplyType {
	return transaction.tokenSupplyType
}

// SetMaxSupply Depends on TokenSupplyType. For tokens of type FUNGIBLE_COMMON - sets the
// maximum number of tokens that can be in circulation. For tokens of type NON_FUNGIBLE_UNIQUE -
// sets the maximum number of NFTs (serial numbers) that can be minted. This field can never be
// changed!
func (transaction *TokenCreateTransaction) SetMaxSupply(maxSupply int64) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.maxSupply = maxSupply
	return transaction
}

// GetMaxSupply returns the max supply
func (transaction *TokenCreateTransaction) GetMaxSupply() int64 {
	return transaction.maxSupply
}

// SetTreasuryAccountID Sets the account which will act as a treasury for the token. This account will receive the specified initial supply
func (transaction *TokenCreateTransaction) SetTreasuryAccountID(treasuryAccountID AccountID) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.treasuryAccountID = &treasuryAccountID
	return transaction
}

// GetTreasuryAccountID returns the treasury account ID
func (transaction *TokenCreateTransaction) GetTreasuryAccountID() AccountID {
	if transaction.treasuryAccountID == nil {
		return AccountID{}
	}

	return *transaction.treasuryAccountID
}

// SetAdminKey Sets the key which can perform update/delete operations on the token. If empty, the token can be perceived as immutable (not being able to be updated/deleted)
func (transaction *TokenCreateTransaction) SetAdminKey(publicKey Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.adminKey = publicKey
	return transaction
}

// GetAdminKey returns the admin key
func (transaction *TokenCreateTransaction) GetAdminKey() Key {
	return transaction.adminKey
}

// SetKycKey Sets the key which can grant or revoke KYC of an account for the token's transactions. If empty, KYC is not required, and KYC grant or revoke operations are not possible.
func (transaction *TokenCreateTransaction) SetKycKey(publicKey Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.kycKey = publicKey
	return transaction
}

func (transaction *TokenCreateTransaction) GetKycKey() Key {
	return transaction.kycKey
}

// SetFreezeKey Sets the key which can sign to freeze or unfreeze an account for token transactions. If empty, freezing is not possible
func (transaction *TokenCreateTransaction) SetFreezeKey(publicKey Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.freezeKey = publicKey
	return transaction
}

// GetFreezeKey returns the freeze key
func (transaction *TokenCreateTransaction) GetFreezeKey() Key {
	return transaction.freezeKey
}

// SetWipeKey Sets the key which can wipe the token balance of an account. If empty, wipe is not possible
func (transaction *TokenCreateTransaction) SetWipeKey(publicKey Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.wipeKey = publicKey
	return transaction
}

// GetWipeKey returns the wipe key
func (transaction *TokenCreateTransaction) GetWipeKey() Key {
	return transaction.wipeKey
}

// SetFeeScheduleKey Set the key which can change the token's custom fee schedule; must sign a TokenFeeScheduleUpdate
// transaction
func (transaction *TokenCreateTransaction) SetFeeScheduleKey(key Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.scheduleKey = key
	return transaction
}

// GetFeeScheduleKey returns the fee schedule key
func (transaction *TokenCreateTransaction) GetFeeScheduleKey() Key {
	return transaction.scheduleKey
}

// SetPauseKey Set the Key which can pause and unpause the Token.
// If Empty the token pause status defaults to PauseNotApplicable, otherwise Unpaused.
func (transaction *TokenCreateTransaction) SetPauseKey(key Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.pauseKey = key
	return transaction
}

// GetPauseKey returns the pause key
func (transaction *TokenCreateTransaction) GetPauseKey() Key {
	return transaction.pauseKey
}

// SetCustomFees Set the custom fees to be assessed during a CryptoTransfer that transfers units of this token
func (transaction *TokenCreateTransaction) SetCustomFees(customFee []Fee) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.customFees = customFee
	return transaction
}

// GetCustomFees returns the custom fees
func (transaction *TokenCreateTransaction) GetCustomFees() []Fee {
	return transaction.customFees
}

func (transaction *TokenCreateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.treasuryAccountID != nil {
		if err := transaction.treasuryAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	if transaction.autoRenewAccountID != nil {
		if err := transaction.autoRenewAccountID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	for _, customFee := range transaction.customFees {
		if err := customFee._ValidateNetworkOnIDs(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TokenCreateTransaction) _Build() *services.TransactionBody {
	body := &services.TokenCreateTransactionBody{
		Name:          transaction.tokenName,
		Symbol:        transaction.tokenSymbol,
		Memo:          transaction.memo,
		Decimals:      transaction.decimals,
		TokenType:     services.TokenType(transaction.tokenType),
		SupplyType:    services.TokenSupplyType(transaction.tokenSupplyType),
		MaxSupply:     transaction.maxSupply,
		InitialSupply: transaction.initialSupply,
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.Expiry = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.treasuryAccountID != nil {
		body.Treasury = transaction.treasuryAccountID._ToProtobuf()
	}

	if transaction.autoRenewAccountID != nil {
		body.AutoRenewAccount = transaction.autoRenewAccountID._ToProtobuf()
	}

	if body.CustomFees == nil {
		body.CustomFees = make([]*services.CustomFee, 0)
	}
	for _, customFee := range transaction.customFees {
		body.CustomFees = append(body.CustomFees, customFee._ToProtobuf())
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey._ToProtoKey()
	}

	if transaction.freezeKey != nil {
		body.FreezeKey = transaction.freezeKey._ToProtoKey()
	}

	if transaction.scheduleKey != nil {
		body.FeeScheduleKey = transaction.scheduleKey._ToProtoKey()
	}

	if transaction.kycKey != nil {
		body.KycKey = transaction.kycKey._ToProtoKey()
	}

	if transaction.wipeKey != nil {
		body.WipeKey = transaction.wipeKey._ToProtoKey()
	}

	if transaction.supplyKey != nil {
		body.SupplyKey = transaction.supplyKey._ToProtoKey()
	}

	if transaction.pauseKey != nil {
		body.PauseKey = transaction.pauseKey._ToProtoKey()
	}

	if transaction.freezeDefault != nil {
		body.FreezeDefault = *transaction.freezeDefault
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenCreation{
			TokenCreation: body,
		},
	}
}

func (transaction *TokenCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenCreateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.TokenCreateTransactionBody{
		Name:          transaction.tokenName,
		Symbol:        transaction.tokenSymbol,
		Memo:          transaction.memo,
		Decimals:      transaction.decimals,
		TokenType:     services.TokenType(transaction.tokenType),
		SupplyType:    services.TokenSupplyType(transaction.tokenSupplyType),
		MaxSupply:     transaction.maxSupply,
		InitialSupply: transaction.initialSupply,
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = _DurationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.Expiry = _TimeToProtobuf(*transaction.expirationTime)
	}

	if transaction.treasuryAccountID != nil {
		body.Treasury = transaction.treasuryAccountID._ToProtobuf()
	}

	if transaction.autoRenewAccountID != nil {
		body.AutoRenewAccount = transaction.autoRenewAccountID._ToProtobuf()
	}

	if body.CustomFees == nil {
		body.CustomFees = make([]*services.CustomFee, 0)
	}
	for _, customFee := range transaction.customFees {
		body.CustomFees = append(body.CustomFees, customFee._ToProtobuf())
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey._ToProtoKey()
	}

	if transaction.freezeKey != nil {
		body.FreezeKey = transaction.freezeKey._ToProtoKey()
	}

	if transaction.scheduleKey != nil {
		body.FeeScheduleKey = transaction.scheduleKey._ToProtoKey()
	}

	if transaction.kycKey != nil {
		body.KycKey = transaction.kycKey._ToProtoKey()
	}

	if transaction.wipeKey != nil {
		body.WipeKey = transaction.wipeKey._ToProtoKey()
	}

	if transaction.supplyKey != nil {
		body.SupplyKey = transaction.supplyKey._ToProtoKey()
	}

	if transaction.pauseKey != nil {
		body.PauseKey = transaction.pauseKey._ToProtoKey()
	}

	if transaction.freezeDefault != nil {
		body.FreezeDefault = *transaction.freezeDefault
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenCreation{
			TokenCreation: body,
		},
	}, nil
}

// The key which can change the supply of a token. The key is used to sign Token Mint/Burn operations
// SetInitialBalance sets the initial number of Hbar to put into the token
func (transaction *TokenCreateTransaction) SetSupplyKey(publicKey Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.supplyKey = publicKey
	return transaction
}

func (transaction *TokenCreateTransaction) GetSupplyKey() Key {
	return transaction.supplyKey
}

// Specifies the initial supply of tokens to be put in circulation. The initial supply is sent to the Treasury Account. The supply is in the lowest denomination possible.
func (transaction *TokenCreateTransaction) SetInitialSupply(initialSupply uint64) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.initialSupply = initialSupply
	return transaction
}

func (transaction *TokenCreateTransaction) GetInitialSupply() uint64 {
	return transaction.initialSupply
}

// The default Freeze status (frozen or unfrozen) of Hedera accounts relative to this token. If true, an account must be unfrozen before it can receive the token
func (transaction *TokenCreateTransaction) SetFreezeDefault(freezeDefault bool) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.freezeDefault = &freezeDefault
	return transaction
}

// GetFreezeDefault returns the freeze default
func (transaction *TokenCreateTransaction) GetFreezeDefault() bool {
	return *transaction.freezeDefault
}

// The epoch second at which the token should expire; if an auto-renew account and period are specified, this is coerced to the current epoch second plus the autoRenewPeriod
func (transaction *TokenCreateTransaction) SetExpirationTime(expirationTime time.Time) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewPeriod = nil
	transaction.expirationTime = &expirationTime

	return transaction
}

func (transaction *TokenCreateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}

	return time.Time{}
}

// An account which will be automatically charged to renew the token's expiration, at autoRenewPeriod interval
func (transaction *TokenCreateTransaction) SetAutoRenewAccount(autoRenewAccountID AccountID) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewAccountID = &autoRenewAccountID
	return transaction
}

func (transaction *TokenCreateTransaction) GetAutoRenewAccount() AccountID {
	if transaction.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *transaction.autoRenewAccountID
}

// The interval at which the auto-renew account will be charged to extend the token's expiry
func (transaction *TokenCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewPeriod = &autoRenewPeriod
	return transaction
}

func (transaction *TokenCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return time.Duration(int64(transaction.autoRenewPeriod.Seconds()) * time.Second.Nanoseconds())
	}

	return time.Duration(0)
}

func _TokenCreateTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().CreateToken,
	}
}

func (transaction *TokenCreateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenCreateTransaction) Sign(
	privateKey PrivateKey,
) *TokenCreateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *TokenCreateTransaction) SignWithOperator(
	client *Client,
) (*TokenCreateTransaction, error) {
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
func (transaction *TokenCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenCreateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenCreateTransaction) Execute(
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
		&transaction.Transaction,
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_TokenCreateTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
		transaction._GetLogID(),
		transaction.grpcDeadline,
		transaction.maxBackoff,
		transaction.minBackoff,
		transaction.maxRetry,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID:  transaction.GetTransactionID(),
			NodeID:         resp.(TransactionResponse).NodeID,
			ValidateStatus: true,
		}, err
	}

	return TransactionResponse{
		TransactionID:  transaction.GetTransactionID(),
		NodeID:         resp.(TransactionResponse).NodeID,
		Hash:           resp.(TransactionResponse).Hash,
		ValidateStatus: true,
	}, nil
}

func (transaction *TokenCreateTransaction) Freeze() (*TokenCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenCreateTransaction) FreezeWith(client *Client) (*TokenCreateTransaction, error) {
	if transaction.autoRenewAccountID == nil && transaction.autoRenewPeriod != nil && client != nil && !client.GetOperatorAccountID()._IsZero() {
		transaction.SetAutoRenewAccount(client.GetOperatorAccountID())
	}

	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TokenCreateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *TokenCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *TokenCreateTransaction) SetMaxTransactionFee(fee Hbar) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *TokenCreateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *TokenCreateTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this	TokenCreateTransaction.
func (transaction *TokenCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenCreateTransaction.
func (transaction *TokenCreateTransaction) SetTransactionMemo(memo string) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *TokenCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenCreateTransaction.
func (transaction *TokenCreateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this	TokenCreateTransaction.
func (transaction *TokenCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenCreateTransaction.
func (transaction *TokenCreateTransaction) SetTransactionID(transactionID TransactionID) *TokenCreateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenCreateTransaction.
func (transaction *TokenCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *TokenCreateTransaction) SetMaxRetry(count int) *TokenCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *TokenCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenCreateTransaction {
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

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (transaction *TokenCreateTransaction) SetMaxBackoff(max time.Duration) *TokenCreateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the maximum amount of time to wait between retries.
func (transaction *TokenCreateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (transaction *TokenCreateTransaction) SetMinBackoff(min time.Duration) *TokenCreateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the minimum amount of time to wait between retries.
func (transaction *TokenCreateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *TokenCreateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("TokenCreateTransaction:%d", timestamp.UnixNano())
}

func (transaction *TokenCreateTransaction) SetLogLevel(level LogLevel) *TokenCreateTransaction {
	transaction.Transaction.SetLogLevel(level)
	return transaction
}
