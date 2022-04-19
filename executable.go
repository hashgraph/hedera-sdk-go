package hedera

/*-
 *
 * Hedera Go SDK
 *
 * Copyright (C) 2020 - 2022 Hedera Hashgraph, LLC
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
	"math"
	"os"
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"

	"github.com/pkg/errors"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var logCtx zerolog.Logger

// A required init function to setup logging at the correct level
func init() { // nolint
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix

	if os.Getenv("HEDERA_SDK_GO_LOG_PRETTY") != "" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	switch os.Getenv("HEDERA_SDK_GO_LOG_LEVEL") {
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "INFO":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.Disabled)
	}

	logCtx = log.With().Str("module", "hedera-sdk-go").Logger()
}

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

func _Execute( // nolint
	client *Client,
	request interface{},
	shouldRetry func(string, interface{}, interface{}) _ExecutionState,
	makeRequest func(interface{}) interface{},
	advanceRequest func(interface{}),
	getNodeAccountID func(interface{}) AccountID,
	getMethod func(interface{}, *_Channel) _Method,
	mapStatusError func(interface{}, interface{}) error,
	mapResponse func(interface{}, interface{}, AccountID, interface{}) (interface{}, error),
	logID string,
	deadline *time.Duration,
	maxBackoff *time.Duration,
	minBackoff *time.Duration,
	maxRetry int,
) (interface{}, error) {
	var maxAttempts int

	if client.maxAttempts != nil {
		maxAttempts = *client.maxAttempts
	} else {
		maxAttempts = maxRetry
	}

	currentBackoff := minBackoff

	var attempt int64
	var errPersistent error

	for attempt = int64(0); attempt < int64(maxAttempts); attempt, *currentBackoff = attempt+1, *currentBackoff*2 {
		var protoRequest interface{}
		var node *_Node

		if transaction, ok := request.(*Transaction); ok {
			if transaction.nodeAccountIDs.locked && transaction.nodeAccountIDs._Length() > 0 {
				protoRequest = makeRequest(request)
				nodeAccountID := getNodeAccountID(request)
				if node, ok = client.network._GetNodeForAccountID(nodeAccountID); !ok {
					return TransactionResponse{}, ErrInvalidNodeAccountIDSet{nodeAccountID}
				}
			} else {
				node = client.network._GetNode()
				transaction.nodeAccountIDs._Set(0, node.accountID)
				protoTransaction, _ := transaction._BuildTransaction(0)
				protoRequest = protoTransaction
			}
		} else if query, ok := request.(*Query); ok {
			if query.nodeAccountIDs.locked && query.nodeAccountIDs._Length() > 0 {
				protoRequest = makeRequest(request)
				nodeAccountID := getNodeAccountID(request)
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
					request.query.paymentTransactions[0].BodyBytes, _ = protobuf.Marshal(&paymentTransaction) // nolint
				}
				query.nodeAccountIDs._Set(0, node.accountID)
				protoRequest = makeRequest(request)
			}
		}

		node._InUse()

		logCtx.Trace().Str("requestId", logID).Str("nodeAccountID", node.accountID.String()).Str("nodeIPAddress", node.address._String())

		if !node._IsHealthy() {
			logCtx.Trace().Str("requestId", logID).Str("delay", node._Wait().String()).Msg("node is unhealthy, waiting before continuing")
			_DelayForAttempt(logID, minBackoff, maxBackoff, attempt)
			continue
		}

		logCtx.Trace().Str("requestId", logID).Msg("updating node account ID index")
		advanceRequest(request)

		channel, err := node._GetChannel()
		if err != nil {
			client.network._IncreaseBackoff(node)
			continue
		}

		method := getMethod(request, channel)

		var resp interface{}

		ctx := context.TODO()
		var cancel context.CancelFunc
		if deadline != nil {
			grpcDeadline := time.Now().Add(*deadline)
			ctx, cancel = context.WithDeadline(ctx, grpcDeadline)
		}

		logCtx.Trace().Str("requestId", logID).Msg("executing gRPC call")
		if method.query != nil {
			resp, err = method.query(ctx, protoRequest.(*services.Query))
		} else {
			resp, err = method.transaction(ctx, protoRequest.(*services.Transaction))
		}

		if cancel != nil {
			cancel()
		}

		if err != nil {
			errPersistent = err
			if _ExecutableDefaultRetryHandler(logID, err) {
				client.network._IncreaseBackoff(node)
				continue
			}
			if errPersistent == nil {
				errPersistent = errors.New("error")
			}

			if _, ok := request.(*Transaction); ok {
				return TransactionResponse{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
			}

			return &services.Response{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
		}

		node._DecreaseBackoff()

		switch shouldRetry(logID, request, resp) {
		case executionStateRetry:
			errPersistent = mapStatusError(request, resp)
			_DelayForAttempt(logID, minBackoff, maxBackoff, attempt)
			continue
		case executionStateExpired:
			if transaction, ok := request.(*Transaction); ok {
				if !client.GetOperatorAccountID()._IsZero() && transaction.regenerateTransactionID && !transaction.transactionIDs.locked {
					logCtx.Trace().Str("requestId", logID).Msg("received `TRANSACTION_EXPIRED` with transaction ID regeneration enabled; regenerating")
					transaction.transactionIDs._Set(transaction.transactionIDs.index, TransactionIDGenerate(client.GetOperatorAccountID()))
					if err != nil {
						panic(err)
					}
					continue
				} else {
					return TransactionResponse{}, mapStatusError(request, resp)
				}
			} else {
				return &services.Response{}, mapStatusError(request, resp)
			}
		case executionStateError:
			if _, ok := request.(*Transaction); ok {
				return TransactionResponse{}, mapStatusError(request, resp)
			}

			return &services.Response{}, mapStatusError(request, resp)
		case executionStateFinished:
			return mapResponse(request, resp, node.accountID, protoRequest)
		}
	}

	if errPersistent == nil {
		errPersistent = errors.New("error")
	}

	if _, ok := request.(*Transaction); ok {
		return TransactionResponse{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
	}

	return &services.Response{}, errors.Wrapf(errPersistent, "retry %d/%d", attempt, maxAttempts)
}

func _DelayForAttempt(logID string, minBackoff *time.Duration, maxBackoff *time.Duration, attempt int64) {
	// 0.1s, 0.2s, 0.4s, 0.8s, ...
	ms := int64(math.Min(float64(minBackoff.Milliseconds())*math.Pow(2, float64(attempt)), float64(maxBackoff.Milliseconds())))
	logCtx.Trace().Str("requestId", logID).Dur("delay", time.Duration(ms)).Int64("attempt", attempt+1).Msg("retrying  request attempt")
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func _ExecutableDefaultRetryHandler(logID string, err error) bool {
	code := status.Code(err)
	logCtx.Trace().Str("requestId", logID).Str("status", code.String()).Msg("received gRPC error with status code")

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
