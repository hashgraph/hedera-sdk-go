package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
)

// AccountAllowanceApproveTransaction
// Creates one or more hbar/token approved allowances <b>relative to the owner account specified in the allowances of
// tx transaction</b>. Each allowance grants a spender the right to transfer a pre-determined amount of the owner's
// hbar/token to any other account of the spender's choice. If the owner is not specified in any allowance, the payer
// of transaction is considered to be the owner for that particular allowance.
// Setting the amount to zero in CryptoAllowance or TokenAllowance will remove the respective allowance for the spender.
//
// (So if account <tt>0.0.X</tt> pays for this transaction and owner is not specified in the allowance,
// then at consensus each spender account will have new allowances to spend hbar or tokens from <tt>0.0.X</tt>).
type AccountAllowanceApproveTransaction struct {
	*Transaction[*AccountAllowanceApproveTransaction]
	hbarAllowances  []*HbarAllowance
	tokenAllowances []*TokenAllowance
	nftAllowances   []*TokenNftAllowance
}

// NewAccountAllowanceApproveTransaction
// Creates an AccountAloowanceApproveTransaction which creates
// one or more hbar/token approved allowances relative to the owner account specified in the allowances of
// tx transaction. Each allowance grants a spender the right to transfer a pre-determined amount of the owner's
// hbar/token to any other account of the spender's choice. If the owner is not specified in any allowance, the payer
// of transaction is considered to be the owner for that particular allowance.
// Setting the amount to zero in CryptoAllowance or TokenAllowance will remove the respective allowance for the spender.
//
// (So if account 0.0.X pays for this transaction and owner is not specified in the allowance,
// then at consensus each spender account will have new allowances to spend hbar or tokens from 0.0.X).
func NewAccountAllowanceApproveTransaction() *AccountAllowanceApproveTransaction {
	tx := &AccountAllowanceApproveTransaction{}
	tx.Transaction = _NewTransaction(tx)
	tx._SetDefaultMaxTransactionFee(NewHbar(2))

	return tx
}

func _AccountAllowanceApproveTransactionFromProtobuf(tx Transaction[*AccountAllowanceApproveTransaction], pb *services.TransactionBody) AccountAllowanceApproveTransaction {
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

	accountAllowanceAppoveTransaction := AccountAllowanceApproveTransaction{
		hbarAllowances:  accountApproval,
		tokenAllowances: tokenApproval,
		nftAllowances:   nftApproval,
	}
	tx.childTransaction = &accountAllowanceAppoveTransaction
	accountAllowanceAppoveTransaction.Transaction = &tx
	return accountAllowanceAppoveTransaction
}

func (tx *AccountAllowanceApproveTransaction) _ApproveHbarApproval(ownerAccountID *AccountID, id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	tx._RequireNotFrozen()
	tx.hbarAllowances = append(tx.hbarAllowances, &HbarAllowance{
		SpenderAccountID: &id,
		Amount:           amount.AsTinybar(),
		OwnerAccountID:   ownerAccountID,
	})

	return tx
}

// AddHbarApproval
// Deprecated - Use ApproveHbarAllowance instead
func (tx *AccountAllowanceApproveTransaction) AddHbarApproval(id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	return tx._ApproveHbarApproval(nil, id, amount)
}

// ApproveHbarApproval
// Deprecated - Use ApproveHbarAllowance instead
func (tx *AccountAllowanceApproveTransaction) ApproveHbarApproval(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	return tx._ApproveHbarApproval(&ownerAccountID, id, amount)
}

// ApproveHbarAllowance
// Approves allowance of hbar transfers for a spender.
func (tx *AccountAllowanceApproveTransaction) ApproveHbarAllowance(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	return tx._ApproveHbarApproval(&ownerAccountID, id, amount)
}

// List of hbar allowance records
func (tx *AccountAllowanceApproveTransaction) GetHbarAllowances() []*HbarAllowance {
	return tx.hbarAllowances
}

