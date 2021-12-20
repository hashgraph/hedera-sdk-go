package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type TokenTransfer struct {
	AccountID AccountID
	Amount    int64
}

type _TokenTransfers struct {
	transfers []TokenTransfer
}

func NewTokenTransfer(accountID AccountID, amount int64) TokenTransfer {
	return TokenTransfer{
		AccountID: accountID,
		Amount:    amount,
	}
}

func _TokenTransferFromProtobuf(pb *services.AccountAmount) TokenTransfer {
	if pb == nil {
		return TokenTransfer{}
	}

	accountID := AccountID{}
	if pb.AccountID != nil {
		accountID = *_AccountIDFromProtobuf(pb.AccountID)
	}

	return TokenTransfer{
		AccountID: accountID,
		Amount:    pb.Amount,
	}
}

func (transfer *TokenTransfer) _ToProtobuf() *services.AccountAmount {
	return &services.AccountAmount{
		AccountID: transfer.AccountID._ToProtobuf(),
		Amount:    transfer.Amount,
	}
}

func (transfer TokenTransfer) ToBytes() []byte {
	data, err := protobuf.Marshal(transfer._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func TokenTransferFromBytes(data []byte) (TokenTransfer, error) {
	if data == nil {
		return TokenTransfer{}, errByteArrayNull
	}
	pb := services.AccountAmount{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return TokenTransfer{}, err
	}

	return _TokenTransferFromProtobuf(&pb), nil
}

func (transfer TokenTransfer) String() string {
	return fmt.Sprintf("accountID: %s, amount: %d", transfer.AccountID.String(), transfer.Amount)
}

func (transfers _TokenTransfers) Len() int {
	return len(transfers.transfers)
}
func (transfers _TokenTransfers) Swap(i, j int) {
	transfers.transfers[i], transfers.transfers[j] = transfers.transfers[j], transfers.transfers[i]
}

func (transfers _TokenTransfers) Less(i, j int) bool {
	if transfers.transfers[i].AccountID.Shard < transfers.transfers[j].AccountID.Shard { //nolint
		return true
	} else if transfers.transfers[i].AccountID.Shard > transfers.transfers[j].AccountID.Shard {
		return false
	}

	if transfers.transfers[i].AccountID.Realm < transfers.transfers[j].AccountID.Realm { //nolint
		return true
	} else if transfers.transfers[i].AccountID.Realm > transfers.transfers[j].AccountID.Realm {
		return false
	}

	if transfers.transfers[i].AccountID.AliasKey != nil && transfers.transfers[j].AccountID.AliasKey != nil {
		if transfers.transfers[i].AccountID.String() < transfers.transfers[j].AccountID.String() { //nolint
			return true
		} else if transfers.transfers[i].AccountID.String() > transfers.transfers[j].AccountID.String() {
			return false
		}
	}

	if transfers.transfers[i].AccountID.Account < transfers.transfers[j].AccountID.Account { //nolint
		return true
	} else { //nolint
		return false
	}
}
