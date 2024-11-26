package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"encoding/hex"
	"strconv"
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/pkg/errors"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
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

type Executable interface {
	GetMaxBackoff() time.Duration
	GetMinBackoff() time.Duration
	GetGrpcDeadline() *time.Duration
	GetMaxRetry() int
	GetNodeAccountIDs() []AccountID
	GetLogLevel() *LogLevel

	shouldRetry(Executable, interface{}) _ExecutionState
	makeRequest() interface{}
	advanceRequest()
	getNodeAccountID() AccountID
	getMethod(*_Channel) _Method
	mapStatusError(Executable, interface{}) error
	mapResponse(interface{}, AccountID, interface{}) (interface{}, error)
	getName() string
	validateNetworkOnIDs(client *Client) error
	isTransaction() bool
	getLogger(Logger) Logger
	getTransactionIDAndMessage() (string, string)
	getLogID(Executable) string // This returns transaction creation timestamp + transaction name
}

type executable struct {
	transactionIDs *_LockableSlice
	nodeAccountIDs *_LockableSlice
	maxBackoff     *time.Duration
	minBackoff     *time.Duration
	grpcDeadline   *time.Duration
	maxRetry       int
	logLevel       *LogLevel
}

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

func (e *executable) GetMaxBackoff() time.Duration {
	if e.maxBackoff != nil {
		return *e.maxBackoff
	}

	return 8 * time.Second
}

func (e *executable) GetMinBackoff() time.Duration {
	if e.minBackoff != nil {
		return *e.minBackoff
	}

	return 250 * time.Millisecond
}

func (e *executable) SetMaxBackoff(max time.Duration) *executable {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < e.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	e.maxBackoff = &max
	return e
}

func (e *executable) SetMinBackoff(min time.Duration) *executable {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if e.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	e.minBackoff = &min
	return e
}

