package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type SchedulableTransactionBody struct {
	TransactionFee uint64
	Memo           string
	Transaction    *Transaction
}

func schedulableTransactionBodyFromProtobuf(pb *proto.SchedulableTransactionBody) *SchedulableTransactionBody {

	pbBody := proto.TransactionBody{
		TransactionID:            nil,
		NodeAccountID:            nil,
		TransactionFee:           pb.TransactionFee,
		TransactionValidDuration: nil,
		GenerateRecord:           false,
		Memo:                     pb.Memo,
		Data:                     nil,
	}

	switch pb.Data.(type) {
	case *proto.SchedulableTransactionBody_ContractCall:
		pbBody.Data = &proto.TransactionBody_ContractCall{
			ContractCall: pb.GetContractCall(),
		}
	case *proto.SchedulableTransactionBody_ContractCreateInstance:
		pbBody.Data = &proto.TransactionBody_ContractCreateInstance{
			ContractCreateInstance: pb.GetContractCreateInstance(),
		}
	case *proto.SchedulableTransactionBody_ContractUpdateInstance:
		pbBody.Data = &proto.TransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: pb.GetContractUpdateInstance(),
		}
	case *proto.SchedulableTransactionBody_ContractDeleteInstance:
		pbBody.Data = &proto.TransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: pb.GetContractDeleteInstance(),
		}
	case *proto.SchedulableTransactionBody_CryptoCreateAccount:
		pbBody.Data = &proto.TransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: pb.GetCryptoCreateAccount(),
		}
	case *proto.SchedulableTransactionBody_CryptoDelete:
		pbBody.Data = &proto.TransactionBody_CryptoDelete{
			CryptoDelete: pb.GetCryptoDelete(),
		}
	case *proto.SchedulableTransactionBody_CryptoTransfer:
		pbBody.Data = &proto.TransactionBody_CryptoTransfer{
			CryptoTransfer: pb.GetCryptoTransfer(),
		}
	case *proto.SchedulableTransactionBody_CryptoUpdateAccount:
		pbBody.Data = &proto.TransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: pb.GetCryptoUpdateAccount(),
		}
	case *proto.SchedulableTransactionBody_FileAppend:
		pbBody.Data = &proto.TransactionBody_FileAppend{
			FileAppend: pb.GetFileAppend(),
		}
	case *proto.SchedulableTransactionBody_FileCreate:
		pbBody.Data = &proto.TransactionBody_FileCreate{
			FileCreate: pb.GetFileCreate(),
		}
	case *proto.SchedulableTransactionBody_FileDelete:
		pbBody.Data = &proto.TransactionBody_FileDelete{
			FileDelete: pb.GetFileDelete(),
		}
	case *proto.SchedulableTransactionBody_FileUpdate:
		pbBody.Data = &proto.TransactionBody_FileUpdate{
			FileUpdate: pb.GetFileUpdate(),
		}
	case *proto.SchedulableTransactionBody_SystemDelete:
		pbBody.Data = &proto.TransactionBody_SystemDelete{
			SystemDelete: pb.GetSystemDelete(),
		}
	case *proto.SchedulableTransactionBody_SystemUndelete:
		pbBody.Data = &proto.TransactionBody_SystemUndelete{
			SystemUndelete: pb.GetSystemUndelete(),
		}
	case *proto.SchedulableTransactionBody_Freeze:
		pbBody.Data = &proto.TransactionBody_Freeze{
			Freeze: pb.GetFreeze(),
		}
	case *proto.SchedulableTransactionBody_ConsensusCreateTopic:
		pbBody.Data = &proto.TransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: pb.GetConsensusCreateTopic(),
		}
	case *proto.SchedulableTransactionBody_ConsensusUpdateTopic:
		pbBody.Data = &proto.TransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: pb.GetConsensusUpdateTopic(),
		}
	case *proto.SchedulableTransactionBody_ConsensusDeleteTopic:
		pbBody.Data = &proto.TransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: pb.GetConsensusDeleteTopic(),
		}
	case *proto.SchedulableTransactionBody_ConsensusSubmitMessage:
		pbBody.Data = &proto.TransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: pb.GetConsensusSubmitMessage(),
		}
	case *proto.SchedulableTransactionBody_TokenCreation:
		pbBody.Data = &proto.TransactionBody_TokenCreation{
			TokenCreation: pb.GetTokenCreation(),
		}
	case *proto.SchedulableTransactionBody_TokenFreeze:
		pbBody.Data = &proto.TransactionBody_TokenFreeze{
			TokenFreeze: pb.GetTokenFreeze(),
		}
	case *proto.SchedulableTransactionBody_TokenUnfreeze:
		pbBody.Data = &proto.TransactionBody_TokenUnfreeze{
			TokenUnfreeze: pb.GetTokenUnfreeze(),
		}
	case *proto.SchedulableTransactionBody_TokenGrantKyc:
		pbBody.Data = &proto.TransactionBody_TokenGrantKyc{
			TokenGrantKyc: pb.GetTokenGrantKyc(),
		}
	case *proto.SchedulableTransactionBody_TokenRevokeKyc:
		pbBody.Data = &proto.TransactionBody_TokenRevokeKyc{
			TokenRevokeKyc: pb.GetTokenRevokeKyc(),
		}
	case *proto.SchedulableTransactionBody_TokenDeletion:
		pbBody.Data = &proto.TransactionBody_TokenDeletion{
			TokenDeletion: pb.GetTokenDeletion(),
		}
	case *proto.SchedulableTransactionBody_TokenUpdate:
		pbBody.Data = &proto.TransactionBody_TokenUpdate{
			TokenUpdate: pb.GetTokenUpdate(),
		}
	case *proto.SchedulableTransactionBody_TokenMint:
		pbBody.Data = &proto.TransactionBody_TokenMint{
			TokenMint: pb.GetTokenMint(),
		}
	case *proto.SchedulableTransactionBody_TokenBurn:
		pbBody.Data = &proto.TransactionBody_TokenBurn{
			TokenBurn: pb.GetTokenBurn(),
		}
	case *proto.SchedulableTransactionBody_TokenWipe:
		pbBody.Data = &proto.TransactionBody_TokenWipe{
			TokenWipe: pb.GetTokenWipe(),
		}
	case *proto.SchedulableTransactionBody_TokenAssociate:
		pbBody.Data = &proto.TransactionBody_TokenAssociate{
			TokenAssociate: pb.GetTokenAssociate(),
		}
	case *proto.SchedulableTransactionBody_TokenDissociate:
		pbBody.Data = &proto.TransactionBody_TokenDissociate{
			TokenDissociate: pb.GetTokenDissociate(),
		}
	case *proto.SchedulableTransactionBody_ScheduleDelete:
		pbBody.Data = &proto.TransactionBody_ScheduleDelete{
			ScheduleDelete: pb.GetScheduleDelete(),
		}
	default:
		pbBody.Data = nil
	}

	return &SchedulableTransactionBody{
		TransactionFee: pb.TransactionFee,
		Memo:           pb.Memo,
		Transaction: &Transaction{
			pbBody:               &pbBody,
			nextNodeIndex:        0,
			nextTransactionIndex: 0,
			maxRetry:             10,
			transactionIDs:       make([]TransactionID, 0),
			transactions:         make([]*proto.Transaction, 0),
			signedTransactions:   make([]*proto.SignedTransaction, 0),
			nodeIDs:              make([]AccountID, 0),
		},
	}
}

