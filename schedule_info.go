package hedera

import (
	"errors"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ScheduleInfo struct {
	ScheduleID       ScheduleID
	CreatorAccountID AccountID
	PayerAccountID   AccountID
	TransactionBody  []byte
	Signatories      []PublicKey
	AdminKey         PublicKey
}

func scheduleInfoFromProtobuf(pb *proto.ScheduleInfo) ScheduleInfo {
	var adminKey PublicKey
	if pb.AdminKey != nil {
		adminKey, _ = publicKeyFromProto(pb.AdminKey)
	}

	var signatories []PublicKey
	if pb.Signatories != nil {
		signatories, _ = publicKeyListFromProto(pb.Signatories)
	}

	return ScheduleInfo{
		ScheduleID:       scheduleIDFromProto(pb.ScheduleID),
		CreatorAccountID: accountIDFromProto(pb.CreatorAccountID),
		PayerAccountID:   accountIDFromProto(pb.PayerAccountID),
		TransactionBody:  pb.TransactionBody,
		Signatories:      signatories,
		AdminKey:         adminKey,
	}
}

func (scheduleInfo *ScheduleInfo) toProtobuf() *proto.ScheduleInfo {
	var adminKey *proto.Key
	if scheduleInfo.AdminKey != nil {
		adminKey = scheduleInfo.AdminKey.toProto()
	}

	var temp KeyList
	if scheduleInfo.Signatories != nil {
		temp.AddAll(scheduleInfo.Signatories)
	}

	var signatories *proto.KeyList
	if temp.keys != nil {
		signatories = &proto.KeyList{Keys: temp.keys}
	}

	return &proto.ScheduleInfo{
		ScheduleID:       scheduleInfo.ScheduleID.toProto(),
		CreatorAccountID: scheduleInfo.CreatorAccountID.toProto(),
		PayerAccountID:   scheduleInfo.PayerAccountID.toProto(),
		TransactionBody:  scheduleInfo.TransactionBody,
		Signatories:      signatories,
		AdminKey:         adminKey,
	}
}

func (scheduleInfo *ScheduleInfo) GetTransaction() (interface{}, error) {
	tx := Transaction{}

	var txBody proto.TransactionBody
	err := protobuf.Unmarshal(scheduleInfo.TransactionBody, &txBody)
	if err != nil {
		return &tx, err
	}

	switch txBody.Data.(type) {
	case *proto.TransactionBody_ContractCall:
		tx := NewContractExecuteTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_ContractCreateInstance:
		tx := NewContractCreateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_ContractUpdateInstance:
		tx := NewContractUpdateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_ContractDeleteInstance:
		tx := NewContractDeleteTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_CryptoCreateAccount:
		tx := NewAccountCreateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_CryptoDelete:
		tx := NewAccountDeleteTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_CryptoTransfer:
		tx := NewTransferTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_CryptoUpdateAccount:
		tx := NewAccountUpdateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_FileAppend:
		tx := NewFileAppendTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_FileCreate:
		tx := NewFileCreateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_FileDelete:
		tx := NewFileDeleteTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_FileUpdate:
		tx := NewFileUpdateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_SystemDelete:
		tx := NewSystemDeleteTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_SystemUndelete:
		tx := NewSystemUndeleteTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_Freeze:
		tx := NewFreezeTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_ConsensusCreateTopic:
		tx := NewConsensusTopicCreateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_ConsensusUpdateTopic:
		tx := NewConsensusTopicUpdateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_ConsensusDeleteTopic:
		tx := NewConsensusTopicDeleteTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_ConsensusSubmitMessage:
		tx := NewConsensusMessageSubmitTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_TokenCreation:
		tx := NewTokenCreateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_TokenFreeze:
		tx := NewTokenFreezeTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_TokenUnfreeze:
		tx := NewTokenUnfreezeTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_TokenGrantKyc:
		tx := NewTokenGrantKycTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_TokenRevokeKyc:
		tx := NewTokenGrantKycTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_TokenDeletion:
		tx := NewTokenDeleteTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_TokenUpdate:
		tx := NewTokenUpdateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_TokenMint:
		tx := NewTokenMintTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_TokenBurn:
		tx := NewTokenBurnTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_TokenWipe:
		tx := NewTokenWipeTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_TokenAssociate:
		tx := NewTokenAssociateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_TokenDissociate:
		tx := NewTokenDissociateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_ScheduleCreate:
		tx := NewScheduleCreateTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_ScheduleSign:
		tx := NewScheduleSignTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	case *proto.TransactionBody_ScheduleDelete:
		tx := NewScheduleDeleteTransaction()
		tx.TransactionBuilder.pb = &txBody
		return tx, nil
	default:
		return Transaction{}, errors.New(" unrecognizable transaction")
	}
}
