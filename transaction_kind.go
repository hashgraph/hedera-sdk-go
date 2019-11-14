package hedera

import (
	"context"
	"errors"

	"github.com/hashgraph/hedera-sdk-go/hedera_proto"
)

type TransactionKind int

const (
	CryptoCreateAccount TransactionKind = iota
	CryptoTransfer
)

func (tk TransactionKind) execute(client Client, tx hedera_proto.Transaction) (*hedera_proto.TransactionResponse, error) {
	switch tk {
	case CryptoCreateAccount:
		return hedera_proto.NewCryptoServiceClient(client.conn).CreateAccount(context.TODO(), &tx)
	case CryptoTransfer:
		return nil, errors.New("not implemented yet")
	default:
		return nil, errors.New("invalid kind provided")
	}
}
