package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

type FreezeTransaction struct {
	Transaction
	startTime time.Time
	endTime   time.Time
	fileID    *FileID
	fileHash  []byte
}

func NewFreezeTransaction() *FreezeTransaction {
	transaction := FreezeTransaction{
		Transaction: newTransaction(),
	}

	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func freezeTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) FreezeTransaction {
	fileId := fileIDFromProtobuf(pb.GetFreeze().GetUpdateFile())

	startTime := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		int(pb.GetFreeze().GetStartHour()), int(pb.GetFreeze().GetStartMin()),
		0, time.Now().Nanosecond(), time.Now().Location(),
	)

	endTime := time.Date(
		time.Now().Year(), time.Now().Month(), time.Now().Day(),
		int(pb.GetFreeze().GetEndHour()), int(pb.GetFreeze().GetEndMin()),
		0, time.Now().Nanosecond(), time.Now().Location(),
	)

	return FreezeTransaction{
		Transaction: transaction,
		startTime:   startTime,
		endTime:     endTime,
		fileID:      &fileId,
		fileHash:    pb.GetFreeze().FileHash,
	}
}

func (transaction *FreezeTransaction) SetStartTime(startTime time.Time) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.startTime = startTime
	return transaction
}

func (transaction *FreezeTransaction) GetStartTime() time.Time {
	return transaction.startTime
}

func (transaction *FreezeTransaction) SetEndTime(endTime time.Time) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.endTime = endTime
	return transaction
}

func (transaction *FreezeTransaction) GetEndTime() time.Time {
	return transaction.endTime
}

func (transaction *FreezeTransaction) SetFileID(id FileID) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.fileID = &id
	return transaction
}

func (transaction *FreezeTransaction) GetFileID() *FileID {
	return transaction.fileID
}

func (transaction *FreezeTransaction) SetFileHash(hash []byte) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.fileHash = hash
	return transaction
}

func (transaction *FreezeTransaction) GetFileHash() []byte {
	return transaction.fileHash
}

func (transaction *FreezeTransaction) build() *proto.TransactionBody {
	body := &proto.FreezeTransactionBody{
		StartHour: int32(transaction.startTime.Hour()),
		StartMin:  int32(transaction.startTime.Minute()),
		EndHour:   int32(transaction.endTime.Hour()),
		EndMin:    int32(transaction.endTime.Minute()),
		FileHash:  transaction.fileHash,
	}

	if !transaction.fileID.isZero() {
		body.UpdateFile = transaction.fileID.toProtobuf()
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_Freeze{
			Freeze: body,
		},
	}
}

func (transaction *FreezeTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *FreezeTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.FreezeTransactionBody{
		StartHour: int32(transaction.startTime.Hour()),
		StartMin:  int32(transaction.startTime.Minute()),
		EndHour:   int32(transaction.endTime.Hour()),
		EndMin:    int32(transaction.endTime.Minute()),
		FileHash:  transaction.fileHash,
	}

	if !transaction.fileID.isZero() {
		body.UpdateFile = transaction.fileID.toProtobuf()
	}
	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_Freeze{
			Freeze: body,
		},
	}, nil
}

//
// The following methods must be copy-pasted/overriden at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

func freezeTransaction_getMethod(request request, channel *channel) method {
	return method{
		transaction: channel.getFreeze().Freeze,
	}
}

func (transaction *FreezeTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *FreezeTransaction) Sign(
	privateKey PrivateKey,
) *FreezeTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *FreezeTransaction) SignWithOperator(
	client *Client,
) (*FreezeTransaction, error) {
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
func (transaction *FreezeTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FreezeTransaction {
	if !transaction.IsFrozen() {
		_, _ = transaction.Freeze()
	}

	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *FreezeTransaction) Execute(
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
		transaction_makeRequest(request{
			transaction: &transaction.Transaction,
		}),
		transaction_advanceRequest,
		transaction_getNodeAccountID,
		freezeTransaction_getMethod,
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

func (transaction *FreezeTransaction) Freeze() (*FreezeTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *FreezeTransaction) FreezeWith(client *Client) (*FreezeTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, transaction_freezeWith(&transaction.Transaction, client, body)
}

func (transaction *FreezeTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this FreezeTransaction.
func (transaction *FreezeTransaction) SetMaxTransactionFee(fee Hbar) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *FreezeTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this FreezeTransaction.
func (transaction *FreezeTransaction) SetTransactionMemo(memo string) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *FreezeTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this FreezeTransaction.
func (transaction *FreezeTransaction) SetTransactionValidDuration(duration time.Duration) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *FreezeTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this FreezeTransaction.
func (transaction *FreezeTransaction) SetTransactionID(transactionID TransactionID) *FreezeTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the node AccountID for this FreezeTransaction.
func (transaction *FreezeTransaction) SetNodeAccountIDs(nodeID []AccountID) *FreezeTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *FreezeTransaction) SetMaxRetry(count int) *FreezeTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *FreezeTransaction) AddSignature(publicKey PublicKey, signature []byte) *FreezeTransaction {
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
