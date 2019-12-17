package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type GetBySolidityIDQuery struct {
	QueryBuilder
	pb *proto.GetBySolidityIDQuery
}

func NewGetBySolidityIDQuery() *GetBySolidityIDQuery {
	pb := &proto.GetBySolidityIDQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_GetBySolidityID{pb}

	return &GetBySolidityIDQuery{inner, pb}
}

func (builder *GetBySolidityIDQuery) SetSolidityID(id string) *GetBySolidityIDQuery {
	builder.pb.SolidityID = id
	return builder
}

func (builder *GetBySolidityIDQuery) Execute(client *Client) (EntityID, error) {
	var id = EntityID{}

	resp, err := builder.execute(client)
	if err != nil {
		return id, err
	}

	if resp.GetGetBySolidityID().GetAccountID() != nil {
		id = EntityID{
			ty: "ACCOUNT",
			id: accountIDFromProto(resp.GetGetBySolidityID().GetAccountID()),
		}
	} else if resp.GetGetBySolidityID().GetFileID() != nil {
		id = EntityID{
			ty: "FILE",
			id: fileIDFromProto(resp.GetGetBySolidityID().GetFileID()),
		}
	} else if resp.GetGetBySolidityID().GetContractID() != nil {
		id = EntityID{
			ty: "CONTRACT",
			id: contractIDFromProto(resp.GetGetBySolidityID().GetContractID()),
		}
	}

	return id, nil
}
