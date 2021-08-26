package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"github.com/pkg/errors"
	"time"
)

type ScheduleInfo struct {
	ScheduleID       ScheduleID
	CreatorAccountID AccountID
	PayerAccountID   AccountID
	ExecutedAt       *time.Time
	DeletedAt        *time.Time
	ExpirationTime   time.Time
	Signatories      *KeyList
	// Deprecated: Use ScheduleInfo.Signatories instead
	Signers                  *KeyList
	AdminKey                 Key
	Memo                     string
	ScheduledTransactionID   *TransactionID
	scheduledTransactionBody *proto.SchedulableTransactionBody
}

func scheduleInfoFromProtobuf(pb *proto.ScheduleInfo) ScheduleInfo {
	if pb == nil {
		return ScheduleInfo{}
	}
	var adminKey Key
	if pb.AdminKey != nil {
		adminKey, _ = keyFromProtobuf(pb.AdminKey)
	}

	var signatories KeyList
	if pb.Signers != nil {
		signatories, _ = keyListFromProtobuf(pb.Signers)
	}

	var scheduledTransactionID TransactionID
	if pb.ScheduledTransactionID != nil {
		scheduledTransactionID = transactionIDFromProtobuf(pb.ScheduledTransactionID)
	}

	var executed *time.Time
	var deleted *time.Time
	switch t := pb.Data.(type) {
	case *proto.ScheduleInfo_ExecutionTime:
		time := timeFromProtobuf(t.ExecutionTime)
		executed = &time
	case *proto.ScheduleInfo_DeletionTime:
		time := timeFromProtobuf(t.DeletionTime)
		deleted = &time
	}

	return ScheduleInfo{
		ScheduleID:               scheduleIDFromProtobuf(pb.ScheduleID),
		CreatorAccountID:         accountIDFromProtobuf(pb.CreatorAccountID),
		PayerAccountID:           accountIDFromProtobuf(pb.PayerAccountID),
		ExecutedAt:               executed,
		DeletedAt:                deleted,
		ExpirationTime:           timeFromProtobuf(pb.ExpirationTime),
		Signatories:              &signatories,
		Signers:                  &signatories,
		AdminKey:                 adminKey,
		Memo:                     pb.Memo,
		ScheduledTransactionID:   &scheduledTransactionID,
		scheduledTransactionBody: pb.ScheduledTransactionBody,
	}
}

func (scheduleInfo *ScheduleInfo) toProtobuf() *proto.ScheduleInfo {
	var adminKey *proto.Key
	if scheduleInfo.AdminKey != nil {
		adminKey = scheduleInfo.AdminKey.toProtoKey()
	}

	var signatories *proto.KeyList
	if scheduleInfo.Signatories != nil {
		signatories = scheduleInfo.Signatories.toProtoKeyList()
	} else if scheduleInfo.Signers != nil {
		signatories = scheduleInfo.Signers.toProtoKeyList()
	}

	info := &proto.ScheduleInfo{
		ScheduleID:               scheduleInfo.ScheduleID.toProtobuf(),
		ExpirationTime:           timeToProtobuf(scheduleInfo.ExpirationTime),
		ScheduledTransactionBody: scheduleInfo.scheduledTransactionBody,
		Memo:                     scheduleInfo.Memo,
		AdminKey:                 adminKey,
		Signers:                  signatories,
		CreatorAccountID:         scheduleInfo.CreatorAccountID.toProtobuf(),
		PayerAccountID:           scheduleInfo.PayerAccountID.toProtobuf(),
		ScheduledTransactionID:   scheduleInfo.ScheduledTransactionID.toProtobuf(),
	}

	if scheduleInfo.ExecutedAt != nil {
		info.Data = &proto.ScheduleInfo_DeletionTime{
			DeletionTime: timeToProtobuf(*scheduleInfo.DeletedAt),
		}
	} else if scheduleInfo.DeletedAt != nil {
		info.Data = &proto.ScheduleInfo_ExecutionTime{
			ExecutionTime: timeToProtobuf(*scheduleInfo.ExecutedAt),
		}
	}

	return info
}

