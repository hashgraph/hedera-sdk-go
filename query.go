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
		maxQueryPayment:           NewHbar(0),
		queryPayment:              NewHbar(0),
	}
}

func (query *Query) SetNodeAccountID(accountID AccountID) *Query {
	query.nodeID = accountID
	return query
}

func (query *Query) GetNodeAccountId() AccountID {
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

func query_shouldRetry(status Status, _ response) bool {
	return status == StatusBusy
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
		request.query.nextPaymentTransactionIndex = (request.query.nextPaymentTransactionIndex + 1) % len(request.query.paymentTransactions)
	}
}

func query_mapResponse(request request, response response, _ AccountID, protoRequest protoRequest) (intermediateResponse, error) {
	return intermediateResponse{
		query: response.query,
	}, nil
}

func query_generatePayments(query *Query, client *Client, cost Hbar) error {
	if len(query.paymentTransactionNodeIDs) == 0 {
		size := client.getNumberOfNodesForTransaction()
		for i := 0; i < size; i++ {
			query.paymentTransactionNodeIDs = append(query.paymentTransactionNodeIDs, client.getNextNode())
		}
	}

	for _, nodeID := range query.paymentTransactionNodeIDs {
		transaction, err := query_makePaymentTransaction(
			query.paymentTransactionID,
			nodeID,
			client.operator,
			cost,
		)
		if err != nil {
			return err
		}

		query.paymentTransactionNodeIDs = append(query.paymentTransactionNodeIDs, nodeID)
		query.paymentTransactions = append(query.paymentTransactions, transaction)
	}

	return nil
}

func query_makePaymentTransaction(transactionID TransactionID, nodeID AccountID, operator *operator, cost Hbar) (*proto.Transaction, error) {
	accountAmounts := make([]*proto.AccountAmount, 0)
	accountAmounts = append(accountAmounts, &proto.AccountAmount{
		AccountID: nodeID.toProtobuf(),
		Amount:    cost.tinybar,
	})
	accountAmounts = append(accountAmounts, &proto.AccountAmount{
		AccountID: operator.accountID.toProtobuf(),
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
