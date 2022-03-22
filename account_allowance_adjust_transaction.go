package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type AccountAllowanceAdjustTransaction struct {
	Transaction
	hbarAllowances  []*HbarAllowance
	tokenAllowances []*TokenAllowance
	nftAllowances   []*TokenNftAllowance
}

func NewAccountAllowanceAdjustTransaction() *AccountAllowanceAdjustTransaction {
	transaction := AccountAllowanceAdjustTransaction{
		Transaction: _NewTransaction(),
	}

	transaction.SetMaxTransactionFee(NewHbar(2))

	return &transaction
}

func _AccountAllowanceAdjustTransactionFromProtobuf(transaction Transaction, pb *services.TransactionBody) *AccountAllowanceAdjustTransaction {
	accountApproval := make([]*HbarAllowance, 0)
	tokenApproval := make([]*TokenAllowance, 0)
	nftApproval := make([]*TokenNftAllowance, 0)

	for _, ap := range pb.GetCryptoAdjustAllowance().GetCryptoAllowances() {
		temp := _HbarAllowanceFromProtobuf(ap)
		accountApproval = append(accountApproval, &temp)
	}

	for _, ap := range pb.GetCryptoAdjustAllowance().GetTokenAllowances() {
		temp := _TokenAllowanceFromProtobuf(ap)
		tokenApproval = append(tokenApproval, &temp)
	}

	for _, ap := range pb.GetCryptoAdjustAllowance().GetNftAllowances() {
		temp := _TokenNftAllowanceFromProtobuf(ap)
		nftApproval = append(nftApproval, &temp)
	}

	return &AccountAllowanceAdjustTransaction{
		Transaction:     transaction,
		hbarAllowances:  accountApproval,
		tokenAllowances: tokenApproval,
		nftAllowances:   nftApproval,
	}
}

func (transaction *AccountAllowanceAdjustTransaction) _AdjustHbarAllowance(ownerAccountID *AccountID, id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	transaction.hbarAllowances = append(transaction.hbarAllowances, &HbarAllowance{
		SpenderAccountID: &id,
		OwnerAccountID:   ownerAccountID,
		Amount:           amount.AsTinybar(),
	})

	return transaction
}

// AddHbarAllowance
// Deprecated: Use `GrantHbarAllowance` instead
func (transaction *AccountAllowanceAdjustTransaction) AddHbarAllowance(id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustHbarAllowance(nil, id, amount)
}

func (transaction *AccountAllowanceAdjustTransaction) GrantHbarAllowance(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustHbarAllowance(&ownerAccountID, id, amount)
}

func (transaction *AccountAllowanceAdjustTransaction) RevokeHbarAllowance(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustHbarAllowance(&ownerAccountID, id, amount.Negated())
}

func (transaction *AccountAllowanceAdjustTransaction) GetHbarAllowances() []*HbarAllowance {
	return transaction.hbarAllowances
}

func (transaction *AccountAllowanceAdjustTransaction) _AdjustTokenAllowance(tokenID TokenID, ownerAccountID *AccountID, accountID AccountID, amount int64) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	tokenApproval := TokenAllowance{
		TokenID:          &tokenID,
		SpenderAccountID: &accountID,
		OwnerAccountID:   ownerAccountID,
		Amount:           amount,
	}

	transaction.tokenAllowances = append(transaction.tokenAllowances, &tokenApproval)
	return transaction
}

// AddTokenAllowance
// Deprecated - Use `GrantTokenAllowance()` instead
func (transaction *AccountAllowanceAdjustTransaction) AddTokenAllowance(tokenID TokenID, accountID AccountID, amount int64) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenAllowance(tokenID, nil, accountID, amount)
}

func (transaction *AccountAllowanceAdjustTransaction) GrantTokenAllowance(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount int64) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenAllowance(tokenID, &ownerAccountID, accountID, amount)
}

func (transaction *AccountAllowanceAdjustTransaction) RevokeTokenAllowance(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount uint64) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenAllowance(tokenID, &ownerAccountID, accountID, -int64(amount))
}

func (transaction *AccountAllowanceAdjustTransaction) GetTokenAllowances() []*TokenAllowance {
	return transaction.tokenAllowances
}

