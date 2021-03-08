package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

// AccountDeleteTransaction creates a new account. After the account is created, the AccountID for it is in the receipt,
// or by asking for a Record of the transaction to be created, and retrieving that. The account can then automatically
// generate records for large transfers into it or out of it, which each last for 25 hours. Records are generated for
// any transfer that exceeds the thresholds given here. This account is charged hbar for each record generated, so the
// thresholds are useful for limiting Record generation to happen only for large transactions.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0,
// with a null key. Future versions of the API will support multiple realms and multiple shards.
type AccountDeleteTransaction struct {
	Transaction
	pb *proto.CryptoDeleteTransactionBody
}

func accountDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) AccountDeleteTransaction {
	return AccountDeleteTransaction{
		Transaction: transaction,
		pb:          pb.GetCryptoDelete(),
	}
}

func NewAccountDeleteTransaction() *AccountDeleteTransaction {
	pb := &proto.CryptoDeleteTransactionBody{}

	transaction := AccountDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

// SetNodeAccountID sets the node AccountID for this AccountCreateTransaction.
func (transaction *AccountDeleteTransaction) SetAccountID(accountID AccountID) *AccountDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.DeleteAccountID = accountID.toProtobuf()
	return transaction
}

func (transaction *AccountDeleteTransaction) GetAccountID() AccountID {
	return accountIDFromProtobuf(transaction.pb.GetDeleteAccountID())
}

// SetTransferAccountID sets the AccountID which will receive all remaining hbars.
func (transaction *AccountDeleteTransaction) SetTransferAccountID(transferAccountID AccountID) *AccountDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.pb.TransferAccountID = transferAccountID.toProtobuf()
	return transaction
}

func (transaction *AccountDeleteTransaction) GetTransferAccountID(transferAccountID AccountID) AccountID {
	return accountIDFromProtobuf(transaction.pb.GetTransferAccountID())
}

func (transaction *AccountDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	body := &proto.TransactionBody{
		TransactionID:            transaction.pbBody.GetTransactionID(),
		NodeAccountID:            transaction.pbBody.GetNodeAccountID(),
		TransactionFee:           transaction.pbBody.GetTransactionFee(),
		TransactionValidDuration: transaction.pbBody.GetTransactionValidDuration(),
		GenerateRecord:           transaction.pbBody.GetGenerateRecord(),
		Memo:                     transaction.pbBody.GetMemo(),
		Data: &proto.TransactionBody_CryptoDelete{
			CryptoDelete: &proto.CryptoDeleteTransactionBody{
				TransferAccountID: transaction.pb.GetTransferAccountID(),
				DeleteAccountID:   transaction.pb.GetDeleteAccountID(),
			},
		},
	}

	txBytes, err := protobuf.Marshal(body)
	if err != nil {
		return &ScheduleCreateTransaction{}, err
	}

	return NewScheduleCreateTransaction().setTransactionBodyBytes(txBytes), nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func accountDeleteTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getCrypto().CryptoDelete,
	}
}

func (transaction *AccountDeleteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *AccountDeleteTransaction) Sign(
	privateKey PrivateKey,
) *AccountDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *AccountDeleteTransaction) SignWithOperator(
	client *Client,
) (*AccountDeleteTransaction, error) {
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
func (transaction *AccountDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountDeleteTransaction {
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
func (transaction *AccountDeleteTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil || client.operator == nil {
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
		accountDeleteTransaction_getMethod,
		transaction_mapResponseStatus,
		transaction_mapResponse,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.transaction.NodeID,
			Hash:          make([]byte, 0),
		}, err
	}

	hash, err := transaction.GetTransactionHash()

	return TransactionResponse{
		TransactionID: transaction.GetTransactionID(),
		NodeID:        resp.transaction.NodeID,
		Hash:          hash,
	}, nil
}

func (transaction *AccountDeleteTransaction) onFreeze(
	pbBody *proto.TransactionBody,
) bool {
	pbBody.Data = &proto.TransactionBody_CryptoDelete{
		CryptoDelete: transaction.pb,
	}

	return true
}

func (transaction *AccountDeleteTransaction) Freeze() (*AccountDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *AccountDeleteTransaction) FreezeWith(client *Client) (*AccountDeleteTransaction, error) {
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

func (transaction *AccountDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetMaxTransactionFee(fee Hbar) *AccountDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *AccountDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetTransactionMemo(memo string) *AccountDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *AccountDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *AccountDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *AccountDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetTransactionID(transactionID TransactionID) *AccountDeleteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the node AccountID for this AccountDeleteTransaction.
func (transaction *AccountDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *AccountDeleteTransaction) SetMaxRetry(count int) *AccountDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *AccountDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountDeleteTransaction {
	if !transaction.IsFrozen() {
		transaction.Freeze()
	}

	transaction.Transaction.AddSignature(publicKey, signature)
	return transaction
}
