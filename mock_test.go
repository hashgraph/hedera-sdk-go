//+build all unit

package hedera

import (
	"context"
	"net"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"google.golang.org/grpc"
)

func TestUnitMock(t *testing.T) {
	responses := []interface{}{
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
		},
		&services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_COST_ANSWER,
					},
				},
			},
		},
		&services.Response{
			Response: &services.Response_TransactionGetReceipt{
				TransactionGetReceipt: &services.TransactionGetReceiptResponse{
					Header: &services.ResponseHeader{
						Cost:         0,
						ResponseType: services.ResponseType_ANSWER_ONLY,
					},
					Receipt: &services.TransactionReceipt{
						Status: services.ResponseCodeEnum_SUCCESS,
						AccountID: &services.AccountID{Account: &services.AccountID_AccountNum{
							AccountNum: 234,
						}},
					},
				},
			},
		},
	}

	client, server := NewMockClientAndServer(responses)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	tran := TransactionIDGenerate(AccountID{Account: 3})

	resp, err := NewAccountCreateTransaction().
		SetKey(newKey).
		SetTransactionID(tran).
		SetInitialBalance(newBalance).
		SetMaxAutomaticTokenAssociations(100).
		Execute(client)
	require.NoError(t, err)

	_, err = resp.GetReceipt(client)
	require.NoError(t, err)
	//assert.Equal(t, balance.Hbars.tinybar, int64(10))

	if server != nil {
		server.GracefulStop()
	}
}

func NewMockClientAndServer(responses []interface{}) (*Client, *grpc.Server) {
	client := ClientForNetwork(map[string]AccountID{
		"0.localhost:50211": {Account: 3},
		"2.localhost:50211": {Account: 4},
		"3.localhost:50211": {Account: 5},
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
		index = index + 1

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

	server.RegisterService(NewServiceDescription(handler, &services.CryptoService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &services.FileService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &services.SmartContractService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &services.ConsensusService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &services.TokenService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &services.ScheduleService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &services.FreezeService_ServiceDesc), nil)

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
