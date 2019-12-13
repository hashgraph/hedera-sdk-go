package hedera

import (
	"errors"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

var ErrNoReceipt = errors.New("response was not `TransactionGetReceipt`")

type TransactionReceipt struct {
	inner *proto.TransactionReceipt
}

func TransactionReceiptFromResponse(response proto.Response) (*TransactionReceipt, error) {
	transactionGetReceipt := response.GetTransactionGetReceipt()

	if transactionGetReceipt == nil {
		return nil, ErrNoReceipt
	}

	receipt := TransactionReceipt{
		transactionGetReceipt.Receipt,
	}

	return &receipt, nil
}

func (transactionReceipt TransactionReceipt) AccountID() *AccountID {
	internalID := transactionReceipt.inner.AccountID

	if internalID == nil {
		return nil
	}

	return &AccountID{
		uint64(internalID.ShardNum),
		uint64(internalID.RealmNum),
		uint64(internalID.AccountNum),
	}
}
