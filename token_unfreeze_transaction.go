package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Unfreezes transfers of the specified token for the account. Must be signed by the Token's freezeKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no Freeze Key is defined, the transaction will resolve to TOKEN_HAS_NO_FREEZE_KEY.
// Once executed the Account is marked as Unfrozen and will be able to receive or send tokens. The
// operation is idempotent.
type TokenUnfreezeTransaction struct {
	Transaction
	tokenID   TokenID
	accountID AccountID
}

func NewTokenUnfreezeTransaction() *TokenUnfreezeTransaction {
	transaction := TokenUnfreezeTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func tokenUnfreezeTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenUnfreezeTransaction {
	return TokenUnfreezeTransaction{
		Transaction: transaction,
		tokenID:     tokenIDFromProtobuf(pb.GetTokenUnfreeze().GetToken()),
		accountID:   accountIDFromProtobuf(pb.GetTokenUnfreeze().GetAccount()),
	}
}

// The token for which this account will be unfrozen. If token does not exist, transaction results in INVALID_TOKEN_ID
func (transaction *TokenUnfreezeTransaction) SetTokenID(id TokenID) *TokenUnfreezeTransaction {
	transaction.requireNotFrozen()
	transaction.tokenID = id
	return transaction
}

func (transaction *TokenUnfreezeTransaction) GetTokenID() TokenID {
	return transaction.tokenID
}

// The account to be unfrozen
func (transaction *TokenUnfreezeTransaction) SetAccountID(id AccountID) *TokenUnfreezeTransaction {
	transaction.requireNotFrozen()
	transaction.accountID = id
	return transaction
}

func (transaction *TokenUnfreezeTransaction) GetAccountID() AccountID {
	return transaction.accountID
}

func (transaction *TokenUnfreezeTransaction) validateNetworkOnIDs(client *Client) error {
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

func (transaction *TokenUnfreezeTransaction) build() *proto.TransactionBody {
	body := &proto.TokenUnfreezeAccountTransactionBody{}
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
		Data: &proto.TransactionBody_TokenUnfreeze{
			TokenUnfreeze: body,
		},
	}
}

func (transaction *TokenUnfreezeTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *TokenUnfreezeTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.TokenUnfreezeAccountTransactionBody{}
	if !transaction.tokenID.isZero() {
		body.Token = transaction.tokenID.toProtobuf()
	}

	if !transaction.accountID.isZero() {
		body.Account = transaction.accountID.toProtobuf()
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_TokenUnfreeze{
			TokenUnfreeze: body,
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func tokenUnfreezeTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().UnfreezeTokenAccount,
	}
}

func (transaction *TokenUnfreezeTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *TokenUnfreezeTransaction) Sign(
	privateKey PrivateKey,
) *TokenUnfreezeTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *TokenUnfreezeTransaction) SignWithOperator(
	client *Client,
) (*TokenUnfreezeTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	//if !transaction.IsFrozen() {
	//	_, err := transaction.FreezeWith(client)
	//	if err != nil {
	//		return transaction, err
	//	}
	//}
	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *TokenUnfreezeTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenUnfreezeTransaction {
	//if !transaction.IsFrozen() {
	//	transaction.Unfreeze()
	//}

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
func (transaction *TokenUnfreezeTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	if !transaction.IsFrozen() {
		transaction.UnfreezeWith(client)
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
		tokenUnfreezeTransaction_getMethod,
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

func (transaction *TokenUnfreezeTransaction) Unfreeze() (*TokenUnfreezeTransaction, error) {
	return transaction.UnfreezeWith(nil)
}

func (transaction *TokenUnfreezeTransaction) UnfreezeWith(client *Client) (*TokenUnfreezeTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &TokenUnfreezeTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
}

func (transaction *TokenUnfreezeTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenUnfreezeTransaction.
func (transaction *TokenUnfreezeTransaction) SetMaxTransactionFee(fee Hbar) *TokenUnfreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *TokenUnfreezeTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenUnfreezeTransaction.
func (transaction *TokenUnfreezeTransaction) SetTransactionMemo(memo string) *TokenUnfreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *TokenUnfreezeTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenUnfreezeTransaction.
func (transaction *TokenUnfreezeTransaction) SetTransactionValidDuration(duration time.Duration) *TokenUnfreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *TokenUnfreezeTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenUnfreezeTransaction.
func (transaction *TokenUnfreezeTransaction) SetTransactionID(transactionID TransactionID) *TokenUnfreezeTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenUnfreezeTransaction.
func (transaction *TokenUnfreezeTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenUnfreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *TokenUnfreezeTransaction) SetMaxRetry(count int) *TokenUnfreezeTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *TokenUnfreezeTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenUnfreezeTransaction {
	transaction.requireOneNodeAccountID()

	if !transaction.isFrozen() {
		transaction.Unfreeze()
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

func (transaction *TokenUnfreezeTransaction) SetMaxBackoff(max time.Duration) *TokenUnfreezeTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *TokenUnfreezeTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *TokenUnfreezeTransaction) SetMinBackoff(min time.Duration) *TokenUnfreezeTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *TokenUnfreezeTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
