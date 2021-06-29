package hedera

import (
	"context"
	"math"
	"time"

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"google.golang.org/grpc"
)

const maxAttempts = 10

type executionState uint32

const (
	executionStateRetry    executionState = 0
	executionStateFinished executionState = 1
	executionStateError    executionState = 2
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
	shouldRetry func(request, response) executionState,
	makeRequest func(request) protoRequest,
	advanceRequest func(request),
	getNodeAccountID func(request) AccountID,
	getMethod func(request, *channel) method,
	mapStatusError func(request, response, *NetworkName) error,
	mapResponse func(request, response, AccountID, protoRequest) (intermediateResponse, error),
) (intermediateResponse, error) {
	maxAttempts := 10
	var attempt int64
	var errPersistent error

	if request.query != nil {
		maxAttempts = request.query.maxRetry
	} else {
		maxAttempts = request.transaction.maxRetry
	}

	for attempt = int64(0); attempt < int64(maxAttempts); attempt++ {
		protoRequest := makeRequest(request)
		nodeAccountID := getNodeAccountID(request)

		node, ok := client.network.networkNodes[nodeAccountID]
		if !ok {
			return intermediateResponse{}, ErrInvalidNodeAccountIDSet{nodeAccountID}
		}

		node.inUse()

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
			resp.query, err = method.query(context.TODO(), protoRequest.query)
		} else {
			resp.transaction, err = method.transaction(context.TODO(), protoRequest.transaction)
		}

		if err != nil {
			errPersistent = err
			if defaultRetryHandler(err) {
				node.increaseDelay()
				continue
			}
			return intermediateResponse{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
		}

		node.decreaseDelay()

		retry := shouldRetry(request, resp)

		switch retry {
		case executionStateRetry:
			if attempt <= int64(maxAttempts) {
				delayForAttempt(attempt)
				continue
			} else {
				errPersistent = mapStatusError(request, resp, client.networkName)
				break
			}
		case executionStateError:
			return intermediateResponse{}, mapStatusError(request, resp, client.networkName)
		case executionStateFinished:
			return mapResponse(request, resp, node.accountID, protoRequest)
		}
	}

	return intermediateResponse{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
}

func delayForAttempt(attempt int64) {
	// 0.1s, 0.2s, 0.4s, 0.8s, ...
	ms := int64(math.Floor(50 * math.Pow(2, float64(attempt))))
	time.Sleep(time.Duration(ms) * time.Millisecond)
}
