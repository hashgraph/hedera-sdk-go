package hedera

import (
	"log"
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ContractCreateTransaction struct {
	Transaction
	pb *proto.ContractCreateTransactionBody
}

func NewContractCreateTransaction() *ContractCreateTransaction {
	pb := &proto.ContractCreateTransactionBody{}

	transaction := ContractCreateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetAutoRenewPeriod(7890000 * time.Second)

	return &transaction
}

func (transaction *ContractCreateTransaction) SetAdminKey(adminKey Key) *ContractCreateTransaction {
	transaction.pb.AdminKey = adminKey.toProtoKey()
	return transaction
}

func (transaction *ContractCreateTransaction) GetAdminKey()  Key{
	var key, err = publicKeyFromProto(transaction.pb.GetAdminKey())
	if err != nil {
		log.Fatal(err)
	}

	return key
}

func (transaction *ContractCreateTransaction) SetGas(gas int64) *ContractCreateTransaction {
	transaction.pb.Gas = gas
	return transaction
}

func (transaction *ContractCreateTransaction) GetGas() int64 {
	return transaction.pb.GetGas()
}

// SetInitialBalance sets the initial number of Hbar to put into the account
func (transaction *ContractCreateTransaction) SetInitialBalance(initialBalance Hbar) *ContractCreateTransaction {
	transaction.pb.InitialBalance = initialBalance.AsTinybar()
	return transaction
}

// GetInitialBalance gets the initial number of Hbar in the account
func (transaction *ContractCreateTransaction) GetInitialBalance() Hbar {
	return HbarFromTinybar(transaction.pb.GetInitialBalance())
}

// SetAutoRenewPeriod sets the time duration for when account is charged to extend its expiration date. When the account
// is created, the payer account is charged enough hbars so that the new account will not expire for the next
// auto renew period. When it reaches the expiration time, the new account will then be automatically charged to
// renew for another auto renew period. If it does not have enough hbars to renew for that long, then the  remaining
// hbars are used to extend its expiration as long as possible. If it is has a zero balance when it expires,
// then it is deleted.
func (transaction *ContractCreateTransaction) SetAutoRenewPeriod(autoRenewPeriod time.Duration) *ContractCreateTransaction {
	transaction.pb.AutoRenewPeriod = durationToProto(autoRenewPeriod)
	return transaction
}

func (transaction *ContractCreateTransaction) GetAutoRenewPeriod() time.Duration {
	return durationFromProto(transaction.pb.GetAutoRenewPeriod())
}

// SetProxyAccountID sets the ID of the account to which this account is proxy staked. If proxyAccountID is not set,
// is an invalid account, or is an account that isn't a node, then this account is automatically proxy staked to a node
// chosen by the network, but without earning payments. If the proxyAccountID account refuses to accept proxy staking ,
// or if it is not currently running a node, then it will behave as if proxyAccountID was not set.
func (transaction *ContractCreateTransaction) SetProxyAccountID(id AccountID) *ContractCreateTransaction {
	transaction.pb.ProxyAccountID = id.toProtobuf()
	return transaction
}

func (transaction *ContractCreateTransaction) GetProxyAccountID() AccountID {
	return accountIDFromProto(transaction.pb.ProxyAccountID)
}

func (transaction *ContractCreateTransaction) SetConstructorParameters(params []byte) *ContractCreateTransaction {
	transaction.pb.ConstructorParameters = params
	return transaction
}

func (transaction *ContractCreateTransaction) GetConstructorParameters() []byte {
	return transaction.pb.GetConstructorParameters()
}

func (transaction *ContractCreateTransaction) SetContractMemo(memo string) *ContractCreateTransaction {
	transaction.pb.Memo = memo
	return transaction
}

func (transaction *ContractCreateTransaction) GetContractMemo() string {
	return transaction.pb.GetMemo()
}
