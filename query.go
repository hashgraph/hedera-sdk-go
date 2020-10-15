package hedera

import (
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

// Transaction contains the protobuf of a prepared transaction which can be signed and executed.
type Query struct {
	pb       *proto.Query
	pbHeader *proto.QueryHeader

	paymentTransactionID        TransactionID
	nodeID                      AccountID
	maxQueryPayment             Hbar
	queryPayment                Hbar
	nextPaymentTransactionIndex int
	nextTransactionIndex        int

	paymentTransactionNodeIDs []AccountID
	paymentTransactions       []*proto.Transaction

	isPaymentRequired bool
}

func newQuery(isPaymentRequired bool, queryHeader *proto.QueryHeader) Query {
	return Query{
		pb:                        &proto.Query{},
		pbHeader:                  queryHeader,
		paymentTransactionID:      TransactionID{},
		nextTransactionIndex:      0,
		paymentTransactions:       make([]*proto.Transaction, 0),
		paymentTransactionNodeIDs: make([]AccountID, 0),
		isPaymentRequired:         isPaymentRequired,
	}
}

func (query *Query) SetNodeId(accountID AccountID) *Query {
	query.nodeID = accountID
	return query
}

func (query *Query) GetNodeId() AccountID {
	return query.nodeID
}

func query_getNodeId(request request, client *Client) AccountID {
	if len(request.query.paymentTransactionNodeIDs) > 0 {
		return request.query.paymentTransactionNodeIDs[request.query.nextPaymentTransactionIndex]
	}

	if request.query.nodeID.isZero() {
		return request.query.nodeID
	} else {
		return client.getNextNode()
	}
}

func (query *Query) SetQueryPayment(queryPayment Hbar) *Query {
	query.queryPayment = queryPayment
	return query
}

func (query *Query) SetMaxQueryPayment(queryMaxPayment Hbar) *Query {
	query.maxQueryPayment = queryMaxPayment
	return query
}

func (query *Query) IsPaymentRequired() bool {
	return true
}

func query_makeRequest(request request) protoRequest {
	if request.query.isPaymentRequired && len(request.query.paymentTransactions) > 0 {
		request.query.pbHeader.Payment = request.query.paymentTransactions[request.query.nextPaymentTransactionIndex]
	}
	request.query.pbHeader.ResponseType = proto.ResponseType_ANSWER_ONLY
	return protoRequest{
		query: request.query.pb,
	}
}

func query_advanceRequest(request request) {
	if request.query.isPaymentRequired && len(request.query.paymentTransactions) > 0 {
		request.query.nextPaymentTransactionIndex++
	}
}

func query_mapResponse(request request, response response, _ AccountID, protoRequest protoRequest) (intermediateResponse, error) {
	return intermediateResponse{
		query: response.query,
		//transaction: TransactionResponse{
		//	TransactionID: request.transaction.id,
		//	NodeID:        request.transaction.nodeIDs[request.transaction.nextTransactionIndex],
		//},
	}, nil
}

//func query_mapResponseHeader(request request, response response) (protoResponseHeader, error) {
//	return protoResponseHeader{
//		responseHeader: proto.ResponseHeader{
//			NodeTransactionPrecheckCode: response.transaction.NodeTransactionPrecheckCode,
//			ResponseType:                request.query.pbHeader.ResponseType,
//			Cost:                        response.transaction.Cost,
//		},
//	}, nil
//}
//
//func query_mapRequestHeader(request request, response response) (QueryHeader, error) {
//	return QueryHeader{
//		header: &proto.QueryHeader{
//			Payment:      request.transaction.,
//			ResponseType: 0,
//		},
//	}, nil
//}

func query_makePaymentTransaction(transactionID TransactionID, nodeID AccountID, operator *operator, cost Hbar) (*proto.Transaction, error) {
	accountAmounts := make([]*proto.AccountAmount, 0)
	accountAmounts = append(accountAmounts, &proto.AccountAmount{
		AccountID: nodeID.toProtobuf(),
		Amount:    cost.tinybar,
	})
	accountAmounts = append(accountAmounts, &proto.AccountAmount{
		AccountID: nodeID.toProtobuf(),
		Amount:    -cost.tinybar,
	})

	body := proto.TransactionBody{
		TransactionID:  transactionID.toProtobuf(),
		NodeAccountID:  nodeID.toProtobuf(),
		TransactionFee: uint64(NewHbar(1).tinybar),
		TransactionValidDuration: &proto.Duration{
			Seconds: 120,
		},
		Data: &proto.TransactionBody_CryptoTransfer{
			CryptoTransfer: &proto.CryptoTransferTransactionBody{
				Transfers: &proto.TransferList{
					AccountAmounts: accountAmounts,
				},
			},
		},
	}

	bodyBytes, err := protobuf.Marshal(&body)
	if err != nil {
		return nil, err
	}

	signature := operator.signer(bodyBytes)
	sigPairs := make([]*proto.SignaturePair, 0)
	sigPairs = append(sigPairs, operator.publicKey.toSignaturePairProtobuf(signature))

	return &proto.Transaction{
		BodyBytes: bodyBytes,
		SigMap: &proto.SignatureMap{
			SigPair: sigPairs,
		},
	}, nil
}