func (tx *AccountAllowanceApproveTransaction) _ApproveTokenApproval(tokenID TokenID, ownerAccountID *AccountID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	tx._RequireNotFrozen()
	tokenApproval := TokenAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &accountID,
		Amount:           amount,
		OwnerAccountID:   ownerAccountID,
	}

	tx.tokenAllowances = append(tx.tokenAllowances, &tokenApproval)
	return tx
}

// Deprecated - Use ApproveTokenAllowance instead
func (tx *AccountAllowanceApproveTransaction) AddTokenApproval(tokenID TokenID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	return tx._ApproveTokenApproval(tokenID, nil, accountID, amount)
}

// ApproveTokenApproval
// Deprecated - Use ApproveTokenAllowance instead
func (tx *AccountAllowanceApproveTransaction) ApproveTokenApproval(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	return tx._ApproveTokenApproval(tokenID, &ownerAccountID, accountID, amount)
}

// ApproveTokenAllowance
// Approve allowance of fungible token transfers for a spender.
func (tx *AccountAllowanceApproveTransaction) ApproveTokenAllowance(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	return tx._ApproveTokenApproval(tokenID, &ownerAccountID, accountID, amount)
}

// List of token allowance records
func (tx *AccountAllowanceApproveTransaction) GetTokenAllowances() []*TokenAllowance {
	return tx.tokenAllowances
}

func (tx *AccountAllowanceApproveTransaction) _ApproveTokenNftApproval(nftID NftID, ownerAccountID *AccountID, spenderAccountID *AccountID, delegatingSpenderAccountId *AccountID) *AccountAllowanceApproveTransaction {
	tx._RequireNotFrozen()

	for _, t := range tx.nftAllowances {
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
				return tx
			}
		}
	}

	tx.nftAllowances = append(tx.nftAllowances, &TokenNftAllowance{
		TokenID:           &nftID.TokenID,
		SpenderAccountID:  spenderAccountID,
		SerialNumbers:     []int64{nftID.SerialNumber},
		AllSerials:        false,
		OwnerAccountID:    ownerAccountID,
		DelegatingSpender: delegatingSpenderAccountId,
	})
	return tx
}

// AddTokenNftApproval
// Deprecated - Use ApproveTokenNftAllowance instead
func (tx *AccountAllowanceApproveTransaction) AddTokenNftApproval(nftID NftID, accountID AccountID) *AccountAllowanceApproveTransaction {
	return tx._ApproveTokenNftApproval(nftID, nil, &accountID, nil)
}

