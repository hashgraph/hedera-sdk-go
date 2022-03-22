package hedera

import (
	"fmt"
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type AccountAllowanceApproveTransaction struct {
	Transaction
	hbarAllowances  []*HbarAllowance
	tokenAllowances []*TokenAllowance
	nftAllowances   []*TokenNftAllowance
}

func NewAccountAllowanceApproveTransaction() *AccountAllowanceApproveTransaction {
	transaction := AccountAllowanceApproveTransaction{
		Transaction: _NewTransaction(),
	}

	transaction.SetMaxTransactionFee(NewHbar(2))

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

func (transaction *AccountAllowanceApproveTransaction) AddHbarApproval(id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	return transaction._ApproveHbarApproval(nil, id, amount)
}

func (transaction *AccountAllowanceApproveTransaction) ApproveHbarApproval(ownerAccountID AccountID, id AccountID, amount Hbar) *AccountAllowanceApproveTransaction {
	return transaction._ApproveHbarApproval(&ownerAccountID, id, amount)
}

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

func (transaction *AccountAllowanceApproveTransaction) AddTokenApproval(tokenID TokenID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenApproval(tokenID, nil, accountID, amount)
}

func (transaction *AccountAllowanceApproveTransaction) ApproveTokenApproval(tokenID TokenID, ownerAccountID AccountID, accountID AccountID, amount int64) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenApproval(tokenID, &ownerAccountID, accountID, amount)
}

func (transaction *AccountAllowanceApproveTransaction) GetTokenAllowances() []*TokenAllowance {
	return transaction.tokenAllowances
}

func (transaction *AccountAllowanceApproveTransaction) _ApproveTokenNftApproval(nftID NftID, ownerAccountID *AccountID, accountID AccountID) *AccountAllowanceApproveTransaction {
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
		SerialNumbers:    []int64{nftID.SerialNumber},
		AllSerials:       false,
		OwnerAccountID:   ownerAccountID,
	})
	return transaction
}

func (transaction *AccountAllowanceApproveTransaction) AddTokenNftApproval(nftID NftID, accountID AccountID) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenNftApproval(nftID, nil, accountID)
}

func (transaction *AccountAllowanceApproveTransaction) ApproveTokenNftApproval(nftID NftID, ownerAccountID AccountID, accountID AccountID) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenNftApproval(nftID, &ownerAccountID, accountID)
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

func (transaction *AccountAllowanceApproveTransaction) AddAllTokenNftApproval(tokenID TokenID, spenderAccount AccountID) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenNftAllowanceAllSerials(tokenID, nil, spenderAccount)
}

func (transaction *AccountAllowanceApproveTransaction) ApproveTokenNftAllowanceAllSerials(tokenID TokenID, ownerAccountID AccountID, spenderAccount AccountID) *AccountAllowanceApproveTransaction {
	return transaction._ApproveTokenNftAllowanceAllSerials(tokenID, nil, spenderAccount)
}

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

func _AccountApproveAllowanceTransactionGetMethod(request _Request, channel *_Channel) _Method {
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
		_Request{
			transaction: &transaction.Transaction,
		},
		_TransactionShouldRetry,
		_TransactionMakeRequest,
		_TransactionAdvanceRequest,
		_TransactionGetNodeAccountID,
		_AccountApproveAllowanceTransactionGetMethod,
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

func (transaction *AccountAllowanceApproveTransaction) GetMaxTransactionFee() Hbar {
	return transaction.Transaction.GetMaxTransactionFee()
}

// SetMaxTransactionFee sets the max transaction fee for this AccountAllowanceApproveTransaction.
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

func (transaction *AccountAllowanceApproveTransaction) GetTransactionMemo() string {
	return transaction.Transaction.GetTransactionMemo()
}

// SetTransactionMemo sets the memo for this AccountAllowanceApproveTransaction.
func (transaction *AccountAllowanceApproveTransaction) SetTransactionMemo(memo string) *AccountAllowanceApproveTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionMemo(memo)
	return transaction
}

func (transaction *AccountAllowanceApproveTransaction) GetTransactionValidDuration() time.Duration {
	return transaction.Transaction.GetTransactionValidDuration()
}

// SetTransactionValidDuration sets the valid duration for this AccountAllowanceApproveTransaction.
func (transaction *AccountAllowanceApproveTransaction) SetTransactionValidDuration(duration time.Duration) *AccountAllowanceApproveTransaction {
	transaction._RequireNotFrozen()
	transaction.Transaction.SetTransactionValidDuration(duration)
	return transaction
}

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

func (transaction *AccountAllowanceApproveTransaction) SetMaxRetry(count int) *AccountAllowanceApproveTransaction {
	transaction.Transaction.SetMaxRetry(count)
	return transaction
}

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

func (transaction *AccountAllowanceApproveTransaction) SetMaxBackoff(max time.Duration) *AccountAllowanceApproveTransaction {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < transaction.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	transaction.maxBackoff = &max
	return transaction
}

func (transaction *AccountAllowanceApproveTransaction) GetMaxBackoff() time.Duration {
	if transaction.maxBackoff != nil {
		return *transaction.maxBackoff
	}

	return 8 * time.Second
}

func (transaction *AccountAllowanceApproveTransaction) SetMinBackoff(min time.Duration) *AccountAllowanceApproveTransaction {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if transaction.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	transaction.minBackoff = &min
	return transaction
}

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
