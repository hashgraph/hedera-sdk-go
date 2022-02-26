package hedera

import (
	"context"
	"math"
	"time"

	"github.com/ethereum/go-ethereum/log"

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

const maxAttempts = 10

type _ExecutionState uint32

const (
	executionStateRetry    _ExecutionState = 0
	executionStateFinished _ExecutionState = 1
	executionStateError    _ExecutionState = 2
	executionStateExpired  _ExecutionState = 3
)

type _Method struct {
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

type _Response struct {
	query       *services.Response
	transaction *services.TransactionResponse
}

type _IntermediateResponse struct {
	query       *services.Response
	transaction TransactionResponse
}

type _ProtoRequest struct {
	query       *services.Query
	transaction *services.Transaction
}

type QueryHeader struct {
	header *services.QueryHeader //nolint
}

type _Request struct {
	query       *Query
	transaction *Transaction
}

func _Execute(
	client *Client,
	request _Request,
	shouldRetry func(string, _Request, _Response) _ExecutionState,
	makeRequest func(request _Request) _ProtoRequest,
	advanceRequest func(_Request),
	getNodeAccountID func(_Request) AccountID,
	getMethod func(_Request, *_Channel) _Method,
	mapStatusError func(_Request, _Response) error,
	mapResponse func(_Request, _Response, AccountID, _ProtoRequest) (_IntermediateResponse, error),
	logID string,
) (_IntermediateResponse, error) {
	var maxAttempts int
	var minBackoff *time.Duration
	var maxBackoff *time.Duration

	if client.maxAttempts != nil {
		maxAttempts = *client.maxAttempts
	} else {
		if request.query != nil {
			maxAttempts = request.query.maxRetry
		} else {
			maxAttempts = request.transaction.maxRetry
		}
	}

	if request.query != nil {
		if request.query.maxBackoff == nil {
			maxBackoff = &client.maxBackoff
		} else {
			maxBackoff = request.query.maxBackoff
		}
		if request.query.minBackoff == nil {
			minBackoff = &client.minBackoff
		} else {
			minBackoff = request.query.minBackoff
		}
	} else {
		if request.transaction.maxBackoff == nil {
			maxBackoff = &client.maxBackoff
		} else {
			maxBackoff = request.transaction.maxBackoff
		}
		if request.transaction.minBackoff == nil {
			minBackoff = &client.minBackoff
		} else {
			minBackoff = request.transaction.minBackoff
		}
	}

	var attempt int64
	var errPersistent error

	for attempt = int64(0); attempt < int64(maxAttempts); attempt++ {
		protoRequest := makeRequest(request)
		nodeAccountID := getNodeAccountID(request)

		node, ok := client.network._GetNodeForAccountID(nodeAccountID)
		if !ok {
			return _IntermediateResponse{}, ErrInvalidNodeAccountIDSet{nodeAccountID}
		}

		node._InUse()

		log.Trace("[%s] [Node] AccountID: %s, IP: %s", logID, node.accountID.String(), node.address._String())

		if !node._IsHealthy() {
			log.Trace("[%s] [Node] Is unhealthy waiting %s before continuing", logID, node._Wait().String())
			delay := node._Wait()
			time.Sleep(delay)
		}

		log.Trace("[%s] [Node] Updating node account ID index", logID)
		advanceRequest(request)

		channel, err := node._GetChannel()
		if err != nil {
			node._IncreaseDelay()
			continue
		}

		method := getMethod(request, channel)

		resp := _Response{}

		log.Trace("[%s] Executing gRPC call", logID)
		if method.query != nil {
			resp.query, err = method.query(context.TODO(), protoRequest.query)
		} else {
			resp.transaction, err = method.transaction(context.TODO(), protoRequest.transaction)
		}

		if err != nil {
			errPersistent = err
			if _ExecutableDefaultRetryHandler(logID, err) {
				node._IncreaseDelay()
				continue
			}
			if errPersistent == nil {
				errPersistent = errors.New("error")
			}
			return _IntermediateResponse{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
		}

		node._DecreaseDelay()

		retry := shouldRetry(logID, request, resp)

		switch retry {
		case executionStateRetry:
			errPersistent = mapStatusError(request, resp)
			_DelayForAttempt(logID, minBackoff, maxBackoff, attempt)
			continue
		case executionStateExpired:
			if !client.GetOperatorAccountID()._IsZero() && request.transaction.regenerateTransactionID && !request.transaction.transactionIDs.locked {
				log.Trace("[%s] Received `TRANSACTION_EXPIRED` with transaction ID regeneration enabled; regenerating", logID)
				_, err = request.transaction.transactionIDs._Set(request.transaction.nextTransactionIndex, TransactionIDGenerate(client.GetOperatorAccountID()))
				if err != nil {
					panic(err)
				}
				continue
			} else {
				return _IntermediateResponse{}, mapStatusError(request, resp)
			}
		case executionStateError:
			return _IntermediateResponse{}, mapStatusError(request, resp)
		case executionStateFinished:
			return mapResponse(request, resp, node.accountID, protoRequest)
		}
	}

	if errPersistent == nil {
		errPersistent = errors.New("error")
	}

	return _IntermediateResponse{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
}

func _DelayForAttempt(logID string, minBackoff *time.Duration, maxBackoff *time.Duration, attempt int64) {
	// 0.1s, 0.2s, 0.4s, 0.8s, ...
	ms := int64(math.Min(float64(minBackoff.Milliseconds())*math.Pow(2, float64(attempt)), float64(maxBackoff.Milliseconds())))
	log.Trace("[%s] Waiting %sms before retrying attempt %d", logID, ms, attempt+1)
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func _ExecutableDefaultRetryHandler(logID string, err error) bool {
	code := status.Code(err)
	log.Trace("[%s] Received gRPC error with status code %s", logID, code.String())

	switch code {
	case codes.ResourceExhausted, codes.Unavailable:
		return true
	case codes.Internal:
		grpcErr, ok := status.FromError(err)

		if !ok {
			return false
		}

		return rstStream.Match([]byte(grpcErr.Message()))
	default:
		return false
	}
}