// GetGrpcDeadline returns the grpc deadline
func (e *executable) GetGrpcDeadline() *time.Duration {
	return e.grpcDeadline
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (e *executable) SetGrpcDeadline(deadline *time.Duration) *executable {
	e.grpcDeadline = deadline
	return e
}

// GetMaxRetry returns the max number of errors before execution will fail.
func (e *executable) GetMaxRetry() int {
	return e.maxRetry
}

func (e *executable) SetMaxRetry(max int) *executable {
	e.maxRetry = max
	return e
}

// GetNodeAccountID returns the node AccountID for this transaction.
func (e *executable) GetNodeAccountIDs() []AccountID {
	nodeAccountIDs := []AccountID{}

	for _, value := range e.nodeAccountIDs.slice {
		nodeAccountIDs = append(nodeAccountIDs, value.(AccountID))
	}

	return nodeAccountIDs
}

func (e *executable) SetNodeAccountIDs(nodeAccountIDs []AccountID) *executable {
	for _, nodeAccountID := range nodeAccountIDs {
		e.nodeAccountIDs._Push(nodeAccountID)
	}
	e.nodeAccountIDs._SetLocked(true)
	return e
}

func (e *executable) GetLogLevel() *LogLevel {
	return e.logLevel
}

func (e *executable) SetLogLevel(level LogLevel) *executable {
	e.logLevel = &level
	return e
}

func (e *executable) getLogger(clientLogger Logger) Logger {
	if e.logLevel != nil {
		return clientLogger.SubLoggerWithLevel(*e.logLevel)
	}
	return clientLogger
}

func (e *executable) getNodeAccountID() AccountID {
	return e.nodeAccountIDs._GetCurrent().(AccountID)
}

// nolint
func _Execute(client *Client, e Executable) (interface{}, error) {
	var maxAttempts int

	if client.maxAttempts != nil {
		maxAttempts = *client.maxAttempts
	} else {
		maxAttempts = e.GetMaxRetry()
	}

	currentBackoff := e.GetMinBackoff()

	var attempt int64
	var errPersistent error
	var marshaledRequest []byte

	txLogger := e.getLogger(client.logger)
	txID, msg := e.getTransactionIDAndMessage()

	for attempt = int64(0); attempt < int64(maxAttempts); attempt++ {
		var protoRequest interface{}
		var node *_Node
		var ok bool

		// If this is not the first attempt, double the backoff time up to the max backoff time
		if attempt > 0 && currentBackoff <= e.GetMaxBackoff() {
			currentBackoff *= 2
		}

		if e.isTransaction() {
			if attempt > 0 && len(e.GetNodeAccountIDs()) > 1 {
				e.advanceRequest()
			}
		}

		protoRequest = e.makeRequest()
		if len(e.GetNodeAccountIDs()) == 0 {
			node = client.network._GetNode()
		} else {
			nodeAccountID := e.getNodeAccountID()
			if node, ok = client.network._GetNodeForAccountID(nodeAccountID); !ok {
				return TransactionResponse{}, ErrInvalidNodeAccountIDSet{nodeAccountID}
			}
		}

		if e.isTransaction() {
			marshaledRequest, _ = protobuf.Marshal(protoRequest.(*services.Transaction))
		} else {
			marshaledRequest, _ = protobuf.Marshal(protoRequest.(*services.Query))
		}

		node._InUse()

		txLogger.Trace("executing", "requestId", e.getLogID(e), "nodeAccountID", node.accountID.String(), "nodeIPAddress", node.address._String(), "Request Proto", hex.EncodeToString(marshaledRequest))

		if !node._IsHealthy() {
			txLogger.Trace("node is unhealthy, waiting before continuing", "requestId", e.getLogID(e), "delay", node._Wait().String())
			_DelayForAttempt(e.getLogID(e), currentBackoff, attempt, txLogger, errNodeIsUnhealthy)
			continue
		}

		txLogger.Trace("updating node account ID index", "requestId", e.getLogID(e))
		channel, err := node._GetChannel(txLogger)
		if err != nil {
			client.network._IncreaseBackoff(node)
			errPersistent = err
			continue
		}

		e.advanceRequest()

		method := e.getMethod(channel)

		var resp interface{}

		ctx := context.TODO()
		var cancel context.CancelFunc

		if e.GetGrpcDeadline() != nil {
			grpcDeadline := time.Now().Add(*e.GetGrpcDeadline())
			ctx, cancel = context.WithDeadline(ctx, grpcDeadline)
		}

		txLogger.Trace("executing gRPC call", "requestId", e.getLogID(e))

		var marshaledResponse []byte
		if method.query != nil {
			resp, err = method.query(ctx, protoRequest.(*services.Query))
			if err == nil {
				marshaledResponse, _ = protobuf.Marshal(resp.(*services.Response))
			}
		} else {
			resp, err = method.transaction(ctx, protoRequest.(*services.Transaction))
			if err == nil {
				marshaledResponse, _ = protobuf.Marshal(resp.(*services.TransactionResponse))
			}
		}

		if cancel != nil {
			cancel()
		}
		if err != nil {
			errPersistent = err
			if _ExecutableDefaultRetryHandler(e.getLogID(e), err, txLogger) {
				client.network._IncreaseBackoff(node)
				continue
			}
			if errPersistent == nil {
				errPersistent = errors.New("error")
			}

			if e.isTransaction() {
				return TransactionResponse{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
			}

			return &services.Response{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
		}

		node._DecreaseBackoff()

		statusError := e.mapStatusError(e, resp)

		txLogger.Trace(
			msg,
			"requestID", e.getLogID(e),
			"nodeID", node.accountID.String(),
			"nodeAddress", node.address._String(),
			"nodeIsHealthy", strconv.FormatBool(node._IsHealthy()),
			"network", client.GetLedgerID().String(),
			"status", statusError.Error(),
			"txID", txID,
		)

		switch e.shouldRetry(e, resp) {
		case executionStateRetry:
			errPersistent = statusError
			_DelayForAttempt(e.getLogID(e), currentBackoff, attempt, txLogger, errPersistent)
			continue
		case executionStateExpired:
			if e.isTransaction() {
				transaction := e.(TransactionInterface)
				if transaction.regenerateID(client) {
					txLogger.Trace("received `TRANSACTION_EXPIRED` with transaction ID regeneration enabled; regenerating", "requestId", e.getLogID(e))
					continue
				} else {
					return TransactionResponse{}, statusError
				}
			} else {
				return &services.Response{}, statusError
			}
		case executionStateError:
			if e.isTransaction() {
				return TransactionResponse{}, statusError
			}

			return &services.Response{}, statusError
		case executionStateFinished:
			txLogger.Trace("finished", "Response Proto", hex.EncodeToString(marshaledResponse))
			return e.mapResponse(resp, node.accountID, protoRequest)
		}
	}

	if errPersistent == nil {
		errPersistent = errors.New("unknown error occurred after max attempts")
	}

	if e.isTransaction() {
		return TransactionResponse{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
	}

	txLogger.Error("exceeded maximum attempts for request", "last exception being", errPersistent)

	return &services.Response{}, errPersistent
}

func _DelayForAttempt(logID string, backoff time.Duration, attempt int64, logger Logger, err error) {
	logger.Trace("retrying request attempt", "requestId", logID, "delay", backoff, "attempt", attempt+1, "error", err)

	time.Sleep(backoff)
}

func _ExecutableDefaultRetryHandler(logID string, err error, logger Logger) bool {
	code := status.Code(err)
	logger.Trace("received gRPC error with status code", "requestId", logID, "status", code.String())
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
