package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Wipes the provided amount of tokens from the specified Account. Must be signed by the Token's
// Wipe key.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If Wipe Key is not present in the Token, transaction results in TOKEN_HAS_NO_WIPE_KEY.
// If the provided account is the Token's Treasury Account, transaction results in
// CANNOT_WIPE_TOKEN_TREASURY_ACCOUNT
// On success, tokens are removed from the account and the total supply of the token is decreased
// by the wiped amount.
//
// The amount provided is in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to wipe 100 tokens from account, one must provide amount of
// 10000. In order to wipe 100.55 tokens, one must provide amount of 10055.
type NftWipeTransaction struct {
	Transaction
	pb *proto.TokenWipeAccountTransactionBody
}

func NewNftWipeTransaction() *NftWipeTransaction {
	pb := &proto.TokenWipeAccountTransactionBody{}

	transaction := NftWipeTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func nftWipeTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenWipeTransaction {
	return TokenWipeTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenWipe(),
	}
}

// The token for which the account will be wiped. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (transaction *NftWipeTransaction) SetTokenID(tokenID TokenID) *NftWipeTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Token = tokenID.toProtobuf()
	return transaction
}

func (transaction *NftWipeTransaction) GetTokenID() TokenID {
	return tokenIDFromProtobuf(transaction.pb.GetToken())
}

// The account to be wiped
func (transaction *NftWipeTransaction) SetAccountID(accountID AccountID) *NftWipeTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Account = accountID.toProtobuf()
	return transaction
}

func (transaction *NftWipeTransaction) GetAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetAccount())
}

func (transaction *NftWipeTransaction) GetSerialNumbers() []int64 {
	return transaction.pb.GetSerialNumbers()
}

func (transaction *NftWipeTransaction) SetSerialNumbers(serial []int64) *NftWipeTransaction {
	transaction.requireNotFrozen()
	transaction.pb.SerialNumbers = serial
	return transaction
}

func (transaction *NftWipeTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *NftWipeTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenWipe{
			TokenWipe: &proto.TokenWipeAccountTransactionBody{
				Token:   transaction.pb.GetToken(),
				Account: transaction.pb.GetAccount(),
				Amount:  transaction.pb.GetAmount(),
			},
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func nftWipeTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().WipeTokenAccount,
	}
}

func (transaction *NftWipeTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *NftWipeTransaction) Sign(
	privateKey PrivateKey,
) *NftWipeTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *NftWipeTransaction) SignWithOperator(
	client *Client,
) (*NftWipeTransaction, error) {
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
func (transaction *NftWipeTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *NftWipeTransaction {
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
func (transaction *NftWipeTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
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
		nftWipeTransaction_getMethod,
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

func (transaction *NftWipeTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenWipe{
		TokenWipe: transaction.pb,
	}

	return true
}

func (transaction *NftWipeTransaction) Freeze() (*NftWipeTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *NftWipeTransaction) FreezeWith(client *Client) (*NftWipeTransaction, error) {
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

func (transaction *NftWipeTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenWipeTransaction.
func (transaction *NftWipeTransaction) SetMaxTransactionFee(fee Hbar) *NftWipeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *NftWipeTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenWipeTransaction.
func (transaction *NftWipeTransaction) SetTransactionMemo(memo string) *NftWipeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *NftWipeTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenWipeTransaction.
func (transaction *NftWipeTransaction) SetTransactionValidDuration(duration time.Duration) *NftWipeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *NftWipeTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenWipeTransaction.
func (transaction *NftWipeTransaction) SetTransactionID(transactionID TransactionID) *NftWipeTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenWipeTransaction.
func (transaction *NftWipeTransaction) SetNodeAccountIDs(nodeID []AccountID) *NftWipeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *NftWipeTransaction) SetMaxRetry(count int) *NftWipeTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *NftWipeTransaction) AddSignature(publicKey PublicKey, signature []byte) *NftWipeTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
