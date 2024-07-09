package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

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

/**
 * A transaction body to "reject" undesired tokens.<br/>
 * This transaction will transfer one or more tokens or token
 * balances held by the requesting account to the treasury
 * for each token type.
 * <p>
 * Each transfer MUST be one of the following:
 * <ul>
 *   <li>A single non-fungible/unique token.</li>
 *   <li>The full balance held for a fungible/common
 *       token type.</li>
 * </ul>
 * When complete, the requesting account SHALL NOT hold the
 * rejected tokens.<br/>
 * Custom fees and royalties defined for the tokens rejected
 * SHALL NOT be charged for this transaction.
 */
type TokenRejectTransaction struct {
	Transaction
	ownerID  *AccountID
	tokenIDs []TokenID
	nftIDs   []NftID
}

func NewTokenRejectTransaction() *TokenRejectTransaction {
	tx := TokenRejectTransaction{
		Transaction: _NewTransaction(),
	}
	return &tx
}

func _TokenRejectTransactionFromProtobuf(tx Transaction, pb *services.TransactionBody) *TokenRejectTransaction {
	return &TokenRejectTransaction{
		Transaction: tx,
		ownerID:     _AccountIDFromProtobuf(pb.GetTokenReject().Owner),
		tokenIDs:    _TokenIDsFromTokenReferenceProtobuf(pb.GetTokenReject().Rejections),
		nftIDs:      _NftIDsFromTokenReferenceProtobuf(pb.GetTokenReject().Rejections),
	}
}

// SetOwnerID Sets the account which owns the tokens to be rejected
func (tx *TokenRejectTransaction) SetOwnerID(ownerID AccountID) *TokenRejectTransaction {
	tx._RequireNotFrozen()
	tx.ownerID = &ownerID
	return tx
}

// GetOwnerID Gets the account which owns the tokens to be rejected
func (tx *TokenRejectTransaction) GetOwnerID() AccountID {
	if tx.ownerID == nil {
		return AccountID{}
	}
	return *tx.ownerID
}

// SetTokenIDs Sets the tokens to be rejected
func (tx *TokenRejectTransaction) SetTokenIDs(ids ...TokenID) *TokenRejectTransaction {
	tx._RequireNotFrozen()
	tx.tokenIDs = make([]TokenID, len(ids))
	copy(tx.tokenIDs, ids)

	return tx
}

// AddTokenID Adds a token to be rejected
func (tx *TokenRejectTransaction) AddTokenID(id TokenID) *TokenRejectTransaction {
	tx._RequireNotFrozen()
	tx.tokenIDs = append(tx.tokenIDs, id)
	return tx
}

// GetTokenIDs Gets the tokens to be rejected
func (tx *TokenRejectTransaction) GetTokenIDs() []TokenID {
	return tx.tokenIDs
}

// SetNftIDs Sets the NFTs to be rejected
func (tx *TokenRejectTransaction) SetNftIDs(ids ...NftID) *TokenRejectTransaction {
	tx._RequireNotFrozen()
	tx.nftIDs = make([]NftID, len(ids))
	copy(tx.nftIDs, ids)

	return tx
}

// AddNftID Adds an NFT to be rejected
func (tx *TokenRejectTransaction) AddNftID(id NftID) *TokenRejectTransaction {
	tx._RequireNotFrozen()
	tx.nftIDs = append(tx.nftIDs, id)
	return tx
}

