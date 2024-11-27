package hiero

// SPDX-License-Identifier: Apache-2.0

type TokenRejectFlow struct {
	ownerID           *AccountID
	tokenIDs          []TokenID
	nftIDs            []NftID
	freezeWithClient  *Client
	signPrivateKey    *PrivateKey
	signPublicKey     *PublicKey
	transactionSigner *TransactionSigner
}

func NewTokenRejectFlow() *TokenRejectFlow {
	tx := TokenRejectFlow{}
	return &tx
}

// SetOwnerID Sets the account which owns the tokens to be rejected
func (tx *TokenRejectFlow) SetOwnerID(ownerID AccountID) *TokenRejectFlow {
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
	tx.tokenIDs = make([]TokenID, len(ids))
	copy(tx.tokenIDs, ids)

	return tx
}

// AddTokenID Adds a token to be rejected
func (tx *TokenRejectFlow) AddTokenID(id TokenID) *TokenRejectFlow {
	tx.tokenIDs = append(tx.tokenIDs, id)
	return tx
}

// GetTokenIDs Gets the tokens to be rejected
func (tx *TokenRejectFlow) GetTokenIDs() []TokenID {
	return tx.tokenIDs
}

// SetNftIDs Sets the NFTs to be rejected
func (tx *TokenRejectFlow) SetNftIDs(ids ...NftID) *TokenRejectFlow {
	tx.nftIDs = make([]NftID, len(ids))
	copy(tx.nftIDs, ids)

	return tx
}

// AddNftID Adds an NFT to be rejected
func (tx *TokenRejectFlow) AddNftID(id NftID) *TokenRejectFlow {
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

func (tx *TokenRejectFlow) FreezeWith(client *Client) (*TokenRejectFlow, error) {
	tx.freezeWithClient = client
	return tx, nil
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
		_, err := tokenDissociateTxn.FreezeWith(tx.freezeWithClient)
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
		_, err := tokenRejectTxn.FreezeWith(tx.freezeWithClient)
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
