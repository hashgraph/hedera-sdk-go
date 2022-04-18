package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type HbarAllowance struct {
	OwnerAccountID   *AccountID
	SpenderAccountID *AccountID
	Amount           int64
}

func NewHbarAllowance(ownerAccountID AccountID, spenderAccountID AccountID, amount int64) HbarAllowance { //nolint
	return HbarAllowance{
		OwnerAccountID:   &ownerAccountID,
		SpenderAccountID: &spenderAccountID,
		Amount:           amount,
	}
}

func _HbarAllowanceFromProtobuf(pb *services.CryptoAllowance) HbarAllowance {
	body := HbarAllowance{
		Amount: pb.Amount,
	}

	if pb.Spender != nil {
		body.SpenderAccountID = _AccountIDFromProtobuf(pb.Spender)
	}

	if pb.Owner != nil {
		body.OwnerAccountID = _AccountIDFromProtobuf(pb.Owner)
	}

	return body
}

func (approval *HbarAllowance) _ToProtobuf() *services.CryptoAllowance {
	body := &services.CryptoAllowance{
		Amount: approval.Amount,
	}

	if approval.SpenderAccountID != nil {
		body.Spender = approval.SpenderAccountID._ToProtobuf()
	}

	if approval.OwnerAccountID != nil {
		body.Owner = approval.OwnerAccountID._ToProtobuf()
	}

	return body
}

func (approval *HbarAllowance) String() string {
	if approval.OwnerAccountID != nil && approval.SpenderAccountID != nil { //nolint
		return fmt.Sprintf("OwnerAccountID: %s, SpenderAccountID: %s, Amount: %s", approval.OwnerAccountID.String(), approval.SpenderAccountID.String(), HbarFromTinybar(approval.Amount).String())
	} else if approval.OwnerAccountID != nil {
		return fmt.Sprintf("OwnerAccountID: %s, Amount: %s", approval.OwnerAccountID.String(), HbarFromTinybar(approval.Amount).String())
	} else if approval.SpenderAccountID != nil {
		return fmt.Sprintf("SpenderAccountID: %s, Amount: %s", approval.SpenderAccountID.String(), HbarFromTinybar(approval.Amount).String())
	}

	return ""
}