func (scheduleInfo *ScheduleInfo) GetScheduledTransaction() (ITransaction, error) {
	pb := scheduleInfo.scheduledTransactionBody

	pbBody := &proto.TransactionBody{
		TransactionFee: pb.TransactionFee,
		Memo:           pb.Memo,
	}

	tx := Transaction{
		transactionFee: pb.GetTransactionFee(),
		memo:           pb.GetMemo(),
	}

	switch pb.Data.(type) {
	case *proto.SchedulableTransactionBody_ContractCall:
		pbBody.Data = &proto.TransactionBody_ContractCall{
			ContractCall: pb.GetContractCall(),
		}

		tx2 := contractExecuteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ContractCreateInstance:
		pbBody.Data = &proto.TransactionBody_ContractCreateInstance{
			ContractCreateInstance: pb.GetContractCreateInstance(),
		}

		tx2 := contractCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ContractUpdateInstance:
		pbBody.Data = &proto.TransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: pb.GetContractUpdateInstance(),
		}

		tx2 := contractUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ContractDeleteInstance:
		pbBody.Data = &proto.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: pb.GetContractDeleteInstance(),
		}

		tx2 := contractDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_CryptoCreateAccount:
		pbBody.Data = &proto.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: pb.GetCryptoCreateAccount(),
		}

		tx2 := accountCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_CryptoDelete:
		pbBody.Data = &proto.TransactionBody_CryptoDelete{
			CryptoDelete: pb.GetCryptoDelete(),
		}

		tx2 := accountDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_CryptoTransfer:
		pbBody.Data = &proto.TransactionBody_CryptoTransfer{
			CryptoTransfer: pb.GetCryptoTransfer(),
		}

		tx2 := transferTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_CryptoUpdateAccount:
		pbBody.Data = &proto.TransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: pb.GetCryptoUpdateAccount(),
		}

		tx2 := accountUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_FileAppend:
		pbBody.Data = &proto.TransactionBody_FileAppend{
			FileAppend: pb.GetFileAppend(),
		}

		tx2 := fileAppendTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_FileCreate:
		pbBody.Data = &proto.TransactionBody_FileCreate{
			FileCreate: pb.GetFileCreate(),
		}

		tx2 := fileCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_FileDelete:
		pbBody.Data = &proto.TransactionBody_FileDelete{
			FileDelete: pb.GetFileDelete(),
		}

		tx2 := fileDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_FileUpdate:
		pbBody.Data = &proto.TransactionBody_FileUpdate{
			FileUpdate: pb.GetFileUpdate(),
		}

		tx2 := fileUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_SystemDelete:
		pbBody.Data = &proto.TransactionBody_SystemDelete{
			SystemDelete: pb.GetSystemDelete(),
		}

		tx2 := systemDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_SystemUndelete:
		pbBody.Data = &proto.TransactionBody_SystemUndelete{
			SystemUndelete: pb.GetSystemUndelete(),
		}

		tx2 := systemUndeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_Freeze:
		pbBody.Data = &proto.TransactionBody_Freeze{
			Freeze: pb.GetFreeze(),
		}

		tx2 := freezeTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ConsensusCreateTopic:
		pbBody.Data = &proto.TransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: pb.GetConsensusCreateTopic(),
		}

		tx2 := topicCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ConsensusUpdateTopic:
		pbBody.Data = &proto.TransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: pb.GetConsensusUpdateTopic(),
		}

		tx2 := topicUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ConsensusDeleteTopic:
		pbBody.Data = &proto.TransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: pb.GetConsensusDeleteTopic(),
		}

		tx2 := topicDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ConsensusSubmitMessage:
		pbBody.Data = &proto.TransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: pb.GetConsensusSubmitMessage(),
		}

		tx2 := topicMessageSubmitTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenCreation:
		pbBody.Data = &proto.TransactionBody_TokenCreation{
			TokenCreation: pb.GetTokenCreation(),
		}

		tx2 := tokenCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenFreeze:
		pbBody.Data = &proto.TransactionBody_TokenFreeze{
			TokenFreeze: pb.GetTokenFreeze(),
		}

		tx2 := tokenFreezeTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenUnfreeze:
		pbBody.Data = &proto.TransactionBody_TokenUnfreeze{
			TokenUnfreeze: pb.GetTokenUnfreeze(),
		}

		tx2 := tokenUnfreezeTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenGrantKyc:
		pbBody.Data = &proto.TransactionBody_TokenGrantKyc{
			TokenGrantKyc: pb.GetTokenGrantKyc(),
		}

		tx2 := tokenGrantKycTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenRevokeKyc:
		pbBody.Data = &proto.TransactionBody_TokenRevokeKyc{
			TokenRevokeKyc: pb.GetTokenRevokeKyc(),
		}

		tx2 := tokenRevokeKycTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenDeletion:
		pbBody.Data = &proto.TransactionBody_TokenDeletion{
			TokenDeletion: pb.GetTokenDeletion(),
		}

		tx2 := tokenDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenUpdate:
		pbBody.Data = &proto.TransactionBody_TokenUpdate{
			TokenUpdate: pb.GetTokenUpdate(),
		}

		tx2 := tokenUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenMint:
		pbBody.Data = &proto.TransactionBody_TokenMint{
			TokenMint: pb.GetTokenMint(),
		}

		tx2 := tokenMintTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenBurn:
		pbBody.Data = &proto.TransactionBody_TokenBurn{
			TokenBurn: pb.GetTokenBurn(),
		}

		tx2 := tokenBurnTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenWipe:
		pbBody.Data = &proto.TransactionBody_TokenWipe{
			TokenWipe: pb.GetTokenWipe(),
		}

		tx2 := tokenWipeTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenAssociate:
		pbBody.Data = &proto.TransactionBody_TokenAssociate{
			TokenAssociate: pb.GetTokenAssociate(),
		}

		tx2 := tokenAssociateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenDissociate:
		pbBody.Data = &proto.TransactionBody_TokenDissociate{
			TokenDissociate: pb.GetTokenDissociate(),
		}

		tx2 := tokenDissociateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ScheduleDelete:
		pbBody.Data = &proto.TransactionBody_ScheduleDelete{
			ScheduleDelete: pb.GetScheduleDelete(),
		}

		tx2 := scheduleDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	default:
		return nil, errors.New("(BUG) non-exhaustive switch statement")
	}
}
