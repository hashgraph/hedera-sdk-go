package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type AccountCreateTransaction struct {
	TransactionBuilder
	PublicKey           *Ed25519PublicKey
	InitialBalance      uint64
	ProxyAccountId      *AccountID
	ReceiverSigRequired bool
}

func NewAccountCreateTransaction(client *Client) AccountCreateTransaction {

	builder := TransactionBuilder{
		client: client,
		kind:   CryptoCreateAccount,
		body:   proto.TransactionBody{},
	}

	return AccountCreateTransaction{
		TransactionBuilder: builder,
	}
}

func (tx AccountCreateTransaction) SetKey(publicKey Ed25519PublicKey) AccountCreateTransaction {

	tx.PublicKey = &publicKey

	return tx
}

func (tx AccountCreateTransaction) SetInitialBalance(balance uint64) AccountCreateTransaction {
	tx.InitialBalance = balance

	return tx
}

func (tx AccountCreateTransaction) SetProxyAccountID(accountID AccountID) AccountCreateTransaction {
	tx.ProxyAccountId = &accountID

	return tx
}

func (tx AccountCreateTransaction) SetReceiverSignatureRequired(required bool) AccountCreateTransaction {
	tx.ReceiverSigRequired = required

	return tx
}

func (tx AccountCreateTransaction) Validate() error {
	if tx.PublicKey == nil {
		return fmt.Errorf("AccountCreateTransaction requires Public Key to be set")
	}

	return nil
}

func (tx AccountCreateTransaction) SetMemo(memo string) AccountCreateTransaction {
	tx.body.Memo = memo

	return tx
}

func (tx AccountCreateTransaction) SetMaxTransactionFee(fee uint64) AccountCreateTransaction {
	tx.MaxTransactionFee = fee

	return tx
}

func (tx AccountCreateTransaction) SetTransactionID(txID TransactionID) AccountCreateTransaction {
	tx.body.TransactionID = txID.proto()

	return tx
}

func (tx AccountCreateTransaction) SetTransactionValidDuration(seconds uint64) AccountCreateTransaction {
	tx.body.TransactionValidDuration = &proto.Duration{Seconds: int64(seconds)}

	return tx
}

func (tx AccountCreateTransaction) SetNodeAccountID(accountID AccountID) AccountCreateTransaction {
	tx.body.NodeAccountID = accountID.proto()

	return tx
}

func (tx AccountCreateTransaction) Build() (*Transaction, error) {

	if err := tx.Validate(); err != nil {
		return nil, err
	}

	protoKey := tx.PublicKey.toProtoKey()

	bodyData := proto.TransactionBody_CryptoCreateAccount{
		CryptoCreateAccount: &proto.CryptoCreateTransactionBody{
			Key:                 &protoKey,
			InitialBalance:      tx.InitialBalance,
			ReceiverSigRequired: tx.ReceiverSigRequired,
		},
	}

	if tx.ProxyAccountId != nil {
		bodyData.CryptoCreateAccount.ProxyAccountID = tx.ProxyAccountId.proto()
	}

	tx.body.Data = &bodyData

	return tx.build()
}

func (tx AccountCreateTransaction) Execute() (*TransactionID, error) {
	transaction, err := tx.Build()

	if err != nil {
		return nil, err
	}

	return transaction.Execute()
}

func (tx AccountCreateTransaction) ExecuteForReceipt() (*TransactionReceipt, error) {
	transaction, err := tx.Build()

	if err != nil {
		return nil, err
	}

	return transaction.ExecuteForReceipt()
}
