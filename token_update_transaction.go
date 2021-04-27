package hedera

import (
	"github.com/golang/protobuf/ptypes/wrappers"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Updates an already created Token.
// If no value is given for a field, that field is left unchanged. For an immutable tokens
// (that is, a token created without an adminKey), only the expiry may be updated. Setting any
// other field in that case will cause the transaction status to resolve to TOKEN_IS_IMMUTABlE.
type TokenUpdateTransaction struct {
	Transaction
	pb *proto.TokenUpdateTransactionBody
}

func NewTokenUpdateTransaction() *TokenUpdateTransaction {
	pb := &proto.TokenUpdateTransactionBody{}

	transaction := TokenUpdateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func tokenUpdateTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenUpdateTransaction {
	return TokenUpdateTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenUpdate(),
	}
}

// The Token to be updated
func (transaction *TokenUpdateTransaction) SetTokenID(tokenID TokenID) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Token = tokenID.toProtobuf()
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTokenID() TokenID {
	return tokenIDFromProtobuf(transaction.pb.GetToken())
}

// The new Symbol of the Token. Must be UTF-8 capitalized alphabetical string identifying the token.
func (transaction *TokenUpdateTransaction) SetTokenSymbol(symbol string) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Symbol = symbol
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTokenSymbol() string {
	return transaction.pb.GetSymbol()
}

// The new Name of the Token. Must be a string of ASCII characters.
func (transaction *TokenUpdateTransaction) SetTokenName(name string) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Name = name
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTokenName() string {
	return transaction.pb.GetName()
}

// The new Treasury account of the Token. If the provided treasury account is not existing or
// deleted, the response will be INVALID_TREASURY_ACCOUNT_FOR_TOKEN. If successful, the Token
// balance held in the previous Treasury Account is transferred to the new one.
func (transaction *TokenUpdateTransaction) SetTreasuryAccountID(treasury AccountID) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Treasury = treasury.toProtobuf()
	return transaction
}

func (transaction *TokenUpdateTransaction) GetTreasuryAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetTreasury())
}

// The new Admin key of the Token. If Token is immutable, transaction will resolve to
// TOKEN_IS_IMMUTABlE.
func (transaction *TokenUpdateTransaction) SetAdminKey(publicKey Key) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AdminKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TokenUpdateTransaction) GetAdminKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetAdminKey())
	if err != nil {
		return PublicKey{}
	}

	return key
}

// The new KYC key of the Token. If Token does not have currently a KYC key, transaction will
// resolve to TOKEN_HAS_NO_KYC_KEY.
func (transaction *TokenUpdateTransaction) SetKycKey(publicKey Key) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.KycKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TokenUpdateTransaction) GetKycKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetKycKey())
	if err != nil {
		return PublicKey{}
	}

	return key
}

// The new Freeze key of the Token. If the Token does not have currently a Freeze key, transaction
// will resolve to TOKEN_HAS_NO_FREEZE_KEY.
func (transaction *TokenUpdateTransaction) SetFreezeKey(publicKey Key) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.FreezeKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TokenUpdateTransaction) GetFreezeKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetFreezeKey())
	if err != nil {
		return PublicKey{}
	}

	return key
}

// The new Wipe key of the Token. If the Token does not have currently a Wipe key, transaction
// will resolve to TOKEN_HAS_NO_WIPE_KEY.
func (transaction *TokenUpdateTransaction) SetWipeKey(publicKey Key) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.WipeKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TokenUpdateTransaction) GetWipeKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetWipeKey())
	if err != nil {
		return PublicKey{}
	}

	return key
}

// The new Supply key of the Token. If the Token does not have currently a Supply key, transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
func (transaction *TokenUpdateTransaction) SetSupplyKey(publicKey Key) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.SupplyKey = publicKey.toProtoKey()
	return transaction
}

func (transaction *TokenUpdateTransaction) GetSupplyKey() Key {
	key, err := keyFromProtobuf(transaction.pb.GetSupplyKey())
	if err != nil {
		return PublicKey{}
	}

	return key
}

// The new account which will be automatically charged to renew the token's expiration, at
// autoRenewPeriod interval.
func (transaction *TokenUpdateTransaction) SetAutoRenewAccount(autoRenewAccount AccountID) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewAccount = autoRenewAccount.toProtobuf()
	return transaction
}

func (transaction *TokenUpdateTransaction) GetAutoRenewAccount() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetAutoRenewAccount())
}

// The new interval at which the auto-renew account will be charged to extend the token's expiry.
func (transaction *TokenUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.AutoRenewPeriod = &proto.Duration{Seconds: int64(autoRenewPeriod.Seconds())}
	return transaction
}

func (transaction *TokenUpdateTransaction) GetAutoRenewPeriod() time.Duration {
	return time.Duration(transaction.pb.GetAutoRenewPeriod().Seconds * time.Second.Nanoseconds())
}

// The new expiry time of the token. Expiry can be updated even if admin key is not set. If the
// provided expiry is earlier than the current token expiry, transaction wil resolve to
// INVALID_EXPIRATION_TIME
func (transaction *TokenUpdateTransaction) SetExpirationTime(expirationTime time.Time) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Expiry = &proto.Timestamp{
		Seconds: expirationTime.Unix(),
		Nanos:   int32(expirationTime.UnixNano()),
	}
	return transaction
}

func (transaction *TokenUpdateTransaction) GetExpirationTime() time.Time {
	return time.Unix(transaction.pb.GetExpiry().Seconds, int64(transaction.pb.GetExpiry().Nanos))
}

func (transaction *TokenUpdateTransaction) SetTokenMemo(memo string) *TokenUpdateTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Memo = &wrappers.StringValue{Value: memo}

	return transaction
}

func (transaction *TokenUpdateTransaction) GeTokenMemo() string {
	if transaction.pb.Memo != nil {
		return transaction.pb.Memo.GetValue()
	}

	return ""
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
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenUpdate{
			TokenUpdate: &proto.TokenUpdateTransactionBody{
				Token:            transaction.pb.GetToken(),
				Symbol:           transaction.pb.GetSymbol(),
				Name:             transaction.pb.GetName(),
				Treasury:         transaction.pb.GetTreasury(),
				AdminKey:         transaction.pb.GetAdminKey(),
				KycKey:           transaction.pb.GetKycKey(),
				FreezeKey:        transaction.pb.GetFreezeKey(),
				WipeKey:          transaction.pb.GetWipeKey(),
				SupplyKey:        transaction.pb.GetSupplyKey(),
				AutoRenewAccount: transaction.pb.GetAutoRenewAccount(),
				AutoRenewPeriod:  transaction.pb.GetAutoRenewPeriod(),
				Expiry:           transaction.pb.GetExpiry(),
				Memo:             transaction.pb.GetMemo(),
			},
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func tokenUpdateTransaction_getMethod(request request, channel *channel) method {
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
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		tokenUpdateTransaction_getMethod,
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

func (transaction *TokenUpdateTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenUpdate{
		TokenUpdate: transaction.pb,
	}

	return true
}

func (transaction *TokenUpdateTransaction) Freeze() (*TokenUpdateTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenUpdateTransaction) FreezeWith(client *Client) (*TokenUpdateTransaction, error) {
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
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
