package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

type LiveHashAddTransaction struct {
	Transaction
	pb *proto.CryptoAddLiveHashTransactionBody
}

func NewLiveHashAddTransaction() *LiveHashAddTransaction {
	pb := &proto.CryptoAddLiveHashTransactionBody{LiveHash: &proto.LiveHash{}}

	transaction := LiveHashAddTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	return &transaction
}

func liveHashAddTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) LiveHashAddTransaction {
	return LiveHashAddTransaction{
		Transaction: transaction,
		pb:          pb.GetCryptoAddLiveHash(),
	}
}

func (transaction *LiveHashAddTransaction) SetHash(hash []byte) *LiveHashAddTransaction {
	transaction.requireNotFrozen()
	transaction.pb.LiveHash.Hash = hash
	return transaction
}

func (transaction *LiveHashAddTransaction) GetHash() []byte {
	return transaction.pb.GetLiveHash().GetHash()
}

func (transaction *LiveHashAddTransaction) SetKeys(keys ...Key) *LiveHashAddTransaction {
	transaction.requireNotFrozen()
	if transaction.pb.LiveHash.Keys == nil {
		transaction.pb.LiveHash.Keys = &proto.KeyList{Keys: []*proto.Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	transaction.pb.LiveHash.Keys = keyList.toProtoKeyList()

	return transaction
}

func (transaction *LiveHashAddTransaction) GetKeys() KeyList {
	keys := transaction.pb.GetLiveHash().GetKeys()
	if keys != nil {
		keyList, err := keyListFromProtobuf(keys)
		if err != nil {
			return KeyList{}
		}

		return keyList
	} else {
		return KeyList{}
	}
}

func (transaction *LiveHashAddTransaction) SetDuration(duration time.Duration) *LiveHashAddTransaction {
	transaction.requireNotFrozen()
	transaction.pb.LiveHash.Duration = durationToProtobuf(duration)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetDuration() time.Duration {
	return durationFromProtobuf(transaction.pb.GetLiveHash().GetDuration())
}

func (transaction *LiveHashAddTransaction) SetAccountID(accountID AccountID) *LiveHashAddTransaction {
	transaction.requireNotFrozen()
	transaction.pb.LiveHash.AccountId = accountID.toProtobuf()
	return transaction
}

func (transaction *LiveHashAddTransaction) GetAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.LiveHash.GetAccountId())
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func liveHashAddTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getCrypto().AddLiveHash,
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
func (transaction *LiveHashAddTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if !transaction.IsFrozen() {
		transaction.FreezeWith(client)
	}

	transactionID := transaction.transactionIDs[0]

	if !client.GetOperatorAccountID().isZero() && client.GetOperatorAccountID().equals(transactionID.AccountID) {
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
		liveHashAddTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{
		TransactionID: transaction.transactionIDs[0],
		NodeID:        resp.transaction.NodeID,
	}, nil
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
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetTransactionMemo(memo string) *LiveHashAddTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetTransactionValidDuration(duration time.Duration) *LiveHashAddTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetTransactionID(transactionID TransactionID) *LiveHashAddTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

func (transaction *LiveHashAddTransaction) GetNodeAccountIDs() []AccountID {
	return transaction.Transaction.GetNodeAccountIDs()
}

// SetNodeAccountID sets the node AccountID for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetNodeAccountIDs(nodeID []AccountID) *LiveHashAddTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}
