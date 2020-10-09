package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TokenUpdateTransaction struct {
	TransactionBuilder
	pb *proto.TokenUpdateTransactionBody
}

func NewTokenUpdateTransaction() TokenUpdateTransaction {
	pb := &proto.TokenUpdateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_TokenUpdate{TokenUpdate: pb}

	builder := TokenUpdateTransaction{inner, pb}

	return builder
}

// The Token to be updated
func (builder TokenUpdateTransaction) SetTokenID(id TokenID) TokenUpdateTransaction {
	builder.pb.Token = id.toProto()
	return builder
}

// The new Symbol of the Token. Must be UTF-8 capitalized alphabetical string identifying the token.
func (builder TokenUpdateTransaction) SetName(name string) TokenUpdateTransaction {
	builder.pb.Name = name
	return builder
}

// The new Name of the Token. Must be a string of ASCII characters.
func (builder TokenUpdateTransaction) SetSymbol(symbol string) TokenUpdateTransaction {
	builder.pb.Symbol = symbol
	return builder
}

// The new Treasury account of the Token. If the provided treasury account is not existing or deleted, the response will be INVALID_TREASURY_ACCOUNT_FOR_TOKEN. If successful, the Token balance held in the previous Treasury Account is transferred to the new one.
func (builder TokenUpdateTransaction) SetTreasury(treasury AccountID) TokenUpdateTransaction {
	builder.pb.Treasury = treasury.toProto()
	return builder
}

// The new KYC key of the Token. If Token does not have currently a KYC key, transaction will resolve to TOKEN_HAS_NO_KYC_KEY.
func (builder TokenUpdateTransaction) SetAdminKey(adminkey PublicKey) TokenUpdateTransaction {
	builder.pb.AdminKey = adminkey.toProto()
	return builder
}

// The new KYC key of the Token. If Token does not have currently a KYC key, transaction will resolve to TOKEN_HAS_NO_KYC_KEY.
func (builder TokenUpdateTransaction) SetKycKey(kyckey PublicKey) TokenUpdateTransaction {
	builder.pb.KycKey = kyckey.toProto()
	return builder
}

// The new Freeze key of the Token. If the Token does not have currently a Freeze key, transaction will resolve to TOKEN_HAS_NO_FREEZE_KEY.
func (builder TokenUpdateTransaction) SetFreezeKey(freezekey PublicKey) TokenUpdateTransaction {
	builder.pb.FreezeKey = freezekey.toProto()
	return builder
}

// The new Wipe key of the Token. If the Token does not have currently a Wipe key, transaction will resolve to TOKEN_HAS_NO_WIPE_KEY.
func (builder TokenUpdateTransaction) SetWipeKey(wipekey PublicKey) TokenUpdateTransaction {
	builder.pb.WipeKey = wipekey.toProto()
	return builder
}

// The new Supply key of the Token. If the Token does not have currently a Supply key, transaction will resolve to TOKEN_HAS_NO_SUPPLY_KEY.
func (builder TokenUpdateTransaction) SetSupplyKey(supplykey PublicKey) TokenUpdateTransaction {
	builder.pb.SupplyKey = supplykey.toProto()
	return builder
}

// The new expiry time of the token. Expiry can be updated even if admin key is not set. If the provided expiry is earlier than the current token expiry, transaction wil resolve to INVALID_EXPIRATION_TIME
func (builder TokenUpdateTransaction) SetExpirationTime(expirationTime uint64) TokenUpdateTransaction {
	builder.pb.Expiry = expirationTime
	return builder
}

// The new account which will be automatically charged to renew the token's expiration, at autoRenewPeriod interval.
func (builder TokenUpdateTransaction) SetAutoRenewAccountID(autoRenewAccountID AccountID) TokenUpdateTransaction {
	builder.pb.AutoRenewAccount = autoRenewAccountID.toProto()
	return builder
}

// The new interval at which the auto-renew account will be charged to extend the token's expiry.
func (builder TokenUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod uint64) TokenUpdateTransaction {
	builder.pb.AutoRenewPeriod = autoRenewPeriod
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder TokenUpdateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) TokenUpdateTransaction {
	return TokenUpdateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder TokenUpdateTransaction) SetTransactionMemo(memo string) TokenUpdateTransaction {
	return TokenUpdateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder TokenUpdateTransaction) SetTransactionValidDuration(validDuration time.Duration) TokenUpdateTransaction {
	return TokenUpdateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder TokenUpdateTransaction) SetTransactionID(transactionID TransactionID) TokenUpdateTransaction {
	return TokenUpdateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeTokenID sets the node TokenID for this Transaction.
func (builder TokenUpdateTransaction) SetNodeAccountID(nodeAccountID AccountID) TokenUpdateTransaction {
	return TokenUpdateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
