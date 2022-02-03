//go:build all || unit
// +build all unit

package hedera

import (
	"context"
	"fmt"
	"math/rand"
	"net"
	"sync"
	"testing"

	"github.com/hashgraph/hedera-protobufs-go/mirror"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/require"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/hashgraph/hedera-protobufs-go/services"
	"google.golang.org/grpc"
)

var wg sync.WaitGroup

func TestUnitMock(t *testing.T) {
	responses := [][]interface{}{{
		status.New(codes.Unavailable, "node is UNAVAILABLE").Err(),
		status.New(codes.Internal, "Received RST_STREAM with code 0").Err(),
	}, {
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
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
						Status: services.ResponseCodeEnum_RECEIPT_NOT_FOUND,
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
	}}

	client, server := NewMockClientAndServer(responses)

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	tran := TransactionIDGenerate(AccountID{Account: 3})

	resp, err := NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
		SetKey(newKey).
		SetTransactionID(tran).
		SetInitialBalance(newBalance).
		SetMaxAutomaticTokenAssociations(100).
		Execute(client)
	require.NoError(t, err)

	receipt, err := resp.GetReceipt(client)
	require.NoError(t, err)
	require.Equal(t, receipt.AccountID, &AccountID{Account: 234})
	server.Close()
}

func TestUnitMockQuery(t *testing.T) {
	responses := [][]interface{}{{
		&services.Response{
			Response: &services.Response_CryptogetAccountBalance{
				CryptogetAccountBalance: &services.CryptoGetAccountBalanceResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY, ResponseType: services.ResponseType_ANSWER_ONLY},
				},
			},
		},
		&services.Response{
			Response: &services.Response_CryptogetAccountBalance{
				CryptogetAccountBalance: &services.CryptoGetAccountBalanceResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY, ResponseType: services.ResponseType_ANSWER_ONLY},
				},
			},
		},
	}, {
		&services.Response{
			Response: &services.Response_CryptogetAccountBalance{
				CryptogetAccountBalance: &services.CryptoGetAccountBalanceResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY, ResponseType: services.ResponseType_ANSWER_ONLY},
				},
			},
		},
		&services.Response{
			Response: &services.Response_CryptogetAccountBalance{
				CryptogetAccountBalance: &services.CryptoGetAccountBalanceResponse{
					Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_ANSWER_ONLY, Cost: 0},
					AccountID: &services.AccountID{ShardNum: 0, RealmNum: 0, Account: &services.AccountID_AccountNum{
						AccountNum: 1800,
					}},
					Balance: 2000,
				},
			},
		},
	}}

	client, server := NewMockClientAndServer(responses)

	_, err := NewAccountBalanceQuery().
		SetAccountID(AccountID{Account: 1800}).
		SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
		Execute(client)
	require.NoError(t, err)

	server.Close()
}

func TestUnitMockAddressBookQuery(t *testing.T) {
	responses := [][]interface{}{{
		&services.NodeAddress{
			RSA_PubKey: "",
			NodeId:     0,
			NodeAccountId: &services.AccountID{
				ShardNum: 0,
				RealmNum: 0,
				Account:  &services.AccountID_AccountNum{AccountNum: 3},
			},
			NodeCertHash: []byte{1},
			ServiceEndpoint: []*services.ServiceEndpoint{
				{
					IpAddressV4: []byte{byte(uint(1)), byte(uint(2)), byte(uint(2)), byte(uint(3))},
					Port:        50123,
				},
				{
					IpAddressV4: []byte{byte(uint(2)), byte(uint(1)), byte(uint(2)), byte(uint(3))},
					Port:        50123,
				},
			},
			Description: "",
			Stake:       0,
		},
		&services.NodeAddress{
			RSA_PubKey: "",
			NodeId:     0,
			NodeAccountId: &services.AccountID{
				ShardNum: 0,
				RealmNum: 0,
				Account:  &services.AccountID_AccountNum{AccountNum: 4},
			},
			NodeCertHash: []byte{1},
			ServiceEndpoint: []*services.ServiceEndpoint{
				{
					IpAddressV4: []byte{byte(uint(1)), byte(uint(2)), byte(uint(2)), byte(uint(9))},
					Port:        50123,
				},
				{
					IpAddressV4: []byte{byte(uint(2)), byte(uint(1)), byte(uint(2)), byte(uint(9))},
					Port:        50123,
				},
			},
			Description: "",
			Stake:       0,
		},
	},
	}

	client, server := NewMockClientAndServer(responses)

	result, err := NewAddressBookQuery().
		SetFileID(FileID{0, 0, 101, nil}).
		Execute(client)
	require.NoError(t, err)

	require.Equal(t, len(result.NodeAddresses), 2)
	require.Equal(t, result.NodeAddresses[0].AccountID.String(), "0.0.3")
	require.Equal(t, result.NodeAddresses[0].Addresses[0].String(), "1.2.2.3:50123")
	require.Equal(t, result.NodeAddresses[0].Addresses[1].String(), "2.1.2.3:50123")
	require.Equal(t, result.NodeAddresses[1].AccountID.String(), "0.0.4")
	require.Equal(t, result.NodeAddresses[1].Addresses[0].String(), "1.2.2.9:50123")
	require.Equal(t, result.NodeAddresses[1].Addresses[1].String(), "2.1.2.9:50123")

	server.Close()
}

