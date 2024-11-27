//go:build all || unit
// +build all unit

package hiero

// SPDX-License-Identifier: Apache-2.0

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/mirror"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/stretchr/testify/require"
	protobuf "google.golang.org/protobuf/proto"

	"github.com/hiero-ledger/hiero-sdk-go/v2/proto/services"
	"google.golang.org/grpc"
)

func TestUnitMockQuery(t *testing.T) {
	t.Parallel()
	responses := [][]interface{}{
		{
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
						Header: &services.ResponseHeader{NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK, ResponseType: services.ResponseType_COST_ANSWER, Cost: 0},
						AccountID: &services.AccountID{ShardNum: 0, RealmNum: 0, Account: &services.AccountID_AccountNum{
							AccountNum: 1800,
						}},
						Balance: 2000,
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
		},
	}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	_, err := NewAccountBalanceQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetAccountID(AccountID{Account: 1800}).
		Execute(client)
	require.NoError(t, err)
}

func DisabledTestUnitMockBackoff(t *testing.T) {
	responses := [][]interface{}{{
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
	}, {
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
		&services.TransactionResponse{
			NodeTransactionPrecheckCode: services.ResponseCodeEnum_BUSY,
		},
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	newKey, err := PrivateKeyGenerateEd25519()
	require.NoError(t, err)

	newBalance := NewHbar(2)

	tran := TransactionIDGenerate(AccountID{Account: 3})

	_, err = NewAccountCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}, {Account: 4}}).
		SetKey(newKey).
		SetTransactionID(tran).
		SetInitialBalance(newBalance).
		SetMaxAutomaticTokenAssociations(100).
		Execute(client)
	require.NoError(t, err)
}

func TestUnitMockAddressBookQuery(t *testing.T) {
	t.Parallel()
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
					DomainName:  "hedera.domain.name",
				},
				{
					IpAddressV4: []byte{byte(uint(2)), byte(uint(1)), byte(uint(2)), byte(uint(3))},
					Port:        50123,
					DomainName:  "hedera.domain.name",
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
					DomainName:  "hedera.domain.name2",
				},
				{
					IpAddressV4: []byte{byte(uint(2)), byte(uint(1)), byte(uint(2)), byte(uint(9))},
					Port:        50123,
					DomainName:  "hedera.domain.name2",
				},
			},
			Description: "",
			Stake:       0,
		},
	},
	}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	result, err := NewAddressBookQuery().
		SetFileID(FileID{0, 0, 101, nil}).
		Execute(client)
	require.NoError(t, err)

	require.Equal(t, len(result.NodeAddresses), 2)
	require.Equal(t, result.NodeAddresses[0].AccountID.String(), "0.0.3")
	require.Equal(t, result.NodeAddresses[0].Addresses[0].String(), "hedera.domain.name:50123")
	require.Equal(t, result.NodeAddresses[0].Addresses[1].String(), "hedera.domain.name:50123")
	require.Equal(t, result.NodeAddresses[1].AccountID.String(), "0.0.4")
	require.Equal(t, result.NodeAddresses[1].Addresses[0].String(), "hedera.domain.name2:50123")
	require.Equal(t, result.NodeAddresses[1].Addresses[1].String(), "hedera.domain.name2:50123")
}

func TestUnitMockGenerateTransactionIDsPerExecution(t *testing.T) {
	t.Parallel()
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
		if count < 2 {
			require.False(t, transactionIds[transactionId])
		}
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
	defer server.Close()

	_, err := NewFileCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetContents([]byte("hello")).
		Execute(client)
	require.NoError(t, err)
}

func TestUnitMockSingleTransactionIDForExecutions(t *testing.T) {
	t.Parallel()
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
	defer server.Close()

	_, err := NewFileCreateTransaction().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetTransactionID(tran).
		SetContents([]byte("hello")).
		Execute(client)
	require.NoError(t, err)
}

func TestUnitMockSingleTransactionIDForExecutionsWithTimeout(t *testing.T) {
	t.Parallel()
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
	defer server.Close()

	_, err := NewFileCreateTransaction().
		SetTransactionID(tran).
		SetContents([]byte("hello")).
		Execute(client)
	require.Error(t, err)
}

