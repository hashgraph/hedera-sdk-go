package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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
	pb                 *proto.TokenCreateTransactionBody
	treasuryAccountID  AccountID
	autoRenewAccountID AccountID
	customFees         []CustomFee
}

func NewTokenCreateTransaction() *TokenCreateTransaction {
	pb := &proto.TokenCreateTransactionBody{}

	transaction := TokenCreateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction.SetMaxTransactionFee(NewHbar(30))
	transaction.SetTokenType(TokenTypeFungibleCommon)

	return &transaction
}

func tokenCreateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenCreateTransaction {
	return TokenCreateTransaction{
		Transaction:        transaction,
		pb:                 pb.GetTokenCreation(),
		treasuryAccountID:  accountIDFromProtobuf(pb.GetTokenCreation().GetTreasury(), nil),
		autoRenewAccountID: accountIDFromProtobuf(pb.GetTokenCreation().GetAutoRenewAccount(), nil),
	}
}

// The publicly visible name of the token, specified as a string of only ASCII characters
func (transaction *TokenCreateTransaction) SetTokenName(name string) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Name = name
	return transaction
}

func (transaction *TokenCreateTransaction) GetTokenName() string {
	return transaction.pb.GetName()
}

// The publicly visible token symbol. It is UTF-8 capitalized alphabetical string identifying the token
func (transaction *TokenCreateTransaction) SetTokenSymbol(symbol string) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Symbol = symbol
	return transaction
}

// The publicly visible token memo. It is max 100 bytes.
func (transaction *TokenCreateTransaction) SetTokenMemo(memo string) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Memo = memo
	return transaction
}

func (transaction *TokenCreateTransaction) GetTokenMemo() string {
	return transaction.pb.GetMemo()
}

func (transaction *TokenCreateTransaction) GetTokenSymbol() string {
	return transaction.pb.GetSymbol()
}

// The number of decimal places a token is divisible by. This field can never be changed!
func (transaction *TokenCreateTransaction) SetDecimals(decimals uint) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Decimals = uint32(decimals)
	return transaction
}

func (transaction *TokenCreateTransaction) GetDecimals() uint {
	return uint(transaction.pb.GetDecimals())
}

func (transaction *TokenCreateTransaction) SetTokenType(t TokenType) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.TokenType = proto.TokenType(t)
	return transaction
}

func (transaction *TokenCreateTransaction) GetTokenType() TokenType {
	return TokenType(transaction.pb.TokenType)
}

func (transaction *TokenCreateTransaction) SetSupplyType(t TokenSupplyType) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.SupplyType = proto.TokenSupplyType(t)
	return transaction
}

func (transaction *TokenCreateTransaction) GetSupplyType() TokenSupplyType {
	return TokenSupplyType(transaction.pb.SupplyType)
}

func (transaction *TokenCreateTransaction) SetMaxSupply(maxSupply int64) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.MaxSupply = maxSupply
	return transaction
}

func (transaction *TokenCreateTransaction) GetMaxSupply() int64 {
	return transaction.pb.MaxSupply
}

// The account which will act as a treasury for the token. This account will receive the specified initial supply
func (transaction *TokenCreateTransaction) SetTreasuryAccountID(treasury AccountID) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Treasury = treasury.toProtobuf()
	return transaction
}

func (transaction *TokenCreateTransaction) GetTreasuryAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetTreasury(), nil)
}

// The key which can perform update/delete operations on the token. If empty, the token can be perceived as immutable (not being able to be updated/deleted)
func (transaction *TokenCreateTransaction) SetAdminKey(publicKey Key) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AdminKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TokenCreateTransaction) GetAdminKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetAdminKey(), nil)
	if err != nil {
		return PublicKey{}
	}

	return key
}

// The key which can grant or revoke KYC of an account for the token's transactions. If empty, KYC is not required, and KYC grant or revoke operations are not possible.
func (transaction *TokenCreateTransaction) SetKycKey(publicKey Key) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.KycKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TokenCreateTransaction) GetKycKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetKycKey(), nil)
	if err != nil {
		return PublicKey{}
	}

	return key
}

