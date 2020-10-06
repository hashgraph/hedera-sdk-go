package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// AccountCreateTransaction creates a new account. After the account is created, the AccountID for it is in the receipt,
// or by asking for a Record of the transaction to be created, and retrieving that. The account can then automatically
// generate records for large transfers into it or out of it, which each last for 25 hours. Records are generated for
// any transfer that exceeds the thresholds given here. This account is charged hbar for each record generated, so the
// thresholds are useful for limiting Record generation to happen only for large transactions.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0,
// with a null key. Future versions of the API will support multiple realms and multiple shards.
type AccountDeleteTransaction struct {
	Transaction
	pb *proto.CryptoDeleteTransactionBody
}

func NewAccountDeleteTransaction() *AccountDeleteTransaction {
	pb := &proto.CryptoDeleteTransactionBody{}

	transaction := AccountDeleteTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}
	return &transaction
}

func (transaction *AccountDeleteTransaction ) SetAccountId(accountId AccountID) *AccountDeleteTransaction {
	transaction.pb.DeleteAccountID = accountId.toProtobuf()
	return transaction
}

func (transaction *AccountDeleteTransaction ) GetAccountId() AccountID {
	return accountIDFromProto(transaction.pb.GetDeleteAccountID())
}

func (transaction *AccountDeleteTransaction ) SetTransferAccountId(transferAccountId AccountID) *AccountDeleteTransaction {
	transaction.pb.TransferAccountID = transferAccountId.toProtobuf()
	return transaction
}

func (transaction *AccountDeleteTransaction ) GetTransferAccountId(transferAccountId AccountID) AccountID {
	return accountIDFromProto(transaction.pb.GetTransferAccountID())
}



