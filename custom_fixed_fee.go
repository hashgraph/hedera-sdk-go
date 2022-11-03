package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

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
		FeeCollectorAccountId:  FeeCollectorAccountID,
		AllCollectorsAreExempt: fee.AllCollectorsAreExempt,
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

func (fee *CustomFixedFee) SetAllCollectorsAreExempt(exempt bool) *CustomFixedFee {
	fee.AllCollectorsAreExempt = exempt
	return fee
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
