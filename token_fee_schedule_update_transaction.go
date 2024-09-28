package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2024 Hedera Hashgraph, LLC
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
	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// TokenFeeScheduleUpdateTransaction
// At consensus, updates a token type's fee schedule to the given list of custom fees.
//
// If the target token type has no fee_schedule_key, resolves to TOKEN_HAS_NO_FEE_SCHEDULE_KEY.
// Otherwise this transaction must be signed to the fee_schedule_key, or the transaction will
// resolve to INVALID_SIGNATURE.
//
// If the custom_fees list is empty, clears the fee schedule or resolves to
// CUSTOM_SCHEDULE_ALREADY_HAS_NO_FEES if the fee schedule was already empty.
type TokenFeeScheduleUpdateTransaction struct {
	*Transaction[*TokenFeeScheduleUpdateTransaction]
	tokenID    *TokenID
	customFees []Fee
}

// NewTokenFeeScheduleUpdateTransaction creates TokenFeeScheduleUpdateTransaction which
// at consensus, updates a token type's fee schedule to the given list of custom fees.
//
// If the target token type has no fee_schedule_key, resolves to TOKEN_HAS_NO_FEE_SCHEDULE_KEY.
// Otherwise this transaction must be signed to the fee_schedule_key, or the transaction will
// resolve to INVALID_SIGNATURE.
//
// If the custom_fees list is empty, clears the fee schedule or resolves to
// CUSTOM_SCHEDULE_ALREADY_HAS_NO_FEES if the fee schedule was already empty.
func NewTokenFeeScheduleUpdateTransaction() *TokenFeeScheduleUpdateTransaction {
	tx := &TokenFeeScheduleUpdateTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(5))

	return tx
}

func _TokenFeeScheduleUpdateTransactionFromProtobuf(pb *services.TransactionBody) *TokenFeeScheduleUpdateTransaction {
	customFees := make([]Fee, 0)

	for _, fee := range pb.GetTokenFeeScheduleUpdate().GetCustomFees() {
		customFees = append(customFees, _CustomFeeFromProtobuf(fee))
	}

	return &TokenFeeScheduleUpdateTransaction{
		tokenID:    _TokenIDFromProtobuf(pb.GetTokenFeeScheduleUpdate().TokenId),
		customFees: customFees,
	}
}

// SetTokenID Sets the token whose fee schedule is to be updated
func (tx *TokenFeeScheduleUpdateTransaction) SetTokenID(tokenID TokenID) *TokenFeeScheduleUpdateTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the token whose fee schedule is to be updated
func (tx *TokenFeeScheduleUpdateTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// SetCustomFees Sets the new custom fees to be assessed during a CryptoTransfer that transfers units of this token
func (tx *TokenFeeScheduleUpdateTransaction) SetCustomFees(fees []Fee) *TokenFeeScheduleUpdateTransaction {
	tx._RequireNotFrozen()
	tx.customFees = fees
	return tx
}

// GetCustomFees returns the new custom fees to be assessed during a CryptoTransfer that transfers units of this token
func (tx *TokenFeeScheduleUpdateTransaction) GetCustomFees() []Fee {
	return tx.customFees
}

// ----------- Overridden functions ----------------

func (tx *TokenFeeScheduleUpdateTransaction) getName() string {
	return "TokenFeeScheduleUpdateTransaction"
}

func (tx *TokenFeeScheduleUpdateTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.tokenID != nil {
		if err := tx.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	for _, customFee := range tx.customFees {
		if err := customFee.validateNetworkOnIDs(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *TokenFeeScheduleUpdateTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenFeeScheduleUpdate{
			TokenFeeScheduleUpdate: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenFeeScheduleUpdateTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return nil, errors.New("cannot schedule `TokenFeeScheduleUpdateTransaction")
}

func (tx *TokenFeeScheduleUpdateTransaction) buildProtoBody() *services.TokenFeeScheduleUpdateTransactionBody {
	body := &services.TokenFeeScheduleUpdateTransactionBody{}
	if tx.tokenID != nil {
		body.TokenId = tx.tokenID._ToProtobuf()
	}

	if len(tx.customFees) > 0 {
		for _, customFee := range tx.customFees {
			if body.CustomFees == nil {
				body.CustomFees = make([]*services.CustomFee, 0)
			}
			body.CustomFees = append(body.CustomFees, customFee._ToProtobuf())
		}
	}

	return body
}

func (tx *TokenFeeScheduleUpdateTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().UpdateTokenFeeSchedule,
	}
}

func (tx *TokenFeeScheduleUpdateTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx *TokenFeeScheduleUpdateTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction)
}

func (tx *TokenFeeScheduleUpdateTransaction) setBaseTransaction(baseTx Transaction[TransactionInterface]) {
	tx.Transaction = castFromBaseToConcreteTransaction[*TokenFeeScheduleUpdateTransaction](baseTx)
}
