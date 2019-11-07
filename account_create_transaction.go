package hedera

import (
	"fmt"

	"github.com/hashgraph/hedera-sdk-go/hedera_proto"
)

type AccountCreateTransaction struct {
	TransactionBuilder
	client            *Client
	body              hedera_proto.CryptoCreateTransactionBody
	maxTransactionFee uint64
}

func NewAccountCreateTransaction(client *Client) AccountCreateTransaction {
	return AccountCreateTransaction{
		client: client,
		body:   hedera_proto.CryptoCreateTransactionBody{},
	}
}

func (transaction AccountCreateTransaction) SetMaxTransactionFee(tinybar uint64) AccountCreateTransaction {
	transaction.maxTransactionFee = tinybar

	return transaction
}

func (transaction AccountCreateTransaction) SetKey(publicKey Ed25519PublicKey) AccountCreateTransaction {
	// fixme: use our own built in function for this
	protoKey := hedera_proto.Key_Ed25519{Ed25519: publicKey.keyData}

	transaction.body.Key = &hedera_proto.Key{Key: &protoKey}

	return transaction
}

func (transaction AccountCreateTransaction) SetInitialBalance(tinybar uint64) AccountCreateTransaction {
	transaction.body.InitialBalance = tinybar

	return transaction
}

func (transaction AccountCreateTransaction) validate() error {
	if transaction.body.Key == nil {
		return fmt.Errorf("AccountCreateTransaction requires .setKey")
	}

	return nil
}

func (transaction AccountCreateTransaction) Build() (*Transaction, error) {
	if err := transaction.validate(); err != nil {
		return nil, err
	}

	// not implemented yet
	return nil, nil
}