type MockServers struct {
	servers []*MockServer
}

func (servers *MockServers) Close() {
	for _, server := range servers.servers {
		if server != nil {
			server.Close()
		}
	}
}

func NewMockClientAndServer(allNodeResponses [][]interface{}) (*Client, *MockServers) {
	network := map[string]AccountID{}
	mirrorNetwork := make([]string, len(allNodeResponses))
	servers := make([]*MockServer, len(allNodeResponses))
	ctx, cancel := context.WithCancel(context.Background())

	logger := NewLogger("hedera client mock", LoggerLevelError)
	var defaultLogger Logger = logger

	client := &Client{
		defaultMaxQueryPayment:          NewHbar(1),
		network:                         _NewNetwork(),
		mirrorNetwork:                   _NewMirrorNetwork(),
		autoValidateChecksums:           false,
		maxAttempts:                     nil,
		minBackoff:                      250 * time.Millisecond,
		maxBackoff:                      8 * time.Second,
		defaultRegenerateTransactionIDs: true,
		defaultNetworkUpdatePeriod:      24 * time.Hour,
		networkUpdateContext:            ctx,
		cancelNetworkUpdate:             cancel,
		logger:                          defaultLogger,
	}

	for i, responses := range allNodeResponses {
		responses := responses

		serverReady := make(chan bool)
		nodeAccountID := AccountID{Account: uint64(3 + i)}
		go func() {
			servers[i] = NewMockServer(responses)
			serverReady <- true
		}()

		<-serverReady

		network[servers[i].listener.Addr().String()] = nodeAccountID
		mirrorNetwork[i] = servers[i].listener.Addr().String()
	}

	client.SetNetwork(network)
	client.SetLedgerID(*NewLedgerIDMainnet())
	client.SetMirrorNetwork(mirrorNetwork)

	key, _ := PrivateKeyFromStringEd25519("302e020100300506032b657004220420d45e1557156908c967804615af59a000be88c7aa7058bfcbe0f46b16c28f887d")
	client.SetOperator(AccountID{Account: 1800}, key)
	client.SetMinBackoff(0)
	client.SetMaxBackoff(0)
	client.SetMinNodeReadmitTime(0)
	client.SetMaxNodeReadmitTime(0)
	client.SetNodeMinBackoff(0)
	client.SetNodeMaxBackoff(0)

	return client, &MockServers{servers}
}

