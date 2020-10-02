package hedera

import (
	"context"
	"math"
	"time"

	"github.com/hashgraph/hedera-sdk-go/proto"
	"google.golang.org/grpc"
)

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
	// query *proto.Response
	transaction TransactionResponse
}

type protoRequest struct {
	query       *proto.Query
	transaction *proto.Transaction
}

type request struct {
    // query Query
    transaction *Transaction
}

func execute(
	client *Client,
    request request,
	shouldRetry func(request, Status) bool,
	makeRequest func(request) protoRequest,
	advanceRequest func(request),
	getNodeId func(request, *Client) AccountID,
	getMethod func(*channel) method,
	mapResponseStatus func(request, response) Status,
	mapResponse func(request, response, AccountID, protoRequest) (intermediateResponse, error),
) (intermediateResponse, error) {
	for attempt := 0; ; /* loop forever */ attempt++ {
		delay := time.Duration(250*int64(math.Pow(2, float64(attempt)))) * time.Millisecond
		protoRequest := makeRequest(request)
		node := getNodeId(request, client)

		channel, err := client.getChannel(node)
		if err != nil {
			return intermediateResponse{}, nil
		}

		method := getMethod(channel)

		advanceRequest(request)

		resp := response{}
		if method.query != nil {
			r, err := method.query(context.TODO(), protoRequest.query)
			if err != nil {
				return intermediateResponse{}, nil
			}

			resp.query = r
		} else {
			r, err := method.transaction(context.TODO(), protoRequest.transaction)
			if err != nil {
				return intermediateResponse{}, nil
			}

			resp.transaction = r
		}

		status := mapResponseStatus(request, resp)

		if shouldRetry(request, status) {
			time.Sleep(delay)
			continue
		}

		if status != StatusOk {
			return intermediateResponse{}, newErrHederaPreCheckStatus(TransactionID{}, status)
		}

		return mapResponse(request, resp, node, protoRequest)
	}
}
