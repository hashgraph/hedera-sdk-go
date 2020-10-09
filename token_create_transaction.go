package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenCreateTransaction struct {
	TransactionBuilder
	pb *proto.TokenCreateTransactionBody
}

func NewTokenCreateTransaction() TokenCreateTransaction {
	pb := &proto.TokenCreateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenCreation{TokenCreation: pb}

	builder := TokenCreateTransaction{inner, pb}
	builder.SetAutoRenewPeriod(7890000)

	return builder
}

// The publicly visible name of the token, specified as a string of only ASCII characters
func (builder TokenCreateTransaction) SetName(name string) TokenCreateTransaction {
	builder.pb.Name = name
	return builder
}

// The publicly visible token symbol. It is UTF-8 capitalized alphabetical string identifying the token
func (builder TokenCreateTransaction) SetSymbol(symbol string) TokenCreateTransaction {
	builder.pb.Symbol = symbol
	return builder
}

// The number of decimal places a token is divisible by. This field can never be changed!
func (builder TokenCreateTransaction) SetDecimals(decimals uint32) TokenCreateTransaction {
	builder.pb.Decimals = decimals
	return builder
}

// Specifies the initial supply of tokens to be put in circulation. The initial supply is sent to the Treasury Account. The supply is in the lowest denomination possible.
func (builder TokenCreateTransaction) SetInitialSupply(initialsupply uint64) TokenCreateTransaction {
	builder.pb.InitialSupply = initialsupply
	return builder
}

// The account which will act as a treasury for the token. This account will receive the specified initial supply
func (builder TokenCreateTransaction) SetTreasury(treasury AccountID) TokenCreateTransaction {
	builder.pb.Treasury = treasury.toProto()
	return builder
}

// The key which can perform update/delete operations on the token. If empty, the token can be perceived as immutable (not being able to be updated/deleted)
func (builder TokenCreateTransaction) SetAdminKey(adminkey PublicKey) TokenCreateTransaction {
	builder.pb.AdminKey = adminkey.toProto()
	return builder
}

// The key which can grant or revoke KYC of an account for the token's transactions. If empty, KYC is not required, and KYC grant or revoke operations are not possible.
func (builder TokenCreateTransaction) SetKycKey(kyckey PublicKey) TokenCreateTransaction {
	builder.pb.KycKey = kyckey.toProto()
	return builder
}

// The key which can sign to freeze or unfreeze an account for token transactions. If empty, freezing is not possible
func (builder TokenCreateTransaction) SetFreezeKey(freezekey PublicKey) TokenCreateTransaction {
	builder.pb.FreezeKey = freezekey.toProto()
	return builder
}

// The key which can wipe the token balance of an account. If empty, wipe is not possible
func (builder TokenCreateTransaction) SetWipeKey(wipekey PublicKey) TokenCreateTransaction {
	builder.pb.WipeKey = wipekey.toProto()
	return builder
}

// The key which can change the supply of a token. The key is used to sign Token Mint/Burn operations
func (builder TokenCreateTransaction) SetSupplyKey(supplykey PublicKey) TokenCreateTransaction {
	builder.pb.SupplyKey = supplykey.toProto()
	return builder
}

// The default Freeze status (frozen or unfrozen) of Hedera accounts relative to this token. If true, an account must be unfrozen before it can receive the token
func (builder TokenCreateTransaction) SetFreezeDefault(freeze bool) TokenCreateTransaction {
	builder.pb.FreezeDefault = freeze
	return builder
}

// The epoch second at which the token should expire; if an auto-renew account and period are specified, this is coerced to the current epoch second plus the autoRenewPeriod
func (builder TokenCreateTransaction) SetExpirationTime(expirationTime uint64) TokenCreateTransaction {
	builder.pb.Expiry = expirationTime
	return builder
}

// An account which will be automatically charged to renew the token's expiration, at autoRenewPeriod interval
func (builder TokenCreateTransaction) SetAutoRenewAccountID(autoRenewAccountID AccountID) TokenCreateTransaction {
	builder.pb.AutoRenewAccount = autoRenewAccountID.toProto()
	return builder
}

// The interval at which the auto-renew account will be charged to extend the token's expiry
func (builder TokenCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod uint64) TokenCreateTransaction {
	builder.pb.AutoRenewPeriod = autoRenewPeriod
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder TokenCreateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenCreateTransaction {
	return TokenCreateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenCreateTransaction) SetTransactionMemo(memo string) TokenCreateTransaction {
	return TokenCreateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenCreateTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenCreateTransaction {
	return TokenCreateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenCreateTransaction) SetTransactionID(transactionID TransactionID) TokenCreateTransaction {
	return TokenCreateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeTokenID sets the node TokenID for this Transaction.
func (builder TokenCreateTransaction) SetNodeTokenID(nodeAccountID AccountID) TokenCreateTransaction {
	return TokenCreateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
