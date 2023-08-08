package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
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
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

// AccountAllowanceApproveTransaction
// Creates one or more hbar/token approved allowances <b>relative to the owner account specified in the allowances of
// this transaction</b>. Each allowance grants a spender the right to transfer a pre-determined amount of the owner's
// hbar/token to any other account of the spender's choice. If the owner is not specified in any allowance, the payer
// of transaction is considered to be the owner for that particular allowance.
// Setting the amount to zero in CryptoAllowance or TokenAllowance will remove the respective allowance for the spender.
//
// (So if account <tt>0.0.X</tt> pays for this transaction and owner is not specified in the allowance,
// then at consensus each spender account will have new allowances to spend hbar or tokens from <tt>0.0.X</tt>).
type AccountAllowanceApproveTransaction struct {
	Transaction
	hbarAllowances  []*HbarAllowance
	tokenAllowances []*TokenAllowance
	nftAllowances   []*TokenNftAllowance
}

// NewAccountAllowanceApproveTransaction
// Creates an AccountAloowanceApproveTransaction which creates
// one or more hbar/token approved allowances relative to the owner account specified in the allowances of
// this transaction. Each allowance grants a spender the right to transfer a pre-determined amount of the owner's
// hbar/token to any other account of the spender's choice. If the owner is not specified in any allowance, the payer
// of transaction is considered to be the owner for that particular allowance.
// Setting the amount to zero in CryptoAllowance or TokenAllowance will remove the respective allowance for the spender.
//
// (So if account 0.0.X pays for this transaction and owner is not specified in the allowance,
// then at consensus each spender account will have new allowances to spend hbar or tokens from 0.0.X).
func NewAccountAllowanceApproveTransaction() *AccountAllowanceApproveTransaction {
	transaction := AccountAllowanceApproveTransaction{
		Transaction: _NewTransaction(),
	}

	transaction._SetDefaultMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _AccountAllowanceApproveTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *AccountAllowanceApproveTransaction {
	accountApproval := make([]*HbarAllowance, 0)
	tokenApproval := make([]*TokenAllowance, 0)
	nftApproval := make([]*TokenNftAllowance, 0)

	for _, ap := range pb.GetCryptoApproveAllowance().GetCryptoAllowances() {
		temp := _HbarAllowanceFromProtobuf(ap)
		accountApproval = append(accountApproval, &temp)
	}

	for _, ap := range pb.GetCryptoApproveAllowance().GetTokenAllowances() {
		temp := _TokenAllowanceFromProtobuf(ap)
		tokenApproval = append(tokenApproval, &temp)
	}

	for _, ap := range pb.GetCryptoApproveAllowance().GetNftAllowances() {
		temp := _TokenNftAllowanceFromProtobuf(ap)
		nftApproval = append(nftApproval, &temp)
	}

	return &AccountAllowanceApproveTransaction{
		Transaction:     transaction,
		hbarAllowances:  accountApproval,
		tokenAllowances: tokenApproval,
		nftAllowances:   nftApproval,
	}
}

func (transaction *AccountAllowanceApproveTransaction) _ApproveHbarApproval(ownerAccountID *AccountID, id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	transaction._RequireNotFrozen()
	transaction.hbarAllowances = append(transaction.hbarAllowances, &HbarAllowance{
		SpenderAccountID: &id,
		Amount:           amount.AsTinybar(),
		OwnerAccountID:   ownerAccountID,
	})

	return transaction
}

// AddHbarApproval
// Deprecated - Use ApproveHbarAllowance instead
func (transaction *AccountAllowanceApproveTransaction) AddHbarApproval(id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	return transaction._ApproveHbarApproval(nil, id, amount)
}

// ApproveHbarApproval
// Deprecated - Use ApproveHbarAllowance instead
func (transaction *AccountAllowanceApproveTransaction) ApproveHbarApproval(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	return transaction._ApproveHbarApproval(&ownerAccountID, id, amount)
}

// ApproveHbarAllowance
// Approves allowance of hbar transfers for a spender.
func (transaction *AccountAllowanceApproveTransaction) ApproveHbarAllowance(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	return transaction._ApproveHbarApproval(&ownerAccountID, id, amount)
}

// List of hbar allowance records
func (transaction *AccountAllowanceApproveTransaction) GetHbarAllowances() []*HbarAllowance {
	return transaction.hbarAllowances
}

func (transaction *AccountAllowanceApproveTransaction) _ApproveTokenApproval(tokenID TokenID, ownerAccountID *AccountID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	transaction._RequireNotFrozen()
	tokenApproval := TokenAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &accountID,
		Amount:           amount,
		OwnerAccountID:   ownerAccountID,
	}

	transaction.tokenAllowances = append(transaction.tokenAllowances, &tokenApproval)
	return transaction
}

