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

func customFixedFeeFromProtobuf(fixedFee *proto.FixedFee, customFee CustomFee, networkName *NetworkName) CustomFixedFee {
	id := tokenIDFromProtobuf(fixedFee.DenominatingTokenId, networkName)
	return CustomFixedFee{
		CustomFee:           customFee,
		Amount:              fixedFee.Amount,
		DenominationTokenID: &id,
	}
}

func (fee CustomFixedFee) validateNetworkOnIDs(client *Client) error {
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

func (fee CustomFixedFee) ToBytes() []byte {
	data, err := protobuf.Marshal(fee.toProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}
