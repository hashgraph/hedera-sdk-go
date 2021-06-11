package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Revokes KYC to the account for the given token. Must be signed by the Token's kycKey.
// If the provided account is not found, the transaction will resolve to INVALID_ACCOUNT_ID.
// If the provided account has been deleted, the transaction will resolve to ACCOUNT_DELETED.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If an Association between the provided token and account is not found, the transaction will
// resolve to TOKEN_NOT_ASSOCIATED_TO_ACCOUNT.
// If no KYC Key is defined, the transaction will resolve to TOKEN_HAS_NO_KYC_KEY.
// Once executed the Account is marked as KYC Revoked
type NftRevokeKycTransaction struct {
	Transaction
	pb *proto.TokenRevokeKycTransactionBody
}

func NewNftRevokeKycTransaction() *NftRevokeKycTransaction {
	pb := &proto.TokenRevokeKycTransactionBody{}

	transaction := NftRevokeKycTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func nftRevokeKycTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenRevokeKycTransaction {
	return TokenRevokeKycTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenRevokeKyc(),
	}
}

// The token for which this account will get his KYC revoked. If token does not exist, transaction results in INVALID_TOKEN_ID
func (transaction *NftRevokeKycTransaction) SetTokenID(tokenID TokenID) *NftRevokeKycTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Token = tokenID.toProtobuf()
	return transaction
}

func (transaction *NftRevokeKycTransaction) GetTokenID() TokenID {
	return tokenIDFromProtobuf(transaction.pb.GetToken())
}

// The account to be KYC Revoked
func (transaction *NftRevokeKycTransaction) SetAccountID(accountID AccountID) *NftRevokeKycTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Account = accountID.toProtobuf()
	return transaction
}

func (transaction *NftRevokeKycTransaction) GetAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.Account)
}

func (transaction *NftRevokeKycTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *NftRevokeKycTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenRevokeKyc{
			TokenRevokeKyc: &proto.TokenRevokeKycTransactionBody{
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

func nftRevokeKycTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().RevokeKycFromTokenAccount,
	}
}

func (transaction *NftRevokeKycTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *NftRevokeKycTransaction) Sign(
	privateKey PrivateKey,
) *NftRevokeKycTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *NftRevokeKycTransaction) SignWithOperator(
	client *Client,
) (*NftRevokeKycTransaction, error) {
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
func (transaction *NftRevokeKycTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *NftRevokeKycTransaction {
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
func (transaction *NftRevokeKycTransaction) Execute(
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
		nftRevokeKycTransaction_getMethod,
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

func (transaction *NftRevokeKycTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenRevokeKyc{
		TokenRevokeKyc: transaction.pb,
	}

	return true
}

func (transaction *NftRevokeKycTransaction) Freeze() (*NftRevokeKycTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *NftRevokeKycTransaction) FreezeWith(client *Client) (*NftRevokeKycTransaction, error) {
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

func (transaction *NftRevokeKycTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenRevokeKycTransaction.
func (transaction *NftRevokeKycTransaction) SetMaxTransactionFee(fee Hbar) *NftRevokeKycTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *NftRevokeKycTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenRevokeKycTransaction.
func (transaction *NftRevokeKycTransaction) SetTransactionMemo(memo string) *NftRevokeKycTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *NftRevokeKycTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenRevokeKycTransaction.
func (transaction *NftRevokeKycTransaction) SetTransactionValidDuration(duration time.Duration) *NftRevokeKycTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *NftRevokeKycTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenRevokeKycTransaction.
func (transaction *NftRevokeKycTransaction) SetTransactionID(transactionID TransactionID) *NftRevokeKycTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenRevokeKycTransaction.
func (transaction *NftRevokeKycTransaction) SetNodeAccountIDs(nodeID []AccountID) *NftRevokeKycTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *NftRevokeKycTransaction) SetMaxRetry(count int) *NftRevokeKycTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *NftRevokeKycTransaction) AddSignature(publicKey PublicKey, signature []byte) *NftRevokeKycTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
