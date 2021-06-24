package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

// TransactionID is the id used to identify a Transaction on the Hedera network. It consists of an AccountID and a
// a valid start time.
type TransactionID struct {
	AccountID  *AccountID
	ValidStart *time.Time
	scheduled  bool
}

// NewTransactionID constructs a new Transaction id struct with the provided AccountID and the valid start time set
// to the current time - 10 seconds.
func TransactionIDGenerate(accountID AccountID) TransactionID {
	allowance := -(time.Duration(rand.Int63n(5*int64(time.Second))) + (8 * time.Second))
	validStart := time.Now().UTC().Add(allowance)

	return TransactionID{&accountID, &validStart, false}
}

// NewTransactionIDWithValidStart constructs a new Transaction id struct with the provided AccountID and the valid start
// time set to a provided time.
func NewTransactionIDWithValidStart(accountID AccountID, validStart time.Time) TransactionID {
	return TransactionID{&accountID, &validStart, false}
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
	var returnString string
	if id.AccountID != nil && id.ValidStart != nil {
		pb = timeToProtobuf(*id.ValidStart)
		returnString = id.AccountID.String() + "@" + strconv.FormatInt(pb.Seconds, 10) + "." + fmt.Sprint(pb.Nanos)
	}

	if id.scheduled {
		returnString = returnString + "?scheduled"
	}

	return returnString
}

func TransactionIdFromString(data string) (TransactionID, error) {
	parts := strings.SplitN(data, "?", 2)

	var accountId *AccountID
	var validStart *time.Time
	scheduled := len(parts) == 2 && strings.Compare(parts[1], "scheduled") == 0

	parts = strings.SplitN(parts[0], "@", 2)

	if len(parts) != 2 {
		return TransactionID{}, errors.New("expecting [{account}@{seconds}.{nanos}|{nonce}][?scheduled]")
	}

	temp, err := AccountIDFromString(parts[0])
	accountId = &temp
	if err != nil {
		return TransactionID{}, err
	}

	validStartParts := strings.SplitN(parts[1], ".", 2)

	if len(validStartParts) != 2 {
		return TransactionID{}, errors.New("expecting {account}@{seconds}.{nanos}")
	}

	sec, err := strconv.ParseInt(validStartParts[0], 10, 64)
	if err != nil {
		return TransactionID{}, err
	}

	nano, err := strconv.ParseInt(validStartParts[1], 10, 64)
	if err != nil {
		return TransactionID{}, err
	}

	temp2 := time.Unix(sec, nano)
	validStart = &temp2

	return TransactionID{
		AccountID:  accountId,
		ValidStart: validStart,
		scheduled:  scheduled,
	}, nil
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
		Scheduled:             id.scheduled,
	}
}

func transactionIDFromProtobuf(pb *proto.TransactionID) TransactionID {
	if pb == nil {
		return TransactionID{}
	}
	var validStart time.Time
	if pb.TransactionValidStart != nil {
		validStart = timeFromProtobuf(pb.TransactionValidStart)
	}

	var accountID AccountID
	if pb.AccountID != nil {
		accountID = accountIDFromProtobuf(pb.AccountID)
	}

	return TransactionID{&accountID, &validStart, pb.Scheduled}
}

func (id TransactionID) ToBytes() []byte {
	data, err := protobuf.Marshal(id.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TransactionIDFromBytes(data []byte) (TransactionID, error) {
	if data == nil {
		return TransactionID{}, errByteArrayNull
	}
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

func TransactionIDValidateNetworkOnIDs(id TransactionID, other AccountID) error {
	if !id.AccountID.isZero() && !other.isZero() && id.AccountID.Network != nil && other.Network != nil && *id.AccountID.Network != *other.Network {
		return errNetworkMismatch
	}

	return nil
}
