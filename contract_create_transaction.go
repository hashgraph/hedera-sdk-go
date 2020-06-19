package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// ContractCreateTransaction starts a new smart contract instance. After the instance is created, the ContractID for it
// is in the receipt, or can be retrieved with a GetByKey query, or by asking for a Record of the transaction to be
// created, and retrieving that. The instance will run the bytecode stored in the given file, referenced either by
// FileID or by the TransactionID of the transaction that created the file. The constructor will be executed using the
// given amount of gas, and any unspent gas will be refunded to the paying account. Constructor inputs come from the
// given constructorParameters.
//
// A smart contract instance normally enforces rules, so "the code is law". For example, an ERC-20 contract prevents a
// transfer from being undone without a signature by the recipient of the transfer. This is always enforced if the
// contract instance was created with the adminKeys being null. But for some uses, it might be desirable to create
// something like an ERC-20 contract that has a specific group of trusted individuals who can act as a "supreme court"
// with the ability to override the normal operation, when a sufficient number of them agree to do so. If adminKeys is
// not null, then they can sign a transaction that can change the state of the smart contract in arbitrary ways, such as
// to reverse a transaction that violates some standard of behavior that is not covered by the code itself. The admin
// keys can also be used to change the autoRenewPeriod, and change the adminKeys field itself. The API currently does
// not implement this ability. But it does allow the adminKeys field to be set and queried, and will in the future
// implement such admin abilities for any instance that has a non-null adminKeys.
//
// An entity (account, file, or smart contract instance) must be created in a particular realm. If the realmID is left
// null, then a new realm will be created with the given admin key. If a new realm has a null adminKey, then anyone can
// create/modify/delete entities in that realm. But if an admin key is given, then any transaction to
// create/modify/delete an entity in that realm must be signed by that key, though anyone can still call functions on
// smart contract instances that exist in that realm. A realm ceases to exist when everything within it has expired and
// no longer exists.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0, with
// a null key. Future versions of the API will support multiple realms and multiple shards.
type ContractCreateTransaction struct {
	TransactionBuilder
	pb *proto.ContractCreateTransactionBody
}

// NewContractCreateTransaction creates a ContractCreateTransaction builder which can be
// used to construct and execute a Contract Create Transaction. This constructor defaults
// the autoRenewPeriod to ~1/4 year to fall within the required range, if desired the value
// can be changed through the SetAutoRenewPeriod method
func NewContractCreateTransaction() ContractCreateTransaction {
	pb := &proto.ContractCreateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_ContractCreateInstance{ContractCreateInstance: pb}

	builder := ContractCreateTransaction{inner, pb}

	// Default autoRenewPeriod to a value within the required range (~1/4 year)
	return builder.SetAutoRenewPeriod(131500 * time.Minute)
}

// SetBytecodeFileID sets the ID of the file containing the smart contract byte code. A copy will be made and held by
// the contract instance, and have the same expiration time as the instance.
func (builder ContractCreateTransaction) SetBytecodeFileID(id FileID) ContractCreateTransaction {
	builder.pb.FileID = id.toProto()
	return builder
}

// SetAdminKey sets the key required to arbitrarily modify the state of the instance and its fields. If this is left
// unset, then such modifications are not possible, and there is no administrator that can override the normal operation
// of the smart contract instance. Note that if it is created with no admin keys, then there is no administrator to
// authorize changing the admin keys, so there can never be any admin keys for that instance.
func (builder ContractCreateTransaction) SetAdminKey(publicKey Ed25519PublicKey) ContractCreateTransaction {
	builder.pb.AdminKey = publicKey.toProto()
	return builder
}

// SetContractMemo sets the optional memo field which can contain a string whose length is up to 100 bytes. That is the
// size after Unicode NFD then UTF-8 conversion. This field can be used to describe the smart contract. It could also be
// used for other purposes.
//
//One recommended purpose is to hold a hexadecimal string that is the SHA-384 hash of a PDF file containing
// a human-readable legal contract. Then, if the admin keys are the public keys of human arbitrators, they can use that
// legal document to guide their decisions during a binding arbitration tribunal, convened to consider any changes to
// the smart contract in the future. The memo field can only be changed using the admin keys. If there are no admin
// keys, then it cannot be changed after the smart contract is created.
func (builder ContractCreateTransaction) SetContractMemo(memo string) ContractCreateTransaction {
	builder.pb.Memo = memo
	return builder
}

// SetGas sets the gas required to run the constructor
func (builder ContractCreateTransaction) SetGas(gas uint64) ContractCreateTransaction {
	builder.pb.Gas = int64(gas)
	return builder
}

// SetInitialBalance sets the initial Hbar to put into the account associated with and owned by the smart contract
func (builder ContractCreateTransaction) SetInitialBalance(initialBalance Hbar) ContractCreateTransaction {
	builder.pb.InitialBalance = initialBalance.AsTinybar()
	return builder
}

// SetProxyAccountID sets the AccountID of the account to which this contract is proxy staked. If proxyAccountID is left
// unset, is an invalid account, or is an account that isn't a node, then this contract is automatically proxy staked
// to a node chosen by the network, but without earning payments. If the proxyAccountID account refuses to accept proxy
// staking , or if it is not currently running a node, then it will behave as if  proxyAccountID was never set.
func (builder ContractCreateTransaction) SetProxyAccountID(id AccountID) ContractCreateTransaction {
	builder.pb.ProxyAccountID = id.toProto()
	return builder
}

// SetAutoRenewPeriod sets the duration the instance will exist for. When that is reached, it will renew itself for
// another autoRenewPeriod duration by charging its associated account. If it has an insufficient balance to extend that
// long, it will extend as long as it can. If its balance is zero, the instance will be deleted.
func (builder ContractCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) ContractCreateTransaction {
	builder.pb.AutoRenewPeriod = durationToProto(autoRenewPeriod)
	return builder
}

// SetConstructorParams sets the ContractFunctionParams to pass to the constructor. If this constructor stores
// information, it is charged gas to store it. There is a fee in hbars to maintain that storage until the expiration
// time, and that fee is added as part of the transaction fee.
func (builder ContractCreateTransaction) SetConstructorParams(params *ContractFunctionParams) ContractCreateTransaction {
	builder.pb.ConstructorParameters = params.build(nil)
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder ContractCreateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) ContractCreateTransaction {
	return ContractCreateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder ContractCreateTransaction) SetTransactionMemo(memo string) ContractCreateTransaction {
	return ContractCreateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder ContractCreateTransaction) SetTransactionValidDuration(validDuration time.Duration) ContractCreateTransaction {
	return ContractCreateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder ContractCreateTransaction) SetTransactionID(transactionID TransactionID) ContractCreateTransaction {
	return ContractCreateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder ContractCreateTransaction) SetNodeAccountID(nodeAccountID AccountID) ContractCreateTransaction {
	return ContractCreateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
