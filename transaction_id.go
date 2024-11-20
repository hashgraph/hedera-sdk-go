package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// TransactionID is the id used to identify a Transaction on the Hiero _Network. It consists of an AccountID and a
// a valid start time.
type TransactionID struct {
	AccountID  *AccountID
	ValidStart *time.Time
	scheduled  bool
	Nonce      *int32
}

// NewTransactionID constructs a new Transaction id struct with the provided AccountID and the valid start time set
// to the current time - 10 seconds.
func TransactionIDGenerate(accountID AccountID) TransactionID {
	allowance := -(time.Duration(rand.Int63n(5*int64(time.Second))) + (8 * time.Second)) // nolint
	validStart := time.Now().UTC().Add(allowance)

	return TransactionID{&accountID, &validStart, false, nil}
}

// NewTransactionIDWithValidStart constructs a new Transaction id struct with the provided AccountID and the valid start
// time set to a provided time.
func NewTransactionIDWithValidStart(accountID AccountID, validStart time.Time) TransactionID {
	return TransactionID{&accountID, &validStart, false, nil}
}

// GetReceipt queries the _Network for a receipt corresponding to the TransactionID's transaction. If the status of the
// receipt is exceptional an ErrHederaReceiptStatus will be returned alongside the receipt, otherwise only the receipt
// will be returned.
func (id TransactionID) GetReceipt(client *Client) (TransactionReceipt, error) {
	return NewTransactionReceiptQuery().
		SetTransactionID(id).
		Execute(client)
}

// GetRecord queries the _Network for a record corresponding to the TransactionID's transaction. If the status of the
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

// String returns a string representation of the TransactionID in `AccountID@ValidStartSeconds.ValidStartNanos?scheduled_bool/nonce` format
func (id TransactionID) String() string {
	var pb *services.Timestamp
	var returnString string
	if id.AccountID != nil && id.ValidStart != nil {
		pb = _TimeToProtobuf(*id.ValidStart)
		// Use fmt.Sprintf to format the string with leading zeros
		returnString = id.AccountID.String() + "@" + fmt.Sprintf("%d.%09d", pb.Seconds, pb.Nanos)
	}

	if id.scheduled {
		returnString += "?scheduled"
	}

	if id.Nonce != nil {
		returnString += "/" + fmt.Sprint(*id.Nonce)
	}

	return returnString
}

// TransactionIDFromString constructs a TransactionID from a string representation
func TransactionIdFromString(data string) (TransactionID, error) { // nolint
	parts := strings.SplitN(data, "/", 2)

	var nonce *int32
	if len(parts) == 2 {
		temp, _ := strconv.ParseInt(parts[1], 10, 32)
		temp32 := int32(temp)
		nonce = &temp32
	}
	parts = strings.SplitN(parts[0], "?", 2)

	var accountID *AccountID
	var validStart *time.Time
	scheduled := len(parts) == 2 && strings.Compare(parts[1], "scheduled") == 0

	parts = strings.SplitN(parts[0], "@", 2)

	if len(parts) != 2 {
		return TransactionID{}, errors.New("expecting [{account}@{seconds}.{nanos}|{nonce}][?scheduled]")
	}

	temp, err := AccountIDFromString(parts[0])
	accountID = &temp
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
		AccountID:  accountID,
		ValidStart: validStart,
		scheduled:  scheduled,
		Nonce:      nonce,
	}, nil
}

func (id TransactionID) _ToProtobuf() *services.TransactionID {
	var validStart *services.Timestamp
	if id.ValidStart != nil {
		validStart = _TimeToProtobuf(*id.ValidStart)
	}

	var accountID *services.AccountID
	if id.AccountID != nil {
		accountID = id.AccountID._ToProtobuf()
	}

	var nonce int32
	if id.Nonce != nil {
		nonce = *id.Nonce
	}

	return &services.TransactionID{
		TransactionValidStart: validStart,
		AccountID:             accountID,
		Scheduled:             id.scheduled,
		Nonce:                 nonce,
	}
}

func _TransactionIDFromProtobuf(pb *services.TransactionID) TransactionID {
	if pb == nil {
		return TransactionID{}
	}
	var validStart time.Time
	if pb.TransactionValidStart != nil {
		validStart = _TimeFromProtobuf(pb.TransactionValidStart)
	}

	var accountID AccountID
	if pb.AccountID != nil {
		accountID = *_AccountIDFromProtobuf(pb.AccountID)
	}

	var nonce *int32
	if pb.Nonce != 0 {
		nonce = &pb.Nonce
	}

	return TransactionID{&accountID, &validStart, pb.Scheduled, nonce}
}

// ToBytes returns a byte array representation of the TransactionID
func (id TransactionID) ToBytes() []byte {
	data, err := protobuf.Marshal(id._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// TransactionIDFromBytes constructs a TransactionID from a byte array
func TransactionIDFromBytes(data []byte) (TransactionID, error) {
	if data == nil {
		return TransactionID{}, errByteArrayNull
	}
	pb := services.TransactionID{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TransactionID{}, err
	}

	return _TransactionIDFromProtobuf(&pb), nil
}

// SetScheduled sets the scheduled flag on the TransactionID
func (id TransactionID) SetScheduled(scheduled bool) TransactionID {
	id.scheduled = scheduled
	return id
}

// GetScheduled returns the scheduled flag on the TransactionID
func (id TransactionID) GetScheduled() bool {
	return id.scheduled
}

// SetNonce sets the nonce on the TransactionID
func (id TransactionID) SetNonce(nonce int32) TransactionID {
	id.Nonce = &nonce
	return id
}

// GetNonce returns the nonce on the TransactionID
func (id TransactionID) GetNonce() int32 {
	if id.Nonce != nil {
		return *id.Nonce
	}
	return 0
}
