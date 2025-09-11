package gorouter

import "errors"

var (
	ErrNoHijacker = errors.New("нет доступного экземпляра http.Hijacker")
)

type ErrorResponse struct {
	Error ErrorInfo `json:"error"`
}

type ErrorInfo struct {
	Message string `json:"message"`
	Code    string `json:"code"`
}

func NewErrorResponse(code, message string) ErrorResponse {
	return ErrorResponse{
		Error: ErrorInfo{
			Code:    code,
			Message: message,
		},
	}
}
