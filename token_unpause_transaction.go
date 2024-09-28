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

// TokenUnpauseTransaction
// Unpauses the Token. Must be signed with the Token's pause key.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If no Pause Key is defined, the transaction will resolve to TOKEN_HAS_NO_PAUSE_KEY.
// Once executed the Token is marked as Unpaused and can be used in Transactions.
// The operation is idempotent - becomes a no-op if the Token is already unpaused.
type TokenUnpauseTransaction struct {
	*Transaction[*TokenUnpauseTransaction]
	tokenID *TokenID
}

// NewTokenUnpauseTransaction creates TokenUnpauseTransaction which unpauses the Token.
// Must be signed with the Token's pause key.
// If the provided token is not found, the transaction will resolve to INVALID_TOKEN_ID.
// If the provided token has been deleted, the transaction will resolve to TOKEN_WAS_DELETED.
// If no Pause Key is defined, the transaction will resolve to TOKEN_HAS_NO_PAUSE_KEY.
// Once executed the Token is marked as Unpaused and can be used in Transactions.
// The operation is idempotent - becomes a no-op if the Token is already unpaused.
func NewTokenUnpauseTransaction() *TokenUnpauseTransaction {
	tx := &TokenUnpauseTransaction{}
	tx.Transaction = _NewTransaction(tx)

	tx._SetDefaultMaxTransactionFee(NewHbar(30))

	return tx
}

func _TokenUnpauseTransactionFromProtobuf(pb *services.TransactionBody) *TokenUnpauseTransaction {
	return &TokenUnpauseTransaction{
		tokenID: _TokenIDFromProtobuf(pb.GetTokenDeletion().GetToken()),
	}
}

// SetTokenID Sets the token to be unpaused.
func (tx *TokenUnpauseTransaction) SetTokenID(tokenID TokenID) *TokenUnpauseTransaction {
	tx._RequireNotFrozen()
	tx.tokenID = &tokenID
	return tx
}

// GetTokenID returns the TokenID of the token to be unpaused.
func (tx *TokenUnpauseTransaction) GetTokenID() TokenID {
	if tx.tokenID == nil {
		return TokenID{}
	}

	return *tx.tokenID
}

// ----------- Overridden functions ----------------

func (tx *TokenUnpauseTransaction) getName() string {
	return "TokenUnpauseTransaction"
}

func (tx *TokenUnpauseTransaction) validateNetworkOnIDs(client *Client) error {
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

func (tx *TokenUnpauseTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenUnpause{
			TokenUnpause: tx.buildProtoBody(),
		},
	}
}

func (tx *TokenUnpauseTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) { //nolint
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenUnpause{
			TokenUnpause: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenUnpauseTransaction) buildProtoBody() *services.TokenUnpauseTransactionBody { //nolint
	body := &services.TokenUnpauseTransactionBody{}
	if tx.tokenID != nil {
		body.Token = tx.tokenID._ToProtobuf()
	}

	return body
}

func (tx *TokenUnpauseTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().DeleteToken,
	}
}

func (tx *TokenUnpauseTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx *TokenUnpauseTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction)
}

func (tx *TokenUnpauseTransaction) setBaseTransaction(baseTx Transaction[TransactionInterface]) {
	tx.Transaction = castFromBaseToConcreteTransaction[*TokenUnpauseTransaction](baseTx)
}
