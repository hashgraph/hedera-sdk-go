package hedera

import "github.com/hashgraph/hedera-sdk-go/proto"

type GetBySolidityIDQuery struct {
	QueryBuilder
	pb *proto.GetBySolidityIDQuery
}

func NewGetBySolidityIDQuery() *GetBySolidityIDQuery {
	pb := &proto.GetBySolidityIDQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_GetBySolidityID{GetBySolidityID: pb}

	return &GetBySolidityIDQuery{inner, pb}
}

func (builder *GetBySolidityIDQuery) SetSolidityID(id string) *GetBySolidityIDQuery {
	builder.pb.SolidityID = id
	return builder
}

func (builder *GetBySolidityIDQuery) Execute(client *Client) (EntityID, error) {
	var id EntityID = nil

	resp, err := builder.execute(client)
	if err != nil {
		return nil, err
	}

	if resp.GetGetBySolidityID().GetAccountID() != nil {
		id = accountIDFromProto(resp.GetGetBySolidityID().GetAccountID())
	} else if resp.GetGetBySolidityID().GetFileID() != nil {
		id = fileIDFromProto(resp.GetGetBySolidityID().GetFileID())
	} else if resp.GetGetBySolidityID().GetContractID() != nil {
		id = contractIDFromProto(resp.GetGetBySolidityID().GetContractID())
	}

	return id, nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (builder *GetBySolidityIDQuery) SetMaxQueryPayment(maxPayment Hbar) *GetBySolidityIDQuery {
	return &GetBySolidityIDQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (builder *GetBySolidityIDQuery) SetQueryPayment(paymentAmount Hbar) *GetBySolidityIDQuery {
	return &GetBySolidityIDQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (builder *GetBySolidityIDQuery) SetQueryPaymentTransaction(tx Transaction) *GetBySolidityIDQuery {
	return &GetBySolidityIDQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
