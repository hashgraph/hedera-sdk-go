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
	"github.com/hashgraph/hedera-protobufs-go/services"
)

// TokenPauseTransaction
// Pauses the Token from being involved in any kind of Transaction until it is unpaused.
// Must be signed with the Token's pause key.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If no Pause Key is defined, the transaction will resolve to TOKEN_HAS_NO_PAUSE_KEY.
// Once executed the Token is marked as paused and will be not able to be a part of any transaction.
// The operation is idempotent - becomes a no-op if the Token is already Paused.
type TokenPauseTransaction struct {
	*Transaction[*TokenPauseTransaction]
	tokenID *TokenID
}

// NewTokenPauseTransaction creates TokenPauseTransaction which
// pauses the Token from being involved in any kind of Transaction until it is unpaused.
// Must be signed with the Token's pause key.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If no Pause Key is defined, the transaction will resolve to TOKEN_HAS_NO_PAUSE_KEY.
// Once executed the Token is marked as paused and will be not able to be a part of any transaction.
// The operation is idempotent - becomes a no-op if the Token is already Paused.
func NewTokenPauseTransaction() *TokenPauseTransaction {
	tx := &TokenPauseTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return tx
}

func _TokenPauseTransactionFromProtobuf(tx Transaction[*TokenPauseTransaction], pb *services.TransactionBody) *TokenPauseTransaction {
	return &TokenPauseTransaction{
		Transaction: &tx,
		tokenID:     _TokenIDFromProtobuf(pb.GetTokenDeletion().GetToken()),
	}
}

// SetTokenID Sets the token to be paused
func (tx *TokenPauseTransaction) SetTokenID(tokenID TokenID) *TokenPauseTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the token to be paused
func (tx *TokenPauseTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// ----------- Overridden functions ----------------

func (tx *TokenPauseTransaction) getName() string {
	return "TokenPauseTransaction"
}

func (tx *TokenPauseTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.tokenID != nil {
		if err := tx.tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *TokenPauseTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenPause{
			TokenPause: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenPauseTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) { //nolint
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenPause{
			TokenPause: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenPauseTransaction) buildProtoBody() *services.TokenPauseTransactionBody {
	body := &services.TokenPauseTransactionBody{}
	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}
	return body
}

func (tx *TokenPauseTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().DeleteToken,
	}
}

func (tx *TokenPauseTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx *TokenPauseTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction[*TokenPauseTransaction](tx.Transaction)
}
