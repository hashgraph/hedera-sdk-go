package methods

// SPDX-License-Identifier: Apache-2.0
import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/hiero-ledger/hiero-sdk-go/tck/param"
	"github.com/hiero-ledger/hiero-sdk-go/tck/response"
	"github.com/hiero-ledger/hiero-sdk-go/tck/utils"
	hiero "github.com/hiero-ledger/hiero-sdk-go/v2/sdk"
)

type FileService struct {
	sdkService *SDKService
}

func (t *FileService) SetSdkService(service *SDKService) {
	t.sdkService = service
}

// CreateFile jRPC method for createFile
func (t *FileService) CreateFile(_ context.Context, params param.CreateFileParams) (*response.FileResponse, error) {
	transaction := hiero.NewFileCreateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.Keys != nil && len(*params.Keys) > 0 {
		var keys []hiero.Key
		for _, keyStr := range *params.Keys {
			key, err := utils.GetKeyFromString(keyStr)
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
		}
		transaction.SetKeys(keys...)
	}

	if params.Contents != nil {
		transaction.SetContents([]byte(*params.Contents))
	}

	if params.ExpirationTime != nil {
		expirationTime, err := strconv.ParseInt(*params.ExpirationTime, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetExpirationTime(time.Unix(expirationTime, 0))
	}

	if params.Memo != nil {
		transaction.SetMemo(*params.Memo)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.FileResponse{
		FileId: receipt.FileID.String(),
		Status: receipt.Status.String(),
	}, nil
}

// UpdateFile jRPC method for updateFile
func (t *FileService) UpdateFile(_ context.Context, params param.UpdateFileParams) (*response.FileResponse, error) {
	transaction := hiero.NewFileUpdateTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.FileId != nil {
		fileId, err := hiero.FileIDFromString(*params.FileId)
		if err != nil {
			return nil, err
		}
		transaction.SetFileID(fileId)
	}

	if params.Keys != nil {
		var keys []hiero.Key
		for _, keyStr := range *params.Keys {
			key, err := utils.GetKeyFromString(keyStr)
			if err != nil {
				return nil, err
			}
			keys = append(keys, key)
		}
		transaction.SetKeys(keys...)
	}

	if params.Contents != nil {
		transaction.SetContents([]byte(*params.Contents))
	}

	if params.ExpirationTime != nil {
		expirationTime, err := strconv.ParseInt(*params.ExpirationTime, 10, 64)
		if err != nil {
			return nil, err
		}
		transaction.SetExpirationTime(time.Unix(expirationTime, 0))
	}

	if params.Memo != nil {
		transaction.SetFileMemo(*params.Memo)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.FileResponse{
		Status: receipt.Status.String(),
	}, nil
}

// DeleteFile jRPC method for deleteFile
func (t *FileService) DeleteFile(_ context.Context, params param.UpdateFileParams) (*response.FileResponse, error) {
	transaction := hiero.NewFileDeleteTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.FileId != nil {
		fileId, err := hiero.FileIDFromString(*params.FileId)
		if err != nil {
			return nil, err
		}
		transaction.SetFileID(fileId)
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.FileResponse{
		Status: receipt.Status.String(),
	}, nil
}

// AppendFile jRPC method for appendFile
func (t *FileService) AppendFile(_ context.Context, params param.AppendFileParams) (*response.FileResponse, error) {
	transaction := hiero.NewFileAppendTransaction().SetGrpcDeadline(&threeSecondsDuration)

	if params.FileId != nil {
		fileId, err := hiero.FileIDFromString(*params.FileId)
		if err != nil {
			return nil, err
		}
		transaction.SetFileID(fileId)
	}

	if params.Contents != nil {
		transaction.SetContents([]byte(*params.Contents))
	}

	if params.ChunkSize != nil {
		if *params.ChunkSize <= 0 {
			return nil, errors.New("internal error")
		}
		transaction.SetMaxChunkSize(*params.ChunkSize)
	}

	if params.MaxChunks != nil {
		if *params.MaxChunks <= 0 {
			return nil, errors.New("internal error")
		}
		transaction.SetMaxChunks(uint64(*params.MaxChunks))
	}

	if params.CommonTransactionParams != nil {
		err := params.CommonTransactionParams.FillOutTransaction(transaction, t.sdkService.GetClient(params.SessionId))
		if err != nil {
			return nil, err
		}
	}

	txResponse, err := transaction.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}
	receipt, err := txResponse.GetReceipt(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	return &response.FileResponse{
		Status: receipt.Status.String(),
	}, nil
}

// GetFileInfo jRPC method for getFileInfo
func (t *FileService) GetFileInfo(_ context.Context, params param.GetFileInfoParams) (*response.FileInfoResponse, error) {
	query := hiero.NewFileInfoQuery().SetGrpcDeadline(&threeSecondsDuration)

	if params.FileId != nil {
		fileId, err := hiero.FileIDFromString(*params.FileId)
		if err != nil {
			return nil, err
		}
		query.SetFileID(fileId)
	}

	info, err := query.Execute(t.sdkService.GetClient(params.SessionId))
	if err != nil {
		return nil, err
	}

	keys := []string{}
	for _, key := range info.Keys.GetKeys() {
		keys = append(keys, key.String())
	}

	return &response.FileInfoResponse{
		FileId:         info.FileID.String(),
		Size:           strconv.FormatInt(info.Size, 10),
		ExpirationTime: strconv.FormatInt(info.ExpirationTime.Unix(), 10),
		IsDeleted:      info.IsDeleted,
		Keys:           keys,
		Memo:           info.FileMemo,
		LedgerId:       info.LedgerID.String(),
	}, nil
}
