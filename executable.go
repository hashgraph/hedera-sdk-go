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
	GetGrpcDeadline() time.Duration
	GetMaxRetry() int
	GetNodeAccountIDs() []AccountID
	GetLogLevel() *LogLevel


	shouldRetry(interface{}, interface{}) _ExecutionState
	makeRequest(interface{}) interface{}
	advanceRequest(interface{})
	getNodeAccountID(interface{}) AccountID
	getMethod(*_Channel) _Method
	mapStatusError(interface{}, interface{}) error
	mapResponse(interface{}, interface{}, AccountID, interface{}) (interface{}, error)
	getName() string
	build() *services.TransactionBody
	buildScheduled() (*services.SchedulableTransactionBody, error)
	validateNetworkOnIDs(client *Client) error
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


func (this *executable) GetMaxBackoff() time.Duration {
	if this.maxBackoff != nil {
		return *this.maxBackoff
	}

	return 8 * time.Second
}

func (this *executable) GetMinBackoff() time.Duration {
	if this.minBackoff != nil {
		return *this.minBackoff
	}

	return 250 * time.Millisecond
}

func (this *executable) SetMaxBackoff(max time.Duration) *executable {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < this.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	this.maxBackoff = &max
	return this
}

func (this *executable) SetMinBackoff(max time.Duration) *executable{
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < this.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	this.maxBackoff = &max
	return this
}

// GetGrpcDeadline returns the grpc deadline
func (this *executable) GetGrpcDeadline() time.Duration {
	return *this.grpcDeadline
}

// When execution is attempted, a single attempt will timeout when this deadline is reached. (The SDK may subsequently retry the execution.)
func (this *executable) SetGrpcDeadline(deadline *time.Duration) *executable {
	this.grpcDeadline = deadline
	return this
}

// GetMaxRetry returns the max number of errors before execution will fail.
func (this *executable) GetMaxRetry() int{
	return this.maxRetry
}

func (this *executable) SetMaxRetry(max int) *executable {
	this.maxRetry = max
	return this
}

// GetNodeAccountID returns the node AccountID for this transaction.
func (this *executable) GetNodeAccountIDs() []AccountID{
	nodeAccountIDs := []AccountID{}

	for _, value := range this.nodeAccountIDs.slice {
		nodeAccountIDs = append(nodeAccountIDs, value.(AccountID))
	}

	return nodeAccountIDs
}

func (this *executable) SetNodeAccountIDs(nodeAccountIDs []AccountID) *executable{
	for _, nodeAccountID := range nodeAccountIDs {
		this.nodeAccountIDs._Push(nodeAccountID)
	}
	this.nodeAccountIDs._SetLocked(true)
	return this
}

func (this *transaction) GetLogLevel() *LogLevel {
	return this.logLevel
}

func (this *executable) SetLogLevel(level LogLevel) *executable{
	this.logLevel = &level
	return this
}


func getLogger(request interface{}, clientLogger Logger) Logger {
	switch req := request.(type) {
	case *transaction:
		if req.logLevel != nil {
			return clientLogger.SubLoggerWithLevel(*req.logLevel)
		}
	case *Query:
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
	case *Query:
		txID := req.GetPaymentTransactionID().String()
		if txID == "" {
			txID = "None"
		}
		return txID, "Query status received"
	default:
		return "", ""
	}
}

func _Execute(client *Client, e Executable) (interface{}, error){
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
	txID, msg := getTransactionIDAndMessage(e)

	for attempt = int64(0); attempt < int64(maxAttempts); attempt, currentBackoff = attempt+1, currentBackoff*2 {
		var protoRequest interface{}
		var node *_Node

		if transaction, ok := e.(*transaction); ok {
			if attempt > 0 && transaction.nodeAccountIDs._Length() > 1 {
				e.advanceRequest(e);
			}
            protoRequest = e.makeRequest(e)
			nodeAccountID := e.getNodeAccountID(e)
			if node, ok = client.network._GetNodeForAccountID(nodeAccountID); !ok {
				return TransactionResponse{}, ErrInvalidNodeAccountIDSet{nodeAccountID}
			}

			marshaledRequest, _ = protobuf.Marshal(protoRequest.(*services.Transaction))
		} else if query, ok := e.(*Query); ok {
			if query.nodeAccountIDs.locked && query.nodeAccountIDs._Length() > 0 {
				protoRequest = e.makeRequest(e)
				nodeAccountID := e.getNodeAccountID(e)
				if node, ok = client.network._GetNodeForAccountID(nodeAccountID); !ok {
					return &services.Response{}, ErrInvalidNodeAccountIDSet{nodeAccountID}
				}
			} else {
				node = client.network._GetNode()
				if len(query.paymentTransactions) > 0 {
					var paymentTransaction services.TransactionBody
					_ = protobuf.Unmarshal(query.paymentTransactions[0].BodyBytes, &paymentTransaction) // nolint
					paymentTransaction.NodeAccountID = node.accountID._ToProtobuf()

					transferTx := paymentTransaction.Data.(*services.TransactionBody_CryptoTransfer)
					transferTx.CryptoTransfer.Transfers.AccountAmounts[0].AccountID = node.accountID._ToProtobuf()
					query.paymentTransactions[0].BodyBytes, _ = protobuf.Marshal(&paymentTransaction) // nolint

					signature := client.operator.signer(query.paymentTransactions[0].BodyBytes) // nolint
					sigPairs := make([]*services.SignaturePair, 0)
					sigPairs = append(sigPairs, client.operator.publicKey._ToSignaturePairProtobuf(signature))

					query.paymentTransactions[0].SigMap = &services.SignatureMap{ // nolint
						SigPair: sigPairs,
					}
				}
				query.nodeAccountIDs._Set(0, node.accountID)
				protoRequest = e.makeRequest(e)
			}
			marshaledRequest, _ = protobuf.Marshal(protoRequest.(*services.Query))
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

		e.advanceRequest(e)

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

			if _, ok := e.(*transaction); ok {
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
			"status",e.mapStatusError(e, resp).Error(),
			"txID", txID,
		)

		switch e.shouldRetry(e, resp) {
		case executionStateRetry:
			errPersistent = e.mapStatusError(e, resp)
			_DelayForAttempt(e.getName(), backOff.NextBackOff(), attempt, txLogger)
			continue
		case executionStateExpired:
			if transaction, ok := e.(*transaction); ok {
				if !client.GetOperatorAccountID()._IsZero() && transaction.regenerateTransactionID && !transaction.transactionIDs.locked {
					txLogger.Trace("received `TRANSACTION_EXPIRED` with transaction ID regeneration enabled; regenerating", "requestId", e.getName())
					transaction.transactionIDs._Set(transaction.transactionIDs.index, TransactionIDGenerate(client.GetOperatorAccountID()))
					if err != nil {
						panic(err)
					}
					continue
				} else {
					return TransactionResponse{}, e.mapStatusError(e, resp)
				}
			} else {
				return &services.Response{}, e.mapStatusError(e, resp)
			}
		case executionStateError:
			if _, ok := e.(*transaction); ok {
				return TransactionResponse{}, e.mapStatusError(e, resp)
			}

			return &services.Response{}, e.mapStatusError(e, resp)
		case executionStateFinished:
			txLogger.Trace("finished", "Response Proto", hex.EncodeToString(marshaledResponse))
			return e.mapResponse(e, resp, node.accountID, protoRequest)
		}
	}

	if errPersistent == nil {
		errPersistent = errors.New("error")
	}

	if _, ok := e.(*transaction); ok {
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