func (body *SchedulableTransactionBody) toProtobuf() *proto.SchedulableTransactionBody {

	pbBody := proto.SchedulableTransactionBody{
		TransactionFee: body.TransactionFee,
		Memo:           body.Memo,
	}

	data := body.Transaction.pbBody.Data

	switch pb := data.(type) {
	case *proto.TransactionBody_ContractCall:
		pbBody.Data = &proto.SchedulableTransactionBody_ContractCall{
			ContractCall: pb.ContractCall,
		}
	case *proto.TransactionBody_ContractCreateInstance:
		pbBody.Data = &proto.SchedulableTransactionBody_ContractCreateInstance{
			ContractCreateInstance: pb.ContractCreateInstance,
		}
	case *proto.TransactionBody_ContractUpdateInstance:
		pbBody.Data = &proto.SchedulableTransactionBody_ContractUpdateInstance{
			ContractUpdateInstance: pb.ContractUpdateInstance,
		}
	case *proto.TransactionBody_ContractDeleteInstance:
		pbBody.Data = &proto.SchedulableTransactionBody_ContractDeleteInstance{
			ContractDeleteInstance: pb.ContractDeleteInstance,
		}
	case *proto.TransactionBody_CryptoCreateAccount:
		pbBody.Data = &proto.SchedulableTransactionBody_CryptoCreateAccount{
			CryptoCreateAccount: pb.CryptoCreateAccount,
		}
	case *proto.TransactionBody_CryptoDelete:
		pbBody.Data = &proto.SchedulableTransactionBody_CryptoDelete{
			CryptoDelete: pb.CryptoDelete,
		}
	case *proto.TransactionBody_CryptoTransfer:
		pbBody.Data = &proto.SchedulableTransactionBody_CryptoTransfer{
			CryptoTransfer: pb.CryptoTransfer,
		}
	case *proto.TransactionBody_CryptoUpdateAccount:
		pbBody.Data = &proto.SchedulableTransactionBody_CryptoUpdateAccount{
			CryptoUpdateAccount: pb.CryptoUpdateAccount,
		}
	case *proto.TransactionBody_FileAppend:
		pbBody.Data = &proto.SchedulableTransactionBody_FileAppend{
			FileAppend: pb.FileAppend,
		}
	case *proto.TransactionBody_FileCreate:
		pbBody.Data = &proto.SchedulableTransactionBody_FileCreate{
			FileCreate: pb.FileCreate,
		}
	case *proto.TransactionBody_FileDelete:
		pbBody.Data = &proto.SchedulableTransactionBody_FileDelete{
			FileDelete: pb.FileDelete,
		}
	case *proto.TransactionBody_FileUpdate:
		pbBody.Data = &proto.SchedulableTransactionBody_FileUpdate{
			FileUpdate: pb.FileUpdate,
		}
	case *proto.TransactionBody_SystemDelete:
		pbBody.Data = &proto.SchedulableTransactionBody_SystemDelete{
			SystemDelete: pb.SystemDelete,
		}
	case *proto.TransactionBody_SystemUndelete:
		pbBody.Data = &proto.SchedulableTransactionBody_SystemUndelete{
			SystemUndelete: pb.SystemUndelete,
		}
	case *proto.TransactionBody_Freeze:
		pbBody.Data = &proto.SchedulableTransactionBody_Freeze{
			Freeze: pb.Freeze,
		}
	case *proto.TransactionBody_ConsensusCreateTopic:
		pbBody.Data = &proto.SchedulableTransactionBody_ConsensusCreateTopic{
			ConsensusCreateTopic: pb.ConsensusCreateTopic,
		}
	case *proto.TransactionBody_ConsensusUpdateTopic:
		pbBody.Data = &proto.SchedulableTransactionBody_ConsensusUpdateTopic{
			ConsensusUpdateTopic: pb.ConsensusUpdateTopic,
		}
	case *proto.TransactionBody_ConsensusDeleteTopic:
		pbBody.Data = &proto.SchedulableTransactionBody_ConsensusDeleteTopic{
			ConsensusDeleteTopic: pb.ConsensusDeleteTopic,
		}
	case *proto.TransactionBody_ConsensusSubmitMessage:
		pbBody.Data = &proto.SchedulableTransactionBody_ConsensusSubmitMessage{
			ConsensusSubmitMessage: pb.ConsensusSubmitMessage,
		}
	case *proto.TransactionBody_TokenCreation:
		pbBody.Data = &proto.SchedulableTransactionBody_TokenCreation{
			TokenCreation: pb.TokenCreation,
		}
	case *proto.TransactionBody_TokenFreeze:
		pbBody.Data = &proto.SchedulableTransactionBody_TokenFreeze{
			TokenFreeze: pb.TokenFreeze,
		}
	case *proto.TransactionBody_TokenUnfreeze:
		pbBody.Data = &proto.SchedulableTransactionBody_TokenUnfreeze{
			TokenUnfreeze: pb.TokenUnfreeze,
		}
	case *proto.TransactionBody_TokenGrantKyc:
		pbBody.Data = &proto.SchedulableTransactionBody_TokenGrantKyc{
			TokenGrantKyc: pb.TokenGrantKyc,
		}
	case *proto.TransactionBody_TokenRevokeKyc:
		pbBody.Data = &proto.SchedulableTransactionBody_TokenRevokeKyc{
			TokenRevokeKyc: pb.TokenRevokeKyc,
		}
	case *proto.TransactionBody_TokenDeletion:
		pbBody.Data = &proto.SchedulableTransactionBody_TokenDeletion{
			TokenDeletion: pb.TokenDeletion,
		}
	case *proto.TransactionBody_TokenUpdate:
		pbBody.Data = &proto.SchedulableTransactionBody_TokenUpdate{
			TokenUpdate: pb.TokenUpdate,
		}
	case *proto.TransactionBody_TokenMint:
		pbBody.Data = &proto.SchedulableTransactionBody_TokenMint{
			TokenMint: pb.TokenMint,
		}
	case *proto.TransactionBody_TokenBurn:
		pbBody.Data = &proto.SchedulableTransactionBody_TokenBurn{
			TokenBurn: pb.TokenBurn,
		}
	case *proto.TransactionBody_TokenWipe:
		pbBody.Data = &proto.SchedulableTransactionBody_TokenWipe{
			TokenWipe: pb.TokenWipe,
		}
	case *proto.TransactionBody_TokenAssociate:
		pbBody.Data = &proto.SchedulableTransactionBody_TokenAssociate{
			TokenAssociate: pb.TokenAssociate,
		}
	case *proto.TransactionBody_TokenDissociate:
		pbBody.Data = &proto.SchedulableTransactionBody_TokenDissociate{
			TokenDissociate: pb.TokenDissociate,
		}
	case *proto.TransactionBody_ScheduleDelete:
		pbBody.Data = &proto.SchedulableTransactionBody_ScheduleDelete{
			ScheduleDelete: pb.ScheduleDelete,
		}
	default:
		pbBody.Data = nil
	}

	return &pbBody
}
