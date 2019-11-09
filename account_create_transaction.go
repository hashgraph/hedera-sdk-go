package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/hedera_proto"
)

type AccountCreateTransaction struct {
	builder        TransactionBuilder
	PublicKey      *Ed25519PublicKey
	InitialBalance uint64
}

func NewAccountCreateTransaction(client *Client) AccountCreateTransaction {

	builder := TransactionBuilder{
		client: client,
		kind:   CryptoCreateAccount,
		body:   hedera_proto.TransactionBody{},
	}

	return AccountCreateTransaction{
		builder: builder,
	}
}

func (tx AccountCreateTransaction) SetMaxTransactionFee(fee uint64) AccountCreateTransaction {
	tx.builder.SetMaxTransactionFee(fee)

	return tx
}

func (tx AccountCreateTransaction) SetKey(publicKey Ed25519PublicKey) AccountCreateTransaction {

	tx.PublicKey = &publicKey

	return tx
}

func (tx AccountCreateTransaction) SetInitialBalance(balance uint64) AccountCreateTransaction {
	tx.InitialBalance = balance

	return tx
}

func (tx AccountCreateTransaction) Validate() error {
	if tx.PublicKey == nil {
		return fmt.Errorf("AccountCreateTransaction requires Public Key to be set")
	}

	return nil
}

func (tx AccountCreateTransaction) Build() (*Transaction, error) {

	if err := tx.Validate(); err != nil {
		return nil, err
	}

	protoKey := tx.PublicKey.toProtoKey()

	tx.builder.body.Data = &hedera_proto.TransactionBody_CryptoCreateAccount{
		CryptoCreateAccount: &hedera_proto.CryptoCreateTransactionBody{
			Key:            &protoKey,
			InitialBalance: tx.InitialBalance,
		},
	}

	return tx.builder.build(CryptoCreateAccount)
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
