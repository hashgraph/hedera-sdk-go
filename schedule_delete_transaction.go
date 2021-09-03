package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"

	"time"
)

type ScheduleDeleteTransaction struct {
	Transaction
	scheduleID *ScheduleID
}

func NewScheduleDeleteTransaction() *ScheduleDeleteTransaction {
	transaction := ScheduleDeleteTransaction{
		Transaction: newTransaction(),
	}
	transaction.SetMaxTransactionFee(NewHbar(5))

	return &transaction
}

func scheduleDeleteTransactionFromProtobuf(transaction Transaction, pb *proto.TransactionBody) ScheduleDeleteTransaction {
	return ScheduleDeleteTransaction{
		Transaction: transaction,
		scheduleID:  scheduleIDFromProtobuf(pb.GetScheduleDelete().GetScheduleID()),
	}
}

func (transaction *ScheduleDeleteTransaction) SetScheduleID(scheduleID ScheduleID) *ScheduleDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.scheduleID = &scheduleID
	return transaction
}

func (transaction *ScheduleDeleteTransaction) GetScheduleID() ScheduleID {
	if transaction.scheduleID == nil {
		return ScheduleID{}
	}

	return *transaction.scheduleID
}

func (transaction *ScheduleDeleteTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if transaction.scheduleID != nil {
		if err := transaction.scheduleID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (transaction *ScheduleDeleteTransaction) build() *proto.TransactionBody {
	body := &proto.ScheduleDeleteTransactionBody{}
	if transaction.scheduleID != nil {
		body.ScheduleID = transaction.scheduleID.toProtobuf()
	}

	return &proto.TransactionBody{
		TransactionFee:           transaction.transactionFee,
		Memo:                     transaction.Transaction.memo,
		TransactionValidDuration: durationToProtobuf(transaction.GetTransactionValidDuration()),
		TransactionID:            transaction.transactionID.toProtobuf(),
		Data: &proto.TransactionBody_ScheduleDelete{
			ScheduleDelete: body,
		},
	}
}

func (transaction *ScheduleDeleteTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction.requireNotFrozen()

	scheduled, err := transaction.constructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction().setSchedulableTransactionBody(scheduled), nil
}

func (transaction *ScheduleDeleteTransaction) constructScheduleProtobuf() (*proto.SchedulableTransactionBody, error) {
	body := &proto.ScheduleDeleteTransactionBody{}
	if !transaction.scheduleID.isZero() {
		body.ScheduleID = transaction.scheduleID.toProtobuf()
	}

	return &proto.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &proto.SchedulableTransactionBody_ScheduleDelete{
			ScheduleDelete: body,
		},
	}, nil
}

func _ScheduleDeleteTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel.getSchedule().DeleteSchedule,
	}
}

func (transaction *ScheduleDeleteTransaction) IsFrozen() bool {
	return transaction.isFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *ScheduleDeleteTransaction) Sign(
	privateKey PrivateKey,
) *ScheduleDeleteTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *ScheduleDeleteTransaction) SignWithOperator(
	client *Client,
) (*ScheduleDeleteTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

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
func (transaction *ScheduleDeleteTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *ScheduleDeleteTransaction {
	if !transaction.keyAlreadySigned(publicKey) {
		transaction.signWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *ScheduleDeleteTransaction) Execute(
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
		_Request{
			transaction: &transaction.Transaction,
		},
		_TransactionShouldRetry,
		_TransactionMakeRequest(_Request{
			transaction: &transaction.Transaction,
		}),
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_ScheduleDeleteTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.transaction.NodeID,
		}, err
	}

	hash, err := transaction.GetTransactionHash()
	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{
		TransactionID: transaction.GetTransactionID(),
		NodeID:        resp.transaction.NodeID,
		Hash:          hash,
	}, nil
}

func (transaction *ScheduleDeleteTransaction) Freeze() (*ScheduleDeleteTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *ScheduleDeleteTransaction) FreezeWith(client *Client) (*ScheduleDeleteTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction.initFee(client)
	err := transaction.validateNetworkOnIDs(client)
	if err != nil {
		return &ScheduleDeleteTransaction{}, err
	}
	if err := transaction.initTransactionID(client); err != nil {
		return transaction, err
	}
	body := transaction.build()

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *ScheduleDeleteTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetMaxTransactionFee(fee Hbar) *ScheduleDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

func (transaction *ScheduleDeleteTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetTransactionMemo(memo string) *ScheduleDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *ScheduleDeleteTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetTransactionValidDuration(duration time.Duration) *ScheduleDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *ScheduleDeleteTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetTransactionID(transactionID TransactionID) *ScheduleDeleteTransaction {
	transaction.requireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountID sets the _Node AccountID for this ScheduleDeleteTransaction.
func (transaction *ScheduleDeleteTransaction) SetNodeAccountIDs(nodeID []AccountID) *ScheduleDeleteTransaction {
	transaction.requireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *ScheduleDeleteTransaction) SetMaxRetry(count int) *ScheduleDeleteTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *ScheduleDeleteTransaction) SetMaxBackoff(max time.Duration) *ScheduleDeleteTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *ScheduleDeleteTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *ScheduleDeleteTransaction) SetMinBackoff(min time.Duration) *ScheduleDeleteTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *ScheduleDeleteTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}
