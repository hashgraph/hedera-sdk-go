package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

// FileCreateTransaction creates a new file, containing the given contents.  It is referenced by its FileID, and does
// not have a filename, so it is important to get and hold onto the FileID. After the file is created, the FileID for
// it can be found in the receipt, or retrieved with a GetByKey query, or by asking for a Record of the transaction to
// be created, and retrieving that.
//
// See FileInfoQuery for more information about files.
//
// The current API ignores shardID, realmID, and newRealmAdminKey, and creates everything in shard 0 and realm 0, with
// a null key. Future versions of the API will support multiple realms and multiple shards.
type FileCreateTransaction struct {
	TransactionBuilder
	pb *proto.FileCreateTransactionBody
}

// NewFileCreateTransaction creates a FileCreateTransaction builder which can be
// used to construct and execute a File Create Transaction.
func NewFileCreateTransaction() FileCreateTransaction {
	pb := &proto.FileCreateTransactionBody{}

	inner := newTransactionBuilder()
	inner.pb.Data = &proto.TransactionBody_FileCreate{FileCreate: pb}

	builder := FileCreateTransaction{inner, pb}
	builder.SetExpirationTime(time.Now().Add(7890000 * time.Second))
	builder.pb.Keys = &proto.KeyList{Keys: []*proto.Key{}}

	return builder
}

// AddKey adds a key to the internal list of keys associated with the file. All of the keys on the list must sign to
// create or modify a file, but only one of them needs to sign in order to delete the file. Each of those "keys" may
// itself be threshold key containing other keys (including other threshold keys). In other words, the behavior is an
// AND for create/modify, OR for delete. This is useful for acting as a revocation server. If it is desired to have the
// behavior be AND for all 3 operations (or OR for all 3), then the list should have only a single Key, which is a
// threshold key, with N=1 for OR, N=M for AND.
//
// If a file is created without adding ANY keys, the file is immutable and ONLY the
// expirationTime of the file can be changed using FileUpdateTransaction. The file contents or its keys will not be
// mutable.
func (builder FileCreateTransaction) AddKey(publicKey PublicKey) FileCreateTransaction {
	builder.pb.Keys.Keys = append(builder.pb.Keys.Keys, publicKey.toProto())
	return builder
}

// SetExpirationTime sets the time at which this file should expire (unless FileUpdateTransaction is used before then to
// extend its life). The file will automatically disappear at the fileExpirationTime, unless its expiration is extended
// by another transaction before that time. If the file is deleted, then its contents will become empty and it will be
// marked as deleted until it expires, and then it will cease to exist.
func (builder FileCreateTransaction) SetExpirationTime(expiration time.Time) FileCreateTransaction {
	builder.pb.ExpirationTime = timeToProto(expiration)
	return builder
}

// SetContents sets the bytes that are the contents of the file (which can be empty). If the size of the file and other
// fields in the transaction exceed the max transaction size then FileAppendTransaction can be used to continue
// uploading the file.
func (builder FileCreateTransaction) SetContents(contents []byte) FileCreateTransaction {
	builder.pb.Contents = contents
	return builder
}

//
// The following _5_ must be copy-pasted at the bottom of **every** _transaction.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxTransactionFee sets the max transaction fee for this Transaction.
func (builder FileCreateTransaction) SetMaxTransactionFee(maxTransactionFee Hbar) FileCreateTransaction {
	return FileCreateTransaction{builder.TransactionBuilder.SetMaxTransactionFee(maxTransactionFee), builder.pb}
}

// SetTransactionMemo sets the memo for this Transaction.
func (builder FileCreateTransaction) SetTransactionMemo(memo string) FileCreateTransaction {
	return FileCreateTransaction{builder.TransactionBuilder.SetTransactionMemo(memo), builder.pb}
}

// SetTransactionValidDuration sets the valid duration for this Transaction.
func (builder FileCreateTransaction) SetTransactionValidDuration(validDuration time.Duration) FileCreateTransaction {
	return FileCreateTransaction{builder.TransactionBuilder.SetTransactionValidDuration(validDuration), builder.pb}
}

// SetTransactionID sets the TransactionID for this Transaction.
func (builder FileCreateTransaction) SetTransactionID(transactionID TransactionID) FileCreateTransaction {
	return FileCreateTransaction{builder.TransactionBuilder.SetTransactionID(transactionID), builder.pb}
}

// SetNodeAccountID sets the node AccountID for this Transaction.
func (builder FileCreateTransaction) SetNodeAccountID(nodeAccountID AccountID) FileCreateTransaction {
	return FileCreateTransaction{builder.TransactionBuilder.SetNodeAccountID(nodeAccountID), builder.pb}
}
