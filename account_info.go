package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	protobuf "google.golang.org/protobuf/proto"
)

// AccountInfo is info about the account returned from an AccountInfoQuery
type AccountInfo struct {
	AccountID                      AccountID
	ContractAccountID              string
	IsDeleted                      bool
	ProxyAccountID                 AccountID
	ProxyReceived                  Hbar
	Key                            Key
	Balance                        Hbar
	GenerateSendRecordThreshold    Hbar
	GenerateReceiveRecordThreshold Hbar
	ReceiverSigRequired            bool
	ExpirationTime                 time.Time
	AutoRenewPeriod                time.Duration
	LiveHashes                     []*LiveHash
	TokenRelationships             []*TokenRelationship
	AccountMemo                    string
	OwnedNfts                      int64
	MaxAutomaticTokenAssociations  uint32
}

func _AccountInfoFromProtobuf(pb *proto.CryptoGetInfoResponse_AccountInfo) (AccountInfo, error) {
	if pb == nil {
		return AccountInfo{}, errParameterNull
	}

	pubKey, err := _KeyFromProtobuf(pb.Key)
	if err != nil {
		return AccountInfo{}, err
	}

	tokenRelationship := make([]*TokenRelationship, len(pb.TokenRelationships))

	if pb.TokenRelationships != nil {
		for i, relationship := range pb.TokenRelationships {
			singleRelationship := _TokenRelationshipFromProtobuf(relationship)
			tokenRelationship[i] = &singleRelationship
		}
	}

	liveHashes := make([]*LiveHash, len(pb.LiveHashes))

	if pb.LiveHashes != nil {
		for i, liveHash := range pb.LiveHashes {
			singleRelationship, err := _LiveHashFromProtobuf(liveHash)
			if err != nil {
				return AccountInfo{}, err
			}
			liveHashes[i] = &singleRelationship
		}
	}

	proxyAccountID := AccountID{}
	if pb.ProxyAccountID != nil {
		proxyAccountID = *_AccountIDFromProtobuf(pb.ProxyAccountID)
	}

	accountID := AccountID{}
	if pb.AccountID != nil {
		accountID = *_AccountIDFromProtobuf(pb.AccountID)
	}

	return AccountInfo{
		AccountID:                      accountID,
		ContractAccountID:              pb.ContractAccountID,
		IsDeleted:                      pb.Deleted,
		ProxyAccountID:                 proxyAccountID,
		ProxyReceived:                  HbarFromTinybar(pb.ProxyReceived),
		Key:                            pubKey,
		Balance:                        HbarFromTinybar(int64(pb.Balance)),
		GenerateSendRecordThreshold:    HbarFromTinybar(int64(pb.GenerateSendRecordThreshold)),    // nolint
		GenerateReceiveRecordThreshold: HbarFromTinybar(int64(pb.GenerateReceiveRecordThreshold)), // nolint
		ReceiverSigRequired:            pb.ReceiverSigRequired,
		TokenRelationships:             tokenRelationship,
		ExpirationTime:                 _TimeFromProtobuf(pb.ExpirationTime),
		AccountMemo:                    pb.Memo,
		AutoRenewPeriod:                _DurationFromProtobuf(pb.AutoRenewPeriod),
		LiveHashes:                     liveHashes,
		OwnedNfts:                      pb.OwnedNfts,
		MaxAutomaticTokenAssociations:  pb.MaxAutomaticTokenAssociations,
	}, nil
}

func (info AccountInfo) _ToProtobuf() *proto.CryptoGetInfoResponse_AccountInfo {
	tokenRelationship := make([]*proto.TokenRelationship, len(info.TokenRelationships))

	for i, relationship := range info.TokenRelationships {
		singleRelationship := relationship._ToProtobuf()
		tokenRelationship[i] = singleRelationship
	}

	liveHashes := make([]*proto.LiveHash, len(info.LiveHashes))

	for i, liveHash := range info.LiveHashes {
		singleRelationship := liveHash._ToProtobuf()
		liveHashes[i] = singleRelationship
	}

	return &proto.CryptoGetInfoResponse_AccountInfo{
		AccountID:                      info.AccountID._ToProtobuf(),
		ContractAccountID:              info.ContractAccountID,
		Deleted:                        info.IsDeleted,
		ProxyAccountID:                 info.ProxyAccountID._ToProtobuf(),
		ProxyReceived:                  info.ProxyReceived.tinybar,
		Key:                            info.Key._ToProtoKey(),
		Balance:                        uint64(info.Balance.tinybar),
		GenerateSendRecordThreshold:    uint64(info.GenerateSendRecordThreshold.tinybar),
		GenerateReceiveRecordThreshold: uint64(info.GenerateReceiveRecordThreshold.tinybar),
		ReceiverSigRequired:            info.ReceiverSigRequired,
		ExpirationTime:                 _TimeToProtobuf(info.ExpirationTime),
		AutoRenewPeriod:                _DurationToProtobuf(info.AutoRenewPeriod),
		LiveHashes:                     liveHashes,
		TokenRelationships:             tokenRelationship,
		Memo:                           info.AccountMemo,
		OwnedNfts:                      info.OwnedNfts,
		MaxAutomaticTokenAssociations:  info.MaxAutomaticTokenAssociations,
	}
}

func (info AccountInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(info._ToProtobuf())
	if err != nil {
		return make([]byte, 0)
	}

	return data
}

func AccountInfoFromBytes(data []byte) (AccountInfo, error) {
	if data == nil {
		return AccountInfo{}, errByteArrayNull
	}
	pb := proto.CryptoGetInfoResponse_AccountInfo{}
	err := protobuf.Unmarshal(data, &pb)
	if err != nil {
		return AccountInfo{}, err
	}

	info, err := _AccountInfoFromProtobuf(&pb)
	if err != nil {
		return AccountInfo{}, err
	}

	return info, nil
}
