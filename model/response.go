package model

type ErrResp struct {
	Message string      `json:"message"`
	Error   interface{} `json:"error,omitempty"`
}
