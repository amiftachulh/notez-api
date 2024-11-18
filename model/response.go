package model

type Response struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}

type PaginationResponse struct {
	Total int         `json:"total"`
	Items interface{} `json:"items"`
}
