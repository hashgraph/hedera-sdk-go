package hedera

import (
	"context"
	"fmt"
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

type request struct {
	query       *proto.Query
	transaction *proto.Transaction
}

func execute(
	client *Client,
	isFrozen func() bool,
	freezeWith func(client *Client) error,
	shouldRetry func(Status) bool,
	makeRequest func() request,
	advanceRequest func(),
	getNodeId func(*Client) AccountID,
	getMethod func(*channel) method,
	mapResponseStatus func(response) Status,
	mapResponse func(response, AccountID, request) (intermediateResponse, error),
) (intermediateResponse, error) {
	println("aaaaaaaaaaaaaaaaaa")
	for attempt := 0; ; /* loop forever */ attempt++ {
		delay := time.Duration(250*int64(math.Pow(2, float64(attempt)))) * time.Millisecond
		println("[execute] delay", delay)
		request := makeRequest()
		fmt.Printf("[execute] request %+v\n", request)
		node := getNodeId(client)
		println("[execute] node", node.String())

		channel, err := client.getChannel(node)
		if err != nil {
			return intermediateResponse{}, nil
		}

		fmt.Printf("[execute] channel %+v\n", channel)

		method := getMethod(channel)

		advanceRequest()

		resp := response{}
		if method.query != nil {
			r, err := method.query(context.TODO(), request.query)
			if err != nil {
				return intermediateResponse{}, nil
			}

			resp.query = r
		} else {
			r, err := method.transaction(context.TODO(), request.transaction)
			if err != nil {
				return intermediateResponse{}, nil
			}

			resp.transaction = r
		}

		println("-----------------RESPONSE-----------------")
		fmt.Printf("%+v\n", resp)

		status := mapResponseStatus(resp)

		if shouldRetry(status) {
			time.Sleep(delay)
			continue
		}

		if status != StatusOk {
			return intermediateResponse{}, newErrHederaPreCheckStatus(TransactionID{}, status)
		}

		return mapResponse(resp, node, request)
	}
}
