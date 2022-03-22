package hedera

import (
	"fmt"
	"time"

	"google.golang.org/protobuf/types/known/wrapperspb"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// Updates an already created Token.
// If no value is given for a field, that field is left unchanged. For an immutable tokens
// (that is, a token created without an adminKey), only the expiry may be updated. Setting any
// other field in that case will cause the transaction status to resolve to TOKEN_IS_IMMUTABlE.
type TokenUpdateTransaction struct {
	Transaction
	tokenID            *TokenID
	treasuryAccountID  *AccountID
	autoRenewAccountID *AccountID
	tokenName          string
	memo               string
	tokenSymbol        string
	adminKey           Key
	kycKey             Key
	freezeKey          Key
	wipeKey            Key
	scheduleKey        Key
	supplyKey          Key
	pauseKey           Key
	expirationTime     *time.Time
	autoRenewPeriod    *time.Duration
}

func NewTokenUpdateTransaction() *TokenUpdateTransaction {
	transaction := TokenUpdateTransaction{
		Transaction: _NewTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func _TokenUpdateTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *TokenUpdateTransaction {
	adminKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetAdminKey())
	kycKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetKycKey())
	freezeKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetFreezeKey())
	wipeKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetWipeKey())
	scheduleKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetFeeScheduleKey())
	supplyKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetSupplyKey())
	pauseKey, _ := _KeyFromProtobuf(pb.GetTokenUpdate().GetPauseKey())

	expirationTime := _TimeFromProtobuf(pb.GetTokenUpdate().GetExpiry())
	autoRenew := _DurationFromProtobuf(pb.GetTokenUpdate().GetAutoRenewPeriod())

	return &TokenUpdateTransaction{
		Transaction:        transaction,
		tokenID:            _TokenIDFromProtobuf(pb.GetTokenUpdate().GetToken()),
		treasuryAccountID:  _AccountIDFromProtobuf(pb.GetTokenUpdate().GetTreasury()),
		autoRenewAccountID: _AccountIDFromProtobuf(pb.GetTokenUpdate().GetAutoRenewAccount()),
		tokenName:          pb.GetTokenUpdate().GetName(),
		memo:               pb.GetTokenUpdate().GetMemo().Value,
		tokenSymbol:        pb.GetTokenUpdate().GetSymbol(),
		adminKey:           adminKey,
		kycKey:             kycKey,
		freezeKey:          freezeKey,
		wipeKey:            wipeKey,
		scheduleKey:        scheduleKey,
		supplyKey:          supplyKey,
		pauseKey:           pauseKey,
		expirationTime:     &expirationTime,
		autoRenewPeriod:    &autoRenew,
	}
}

// The Token to be updated
func (transaction *TokenUpdateTransaction) SetTokenID(tokenID TokenID) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenID = &tokenID
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTokenID() TokenID {
	if transaction.tokenID == nil {
		return TokenID{}
	}

	return *transaction.tokenID
}

// The new Symbol of the Token. Must be UTF-8 capitalized alphabetical string identifying the token.
func (transaction *TokenUpdateTransaction) SetTokenSymbol(symbol string) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenSymbol = symbol
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTokenSymbol() string {
	return transaction.tokenSymbol
}

// The new Name of the Token. Must be a string of ASCII characters.
func (transaction *TokenUpdateTransaction) SetTokenName(name string) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.tokenName = name
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTokenName() string {
	return transaction.tokenName
}

// The new Treasury account of the Token. If the provided treasury account is not existing or
// deleted, the _Response will be INVALID_TREASURY_ACCOUNT_FOR_TOKEN. If successful, the Token
// balance held in the previous Treasury Account is transferred to the new one.
func (transaction *TokenUpdateTransaction) SetTreasuryAccountID(treasuryAccountID AccountID) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.treasuryAccountID = &treasuryAccountID
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTreasuryAccountID() AccountID {
	if transaction.treasuryAccountID == nil {
		return AccountID{}
	}

	return *transaction.treasuryAccountID
}

// The new Admin key of the Token. If Token is immutable, transaction will resolve to
// TOKEN_IS_IMMUTABlE.
func (transaction *TokenUpdateTransaction) SetAdminKey(publicKey Key) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.adminKey = publicKey
	return transaction
}

func (transaction *TokenUpdateTransaction) GetAdminKey() Key {
	return transaction.adminKey
}

func (transaction *TokenUpdateTransaction) SetPauseKey(publicKey Key) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.pauseKey = publicKey
	return transaction
}

func (transaction *TokenUpdateTransaction) GetPauseKey() Key {
	return transaction.pauseKey
}

// The new KYC key of the Token. If Token does not have currently a KYC key, transaction will
// resolve to TOKEN_HAS_NO_KYC_KEY.
func (transaction *TokenUpdateTransaction) SetKycKey(publicKey Key) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.kycKey = publicKey
	return transaction
}

func (transaction *TokenUpdateTransaction) GetKycKey() Key {
	return transaction.kycKey
}

// The new Freeze key of the Token. If the Token does not have currently a Freeze key, transaction
// will resolve to TOKEN_HAS_NO_FREEZE_KEY.
func (transaction *TokenUpdateTransaction) SetFreezeKey(publicKey Key) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.freezeKey = publicKey
	return transaction
}

func (transaction *TokenUpdateTransaction) GetFreezeKey() Key {
	return transaction.freezeKey
}

