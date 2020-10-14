package hedera

import (
	"fmt"
	protobuf "github.com/golang/protobuf/proto"
	"github.com/hashgraph/hedera-sdk-go/proto"
)

type TransactionReceiptQuery struct {
	Query
	pb *proto.TransactionGetReceiptQuery
}

func NewTransactionReceiptQuery() *TransactionReceiptQuery {
	header := proto.QueryHeader{}
	query := newQuery(false, &header)
	pb := &proto.TransactionGetReceiptQuery{Header: &header}

	return &TransactionReceiptQuery{
		Query: query,
		pb:    pb,
	}
}

func transactionReceiptQuery_shouldRetry(status Status, response response) bool {
	if status == StatusBusy {
		return true
	}

	fmt.Printf("%+v\n", response.query)

	status = Status(response.query.GetTransactionGetReceipt().Receipt.Status)

	switch status {
	case StatusBusy:
	case StatusUnknown:
	case StatusOk:
	case StatusReceiptNotFound:
		return true
	default:
		return false
	}

	return false
}

func transactionReceiptQuery_mapResponseStatus(_ request, response response) Status {
	return Status(response.query.GetTransactionGetReceipt().Header.NodeTransactionPrecheckCode)
}

func transactionReceiptQuery_getMethod(channel *channel) method {
	return method{
		query: channel.getCrypto().GetTransactionReceipts,
	}
}

func (query *TransactionReceiptQuery) SetTransactionID(transactionID TransactionID) *TransactionReceiptQuery {
	query.pb.TransactionID = transactionID.toProtobuf()
	return query
}

func (query *TransactionReceiptQuery) SetNodeId(accountID AccountID) *TransactionReceiptQuery {
	query.paymentTransactionNodeIDs = make([]AccountID, 0)
	query.paymentTransactionNodeIDs = append(query.paymentTransactionNodeIDs, accountID)
	return query
}

func (query *TransactionReceiptQuery) GetNodeId(client *Client) AccountID {
	if query.paymentTransactionNodeIDs != nil {
		return query.paymentTransactionNodeIDs[query.nextPaymentTransactionIndex]
	}

	if query.nodeID.isZero() {
		return query.nodeID
	} else {
		return client.getNextNode()
	}
}

func (query *TransactionReceiptQuery) SetQueryPayment(queryPayment Hbar) *TransactionReceiptQuery {
	query.queryPayment = queryPayment
	return query
}

func (query *TransactionReceiptQuery) SetMaxQueryPayment(queryMaxPayment Hbar) *TransactionReceiptQuery {
	query.maxQueryPayment = queryMaxPayment
	return query
}

func (query *TransactionReceiptQuery) Execute(client *Client) (TransactionReceipt, error) {
	if client == nil || client.operator == nil {
		return TransactionReceipt{}, errNoClientProvided
	}

	query.queryPayment = NewHbar(2)
	query.paymentTransactionID = TransactionIDGenerate(client.operator.accountID)

	cost := query.queryPayment

	if len(query.paymentTransactionNodeIDs) == 0 {
		size := client.getNumberOfNodesForTransaction()
		for i := 0; i < size; i++ {
			query.paymentTransactionNodeIDs = append(query.paymentTransactionNodeIDs, client.getNextNode())
		}
	}

	for _, nodeID := range query.paymentTransactionNodeIDs {
		transaction, err := makePaymentTransaction(
			query.paymentTransactionID,
			nodeID,
			client.operator,
			cost,
		)
		if err != nil {
			return TransactionReceipt{}, err
		}

		query.paymentTransactionNodeIDs = append(query.paymentTransactionNodeIDs, nodeID)
		query.paymentTransactions = append(query.paymentTransactions, transaction)
	}

	resp, err := execute(
		client,
		request{
			query: &query.Query,
		},
		transactionReceiptQuery_shouldRetry,
		query_makeRequest,
		query_advanceRequest,
		query_getNodeId,
		transactionReceiptQuery_getMethod,
		transactionReceiptQuery_mapResponseStatus,
		query_mapResponse,
	)

	if err != nil {
		return TransactionReceipt{}, err
	}

	return transactionReceiptFromProtobuf(resp.query.GetTransactionGetReceipt().Receipt), nil
}

func makePaymentTransaction(transactionID TransactionID, nodeID AccountID, operator *operator, cost Hbar) (*proto.Transaction, error) {
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
