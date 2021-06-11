package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Burns tokens from the Token's treasury Account. If no Supply Key is defined, the transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
// The operation decreases the Total Supply of the Token. Total supply cannot go below
// zero.
// The amount provided must be in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to burn 100 tokens, one must provide amount of 10000. In order
// to burn 100.55 tokens, one must provide amount of 10055.
type NftBurnTransaction struct {
	Transaction
	pb *proto.TokenBurnTransactionBody
}

func NewNftBurnTransaction() *NftBurnTransaction {
	pb := &proto.TokenBurnTransactionBody{}

	transaction := NftBurnTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func nftBurnTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenBurnTransaction {
	return TokenBurnTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenBurn(),
	}
}

// The token for which to burn tokens. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (transaction *NftBurnTransaction) SetTokenID(tokenID TokenID) *NftBurnTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Token = tokenID.toProtobuf()
	return transaction
}

func (transaction *NftBurnTransaction) GetTokenID() TokenID {
	return tokenIDFromProtobuf(transaction.pb.GetToken())
}

func (transaction *NftBurnTransaction) SetSerialNumbers(serial []int64) *NftBurnTransaction {
	transaction.requireNotFrozen()
	transaction.pb.SerialNumbers = serial
	return transaction
}

func (transaction *NftBurnTransaction) GetSerialNumbers() []int64 {
	return transaction.pb.GetSerialNumbers()
}

func (transaction *NftBurnTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *NftBurnTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenBurn{
			TokenBurn: &proto.TokenBurnTransactionBody{
				Token:  transaction.pb.GetToken(),
				Amount: transaction.pb.GetAmount(),
			},
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func nftBurnTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().BurnToken,
	}
}

func (transaction *NftBurnTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *NftBurnTransaction) Sign(
	privateKey PrivateKey,
) *NftBurnTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *NftBurnTransaction) SignWithOperator(
	client *Client,
) (*NftBurnTransaction, error) {
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
func (transaction *NftBurnTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *NftBurnTransaction {
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
func (transaction *NftBurnTransaction) Execute(
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

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
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
		nftBurnTransaction_getMethod,
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

func (transaction *NftBurnTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenBurn{
		TokenBurn: transaction.pb,
	}

	return true
}

func (transaction *NftBurnTransaction) Freeze() (*NftBurnTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *NftBurnTransaction) FreezeWith(client *Client) (*NftBurnTransaction, error) {
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

func (transaction *NftBurnTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenBurnTransaction.
func (transaction *NftBurnTransaction) SetMaxTransactionFee(fee Hbar) *NftBurnTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *NftBurnTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenBurnTransaction.
func (transaction *NftBurnTransaction) SetTransactionMemo(memo string) *NftBurnTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *NftBurnTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenBurnTransaction.
func (transaction *NftBurnTransaction) SetTransactionValidDuration(duration time.Duration) *NftBurnTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *NftBurnTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenBurnTransaction.
func (transaction *NftBurnTransaction) SetTransactionID(transactionID TransactionID) *NftBurnTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenBurnTransaction.
func (transaction *NftBurnTransaction) SetNodeAccountIDs(nodeID []AccountID) *NftBurnTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *NftBurnTransaction) SetMaxRetry(count int) *NftBurnTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *NftBurnTransaction) AddSignature(publicKey PublicKey, signature []byte) *NftBurnTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
