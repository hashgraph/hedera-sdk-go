package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type Fee interface {
	toProtobuf() *proto.CustomFee
	validateNetworkOnIDs(client *Client) error
}

type CustomFixedFee struct {
	CustomFee
	Amount              int64
	DenominationTokenID *TokenID
}

func customFixedFeeFromProtobuf(fixedFee *proto.FixedFee, customFee CustomFee) CustomFixedFee {
	id := tokenIDFromProtobuf(fixedFee.DenominatingTokenId)
	return CustomFixedFee{
		CustomFee:           customFee,
		Amount:              fixedFee.Amount,
		DenominationTokenID: &id,
	}
}

func (fee CustomFixedFee) validateNetworkOnIDs(client *Client) error {
	if client == nil {
		return nil
	}
	if fee.DenominationTokenID != nil {
		if err := fee.DenominationTokenID.Validate(client); err != nil {
			return err
		}
	}

	if fee.FeeCollectorAccountID != nil {
		if err := fee.FeeCollectorAccountID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (fee CustomFixedFee) toProtobuf() *proto.CustomFee {
	var tokenID *proto.TokenID
	if fee.DenominationTokenID != nil {
		tokenID = fee.DenominationTokenID.toProtobuf()
	}

	var FeeCollectorAccountID *proto.AccountID
	if fee.FeeCollectorAccountID != nil {
		FeeCollectorAccountID = fee.CustomFee.FeeCollectorAccountID.toProtobuf()
	}

	return &proto.CustomFee{
		Fee: &proto.CustomFee_FixedFee{
			FixedFee: &proto.FixedFee{
				Amount:              fee.Amount,
				DenominatingTokenId: tokenID,
			},
		},
		FeeCollectorAccountId: FeeCollectorAccountID,
	}
}

func (fee *CustomFixedFee) SetHbarAmount(hbar Hbar) {
	fee.Amount = int64(hbar.As(HbarUnits.Hbar))
	fee.DenominationTokenID = nil
}

func (fee *CustomFixedFee) GetHbarAmount() Hbar {
	return NewHbar(float64(fee.Amount))
}

func (fee *CustomFixedFee) SetDenominatingTokenToSameToken() {
	fee.DenominationTokenID = &TokenID{0, 0, 0, nil}
}

func (fee CustomFixedFee) ToBytes() []byte {
	data, err := protobuf.Marshal(fee.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}
