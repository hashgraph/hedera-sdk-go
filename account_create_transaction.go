//go:generate go build ./generator/main.go
//go:generate ./main
package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// AccountCreateTransaction creates a new account. After the account is created, the AccountID for it is in the receipt,
// or by asking for a Record of the transaction to be created, and retrieving that. The account can then automatically
// generate records for large transfers into it or out of it, which each last for 25 hours. Records are generated for
// any transfer that exceeds the thresholds given here. This account is charged hbar for each record generated, so the
// thresholds are useful for limiting Record generation to happen only for large transactions.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0,
// with a null key. Future versions of the API will support multiple realms and multiple shards.
//
// ProtobufType: CryptoCreateTransactionBody
// ProtobufAccessor: CryptoCreateAccount
type AccountCreateTransaction struct {
	Transaction
	pb                       *proto.CryptoCreateTransactionBody
	key                      Key             `hedera:"setter,getter,fromProtobuf,toProtobuf"`
	initialBalance           Hbar            `hedera:"setter,getter,fromProtobuf,toProtobuf"`
	accountMemo              string          `hedera:"setter,getter,fromProtobuf,toProtobuf"`
	sendRecordThreshold      Hbar            `hedera:"setter,getter,fromProtobuf,toProtobuf"`
	receiveRecordThreshold   Hbar            `hedera:"setter,getter,fromProtobuf,toProtobuf"`
	proxyAccountID           AccountID       `hedera:"setter,getter,fromProtobuf,toProtobuf"`
	autoRenewPeriod          time.Duration   `hedera:"setter,getter,fromProtobuf,toProtobuf"`
	receiverSigRequired      bool            `hedera:"setter,getter,fromProtobuf,toProtobuf"`
	transactionIDs           []TransactionID `hedera:"setter,getter,singular"`
	nodeAccountIDs           []AccountID     `hedera:"setter,getter"`
	maxRetry                 int             `hedera:"setter,getter"`
	transactionValidDuration time.Duration   `hedera:"setter,getter"`
	transactionMemo          string          `hedera:"setter,getter"`
	maxTransactionFee        Hbar            `hedera:"setter,getter"`
}

// NewAccountCreateTransaction creates an AccountCreateTransaction transaction which can be used to construct and
// execute a Crypto Create Transaction.
func NewAccountCreateTransaction() *AccountCreateTransaction {
	pb := &proto.CryptoCreateTransactionBody{}

	transaction := AccountCreateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	// TODO: Undo this change, this should use setters
	transaction.autoRenewPeriod = 7890000 * time.Second
	transaction.SetMaxTransactionFee(NewHbar(2))

	// Default to maximum values for record thresholds. Without this records would be
	// auto-created whenever a send or receive transaction takes place for this new account.
	// This should be an explicit ask.
	transaction.receiveRecordThreshold = MaxHbar
	transaction.sendRecordThreshold = MaxHbar

	return &transaction
}

func (transaction *AccountCreateTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *AccountCreateTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	pb := transaction.toProtobuf()
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.pbBody.GetTransactionFee(),
		Memo:           transaction.pbBody.GetMemo(),
		Data: &proto.SchedulableTransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: pb,
		},
	}, nil
}

func accountCreateTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getCrypto().CreateAccount,
	}
}

// Execute executes the Transaction with the provided client
func (transaction *AccountCreateTransaction) Execute(
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
		accountCreateTransaction_getMethod,
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

func (transaction *AccountCreateTransaction) FreezeWith(client *Client) (*AccountCreateTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	err := transaction.validateChecksums(client)
	if err != nil {
		return &AccountCreateTransaction{}, err
	}

	if !transaction.onFreeze(transaction.pbBody, transaction.toProtobuf()) {
		return transaction, nil
	}

	return transaction, transaction_freezeWith(&transaction.Transaction, client)
}

