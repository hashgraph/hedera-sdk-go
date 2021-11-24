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

	client, server := NewMockClientAndServer(responses)

	balance, err := NewAccountBalanceQuery().SetAccountID(AccountID{Account: 3}).Execute(client)
	assert.NoError(t, err)
	assert.Equal(t, balance.Hbars.tinybar, int64(10))

	if server != nil {
		server.GracefulStop()
	}
}

func NewMockClientAndServer(responses []interface{}) (*Client, *grpc.Server) {
	client := ClientForNetwork(map[string]AccountID{
		"localhost:50211": {Account: 3},
	})

	var server *grpc.Server
	go func() {
		server = NewServer(responses)
	}()

	return client, server
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

func NewServer(responses []interface{}) *grpc.Server {
	server := grpc.NewServer()
	handler := NewMockHandler(responses)

	server.RegisterService(NewServiceDescription(handler, &proto.CryptoService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &proto.FileService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &proto.SmartContractService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &proto.ConsensusService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &proto.TokenService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &proto.ScheduleService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &proto.FreezeService_ServiceDesc), nil)

	lis, err := net.Listen("tcp", "localhost:50211")
	if err != nil {
		panic(err)
	}

	if err = server.Serve(lis); err != nil {
		panic(err)
	}

	return server
}

func NewServiceDescription(handler func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error), service *grpc.ServiceDesc) *grpc.ServiceDesc {
	var methods []grpc.MethodDesc
	for _, desc := range service.Methods {
		methods = append(methods, grpc.MethodDesc{
			MethodName: desc.MethodName,
			Handler:    handler,
		})
	}

	return &grpc.ServiceDesc{
		ServiceName: service.ServiceName,
		HandlerType: service.HandlerType,
		Methods:     methods,
		Streams:     []grpc.StreamDesc{},
		Metadata:    service.Metadata,
	}
}
