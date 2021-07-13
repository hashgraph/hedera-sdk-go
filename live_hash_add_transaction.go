package hedera

import (
	"github.com/hashgraph/hedera-protobufs-go/services"
	"github.com/pkg/errors"

	"time"
)

type LiveHashAddTransaction struct {
	Transaction
	pb        *services.CryptoAddLiveHashTransactionBody
	accountID AccountID
}

func NewLiveHashAddTransaction() *LiveHashAddTransaction {
	pb := &services.CryptoAddLiveHashTransactionBody{LiveHash: &services.LiveHash{}}

	transaction := LiveHashAddTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func liveHashAddTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) LiveHashAddTransaction {
	return LiveHashAddTransaction{
		Transaction: transaction,
		pb:          pb.GetCryptoAddLiveHash(),
		accountID:   accountIDFromProtobuf(pb.GetCryptoAddLiveHash().GetLiveHash().GetAccountId(), nil),
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
		transaction.pb.LiveHash.Keys = &services.KeyList{Keys: []*services.Key{}}
	}
	keyList := NewKeyList()
	keyList.AddAll(keys)

	transaction.pb.LiveHash.Keys = keyList.toProtoKeyList()

	return transaction
}

func (transaction *LiveHashAddTransaction) GetKeys() KeyList {
	keys := transaction.pb.GetLiveHash().GetKeys()
	if keys != nil {
		keyList, err := keyListFromProtobuf(keys, nil)
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

func (transaction *LiveHashAddTransaction) SetAccountID(id AccountID) *LiveHashAddTransaction {
	transaction.requireNotFrozen()
	transaction.accountID = id
	return transaction
}

func (transaction *LiveHashAddTransaction) GetAccountID() AccountID {
	return transaction.accountID
}

func (transaction *LiveHashAddTransaction) validateNetworkOnIDs(client *Client) error {
	var err error
	err = transaction.accountID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *LiveHashAddTransaction) build() *LiveHashAddTransaction {
	if !transaction.accountID.isZero() {
		transaction.pb.LiveHash.AccountId = transaction.accountID.toProtobuf()
	}

	return transaction
}

func (transaction *LiveHashAddTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `LiveHashAddTransaction`")
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
func (transaction *LiveHashAddTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *LiveHashAddTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *LiveHashAddTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
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
		liveHashAddTransaction_getMethod,
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

func (transaction *LiveHashAddTransaction) onFreeze(
	pbBody *services.TransactionBody,
) bool {
	pbBody.Data = &services.TransactionBody_CryptoAddLiveHash{
		CryptoAddLiveHash: transaction.pb,
	}

	return true
}

func (transaction *LiveHashAddTransaction) Freeze() (*LiveHashAddTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *LiveHashAddTransaction) FreezeWith(client *Client) (*LiveHashAddTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &LiveHashAddTransaction{}, err
	}
	transaction.build()

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

// SetNodeAccountID sets the node AccountID for this LiveHashAddTransaction.
func (transaction *LiveHashAddTransaction) SetNodeAccountIDs(nodeID []AccountID) *LiveHashAddTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *LiveHashAddTransaction) SetMaxRetry(count int) *LiveHashAddTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *LiveHashAddTransaction) AddSignature(publicKey PublicKey, signature []byte) *LiveHashAddTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
