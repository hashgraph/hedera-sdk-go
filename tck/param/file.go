package param

// SPDX-License-Identifier: Apache-2.0
type CreateFileParams struct {
	BaseParams
	Keys           *[]string `json:"keys"`
	Contents       *string   `json:"contents"`
	ExpirationTime *string   `json:"expirationTime"`
	Memo           *string   `json:"memo"`
}

type UpdateFileParams struct {
	BaseParams
	FileId         *string   `json:"fileId"`
	Keys           *[]string `json:"keys"`
	Contents       *string   `json:"contents"`
	ExpirationTime *string   `json:"expirationTime"`
	Memo           *string   `json:"memo"`
}

type DeleteFileParams struct {
	BaseParams
	FileId *string `json:"fileId"`
}

type AppendFileParams struct {
	BaseParams
	FileId    *string `json:"fileId"`
	Contents  *string `json:"contents"`
	MaxChunks *int    `json:"maxChunks"`
	ChunkSize *int    `json:"chunkSize"`
}

type GetFileInfoParams struct {
	BaseParams
	FileId *string `json:"fileId"`
}
