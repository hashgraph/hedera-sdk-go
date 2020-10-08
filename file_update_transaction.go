package hedera

import (
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
)

type FileUpdateTransaction struct {
	Transaction
	pb *proto.FileUpdateTransactionBody
}

func NewFileUpdateTransaction() *FileUpdateTransaction {
	pb := &proto.FileUpdateTransactionBody{}

	transaction := FileUpdateTransaction{
		pb:          pb,
		Transaction: newTransaction(),
	}

	return &transaction
}

func (transaction *FileUpdateTransaction) SetFileID(id FileID) *FileUpdateTransaction {
	transaction.pb.FileID = id.toProto()
	return transaction
}

func (transaction *FileUpdateTransaction) GetFileID() FileID {
	return fileIDFromProto(transaction.pb.GetFileID())
}

func (transaction *FileUpdateTransaction) SetKeys(keys ...Key) *FileUpdateTransaction {
	if transaction.pb.Keys == nil {
		transaction.pb.Keys = &proto.KeyList{Keys: []*proto.Key{}}
	}
	keyList := KeyList{keys: []*proto.Key{}}
	keyList.AddAll(keys)

	transaction.pb.Keys = keyList.toProtoKeyList()

	return transaction
}

func (transaction *FileUpdateTransaction) GetKeys() KeyList {
	return keyListFromProto(transaction.pb.Keys)
}

func (transaction *FileUpdateTransaction) SetExpirationTime(expiration time.Time) *FileUpdateTransaction {
	transaction.pb.ExpirationTime = timeToProto(expiration)
	return transaction
}

func (transaction *FileUpdateTransaction) GetExpirationTime() time.Time {
	return timeFromProto(transaction.pb.ExpirationTime)
}

func (transaction *FileUpdateTransaction) SetContents(contents []byte) *FileUpdateTransaction {
	transaction.pb.Contents = contents
	return transaction
}

func (transaction *FileUpdateTransaction) GetContents() []byte {
	return transaction.pb.Contents
}
