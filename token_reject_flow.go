package hedera

import (
	"time"
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

type TokenRejectFlow struct {
	Transaction
	ownerID           *AccountID
	tokenIDs          []TokenID
	nftIDs            []NftID
	freezeWithClient  *Client
	signPrivateKey    *PrivateKey
	signPublicKey     *PublicKey
	transactionSigner *TransactionSigner
}

func NewTokenRejectFlow() *TokenRejectFlow {
	tx := TokenRejectFlow{
		Transaction: _NewTransaction(),
	}
	return &tx
}

// SetOwnerID Sets the account which owns the tokens to be rejected
func (tx *TokenRejectFlow) SetOwnerID(ownerID AccountID) *TokenRejectFlow {
	tx._RequireNotFrozen()
	tx.ownerID = &ownerID
	return tx
}

// GetOwnerID Gets the account which owns the tokens to be rejected
func (tx *TokenRejectFlow) GetOwnerID() AccountID {
	if tx.ownerID == nil {
		return AccountID{}
	}
	return *tx.ownerID
}

// SetTokenIDs Sets the tokens to be rejected
func (tx *TokenRejectFlow) SetTokenIDs(ids ...TokenID) *TokenRejectFlow {
	tx._RequireNotFrozen()
	tx.tokenIDs = make([]TokenID, len(ids))
	copy(tx.tokenIDs, ids)

	return tx
}

// AddTokenID Adds a token to be rejected
func (tx *TokenRejectFlow) AddTokenID(id TokenID) *TokenRejectFlow {
	tx._RequireNotFrozen()
	tx.tokenIDs = append(tx.tokenIDs, id)
	return tx
}

// GetTokenIDs Gets the tokens to be rejected
func (tx *TokenRejectFlow) GetTokenIDs() []TokenID {
	return tx.tokenIDs
}

// SetNftIDs Sets the NFTs to be rejected
func (tx *TokenRejectFlow) SetNftIDs(ids ...NftID) *TokenRejectFlow {
	tx._RequireNotFrozen()
	tx.nftIDs = make([]NftID, len(ids))
	copy(tx.nftIDs, ids)

	return tx
}

// AddNftID Adds an NFT to be rejected
func (tx *TokenRejectFlow) AddNftID(id NftID) *TokenRejectFlow {
	tx._RequireNotFrozen()
	tx.nftIDs = append(tx.nftIDs, id)
	return tx
}

// GetNftIDs Gets the NFTs to be rejected
func (tx *TokenRejectFlow) GetNftIDs() []NftID {
	return tx.nftIDs
}

// ---- Required Interfaces ---- //

// Sign uses the provided privateKey to sign the transaction.
func (tx *TokenRejectFlow) Sign(privateKey PrivateKey) *TokenRejectFlow {
	tx.signPrivateKey = &privateKey
	return tx
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (tx *TokenRejectFlow) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *TokenRejectFlow {
	tx.signPublicKey = &publicKey
	tx.transactionSigner = &signer
	return tx
}

// AddSignature adds a signature to the transaction.
func (tx *TokenRejectFlow) AddSignature(publicKey PublicKey, signature []byte) *TokenRejectFlow {
	tx.Transaction.AddSignature(publicKey, signature)
	return tx
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (tx *TokenRejectFlow) SetGrpcDeadline(deadline *time.Duration) *TokenRejectFlow {
	tx.Transaction.SetGrpcDeadline(deadline)
	return tx
}

func (tx *TokenRejectFlow) Freeze() (*TokenRejectFlow, error) {
	return tx.FreezeWith(nil), nil
}

func (tx *TokenRejectFlow) FreezeWith(client *Client) *TokenRejectFlow {
	tx.freezeWithClient = client
	return tx
}

// SetMaxTransactionFee sets the max transaction fee for this TokenRejectFlow.
func (tx *TokenRejectFlow) SetMaxTransactionFee(fee Hbar) *TokenRejectFlow {
	tx.Transaction.SetMaxTransactionFee(fee)
	return tx
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (tx *TokenRejectFlow) SetRegenerateTransactionID(regenerateTransactionID bool) *TokenRejectFlow {
	tx.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return tx
}

// SetTransactionMemo sets the memo for this TokenRejectFlow.
func (tx *TokenRejectFlow) SetTransactionMemo(memo string) *TokenRejectFlow {
	tx.Transaction.SetTransactionMemo(memo)
	return tx
}

// SetTransactionValidDuration sets the valid duration for this TokenRejectFlow.
func (tx *TokenRejectFlow) SetTransactionValidDuration(duration time.Duration) *TokenRejectFlow {
	tx.Transaction.SetTransactionValidDuration(duration)
	return tx
}

// ToBytes serialise the tx to bytes, no matter if it is signed (locked), or not
func (tx *TokenRejectFlow) ToBytes() ([]byte, error) {
	bytes, err := tx.Transaction.toBytes(tx)
	if err != nil {
		return nil, err
	}
	return bytes, nil
}

// SetTransactionID sets the TransactionID for this TokenRejectFlow.
func (tx *TokenRejectFlow) SetTransactionID(transactionID TransactionID) *TokenRejectFlow {
	tx.Transaction.SetTransactionID(transactionID)
	return tx
}

// SetNodeAccountIDs sets the _Node AccountID for this TokenRejectFlow.
func (tx *TokenRejectFlow) SetNodeAccountIDs(nodeID []AccountID) *TokenRejectFlow {
	tx.Transaction.SetNodeAccountIDs(nodeID)
	return tx
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (tx *TokenRejectFlow) SetMaxRetry(count int) *TokenRejectFlow {
	tx.Transaction.SetMaxRetry(count)
	return tx
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (tx *TokenRejectFlow) SetMaxBackoff(max time.Duration) *TokenRejectFlow {
	tx.Transaction.SetMaxBackoff(max)
	return tx
}

// SetMinBackoff sets the minimum amount of time to wait between retries.
func (tx *TokenRejectFlow) SetMinBackoff(min time.Duration) *TokenRejectFlow {
	tx.Transaction.SetMinBackoff(min)
	return tx
}

func (tx *TokenRejectFlow) SetLogLevel(level LogLevel) *TokenRejectFlow {
	tx.Transaction.SetLogLevel(level)
	return tx
}

func (tx *TokenRejectFlow) _CreateTokenDissociateTransaction(client *Client) (*TokenDissociateTransaction, error) {
	if client == nil {
		return &TokenDissociateTransaction{}, nil
	}
	tokenDissociateTxn := NewTokenDissociateTransaction()

	if tx.ownerID != nil {
		tokenDissociateTxn.SetAccountID(*tx.ownerID)
	}

	tokenIDs := make([]TokenID, 0)
	if tx.tokenIDs != nil {
		tokenIDs = append(tokenIDs, tx.tokenIDs...)
	}

	if tx.nftIDs != nil {
		seenTokenIDs := make(map[TokenID]struct{})
		for _, nftID := range tx.nftIDs {
			if _, exists := seenTokenIDs[nftID.TokenID]; !exists {
				seenTokenIDs[nftID.TokenID] = struct{}{}
				tokenIDs = append(tokenIDs, nftID.TokenID)
			}
		}
	}

	if len(tokenIDs) != 0 {
		tokenDissociateTxn.SetTokenIDs(tokenIDs...)
	}

	if tx.freezeWithClient != nil {
		_, err := tokenDissociateTxn.freezeWith(tx.freezeWithClient, tokenDissociateTxn)
		if err != nil {
			return nil, err
		}
	}

	if tx.signPrivateKey != nil {
		tokenDissociateTxn = tokenDissociateTxn.Sign(*tx.signPrivateKey)
	}

	if tx.signPublicKey != nil && tx.transactionSigner != nil {
		tokenDissociateTxn = tokenDissociateTxn.SignWith(*tx.signPublicKey, *tx.transactionSigner)
	}

	return tokenDissociateTxn, nil
}

func (tx *TokenRejectFlow) _CreateTokenRejectTransaction(client *Client) (*TokenRejectTransaction, error) {
	if client == nil {
		return &TokenRejectTransaction{}, nil
	}
	tokenRejectTxn := NewTokenRejectTransaction()

	if tx.ownerID != nil {
		tokenRejectTxn.SetOwnerID(*tx.ownerID)
	}

	if tx.tokenIDs != nil {
		tokenRejectTxn.SetTokenIDs(tx.tokenIDs...)
	}

	if tx.nftIDs != nil {
		tokenRejectTxn.SetNftIDs(tx.nftIDs...)
	}

	if tx.freezeWithClient != nil {
		_, err := tokenRejectTxn.freezeWith(tx.freezeWithClient, tokenRejectTxn)
		if err != nil {
			return nil, err
		}
	}

	if tx.signPrivateKey != nil {
		tokenRejectTxn = tokenRejectTxn.Sign(*tx.signPrivateKey)
	}

	if tx.signPublicKey != nil && tx.transactionSigner != nil {
		tokenRejectTxn = tokenRejectTxn.SignWith(*tx.signPublicKey, *tx.transactionSigner)
	}

	return tokenRejectTxn, nil
}

func (tx *TokenRejectFlow) Execute(client *Client) (TransactionResponse, error) {
	tokenRejectTxn, err := tx._CreateTokenRejectTransaction(client)
	if err != nil {
		return TransactionResponse{}, err
	}
	tokenRejectResponse, err := tokenRejectTxn.Execute(client)
	if err != nil {
		return TransactionResponse{}, err
	}
	_, err = tokenRejectResponse.GetReceipt(client)
	if err != nil {
		return TransactionResponse{}, err
	}

	tokenDissociateTxn, err := tx._CreateTokenDissociateTransaction(client)
	if err != nil {
		return TransactionResponse{}, err
	}
	tokenDissociateResponse, err := tokenDissociateTxn.Execute(client)
	if err != nil {
		return TransactionResponse{}, err
	}
	_, err = tokenDissociateResponse.GetReceipt(client)
	if err != nil {
		return TransactionResponse{}, err
	}

	return tokenRejectResponse, nil
}

func (tx *TokenRejectFlow) Schedule() (*ScheduleCreateTransaction, error) {
	return tx.Transaction.schedule(tx)
}
func (tx *TokenRejectFlow) validateNetworkOnIDs(client *Client) error {
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