// The key which can sign to freeze or unfreeze an account for token transactions. If empty, freezing is not possible
func (transaction *TokenCreateTransaction) SetFreezeKey(publicKey Key) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.FreezeKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TokenCreateTransaction) GetFreezeKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetFreezeKey(), nil)
	if err != nil {
		return PublicKey{}
	}

	return key
}

// The key which can wipe the token balance of an account. If empty, wipe is not possible
func (transaction *TokenCreateTransaction) SetWipeKey(publicKey Key) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.WipeKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TokenCreateTransaction) GetWipeKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetWipeKey(), nil)
	if err != nil {
		return PublicKey{}
	}

	return key
}

func (transaction *TokenCreateTransaction) SetFeeScheduleKey(key Key) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.FeeScheduleKey = key.toProtoKey()
	return transaction
}

func (transaction *TokenCreateTransaction) GetFeeScheduleKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetFeeScheduleKey(), nil)
	if err != nil {
		return PublicKey{}
	}

	return key
}

func (transaction *TokenCreateTransaction) SetCustomFees(customFee []CustomFee) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.customFees = customFee
	return transaction
}

func (transaction *TokenCreateTransaction) AddCustomFee(customFee CustomFee) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	if transaction.customFees == nil {
		transaction.customFees = make([]CustomFee, 0)
	}

	transaction.customFees = append(transaction.customFees, customFee)
	return transaction
}

func (transaction *TokenCreateTransaction) GetCustomFees() []CustomFee {
	return transaction.customFees
}

func (transaction *TokenCreateTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	err = transaction.treasuryAccountID.Validate(client)
	if err != nil {
		return err
	}
	err = transaction.autoRenewAccountID.Validate(client)
	if err != nil {
		return err
	}
	for _, customFee := range transaction.customFees {
		err = customFee.FeeCollectorAccountID.Validate(client)
		if err != nil {
			return err
		}
		switch t := customFee.Fee.(type) {
		case CustomFixedFee:
			if t.DenominationTokenID != nil {
				err = t.DenominationTokenID.Validate(client)
				if err != nil {
					return err
				}
			}
		}
	}

	return nil
}

func (transaction *TokenCreateTransaction) build() *TokenCreateTransaction {
	if !transaction.treasuryAccountID.isZero() {
		transaction.pb.Treasury = transaction.treasuryAccountID.toProtobuf()
	}

	if !transaction.autoRenewAccountID.isZero() {
		transaction.pb.AutoRenewAccount = transaction.autoRenewAccountID.toProtobuf()
	}

	if len(transaction.customFees) > 0 {
		for _, customFee := range transaction.customFees {
			if transaction.pb.CustomFees == nil {
				transaction.pb.CustomFees = make([]*proto.CustomFee, 0)
			}
			transaction.pb.CustomFees = append(transaction.pb.CustomFees, customFee.toProtobuf())
		}
	}

	return transaction
}

func (transaction *TokenCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenCreateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	transaction.build()
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenCreation{
			TokenCreation: &proto.TokenCreateTransactionBody{
				Name:             transaction.pb.GetName(),
				Symbol:           transaction.pb.GetSymbol(),
				Decimals:         transaction.pb.GetDecimals(),
				InitialSupply:    transaction.pb.GetInitialSupply(),
				Treasury:         transaction.pb.GetTreasury(),
				AdminKey:         transaction.pb.GetAdminKey(),
				KycKey:           transaction.pb.GetKycKey(),
				FreezeKey:        transaction.pb.GetFreezeKey(),
				WipeKey:          transaction.pb.GetWipeKey(),
				SupplyKey:        transaction.pb.GetSupplyKey(),
				FreezeDefault:    transaction.pb.GetFreezeDefault(),
				Expiry:           transaction.pb.GetExpiry(),
				AutoRenewAccount: transaction.pb.GetAutoRenewAccount(),
				AutoRenewPeriod:  transaction.pb.GetAutoRenewPeriod(),
				Memo:             transaction.pb.GetMemo(),
				TokenType:        transaction.pb.GetTokenType(),
				SupplyType:       transaction.pb.GetSupplyType(),
				MaxSupply:        transaction.pb.GetMaxSupply(),
			},
		},
	}, nil
}

