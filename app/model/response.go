package model

type SuccessResponse struct {
	UserID  int    `json:"user_id"`
	Message string `json:"message"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func NewSuccessResponse(id int, msg string) SuccessResponse {
	return SuccessResponse{
		UserID:  id,
		Message: msg,
	}
}

func NewErrorResponse(msg string) ErrorResponse {
	return ErrorResponse{
		Error: msg,
	}
}