// GetNftIDs Gets the NFTs to be rejected
func (tx *TokenRejectTransaction) GetNftIDs() []NftID {
	return tx.nftIDs
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenRejectTransaction) Sign(privateKey PrivateKey) *TokenRejectTransaction {
	tx.Transaction.Sign(privateKey)
	return tx
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (tx *TokenRejectTransaction) SignWithOperator(client *Client) (*TokenRejectTransaction, error) {
	_, err := tx.Transaction.signWithOperator(client, tx)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenRejectTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenRejectTransaction {
	tx.Transaction.SignWith(publicKey, signer)
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenRejectTransaction) AddSignature(publicKey PublicKey, signature []byte) *TokenRejectTransaction {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenRejectTransaction) SetGrpcDeadline(deadline *time.Duration) *TokenRejectTransaction {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenRejectTransaction) Freeze() (*TokenRejectTransaction, error) {
	return tx.FreezeWith(nil)
}

func (tx *TokenRejectTransaction) FreezeWith(client *Client) (*TokenRejectTransaction, error) {
	_, err := tx.Transaction.freezeWith(client, tx)
	return tx, err
}

// SetMaxTransactionFee sets the max transaction fee for this TokenRejectTransaction.
func (tx *TokenRejectTransaction) SetMaxTransactionFee(fee Hbar) *TokenRejectTransaction {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenRejectTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenRejectTransaction {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenRejectTransaction.
func (tx *TokenRejectTransaction) SetTransactionMemo(memo string) *TokenRejectTransaction {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenRejectTransaction.
func (tx *TokenRejectTransaction) SetTransactionValidDuration(duration time.Duration) *TokenRejectTransaction {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *TokenRejectTransaction) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this TokenRejectTransaction.
func (tx *TokenRejectTransaction) SetTransactionID(transactionID TransactionID) *TokenRejectTransaction {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenRejectTransaction.
func (tx *TokenRejectTransaction) SetNodeAccountIDs(nodeID []AccountID) *TokenRejectTransaction {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenRejectTransaction) SetMaxRetry(count int) *TokenRejectTransaction {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenRejectTransaction) SetMaxBackoff(max time.Duration) *TokenRejectTransaction {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenRejectTransaction) SetMinBackoff(min time.Duration) *TokenRejectTransaction {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenRejectTransaction) SetLogLevel(level LogLevel) *TokenRejectTransaction {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenRejectTransaction) Execute(client *Client) (TransactionResponse, error) {
	return tx.Transaction.execute(client, tx)
}

func (tx *TokenRejectTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}

// ----------- Overridden functions ----------------

func (tx *TokenRejectTransaction) getName() string {
	return "TokenRejectTransaction"
}

func (tx *TokenRejectTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if tx.ownerID != nil {
		if err := tx.ownerID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	for _, tokenID := range tx.tokenIDs {
		if err := tokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	for _, nftID := range tx.nftIDs {
		if err := nftID.TokenID.ValidateChecksum(client); err != nil {
			return err
		}
	}

	return nil
}

func (tx *TokenRejectTransaction) build() *services.TransactionBody {
	body := tx.buildProtoBody()

	return &services.TransactionBody{
		TransactionFee:           tx.transactionFee,
		Memo:                     tx.Transaction.memo,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		TransactionID:            tx.transactionID._ToProtobuf(),
		Data: &services.TransactionBody_TokenReject{
			TokenReject: body,
		},
	}
}

func (tx *TokenRejectTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_TokenReject{
			TokenReject: tx.buildProtoBody(),
		},
	}, nil
}

func (tx *TokenRejectTransaction) buildProtoBody() *services.TokenRejectTransactionBody {
	body := &services.TokenRejectTransactionBody{}

	if tx.ownerID != nil {
		body.Owner = tx.ownerID._ToProtobuf()
	}

	if len(tx.tokenIDs) > 0 {
		for _, tokenID := range tx.tokenIDs {
			tokenReference := &services.TokenReference_FungibleToken{
				FungibleToken: tokenID._ToProtobuf(),
			}

			body.Rejections = append(body.Rejections, &services.TokenReference{
				TokenIdentifier: tokenReference,
			})
		}
	}

	if len(tx.nftIDs) > 0 {
		for _, nftID := range tx.nftIDs {
			tokenReference := &services.TokenReference_Nft{
				Nft: nftID._ToProtobuf(),
			}

			body.Rejections = append(body.Rejections, &services.TokenReference{
				TokenIdentifier: tokenReference,
			})
		}
	}

	return body
}

func (tx *TokenRejectTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetToken().RejectToken,
	}
}

func (tx *TokenRejectTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}