// Deprecated - Use ApproveTokenAllowance instead
func (transaction *AccountAllowanceApproveTransaction) AddTokenApproval(tokenID TokenID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenApproval(tokenID, nil, accountID, amount)
}

// ApproveTokenApproval
// Deprecated - Use ApproveTokenAllowance instead
func (transaction *AccountAllowanceApproveTransaction) ApproveTokenApproval(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenApproval(tokenID, &ownerAccountID, accountID, amount)
}

// ApproveTokenAllowance
// Approve allowance of fungible token transfers for a spender.
func (transaction *AccountAllowanceApproveTransaction) ApproveTokenAllowance(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenApproval(tokenID, &ownerAccountID, accountID, amount)
}

// List of token allowance records
func (transaction *AccountAllowanceApproveTransaction) GetTokenAllowances() []*TokenAllowance {
	return transaction.tokenAllowances
}

func (transaction *AccountAllowanceApproveTransaction) _ApproveTokenNftApproval(nftID NftID, ownerAccountID *AccountID, spenderAccountID *AccountID, delegatingSpenderAccountId *AccountID) *AccountAllowanceApproveTransaction {
	transaction._RequireNotFrozen()

	for _, t := range transaction.nftAllowances {
		if t.TokenID.String() == nftID.TokenID.String() {
			if t.SpenderAccountID.String() == spenderAccountID.String() {
				b := false
				for _, s := range t.SerialNumbers {
					if s == nftID.SerialNumber {
						b = true
					}
				}
				if !b {
					t.SerialNumbers = append(t.SerialNumbers, nftID.SerialNumber)
				}
				return transaction
			}
		}
	}

	transaction.nftAllowances = append(transaction.nftAllowances, &TokenNftAllowance{
		TokenID:           &nftID.TokenID,
		SpenderAccountID:  spenderAccountID,
		SerialNumbers:     []int64{nftID.SerialNumber},
		AllSerials:        false,
		OwnerAccountID:    ownerAccountID,
		DelegatingSpender: delegatingSpenderAccountId,
	})
	return transaction
}

// AddTokenNftApproval
// Deprecated - Use ApproveTokenNftAllowance instead
func (transaction *AccountAllowanceApproveTransaction) AddTokenNftApproval(nftID NftID, accountID AccountID) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenNftApproval(nftID, nil, &accountID, nil)
}

// ApproveTokenNftApproval
// Deprecated - Use ApproveTokenNftAllowance instead
func (transaction *AccountAllowanceApproveTransaction) ApproveTokenNftApproval(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenNftApproval(nftID, &ownerAccountID, &accountID, nil)
}

func (transaction *AccountAllowanceApproveTransaction) ApproveTokenNftAllowanceWithDelegatingSpender(nftID NftID, ownerAccountID AccountID, spenderAccountId AccountID, delegatingSpenderAccountID AccountID) *AccountAllowanceApproveTransaction {
	transaction._RequireNotFrozen()
	return transaction._ApproveTokenNftApproval(nftID, &ownerAccountID, &spenderAccountId, &delegatingSpenderAccountID)
}

// ApproveTokenNftAllowance
// Approve allowance of non-fungible token transfers for a spender.
func (transaction *AccountAllowanceApproveTransaction) ApproveTokenNftAllowance(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenNftApproval(nftID, &ownerAccountID, &accountID, nil)
}

func (transaction *AccountAllowanceApproveTransaction) _ApproveTokenNftAllowanceAllSerials(tokenID TokenID, ownerAccountID *AccountID, spenderAccount AccountID) *AccountAllowanceApproveTransaction {
	for _, t := range transaction.nftAllowances {
		if t.TokenID.String() == tokenID.String() {
			if t.SpenderAccountID.String() == spenderAccount.String() {
				t.SerialNumbers = []int64{}
				t.AllSerials = true
				return transaction
			}
		}
	}

	transaction.nftAllowances = append(transaction.nftAllowances, &TokenNftAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &spenderAccount,
		SerialNumbers:    []int64{},
		AllSerials:       true,
		OwnerAccountID:   ownerAccountID,
	})
	return transaction
}

// AddAllTokenNftApproval
// Approve allowance of non-fungible token transfers for a spender.
// Spender has access to all of the owner's NFT units of type tokenId (currently
// owned and any in the future).
func (transaction *AccountAllowanceApproveTransaction) AddAllTokenNftApproval(tokenID TokenID, spenderAccount AccountID) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenNftAllowanceAllSerials(tokenID, nil, spenderAccount)
}

// ApproveTokenNftAllowanceAllSerials
// Approve allowance of non-fungible token transfers for a spender.
// Spender has access to all of the owner's NFT units of type tokenId (currently
// owned and any in the future).
func (transaction *AccountAllowanceApproveTransaction) ApproveTokenNftAllowanceAllSerials(tokenID TokenID, ownerAccountID AccountID, spenderAccount AccountID) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenNftAllowanceAllSerials(tokenID, &ownerAccountID, spenderAccount)
}