func (transaction *AccountAllowanceAdjustTransaction) _AdjustTokenNftAllowance(nftID NftID, ownerAccountID *AccountID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()

	for _, t := range transaction.nftAllowances {
		if t.TokenID.String() == nftID.TokenID.String() {
			if t.SpenderAccountID.String() == accountID.String() {
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
		TokenID:          &nftID.TokenID,
		SpenderAccountID: &accountID,
		OwnerAccountID:   ownerAccountID,
		SerialNumbers:    []int64{nftID.SerialNumber},
		AllSerials:       false,
	})
	return transaction
}

// AddTokenNftAllowance
// Deprecated: Use `GrantTokenNftAllowance()` instead
func (transaction *AccountAllowanceAdjustTransaction) AddTokenNftAllowance(nftID NftID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenNftAllowance(nftID, nil, accountID)
}

func (transaction *AccountAllowanceAdjustTransaction) GrantTokenNftAllowance(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenNftAllowance(nftID, &ownerAccountID, accountID)
}

func (transaction *AccountAllowanceAdjustTransaction) RevokeTokenNftAllowance(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenNftAllowance(nftID, &ownerAccountID, accountID)
}

func (transaction *AccountAllowanceAdjustTransaction) _AdjustTokenNftAllowanceAllSerials(tokenID TokenID, ownerAccountID *AccountID, spenderAccount AccountID, allSerials bool) *AccountAllowanceAdjustTransaction {
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
		OwnerAccountID:   ownerAccountID,
		SerialNumbers:    []int64{},
		AllSerials:       allSerials,
	})
	return transaction
}

// AddAllTokenNftAllowance
// Deprecated: Use `GrantTokenNftAllowanceAllSerials()` instead
func (transaction *AccountAllowanceAdjustTransaction) AddAllTokenNftAllowance(tokenID TokenID, spenderAccount AccountID) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenNftAllowanceAllSerials(tokenID, nil, spenderAccount, true)
}

func (transaction *AccountAllowanceAdjustTransaction) GrantTokenNftAllowanceAllSerials(ownerAccountID AccountID, tokenID TokenID, spenderAccount AccountID) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenNftAllowanceAllSerials(tokenID, &ownerAccountID, spenderAccount, true)
}

func (transaction *AccountAllowanceAdjustTransaction) RevokeTokenNftAllowanceAllSerials(ownerAccountID AccountID, tokenID TokenID, spenderAccount AccountID) *AccountAllowanceAdjustTransaction {
	return transaction._AdjustTokenNftAllowanceAllSerials(tokenID, &ownerAccountID, spenderAccount, false)
}

func (transaction *AccountAllowanceAdjustTransaction) GetTokenNftAllowances() []*TokenNftAllowance {
	return transaction.nftAllowances
}

func (transaction *AccountAllowanceAdjustTransaction) _ValidateNetworkOnIDs(client *Client) error {
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

func (transaction *AccountAllowanceAdjustTransaction) _Build() *services.TransactionBody {
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
		Data: &services.TransactionBody_CryptoAdjustAllowance{
			CryptoAdjustAllowance: &services.CryptoAdjustAllowanceTransactionBody{
				CryptoAllowances: accountApproval,
				NftAllowances:    nftApproval,
				TokenAllowances:  tokenApproval,
			},
		},
	}
}

func (transaction *AccountAllowanceAdjustTransaction) Schedule() (*ScheduleCreateTransaction, error) {
	transaction._RequireNotFrozen()

	scheduled, err := transaction._ConstructScheduleProtobuf()
	if err != nil {
		return nil, err
	}

	return NewScheduleCreateTransaction()._SetSchedulableTransactionBody(scheduled), nil
}

func (transaction *AccountAllowanceAdjustTransaction) _ConstructScheduleProtobuf() (*services.SchedulableTransactionBody, error) {
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
		Data: &services.SchedulableTransactionBody_CryptoAdjustAllowance{
			CryptoAdjustAllowance: &services.CryptoAdjustAllowanceTransactionBody{
				CryptoAllowances: accountApproval,
				NftAllowances:    nftApproval,
				TokenAllowances:  tokenApproval,
			},
		},
	}, nil
}

func _AccountAdjustAllowanceTransactionGetMethod(request _Request, channel *_Channel) _Method {
	return _Method{
		transaction: channel._GetCrypto().AdjustAllowance,
	}
}

