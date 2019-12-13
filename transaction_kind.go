package hedera

import (
	"context"
	"errors"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TransactionKind int

const (
	CryptoCreateAccount TransactionKind = iota
	CryptoTransfer
)

func (tk TransactionKind) execute(client Client, tx proto.Transaction) (*proto.TransactionResponse, error) {
	switch tk {
	case CryptoCreateAccount:
		return proto.NewCryptoServiceClient(client.conn).CreateAccount(context.TODO(), &tx)
	case CryptoTransfer:
		return nil, errors.New("not implemented yet")
	default:
		return nil, errors.New("invalid kind provided")
	}
}
