package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type GetByKeyQuery struct {
	QueryBuilder
	pb *proto.GetByKeyQuery
}

func NewGetByKeyQuery() *GetByKeyQuery {
	pb := &proto.GetByKeyQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_GetByKey{pb}

	return &GetByKeyQuery{inner, pb}
}

func (builder *GetByKeyQuery) SetKey(key Ed25519PublicKey) *GetByKeyQuery {
	builder.pb.Key = key.toProto()
	return builder
}

func (builder *GetByKeyQuery) Execute(client *Client) ([]EntityID, error) {
	var ids = []EntityID{}

	resp, err := builder.execute(client)
	if err != nil {
		return ids, err
	}

	for _, element := range resp.GetGetByKey().Entities {
		if element.GetAccountID() != nil {
			ids = append(ids, EntityID{
				ty: "ACCOUNT",
				id: accountIDFromProto(element.GetAccountID()),
			})
		} else if element.GetFileID() != nil {
			ids = append(ids, EntityID{
				ty: "FILE",
				id: fileIDFromProto(element.GetFileID()),
			})
		} else if element.GetContractID() != nil {
			ids = append(ids, EntityID{
				ty: "CONTRACT",
				id: contractIDFromProto(element.GetContractID()),
			})
		}
	}

	return ids, nil
}
