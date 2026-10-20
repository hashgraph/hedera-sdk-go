//go:build all || unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"bytes"
	"net"
	"runtime/pprof"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/mirror"
	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func TestUnitTopicMessageQueryValidate(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	topicID, err := TopicIDFromString("0.0.123-esxsf")
	require.NoError(t, err)

	topicInfo := NewTopicMessageQuery().
		SetTopicID(topicID)

	err = topicInfo.validateNetworkOnIDs(client)
	require.NoError(t, err)
}

func TestUnitTopicMessageQueryValidateWrong(t *testing.T) {
	t.Parallel()

	client, err := _NewMockClient()
	client.SetLedgerID(*NewLedgerIDTestnet())
	require.NoError(t, err)
	client.SetAutoValidateChecksums(true)
	topicID, err := TopicIDFromString("0.0.123-rmkykd")
	require.NoError(t, err)

	topicInfo := NewTopicMessageQuery().
		SetTopicID(topicID)

	err = topicInfo.validateNetworkOnIDs(client)
	assert.Error(t, err)
	if err != nil {
		assert.Equal(t, "network mismatch or wrong checksum given, given checksum: rmkykd, correct checksum esxsf, network: testnet", err.Error())
	}
}

func TestUnitTopicMessageQueryGet(t *testing.T) {
	t.Parallel()

	topicID := TopicID{Topic: 7}

	balance := NewTopicMessageQuery().
		SetTopicID(topicID).
		SetStartTime(time.Now()).
		SetMaxAttempts(23).
		SetEndTime(time.Now()).
		SetLimit(32).
		SetCompletionHandler(func() {}).
		SetErrorHandler(func(stat status.Status) {}).
		SetRetryHandler(func(err error) bool { return false })

	balance.GetTopicID()
	balance.GetStartTime()
	balance.GetMaxAttempts()
	balance.GetEndTime()
	balance.GetLimit()
}

func TestUnitTopicMessageQueryNothingSet(t *testing.T) {
	t.Parallel()

	balance := NewTopicMessageQuery()

	balance.GetTopicID()
	balance.GetStartTime()
	balance.GetMaxAttempts()
	balance.GetEndTime()
	balance.GetLimit()
}

type topicMessageMirror struct {
	mirror.UnimplementedConsensusServiceServer
	mu        sync.Mutex
	queries   []*mirror.ConsensusTopicQuery
	subscribe func(call int, stream grpc.ServerStreamingServer[mirror.ConsensusTopicResponse]) error
}

func (m *topicMessageMirror) SubscribeTopic(query *mirror.ConsensusTopicQuery, stream grpc.ServerStreamingServer[mirror.ConsensusTopicResponse]) error {
	m.mu.Lock()
	m.queries = append(m.queries, query)
	call := len(m.queries)
	m.mu.Unlock()

	return m.subscribe(call, stream)
}

func (m *topicMessageMirror) calls() int {
	m.mu.Lock()
	defer m.mu.Unlock()

	return len(m.queries)
}

func (m *topicMessageMirror) query(call int) *mirror.ConsensusTopicQuery {
	m.mu.Lock()
	defer m.mu.Unlock()

	return m.queries[call-1]
}

func setupTopicMessageMirror(t *testing.T, subscribe func(call int, stream grpc.ServerStreamingServer[mirror.ConsensusTopicResponse]) error) (*Client, *topicMessageMirror) {
	t.Helper()

	mirrorServer := &topicMessageMirror{subscribe: subscribe}
	client := setupMirrorServer(t, func(server *grpc.Server) {
		mirror.RegisterConsensusServiceServer(server, mirrorServer)
	})

	return client, mirrorServer
}

func topicMessageResponse(sequence uint64) *mirror.ConsensusTopicResponse {
	return &mirror.ConsensusTopicResponse{
		ConsensusTimestamp: &services.Timestamp{Seconds: int64(1000 + sequence)},
		Message:            []byte{byte(sequence)},
		SequenceNumber:     sequence,
	}
}

