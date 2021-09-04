package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"github.com/pkg/errors"
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

func _ScheduleInfoFromProtobuf(pb *proto.ScheduleInfo) ScheduleInfo {
	if pb == nil {
		return ScheduleInfo{}
	}
	var adminKey Key
	if pb.AdminKey != nil {
		adminKey, _ = _KeyFromProtobuf(pb.AdminKey)
	}

	var signatories KeyList
	if pb.Signers != nil {
		signatories, _ = _KeyListFromProtobuf(pb.Signers)
	}

	var scheduledTransactionID TransactionID
	if pb.ScheduledTransactionID != nil {
		scheduledTransactionID = _TransactionIDFromProtobuf(pb.ScheduledTransactionID)
	}

	var executed *time.Time
	var deleted *time.Time
	switch t := pb.Data.(type) {
	case *proto.ScheduleInfo_ExecutionTime:
		time := _TimeFromProtobuf(t.ExecutionTime)
		executed = &time
	case *proto.ScheduleInfo_DeletionTime:
		time := _TimeFromProtobuf(t.DeletionTime)
		deleted = &time
	}

	creatorAccountID := AccountID{}
	if pb.CreatorAccountID != nil {
		creatorAccountID = *_AccountIDFromProtobuf(pb.CreatorAccountID)
	}

	payerAccountID := AccountID{}
	if pb.PayerAccountID != nil {
		payerAccountID = *_AccountIDFromProtobuf(pb.PayerAccountID)
	}

	scheduleID := ScheduleID{}
	if pb.ScheduleID != nil {
		scheduleID = *_ScheduleIDFromProtobuf(pb.ScheduleID)
	}

	return ScheduleInfo{
		ScheduleID:               scheduleID,
		CreatorAccountID:         creatorAccountID,
		PayerAccountID:           payerAccountID,
		ExecutedAt:               executed,
		DeletedAt:                deleted,
		ExpirationTime:           _TimeFromProtobuf(pb.ExpirationTime),
		Signatories:              &signatories,
		Signers:                  &signatories,
		AdminKey:                 adminKey,
		Memo:                     pb.Memo,
		ScheduledTransactionID:   &scheduledTransactionID,
		scheduledTransactionBody: pb.ScheduledTransactionBody,
	}
}

func (scheduleInfo *ScheduleInfo) _ToProtobuf() *proto.ScheduleInfo { // nolint
	var adminKey *proto.Key
	if scheduleInfo.AdminKey != nil {
		adminKey = scheduleInfo.AdminKey._ToProtoKey()
	}

	var signatories *proto.KeyList
	if scheduleInfo.Signatories != nil {
		signatories = scheduleInfo.Signatories._ToProtoKeyList()
	} else if scheduleInfo.Signers != nil {
		signatories = scheduleInfo.Signers._ToProtoKeyList()
	}

	info := &proto.ScheduleInfo{
		ScheduleID:               scheduleInfo.ScheduleID._ToProtobuf(),
		ExpirationTime:           _TimeToProtobuf(scheduleInfo.ExpirationTime),
		ScheduledTransactionBody: scheduleInfo.scheduledTransactionBody,
		Memo:                     scheduleInfo.Memo,
		AdminKey:                 adminKey,
		Signers:                  signatories,
		CreatorAccountID:         scheduleInfo.CreatorAccountID._ToProtobuf(),
		PayerAccountID:           scheduleInfo.PayerAccountID._ToProtobuf(),
		ScheduledTransactionID:   scheduleInfo.ScheduledTransactionID._ToProtobuf(),
	}

	if scheduleInfo.ExecutedAt != nil {
		info.Data = &proto.ScheduleInfo_DeletionTime{
			DeletionTime: _TimeToProtobuf(*scheduleInfo.DeletedAt),
		}
	} else if scheduleInfo.DeletedAt != nil {
		info.Data = &proto.ScheduleInfo_ExecutionTime{
			ExecutionTime: _TimeToProtobuf(*scheduleInfo.ExecutedAt),
		}
	}

	return info
}

