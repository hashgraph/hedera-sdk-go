package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2023 Hedera Hashgraph, LLC
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

import (
	"context"
	"encoding/hex"
	"strconv"
	"time"

	"github.com/cenkalti/backoff/v4"

	protobuf "google.golang.org/protobuf/proto"

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

type Executable interface {
	GetMaxBackoff() time.Duration
	GetMinBackoff() time.Duration
	GetGrpcDeadline() *time.Duration
	GetMaxRetry() int
	GetNodeAccountIDs() []AccountID
	GetLogLevel() *LogLevel

	shouldRetry(interface{}) _ExecutionState
	makeRequest() interface{}
	advanceRequest()
	getNodeAccountID() AccountID
	getMethod(*_Channel) _Method
	mapStatusError(interface{}) error
	mapResponse(interface{}, AccountID, interface{}) (interface{}, error)
	getName() string
	validateNetworkOnIDs(client *Client) error
	build() *services.TransactionBody
	buildScheduled() (*services.SchedulableTransactionBody, error)
	isTransaction() bool
}

type executable struct {
	e              Executable
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

func getLogger(request interface{}, clientLogger Logger) Logger {
	switch req := request.(type) {
	case *transaction:
		if req.logLevel != nil {
			return clientLogger.SubLoggerWithLevel(*req.logLevel)
		}
	case *query:
		if req.logLevel != nil {
			return clientLogger.SubLoggerWithLevel(*req.logLevel)
		}
	}
	return clientLogger
}

func getTransactionIDAndMessage(request interface{}) (string, string) {
	switch req := request.(type) {
	case *transaction:
		return req.GetTransactionID().String(), "transaction status received"
	case *query:
		txID := req.GetPaymentTransactionID().String()
		if txID == "" {
			txID = "None"
		}
		return txID, "Query status received"
	default:
		return "", ""
	}
}

func _Execute(client *Client, e Executable) (interface{}, error) {
	var maxAttempts int
	backOff := backoff.NewExponentialBackOff()
	backOff.InitialInterval = e.GetMinBackoff()
	backOff.MaxInterval = e.GetMaxBackoff()
	backOff.Multiplier = 2

	if client.maxAttempts != nil {
		maxAttempts = *client.maxAttempts
	} else {
		maxAttempts = e.GetMaxRetry()
	}

	currentBackoff := e.GetMinBackoff()

	var attempt int64
	var errPersistent error
	var marshaledRequest []byte

	txLogger := getLogger(e, client.logger)
	//txID, msg := getTransactionIDAndMessage()
	txID, msg := "TODO", "TODO"

	for attempt = int64(0); attempt < int64(maxAttempts); attempt, currentBackoff = attempt+1, currentBackoff*2 {
		var protoRequest interface{}
		var node *_Node
		var ok bool

		if e.isTransaction() {
			if attempt > 0 && len(e.GetNodeAccountIDs()) > 1 {
				e.advanceRequest()
			}
			protoRequest = e.makeRequest()
			nodeAccountID := e.getNodeAccountID()
			if node, ok = client.network._GetNodeForAccountID(nodeAccountID); !ok {
				return TransactionResponse{}, ErrInvalidNodeAccountIDSet{nodeAccountID}
			}

			marshaledRequest, _ = protobuf.Marshal(protoRequest.(*services.Transaction))
			//} else {
			//	if query.nodeAccountIDs.locked && query.nodeAccountIDs._Length() > 0 {
			//		protoRequest = e.makeRequest()
			//		nodeAccountID := e.getNodeAccountID()
			//		if node, ok = client.network._GetNodeForAccountID(nodeAccountID); !ok {
			//			return &services.Response{}, ErrInvalidNodeAccountIDSet{nodeAccountID}
			//		}
			//	} else {
			//		node = client.network._GetNode()
			//		if len(query.paymentTransactions) > 0 {
			//			var paymentTransaction services.TransactionBody
			//			_ = protobuf.Unmarshal(query.paymentTransactions[0].BodyBytes, &paymentTransaction) // nolint
			//			paymentTransaction.NodeAccountID = node.accountID._ToProtobuf()
			//
			//			transferTx := paymentTransaction.Data.(*services.TransactionBody_CryptoTransfer)
			//			transferTx.CryptoTransfer.Transfers.AccountAmounts[0].AccountID = node.accountID._ToProtobuf()
			//			query.paymentTransactions[0].BodyBytes, _ = protobuf.Marshal(&paymentTransaction) // nolint
			//
			//			signature := client.operator.signer(query.paymentTransactions[0].BodyBytes) // nolint
			//			sigPairs := make([]*services.SignaturePair, 0)
			//			sigPairs = append(sigPairs, client.operator.publicKey._ToSignaturePairProtobuf(signature))
			//
			//			query.paymentTransactions[0].SigMap = &services.SignatureMap{ // nolint
			//				SigPair: sigPairs,
			//			}
			//		}
			//		query.nodeAccountIDs._Set(0, node.accountID)
			//		protoRequest = e.makeRequest()
			//	}
			//	marshaledRequest, _ = protobuf.Marshal(protoRequest.(*services.Query))
		}

		node._InUse()

		txLogger.Trace("executing", "requestId", e.getName(), "nodeAccountID", node.accountID.String(), "nodeIPAddress", node.address._String(), "Request Proto", hex.EncodeToString(marshaledRequest))

		if !node._IsHealthy() {
			txLogger.Trace("node is unhealthy, waiting before continuing", "requestId", e.getName(), "delay", node._Wait().String())
			_DelayForAttempt(e.getName(), backOff.NextBackOff(), attempt, txLogger)
			continue
		}

		txLogger.Trace("updating node account ID index", "requestId", e.getName())
		channel, err := node._GetChannel(txLogger)
		if err != nil {
			client.network._IncreaseBackoff(node)
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

		txLogger.Trace("executing gRPC call", "requestId", e.getName())

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
			if _ExecutableDefaultRetryHandler(e.getName(), err, txLogger) {
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

		txLogger.Trace(
			msg,
			"requestID", e.getName(),
			"nodeID", node.accountID.String(),
			"nodeAddress", node.address._String(),
			"nodeIsHealthy", strconv.FormatBool(node._IsHealthy()),
			"network", client.GetLedgerID().String(),
			"status", e.mapStatusError(resp).Error(),
			"txID", txID,
		)

		switch e.shouldRetry(resp) {
		case executionStateRetry:
			errPersistent = e.mapStatusError(resp)
			_DelayForAttempt(e.getName(), backOff.NextBackOff(), attempt, txLogger)
			continue
		case executionStateExpired:
			//if e.isTransaction() {
			//	if !client.GetOperatorAccountID()._IsZero() && transaction.regenerateTransactionID && !transaction.transactionIDs.locked {
			//		txLogger.Trace("received `TRANSACTION_EXPIRED` with transaction ID regeneration enabled; regenerating", "requestId", e.getName())
			//		transaction.transactionIDs._Set(transaction.transactionIDs.index, TransactionIDGenerate(client.GetOperatorAccountID()))
			//		if err != nil {
			//			panic(err)
			//		}
			//		continue
			//	} else {
			//		return TransactionResponse{}, e.mapStatusError(resp)
			//	}
			//} else {
			//	return &services.Response{}, e.mapStatusError(resp)
			//}
		case executionStateError:
			if e.isTransaction() {
				return TransactionResponse{}, e.mapStatusError(resp)
			}

			return &services.Response{}, e.mapStatusError(resp)
		case executionStateFinished:
			txLogger.Trace("finished", "Response Proto", hex.EncodeToString(marshaledResponse))
			return e.mapResponse(resp, node.accountID, protoRequest)
		}
	}

	if errPersistent == nil {
		errPersistent = errors.New("error")
	}

	if e.isTransaction() {
		return TransactionResponse{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
	}

	return &services.Response{}, errPersistent
}

func _DelayForAttempt(logID string, backoff time.Duration, attempt int64, logger Logger) {
	logger.Trace("retrying request attempt", "requestId", logID, "delay", backoff, "attempt", attempt+1)

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
