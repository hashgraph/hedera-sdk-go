package hedera

import (
	"time"

	"github.com/golang/protobuf/ptypes/wrappers"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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
	expirationTime     *time.Time
	autoRenewPeriod    *time.Duration
}

func NewTokenUpdateTransaction() *TokenUpdateTransaction {
	transaction := TokenUpdateTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func tokenUpdateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenUpdateTransaction {
	adminKey, _ := keyFromProtobuf(pb.GetTokenUpdate().GetAdminKey())
	kycKey, _ := keyFromProtobuf(pb.GetTokenUpdate().GetKycKey())
	freezeKey, _ := keyFromProtobuf(pb.GetTokenUpdate().GetFreezeKey())
	wipeKey, _ := keyFromProtobuf(pb.GetTokenUpdate().GetWipeKey())
	scheduleKey, _ := keyFromProtobuf(pb.GetTokenUpdate().GetFeeScheduleKey())
	supplyKey, _ := keyFromProtobuf(pb.GetTokenUpdate().GetSupplyKey())

	expirationTime := timeFromProtobuf(pb.GetTokenUpdate().GetExpiry())
	autoRenew := durationFromProtobuf(pb.GetTokenUpdate().GetAutoRenewPeriod())

	return TokenUpdateTransaction{
		Transaction:        transaction,
		tokenID:            tokenIDFromProtobuf(pb.GetTokenUpdate().GetToken()),
		treasuryAccountID:  accountIDFromProtobuf(pb.GetTokenUpdate().GetTreasury()),
		autoRenewAccountID: accountIDFromProtobuf(pb.GetTokenUpdate().GetAutoRenewAccount()),
		tokenName:          pb.GetTokenUpdate().GetName(),
		memo:               pb.GetTokenUpdate().GetMemo().Value,
		tokenSymbol:        pb.GetTokenUpdate().GetSymbol(),
		adminKey:           adminKey,
		kycKey:             kycKey,
		freezeKey:          freezeKey,
		wipeKey:            wipeKey,
		scheduleKey:        scheduleKey,
		supplyKey:          supplyKey,
		expirationTime:     &expirationTime,
		autoRenewPeriod:    &autoRenew,
	}
}

// The Token to be updated
func (transaction *TokenUpdateTransaction) SetTokenID(tokenID TokenID) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
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
	transaction.requireNotFrozen()
	transaction.tokenSymbol = symbol
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTokenSymbol() string {
	return transaction.tokenSymbol
}

// The new Name of the Token. Must be a string of ASCII characters.
func (transaction *TokenUpdateTransaction) SetTokenName(name string) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.tokenName = name
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTokenName() string {
	return transaction.tokenName
}

// The new Treasury account of the Token. If the provided treasury account is not existing or
// deleted, the response will be INVALID_TREASURY_ACCOUNT_FOR_TOKEN. If successful, the Token
// balance held in the previous Treasury Account is transferred to the new one.
func (transaction *TokenUpdateTransaction) SetTreasuryAccountID(treasuryAccountID AccountID) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
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
	transaction.requireNotFrozen()
	transaction.adminKey = publicKey
	return transaction
}

func (transaction *TokenUpdateTransaction) GetAdminKey() Key {
	return transaction.adminKey
}

// The new KYC key of the Token. If Token does not have currently a KYC key, transaction will
// resolve to TOKEN_HAS_NO_KYC_KEY.
func (transaction *TokenUpdateTransaction) SetKycKey(publicKey Key) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.kycKey = publicKey
	return transaction
}

func (transaction *TokenUpdateTransaction) GetKycKey() Key {
	return transaction.kycKey
}

// The new Freeze key of the Token. If the Token does not have currently a Freeze key, transaction
// will resolve to TOKEN_HAS_NO_FREEZE_KEY.
func (transaction *TokenUpdateTransaction) SetFreezeKey(publicKey Key) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.freezeKey = publicKey
	return transaction
}

func (transaction *TokenUpdateTransaction) GetFreezeKey() Key {
	return transaction.freezeKey
}

// The new Wipe key of the Token. If the Token does not have currently a Wipe key, transaction
// will resolve to TOKEN_HAS_NO_WIPE_KEY.
func (transaction *TokenUpdateTransaction) SetWipeKey(publicKey Key) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.wipeKey = publicKey
	return transaction
}

func (transaction *TokenUpdateTransaction) GetWipeKey() Key {
	return transaction.wipeKey
}

// The new Supply key of the Token. If the Token does not have currently a Supply key, transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
func (transaction *TokenUpdateTransaction) SetSupplyKey(publicKey Key) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.supplyKey = publicKey
	return transaction
}

func (transaction *TokenUpdateTransaction) GetSupplyKey() Key {
	return transaction.supplyKey
}

func (transaction *TokenUpdateTransaction) SetFeeScheduleKey(key Key) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.scheduleKey = key
	return transaction
}

func (transaction *TokenUpdateTransaction) GetFeeScheduleKey() Key {
	return transaction.scheduleKey
}

// The new account which will be automatically charged to renew the token's expiration, at
// autoRenewPeriod interval.
func (transaction *TokenUpdateTransaction) SetAutoRenewAccount(autoRenewAccountID AccountID) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
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
	transaction.requireNotFrozen()
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
	transaction.requireNotFrozen()
	transaction.expirationTime = &expirationTime
	return transaction
}