type topicMessageRecorder struct {
	messages  chan TopicMessage
	completed chan struct{}
	errors    chan status.Status
}

func newTopicMessageRecorder() *topicMessageRecorder {
	return &topicMessageRecorder{
		messages:  make(chan TopicMessage, 10),
		completed: make(chan struct{}),
		errors:    make(chan status.Status, 1),
	}
}

func (r *topicMessageRecorder) query() *TopicMessageQuery {
	return NewTopicMessageQuery().
		SetTopicID(TopicID{Topic: 5}).
		SetCompletionHandler(func() { close(r.completed) }).
		SetErrorHandler(func(stat status.Status) { r.errors <- stat })
}

func (r *topicMessageRecorder) onNext(message TopicMessage) {
	r.messages <- message
}

func (r *topicMessageRecorder) waitMessages(t *testing.T, count int) []TopicMessage {
	t.Helper()

	messages := make([]TopicMessage, 0, count)
	for len(messages) < count {
		select {
		case message := <-r.messages:
			messages = append(messages, message)
		case <-time.After(5 * time.Second):
			require.FailNow(t, "timed out waiting for topic messages", "received %d of %d", len(messages), count)
		}
	}

	return messages
}

func (r *topicMessageRecorder) waitCompleted(t *testing.T) {
	t.Helper()

	select {
	case <-r.completed:
	case <-time.After(5 * time.Second):
		require.FailNow(t, "timed out waiting for the completion handler")
	}
}

func (r *topicMessageRecorder) waitError(t *testing.T) status.Status {
	t.Helper()

	select {
	case stat := <-r.errors:
		return stat
	case <-time.After(5 * time.Second):
		require.FailNow(t, "timed out waiting for the error handler")
		return status.Status{}
	}
}

func sequenceNumbers(messages []TopicMessage) []uint64 {
	sequences := make([]uint64, 0, len(messages))
	for _, message := range messages {
		sequences = append(sequences, message.SequenceNumber)
	}

	return sequences
}

func TestUnitTopicMessageQueryReconnectsFromLastMessage(t *testing.T) {
	t.Parallel()

	client, mirrorServer := setupTopicMessageMirror(t, func(call int, stream grpc.ServerStreamingServer[mirror.ConsensusTopicResponse]) error {
		if call == 1 {
			_ = stream.Send(topicMessageResponse(1))
			return status.Error(codes.Unavailable, "mirror node restarting")
		}
		_ = stream.Send(topicMessageResponse(2))
		_ = stream.Send(topicMessageResponse(3))
		return nil
	})
	recorder := newTopicMessageRecorder()

	handle, err := recorder.query().SetLimit(3).Subscribe(client, recorder.onNext)
	require.NoError(t, err)
	defer handle.Unsubscribe()

	require.Equal(t, []uint64{1, 2, 3}, sequenceNumbers(recorder.waitMessages(t, 3)))
	recorder.waitCompleted(t)
	require.Equal(t, 2, mirrorServer.calls())

	resumed := mirrorServer.query(2)
	require.Equal(t, int64(1001), resumed.ConsensusStartTime.Seconds)
	require.Equal(t, int32(1), resumed.ConsensusStartTime.Nanos)
	require.Equal(t, uint64(2), resumed.Limit)
}

func TestUnitTopicMessageQueryCompletesWhenLimitReachedBeforeError(t *testing.T) {
	t.Parallel()

	client, mirrorServer := setupTopicMessageMirror(t, func(_ int, stream grpc.ServerStreamingServer[mirror.ConsensusTopicResponse]) error {
		_ = stream.Send(topicMessageResponse(1))
		return status.Error(codes.Unavailable, "mirror node restarting")
	})
	recorder := newTopicMessageRecorder()

	handle, err := recorder.query().SetLimit(1).Subscribe(client, recorder.onNext)
	require.NoError(t, err)
	defer handle.Unsubscribe()

	recorder.waitMessages(t, 1)
	recorder.waitCompleted(t)
	require.Equal(t, 1, mirrorServer.calls())
}

