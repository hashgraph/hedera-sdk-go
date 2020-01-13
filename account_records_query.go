package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type AccountRecordQuery struct {
	QueryBuilder
	pb *proto.CryptoGetAccountRecordsQuery
}

func NewAccountRecordQuery() *AccountRecordQuery {
	pb := &proto.CryptoGetAccountRecordsQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_CryptoGetAccountRecords{CryptoGetAccountRecords: pb}

	return &AccountRecordQuery{inner, pb}
}

func (builder *AccountRecordQuery) SetAccountID(id AccountID) *AccountRecordQuery {
	builder.pb.AccountID = id.toProto()
	return builder
}

func (builder *AccountRecordQuery) Execute(client *Client) ([]TransactionRecord, error) {
	var records = []TransactionRecord{}

	resp, err := builder.execute(client)
	if err != nil {
		return records, err
	}

	for _, element := range resp.GetCryptoGetAccountRecords().Records {
		records = append(records, transactionRecordFromProto(element))
	}

	return records, nil
}