// List of NFT allowance records
func (transaction *AccountAllowanceApproveTransaction) GetTokenNftAllowances() []*TokenNftAllowance {
	return transaction.nftAllowances
}

func (transaction *AccountAllowanceApproveTransaction) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	for _, ap := range transaction.hbarAllowances {
		if ap.SpenderAccountID != nil {
			if err := ap.SpenderAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if ap.OwnerAccountID != nil {
			if err := ap.OwnerAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}

	for _, ap := range transaction.tokenAllowances {
		if ap.SpenderAccountID != nil {
			if err := ap.SpenderAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if ap.TokenID != nil {
			if err := ap.TokenID.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if ap.OwnerAccountID != nil {
			if err := ap.OwnerAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}

	for _, ap := range transaction.nftAllowances {
		if ap.SpenderAccountID != nil {
			if err := ap.SpenderAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if ap.TokenID != nil {
			if err := ap.TokenID.ValidateChecksum(client); err != nil {
				return err
			}
		}

		if ap.OwnerAccountID != nil {
			if err := ap.OwnerAccountID.ValidateChecksum(client); err != nil {
				return err
			}
		}
	}

	return nil
}

func (transaction *AccountAllowanceApproveTransaction) _Build() *services.TransactionBody {
	accountApproval := make([]*services.CryptoAllowance, 0)
	tokenApproval := make([]*services.TokenAllowance, 0)
	nftApproval := make([]*services.NftAllowance, 0)

	for _, ap := range transaction.hbarAllowances {
		accountApproval = append(accountApproval, ap._ToProtobuf())
	}

	for _, ap := range transaction.tokenAllowances {
		tokenApproval = append(tokenApproval, ap._ToProtobuf())
	}

	for _, ap := range transaction.nftAllowances {
		nftApproval = append(nftApproval, ap._ToProtobuf())
	}

	return &services.TransactionBody{
		TransactionID:            transaction.transactionID._ToProtobuf(),
		TransactionFee:           transaction.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(transaction.GetTransactionValidDuration()),
		Memo:                     transaction.Transaction.memo,
		Data: &services.TransactionBody_CryptoApproveAllowance{
			CryptoApproveAllowance: &services.CryptoApproveAllowanceTransactionBody{
				CryptoAllowances: accountApproval,
				NftAllowances:    nftApproval,
				TokenAllowances:  tokenApproval,
			},
		},
	}
}

func (transaction *AccountAllowanceApproveTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *AccountAllowanceApproveTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	accountApproval := make([]*services.CryptoAllowance, 0)
	tokenApproval := make([]*services.TokenAllowance, 0)
	nftApproval := make([]*services.NftAllowance, 0)

	for _, ap := range transaction.hbarAllowances {
		accountApproval = append(accountApproval, ap._ToProtobuf())
	}

	for _, ap := range transaction.tokenAllowances {
		tokenApproval = append(tokenApproval, ap._ToProtobuf())
	}

	for _, ap := range transaction.nftAllowances {
		nftApproval = append(nftApproval, ap._ToProtobuf())
	}

	return &services.SchedulableTransactionBody{
		TransactionFee: transaction.transactionFee,
		Memo:           transaction.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoApproveAllowance{
			CryptoApproveAllowance: &services.CryptoApproveAllowanceTransactionBody{
				CryptoAllowances: accountApproval,
				NftAllowances:    nftApproval,
				TokenAllowances:  tokenApproval,
			},
		},
	}, nil
}

func _AccountApproveAllowanceTransactionGetMethod(request interface{}, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().ApproveAllowances,
	}
}

func (transaction *AccountAllowanceApproveTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *AccountAllowanceApproveTransaction) Sign(
	privateKey PrivateKey,
) *AccountAllowanceApproveTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

// SignWithOperator signs the transaction with client's operator privateKey.
func (transaction *AccountAllowanceApproveTransaction) SignWithOperator(
	client *Client,
) (*AccountAllowanceApproveTransaction, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return transaction, err
		}
	}
	return transaction.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (transaction *AccountAllowanceApproveTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountAllowanceApproveTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *AccountAllowanceApproveTransaction) Execute(
	client *Client,
) (TransactionResponse, error) {
	if client == nil {
		return TransactionResponse{}, errNoClientProvided
	}

	if transaction.freezeError != nil {
		return TransactionResponse{}, transaction.freezeError
	}

	if !transaction.IsFrozen() {
		_, err := transaction.FreezeWith(client)
		if err != nil {
			return TransactionResponse{}, err
		}
	}

	transactionID := transaction.transactionIDs._GetCurrent().(TransactionID)

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		transaction.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := _Execute(
		client,
		&transaction.Transaction,
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_AccountApproveAllowanceTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
		transaction._GetLogID(),
		transaction.grpcDeadline,
		transaction.maxBackoff,
		transaction.minBackoff,
		transaction.maxRetry,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID:  transaction.GetTransactionID(),
			NodeID:         resp.(TransactionResponse).NodeID,
			ValidateStatus: true,
		}, err
	}

	return TransactionResponse{
		TransactionID:  transaction.GetTransactionID(),
		NodeID:         resp.(TransactionResponse).NodeID,
		Hash:           resp.(TransactionResponse).Hash,
		ValidateStatus: true,
	}, nil
}

func (transaction *AccountAllowanceApproveTransaction) Freeze() (*AccountAllowanceApproveTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *AccountAllowanceApproveTransaction) FreezeWith(client *Client) (*AccountAllowanceApproveTransaction, error) {
	if transaction.IsFrozen() {
		return transaction, nil
	}
	transaction._InitFee(client)
	if err := transaction._InitTransactionID(client); err != nil {
		return transaction, err
	}
	err := transaction._ValidateNetworkOnIDs(client)
	body := transaction._Build()
	if err != nil {
		return &AccountAllowanceApproveTransaction{}, err
	}

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

// GetMaxTransactionFee returns the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *AccountAllowanceApproveTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the maximum transaction fee the operator (paying account) is willing to pay.
func (transaction *AccountAllowanceApproveTransaction) SetMaxTransactionFee(fee Hbar) *AccountAllowanceApproveTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *AccountAllowanceApproveTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountAllowanceApproveTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *AccountAllowanceApproveTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

// GetTransactionMemo returns the memo for this AccountAllowanceApproveTransaction.
func (transaction *AccountAllowanceApproveTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountAllowanceApproveTransaction.
func (transaction *AccountAllowanceApproveTransaction) SetTransactionMemo(memo string) *AccountAllowanceApproveTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

// GetTransactionValidDuration returns the duration that this transaction is valid for.
func (transaction *AccountAllowanceApproveTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountAllowanceApproveTransaction.
func (transaction *AccountAllowanceApproveTransaction) SetTransactionValidDuration(duration time.Duration) *AccountAllowanceApproveTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

// GetTransactionID gets the TransactionID for this AccountAllowanceApproveTransaction.
func (transaction *AccountAllowanceApproveTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountAllowanceApproveTransaction.
func (transaction *AccountAllowanceApproveTransaction) SetTransactionID(transactionID TransactionID) *AccountAllowanceApproveTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountAllowanceApproveTransaction.
func (transaction *AccountAllowanceApproveTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountAllowanceApproveTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

// SetMaxRetry sets the max number of errors before execution will fail.
func (transaction *AccountAllowanceApproveTransaction) SetMaxRetry(count int) *AccountAllowanceApproveTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

// AddSignature adds a signature to the Transaction.
func (transaction *AccountAllowanceApproveTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountAllowanceApproveTransaction {
	transaction._RequireOneNodeAccountID()

	if transaction._KeyAlreadySigned(publicKey) {
		return transaction
	}

	if transaction.signedTransactions._Length() == 0 {
		return transaction
	}

	transaction.transactions = _NewLockableSlice()
	transaction.publicKeys = append(transaction.publicKeys, publicKey)
	transaction.transactionSigners = append(transaction.transactionSigners, nil)
	transaction.transactionIDs.locked = true

	for index := 0; index < transaction.signedTransactions._Length(); index++ {
		var temp *services.SignedTransaction
		switch t := transaction.signedTransactions._Get(index).(type) { //nolint
		case *services.SignedTransaction:
			temp = t
		}
		temp.SigMap.SigPair = append(
			temp.SigMap.SigPair,
			publicKey._ToSignaturePairProtobuf(signature),
		)
		transaction.signedTransactions._Set(index, temp)
	}

	return transaction
}

// SetMaxBackoff The maximum amount of time to wait between retries.
// Every retry attempt will increase the wait time exponentially until it reaches this time.
func (transaction *AccountAllowanceApproveTransaction) SetMaxBackoff(max time.Duration) *AccountAllowanceApproveTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

// GetMaxBackoff returns the max back off for this AccountAllowanceApproveTransaction.
func (transaction *AccountAllowanceApproveTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

// SetMinBackoff sets the min back off for this AccountAllowanceApproveTransaction.
func (transaction *AccountAllowanceApproveTransaction) SetMinBackoff(min time.Duration) *AccountAllowanceApproveTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

// GetMinBackoff returns the min back off for this AccountAllowanceApproveTransaction.
func (transaction *AccountAllowanceApproveTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *AccountAllowanceApproveTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountAllowanceApproveTransaction:%d", timestamp.UnixNano())
}
