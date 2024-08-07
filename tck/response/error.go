package response

import (
	"github.com/creachadair/jrpc2"
	"github.com/hashgraph/hedera-sdk-go/v2"
)

// see https://json-rpc.readthedocs.io/en/latest/exceptions.html
// some of the response codes are handled internally in the jrpc2 lib
const (
	InvalidParamsCode = -32602
	HederaErrorCode   = -32001
)

var (
	InvalidParams = jrpc2.Errorf(InvalidParamsCode, "Invalid params")
	InternalError = jrpc2.Errorf(jrpc2.InternalError, "Internal error")
	HederaError   = jrpc2.Errorf(HederaErrorCode, "Hedera error")
)

type ErrorData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}

func NewHederaReceiptError(err hedera.ErrHederaReceiptStatus) error {
	return HederaError.WithData(&ErrorData{Status: err.Status.String(), Message: err.Error()})
}

func NewHederaPrecheckError(err hedera.ErrHederaPreCheckStatus) error {
	return HederaError.WithData(&ErrorData{Status: err.Status.String(), Message: err.Error()})
}
