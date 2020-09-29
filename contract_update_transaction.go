package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// ContractUpdateTransaction is used to modify a smart contract instance to have the given parameter values. Any nil
// field is ignored (left unchanged). If only the contractInstanceExpirationTime is being modified, then no signature is
// needed on this transaction other than for the account paying for the transaction itself. But if any of the other
// fields are being modified, then it must be signed by the adminKey. The use of adminKey is not currently supported in
// this API, but in the future will be implemented to allow these fields to be modified, and also to make modifications
// to the state of the instance. If the contract is created with no admin key, then none of the fields can be changed
// that need an admin signature, and therefore no admin key can ever be added. So if there is no admin key, then things
// like the bytecode are immutable. But if there is an admin key, then they can be changed.
//
// For example, the admin key might be a threshold key, which requires 3 of 5 binding arbitration judges to agree before
// the bytecode can be changed. This can be used to add flexibility to the management of smart contract behavior. But
// this is optional. If the smart contract is created without an admin key, then such a key can never be added, and its
// bytecode will be immutable.
type ContractUpdateTransaction struct {
	TransactionBuilder
	pb *proto.ContractUpdateTransactionBody
}

// NewContractUpdateTransaction creates a ContractUpdateTransaction transaction which can be
// used to construct and execute a Contract Update Transaction.
func NewContractUpdateTransaction() ContractUpdateTransaction {
	pb := &proto.ContractUpdateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ContractUpdateInstance{ContractUpdateInstance: pb}

	transaction := ContractUpdateTransaction{inner, pb}

	return transaction
}

// SetContractID sets The Contract ID instance to update (this can't be changed on the contract)
func (transaction ContractUpdateTransaction) SetContractID(id ContractID) ContractUpdateTransaction {
	transaction.pb.ContractID = id.toProto()
	return transaction
}

// SetBytecodeFileID sets the file ID of file containing the smart contract byte code. A copy will be made and held by
// the contract instance, and have the same expiration time as the instance.
func (transaction ContractUpdateTransaction) SetBytecodeFileID(id FileID) ContractUpdateTransaction {
	transaction.pb.FileID = id.toProto()
	return transaction
}

// SetAdminKey sets the key which can be used to arbitrarily modify the state of the instance by signing a
// ContractUpdateTransaction to modify it. If the admin key was never set then such modifications are not possible,
// and there is no administrator that can override the normal operation of the smart contract instance.
func (transaction ContractUpdateTransaction) SetAdminKey(publicKey PublicKey) ContractUpdateTransaction {
	transaction.pb.AdminKey = publicKey.toProto()
	return transaction
}

// SetProxyAccountID sets the ID of the account to which this contract is proxy staked. If proxyAccountID is left unset,
// is an invalid account, or is an account that isn't a node, then this contract is automatically proxy staked to a node
// chosen by the network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking,
// or if it is not currently running a node, then it will behave as if proxyAccountID was never set.
func (transaction ContractUpdateTransaction) SetProxyAccountID(id AccountID) ContractUpdateTransaction {
	transaction.pb.ProxyAccountID = id.toProto()
	return transaction
}

// SetAutoRenewPeriod sets the duration for which the contract instance will automatically charge its account to
// renew for.
func (transaction ContractUpdateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) ContractUpdateTransaction {
	transaction.pb.AutoRenewPeriod = durationToProto(autoRenewPeriod)
	return transaction
}

// SetExpirationTime extends the expiration of the instance and its account to the provided time. If the time provided
// is the current or past time, then there will be no effect.
func (transaction ContractUpdateTransaction) SetExpirationTime(expiration time.Time) ContractUpdateTransaction {
	transaction.pb.ExpirationTime = timeToProto(expiration)
	return transaction
}

// SetContractMemo sets the memo associated with the contract (max 100 bytes)
func (transaction ContractUpdateTransaction) SetContractMemo(memo string) ContractUpdateTransaction {
	transaction.pb.Memo = memo
	return transaction
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (transaction ContractUpdateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ContractUpdateTransaction {
	return ContractUpdateTransaction{transaction.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), transaction.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (transaction ContractUpdateTransaction) SetTransactionMemo(memo string) ContractUpdateTransaction {
	return ContractUpdateTransaction{transaction.TransactionBuilder.SetTransactionMemo(memo), transaction.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (transaction ContractUpdateTransaction) SetTransactionValidDuration(validDuration time.Duration) ContractUpdateTransaction {
	return ContractUpdateTransaction{transaction.TransactionBuilder.SetTransactionValidDuration(validDuration), transaction.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (transaction ContractUpdateTransaction) SetTransactionID(transactionID TransactionID) ContractUpdateTransaction {
	return ContractUpdateTransaction{transaction.TransactionBuilder.SetTransactionID(transactionID), transaction.pb}
}

// SetNodeID sets the node AccountID for this Transaction.
func (transaction ContractUpdateTransaction) SetNodeID(nodeAccountID AccountID) ContractUpdateTransaction {
	return ContractUpdateTransaction{transaction.TransactionBuilder.SetNodeID(nodeAccountID), transaction.pb}
}
