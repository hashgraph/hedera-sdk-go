package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// Create a new token. After the token is created, the Token ID for it is in the receipt.
// The specified Treasury Account is receiving the initial supply of tokens as-well as the tokens
// from the Token Mint operation once executed. The balance of the treasury account is decreased
// when the Token Burn operation is executed.
//
// The supply that is going to be put in circulation is going to be the initial supply provided.
// The maximum supply a token can have is 2^63-1.
//
// Example:
// Token A has initial supply set to 10_000 and decimals set to 2. The tokens that will be put
// into circulation are going be 100.
// Token B has initial supply set to 10_012_345_678 and decimals set to 8. The number of tokens
// that will be put into circulation are going to be 100.12345678
//
// Creating immutable token: Token can be created as immutable if the adminKey is omitted. In this
// case, the name, symbol, treasury, management keys, expiry and renew properties cannot be
// updated. If a token is created as immutable, anyone is able to extend the expiry time by paying the fee.
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

func NewTokenCreateTransaction() *TokenCreateTransaction {
	transaction := TokenCreateTransaction{
		Transaction: _NewTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction.SetMaxTransactionFee(NewHbar(30))
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

func (transaction *TokenCreateTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenCreateTransaction {
	transaction.Transaction.SetGrpcDeadline(deadline)
	return transaction
}

// The publicly visible name of the token, specified as a string of only ASCII characters
func (transaction *TokenCreateTransaction) SetTokenName(name string) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenName = name
	return transaction
}

func (transaction *TokenCreateTransaction) GetTokenName() string {
	return transaction.tokenName
}

// The publicly visible token symbol. It is UTF-8 capitalized alphabetical string identifying the token
func (transaction *TokenCreateTransaction) SetTokenSymbol(symbol string) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenSymbol = symbol
	return transaction
}

// The publicly visible token memo. It is max 100 bytes.
func (transaction *TokenCreateTransaction) SetTokenMemo(memo string) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo
	return transaction
}

func (transaction *TokenCreateTransaction) GetTokenMemo() string {
	return transaction.memo
}

func (transaction *TokenCreateTransaction) GetTokenSymbol() string {
	return transaction.tokenSymbol
}

// The number of decimal places a token is divisible by. This field can never be changed!
func (transaction *TokenCreateTransaction) SetDecimals(decimals uint) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.decimals = uint32(decimals)
	return transaction
}

func (transaction *TokenCreateTransaction) GetDecimals() uint {
	return uint(transaction.decimals)
}

func (transaction *TokenCreateTransaction) SetTokenType(t TokenType) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenType = t
	return transaction
}

func (transaction *TokenCreateTransaction) GetTokenType() TokenType {
	return transaction.tokenType
}

func (transaction *TokenCreateTransaction) SetSupplyType(tokenSupply TokenSupplyType) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenSupplyType = tokenSupply
	return transaction
}

func (transaction *TokenCreateTransaction) GetSupplyType() TokenSupplyType {
	return transaction.tokenSupplyType
}

func (transaction *TokenCreateTransaction) SetMaxSupply(maxSupply int64) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.maxSupply = maxSupply
	return transaction
}

func (transaction *TokenCreateTransaction) GetMaxSupply() int64 {
	return transaction.maxSupply
}

// The account which will act as a treasury for the token. This account will receive the specified initial supply
func (transaction *TokenCreateTransaction) SetTreasuryAccountID(treasuryAccountID AccountID) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.treasuryAccountID = &treasuryAccountID
	return transaction
}

func (transaction *TokenCreateTransaction) GetTreasuryAccountID() AccountID {
	if transaction.treasuryAccountID == nil {
		return AccountID{}
	}

	return *transaction.treasuryAccountID
}

// The key which can perform update/delete operations on the token. If empty, the token can be perceived as immutable (not being able to be updated/deleted)
func (transaction *TokenCreateTransaction) SetAdminKey(publicKey Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.adminKey = publicKey
	return transaction
}

func (transaction *TokenCreateTransaction) GetAdminKey() Key {
	return transaction.adminKey
}

// The key which can grant or revoke KYC of an account for the token's transactions. If empty, KYC is not required, and KYC grant or revoke operations are not possible.
func (transaction *TokenCreateTransaction) SetKycKey(publicKey Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.kycKey = publicKey
	return transaction
}

func (transaction *TokenCreateTransaction) GetKycKey() Key {
	return transaction.kycKey
}

// The key which can sign to freeze or unfreeze an account for token transactions. If empty, freezing is not possible
func (transaction *TokenCreateTransaction) SetFreezeKey(publicKey Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.freezeKey = publicKey
	return transaction
}

func (transaction *TokenCreateTransaction) GetFreezeKey() Key {
	return transaction.freezeKey
}

// The key which can wipe the token balance of an account. If empty, wipe is not possible
func (transaction *TokenCreateTransaction) SetWipeKey(publicKey Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.wipeKey = publicKey
	return transaction
}

func (transaction *TokenCreateTransaction) GetWipeKey() Key {
	return transaction.wipeKey
}

func (transaction *TokenCreateTransaction) SetFeeScheduleKey(key Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.scheduleKey = key
	return transaction
}

func (transaction *TokenCreateTransaction) GetFeeScheduleKey() Key {
	return transaction.scheduleKey
}

func (transaction *TokenCreateTransaction) SetPauseKey(key Key) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.pauseKey = key
	return transaction
}

func (transaction *TokenCreateTransaction) GetPauseKey() Key {
	return transaction.pauseKey
}

func (transaction *TokenCreateTransaction) SetCustomFees(customFee []Fee) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.customFees = customFee
	return transaction
}

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

	if transaction.adminKey != nil {
		body.KycKey = transaction.kycKey._ToProtoKey()
	}

	if transaction.wipeKey != nil {
		body.WipeKey = transaction.wipeKey._ToProtoKey()
	}

	if transaction.supplyKey != nil {
		body.SupplyKey = transaction.supplyKey._ToProtoKey()
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

func _TokenCreateTransactionGetMethod(request _Request, channel *_Channel) _Method {
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
		_Request{
			transaction: &transaction.Transaction,
		},
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_TokenCreateTransactionGetMethod,
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

func (transaction *TokenCreateTransaction) Freeze() (*TokenCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenCreateTransaction) FreezeWith(client *Client) (*TokenCreateTransaction, error) {
	if transaction.autoRenewPeriod != nil && client != nil && !client.GetOperatorAccountID()._IsZero() {
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

func (transaction *TokenCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenCreateTransaction.
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

func (transaction *TokenCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenCreateTransaction.
func (transaction *TokenCreateTransaction) SetTransactionMemo(memo string) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenCreateTransaction.
func (transaction *TokenCreateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenCreateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

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

func (transaction *TokenCreateTransaction) SetMaxRetry(count int) *TokenCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

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

func (transaction *TokenCreateTransaction) SetMaxBackoff(max time.Duration) *TokenCreateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TokenCreateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TokenCreateTransaction) SetMinBackoff(min time.Duration) *TokenCreateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

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
