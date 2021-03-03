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
	AccountID  *AccountID
	ValidStart *time.Time
	Nonce      []byte
	scheduled  bool
}

// NewTransactionID constructs a new Transaction id struct with the provided AccountID and the valid start time set
// to the current time - 10 seconds.
func TransactionIDGenerate(accountID AccountID) TransactionID {
	allowance := -(time.Duration(rand.Int63n(5*int64(time.Second))) + (8 * time.Second))
	validStart := time.Now().UTC().Add(allowance)

	return TransactionID{&accountID, &validStart, nil, false}
}

func TransactionIDWithNonce(nonce []byte) TransactionID {
	return TransactionID{nil, nil, nonce, false}
}

// NewTransactionIDWithValidStart constructs a new Transaction id struct with the provided AccountID and the valid start
// time set to a provided time.
func NewTransactionIDWithValidStart(accountID AccountID, validStart time.Time) TransactionID {
	return TransactionID{&accountID, &validStart, nil, false}
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
	var pb *proto.Timestamp
	if id.ValidStart != nil {
		pb = timeToProtobuf(*id.ValidStart)
	} else {
		return fmt.Sprintf("%v@%v.%v", id.AccountID, 0, 0)
	}

	return fmt.Sprintf("%v@%v.%v", id.AccountID, pb.Seconds, pb.Nanos)
}

func (id TransactionID) toProtobuf() *proto.TransactionID {
	var validStart *proto.Timestamp
	if id.ValidStart != nil {
		validStart = timeToProtobuf(*id.ValidStart)
	}

	var accountID *proto.AccountID
	if id.AccountID != nil {
		accountID = id.AccountID.toProtobuf()
	}

	return &proto.TransactionID{
		TransactionValidStart: validStart,
		AccountID:             accountID,
		Nonce:                 id.Nonce,
		Scheduled:             id.scheduled,
	}
}

func transactionIDFromProtobuf(pb *proto.TransactionID) TransactionID {
	var validStart time.Time
	if pb.TransactionValidStart != nil {
		validStart = timeFromProtobuf(pb.TransactionValidStart)
	}

	var accountID AccountID
	if pb.AccountID != nil {
		accountID = accountIDFromProtobuf(pb.AccountID)
	}

	return TransactionID{&accountID, &validStart, pb.Nonce, pb.Scheduled}
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

func (id TransactionID) SetScheduled(scheduled bool) TransactionID {
	id.scheduled = scheduled
	return id
}

func (id TransactionID) GetScheduled() bool {
	return id.scheduled
}
