package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"fmt"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	protobuf "google.golang.org/protobuf/proto"
)

type AssessedCustomFee struct {
	Amount                int64
	TokenID               *TokenID
	FeeCollectorAccountId *AccountID // nolint
	PayerAccountIDs       []*AccountID
}

func _AssessedCustomFeeFromProtobuf(assessedFee *services.AssessedCustomFee) AssessedCustomFee {
	accountID := _AccountIDFromProtobuf(assessedFee.FeeCollectorAccountId)
	tokenID := _TokenIDFromProtobuf(assessedFee.TokenId)

	payerAccountIds := make([]*AccountID, 0)

	for _, id := range assessedFee.EffectivePayerAccountId {
		payerAccountIds = append(payerAccountIds, _AccountIDFromProtobuf(id))
	}

	return AssessedCustomFee{
		Amount:                assessedFee.Amount,
		TokenID:               tokenID,
		FeeCollectorAccountId: accountID,
		PayerAccountIDs:       payerAccountIds,
	}
}

func (fee *AssessedCustomFee) _ToProtobuf() *services.AssessedCustomFee {
	var tokenID *services.TokenID
	if fee.TokenID != nil {
		tokenID = fee.TokenID._ToProtobuf()
	}

	var accountID *services.AccountID
	if fee.FeeCollectorAccountId != nil {
		accountID = fee.FeeCollectorAccountId._ToProtobuf()
	}

	payerAccountIds := make([]*services.AccountID, len(fee.PayerAccountIDs))

	for _, id := range fee.PayerAccountIDs {
		payerAccountIds = append(payerAccountIds, id._ToProtobuf())
	}

	return &services.AssessedCustomFee{
		Amount:                  fee.Amount,
		TokenId:                 tokenID,
		FeeCollectorAccountId:   accountID,
		EffectivePayerAccountId: payerAccountIds,
	}
}

// ToBytes returns the serialized bytes of a AssessedCustomFee
func (fee *AssessedCustomFee) ToBytes() []byte {
	data, err := protobuf.Marshal(fee._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

// AssessedCustomFeeFromBytes returns a AssessedCustomFee from bytes
func AssessedCustomFeeFromBytes(data []byte) (AssessedCustomFee, error) {
	if data == nil {
		return AssessedCustomFee{}, errByteArrayNull
	}
	pb := services.AssessedCustomFee{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return AssessedCustomFee{}, err
	}

	return _AssessedCustomFeeFromProtobuf(&pb), nil
}

// String returns a string representation of a AssessedCustomFee
func (fee AssessedCustomFee) String() string {
	accountIDs := ""
	for _, s := range fee.PayerAccountIDs {
		accountIDs = accountIDs + " " + s.String()
	}
	if fee.TokenID != nil {
		return fmt.Sprintf("feeCollectorAccountID: %s, amount: %d, tokenID: %s, payerAccountIds: %s", fee.FeeCollectorAccountId.String(), fee.Amount, fee.TokenID.String(), accountIDs)
	}

	return fmt.Sprintf("feeCollectorAccountID: %s, amount: %d, payerAccountIds: %s", fee.FeeCollectorAccountId.String(), fee.Amount, accountIDs)
}
