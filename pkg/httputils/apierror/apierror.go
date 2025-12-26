package apierror

type APIError struct {
	Code       string
	StatusCode int
	Message    string
}

func (e *APIError) Error() string {
	return e.Message
}

func Error(code, message string, statusCode int) error {
	return &APIError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}
