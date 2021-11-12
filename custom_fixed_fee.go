package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type Fee interface {
	_ToProtobuf() *services.CustomFee
	_ValidateNetworkOnIDs(client *Client) error
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

func _CustomFixedFeeFromProtobuf(fixedFee *services.FixedFee, customFee CustomFee) CustomFixedFee {
	return CustomFixedFee{
		CustomFee:           customFee,
		Amount:              fixedFee.Amount,
		DenominationTokenID: _TokenIDFromProtobuf(fixedFee.DenominatingTokenId),
	}
}

func (fee CustomFixedFee) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}
	if fee.DenominationTokenID != nil {
		if fee.DenominationTokenID != nil {
			if err := fee.DenominationTokenID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}

	if fee.FeeCollectorAccountID != nil {
		if fee.FeeCollectorAccountID != nil {
			if err := fee.FeeCollectorAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}

	return nil
}

func (fee CustomFixedFee) _ToProtobuf() *services.CustomFee {
	var tokenID *services.TokenID
	if fee.DenominationTokenID != nil {
		tokenID = fee.DenominationTokenID._ToProtobuf()
	}

	var FeeCollectorAccountID *services.AccountID
	if fee.FeeCollectorAccountID != nil {
		FeeCollectorAccountID = fee.CustomFee.FeeCollectorAccountID._ToProtobuf()
	}

	return &services.CustomFee{
		Fee: &services.CustomFee_FixedFee{
			FixedFee: &services.FixedFee{
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

func (fee *CustomFixedFee) SetHbarAmount(hbar Hbar) *CustomFixedFee {
	fee.Amount = int64(hbar.As(HbarUnits.Tinybar))
	fee.DenominationTokenID = nil
	return fee
}

func (fee *CustomFixedFee) GetHbarAmount() Hbar {
	return NewHbar(float64(fee.Amount))
}

func (fee *CustomFixedFee) SetDenominatingTokenToSameToken() *CustomFixedFee {
	fee.DenominationTokenID = &TokenID{0, 0, 0, nil}
	return fee
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
	data, err := protobuf.Marshal(fee._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func (fee CustomFixedFee) String() string {
	if fee.DenominationTokenID != nil {
		return fmt.Sprintf("feeCollectorAccountID: %s, amount: %d, denominatingTokenID: %s", fee.FeeCollectorAccountID.String(), fee.Amount, fee.DenominationTokenID.String())
	}

	return fmt.Sprintf("feeCollectorAccountID: %s, amount: %d", fee.FeeCollectorAccountID.String(), fee.Amount)
}