func TestUnitMockAccountInfoQuery(t *testing.T) {
	t.Skip("Skipping test as it is currently broken with the addition of generating new payment transactions for queries")
	call := func(request *services.Query) *services.Response {
		require.NotNil(t, request.Query)
		accountInfoQuery := request.Query.(*services.Query_CryptoGetInfo).CryptoGetInfo

		require.Equal(t, accountInfoQuery.AccountID.String(), AccountID{Account: 5}._ToProtobuf().String())

		var payment services.TransactionBody
		require.NotEmpty(t, accountInfoQuery.Header.Payment.BodyBytes)
		err := protobuf.Unmarshal(accountInfoQuery.Header.Payment.BodyBytes, &payment)
		require.NoError(t, err)

		require.NotNil(t, payment.TransactionID)
		require.Equal(t, payment.TransactionID.AccountID.String(), AccountID{Account: 1800}._ToProtobuf().String())
		require.NotNil(t, payment.NodeAccountID)
		require.Equal(t, payment.NodeAccountID.String(), AccountID{Account: 3}._ToProtobuf().String())

		require.Equal(t, payment.Data, &services.TransactionBody_CryptoTransfer{
			CryptoTransfer: &services.CryptoTransferTransactionBody{
				Transfers: &services.TransferList{
					AccountAmounts: []*services.AccountAmount{
						{
							AccountID: AccountID{Account: 3}._ToProtobuf(),
							Amount:    HbarFromTinybar(35).AsTinybar(),
						},
						{
							AccountID: AccountID{Account: 1800}._ToProtobuf(),
							Amount:    -HbarFromTinybar(35).AsTinybar(),
						},
					},
				},
			},
		})

		key, _ := PrivateKeyFromStringEd25519(mockPrivateKey)

		return &services.Response{
			Response: &services.Response_CryptoGetInfo{
				CryptoGetInfo: &services.CryptoGetInfoResponse{
					Header: &services.ResponseHeader{
						NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
						ResponseType:                services.ResponseType_ANSWER_ONLY,
						Cost:                        35,
					},
					AccountInfo: &services.CryptoGetInfoResponse_AccountInfo{
						AccountID:         &services.AccountID{Account: &services.AccountID_AccountNum{5}},
						ContractAccountID: "",
						Deleted:           false,
						ProxyAccountID:    &services.AccountID{Account: &services.AccountID_AccountNum{5}},
						ProxyReceived:     0,
						Key:               key._ToProtoKey(),
						Balance:           0,
					},
				},
			},
		}
	}

	costCall := func(request *services.Query) *services.Response {
		require.NotNil(t, request.Query)
		accountInfoQuery := request.Query.(*services.Query_CryptoGetInfo).CryptoGetInfo

		require.Equal(t, accountInfoQuery.Header.ResponseType, services.ResponseType_COST_ANSWER)

		require.Equal(t, accountInfoQuery.AccountID.String(), AccountID{Account: 5}._ToProtobuf().String())

		var payment services.TransactionBody
		require.NotEmpty(t, accountInfoQuery.Header.Payment.BodyBytes)
		err := protobuf.Unmarshal(accountInfoQuery.Header.Payment.BodyBytes, &payment)
		require.NoError(t, err)

		return &services.Response{
			Response: &services.Response_CryptoGetInfo{
				CryptoGetInfo: &services.CryptoGetInfoResponse{
					Header: &services.ResponseHeader{
						NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
						ResponseType:                services.ResponseType_COST_ANSWER,
						Cost:                        35,
					},
				},
			},
		}
	}

	responses := [][]interface{}{{
		costCall, call,
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	_, err := NewAccountInfoQuery().
		SetNodeAccountIDs([]AccountID{{Account: 3}}).
		SetAccountID(AccountID{Account: 5}).
		Execute(client)
	require.NoError(t, err)
}

func TestUnitMockAccountInfoQueryNoNodeSet(t *testing.T) {
	t.Skip("Skipping test as it is currently broken with the addition of generating new payment transactions for queries")
	call := func(request *services.Query) *services.Response {
		require.NotNil(t, request.Query)
		accountInfoQuery := request.Query.(*services.Query_CryptoGetInfo).CryptoGetInfo

		require.Equal(t, accountInfoQuery.AccountID.String(), AccountID{Account: 5}._ToProtobuf().String())

		var payment services.TransactionBody
		require.NotEmpty(t, accountInfoQuery.Header.Payment.BodyBytes)
		err := protobuf.Unmarshal(accountInfoQuery.Header.Payment.BodyBytes, &payment)
		require.NoError(t, err)

		require.NotNil(t, payment.TransactionID)
		require.Equal(t, payment.TransactionID.AccountID.String(), AccountID{Account: 1800}._ToProtobuf().String())
		require.NotNil(t, payment.NodeAccountID)
		require.Equal(t, payment.NodeAccountID.String(), AccountID{Account: 3}._ToProtobuf().String())

		require.Equal(t, payment.Data, &services.TransactionBody_CryptoTransfer{
			CryptoTransfer: &services.CryptoTransferTransactionBody{
				Transfers: &services.TransferList{
					AccountAmounts: []*services.AccountAmount{
						{
							AccountID: AccountID{Account: 3}._ToProtobuf(),
							Amount:    HbarFromTinybar(35).AsTinybar(),
						},
						{
							AccountID: AccountID{Account: 1800}._ToProtobuf(),
							Amount:    -HbarFromTinybar(35).AsTinybar(),
						},
					},
				},
			},
		})

		key, _ := PrivateKeyFromStringEd25519(mockPrivateKey)

		return &services.Response{
			Response: &services.Response_CryptoGetInfo{
				CryptoGetInfo: &services.CryptoGetInfoResponse{
					Header: &services.ResponseHeader{
						NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
						ResponseType:                services.ResponseType_ANSWER_ONLY,
						Cost:                        35,
					},
					AccountInfo: &services.CryptoGetInfoResponse_AccountInfo{
						AccountID:         &services.AccountID{Account: &services.AccountID_AccountNum{5}},
						ContractAccountID: "",
						Deleted:           false,
						ProxyAccountID:    &services.AccountID{Account: &services.AccountID_AccountNum{5}},
						ProxyReceived:     0,
						Key:               key._ToProtoKey(),
						Balance:           0,
					},
				},
			},
		}
	}

	costCall := func(request *services.Query) *services.Response {
		require.NotNil(t, request.Query)
		accountInfoQuery := request.Query.(*services.Query_CryptoGetInfo).CryptoGetInfo

		require.Equal(t, accountInfoQuery.Header.ResponseType, services.ResponseType_COST_ANSWER)

		require.Equal(t, accountInfoQuery.AccountID.String(), AccountID{Account: 5}._ToProtobuf().String())

		var payment services.TransactionBody
		require.NotEmpty(t, accountInfoQuery.Header.Payment.BodyBytes)
		err := protobuf.Unmarshal(accountInfoQuery.Header.Payment.BodyBytes, &payment)
		require.NoError(t, err)

		return &services.Response{
			Response: &services.Response_CryptoGetInfo{
				CryptoGetInfo: &services.CryptoGetInfoResponse{
					Header: &services.ResponseHeader{
						NodeTransactionPrecheckCode: services.ResponseCodeEnum_OK,
						ResponseType:                services.ResponseType_COST_ANSWER,
						Cost:                        35,
					},
				},
			},
		}
	}

	responses := [][]interface{}{{
		costCall, call,
	}}

	client, server := NewMockClientAndServer(responses)
	defer server.Close()

	_, err := NewAccountInfoQuery().
		SetAccountID(AccountID{Account: 5}).
		Execute(client)
	require.NoError(t, err)
}

func NewMockHandler(responses []interface{}) func(interface{}, context.Context, func(interface{}) error, grpc.UnaryServerInterceptor) (interface{}, error) {
	index := 0
	return func(_srv interface{}, _ctx context.Context, dec func(interface{}) error, _interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
		if index >= len(responses) {
			return nil, status.New(codes.Aborted, "No response found").Err()
		}
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

type MockServer struct {
	listener net.Listener
	server   *grpc.Server
}

func NewMockServer(responses []interface{}) (server *MockServer) {
	var err error
	server = &MockServer{
		server: grpc.NewServer(),
	}
	handler := NewMockHandler(responses)
	streamHandler := NewMockStreamHandler(responses)

	server.server.RegisterService(NewServiceDescription(handler, &services.CryptoService_ServiceDesc), nil)
	server.server.RegisterService(NewServiceDescription(handler, &services.FileService_ServiceDesc), nil)
	server.server.RegisterService(NewServiceDescription(handler, &services.SmartContractService_ServiceDesc), nil)
	server.server.RegisterService(NewServiceDescription(handler, &services.ConsensusService_ServiceDesc), nil)
	server.server.RegisterService(NewServiceDescription(handler, &services.TokenService_ServiceDesc), nil)
	server.server.RegisterService(NewServiceDescription(handler, &services.ScheduleService_ServiceDesc), nil)
	server.server.RegisterService(NewServiceDescription(handler, &services.FreezeService_ServiceDesc), nil)
	server.server.RegisterService(NewServiceDescription(handler, &services.NetworkService_ServiceDesc), nil)
	server.server.RegisterService(NewServiceDescription(handler, &services.AddressBookService_ServiceDesc), nil)
	server.server.RegisterService(NewMirrorServiceDescription(streamHandler, &mirror.NetworkService_ServiceDesc), nil)

	server.listener, err = net.Listen("tcp", "localhost:0")
	if err != nil {
		panic(err)
	}

	go func() {
		if err = server.server.Serve(server.listener); err != nil {
			panic(err)
		}
	}()

	return server
}

func (server *MockServer) Close() {
	if server.server != nil {
		server.server.GracefulStop()
	}
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
