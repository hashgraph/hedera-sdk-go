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
type NftUnfreezeTransaction struct {
	Transaction
	pb *proto.TokenUnfreezeAccountTransactionBody
}

func NewNftUnfreezeTransaction() *NftUnfreezeTransaction {
	pb := &proto.TokenUnfreezeAccountTransactionBody{}

	transaction := NftUnfreezeTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func nftUnfreezeTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenUnfreezeTransaction {
	return TokenUnfreezeTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenUnfreeze(),
	}
}

// The token for which this account will be unfrozen. If token does not exist, transaction results in INVALID_TOKEN_ID
func (transaction *NftUnfreezeTransaction) SetTokenID(tokenID TokenID) *NftUnfreezeTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Token = tokenID.toProtobuf()
	return transaction
}

func (transaction *NftUnfreezeTransaction) GetTokenID() TokenID {
	return tokenIDFromProtobuf(transaction.pb.Token)
}

// The account to be unfrozen
func (transaction *NftUnfreezeTransaction) SetAccountID(accountID AccountID) *NftUnfreezeTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Account = accountID.toProtobuf()
	return transaction
}

func (transaction *NftUnfreezeTransaction) GetAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetAccount())
}

func (transaction *NftUnfreezeTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *NftUnfreezeTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenUnfreeze{
			TokenUnfreeze: &proto.TokenUnfreezeAccountTransactionBody{
				Token:   transaction.pb.GetToken(),
				Account: transaction.pb.GetAccount(),
			},
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func nftUnfreezeTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().UnfreezeTokenAccount,
	}
}

func (transaction *NftUnfreezeTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *NftUnfreezeTransaction) Sign(
	privateKey PrivateKey,
) *NftUnfreezeTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *NftUnfreezeTransaction) SignWithOperator(
	client *Client,
) (*NftUnfreezeTransaction, error) {
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
func (transaction *NftUnfreezeTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *NftUnfreezeTransaction {
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
func (transaction *NftUnfreezeTransaction) Execute(
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
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		nftUnfreezeTransaction_getMethod,
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

func (transaction *NftUnfreezeTransaction) onUnfreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenUnfreeze{
		TokenUnfreeze: transaction.pb,
	}

	return true
}

func (transaction *NftUnfreezeTransaction) Unfreeze() (*NftUnfreezeTransaction, error) {
	return transaction.UnfreezeWith(nil)
}

func (transaction *NftUnfreezeTransaction) UnfreezeWith(client *Client) (*NftUnfreezeTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onUnfreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *NftUnfreezeTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenUnfreezeTransaction.
func (transaction *NftUnfreezeTransaction) SetMaxTransactionFee(fee Hbar) *NftUnfreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *NftUnfreezeTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenUnfreezeTransaction.
func (transaction *NftUnfreezeTransaction) SetTransactionMemo(memo string) *NftUnfreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *NftUnfreezeTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenUnfreezeTransaction.
func (transaction *NftUnfreezeTransaction) SetTransactionValidDuration(duration time.Duration) *NftUnfreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *NftUnfreezeTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenUnfreezeTransaction.
func (transaction *NftUnfreezeTransaction) SetTransactionID(transactionID TransactionID) *NftUnfreezeTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenUnfreezeTransaction.
func (transaction *NftUnfreezeTransaction) SetNodeAccountIDs(nodeID []AccountID) *NftUnfreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *NftUnfreezeTransaction) SetMaxRetry(count int) *NftUnfreezeTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *NftUnfreezeTransaction) AddSignature(publicKey PublicKey, signature []byte) *NftUnfreezeTransaction {
	if !transaction.IsFrozen() {
		transaction.Unfreeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
