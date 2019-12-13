package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TransactionReceipt struct {
	AccountID *AccountID
}

func transactionReceiptFromResponse(response *proto.Response) TransactionReceipt {
	pb := response.GetTransactionGetReceipt()

	var accountID *AccountID
	if pb.Receipt.AccountID != nil {
		accountIDValue := accountIDFromProto(pb.Receipt.AccountID)
		accountID = &accountIDValue
	}

	return TransactionReceipt{
		AccountID: accountID,
	}
}

//func (transactionReceipt TransactionReceipt) AccountID() *AccountID {
//	internalID := transactionReceipt.inner.AccountID
//
//	if internalID == nil {
//		return nil
//	}
//
//	return &AccountID{
//		uint64(internalID.ShardNum),
//		uint64(internalID.RealmNum),
//		uint64(internalID.AccountNum),
//	}
//}
