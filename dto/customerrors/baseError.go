package customerrors

import "fmt"

type baseError struct {
	code         string
	errorMessage string
}

func (e *baseError) Error() string {
	return fmt.Sprintf("\n@Error:: ErrorCode (%s) => ErrorMessage -> (%s)", e.code, e.errorMessage)
}

func GetBaseError(code, errorMessage string) *baseError {
	return &baseError{code, errorMessage}
}

func GetBaseErrorWithDefaultMessage(code string) *baseError {
	return &baseError{code, "Unable to complete your request. Please try again later."}
}
