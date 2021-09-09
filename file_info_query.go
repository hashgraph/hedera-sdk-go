package hedera

import (
	"time"

	protobuf "google.golang.org/protobuf/proto"

	"github.com/hashgraph/hedera-sdk-go/v2/proto"
)

type FileInfoQuery struct {
	Query
	fileID *FileID
}

func NewFileInfoQuery() *FileInfoQuery {
	return &FileInfoQuery{
		Query: _NewQuery(true),
	}
}

func (query *FileInfoQuery) SetFileID(fileID FileID) *FileInfoQuery {
	query.fileID = &fileID
	return query
}

func (query *FileInfoQuery) GetFileID() FileID {
	if query.fileID == nil {
		return FileID{}
	}

	return *query.fileID
}

func (query *FileInfoQuery) _ValidateNetworkOnIDs(client *Client) error {
	if client == nil || !client.autoValidateChecksums {
		return nil
	}

	if query.fileID != nil {
		if err := query.fileID.Validate(client); err != nil {
			return err
		}
	}

	return nil
}

func (query *FileInfoQuery) _Build() *proto.Query_FileGetInfo {
	body := &proto.FileGetInfoQuery{
		Header: &proto.QueryHeader{},
	}

	if query.fileID != nil {
		body.FileID = query.fileID._ToProtobuf()
	}

	return &proto.Query_FileGetInfo{
		FileGetInfo: body,
	}
}

func (query *FileInfoQuery) _QueryMakeRequest() _ProtoRequest {
	pb := query._Build()
	_ = query._BuildAllPaymentTransactions()
	if query.isPaymentRequired && len(query.paymentTransactions) > 0 {
		pb.FileGetInfo.Header.Payment = query.paymentTransactions[query.nextPaymentTransactionIndex]
	}
	pb.FileGetInfo.Header.ResponseType = proto.ResponseType_ANSWER_ONLY

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}
}

func (query *FileInfoQuery) _CostQueryMakeRequest(client *Client) (_ProtoRequest, error) {
	pb := query._Build()

	paymentTransaction, err := _QueryMakePaymentTransaction(TransactionIDGenerate(client.GetOperatorAccountID()), AccountID{}, Hbar{})
	if err != nil {
		return _ProtoRequest{}, err
	}

	paymentBytes, err := protobuf.Marshal(paymentTransaction)
	if err != nil {
		return _ProtoRequest{}, err
	}

	pb.FileGetInfo.Header.Payment = &proto.Transaction{
		SignedTransactionBytes: paymentBytes,
	}
	pb.FileGetInfo.Header.ResponseType = proto.ResponseType_COST_ANSWER

	return _ProtoRequest{
		query: &proto.Query{
			Query: pb,
		},
	}, nil
}

func (query *FileInfoQuery) GetCost(client *Client) (Hbar, error) {
	if client == nil || client.operator == nil {
		return Hbar{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.nodeIDs = client.network._GetNodeAccountIDsForExecute()
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return Hbar{}, err
	}

	protoReq, err := query._CostQueryMakeRequest(client)
	if err != nil {
		return Hbar{}, err
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_FileInfoQueryShouldRetry,
		protoReq,
		_CostQueryAdvanceRequest,
		_CostQueryGetNodeAccountID,
		_FileInfoQueryGetMethod,
		_FileInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return Hbar{}, err
	}

	cost := int64(resp.query.GetFileGetInfo().Header.Cost)
	if cost < 25 {
		return HbarFromTinybar(25), nil
	}
	return HbarFromTinybar(cost), nil
}

func _FileInfoQueryShouldRetry(_ _Request, response _Response) _ExecutionState {
	return _QueryShouldRetry(Status(response.query.GetFileGetInfo().Header.NodeTransactionPrecheckCode))
}

func _FileInfoQueryMapStatusError(_ _Request, response _Response) error {
	return ErrHederaPreCheckStatus{
		Status: Status(response.query.GetFileGetInfo().Header.NodeTransactionPrecheckCode),
	}
}

func _FileInfoQueryGetMethod(_ _Request, channel *_Channel) _Method {
	return _Method{
		query: channel._GetFile().GetFileInfo,
	}
}

func (query *FileInfoQuery) Execute(client *Client) (FileInfo, error) {
	if client == nil || client.operator == nil {
		return FileInfo{}, errNoClientProvided
	}

	if len(query.Query.GetNodeAccountIDs()) == 0 {
		query.nodeIDs = client.network._GetNodeAccountIDsForExecute()
	}

	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return FileInfo{}, err
	}

	var cost Hbar
	if query.queryPayment.tinybar != 0 {
		cost = query.queryPayment
	} else {
		if query.maxQueryPayment.tinybar == 0 {
			cost = client.maxQueryPayment
		} else {
			cost = query.maxQueryPayment
		}

		actualCost, err := query.GetCost(client)
		if err != nil {
			return FileInfo{}, err
		}

		if cost.tinybar < actualCost.tinybar {
			return FileInfo{}, ErrMaxQueryPaymentExceeded{
				QueryCost:       actualCost,
				MaxQueryPayment: cost,
				query:           "FileInfoQuery",
			}
		}

		cost = actualCost
	}

	query.actualCost = cost

	if !query.IsFrozen() {
		_, err := query.FreezeWith(client)
		if err != nil {
			return FileInfo{}, err
		}
	}

	transactionID := query.paymentTransactionID

	if !client.GetOperatorAccountID()._IsZero() && client.GetOperatorAccountID()._Equals(*transactionID.AccountID) {
		query.SignWith(
			client.GetOperatorPublicKey(),
			client.operator.signer,
		)
	}

	resp, err := _Execute(
		client,
		_Request{
			query: &query.Query,
		},
		_FileInfoQueryShouldRetry,
		query._QueryMakeRequest(),
		_QueryAdvanceRequest,
		_QueryGetNodeAccountID,
		_FileInfoQueryGetMethod,
		_FileInfoQueryMapStatusError,
		_QueryMapResponse,
	)

	if err != nil {
		return FileInfo{}, err
	}

	info, err := _FileInfoFromProtobuf(resp.query.GetFileGetInfo().FileInfo)
	if err != nil {
		return FileInfo{}, err
	}

	return info, nil
}

