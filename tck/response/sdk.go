package response

type SetupResponse struct {
	Message string
	Status  string
}

func NewSetupReponse(message string) *SetupResponse {
	if message != "" {
		return &SetupResponse{
			Message: message,
			Status:  "SUCCESS",
		}
	}
	return &SetupResponse{
		Status: "SUCCESS",
	}
}