func TestUnitMockGenerateTransactionIDsPerExecution(t *testing.T) {
	count := 0
	transactionIds := make(map[string]bool)

	call := func(request *services.Transaction) *services.TransactionResponse {
		var response *services.TransactionResponse

		require.NotEmpty(t, request.SignedTransactionBytes)
		signedTransaction := services.SignedTransaction{}
		_ = protobuf.Unmarshal(request.SignedTransactionBytes, &signedTransaction)

		require.NotEmpty(t, signedTransaction.BodyBytes)
		transactionBody := services.TransactionBody{}
		_ = protobuf.Unmarshal(signedTransaction.BodyBytes, &transactionBody)

		require.NotNil(t, transactionBody.TransactionID)
		transactionId := transactionBody.TransactionID.String()
		require.NotEqual(t, "", transactionId)
		require.False(t, transactionIds[transactionId])
		transactionIds[transactionId] = true

		sigMap := signedTransaction.GetSigMap()
		require.NotNil(t, sigMap)
		require.NotEqual(t, 0, len(sigMap.SigPair))

		for _, sigPair := range sigMap.SigPair {
			verified := false

			switch k := sigPair.Signature.(type) {
			case *services.SignaturePair_Ed25519:
				pbTemp, _ := PublicKeyFromBytesEd25519(sigPair.PubKeyPrefix)
				verified = pbTemp.Verify(signedTransaction.BodyBytes, k.Ed25519)
			case *services.SignaturePair_ECDSASecp256K1:
				pbTemp, _ := PublicKeyFromBytesECDSA(sigPair.PubKeyPrefix)
				verified = pbTemp.Verify(signedTransaction.BodyBytes, k.ECDSASecp256K1)
			}
			require.True(t, verified)
		}

		if count < 2 {
			response = &services.TransactionResponse{
				NodeTransactionPrecheckCode: services.ResponseCodeEnum_TRANSACTION_EXPIRED,
			}
		} else {
			response = &services.TransactionResponse{
				NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
			}
		}

		count += 1

		return response
	}
	responses := [][]interface{}{{
		call, call, call,
	}}

	client, server := NewMockClientAndServer(responses)

	_, err := NewFileCreateTransaction().
		SetContents([]byte("hello")).
		Execute(client)
	require.NoError(t, err)

	server.Close()
}