func TestUnitTopicMessageQuerySubscriptionIgnoresLaterQueryChanges(t *testing.T) {
	t.Parallel()

	release := make(chan struct{})
	client, mirrorServer := setupTopicMessageMirror(t, func(call int, stream grpc.ServerStreamingServer[mirror.ConsensusTopicResponse]) error {
		if call == 1 {
			_ = stream.Send(topicMessageResponse(1))
			<-release
			return status.Error(codes.Unavailable, "mirror node restarting")
		}
		_ = stream.Send(topicMessageResponse(2))
		return nil
	})
	recorder := newTopicMessageRecorder()
	query := recorder.query()

	handle, err := query.Subscribe(client, recorder.onNext)
	require.NoError(t, err)
	defer handle.Unsubscribe()

	recorder.waitMessages(t, 1)
	query.SetLimit(1).SetMaxAttempts(0)
	close(release)

	require.Equal(t, []uint64{2}, sequenceNumbers(recorder.waitMessages(t, 1)))
	recorder.waitCompleted(t)
	require.Equal(t, uint64(0), mirrorServer.query(2).Limit)
}

func TestUnitTopicMessageQueryReassemblesChunksAcrossReconnect(t *testing.T) {
	t.Parallel()

	transactionID := TransactionIDGenerate(AccountID{Account: 3})._ToProtobuf()
	chunk := func(number int32, contents string) *mirror.ConsensusTopicResponse {
		resp := topicMessageResponse(uint64(number))
		resp.Message = []byte(contents)
		resp.ChunkInfo = &services.ConsensusMessageChunkInfo{InitialTransactionID: transactionID, Total: 2, Number: number}
		return resp
	}

	client, _ := setupTopicMessageMirror(t, func(call int, stream grpc.ServerStreamingServer[mirror.ConsensusTopicResponse]) error {
		if call == 1 {
			_ = stream.Send(chunk(1, "hello "))
			return status.Error(codes.Unavailable, "mirror node restarting")
		}
		_ = stream.Send(chunk(2, "world"))
		return nil
	})
	recorder := newTopicMessageRecorder()

	handle, err := recorder.query().Subscribe(client, recorder.onNext)
	require.NoError(t, err)
	defer handle.Unsubscribe()

	message := recorder.waitMessages(t, 1)[0]
	recorder.waitCompleted(t)
	require.Equal(t, "hello world", string(message.Contents))
	require.Len(t, message.Chunks, 2)
	require.Empty(t, recorder.messages)
}

func TestUnitTopicMessageQueryResetsAttemptsAfterReceivingMessage(t *testing.T) {
	t.Parallel()

	client, mirrorServer := setupTopicMessageMirror(t, func(call int, stream grpc.ServerStreamingServer[mirror.ConsensusTopicResponse]) error {
		_ = stream.Send(topicMessageResponse(uint64(call)))
		if call < 3 {
			return status.Error(codes.Unavailable, "mirror node restarting")
		}
		return nil
	})
	recorder := newTopicMessageRecorder()

	handle, err := recorder.query().SetMaxAttempts(1).Subscribe(client, recorder.onNext)
	require.NoError(t, err)
	defer handle.Unsubscribe()

	require.Equal(t, []uint64{1, 2, 3}, sequenceNumbers(recorder.waitMessages(t, 3)))
	recorder.waitCompleted(t)
	require.Equal(t, 3, mirrorServer.calls())
}

func TestUnitTopicMessageQueryMaxAttemptsReached(t *testing.T) {
	t.Parallel()

	client, mirrorServer := setupTopicMessageMirror(t, func(int, grpc.ServerStreamingServer[mirror.ConsensusTopicResponse]) error {
		return status.Error(codes.Unavailable, "mirror node restarting")
	})
	recorder := newTopicMessageRecorder()

	handle, err := recorder.query().SetMaxAttempts(1).Subscribe(client, recorder.onNext)
	require.NoError(t, err)
	defer handle.Unsubscribe()

	stat := recorder.waitError(t)
	require.Equal(t, codes.Unavailable, stat.Code())
	require.Contains(t, stat.Message(), "max attempts reached")
	require.Equal(t, 2, mirrorServer.calls())
}

