package hedera

import (
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type ScheduleInfoQuery struct {
	QueryBuilder
	pb *proto.ScheduleGetInfoQuery
}

func NewScheduleInfoQuery() *ScheduleInfoQuery {
	pb := &proto.ScheduleGetInfoQuery{Header: &proto.QueryHeader{}}

	inner := newQueryBuilder(pb.Header)
	inner.pb.Query = &proto.Query_ScheduleGetInfo{ScheduleGetInfo: pb}

	return &ScheduleInfoQuery{inner, pb}
}

func (builder *ScheduleInfoQuery) SetScheduleID(id ScheduleID) *ScheduleInfoQuery {
	builder.pb.ScheduleID = id.toProto()
	return builder
}

func (builder *ScheduleInfoQuery) GetScheduleID(id ScheduleID) ScheduleID {
	return scheduleIDFromProto(builder.pb.GetScheduleID())
}

func (builder *ScheduleInfoQuery) Execute(client *Client) (ScheduleInfo, error) {
	resp, err := builder.execute(client)
	if err != nil {
		return ScheduleInfo{}, err
	}

	keyList, err := publicKeyListFromProto(resp.GetScheduleGetInfo().ScheduleInfo.GetSignatories())
	if err != nil {
		return ScheduleInfo{}, err
	}

	adminKey, err := publicKeyFromProto(resp.GetScheduleGetInfo().ScheduleInfo.GetAdminKey())
	if err != nil {
		return ScheduleInfo{}, err
	}

	return ScheduleInfo{
		ScheduleID:       scheduleIDFromProto(resp.GetScheduleGetInfo().ScheduleInfo.GetScheduleID()),
		CreatorAccountID: accountIDFromProto(resp.GetScheduleGetInfo().ScheduleInfo.GetCreatorAccountID()),
		PayerAccountID:   accountIDFromProto(resp.GetScheduleGetInfo().ScheduleInfo.GetPayerAccountID()),
		TransactionBody:  resp.GetScheduleGetInfo().ScheduleInfo.GetTransactionBody(),
		Signers:          keyList,
		AdminKey:         adminKey,
	}, nil
}

func (builder *ScheduleInfoQuery) Cost(client *Client) (Hbar, error) {
	// deleted files return a COST_ANSWER of zero which triggers `INSUFFICIENT_TX_FEE`
	// if you set that as the query payment; 25 tinybar seems to be enough to get
	// `FILE_DELETED` back instead.
	cost, err := builder.QueryBuilder.GetCost(client)
	if err != nil {
		return ZeroHbar, err
	}

	// math.Min requires float64 and returns float64
	if cost.AsTinybar() > 25 {
		return cost, nil
	}

	return HbarFromTinybar(25), nil
}

//
// The following _3_ must be copy-pasted at the bottom of **every** _query.go file
// We override the embedded fluent setter methods to return the outer type
//

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (builder *ScheduleInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *ScheduleInfoQuery {
	return &ScheduleInfoQuery{*builder.QueryBuilder.SetMaxQueryPayment(maxPayment), builder.pb}
}

// SetQueryPayment sets the payment amount for this Query.
func (builder *ScheduleInfoQuery) SetQueryPayment(paymentAmount Hbar) *ScheduleInfoQuery {
	return &ScheduleInfoQuery{*builder.QueryBuilder.SetQueryPayment(paymentAmount), builder.pb}
}

// SetQueryPaymentTransaction sets the payment Transaction for this Query.
func (builder *ScheduleInfoQuery) SetQueryPaymentTransaction(tx Transaction) *ScheduleInfoQuery {
	return &ScheduleInfoQuery{*builder.QueryBuilder.SetQueryPaymentTransaction(tx), builder.pb}
}
