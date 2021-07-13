package hedera

import (
	"context"
	"github.com/hashgraph/hedera-protobufs-go/services"
	"math"
	"time"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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
		*services.Query,
		...grpc.CallOption,
	) (*services.Response, error)
	transaction func(
		context.Context,
		*services.Transaction,
		...grpc.CallOption,
	) (*services.TransactionResponse, error)
}

type response struct {
	query       *services.Response
	transaction *services.TransactionResponse
}

type intermediateResponse struct {
	query       *services.Response
	transaction TransactionResponse
}

type protoRequest struct {
	query       *services.Query
	transaction *services.Transaction
}

type QueryHeader struct {
	header *services.QueryHeader
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
			if executableDefaultRetryHandler(err) {
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

func executableDefaultRetryHandler(err error) bool {
	code := status.Code(err)

	switch code {
	case codes.ResourceExhausted, codes.Unavailable:
		return true
	case codes.Internal:
		grpcErr, ok := status.FromError(err)

		if !ok {
			return false
		}

		return RST_STREAM.FindIndex([]byte(grpcErr.Message())) != nil
	default:
		return false
	}
}