func TestUnitMockSingleTransactionIDForExecutions(t *testing.T) {
	count := 0
	tran := TransactionIDGenerate(AccountID{Account: 1800})
	transactionIds := make(map[string]bool)
	transactionIds[tran._ToProtobuf().String()] = true

	call := func(request *services.Transaction) *services.TransactionResponse {
		var response *services.TransactionResponse

		require.NotEmpty(t, request.SignedTransactionBytes)
		signedTransaction := services.SignedTransaction{}
		_ = protobuf.Unmarshal(request.SignedTransactionBytes, &signedTransaction)

		require.NotEmpty(t, signedTransaction.BodyBytes)
		transactionBody := services.TransactionBody{}
		_ = protobuf.Unmarshal(signedTransaction.BodyBytes, &transactionBody)

		require.NotNil(t, transactionBody.TransactionID)
		transactionId := transactionBody.TransactionID.String()
		require.NotEqual(t, "", transactionId)
		require.True(t, transactionIds[transactionId])
		transactionIds[transactionId] = true

		sigMap := signedTransaction.GetSigMap()
		require.NotNil(t, sigMap)
		require.NotEqual(t, 0, len(sigMap.SigPair))

		for _, sigPair := range sigMap.SigPair {
			verified := false

			switch k := sigPair.Signature.(type) {
			case *services.SignaturePair_Ed25519:
				pbTemp, _ := PublicKeyFromBytesEd25519(sigPair.PubKeyPrefix)
				verified = pbTemp.Verify(signedTransaction.BodyBytes, k.Ed25519)
			case *services.SignaturePair_ECDSASecp256K1:
				pbTemp, _ := PublicKeyFromBytesECDSA(sigPair.PubKeyPrefix)
				verified = pbTemp.Verify(signedTransaction.BodyBytes, k.ECDSASecp256K1)
			}
			require.True(t, verified)
		}

		if count < 2 {
			response = &services.TransactionResponse{
				NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
			}
		} else {
			response = &services.TransactionResponse{
				NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
			}
		}

		count += 1

		return response
	}
	responses := [][]interface{}{{
		call, call, call,
	}}

	client, server := NewMockClientAndServer(responses)

	_, err := NewFileCreateTransaction().
		SetTransactionID(tran).
		SetContents([]byte("hello")).
		Execute(client)
	require.NoError(t, err)

	server.Close()
}

func TestUnitMockSingleTransactionIDForExecutionsWithTimeout(t *testing.T) {
	count := 0
	tran := TransactionIDGenerate(AccountID{Account: 1800})
	transactionIds := make(map[string]bool)
	transactionIds[tran._ToProtobuf().String()] = true

	call := func(request *services.Transaction) *services.TransactionResponse {
		var response *services.TransactionResponse

		require.NotEmpty(t, request.SignedTransactionBytes)
		signedTransaction := services.SignedTransaction{}
		_ = protobuf.Unmarshal(request.SignedTransactionBytes, &signedTransaction)

		require.NotEmpty(t, signedTransaction.BodyBytes)
		transactionBody := services.TransactionBody{}
		_ = protobuf.Unmarshal(signedTransaction.BodyBytes, &transactionBody)

		require.NotNil(t, transactionBody.TransactionID)
		transactionId := transactionBody.TransactionID.String()
		require.NotEqual(t, "", transactionId)
		require.True(t, transactionIds[transactionId])
		transactionIds[transactionId] = true

		sigMap := signedTransaction.GetSigMap()
		require.NotNil(t, sigMap)
		require.NotEqual(t, 0, len(sigMap.SigPair))

		for _, sigPair := range sigMap.SigPair {
			verified := false

			switch k := sigPair.Signature.(type) {
			case *services.SignaturePair_Ed25519:
				pbTemp, _ := PublicKeyFromBytesEd25519(sigPair.PubKeyPrefix)
				verified = pbTemp.Verify(signedTransaction.BodyBytes, k.Ed25519)
			case *services.SignaturePair_ECDSASecp256K1:
				pbTemp, _ := PublicKeyFromBytesECDSA(sigPair.PubKeyPrefix)
				verified = pbTemp.Verify(signedTransaction.BodyBytes, k.ECDSASecp256K1)
			}
			require.True(t, verified)
		}

		if count < 2 {
			response = &services.TransactionResponse{
				NodeTransactionPrecheckCode: services.ResponseCodeEnum_TRANSACTION_EXPIRED,
			}
		} else {
			response = &services.TransactionResponse{
				NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
			}
		}

		count += 1

		return response
	}
	responses := [][]interface{}{{
		call, call, call,
	}}

	client, server := NewMockClientAndServer(responses)

	_, err := NewFileCreateTransaction().
		SetTransactionID(tran).
		SetContents([]byte("hello")).
		Execute(client)
	require.Error(t, err)

	server.Close()
}

