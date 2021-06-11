package hedera

import (
	"time"

	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/v2/proto"
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
}

func accountInfoFromProtobuf(pb *proto.CryptoGetInfoResponse_AccountInfo, networkName *NetworkName) (AccountInfo, error) {
	if pb == nil {
		return AccountInfo{}, errParameterNull
	}

	pubKey, err := keyFromProtobuf(pb.Key, networkName)
	if err != nil {
		return AccountInfo{}, err
	}

	tokenRelationship := make([]*TokenRelationship, len(pb.TokenRelationships))

	if pb.TokenRelationships != nil {
		for i, relationship := range pb.TokenRelationships {
			singleRelationship := tokenRelationshipFromProtobuf(relationship, networkName)
			tokenRelationship[i] = &singleRelationship
		}
	}

	liveHashes := make([]*LiveHash, len(pb.LiveHashes))

	if pb.LiveHashes != nil {
		for i, liveHash := range pb.LiveHashes {
			singleRelationship, err := liveHashFromProtobuf(liveHash)
			if err != nil {
				return AccountInfo{}, err
			}
			liveHashes[i] = &singleRelationship
		}
	}

	var proxyAccountID AccountID
	if pb.ProxyAccountID != nil {
		proxyAccountID = accountIDFromProtobuf(pb.ProxyAccountID, networkName)
	}

	return AccountInfo{
		AccountID:                      accountIDFromProtobuf(pb.AccountID, networkName),
		ContractAccountID:              pb.ContractAccountID,
		IsDeleted:                      pb.Deleted,
		ProxyAccountID:                 proxyAccountID,
		ProxyReceived:                  HbarFromTinybar(pb.ProxyReceived),
		Key:                            pubKey,
		Balance:                        HbarFromTinybar(int64(pb.Balance)),
		GenerateSendRecordThreshold:    HbarFromTinybar(int64(pb.GenerateSendRecordThreshold)),
		GenerateReceiveRecordThreshold: HbarFromTinybar(int64(pb.GenerateReceiveRecordThreshold)),
		ReceiverSigRequired:            pb.ReceiverSigRequired,
		TokenRelationships:             tokenRelationship,
		ExpirationTime:                 timeFromProtobuf(pb.ExpirationTime),
		AccountMemo:                    pb.Memo,
		AutoRenewPeriod:                durationFromProtobuf(pb.AutoRenewPeriod),
		LiveHashes:                     liveHashes,
		OwnedNfts:                      pb.OwnedNfts,
	}, nil
}

func (info AccountInfo) toProtobuf() *proto.CryptoGetInfoResponse_AccountInfo {

	tokenRelationship := make([]*proto.TokenRelationship, len(info.TokenRelationships))

	for i, relationship := range info.TokenRelationships {
		singleRelationship := relationship.toProtobuf()
		tokenRelationship[i] = singleRelationship
	}

	liveHashes := make([]*proto.LiveHash, len(info.LiveHashes))

	for i, liveHash := range info.LiveHashes {
		singleRelationship := liveHash.toProtobuf()
		liveHashes[i] = singleRelationship
	}

	return &proto.CryptoGetInfoResponse_AccountInfo{
		AccountID:                      info.AccountID.toProtobuf(),
		ContractAccountID:              info.ContractAccountID,
		Deleted:                        info.IsDeleted,
		ProxyAccountID:                 info.ProxyAccountID.toProtobuf(),
		ProxyReceived:                  info.ProxyReceived.tinybar,
		Key:                            info.Key.toProtoKey(),
		Balance:                        uint64(info.Balance.tinybar),
		GenerateSendRecordThreshold:    uint64(info.GenerateSendRecordThreshold.tinybar),
		GenerateReceiveRecordThreshold: uint64(info.GenerateReceiveRecordThreshold.tinybar),
		ReceiverSigRequired:            info.ReceiverSigRequired,
		ExpirationTime:                 timeToProtobuf(info.ExpirationTime),
		AutoRenewPeriod:                durationToProtobuf(info.AutoRenewPeriod),
		LiveHashes:                     liveHashes,
		TokenRelationships:             tokenRelationship,
		Memo:                           info.AccountMemo,
		OwnedNfts:                      info.OwnedNfts,
	}
}

func (info AccountInfo) ToBytes() []byte {
	data, err := protobuf.Marshal(info.toProtobuf())
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

	info, err := accountInfoFromProtobuf(&pb, nil)
	if err != nil {
		return AccountInfo{}, err
	}

	return info, nil
}
