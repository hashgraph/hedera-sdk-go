package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type SystemDeleteTransaction struct {
	Transaction
	contractID     ContractID
	fileID         FileID
	expirationTime *time.Time
}

func NewSystemDeleteTransaction() *SystemDeleteTransaction {
	transaction := SystemDeleteTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func systemDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) SystemDeleteTransaction {
	expiration := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		time.Now().Hour(), time.Now().Minute(),
		int(pb.GetSystemDelete().ExpirationTime.Seconds), time.Now().Nanosecond(), time.Now().Location(),
	)
	return SystemDeleteTransaction{
		Transaction:    transaction,
		contractID:     contractIDFromProtobuf(pb.GetSystemDelete().GetContractID()),
		fileID:         fileIDFromProtobuf(pb.GetSystemDelete().GetFileID()),
		expirationTime: &expiration,
	}
}

func (transaction *SystemDeleteTransaction) SetExpirationTime(expiration time.Time) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.expirationTime = &expiration
	return transaction
}

func (transaction *SystemDeleteTransaction) GetExpirationTime() int64 {
	if transaction.expirationTime != nil {
		return transaction.expirationTime.Unix()
	}

	return 0
}

func (transaction *SystemDeleteTransaction) SetContractID(id ContractID) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.contractID = id
	return transaction
}

func (transaction *SystemDeleteTransaction) GetContract() ContractID {
	return transaction.contractID
}

func (transaction *SystemDeleteTransaction) SetFileID(id FileID) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.fileID = id
	return transaction
}

func (transaction *SystemDeleteTransaction) GetFileID() FileID {
	return transaction.fileID
}

func (transaction *SystemDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	var err error
	err = transaction.contractID.Validate(client)
	if err != nil {
		return err
	}
	err = transaction.fileID.Validate(client)
	if err != nil {
		return err
	}

	return nil
}

func (transaction *SystemDeleteTransaction) build() *proto.TransactionBody {
	body := &proto.SystemDeleteTransactionBody{}

	if transaction.expirationTime != nil {
		body.ExpirationTime = &proto.TimestampSeconds{
			Seconds: transaction.expirationTime.Unix(),
		}
	}

	if !transaction.contractID.isZero() {
		body.Id = &proto.SystemDeleteTransactionBody_ContractID{
			ContractID: transaction.contractID.toProtobuf(),
		}
	}

	if !transaction.fileID.isZero() {
		body.Id = &proto.SystemDeleteTransactionBody_FileID{
			FileID: transaction.fileID.toProtobuf(),
		}
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_SystemDelete{
			SystemDelete: body,
		},
	}
}

func (transaction *SystemDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *SystemDeleteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.SystemDeleteTransactionBody{}

	if transaction.expirationTime != nil {
		body.ExpirationTime = &proto.TimestampSeconds{
			Seconds: transaction.expirationTime.Unix(),
		}
	}

	if !transaction.contractID.isZero() {
		body.Id = &proto.SystemDeleteTransactionBody_ContractID{
			ContractID: transaction.contractID.toProtobuf(),
		}
	}

	if !transaction.fileID.isZero() {
		body.Id = &proto.SystemDeleteTransactionBody_FileID{
			FileID: transaction.fileID.toProtobuf(),
		}
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_SystemDelete{
			SystemDelete: body,
		},
	}, nil

}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func systemDeleteTransaction_getMethod(request request, channel *channel) method {
	//switch os := runtime.GOOS; os {
	//case "darwin":
	//	fmt.Println("OS X.")
	//}
	if channel.getContract() == nil {
		return method{
			transaction: channel.getFile().SystemDelete,
		}
	} else {
		return method{
			transaction: channel.getContract().SystemDelete,
		}
	}
}

func (transaction *SystemDeleteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *SystemDeleteTransaction) Sign(
	privateKey PrivateKey,
) *SystemDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *SystemDeleteTransaction) SignWithOperator(
	client *Client,
) (*SystemDeleteTransaction, error) {
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
func (transaction *SystemDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *SystemDeleteTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *SystemDeleteTransaction) Execute(
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
		transaction_makeRequest(request{
			transaction: &transaction.Transaction,
		}),
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		systemDeleteTransaction_getMethod,
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

func (transaction *SystemDeleteTransaction) Freeze() (*SystemDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *SystemDeleteTransaction) FreezeWith(client *Client) (*SystemDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &SystemDeleteTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
}

func (transaction *SystemDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetMaxTransactionFee(fee Hbar) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *SystemDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetTransactionMemo(memo string) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *SystemDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *SystemDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetTransactionID(transactionID TransactionID) *SystemDeleteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this SystemDeleteTransaction.
func (transaction *SystemDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *SystemDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *SystemDeleteTransaction) SetMaxRetry(count int) *SystemDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *SystemDeleteTransaction) AddSignature(publicKey PublicKey, signature []byte) *SystemDeleteTransaction {
	transaction.requireOneNodeAccountID()

	if !transaction.isFrozen() {
		transaction.Freeze()
	}

	if transaction.keyAlreadySigned(publicKey) {
		return transaction
	}

	if len(transaction.signedTransactions) == 0 {
		return transaction
	}

	transaction.transactions = make([]*proto.Transaction, 0)
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)

	for index := 0; index < len(transaction.signedTransactions); index++ {
		transaction.signedTransactions[index].SigMap.SigPair = append(
			transaction.signedTransactions[index].SigMap.SigPair,
			publicKey.toSignaturePairProtobuf(signature),
		)
	}

	//transaction.signedTransactions[0].SigMap.SigPair = append(transaction.signedTransactions[0].SigMap.SigPair, publicKey.toSignaturePairProtobuf(signature))
	return transaction
}
