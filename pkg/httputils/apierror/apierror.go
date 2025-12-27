package apierror

import "net/http"

type APIError struct {
	Code       string
	StatusCode int
	Message    string
}

func (e APIError) Error() string {
	return e.Message
}

func Error(code, message string, statusCode int) error {
	return APIError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

var (
	ErrUnauthorized APIError = APIError{
		Code:       "unauthorized",
		StatusCode: http.StatusUnauthorized,
		Message:    "Unauthorized",
	}
	ErrBadRequest APIError = APIError{
		Code:       "bad_request",
		StatusCode: http.StatusBadRequest,
		Message:    "Bad Request",
	}
	ErrInternalServer APIError = APIError{
		Code:       "internal_server_error",
		StatusCode: http.StatusInternalServerError,
		Message:    "Internal Server Error",
	}
	ErrResourceNotFound APIError = APIError{
		Code:       "resource_not_found",
		StatusCode: http.StatusNotFound,
		Message:    "Resource Not Found",
	}
)
