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
	Transaction
	pb *proto.FileCreateTransactionBody
}

// NewFileCreateTransaction creates a FileCreateTransaction transaction which can be
// used to construct and execute a File Create Transaction.
func NewFileCreateTransaction() *FileCreateTransaction {
	pb := &proto.FileCreateTransactionBody{}

	transaction := FileCreateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	transaction.SetExpirationTime(time.Now().Add(7890000))

	return &transaction
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
func (transaction *FileCreateTransaction) SetKeys(keys ...Key) *FileCreateTransaction {
	if transaction.pb.Keys == nil {
		transaction.pb.Keys = &proto.KeyList{Keys: []*proto.Key{}}
	}
	keyList := KeyList{keys: []*proto.Key{}}
	keyList.AddAll(keys)

	transaction.pb.Keys = keyList.toProtoKeyList()

	return transaction
}

func (transaction *FileCreateTransaction) GetKeys() KeyList {
	return keyListFromProto(transaction.pb.GetKeys())
}

// SetExpirationTime sets the time at which this file should expire (unless FileUpdateTransaction is used before then to
// extend its life). The file will automatically disappear at the fileExpirationTime, unless its expiration is extended
// by another transaction before that time. If the file is deleted, then its contents will become empty and it will be
// marked as deleted until it expires, and then it will cease to exist.
func (transaction *FileCreateTransaction) SetExpirationTime(expiration time.Time) *FileCreateTransaction {
	transaction.pb.ExpirationTime = timeToProto(expiration)
	return transaction
}

func (transaction *FileCreateTransaction) GetExpirationTime() time.Time {
	return timeFromProto(transaction.pb.GetExpirationTime())
}

// SetContents sets the bytes that are the contents of the file (which can be empty). If the size of the file and other
// fields in the transaction exceed the max transaction size then FileAppendTransaction can be used to continue
// uploading the file.
func (transaction *FileCreateTransaction) SetContents(contents []byte) *FileCreateTransaction {
	transaction.pb.Contents = contents
	return transaction
}

func (transaction *FileCreateTransaction) GetContents() []byte {
	return transaction.pb.GetContents()
}

