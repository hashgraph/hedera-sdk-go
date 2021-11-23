//+build all unit

package hedera

import (
	"context"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"net"
	"testing"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func TestUnitMock(t *testing.T) {
	var network = map[string]AccountID{
		"localhost:50211": {Account: 3},
	}

	client := ClientForNetwork(network)
	server := grpc.NewServer()

	responses := []interface{}{
		status.Error(codes.Unavailable, ""),
		&proto.Response{
			Response: &proto.Response_CryptogetAccountBalance{
				CryptogetAccountBalance: &proto.CryptoGetAccountBalanceResponse{
					Header: &proto.ResponseHeader{
						NodeTransactionPrecheckCode: 0,
					},
					Balance: 10,
				},
			},
		},
	}

	go func() {
		server.RegisterService(NewServiceDescription(responses), nil)
		lis, err := net.Listen("tcp", "localhost:50211")
		assert.NoError(t, err)

		err = server.Serve(lis)
		assert.NoError(t, err)
	}()

	balance, err := NewAccountBalanceQuery().SetAccountID(AccountID{Account: 3}).Execute(client)
	assert.NoError(t, err)
	assert.Equal(t, balance.Hbars.tinybar, int64(10))
}

func NewMockHandler(responses []interface{}) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	index := 0
	return func(_srv interface{}, _ctx context.Context, _dec func(interface{}) error, _interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		response := responses[index]
		index += 1

		switch response := response.(type) {
		case error:
			return nil, response
		default:
			return response, nil
		}
	}
}

// TODO: Create method that encompasses all services
func NewServiceDescription(responses []interface{}) *grpc.ServiceDesc {
	handler := NewMockHandler(responses)
	return &grpc.ServiceDesc{
		ServiceName: "proto.CryptoService",
		HandlerType: (*proto.CryptoServiceServer)(nil),
		Methods: []grpc.MethodDesc{
			{
				MethodName: "createAccount",
				Handler:    handler,
			},
			{
				MethodName: "updateAccount",
				Handler:    handler,
			},
			{
				MethodName: "cryptoTransfer",
				Handler:    handler,
			},
			{
				MethodName: "cryptoDelete",
				Handler:    handler,
			},
			{
				MethodName: "addLiveHash",
				Handler:    handler,
			},
			{
				MethodName: "deleteLiveHash",
				Handler:    handler,
			},
			{
				MethodName: "getLiveHash",
				Handler:    handler,
			},
			{
				MethodName: "getAccountRecords",
				Handler:    handler,
			},
			{
				MethodName: "cryptoGetBalance",
				Handler:    handler,
			},
			{
				MethodName: "getAccountInfo",
				Handler:    handler,
			},
			{
				MethodName: "getTransactionReceipts",
				Handler:    handler,
			},
			{
				MethodName: "getFastTransactionRecord",
				Handler:    handler,
			},
			{
				MethodName: "getTxRecordByTxID",
				Handler:    handler,
			},
			{
				MethodName: "getStakersByAccountID",
				Handler:    handler,
			},
		},
		Streams:  []grpc.StreamDesc{},
		Metadata: "proto/crypto_service.proto",
	}
}
