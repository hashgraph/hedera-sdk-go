package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"time"
)

type ScheduleInfo struct {
	ScheduleID               ScheduleID
	CreatorAccountID         AccountID
	PayerAccountID           AccountID
	Executed                 time.Time
	Deleted                  time.Time
	ExpirationTime           time.Time
	ScheduledTransactionBody *SchedulableTransactionBody
	Signers                  *KeyList
	AdminKey                 Key
	Memo                     string
	ScheduledTransactionID   *TransactionID
}

func scheduleInfoFromProtobuf(pb *proto.ScheduleInfo) ScheduleInfo {
	var adminKey Key
	if pb.AdminKey != nil {
		adminKey, _ = keyFromProtobuf(pb.AdminKey)
	}

	var signers KeyList
	if pb.Signers != nil {
		signers, _ = keyListFromProtobuf(pb.Signers)
	}

	var scheduledTransactionID TransactionID
	if pb.ScheduledTransactionID != nil {
		scheduledTransactionID = transactionIDFromProtobuf(pb.ScheduledTransactionID)
	}

	var executed time.Time
	var deleted time.Time
	switch t := pb.Data.(type) {
	case *proto.ScheduleInfo_ExecutionTime:
		executed = timeFromProtobuf(t.ExecutionTime)
	case *proto.ScheduleInfo_DeletionTime:
		deleted = timeFromProtobuf(t.DeletionTime)
	}

	return ScheduleInfo{
		ScheduleID:               scheduleIDFromProtobuf(pb.ScheduleID),
		CreatorAccountID:         accountIDFromProtobuf(pb.CreatorAccountID),
		PayerAccountID:           accountIDFromProtobuf(pb.PayerAccountID),
		Executed:                 executed,
		Deleted:                  deleted,
		ExpirationTime:           timeFromProtobuf(pb.ExpirationTime),
		ScheduledTransactionBody: schedulableTransactionBodyFromProtobuf(pb.ScheduledTransactionBody),
		Signers:                  &signers,
		AdminKey:                 adminKey,
		Memo:                     pb.Memo,
		ScheduledTransactionID:   &scheduledTransactionID,
	}
}

func (scheduleInfo *ScheduleInfo) toProtobuf() *proto.ScheduleInfo {
	var adminKey *proto.Key
	if scheduleInfo.AdminKey != nil {
		adminKey = scheduleInfo.AdminKey.toProtoKey()
	}

	var signers *proto.KeyList
	if scheduleInfo.Signers != nil {
		signers = scheduleInfo.Signers.toProtoKeyList()
	}

	info := &proto.ScheduleInfo{
		ScheduleID:               scheduleInfo.ScheduleID.toProtobuf(),
		ExpirationTime:           timeToProtobuf(scheduleInfo.ExpirationTime),
		ScheduledTransactionBody: scheduleInfo.ScheduledTransactionBody.toProtobuf(),
		Memo:                     scheduleInfo.Memo,
		AdminKey:                 adminKey,
		Signers:                  signers,
		CreatorAccountID:         scheduleInfo.CreatorAccountID.toProtobuf(),
		PayerAccountID:           scheduleInfo.PayerAccountID.toProtobuf(),
		ScheduledTransactionID:   scheduleInfo.ScheduledTransactionID.toProtobuf(),
	}

	if scheduleInfo.Executed.IsZero() {
		info.Data = &proto.ScheduleInfo_DeletionTime{
			DeletionTime: timeToProtobuf(scheduleInfo.Deleted),
		}
	} else {
		info.Data = &proto.ScheduleInfo_ExecutionTime{
			ExecutionTime: timeToProtobuf(scheduleInfo.Executed),
		}
	}

	return info
}

func (scheduleInfo *ScheduleInfo) GetTransaction() (interface{}, error) {
	pbBody := scheduleInfo.ScheduledTransactionBody.Transaction.pbBody
	tx := *scheduleInfo.ScheduledTransactionBody.Transaction
	switch pbBody.Data.(type) {
	case *proto.TransactionBody_ContractCall:
		return contractExecuteTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_ContractCreateInstance:
		return contractCreateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_ContractUpdateInstance:
		return contractUpdateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_ContractDeleteInstance:
		return contractDeleteTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_CryptoAddLiveHash:
		return liveHashAddTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_CryptoCreateAccount:
		return accountCreateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_CryptoDelete:
		return accountDeleteTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_CryptoDeleteLiveHash:
		return liveHashDeleteTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_CryptoTransfer:
		return transferTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_CryptoUpdateAccount:
		return accountUpdateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_FileAppend:
		return fileAppendTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_FileCreate:
		return fileCreateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_FileDelete:
		return fileDeleteTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_FileUpdate:
		return fileUpdateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_SystemDelete:
		return systemDeleteTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_SystemUndelete:
		return systemUndeleteTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_Freeze:
		return freezeTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_ConsensusCreateTopic:
		return topicCreateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_ConsensusUpdateTopic:
		return topicUpdateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_ConsensusDeleteTopic:
		return topicDeleteTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_ConsensusSubmitMessage:
		return topicMessageSubmitTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_TokenCreation:
		return tokenCreateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_TokenFreeze:
		return tokenFreezeTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_TokenUnfreeze:
		return tokenUnfreezeTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_TokenGrantKyc:
		return tokenGrantKycTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_TokenRevokeKyc:
		return tokenRevokeKycTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_TokenDeletion:
		return tokenDeleteTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_TokenUpdate:
		return tokenUpdateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_TokenMint:
		return tokenMintTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_TokenBurn:
		return tokenBurnTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_TokenWipe:
		return tokenWipeTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_TokenAssociate:
		return tokenAssociateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_TokenDissociate:
		return tokenDissociateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_ScheduleCreate:
		return scheduleCreateTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_ScheduleSign:
		return scheduleSignTransactionFromProtobuf(tx, pbBody), nil
	case *proto.TransactionBody_ScheduleDelete:
		return scheduleDeleteTransactionFromProtobuf(tx, pbBody), nil
	default:
		return Transaction{}, errFailedToDeserializeBytes
	}
}