func (transaction *AccountAllowanceAdjustTransaction) IsFrozen() bool {
	return transaction._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (transaction *AccountAllowanceAdjustTransaction) Sign(
	privateKey PrivateKey,
) *AccountAllowanceAdjustTransaction {
	return transaction.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (transaction *AccountAllowanceAdjustTransaction) SignWithOperator(
	client *Client,
) (*AccountAllowanceAdjustTransaction, error) {
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
func (transaction *AccountAllowanceAdjustTransaction) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *AccountAllowanceAdjustTransaction {
	if !transaction._KeyAlreadySigned(publicKey) {
		transaction._SignWith(publicKey, signer)
	}

	return transaction
}

// Execute executes the Transaction with the provided client
func (transaction *AccountAllowanceAdjustTransaction) Execute(
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
		_Request{
			transaction: &transaction.Transaction,
		},
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_AccountAdjustAllowanceTransactionGetMethod,
		_TransactionMapStatusError,
		_TransactionMapResponse,
		transaction._GetLogID(),
		transaction.grpcDeadline,
	)

	if err != nil {
		return TransactionResponse{
			TransactionID: transaction.GetTransactionID(),
			NodeID:        resp.transaction.NodeID,
		}, err
	}

	hash, err := transaction.GetTransactionHash()
	if err != nil {
		return TransactionResponse{}, err
	}

	return TransactionResponse{
		TransactionID: transaction.GetTransactionID(),
		NodeID:        resp.transaction.NodeID,
		Hash:          hash,
	}, nil
}

func (transaction *AccountAllowanceAdjustTransaction) Freeze() (*AccountAllowanceAdjustTransaction, error) {
	return transaction.FreezeWith(nil)
}

func (transaction *AccountAllowanceAdjustTransaction) FreezeWith(client *Client) (*AccountAllowanceAdjustTransaction, error) {
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
		return &AccountAllowanceAdjustTransaction{}, err
	}

	return transaction, _TransactionFreezeWith(&transaction.Transaction, client, body)
}

func (transaction *AccountAllowanceAdjustTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this AccountAllowanceAdjustTransaction.
func (transaction *AccountAllowanceAdjustTransaction) SetMaxTransactionFee(fee Hbar) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetMaxTransactionFee(fee)
	return transaction
}

// SetRegenerateTransactionID sets if transaction IDs should be regenerated when `TRANSACTION_EXPIRED` is received
func (transaction *AccountAllowanceAdjustTransaction) SetRegenerateTransactionID(regenerateTransactionID bool) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetRegenerateTransactionID(regenerateTransactionID)
	return transaction
}

// GetRegenerateTransactionID returns true if transaction ID regeneration is enabled.
func (transaction *AccountAllowanceAdjustTransaction) GetRegenerateTransactionID() bool {
	return transaction.Transaction.GetRegenerateTransactionID()
}

func (transaction *AccountAllowanceAdjustTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountAllowanceAdjustTransaction.
func (transaction *AccountAllowanceAdjustTransaction) SetTransactionMemo(memo string) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *AccountAllowanceAdjustTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountAllowanceAdjustTransaction.
func (transaction *AccountAllowanceAdjustTransaction) SetTransactionValidDuration(duration time.Duration) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

func (transaction *AccountAllowanceAdjustTransaction) GetTransactionID() TransactionID {
	return transaction.Transaction.GetTransactionID()
}

// SetTransactionID sets the TransactionID for this AccountAllowanceAdjustTransaction.
func (transaction *AccountAllowanceAdjustTransaction) SetTransactionID(transactionID TransactionID) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()

	transaction.Transaction.SetTransactionID(transactionID)
	return transaction
}

// SetNodeAccountIDs sets the _Node AccountID for this AccountAllowanceAdjustTransaction.
func (transaction *AccountAllowanceAdjustTransaction) SetNodeAccountIDs(nodeID []AccountID) *AccountAllowanceAdjustTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetNodeAccountIDs(nodeID)
	return transaction
}

func (transaction *AccountAllowanceAdjustTransaction) SetMaxRetry(count int) *AccountAllowanceAdjustTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

func (transaction *AccountAllowanceAdjustTransaction) AddSignature(publicKey PublicKey, signature []byte) *AccountAllowanceAdjustTransaction {
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

func (transaction *AccountAllowanceAdjustTransaction) SetMaxBackoff(max time.Duration) *AccountAllowanceAdjustTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *AccountAllowanceAdjustTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *AccountAllowanceAdjustTransaction) SetMinBackoff(min time.Duration) *AccountAllowanceAdjustTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

func (transaction *AccountAllowanceAdjustTransaction) GetMinBackoff() time.Duration {
	if transaction.minBackoff != nil {
		return *transaction.minBackoff
	}

	return 250 * time.Millisecond
}

func (transaction *AccountAllowanceAdjustTransaction) _GetLogID() string {
	timestamp := transaction.transactionIDs._GetCurrent().(TransactionID).ValidStart
	return fmt.Sprintf("AccountAllowanceAdjustTransaction:%d", timestamp.UnixNano())
}