func (transaction *TokenUpdateTransaction) GetExpirationTime() time.Time {
	if transaction.expirationTime != nil {
		return time.Unix(transaction.expirationTime.Unix(), transaction.expirationTime.UnixNano())
	}

	return time.Time{}
}

func (transaction *TokenUpdateTransaction) SetTokenMemo(memo string) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.memo = memo

	return transaction
}

func (transaction *TokenUpdateTransaction) GeTokenMemo() string {
	return transaction.memo
}

func (transaction *TokenUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.tokenID != nil {
		if err := transaction.tokenID.Validate(client); err != nil {
			return err
		}
	}

	if transaction.treasuryAccountID != nil {
		if err := transaction.treasuryAccountID.Validate(client); err != nil {
			return err
		}
	}

	if transaction.autoRenewAccountID != nil {
		if err := transaction.autoRenewAccountID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *TokenUpdateTransaction) build() *proto.TransactionBody {
	body := &proto.TokenUpdateTransactionBody{
		Name:   transaction.tokenName,
		Symbol: transaction.tokenSymbol,
		Memo:   &wrappers.StringValue{Value: transaction.memo},
	}

	if !transaction.tokenID.isZero() {
		body.Token = transaction.tokenID.toProtobuf()
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = durationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.Expiry = timeToProtobuf(*transaction.expirationTime)
	}

	if !transaction.treasuryAccountID.isZero() {
		body.Treasury = transaction.treasuryAccountID.toProtobuf()
	}

	if !transaction.autoRenewAccountID.isZero() {
		body.AutoRenewAccount = transaction.autoRenewAccountID.toProtobuf()
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey.toProtoKey()
	}

	if transaction.freezeKey != nil {
		body.FreezeKey = transaction.freezeKey.toProtoKey()
	}

	if transaction.scheduleKey != nil {
		body.FeeScheduleKey = transaction.scheduleKey.toProtoKey()
	}

	if transaction.adminKey != nil {
		body.KycKey = transaction.kycKey.toProtoKey()
	}

	if transaction.wipeKey != nil {
		body.WipeKey = transaction.wipeKey.toProtoKey()
	}

	if transaction.supplyKey != nil {
		body.SupplyKey = transaction.supplyKey.toProtoKey()
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_TokenUpdate{
			TokenUpdate: body,
		},
	}
}

func (transaction *TokenUpdateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenUpdateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.TokenUpdateTransactionBody{
		Name:   transaction.tokenName,
		Symbol: transaction.tokenSymbol,
		Memo:   &wrappers.StringValue{Value: transaction.memo},
	}

	if !transaction.tokenID.isZero() {
		body.Token = transaction.tokenID.toProtobuf()
	}

	if transaction.autoRenewPeriod != nil {
		body.AutoRenewPeriod = durationToProtobuf(*transaction.autoRenewPeriod)
	}

	if transaction.expirationTime != nil {
		body.Expiry = timeToProtobuf(*transaction.expirationTime)
	}

	if !transaction.treasuryAccountID.isZero() {
		body.Treasury = transaction.treasuryAccountID.toProtobuf()
	}

	if !transaction.autoRenewAccountID.isZero() {
		body.AutoRenewAccount = transaction.autoRenewAccountID.toProtobuf()
	}

	if transaction.adminKey != nil {
		body.AdminKey = transaction.adminKey.toProtoKey()
	}

	if transaction.freezeKey != nil {
		body.FreezeKey = transaction.freezeKey.toProtoKey()
	}

	if transaction.scheduleKey != nil {
		body.FeeScheduleKey = transaction.scheduleKey.toProtoKey()
	}

	if transaction.adminKey != nil {
		body.KycKey = transaction.kycKey.toProtoKey()
	}

	if transaction.wipeKey != nil {
		body.WipeKey = transaction.wipeKey.toProtoKey()
	}

	if transaction.supplyKey != nil {
		body.SupplyKey = transaction.supplyKey.toProtoKey()
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_TokenUpdate{
			TokenUpdate: body,
		},
	}, nil
}

func _TokenUpdateTransactionGetMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().UpdateToken,
	}
}

func (transaction *TokenUpdateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
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
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

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
	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
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

	transactionID := transaction.GetTransactionID()

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		_TransactionShouldRetry,
		_TransactionMakeRequest(request{
			transaction: &transaction.Transaction,
		}),
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_TokenUpdateTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
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
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TokenUpdateTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenUpdateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenUpdateTransaction.
func (transaction *TokenUpdateTransaction) SetMaxTransactionFee(fee Hbar) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenUpdateTransaction.
func (transaction *TokenUpdateTransaction) SetTransactionMemo(memo string) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenUpdateTransaction.
func (transaction *TokenUpdateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenUpdateTransaction.
func (transaction *TokenUpdateTransaction) SetTransactionID(transactionID TransactionID) *TokenUpdateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenUpdateTransaction.
func (transaction *TokenUpdateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenUpdateTransaction) SetMaxRetry(count int) *TokenUpdateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenUpdateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenUpdateTransaction {
	transaction.requireOneNodeAccountID()

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	if len(transaction.signedTransactions) == 0 {
		return transaction
	}

	transaction.transactions = make([]*proto.Transaction, 0)
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)

	for index := 0; index < len(transaction.signedTransactions); index++ {
		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
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
