package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
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
	AliasKey                       *PublicKey
	LedgerID                       LedgerID
	HbarAllowances                 []HbarAllowance
	NftAllowances                  []TokenNftAllowance
	TokenAllowances                []TokenAllowance
}

func _AccountInfoFromProtobuf(pb *services.CryptoGetInfoResponse_AccountInfo) (AccountInfo, error) {
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

	var alias *PublicKey
	if len(pb.Alias) != 0 {
		pbKey := services.Key{}
		_ = protobuf.Unmarshal(pb.Alias, &pbKey)
		initialKey, _ := _KeyFromProtobuf(&pbKey)
		switch t2 := initialKey.(type) { //nolint
		case PublicKey:
			alias = &t2
		}
	}

	hbarAllowances := make([]HbarAllowance, len(pb.GrantedCryptoAllowances))
	if len(pb.GrantedCryptoAllowances) > 0 {
		for _, allowance := range pb.GrantedCryptoAllowances {
			hbarAllowances = append(hbarAllowances, _HbarAllowanceFromGrantedProtobuf(allowance))
		}
	}

	tokenAllowances := make([]TokenAllowance, len(pb.GrantedTokenAllowances))
	if len(pb.GrantedTokenAllowances) > 0 {
		for _, allowance := range pb.GrantedTokenAllowances {
			tokenAllowances = append(tokenAllowances, _TokenAllowanceFromGrantedProtobuf(allowance))
		}
	}

	nftAllowances := make([]TokenNftAllowance, len(pb.GrantedNftAllowances))
	if len(pb.GrantedNftAllowances) > 0 {
		for _, allowance := range pb.GrantedNftAllowances {
			nftAllowances = append(nftAllowances, _TokenNftAllowanceFromGrantedProtobuf(allowance))
		}
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
		MaxAutomaticTokenAssociations:  uint32(pb.MaxAutomaticTokenAssociations),
		AliasKey:                       alias,
		LedgerID:                       LedgerID{pb.LedgerId},
		HbarAllowances:                 hbarAllowances,
		NftAllowances:                  nftAllowances,
		TokenAllowances:                tokenAllowances,
	}, nil
}

func (info AccountInfo) _ToProtobuf() *services.CryptoGetInfoResponse_AccountInfo {
	tokenRelationship := make([]*services.TokenRelationship, len(info.TokenRelationships))

	for i, relationship := range info.TokenRelationships {
		singleRelationship := relationship._ToProtobuf()
		tokenRelationship[i] = singleRelationship
	}

	liveHashes := make([]*services.LiveHash, len(info.LiveHashes))

	for i, liveHash := range info.LiveHashes {
		singleRelationship := liveHash._ToProtobuf()
		liveHashes[i] = singleRelationship
	}

	var alias []byte
	if info.AliasKey != nil {
		alias, _ = protobuf.Marshal(info.AliasKey._ToProtoKey())
	}

	body := &services.CryptoGetInfoResponse_AccountInfo{
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
		MaxAutomaticTokenAssociations:  int32(info.MaxAutomaticTokenAssociations),
		Alias:                          alias,
		LedgerId:                       info.LedgerID.ToBytes(),
	}

	hbarAllowances := make([]*services.GrantedCryptoAllowance, 0)
	if len(info.HbarAllowances) > 0 {
		for _, allowance := range info.HbarAllowances {
			hbarAllowances = append(hbarAllowances, allowance._ToGrantedProtobuf())
		}
		body.GrantedCryptoAllowances = hbarAllowances
	}

	tokenAllowances := make([]*services.GrantedTokenAllowance, 0)
	if len(info.TokenAllowances) > 0 {
		for _, allowance := range info.TokenAllowances {
			tokenAllowances = append(tokenAllowances, allowance._ToGrantedProtobuf())
		}
		body.GrantedTokenAllowances = tokenAllowances
	}

	nftAllowances := make([]*services.GrantedNftAllowance, 0)
	if len(info.NftAllowances) > 0 {
		for _, allowance := range info.NftAllowances {
			nftAllowances = append(nftAllowances, allowance._ToGrantedProtobuf())
		}
		body.GrantedNftAllowances = nftAllowances
	}

	return body
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
	pb := services.CryptoGetInfoResponse_AccountInfo{}
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
