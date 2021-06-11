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
type NftCreateTransaction struct {
	Transaction
	pb *proto.TokenCreateTransactionBody
}

func NewNftCreateTransaction() *NftCreateTransaction {
	pb := &proto.TokenCreateTransactionBody{}

	transaction := NftCreateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)
	transaction.SetMaxTransactionFee(NewHbar(30))
	transaction.pb.TokenType = proto.TokenType_NON_FUNGIBLE_UNIQUE

	return &transaction
}

func nftCreateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenCreateTransaction {
	return TokenCreateTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenCreation(),
	}
}

// The publicly visible name of the token, specified as a string of only ASCII characters
func (transaction *NftCreateTransaction) SetTokenName(name string) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Name = name
	return transaction
}

func (transaction *NftCreateTransaction) GetTokenName() string {
	return transaction.pb.GetName()
}

// The publicly visible token symbol. It is UTF-8 capitalized alphabetical string identifying the token
func (transaction *NftCreateTransaction) SetTokenSymbol(symbol string) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Symbol = symbol
	return transaction
}

// The publicly visible token memo. It is max 100 bytes.
func (transaction *NftCreateTransaction) SetTokenMemo(memo string) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Memo = memo
	return transaction
}

func (transaction *NftCreateTransaction) GetTokenMemo() string {
	return transaction.pb.GetMemo()
}

func (transaction *NftCreateTransaction) GetTokenSymbol() string {
	return transaction.pb.GetSymbol()
}

// The number of decimal places a token is divisible by. This field can never be changed!
func (transaction *NftCreateTransaction) SetDecimals(decimals uint) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Decimals = uint32(decimals)
	return transaction
}

func (transaction *NftCreateTransaction) GetDecimals() uint {
	return uint(transaction.pb.GetDecimals())
}

func (transaction *NftCreateTransaction) GetTokenType() TokenType {
	return TokenType(transaction.pb.TokenType)
}

func (transaction *NftCreateTransaction) SetSupplyType(t TokenSupplyType) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.SupplyType = proto.TokenSupplyType(t)
	return transaction
}

func (transaction *NftCreateTransaction) GetSupplyType() TokenSupplyType {
	return TokenSupplyType(transaction.pb.SupplyType)
}

func (transaction *NftCreateTransaction) SetMaxSupply(maxSupply int64) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.MaxSupply = maxSupply
	return transaction
}

func (transaction *NftCreateTransaction) GetMaxSupply() int64 {
	return transaction.pb.MaxSupply
}

// The account which will act as a treasury for the token. This account will receive the specified initial supply
func (transaction *NftCreateTransaction) SetTreasuryAccountID(treasury AccountID) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Treasury = treasury.toProtobuf()
	return transaction
}

func (transaction *NftCreateTransaction) GetTreasuryAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetTreasury())
}

// The key which can perform update/delete operations on the token. If empty, the token can be perceived as immutable (not being able to be updated/deleted)
func (transaction *NftCreateTransaction) SetAdminKey(publicKey Key) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AdminKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *NftCreateTransaction) GetAdminKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetAdminKey())
	if err != nil {
		return PublicKey{}
	}

	return key
}

// The key which can grant or revoke KYC of an account for the token's transactions. If empty, KYC is not required, and KYC grant or revoke operations are not possible.
func (transaction *NftCreateTransaction) SetKycKey(publicKey Key) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.KycKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *NftCreateTransaction) GetKycKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetKycKey())
	if err != nil {
		return PublicKey{}
	}

	return key
}

// The key which can sign to freeze or unfreeze an account for token transactions. If empty, freezing is not possible
func (transaction *NftCreateTransaction) SetFreezeKey(publicKey Key) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.FreezeKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *NftCreateTransaction) GetFreezeKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetFreezeKey())
	if err != nil {
		return PublicKey{}
	}

	return key
}

// The key which can wipe the token balance of an account. If empty, wipe is not possible
func (transaction *NftCreateTransaction) SetWipeKey(publicKey Key) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.WipeKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *NftCreateTransaction) GetWipeKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetWipeKey())
	if err != nil {
		return PublicKey{}
	}

	return key
}

func (transaction *NftCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *NftCreateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
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
func (transaction *NftCreateTransaction) SetSupplyKey(publicKey Key) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.SupplyKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *NftCreateTransaction) GetSupplyKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetSupplyKey())
	if err != nil {
		return PublicKey{}
	}

	return key
}

