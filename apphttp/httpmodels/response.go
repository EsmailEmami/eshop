package httpmodels

type SuccessDBResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func NewSuccessDBResponse(message string, data interface{}) *SuccessDBResponse {
	return &SuccessDBResponse{
		Message: message,
		Data:    data,
	}
}

type SuccessResponse struct {
	Message string `json:"message"`
}

func NewSuccessResponse(message string) *SuccessResponse {
	return &SuccessResponse{
		Message: message,
	}
}