// The new Wipe key of the Token. If the Token does not have currently a Wipe key, transaction
// will resolve to TOKEN_HAS_NO_WIPE_KEY.
func (transaction *TokenUpdateTransaction) SetWipeKey(publicKey Key) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.wipeKey = publicKey
	return transaction
}

func (transaction *TokenUpdateTransaction) GetWipeKey() Key {
	return transaction.wipeKey
}

// The new Supply key of the Token. If the Token does not have currently a Supply key, transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
func (transaction *TokenUpdateTransaction) SetSupplyKey(publicKey Key) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.supplyKey = publicKey
	return transaction
}

func (transaction *TokenUpdateTransaction) GetSupplyKey() Key {
	return transaction.supplyKey
}

func (transaction *TokenUpdateTransaction) SetFeeScheduleKey(key Key) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.scheduleKey = key
	return transaction
}

func (transaction *TokenUpdateTransaction) GetFeeScheduleKey() Key {
	return transaction.scheduleKey
}

// The new account which will be automatically charged to renew the token's expiration, at
// autoRenewPeriod interval.
func (transaction *TokenUpdateTransaction) SetAutoRenewAccount(autoRenewAccountID AccountID) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewAccountID = &autoRenewAccountID
	return transaction
}

func (transaction *TokenUpdateTransaction) GetAutoRenewAccount() AccountID {
	if transaction.autoRenewAccountID == nil {
		return AccountID{}
	}

	return *transaction.autoRenewAccountID
}

// The new interval at which the auto-renew account will be charged to extend the token's expiry.
func (transaction *TokenUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.autoRenewPeriod = &autoRenewPeriod
	return transaction
}

func (transaction *TokenUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.autoRenewPeriod != nil {
		return time.Duration(int64(transaction.autoRenewPeriod.Seconds()) * time.Second.Nanoseconds())
	}

	return time.Duration(0)
}

// The new expiry time of the token. Expiry can be updated even if admin key is not set. If the
// provided expiry is earlier than the current token expiry, transaction wil resolve to
// INVALID_EXPIRATION_TIME
func (transaction *TokenUpdateTransaction) SetExpirationTime(expirationTime time.Time) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.expirationTime = &expirationTime
	return transaction
}

func (transaction *TokenUpdateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return *transaction.expirationTime
	}
	return time.Time{}
}

func (transaction *TokenUpdateTransaction) SetTokenMemo(memo string) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.memo = memo

	return transaction
}

func (transaction *TokenUpdateTransaction) GeTokenMemo() string {
	return transaction.memo
}

func (transaction *TokenUpdateTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.tokenID != nil {
		if err := transaction.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
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

	return nil
}

func (transaction *TokenUpdateTransaction) _Build() *services.TransactionBody {
	body := &services.TokenUpdateTransactionBody{
		Name:   transaction.tokenName,
		Symbol: transaction.tokenSymbol,
		Memo:   &wrapperspb.StringValue{Value: transaction.memo},
	}

	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
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

	if transaction.pauseKey != nil {
		body.PauseKey = transaction.pauseKey._ToProtoKey()
	}

	return &services.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenUpdate{
			TokenUpdate: body,
		},
	}
}

func (transaction *TokenUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenUpdateTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	body := &services.TokenUpdateTransactionBody{
		Name:   transaction.tokenName,
		Symbol: transaction.tokenSymbol,
		Memo:   &wrapperspb.StringValue{Value: transaction.memo},
	}

	if transaction.tokenID != nil {
		body.Token = transaction.tokenID._ToProtobuf()
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

	if transaction.pauseKey != nil {
		body.PauseKey = transaction.pauseKey._ToProtoKey()
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenUpdate{
			TokenUpdate: body,
		},
	}, nil
}

func _TokenUpdateTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().UpdateToken,
	}
}

func (transaction *TokenUpdateTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenUpdateTransaction) Sign(
	privateKey PrivateKey,
) *TokenUpdateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenUpdateTransaction) SignWithOperator(
	client *Client,
) (*TokenUpdateTransaction, error) {
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
func (transaction *TokenUpdateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenUpdateTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenUpdateTransaction) Execute(
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
		_TokenUpdateTransactionGetMethod,
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

func (transaction *TokenUpdateTransaction) Freeze() (*TokenUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenUpdateTransaction) FreezeWith(client *Client) (*TokenUpdateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	err := transaction._ValidateNetworkOnIDs(client)
	if err != nil {
		return &TokenUpdateTransaction{}, err
	}
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction._Build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenUpdateTransaction.
func (transaction *TokenUpdateTransaction) SetMaxTransactionFee(fee Hbar) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *TokenUpdateTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *TokenUpdateTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

func (transaction *TokenUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenUpdateTransaction.
func (transaction *TokenUpdateTransaction) SetTransactionMemo(memo string) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenUpdateTransaction.
func (transaction *TokenUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenUpdateTransaction.
func (transaction *TokenUpdateTransaction) SetTransactionID(transactionID TransactionID) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the _Node TokenID for this TokenUpdateTransaction.
func (transaction *TokenUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenUpdateTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenUpdateTransaction) SetMaxRetry(count int) *TokenUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenUpdateTransaction {
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

func (transaction *TokenUpdateTransaction) SetMaxBackoff(max time.Duration) *TokenUpdateTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TokenUpdateTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TokenUpdateTransaction) SetMinBackoff(min time.Duration) *TokenUpdateTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TokenUpdateTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *TokenUpdateTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("TokenUpdateTransaction:%d", timestamp.UnixNano())
}
