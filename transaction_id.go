package hedera

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// TransactionID is the id used to identify a Transaction on the Hedera network. It consists of an AccountID and a
// a valid start time.
type TransactionID struct {
	AccountID  AccountID
	ValidStart time.Time
}

// NewTransactionID constructs a new Transaction id struct with the provided AccountID and the valid start time set
// to the current time - 10 seconds.
func TransactionIDGenerate(accountID AccountID) TransactionID {
	allowance := -(time.Duration(rand.Intn(5*int(time.Second))) + (8 * time.Second))
	validStart := time.Now().UTC().Add(allowance)

	return TransactionID{accountID, validStart}
}

// NewTransactionIDWithValidStart constructs a new Transaction id struct with the provided AccountID and the valid start
// time set to a provided time.
func NewTransactionIDWithValidStart(accountID AccountID, validStart time.Time) TransactionID {
	return TransactionID{accountID, validStart}
}

// GetReceipt queries the network for a receipt corresponding to the TransactionID's transaction. If the status of the
// receipt is exceptional an ErrHederaReceiptStatus will be returned alongside the receipt, otherwise only the receipt
// will be returned.
func (id TransactionID) GetReceipt(client *Client) (TransactionReceipt, error) {
	receipt, err := NewTransactionReceiptQuery().
		SetTransactionID(id).
		Execute(client)

	if err != nil {
		// something went wrong with the query
		return TransactionReceipt{}, err
	}

	return receipt, nil
}

// // GetRecord queries the network for a record corresponding to the TransactionID's transaction. If the status of the
// // record's receipt is exceptional an ErrHederaRecordStatus will be returned alongside the record, otherwise, only the
// // record will be returned. If consensus has not been reached, this function will return a HederaReceiptError with a
// // status of StatusBusy.
// func (id TransactionID) GetRecord(client *Client) (TransactionRecord, error) {
// 	receipt, err := NewTransactionReceiptQuery().
// 		SetTransactionID(id).
// 		Execute(client)

// 	if err != nil {
// 		// something went wrong with the receipt query
// 		return TransactionRecord{}, err
// 	}

// 	if receipt.Status == StatusBusy {
// 		// Consensus has not been reached.
// 		return TransactionRecord{}, newErrHederaReceiptStatus(id, receipt.Status)
// 	}

// 	record, err := NewTransactionRecordQuery().SetTransactionID(id).Execute(client)

// 	if err != nil {
// 		// something went wrong with the record query
// 		return TransactionRecord{}, err
// 	}

// 	if record.Receipt.Status.isExceptional(true) {
// 		// The receipt status of the record's receipt was exceptional, return the record AND the error
// 		return record, newErrHederaRecordStatus(id, record.Receipt.Status)
// 	}

// 	return record, nil
// }

// String returns a string representation of the TransactionID in `AccountID@ValidStartSeconds.ValidStartNanos` format
func (id TransactionID) String() string {
	pb := timeToProto(id.ValidStart)
	return fmt.Sprintf("%v@%v.%v", id.AccountID, pb.Seconds, pb.Nanos)
}

func (id TransactionID) toProtobuf() *proto.TransactionID {
	return &proto.TransactionID{
		TransactionValidStart: timeToProto(id.ValidStart),
		AccountID:             id.AccountID.toProtobuf(),
	}
}

func transactionIDFromProto(pb *proto.TransactionID) TransactionID {
	validStart := timeFromProto(pb.TransactionValidStart)
	accountID := accountIDFromProto(pb.AccountID)

	return TransactionID{accountID, validStart}
}
