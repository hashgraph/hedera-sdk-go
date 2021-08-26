package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

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
	Transaction
	tokenID   TokenID
	accountID AccountID
}

func NewTokenGrantKycTransaction() *TokenGrantKycTransaction {
	transaction := TokenGrantKycTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func tokenGrantKycTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenGrantKycTransaction {
	return TokenGrantKycTransaction{
		Transaction: transaction,
		tokenID:     tokenIDFromProtobuf(pb.GetTokenGrantKyc().GetToken()),
		accountID:   accountIDFromProtobuf(pb.GetTokenGrantKyc().GetAccount()),
	}
}

// The token for which this account will be granted KYC. If token does not exist, transaction results in INVALID_TOKEN_ID
func (transaction *TokenGrantKycTransaction) SetTokenID(id TokenID) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()
	transaction.tokenID = id
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetTokenID() TokenID {
	return transaction.tokenID
}

// The account to be KYCed
func (transaction *TokenGrantKycTransaction) SetAccountID(id AccountID) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()
	transaction.accountID = id
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetAccountID() AccountID {
	return transaction.accountID
}

func (transaction *TokenGrantKycTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = transaction.tokenID.Validate(client)
	if err != nil {
		return err
	}
	err = transaction.accountID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *TokenGrantKycTransaction) build() *proto.TransactionBody {
	body := &proto.TokenGrantKycTransactionBody{}
	if !transaction.tokenID.isZero() {
		body.Token = transaction.tokenID.toProtobuf()
	}

	if !transaction.accountID.isZero() {
		body.Account = transaction.accountID.toProtobuf()
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_TokenGrantKyc{
			TokenGrantKyc: body,
		},
	}
}

func (transaction *TokenGrantKycTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenGrantKycTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.TokenGrantKycTransactionBody{}
	if !transaction.tokenID.isZero() {
		body.Token = transaction.tokenID.toProtobuf()
	}

	if !transaction.accountID.isZero() {
		body.Account = transaction.accountID.toProtobuf()
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_TokenGrantKyc{
			TokenGrantKyc: body,
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func tokenGrantKycTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().GrantKycToTokenAccount,
	}
}

func (transaction *TokenGrantKycTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenGrantKycTransaction) Sign(
	privateKey PrivateKey,
) *TokenGrantKycTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenGrantKycTransaction) SignWithOperator(
	client *Client,
) (*TokenGrantKycTransaction, error) {
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
func (transaction *TokenGrantKycTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenGrantKycTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *TokenGrantKycTransaction) Execute(
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
		transaction_makeRequest(request{
			transaction: &transaction.Transaction,
		}),
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		tokenGrantKycTransaction_getMethod,
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

func (transaction *TokenGrantKycTransaction) Freeze() (*TokenGrantKycTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *TokenGrantKycTransaction) FreezeWith(client *Client) (*TokenGrantKycTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TokenGrantKycTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenGrantKycTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetMaxTransactionFee(fee Hbar) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetTransactionMemo(memo string) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetTransactionValidDuration(duration time.Duration) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetTransactionID(transactionID TransactionID) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenGrantKycTransaction.
func (transaction *TokenGrantKycTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenGrantKycTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenGrantKycTransaction) SetMaxRetry(count int) *TokenGrantKycTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenGrantKycTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenGrantKycTransaction {
	transaction.requireOneNodeAccountID()

	if !transaction.isFrozen() {
		transaction.Freeze()
	}

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

	//transaction.signedTransactions[0].SigMap.SigPair = append(transaction.signedTransactions[0].SigMap.SigPair, publicKey.toSignaturePairProtobuf(signature))
	return transaction
}

func (transaction *TokenGrantKycTransaction) SetMaxBackoff(max time.Duration) *TokenGrantKycTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TokenGrantKycTransaction) SetMinBackoff(min time.Duration) *TokenGrantKycTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TokenGrantKycTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