func (scheduleInfo *ScheduleInfo) GetScheduledTransaction() (ITransaction, error) { // nolint
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

		tx2 := _ContractExecuteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ContractCreateInstance:
		pbBody.Data = &proto.TransactionBody_ContractCreateInstance{
			ContractCreateInstance: pb.GetContractCreateInstance(),
		}

		tx2 := _ContractCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ContractUpdateInstance:
		pbBody.Data = &proto.TransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: pb.GetContractUpdateInstance(),
		}

		tx2 := _ContractUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ContractDeleteInstance:
		pbBody.Data = &proto.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: pb.GetContractDeleteInstance(),
		}

		tx2 := _ContractDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_CryptoCreateAccount:
		pbBody.Data = &proto.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: pb.GetCryptoCreateAccount(),
		}

		tx2 := _AccountCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_CryptoDelete:
		pbBody.Data = &proto.TransactionBody_CryptoDelete{
			CryptoDelete: pb.GetCryptoDelete(),
		}

		tx2 := _AccountDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_CryptoTransfer:
		pbBody.Data = &proto.TransactionBody_CryptoTransfer{
			CryptoTransfer: pb.GetCryptoTransfer(),
		}

		tx2 := _TransferTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_CryptoUpdateAccount:
		pbBody.Data = &proto.TransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: pb.GetCryptoUpdateAccount(),
		}

		tx2 := _AccountUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_FileAppend:
		pbBody.Data = &proto.TransactionBody_FileAppend{
			FileAppend: pb.GetFileAppend(),
		}

		tx2 := _FileAppendTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_FileCreate:
		pbBody.Data = &proto.TransactionBody_FileCreate{
			FileCreate: pb.GetFileCreate(),
		}

		tx2 := _FileCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_FileDelete:
		pbBody.Data = &proto.TransactionBody_FileDelete{
			FileDelete: pb.GetFileDelete(),
		}

		tx2 := _FileDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_FileUpdate:
		pbBody.Data = &proto.TransactionBody_FileUpdate{
			FileUpdate: pb.GetFileUpdate(),
		}

		tx2 := _FileUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_SystemDelete:
		pbBody.Data = &proto.TransactionBody_SystemDelete{
			SystemDelete: pb.GetSystemDelete(),
		}

		tx2 := _SystemDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_SystemUndelete:
		pbBody.Data = &proto.TransactionBody_SystemUndelete{
			SystemUndelete: pb.GetSystemUndelete(),
		}

		tx2 := _SystemUndeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_Freeze:
		pbBody.Data = &proto.TransactionBody_Freeze{
			Freeze: pb.GetFreeze(),
		}

		tx2 := _FreezeTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ConsensusCreateTopic:
		pbBody.Data = &proto.TransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: pb.GetConsensusCreateTopic(),
		}

		tx2 := _TopicCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ConsensusUpdateTopic:
		pbBody.Data = &proto.TransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: pb.GetConsensusUpdateTopic(),
		}

		tx2 := _TopicUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ConsensusDeleteTopic:
		pbBody.Data = &proto.TransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: pb.GetConsensusDeleteTopic(),
		}

		tx2 := _TopicDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ConsensusSubmitMessage:
		pbBody.Data = &proto.TransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: pb.GetConsensusSubmitMessage(),
		}

		tx2 := _TopicMessageSubmitTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenCreation:
		pbBody.Data = &proto.TransactionBody_TokenCreation{
			TokenCreation: pb.GetTokenCreation(),
		}

		tx2 := _TokenCreateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenFreeze:
		pbBody.Data = &proto.TransactionBody_TokenFreeze{
			TokenFreeze: pb.GetTokenFreeze(),
		}

		tx2 := _TokenFreezeTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenUnfreeze:
		pbBody.Data = &proto.TransactionBody_TokenUnfreeze{
			TokenUnfreeze: pb.GetTokenUnfreeze(),
		}

		tx2 := _TokenUnfreezeTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenGrantKyc:
		pbBody.Data = &proto.TransactionBody_TokenGrantKyc{
			TokenGrantKyc: pb.GetTokenGrantKyc(),
		}

		tx2 := _TokenGrantKycTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenRevokeKyc:
		pbBody.Data = &proto.TransactionBody_TokenRevokeKyc{
			TokenRevokeKyc: pb.GetTokenRevokeKyc(),
		}

		tx2 := _TokenRevokeKycTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenDeletion:
		pbBody.Data = &proto.TransactionBody_TokenDeletion{
			TokenDeletion: pb.GetTokenDeletion(),
		}

		tx2 := _TokenDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenUpdate:
		pbBody.Data = &proto.TransactionBody_TokenUpdate{
			TokenUpdate: pb.GetTokenUpdate(),
		}

		tx2 := _TokenUpdateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenMint:
		pbBody.Data = &proto.TransactionBody_TokenMint{
			TokenMint: pb.GetTokenMint(),
		}

		tx2 := _TokenMintTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenBurn:
		pbBody.Data = &proto.TransactionBody_TokenBurn{
			TokenBurn: pb.GetTokenBurn(),
		}

		tx2 := _TokenBurnTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenWipe:
		pbBody.Data = &proto.TransactionBody_TokenWipe{
			TokenWipe: pb.GetTokenWipe(),
		}

		tx2 := _TokenWipeTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenAssociate:
		pbBody.Data = &proto.TransactionBody_TokenAssociate{
			TokenAssociate: pb.GetTokenAssociate(),
		}

		tx2 := _TokenAssociateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_TokenDissociate:
		pbBody.Data = &proto.TransactionBody_TokenDissociate{
			TokenDissociate: pb.GetTokenDissociate(),
		}

		tx2 := _TokenDissociateTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	case *proto.SchedulableTransactionBody_ScheduleDelete:
		pbBody.Data = &proto.TransactionBody_ScheduleDelete{
			ScheduleDelete: pb.GetScheduleDelete(),
		}

		tx2 := _ScheduleDeleteTransactionFromProtobuf(tx, pbBody)
		return &tx2, nil
	default:
		return nil, errors.New("(BUG) non-exhaustive switch statement")
	}
}