// ApproveTokenNftApproval
// Deprecated - Use ApproveTokenNftAllowance instead
func (tx *AccountAllowanceApproveTransaction) ApproveTokenNftApproval(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceApproveTransaction {
	return tx._ApproveTokenNftApproval(nftID, &ownerAccountID, &accountID, nil)
}

func (tx *AccountAllowanceApproveTransaction) ApproveTokenNftAllowanceWithDelegatingSpender(nftID NftID, ownerAccountID AccountID, spenderAccountId AccountID, delegatingSpenderAccountID AccountID) *AccountAllowanceApproveTransaction {
	tx._RequireNotFrozen()
	return tx._ApproveTokenNftApproval(nftID, &ownerAccountID, &spenderAccountId, &delegatingSpenderAccountID)
}

// ApproveTokenNftAllowance
// Approve allowance of non-fungible token transfers for a spender.
func (tx *AccountAllowanceApproveTransaction) ApproveTokenNftAllowance(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceApproveTransaction {
	return tx._ApproveTokenNftApproval(nftID, &ownerAccountID, &accountID, nil)
}

func (tx *AccountAllowanceApproveTransaction) _ApproveTokenNftAllowanceAllSerials(tokenID TokenID, ownerAccountID *AccountID, spenderAccount AccountID) *AccountAllowanceApproveTransaction {
	for _, t := range tx.nftAllowances {
		if t.TokenID.String() == tokenID.String() {
			if t.SpenderAccountID.String() == spenderAccount.String() {
				t.SerialNumbers = []int64{}
				t.AllSerials = true
				return tx
			}
		}
	}

	tx.nftAllowances = append(tx.nftAllowances, &TokenNftAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &spenderAccount,
		SerialNumbers:    []int64{},
		AllSerials:       true,
		OwnerAccountID:   ownerAccountID,
	})
	return tx
}

// AddAllTokenNftApproval
// Approve allowance of non-fungible token transfers for a spender.
// Spender has access to all of the owner's NFT units of type tokenId (currently
// owned and any in the future).
func (tx *AccountAllowanceApproveTransaction) AddAllTokenNftApproval(tokenID TokenID, spenderAccount AccountID) *AccountAllowanceApproveTransaction {
	return tx._ApproveTokenNftAllowanceAllSerials(tokenID, nil, spenderAccount)
}

// ApproveTokenNftAllowanceAllSerials
// Approve allowance of non-fungible token transfers for a spender.
// Spender has access to all of the owner's NFT units of type tokenId (currently
// owned and any in the future).
func (tx *AccountAllowanceApproveTransaction) ApproveTokenNftAllowanceAllSerials(tokenID TokenID, ownerAccountID AccountID, spenderAccount AccountID) *AccountAllowanceApproveTransaction {
	return tx._ApproveTokenNftAllowanceAllSerials(tokenID, &ownerAccountID, spenderAccount)
}

// List of NFT allowance records
func (tx *AccountAllowanceApproveTransaction) GetTokenNftAllowances() []*TokenNftAllowance {
	return tx.nftAllowances
}

// ----------- Overridden functions ----------------

func (tx AccountAllowanceApproveTransaction) getName() string {
	return "AccountAllowanceApproveTransaction"
}
func (tx AccountAllowanceApproveTransaction) validateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	for _, ap := range tx.hbarAllowances {
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

	for _, ap := range tx.tokenAllowances {
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

	for _, ap := range tx.nftAllowances {
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

func (tx AccountAllowanceApproveTransaction) build() *services.TransactionBody {
	return &services.TransactionBody{
		TransactionID:            tx.transactionID._ToProtobuf(),
		TransactionFee:           tx.transactionFee,
		TransactionValidDuration: _DurationToProtobuf(tx.GetTransactionValidDuration()),
		Memo:                     tx.Transaction.memo,
		Data: &services.TransactionBody_CryptoApproveAllowance{
			CryptoApproveAllowance: tx.buildProtoBody(),
		},
	}
}

func (tx AccountAllowanceApproveTransaction) buildScheduled() (*services.SchedulableTransactionBody, error) {
	return &services.SchedulableTransactionBody{
		TransactionFee: tx.transactionFee,
		Memo:           tx.Transaction.memo,
		Data: &services.SchedulableTransactionBody_CryptoApproveAllowance{
			CryptoApproveAllowance: tx.buildProtoBody(),
		},
	}, nil
}

func (tx AccountAllowanceApproveTransaction) buildProtoBody() *services.CryptoApproveAllowanceTransactionBody {
	body := &services.CryptoApproveAllowanceTransactionBody{
		CryptoAllowances: make([]*services.CryptoAllowance, 0),
		TokenAllowances:  make([]*services.TokenAllowance, 0),
		NftAllowances:    make([]*services.NftAllowance, 0),
	}

	for _, ap := range tx.hbarAllowances {
		body.CryptoAllowances = append(body.CryptoAllowances, ap._ToProtobuf())
	}

	for _, ap := range tx.tokenAllowances {
		body.TokenAllowances = append(body.TokenAllowances, ap._ToProtobuf())
	}

	for _, ap := range tx.nftAllowances {
		body.NftAllowances = append(body.NftAllowances, ap._ToProtobuf())
	}

	return body
}

func (tx AccountAllowanceApproveTransaction) getMethod(channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().ApproveAllowances,
	}
}

func (tx AccountAllowanceApproveTransaction) constructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
	return tx.buildScheduled()
}

func (tx AccountAllowanceApproveTransaction) getBaseTransaction() *Transaction[TransactionInterface] {
	return castFromConcreteToBaseTransaction(tx.Transaction, &tx)
}