// The key which can change the supply of a token. The key is used to sign Token Mint/Burn operations
// SetInitialBalance sets the initial number of Hbar to put into the token
func (transaction *TokenCreateTransaction) SetSupplyKey(publicKey Key) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.SupplyKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TokenCreateTransaction) GetSupplyKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetSupplyKey(), nil)
	if err != nil {
		return PublicKey{}
	}

	return key
}

// Specifies the initial supply of tokens to be put in circulation. The initial supply is sent to the Treasury Account. The supply is in the lowest denomination possible.
func (transaction *TokenCreateTransaction) SetInitialSupply(initialSupply uint64) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.InitialSupply = initialSupply
	return transaction
}

func (transaction *TokenCreateTransaction) GetInitialSupply() uint64 {
	return transaction.pb.GetInitialSupply()
}

// The default Freeze status (frozen or unfrozen) of Hedera accounts relative to this token. If true, an account must be unfrozen before it can receive the token
func (transaction *TokenCreateTransaction) SetFreezeDefault(freezeDefault bool) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.FreezeDefault = freezeDefault
	return transaction
}

func (transaction *TokenCreateTransaction) GetFreezeDefault() bool {
	return transaction.pb.GetFreezeDefault()
}

// The epoch second at which the token should expire; if an auto-renew account and period are specified, this is coerced to the current epoch second plus the autoRenewPeriod
func (transaction *TokenCreateTransaction) SetExpirationTime(expirationTime time.Time) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewPeriod = nil
	transaction.pb.Expiry = timeToProtobuf(expirationTime)

	return transaction
}

func (transaction *TokenCreateTransaction) GetExpirationTime() time.Time {
	if transaction.pb.GetExpiry() != nil {
		return time.Unix(transaction.pb.GetExpiry().Seconds, int64(transaction.pb.GetExpiry().Nanos))
	}

	return time.Time{}
}

// An account which will be automatically charged to renew the token's expiration, at autoRenewPeriod interval
func (transaction *TokenCreateTransaction) SetAutoRenewAccount(id AccountID) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.autoRenewAccountID = id
	return transaction
}

func (transaction *TokenCreateTransaction) GetAutoRenewAccount() AccountID {
	return transaction.autoRenewAccountID
}

// The interval at which the auto-renew account will be charged to extend the token's expiry
func (transaction *TokenCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewPeriod = durationToProtobuf(autoRenewPeriod)
	return transaction
}

func (transaction *TokenCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.pb.GetAutoRenewPeriod() != nil {
		return time.Duration(transaction.pb.GetAutoRenewPeriod().Seconds * time.Second.Nanoseconds())
	}

	return time.Duration(0)
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func tokenCreateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().CreateToken,
	}
}

func (transaction *TokenCreateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
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
func (transaction *TokenCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenCreateTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
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
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		tokenCreateTransaction_getMethod,
		transaction_mapStatusError,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.transaction.NodeID,
		}, err
	}

	hash, err := transaction.GetTransactionHash()

	return TransactionResponse{
		TransactionID: transaction.GetTransactionID(),
		NodeID:        resp.transaction.NodeID,
		Hash:          hash,
	}, nil
}

func (transaction *TokenCreateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenCreation{
		TokenCreation: transaction.pb,
	}

	return true
}

func (transaction *TokenCreateTransaction) Freeze() (*TokenCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenCreateTransaction) FreezeWith(client *Client) (*TokenCreateTransaction, error) {
	if transaction.pb.AutoRenewPeriod != nil && client != nil && !client.GetOperatorAccountID().isZero() {
		transaction.SetAutoRenewAccount(client.GetOperatorAccountID())
	}

	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TokenCreateTransaction{}, err
	}
	transaction.build()

	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *TokenCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenCreateTransaction.
func (transaction *TokenCreateTransaction) SetMaxTransactionFee(fee Hbar) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenCreateTransaction.
func (transaction *TokenCreateTransaction) SetTransactionMemo(memo string) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenCreateTransaction.
func (transaction *TokenCreateTransaction) SetTransactionValidDuration(duration time.Duration) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenCreateTransaction.
func (transaction *TokenCreateTransaction) SetTransactionID(transactionID TransactionID) *TokenCreateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenCreateTransaction.
func (transaction *TokenCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenCreateTransaction) SetMaxRetry(count int) *TokenCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenCreateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
