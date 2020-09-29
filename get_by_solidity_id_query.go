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

func (transaction *GetBySolidityIDQuery) SetSolidityID(id string) *GetBySolidityIDQuery {
	transaction.pb.SolidityID = id
	return transaction
}

func (transaction *GetBySolidityIDQuery) Execute(client *Client) (EntityID, error) {
	var id EntityID = nil

	resp, err := transaction.execute(client)
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
func (transaction *GetBySolidityIDQuery) SetMaxQueryPayment(maxPayment Hbar) *GetBySolidityIDQuery {
	return &GetBySolidityIDQuery{*transaction.QueryBuilder.SetMaxQueryPayment(maxPayment), transaction.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (transaction *GetBySolidityIDQuery) SetQueryPayment(paymentAmount Hbar) *GetBySolidityIDQuery {
	return &GetBySolidityIDQuery{*transaction.QueryBuilder.SetQueryPayment(paymentAmount), transaction.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (transaction *GetBySolidityIDQuery) SetQueryPaymentTransaction(tx Transaction) *GetBySolidityIDQuery {
	return &GetBySolidityIDQuery{*transaction.QueryBuilder.SetQueryPaymentTransaction(tx), transaction.pb}
}
