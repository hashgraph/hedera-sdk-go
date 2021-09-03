package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
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

func NewCustomFixedFee() *CustomFixedFee {
	return &CustomFixedFee{
		CustomFee:           CustomFee{},
		Amount:              0,
		DenominationTokenID: nil,
	}
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
		if fee.DenominationTokenID != nil {
			if err := fee.DenominationTokenID.Validate(client); err != nil {
				return err
			}
		}

	}

	if fee.FeeCollectorAccountID != nil {
		if fee.FeeCollectorAccountID != nil {
			if err := fee.FeeCollectorAccountID.Validate(client); err != nil {
				return err
			}
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

func (fee *CustomFixedFee) SetAmount(tinybar int64) *CustomFixedFee {
	fee.Amount = tinybar
	return fee
}

func (fee *CustomFixedFee) GetAmount() Hbar {
	return NewHbar(float64(fee.Amount))
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

func (fee *CustomFixedFee) SetDenominatingTokenID(id TokenID) *CustomFixedFee {
	fee.DenominationTokenID = &id
	return fee
}

func (fee *CustomFixedFee) GetDenominatingTokenID() TokenID {
	if fee.DenominationTokenID != nil {
		return *fee.DenominationTokenID
	}

	return TokenID{}
}

func (fee *CustomFixedFee) SetFeeCollectorAccountID(id AccountID) *CustomFixedFee {
	fee.FeeCollectorAccountID = &id
	return fee
}

func (fee *CustomFixedFee) GetFeeCollectorAccountID() AccountID {
	return *fee.FeeCollectorAccountID
}

func (fee CustomFixedFee) ToBytes() []byte {
	data, err := protobuf.Marshal(fee.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}
