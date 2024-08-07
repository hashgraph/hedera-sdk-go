package response

type GenerateKeyResponse struct {
	Key         string   `json:"key"`
	PrivateKeys []string `json:"privateKeys"`
}
