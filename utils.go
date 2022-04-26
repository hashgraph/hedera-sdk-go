package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
)

type AccountCreateOption func() (func(*AccountCreateTransaction), error)

func WithMemo(s string) AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		return func(a *AccountCreateTransaction) {
			a.memo = s
		}, nil
	}
}

func WithInitBalance(u uint64) AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		return func(a *AccountCreateTransaction) {
			a._RequireNotFrozen()
			a.initialBalance = u
		}, nil
	}
}

func WithReceiveerSigRequired() AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		return func(a *AccountCreateTransaction) {
			a.receiverSignatureRequired = true
		}, nil
	}
}

func WithProxyAccountIDStr(s string) AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		if accountID, err := AccountIDFromString(s); err != nil {
			return nil, err
		} else {
			return func(a *AccountCreateTransaction) {
				a.proxyAccountID = &accountID
			}, nil
		}
	}
}

func WithMaxAutoTokenAssociations(u uint32) AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		return func(a *AccountCreateTransaction) {
			a._RequireNotFrozen()
			a.maxAutomaticTokenAssociations = u
		}, nil
	}
}

func WithAutoRenewPeriod(t time.Duration) AccountCreateOption {
	return func() (func(a *AccountCreateTransaction), error) {
		return func(a *AccountCreateTransaction) {
			a._RequireNotFrozen()
			a.autoRenewPeriod = &t
		}, nil
	}
}

type TransferOption func() (func(*TransferTransaction), error)

func BuildAccountCreateTransactionBody(keyByte []byte, opts ...AccountCreateOption) (*services.TransactionBody, error) {
	key, err := PublicKeyFromBytesEd25519(keyByte)
	if err != nil {
		return &services.TransactionBody{}, err
	}
	tx := NewAccountCreateTransaction().SetKey(key)

	for _, opt := range opts {
		f, err2 := opt()
		if err2 != nil {
			return &services.TransactionBody{}, err2
		}
		f(tx)
	}

	return tx._Build(), nil
}

func BuildTransferHbarTransactionBody(
	SenderActIDStr,
	ReceiverActIDStr string,
	amount float64,
	opts ...TransferOption) (*services.TransactionBody, error) {
	SenderActID, err := AccountIDFromString(SenderActIDStr)
	if err != nil {
		return &services.TransactionBody{}, err
	}

	ReceiverActID, err2 := AccountIDFromString(ReceiverActIDStr)
	if err2 != nil {
		return &services.TransactionBody{}, err2
	}

	tx := NewTransferTransaction().
		AddHbarTransfer(SenderActID, NewHbar(amount*-1)).
		AddHbarTransfer(ReceiverActID, NewHbar(amount))

	for _, opt := range opts {
		f, err2 := opt()
		if err2 != nil {
			return &services.TransactionBody{}, err2
		}
		f(tx)
	}

	return tx._Build(), nil
}
