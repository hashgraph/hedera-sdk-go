package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// Mints tokens from the Token's treasury Account. If no Supply Key is defined, the transaction
// will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
// The operation decreases the Total Supply of the Token. Total supply cannot go below
// zero.
// The amount provided must be in the lowest denomination possible. Example:
// Token A has 2 decimals. In order to mint 100 tokens, one must provide amount of 10000. In order
// to mint 100.55 tokens, one must provide amount of 10055.
type NftMintTransaction struct {
	Transaction
	pb *proto.TokenMintTransactionBody
}

func NewNftMintTransaction() *NftMintTransaction {
	pb := &proto.TokenMintTransactionBody{}

	transaction := NftMintTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(30))

	return &transaction
}

func nftMintTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) TokenMintTransaction {
	return TokenMintTransaction{
		Transaction: transaction,
		pb:          pb.GetTokenMint(),
	}
}

// The token for which to mint tokens. If token does not exist, transaction results in
// INVALID_TOKEN_ID
func (transaction *NftMintTransaction) SetTokenID(tokenID TokenID) *NftMintTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Token = tokenID.toProtobuf()
	return transaction
}

func (transaction *NftMintTransaction) GetTokenID() TokenID {
	return tokenIDFromProtobuf(transaction.pb.Token)
}

func (transaction *NftMintTransaction) SetMetadata(meta [][]byte) *NftMintTransaction {
	transaction.requireNotFrozen()
	transaction.pb.Metadata = meta
	return transaction
}

func (transaction *NftMintTransaction) GetMetadata() [][]byte {
	return transaction.pb.GetMetadata()
}

func (transaction *NftMintTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *NftMintTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_TokenMint{
			TokenMint: &proto.TokenMintTransactionBody{
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

func nftMintTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getToken().MintToken,
	}
}

func (transaction *NftMintTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *NftMintTransaction) Sign(
	privateKey PrivateKey,
) *NftMintTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *NftMintTransaction) SignWithOperator(
	client *Client,
) (*NftMintTransaction, error) {
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
func (transaction *NftMintTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *NftMintTransaction {
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
func (transaction *NftMintTransaction) Execute(
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
		nftMintTransaction_getMethod,
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

func (transaction *NftMintTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_TokenMint{
		TokenMint: transaction.pb,
	}

	return true
}

func (transaction *NftMintTransaction) Freeze() (*NftMintTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *NftMintTransaction) FreezeWith(client *Client) (*NftMintTransaction, error) {
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

func (transaction *NftMintTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this TokenMintTransaction.
func (transaction *NftMintTransaction) SetMaxTransactionFee(fee Hbar) *NftMintTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *NftMintTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this TokenMintTransaction.
func (transaction *NftMintTransaction) SetTransactionMemo(memo string) *NftMintTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *NftMintTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this TokenMintTransaction.
func (transaction *NftMintTransaction) SetTransactionValidDuration(duration time.Duration) *NftMintTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *NftMintTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this TokenMintTransaction.
func (transaction *NftMintTransaction) SetTransactionID(transactionID TransactionID) *NftMintTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeTokenID sets the node TokenID for this TokenMintTransaction.
func (transaction *NftMintTransaction) SetNodeAccountIDs(nodeID []AccountID) *NftMintTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *NftMintTransaction) SetMaxRetry(count int) *NftMintTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *NftMintTransaction) AddSignature(publicKey PublicKey, signature []byte) *NftMintTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
