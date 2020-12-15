package hedera

import (
	"context"
	"github.com/pkg/errors"
	"math"
	"time"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"google.golang.org/grpc"
)

const maxAttempts = 10

type method struct {
	query func(
		context.Context,
		*proto.Query,
		...grpc.CallOption,
	) (*proto.Response, error)
	transaction func(
		context.Context,
		*proto.Transaction,
		...grpc.CallOption,
	) (*proto.TransactionResponse, error)
}

type response struct {
	query       *proto.Response
	transaction *proto.TransactionResponse
}

type intermediateResponse struct {
	query       *proto.Response
	transaction TransactionResponse
}

type protoRequest struct {
	query       *proto.Query
	transaction *proto.Transaction
}

type QueryHeader struct {
	header *proto.QueryHeader
}

type request struct {
	query       *Query
	transaction *Transaction
}

func execute(
	client *Client,
	request request,
	shouldRetry func(Status, response) bool,
	makeRequest func(request) protoRequest,
	advanceRequest func(request),
	getNodeAccountID func(request) AccountID,
	getMethod func(request, *channel) method,
	mapResponseStatus func(request, response) Status,
	mapResponse func(request, response, AccountID, protoRequest) (intermediateResponse, error),
) (intermediateResponse, error) {
	maxAttempts := 10
	var attempt int64
	var errPersistent error
	var stat Status
	if request.query != nil {
		maxAttempts = request.query.maxRetry
	} else {
		maxAttempts = request.transaction.maxRetry
	}

	for attempt := int64(0); attempt < int64(maxAttempts); attempt++ {
		protoRequest := makeRequest(request)
		nodeAccountID := getNodeAccountID(request)
		node, ok := client.network.networkNodes[nodeAccountID]
		//grpcErr, _ := status.FromError(errPersistent)
		//if grpcErr.Code() == codes.Unavailable {
		//	println("in trans grpc error")
		//	fmt.Printf("%+v\n", grpcErr.Message())
		//	time.Sleep(60 * time.Second)
		//	errPersistent = nil
		//}
		if !ok {
			return intermediateResponse{}, ErrInvalidNodeAccountIDSet{nodeAccountID}
		}

		channel, err := node.getChannel()
		if err != nil {
			return intermediateResponse{}, err
		}

		method := getMethod(request, channel)

		advanceRequest(request)

		resp := response{}

		if !node.isHealthy() {
			node.wait()
		}

		if method.query != nil {
			r, err := method.query(context.TODO(), protoRequest.query)
			if err != nil {
				node.increaseDelay()
				errPersistent = err
				continue
			}

			resp.query = r
		} else {
			r, err := method.transaction(context.TODO(), protoRequest.transaction)
			if err != nil {
				errPersistent = err
				//grpcErr, _ := status.FromError(err)
				//if grpcErr.Code() == codes.Unavailable{
				//	println("in trans grpc error")
				//	fmt.Printf("%+v\n", grpcErr.Message())
				//	node.increaseDelay()
				//	delayForAttempt(attempt)
				//	continue
				//}
				node.increaseDelay()
			}

			resp.transaction = r
		}

		node.decreaseDelay()

		stat = mapResponseStatus(request, resp)

		if shouldRetry(stat, resp) && attempt <= int64(maxAttempts) {
			delayForAttempt(attempt)
			continue
		}

		if stat != StatusOk && stat != StatusSuccess {
			if request.query != nil {
				return intermediateResponse{}, newErrHederaPreCheckStatus(TransactionID{}, stat)
			} else {
				return intermediateResponse{}, newErrHederaPreCheckStatus(request.transaction.GetTransactionID(), stat)
			}
		}

		return mapResponse(request, resp, node.accountID, protoRequest)
	}

	return intermediateResponse{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
}

func delayForAttempt(attempt int64) {
	// 0.1s, 0.2s, 0.4s, 0.8s, ...
	ms := int64(math.Floor(50 * math.Pow(2, float64(attempt))))
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
