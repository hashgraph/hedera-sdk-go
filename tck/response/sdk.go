package response

type SetupResponse struct {
	Message string
	Status  string
}

func NewSetupReponse(message string) *SetupResponse {
	response := &SetupResponse{
		Status: "SUCCESS",
	}
	if message != "" {
		response.Message = message
	}
	return response
}