// SetMaxQueryPayment sets the maximum payment allowed for this Query.
func (query *FileInfoQuery) SetMaxQueryPayment(maxPayment Hbar) *FileInfoQuery {
	query.Query.SetMaxQueryPayment(maxPayment)
	return query
}

// SetQueryPayment sets the payment amount for this Query.
func (query *FileInfoQuery) SetQueryPayment(paymentAmount Hbar) *FileInfoQuery {
	query.Query.SetQueryPayment(paymentAmount)
	return query
}

func (query *FileInfoQuery) SetNodeAccountIDs(accountID []AccountID) *FileInfoQuery {
	query.Query.SetNodeAccountIDs(accountID)
	return query
}

func (query *FileInfoQuery) GetNodeAccountIDs() []AccountID {
	return query.Query.GetNodeAccountIDs()
}

func (query *FileInfoQuery) SetMaxRetry(count int) *FileInfoQuery {
	query.Query.SetMaxRetry(count)
	return query
}

func (query *FileInfoQuery) SetMaxBackoff(max time.Duration) *FileInfoQuery {
	if max.Nanoseconds() < 0 {
		panic("maxBackoff must be a positive duration")
	} else if max.Nanoseconds() < query.minBackoff.Nanoseconds() {
		panic("maxBackoff must be greater than or equal to minBackoff")
	}
	query.maxBackoff = &max
	return query
}

func (query *FileInfoQuery) GetMaxBackoff() time.Duration {
	if query.maxBackoff != nil {
		return *query.maxBackoff
	}

	return 8 * time.Second
}

func (query *FileInfoQuery) SetMinBackoff(min time.Duration) *FileInfoQuery {
	if min.Nanoseconds() < 0 {
		panic("minBackoff must be a positive duration")
	} else if query.maxBackoff.Nanoseconds() < min.Nanoseconds() {
		panic("minBackoff must be less than or equal to maxBackoff")
	}
	query.minBackoff = &min
	return query
}

func (query *FileInfoQuery) GetMinBackoff() time.Duration {
	if query.minBackoff != nil {
		return *query.minBackoff
	}

	return 250 * time.Millisecond
}

func (query *FileInfoQuery) IsFrozen() bool {
	return query._IsFrozen()
}

// Sign uses the provided privateKey to sign the transaction.
func (query *FileInfoQuery) Sign(
	privateKey PrivateKey,
) *FileInfoQuery {
	return query.SignWith(privateKey.PublicKey(), privateKey.Sign)
}

func (query *FileInfoQuery) SignWithOperator(
	client *Client,
) (*FileInfoQuery, error) {
	// If the transaction is not signed by the _Operator, we need
	// to sign the transaction with the _Operator

	if client == nil {
		return nil, errNoClientProvided
	} else if client.operator == nil {
		return nil, errClientOperatorSigning
	}

	if !query.IsFrozen() {
		_, err := query.FreezeWith(client)
		if err != nil {
			return query, err
		}
	}
	return query.SignWith(client.operator.publicKey, client.operator.signer), nil
}

// SignWith executes the TransactionSigner and adds the resulting signature data to the Transaction's signature map
// with the publicKey as the map key.
func (query *FileInfoQuery) SignWith(
	publicKey PublicKey,
	signer TransactionSigner,
) *FileInfoQuery {
	if !query._KeyAlreadySigned(publicKey) {
		query._SignWith(publicKey, signer)
	}

	return query
}

func (query *FileInfoQuery) Freeze() (*FileInfoQuery, error) {
	return query.FreezeWith(nil)
}

func (query *FileInfoQuery) FreezeWith(client *Client) (*FileInfoQuery, error) {
	if query.IsFrozen() {
		return query, nil
	}
	if query.actualCost.AsTinybar() == 0 {
		if query.queryPayment.tinybar != 0 {
			query.actualCost = query.queryPayment
		} else {
			if query.maxQueryPayment.tinybar == 0 {
				query.actualCost = client.maxQueryPayment
			} else {
				query.actualCost = query.maxQueryPayment
			}

			actualCost, err := query.GetCost(client)
			if err != nil {
				return &FileInfoQuery{}, err
			}

			if query.actualCost.tinybar < actualCost.tinybar {
				return &FileInfoQuery{}, ErrMaxQueryPaymentExceeded{
					QueryCost:       actualCost,
					MaxQueryPayment: query.actualCost,
					query:           "FileInfoQuery",
				}
			}

			query.actualCost = actualCost
		}
	}
	err := query._ValidateNetworkOnIDs(client)
	if err != nil {
		return &FileInfoQuery{}, err
	}
	if err = query._InitPaymentTransactionID(client); err != nil {
		return query, err
	}

	return query, _QueryGeneratePayments(&query.Query, query.actualCost)
}
