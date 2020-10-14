package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
	"time"
)

type LiveHashAddTransaction struct {
	Transaction
	pb *proto.CryptoAddLiveHashTransactionBody
}

func NewLiveHashAddTransaction() *LiveHashAddTransaction {
	pb := &proto.CryptoAddLiveHashTransactionBody{}

	transaction := LiveHashAddTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	return &transaction
}

func (transaction *LiveHashAddTransaction) SetHash(hash []byte) *LiveHashAddTransaction {
	transaction.pb.LiveHash.Hash = hash
	return transaction
}

func (transaction *LiveHashAddTransaction) GetHash() []byte {
	return transaction.pb.GetLiveHash().GetHash()
}

func (transaction *LiveHashAddTransaction) SetKeys(keys ...Key) *LiveHashAddTransaction {
	if transaction.pb.LiveHash.Keys == nil {
		transaction.pb.LiveHash.Keys = &proto.KeyList{Keys: []*proto.Key{}}
	}
	keyList := KeyList{keys: []*proto.Key{}}
	keyList.AddAll(keys)

	transaction.pb.LiveHash.Keys = keyList.toProtoKeyList()

	return transaction
}

func (transaction *LiveHashAddTransaction) GetKeys() KeyList {
	return keyListFromProto(transaction.pb.GetLiveHash().GetKeys())
}

func (transaction *LiveHashAddTransaction) SetDuration(duration time.Duration) *LiveHashAddTransaction {
	transaction.pb.LiveHash.Duration = durationToProto(duration)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetDuration() time.Duration {
	return durationFromProto(transaction.pb.GetLiveHash().GetDuration())
}

func (transaction *LiveHashAddTransaction) SetAccountID(accountID AccountID) *LiveHashAddTransaction {
	transaction.pb.LiveHash.AccountId = accountID.toProtobuf()
	return transaction
}

func (transaction *LiveHashAddTransaction) GetAccountID() AccountID {
	return accountIDFromProto(transaction.pb.LiveHash.GetAccountId())
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func liveHashAddTransaction_getMethod(channel *channel) method {
	return method{
		transaction: channel.getCrypto().CreateAccount,
	}
}

func (transaction *LiveHashAddTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *LiveHashAddTransaction) Sign(
	privateKey PrivateKey,
) *LiveHashAddTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *LiveHashAddTransaction) SignWithOperator(
	client *Client,
) (*LiveHashAddTransaction, error) {
	// If the transaction is not signed by the operator, we need
	// to sign the transaction with the operator

	if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *LiveHashAddTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *LiveHashAddTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	for index := 0; index < len(transaction.transactions); index++ {
		signature := signer(transaction.transactions[index].GetBodyBytes())

		transaction.signatures[index].SigPair = append(
			transaction.signatures[index].SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *LiveHashAddTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	transactionID := transaction.id

	if !client.GetOperatorID().isZero() && client.GetOperatorID().equals(transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorKey(),
			client.operator.signer,
		)
	}

	_, err := execute(
		client,
		request{
			transaction: &transaction.Transaction,
		},
		transaction_shouldRetry,
		transaction_makeRequest,
		transaction_advanceRequest,
		transaction_getNodeId,
		liveHashAddTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{TransactionID: transaction.id}, nil
}

func (transaction *LiveHashAddTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_CryptoAddLiveHash{
		CryptoAddLiveHash: transaction.pb,
	}

	return true
}

func (transaction *LiveHashAddTransaction) Freeze() (*LiveHashAddTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *LiveHashAddTransaction) FreezeWith(client *Client) (*LiveHashAddTransaction, error) {
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}

	if !transaction.onFreeze(transaction.pbBody) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

func (transaction *LiveHashAddTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetMaxTransactionFee(fee Hbar) *LiveHashAddTransaction {
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetTransactionMemo(memo string) *LiveHashAddTransaction {
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetTransactionValidDuration(duration time.Duration) *LiveHashAddTransaction {
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetTransactionID(transactionID TransactionID) *LiveHashAddTransaction {
	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetNodeID() AccountID {
	return transaction.Transaction.GetNodeID()
}

// SetNodeID sets the node AccountID for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetNodeID(nodeID AccountID) *LiveHashAddTransaction {
	transaction.Transaction.SetNodeID(nodeID)
	return transaction
}
