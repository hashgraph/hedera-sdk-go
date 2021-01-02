package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"math/rand"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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
	allowance := -(time.Duration(rand.Int63n(5*int64(time.Second))) + (8 * time.Second))
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
	return NewTransactionReceiptQuery().
		SetTransactionID(id).
		Execute(client)
}

// GetRecord queries the network for a record corresponding to the TransactionID's transaction. If the status of the
// record's receipt is exceptional an ErrHederaRecordStatus will be returned alongside the record, otherwise, only the
// record will be returned. If consensus has not been reached, this function will return a HederaReceiptError with a
// status of StatusBusy.
func (id TransactionID) GetRecord(client *Client) (TransactionRecord, error) {
	_, err := NewTransactionReceiptQuery().
		SetTransactionID(id).
		Execute(client)

	if err != nil {
		return TransactionRecord{}, err
	}

	return NewTransactionRecordQuery().
		SetTransactionID(id).
		Execute(client)
}

// String returns a string representation of the TransactionID in `AccountID@ValidStartSeconds.ValidStartNanos` format
func (id TransactionID) String() string {
	pb := timeToProtobuf(id.ValidStart)
	return fmt.Sprintf("%v@%v.%v", id.AccountID, pb.Seconds, pb.Nanos)
}

func (id TransactionID) toProtobuf() *proto.TransactionID {
	return &proto.TransactionID{
		TransactionValidStart: timeToProtobuf(id.ValidStart),
		AccountID:             id.AccountID.toProtobuf(),
	}
}

func transactionIDFromProtobuf(pb *proto.TransactionID) TransactionID {
	validStart := timeFromProtobuf(pb.TransactionValidStart)
	accountID := accountIDFromProtobuf(pb.AccountID)

	return TransactionID{accountID, validStart}
}

func (id TransactionID) ToBytes() []byte {
	data, err := protobuf.Marshal(id.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TransactionIDFromBytes(data []byte) (TransactionID, error) {
	pb := proto.TransactionID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TransactionID{}, err
	}

	return transactionIDFromProtobuf(&pb), nil
}
