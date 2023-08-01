package helpers

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

type ListResponse[T any] struct {
	Total    int64 `json:"total"`
	Page     int64 `json:"page"`
	Limit    int64 `json:"limit"`
	LastPage int64 `json:"last_page"`
	From     int64 `json:"from"`
	To       int64 `json:"to"`
	Data     []T   `json:"data"`
}

func NewListResponse[T any](page, limit int, total int64, data []T) *ListResponse[T] {
	response := new(ListResponse[T])
	response.Page = int64(page)
	response.Limit = int64(limit)
	response.From = ((response.Page - 1) * response.Limit) + 1
	response.To = response.From + response.Limit - 1
	response.Total = total
	response.Data = data

	// calculate last page
	lp := float64(total) / float64(limit)
	lastPage := int64(lp)
	if lp > float64(lastPage) {
		lastPage++
	}
	response.LastPage = lastPage

	return response

}
