package response

// SPDX-License-Identifier: Apache-2.0

type FileResponse struct {
	FileId string `json:"fileId"`
	Status string `json:"status"`
}

type FileInfoResponse struct {
	FileId         string   `json:"fileId"`
	Size           string   `json:"size"`
	ExpirationTime string   `json:"expirationTime"`
	IsDeleted      bool     `json:"isDeleted"`
	Keys           []string `json:"keys"`
	Memo           string   `json:"memo"`
	LedgerId       string   `json:"ledgerId"`
}
