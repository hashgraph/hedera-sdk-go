package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

type AssessedCustomFee struct {
	Amount                int64
	TokenID               *TokenID
	FeeCollectorAccountId *AccountID // nolint
	PayerAccountIDs       []*AccountID
}

func _AssessedCustomFeeFromProtobuf(assessedFee *proto.AssessedCustomFee) AssessedCustomFee {
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

func (fee *AssessedCustomFee) _ToProtobuf() *proto.AssessedCustomFee {
	var tokenID *proto.TokenID
	if fee.TokenID != nil {
		tokenID = fee.TokenID._ToProtobuf()
	}

	var accountID *proto.AccountID
	if fee.TokenID != nil {
		accountID = fee.FeeCollectorAccountId._ToProtobuf()
	}

	payerAccountIds := make([]*proto.AccountID, len(fee.PayerAccountIDs))

	for _, id := range fee.PayerAccountIDs {
		payerAccountIds = append(payerAccountIds, id._ToProtobuf())
	}

	return &proto.AssessedCustomFee{
		Amount:                  fee.Amount,
		TokenId:                 tokenID,
		FeeCollectorAccountId:   accountID,
		EffectivePayerAccountId: payerAccountIds,
	}
}

func (fee *AssessedCustomFee) ToBytes() []byte {
	data, err := protobuf.Marshal(fee._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func AssessedCustomFeeFromBytes(data []byte) (AssessedCustomFee, error) {
	if data == nil {
		return AssessedCustomFee{}, errByteArrayNull
	}
	pb := proto.AssessedCustomFee{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return AssessedCustomFee{}, err
	}

	return _AssessedCustomFeeFromProtobuf(&pb), nil
}

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
