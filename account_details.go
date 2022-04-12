package hedera

import (
	"time"

	"github.com/hashgraph/hedera-protobufs-go/services"
	protobuf "google.golang.org/protobuf/proto"
)

type AccountDetails struct {
	AccountId                     *AccountID
	ContractAccountId             string
	Deleted                       bool
	ProxyAccountId                *AccountID
	ProxyReceived                 int64
	Key                           Key
	Balance                       uint64
	ReceiverSigRequired           bool
	ExpirationTime                *time.Time
	AutoRenewPeriod               *time.Duration
	TokenRelationships            []TokenRelationship
	Memo                          string
	OwnedNfts                     int64
	MaxAutomaticTokenAssociations int32
	Alias                         *PublicKey
	LedgerId                      LedgerID
	HbarAllowances                []HbarAllowance
	TokenNftAllowances            []TokenNftAllowance
	TokenAllowances               []TokenAllowance
}

func _AccountDetailsFromProtobuf(pb *services.GetAccountDetailsResponse_AccountDetails) (AccountDetails, error) {
	var err error
	var accountID *AccountID
	if pb.AccountId != nil {
		accountID = _AccountIDFromProtobuf(pb.AccountId)
	}

	var proxyAccountID *AccountID
	if pb.ProxyAccountId != nil {
		proxyAccountID = _AccountIDFromProtobuf(pb.ProxyAccountId)
	}

	var key Key
	if pb.Key != nil {
		key, err = _KeyFromProtobuf(pb.Key)
		if err != nil {
			return AccountDetails{}, err
		}
	}

	var expirationTime time.Time
	if pb.ExpirationTime != nil {
		expirationTime = _TimeFromProtobuf(pb.ExpirationTime)
	}

	var autoRenewPeriod time.Duration
	if pb.AutoRenewPeriod != nil {
		autoRenewPeriod = _DurationFromProtobuf(pb.AutoRenewPeriod)
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

	tokenRelationship := make([]TokenRelationship, 0)
	for _, relation := range pb.TokenRelationships {
		token := _TokenRelationshipFromProtobuf(relation)
		tokenRelationship = append(tokenRelationship, token)
	}

	hbarAllowances := make([]HbarAllowance, len(pb.GrantedCryptoAllowances))
	if len(pb.GrantedCryptoAllowances) > 0 {
		for _, allowance := range pb.GrantedCryptoAllowances {
			hbarAllowance := _HbarAllowanceFromGrantedProtobuf(allowance)
			hbarAllowance.OwnerAccountID = accountID
			hbarAllowances = append(hbarAllowances, hbarAllowance)
		}
	}

	tokenAllowances := make([]TokenAllowance, len(pb.GrantedTokenAllowances))
	if len(pb.GrantedTokenAllowances) > 0 {
		for _, allowance := range pb.GrantedTokenAllowances {
			tokenAllowance := _TokenAllowanceFromGrantedProtobuf(allowance)
			tokenAllowance.OwnerAccountID = accountID
			tokenAllowances = append(tokenAllowances, tokenAllowance)
		}
	}

	nftAllowances := make([]TokenNftAllowance, len(pb.GrantedNftAllowances))
	if len(pb.GrantedNftAllowances) > 0 {
		for _, allowance := range pb.GrantedNftAllowances {
			nftAllowance := _TokenNftAllowanceFromGrantedProtobuf(allowance)
			nftAllowance.OwnerAccountID = accountID
			nftAllowances = append(nftAllowances, nftAllowance)
		}
	}

	return AccountDetails{
		AccountId:                     accountID,
		ContractAccountId:             pb.ContractAccountId,
		Deleted:                       pb.Deleted,
		ProxyAccountId:                proxyAccountID,
		ProxyReceived:                 pb.ProxyReceived,
		Key:                           key,
		Balance:                       pb.Balance,
		ReceiverSigRequired:           pb.ReceiverSigRequired,
		ExpirationTime:                &expirationTime,
		AutoRenewPeriod:               &autoRenewPeriod,
		TokenRelationships:            tokenRelationship,
		Memo:                          pb.Memo,
		OwnedNfts:                     pb.OwnedNfts,
		MaxAutomaticTokenAssociations: pb.MaxAutomaticTokenAssociations,
		Alias:                         alias,
		LedgerId:                      LedgerID{pb.LedgerId},
		HbarAllowances:                hbarAllowances,
		TokenNftAllowances:            nftAllowances,
		TokenAllowances:               tokenAllowances,
	}, nil
}

func (info AccountDetails) _ToProtobuf() *services.GetAccountDetailsResponse_AccountDetails {
	var accountID *services.AccountID
	if info.AccountId != nil {
		accountID = info.AccountId._ToProtobuf()
	}

	var proxyAccountID *services.AccountID
	if info.ProxyAccountId != nil {
		proxyAccountID = info.ProxyAccountId._ToProtobuf()
	}

	var key *services.Key
	if info.Key != nil {
		key = info.Key._ToProtoKey()
	}

	var expirationTime *services.Timestamp
	if info.ExpirationTime != nil {
		expirationTime = _TimeToProtobuf(*info.ExpirationTime)
	}

	var autoRenewPeriod *services.Duration
	if info.AutoRenewPeriod != nil {
		autoRenewPeriod = _DurationToProtobuf(*info.AutoRenewPeriod)
	}

	var alias []byte
	if info.Alias != nil {
		alias, _ = protobuf.Marshal(info.Alias._ToProtoKey())
	}

	tokenRelationship := make([]*services.TokenRelationship, 0)
	for _, relation := range info.TokenRelationships {
		tokenRelationship = append(tokenRelationship, relation._ToProtobuf())
	}

	hbarAllowances := make([]*services.GrantedCryptoAllowance, 0)
	if len(info.HbarAllowances) > 0 {
		for _, allowance := range info.HbarAllowances {
			hbarAllowances = append(hbarAllowances, allowance._ToGrantedProtobuf())
		}
	}

	tokenAllowances := make([]*services.GrantedTokenAllowance, 0)
	if len(info.TokenAllowances) > 0 {
		for _, allowance := range info.TokenAllowances {
			tokenAllowances = append(tokenAllowances, allowance._ToGrantedProtobuf())
		}
	}

	nftAllowances := make([]*services.GrantedNftAllowance, 0)
	if len(info.TokenNftAllowances) > 0 {
		for _, allowance := range info.TokenNftAllowances {
			nftAllowances = append(nftAllowances, allowance._ToGrantedProtobuf())
		}
	}

	return &services.GetAccountDetailsResponse_AccountDetails{
		AccountId:                     accountID,
		ContractAccountId:             info.ContractAccountId,
		Deleted:                       info.Deleted,
		ProxyAccountId:                proxyAccountID,
		ProxyReceived:                 info.ProxyReceived,
		Key:                           key,
		Balance:                       info.Balance,
		ReceiverSigRequired:           info.ReceiverSigRequired,
		ExpirationTime:                expirationTime,
		AutoRenewPeriod:               autoRenewPeriod,
		TokenRelationships:            tokenRelationship,
		Memo:                          info.Memo,
		OwnedNfts:                     info.OwnedNfts,
		MaxAutomaticTokenAssociations: info.MaxAutomaticTokenAssociations,
		Alias:                         alias,
		LedgerId:                      info.LedgerId.ToBytes(),
		GrantedCryptoAllowances:       hbarAllowances,
		GrantedNftAllowances:          nftAllowances,
		GrantedTokenAllowances:        tokenAllowances,
	}
}