func TestUnitTopicMessageQueryNonRetryableError(t *testing.T) {
	t.Parallel()

	client, mirrorServer := setupTopicMessageMirror(t, func(int, grpc.ServerStreamingServer[mirror.ConsensusTopicResponse]) error {
		return status.Error(codes.InvalidArgument, "invalid topic")
	})
	recorder := newTopicMessageRecorder()

	handle, err := recorder.query().Subscribe(client, recorder.onNext)
	require.NoError(t, err)
	defer handle.Unsubscribe()

	stat := recorder.waitError(t)
	require.Equal(t, codes.InvalidArgument, stat.Code())
	require.Equal(t, "invalid topic", stat.Message())
	require.Equal(t, 1, mirrorServer.calls())
}

func TestUnitTopicMessageQueryUnsubscribeDuringBackoff(t *testing.T) {
	t.Parallel()

	client, mirrorServer := setupTopicMessageMirror(t, func(int, grpc.ServerStreamingServer[mirror.ConsensusTopicResponse]) error {
		return status.Error(codes.Unavailable, "mirror node restarting")
	})
	recorder := newTopicMessageRecorder()
	var retries atomic.Int32
	backingOff := make(chan struct{})

	handle, err := recorder.query().
		SetMaxAttempts(5).
		SetRetryHandler(func(error) bool {
			if retries.Add(1) == 3 {
				close(backingOff)
			}
			return true
		}).
		Subscribe(client, recorder.onNext)
	require.NoError(t, err)

	select {
	case <-backingOff:
	case <-time.After(5 * time.Second):
		require.FailNow(t, "timed out waiting for the third retry")
	}

	start := time.Now()
	handle.Unsubscribe()
	recorder.waitCompleted(t)
	require.Less(t, time.Since(start), 500*time.Millisecond)

	time.Sleep(1500 * time.Millisecond)
	require.Equal(t, 3, mirrorServer.calls())
}

func TestUnitTopicMessageQueryUnsubscribeStopsSubscription(t *testing.T) {
	client, _ := setupTopicMessageMirror(t, func(_ int, stream grpc.ServerStreamingServer[mirror.ConsensusTopicResponse]) error {
		_ = stream.Send(topicMessageResponse(1))
		<-stream.Context().Done()
		return stream.Context().Err()
	})
	recorder := newTopicMessageRecorder()
	before := countTopicSubscriptionGoroutines(t)

	handle, err := recorder.query().Subscribe(client, recorder.onNext)
	require.NoError(t, err)

	recorder.waitMessages(t, 1)
	require.Greater(t, countTopicSubscriptionGoroutines(t), before)

	handle.Unsubscribe()
	recorder.waitCompleted(t)
	require.Eventually(t, func() bool {
		return countTopicSubscriptionGoroutines(t) == before
	}, 2*time.Second, 10*time.Millisecond)
}

func countTopicSubscriptionGoroutines(t *testing.T) int {
	var stacks bytes.Buffer
	require.NoError(t, pprof.Lookup("goroutine").WriteTo(&stacks, 2))

	count := 0
	for _, goroutine := range strings.Split(stacks.String(), "\n\n") {
		if strings.Contains(goroutine, "(*TopicMessageQuery).Subscribe.func1(") || strings.Contains(goroutine, "processProtoMessageStream") {
			count++
		}
	}

	return count
}

func TestUnitTopicMessageQuerySubscribeReturnsOpenError(t *testing.T) {
	t.Parallel()

	listener, err := net.Listen("tcp", "localhost:0")
	require.NoError(t, err)
	address := listener.Addr().String()
	require.NoError(t, listener.Close())

	client, err := _NewMockClient()
	require.NoError(t, err)
	defer client.Close()
	client.SetMirrorNetwork([]string{address})
	recorder := newTopicMessageRecorder()

	_, err = recorder.query().Subscribe(client, recorder.onNext)
	require.Equal(t, codes.Unavailable, status.Code(err))
	require.Empty(t, recorder.errors)
}
