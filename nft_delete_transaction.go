package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Deletes an already created Token.
// If no value is given for a field, that field is left unchanged. For an immutable tokens
// (that is, a token created without an adminKey), only the expiry may be deleted. Setting any
// other field in that case will cause the transaction status to resolve to TOKEN_IS_IMMUTABlE.
type NftDeleteTransaction struct {
	Transaction
	pb *proto.TokenDeleteTransactionBody
}

func NewNftDeleteTransaction() *NftDeleteTransaction {
	pb := &proto.TokenDeleteTransactionBody{}

	transaction := NftDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func nftDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenDeleteTransaction {
	return TokenDeleteTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenDeletion(),
	}
}

// The Token to be deleted
func (transaction *NftDeleteTransaction) SetTokenID(tokenID TokenID) *NftDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Token = tokenID.toProtobuf()
	return transaction
}

func (transaction *NftDeleteTransaction) GetTokenID() TokenID {
	return tokenIDFromProtobuf(transaction.pb.GetToken())
}

func (transaction *NftDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *NftDeleteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenDeletion{
			TokenDeletion: &proto.TokenDeleteTransactionBody{
				Token: transaction.pb.GetToken(),
			},
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func nftDeleteTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().DeleteToken,
	}
}

func (transaction *NftDeleteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *NftDeleteTransaction) Sign(
	privateKey PrivateKey,
) *NftDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *NftDeleteTransaction) SignWithOperator(
	client *Client,
) (*NftDeleteTransaction, error) {
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
func (transaction *NftDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *NftDeleteTransaction {
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
func (transaction *NftDeleteTransaction) Execute(
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
		nftDeleteTransaction_getMethod,
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

func (transaction *NftDeleteTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenDeletion{
		TokenDeletion: transaction.pb,
	}

	return true
}

func (transaction *NftDeleteTransaction) Freeze() (*NftDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *NftDeleteTransaction) FreezeWith(client *Client) (*NftDeleteTransaction, error) {
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

func (transaction *NftDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenDeleteTransaction.
func (transaction *NftDeleteTransaction) SetMaxTransactionFee(fee Hbar) *NftDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *NftDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenDeleteTransaction.
func (transaction *NftDeleteTransaction) SetTransactionMemo(memo string) *NftDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *NftDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenDeleteTransaction.
func (transaction *NftDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *NftDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *NftDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenDeleteTransaction.
func (transaction *NftDeleteTransaction) SetTransactionID(transactionID TransactionID) *NftDeleteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenDeleteTransaction.
func (transaction *NftDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *NftDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *NftDeleteTransaction) SetMaxRetry(count int) *NftDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *NftDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *NftDeleteTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