// Specifies the initial supply of tokens to be put in circulation. The initial supply is sent to the Treasury Account. The supply is in the lowest denomination possible.
func (transaction *NftCreateTransaction) SetInitialSupply(initialSupply uint64) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.InitialSupply = initialSupply
	return transaction
}

func (transaction *NftCreateTransaction) GetInitialSupply() uint64 {
	return transaction.pb.GetInitialSupply()
}

// The default Freeze status (frozen or unfrozen) of Hedera accounts relative to this token. If true, an account must be unfrozen before it can receive the token
func (transaction *NftCreateTransaction) SetFreezeDefault(freezeDefault bool) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.FreezeDefault = freezeDefault
	return transaction
}

func (transaction *NftCreateTransaction) GetFreezeDefault() bool {
	return transaction.pb.GetFreezeDefault()
}

// The epoch second at which the token should expire; if an auto-renew account and period are specified, this is coerced to the current epoch second plus the autoRenewPeriod
func (transaction *NftCreateTransaction) SetExpirationTime(expirationTime time.Time) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewPeriod = nil
	transaction.pb.Expiry = timeToProtobuf(expirationTime)

	return transaction
}

func (transaction *NftCreateTransaction) GetExpirationTime() time.Time {
	if transaction.pb.GetExpiry() != nil {
		return time.Unix(transaction.pb.GetExpiry().Seconds, int64(transaction.pb.GetExpiry().Nanos))
	}

	return time.Time{}
}

// An account which will be automatically charged to renew the token's expiration, at autoRenewPeriod interval
func (transaction *NftCreateTransaction) SetAutoRenewAccount(autoRenewAccount AccountID) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewAccount = autoRenewAccount.toProtobuf()
	return transaction
}

func (transaction *NftCreateTransaction) GetAutoRenewAccount() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetAutoRenewAccount())
}

// The interval at which the auto-renew account will be charged to extend the token's expiry
func (transaction *NftCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewPeriod = durationToProtobuf(autoRenewPeriod)
	return transaction
}

func (transaction *NftCreateTransaction) GetAutoRenewPeriod() time.Duration {
	if transaction.pb.GetAutoRenewPeriod() != nil {
		return time.Duration(transaction.pb.GetAutoRenewPeriod().Seconds * time.Second.Nanoseconds())
	}

	return time.Duration(0)
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func nftCreateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().CreateToken,
	}
}

func (transaction *NftCreateTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *NftCreateTransaction) Sign(
	privateKey PrivateKey,
) *NftCreateTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *NftCreateTransaction) SignWithOperator(
	client *Client,
) (*NftCreateTransaction, error) {
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
func (transaction *NftCreateTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *NftCreateTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	} else {
		transaction.transactions = make([]*proto.Transaction, 0)
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.signedTransactions); index++ {
		signature := signer(transaction.signedTransactions[index].GetBodyBytes())

		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *NftCreateTransaction) Execute(
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
		nftCreateTransaction_getMethod,
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

func (transaction *NftCreateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenCreation{
		TokenCreation: transaction.pb,
	}

	return true
}

func (transaction *NftCreateTransaction) Freeze() (*NftCreateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *NftCreateTransaction) FreezeWith(client *Client) (*NftCreateTransaction, error) {
	if transaction.pb.AutoRenewPeriod != nil && client != nil && !client.GetOperatorAccountID().isZero() {
		transaction.SetAutoRenewAccount(client.GetOperatorAccountID())
	}

	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *NftCreateTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenCreateTransaction.
func (transaction *NftCreateTransaction) SetMaxTransactionFee(fee Hbar) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *NftCreateTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenCreateTransaction.
func (transaction *NftCreateTransaction) SetTransactionMemo(memo string) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *NftCreateTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenCreateTransaction.
func (transaction *NftCreateTransaction) SetTransactionValidDuration(duration time.Duration) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *NftCreateTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenCreateTransaction.
func (transaction *NftCreateTransaction) SetTransactionID(transactionID TransactionID) *NftCreateTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenCreateTransaction.
func (transaction *NftCreateTransaction) SetNodeAccountIDs(nodeID []AccountID) *NftCreateTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *NftCreateTransaction) SetMaxRetry(count int) *NftCreateTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *NftCreateTransaction) AddSignature(publicKey PublicKey, signature []byte) *NftCreateTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