type MockServers struct {
	servers []*grpc.Server
}

func (servers MockServers) Close() {
	for _, server := range servers.servers {
		if server != nil {
			server.Stop()
		}
	}
}

func NewMockClientAndServer(allNodeResponses [][]interface{}) (*Client, MockServers) {
	network := map[string]AccountID{}
	mirrowNetwork := make([]string, 0)
	servers := make([]*grpc.Server, 0)

	for i, responses := range allNodeResponses {
		address := fmt.Sprintf("localhost:%d", 5000+rand.Intn(999))
		nodeAccountID := AccountID{Account: uint64(3 + i)}

		network[address] = nodeAccountID

		responses := responses

		go func() {
			server := NewServer(responses, address)
			servers = append(servers, server)
		}()
	}

	for _, responses := range allNodeResponses {
		address := fmt.Sprintf("localhost:%d", 5000+rand.Intn(999))

		mirrowNetwork = append(mirrowNetwork, address)

		responses := responses

		go func() {
			server := NewServer(responses, address)
			servers = append(servers, server)
		}()
	}

	client := _NewClient(network, mirrowNetwork, "mainnet")

	key, _ := PrivateKeyFromStringEd25519("302e020100300506032b657004220420d45e1557156908c967804615af59a000be88c7aa7058bfcbe0f46b16c28f887d")
	client.SetOperator(AccountID{Account: 1800}, key)

	return client, MockServers{servers}
}

func NewMockHandler(responses []interface{}) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	index := 0
	return func(_srv interface{}, _ctx context.Context, dec func(interface{}) error, _interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		response := responses[index]
		index = index + 1

		switch response := response.(type) {
		case error:
			return nil, response
		case *services.TransactionResponse:
			return response, nil
		case *services.Response:
			return response, nil
		case *services.NodeAddress:
			return response, nil
		case func(request *services.Transaction) *services.TransactionResponse:
			request := new(services.Transaction)
			if err := dec(request); err != nil {
				return nil, err
			}
			return response(request), nil
		case func(request *services.Query) *services.Response:
			request := new(services.Query)
			if err := dec(request); err != nil {
				return nil, err
			}
			return response(request), nil
		case func(request *services.Query) *services.NodeAddress:
			request := new(services.Query)
			if err := dec(request); err != nil {
				return nil, err
			}
			return response(request), nil
		default:
			return response, nil
		}
	}
}

func NewMockStreamHandler(responses []interface{}) func(interface{}, grpc.ServerStream) error {
	return func(_ interface{}, stream grpc.ServerStream) error {
		for _, resp := range responses {
			err := stream.SendMsg(resp)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func NewServer(responses []interface{}, address string) *grpc.Server {
	server := grpc.NewServer()
	handler := NewMockHandler(responses)
	streamHandler := NewMockStreamHandler(responses)

	server.RegisterService(NewServiceDescription(handler, &services.CryptoService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &services.FileService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &services.SmartContractService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &services.ConsensusService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &services.TokenService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &services.ScheduleService_ServiceDesc), nil)
	server.RegisterService(NewServiceDescription(handler, &services.FreezeService_ServiceDesc), nil)
	server.RegisterService(NewMirrorServiceDescription(streamHandler, &mirror.NetworkService_ServiceDesc), nil)

	lis, err := net.Listen("tcp", address)
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

func NewMirrorServiceDescription(handler func(interface{}, grpc.ServerStream) error, service *grpc.ServiceDesc) *grpc.ServiceDesc {
	var streams []grpc.StreamDesc
	for _, stream := range service.Streams {
		streams = append(streams, grpc.StreamDesc{
			StreamName:    stream.StreamName,
			Handler:       handler,
			ServerStreams: stream.ServerStreams,
			ClientStreams: stream.ClientStreams,
		})
	}

	return &grpc.ServiceDesc{
		ServiceName: service.ServiceName,
		HandlerType: service.HandlerType,
		Methods:     []grpc.MethodDesc{},
		Streams:     streams,
		Metadata:    service.Metadata,
	}
}
